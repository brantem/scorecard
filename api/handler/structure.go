package handler

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/brantem/scorecard/constant"
	"github.com/brantem/scorecard/model"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (h *Handler) structures(c *fiber.Ctx) error {
	var result struct {
		Nodes []*model.Structure `json:"nodes"`
		Error any                `json:"error"`
	}
	result.Nodes = []*model.Structure{}

	rows, err := h.db.QueryxContext(c.Context(), `SELECT id, parent_id, title, syllabus_id FROM structures`)
	if err != nil {
		log.Error().Err(err).Msg("structure.structures")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}
	defer rows.Close()

	var syllabusIds []int
	for rows.Next() {
		var node model.Structure
		if err := rows.StructScan(&node); err != nil {
			log.Error().Err(err).Msg("structure.structures")
			result.Error = constant.RespInternalServerError
			return c.Status(fiber.StatusInternalServerError).JSON(result)
		}
		result.Nodes = append(result.Nodes, &node)
		if node.SyllabusID != nil {
			syllabusIds = append(syllabusIds, *node.SyllabusID)
		}
	}
	c.Set("X-Total-Count", strconv.Itoa(len(result.Nodes)))

	syllabuses, err := h.getSyllabuses(c.Context(), syllabusIds)
	if err == nil {
		for _, node := range result.Nodes {
			if node.SyllabusID != nil {
				node.Syllabus = syllabuses[*node.SyllabusID]
			}
		}
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func (h *Handler) generateStructures(c *fiber.Ctx) error {
	var result struct {
		Success bool `json:"success"`
		Error   any  `json:"error"`
	}

	rows, err := h.db.QueryContext(c.Context(), `SELECT id, parent_id FROM syllabuses`)
	if err != nil {
		log.Error().Err(err).Msg("structure.generateStructures")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}
	defer rows.Close()

	syllabuses := make(map[int]*int)
	for rows.Next() {
		var syllabusID int
		var parentID *int
		if err := rows.Scan(&syllabusID, &parentID); err != nil {
			log.Error().Err(err).Msg("structure.generateStructures")
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
		log.Error().Err(err).Msg("structure.generateStructures")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}
	defer rows2.Close()

	structures := make(map[int]int, len(syllabuses))
	for rows2.Next() {
		var structureID, syllabusID int
		if err := rows2.Scan(&structureID, &syllabusID); err != nil {
			tx.Rollback()
			log.Error().Err(err).Msg("structure.generateStructures")
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
		log.Error().Err(err).Msg("structure.generateStructures")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}

	tx.Commit()

	result.Success = true
	return c.Status(fiber.StatusOK).JSON(result)
}

func (h *Handler) saveStructure(c *fiber.Ctx) error {
	var result struct {
		Success bool `json:"success"`
		Error   any  `json:"error"`
	}

	var body struct {
		ParentID   *int   `json:"parentID"`
		Title      string `json:"title"`
		SyllabusID *int   `json:"syllabusID"`
	}
	if err := c.BodyParser(&body); err != nil {
		log.Error().Err(err).Msg("structure.saveStructure")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}

	if structureID := c.Params("structureID"); structureID != "" {
		_, err := h.db.ExecContext(c.Context(), `
			UPDATE structures
			SET parent_id = ?, title = ?, syllabus_id,
			WHERE id = ?
		`, body.ParentID, body.Title, body.SyllabusID, structureID)
		if err != nil {
			log.Error().Err(err).Msg("structure.saveStructure")
			result.Error = constant.RespInternalServerError
			return c.Status(fiber.StatusInternalServerError).JSON(result)
		}
	} else {
		_, err := h.db.ExecContext(c.Context(), `
			INSERT INTO structures (parent_id, title, syllabus_id)
			VALUES (?, ?, ?)
		`, body.ParentID, body.Title, body.SyllabusID)
		if err != nil {
			log.Error().Err(err).Msg("structure.saveStructure")
			result.Error = constant.RespInternalServerError
			return c.Status(fiber.StatusInternalServerError).JSON(result)
		}
	}

	result.Success = true
	return c.Status(fiber.StatusOK).JSON(result)
}

func (h *Handler) deleteStructure(c *fiber.Ctx) error {
	var result struct {
		Success bool `json:"success"`
		Error   any  `json:"error"`
	}

	structureID := c.Params("structureId")
	if structureID == "all" {
		_, err := h.db.ExecContext(c.Context(), `DELETE FROM structures`)
		if err != nil {
			log.Error().Err(err).Msg("structure.deleteStructure")
			result.Error = constant.RespInternalServerError
			return c.Status(fiber.StatusInternalServerError).JSON(result)
		}
		result.Success = true
	} else if v, err := strconv.Atoi(structureID); err == nil {
		_, err := h.db.ExecContext(c.Context(), `DELETE FROM structures WHERE id = ?`, v)
		if err != nil {
			log.Error().Err(err).Msg("structure.deleteStructure")
			result.Error = constant.RespInternalServerError
			return c.Status(fiber.StatusInternalServerError).JSON(result)
		}
		result.Success = true
	}

	return c.Status(fiber.StatusOK).JSON(result)
}
