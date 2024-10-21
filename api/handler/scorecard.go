package handler

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/brantem/scorecard/constant"
	"github.com/brantem/scorecard/model"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (h *Handler) getScorecardStructureSyllabuses(ctx context.Context, ids []int) (map[int]*model.ScorecardStructureSyllabus, error) {
	if len(ids) == 0 {
		return make(map[int]*model.ScorecardStructureSyllabus), nil
	}

	query, args, err := sqlx.In(`
		SELECT s.id, s.title, COALESCE(ss.prev_id, 0) = -1 AS is_assignment
		FROM syllabuses s
		JOIN syllabus_structures ss ON ss.id = s.structure_id
		WHERE s.id IN (?)
	`, ids)
	if err != nil {
		log.Error().Err(err).Msg("syllabus.getScorecardStructureSyllabuses")
		return nil, constant.ErrInternalServerError
	}

	rows, err := h.db.QueryxContext(ctx, query, args...)
	if err != nil {
		log.Error().Err(err).Msg("syllabus.getScorecardStructureSyllabuses")
		return nil, constant.ErrInternalServerError
	}
	defer rows.Close()

	m := make(map[int]*model.ScorecardStructureSyllabus, len(ids))
	for rows.Next() {
		var syllabus model.ScorecardStructureSyllabus
		if err := rows.StructScan(&syllabus); err != nil {
			log.Error().Err(err).Msg("syllabus.getScorecardStructureSyllabuses")
			return nil, constant.ErrInternalServerError
		}
		m[syllabus.ID] = &syllabus
	}

	return m, nil
}

func (h *Handler) scorecardStructures(c *fiber.Ctx) error {
	var result struct {
		Nodes []*model.ScorecardStructure `json:"nodes"`
		Error any                         `json:"error"`
	}
	result.Nodes = []*model.ScorecardStructure{}

	rows, err := h.db.QueryxContext(c.Context(), `SELECT id, parent_id, title, syllabus_id FROM structures`)
	if err != nil {
		log.Error().Err(err).Msg("scorecard.scorecardStructures")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}
	defer rows.Close()

	var syllabusIds []int
	for rows.Next() {
		var node model.ScorecardStructure
		if err := rows.StructScan(&node); err != nil {
			log.Error().Err(err).Msg("scorecard.scorecardStructures")
			result.Error = constant.RespInternalServerError
			return c.Status(fiber.StatusInternalServerError).JSON(result)
		}
		result.Nodes = append(result.Nodes, &node)
		if node.SyllabusID != nil {
			syllabusIds = append(syllabusIds, *node.SyllabusID)
		}
	}
	c.Set("X-Total-Count", strconv.Itoa(len(result.Nodes)))

	syllabuses, err := h.getScorecardStructureSyllabuses(c.Context(), syllabusIds)
	if err == nil {
		for _, node := range result.Nodes {
			if node.SyllabusID != nil {
				node.Syllabus = syllabuses[*node.SyllabusID]
			}
		}
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func (h *Handler) generateScorecardStructures(c *fiber.Ctx) error {
	var result struct {
		Success bool `json:"success"`
		Error   any  `json:"error"`
	}

	rows, err := h.db.QueryContext(c.Context(), `SELECT id, parent_id FROM syllabuses`)
	if err != nil {
		log.Error().Err(err).Msg("scorecard.generateScorecardStructures")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}
	defer rows.Close()

	syllabuses := make(map[int]*int)
	for rows.Next() {
		var syllabusID int
		var parentID *int
		if err := rows.Scan(&syllabusID, &parentID); err != nil {
			log.Error().Err(err).Msg("scorecard.generateScorecardStructures")
			result.Error = constant.RespInternalServerError
			return c.Status(fiber.StatusInternalServerError).JSON(result)
		}
		syllabuses[syllabusID] = parentID
	}

	tx := h.db.MustBeginTx(c.Context(), nil)

	rows2, err := tx.QueryContext(c.Context(), `
		INSERT INTO structures (title, syllabus_id)
		SELECT title, id
		FROM syllabuses
		RETURNING id, syllabus_id
	`)
	if err != nil {
		tx.Rollback()
		log.Error().Err(err).Msg("scorecard.generateScorecardStructures")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}
	defer rows2.Close()

	structures := make(map[int]int, len(syllabuses))
	for rows2.Next() {
		var structureID, syllabusID int
		if err := rows2.Scan(&structureID, &syllabusID); err != nil {
			tx.Rollback()
			log.Error().Err(err).Msg("scorecard.generateScorecardStructures")
			result.Error = constant.RespInternalServerError
			return c.Status(fiber.StatusInternalServerError).JSON(result)
		}
		structures[syllabusID] = structureID
	}

	values := make([]string, 0, len(structures))
	for syllabusID, structureID := range structures {
		if v := syllabuses[syllabusID]; v != nil {
			values = append(values, fmt.Sprintf("(%d, %d)", structureID, structures[*v]))
		} else {
			values = append(values, fmt.Sprintf("(%d, NULL)", structureID))
		}
	}

	_, err = tx.ExecContext(c.Context(), fmt.Sprintf(`
		WITH t(id, parent_id) AS (
		  VALUES %s
		)
		UPDATE structures
		SET parent_id = t.parent_id
		FROM t
		WHERE structures.id = t.id
	`, strings.Join(values, ", ")))
	if err != nil {
		tx.Rollback()
		log.Error().Err(err).Msg("scorecard.generateScorecardStructures")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}

	tx.Commit()

	result.Success = true
	return c.Status(fiber.StatusOK).JSON(result)
}

func (h *Handler) saveScorecardStructure(c *fiber.Ctx) error {
	var result struct {
		Success bool `json:"success"`
		Error   any  `json:"error"`
	}

	var body struct {
		ParentID   *int    `json:"parentId"`
		Title      *string `json:"title"`
		SyllabusID *int    `json:"syllabusId"`
	}
	if err := c.BodyParser(&body); err != nil {
		log.Error().Err(err).Msg("scorecard.saveScorecardStructure")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}

	structureID := c.Params("structureID")

	if body.Title == nil && body.SyllabusID == nil {
		result.Error = fiber.Map{
			"title":      "REQUIRED",
			"syllabusId": "REQUIRED",
		}
		return c.Status(fiber.StatusBadRequest).JSON(result)
	}

	if structureID != "" {
		if body.Title == nil {
			result.Error = fiber.Map{"title": "REQUIRED"}
			return c.Status(fiber.StatusBadRequest).JSON(result)
		}

		_, err := h.db.ExecContext(c.Context(), `
			UPDATE structures
			SET parent_id = ?, title = ?
			WHERE id = ?
		`, body.ParentID, body.Title, structureID)
		if err != nil {
			log.Error().Err(err).Msg("scorecard.saveScorecardStructure")
			result.Error = constant.RespInternalServerError
			return c.Status(fiber.StatusInternalServerError).JSON(result)
		}
	} else if body.SyllabusID != nil {
		_, err := h.db.ExecContext(c.Context(), `
			INSERT INTO structures (parent_id, title, syllabus_id)
			VALUES (?, (SELECT title from syllabuses WHERE id = ?), ?)
		`, body.ParentID, body.SyllabusID, body.SyllabusID)
		if err != nil {
			log.Error().Err(err).Msg("scorecard.saveScorecardStructure")
			result.Error = constant.RespInternalServerError
			return c.Status(fiber.StatusInternalServerError).JSON(result)
		}
	} else {
		_, err := h.db.ExecContext(c.Context(), `
			INSERT INTO structures (parent_id, title)
			VALUES (?, ?)
		`, body.ParentID, body.Title)
		if err != nil {
			log.Error().Err(err).Msg("scorecard.saveScorecardStructure")
			result.Error = constant.RespInternalServerError
			return c.Status(fiber.StatusInternalServerError).JSON(result)
		}
	}

	result.Success = true
	return c.Status(fiber.StatusOK).JSON(result)
}

func (h *Handler) deleteScorecardStructure(c *fiber.Ctx) error {
	var result struct {
		Success bool `json:"success"`
		Error   any  `json:"error"`
	}

	structureID := c.Params("structureId")
	if structureID == "all" {
		_, err := h.db.ExecContext(c.Context(), `DELETE FROM structures`)
		if err != nil {
			log.Error().Err(err).Msg("scorecard.deleteScorecardStructure")
			result.Error = constant.RespInternalServerError
			return c.Status(fiber.StatusInternalServerError).JSON(result)
		}
		result.Success = true
	} else if v, err := strconv.Atoi(structureID); err == nil {
		_, err := h.db.ExecContext(c.Context(), `DELETE FROM structures WHERE id = ?`, v)
		if err != nil {
			log.Error().Err(err).Msg("scorecard.deleteScorecardStructure")
			result.Error = constant.RespInternalServerError
			return c.Status(fiber.StatusInternalServerError).JSON(result)
		}
		result.Success = true
	}

	return c.Status(fiber.StatusOK).JSON(result)
}
