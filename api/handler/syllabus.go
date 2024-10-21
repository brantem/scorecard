package handler

import (
	"strconv"

	"github.com/brantem/scorecard/constant"
	"github.com/brantem/scorecard/model"
	"github.com/gofiber/fiber/v2"
	"github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
)

func (h *Handler) syllabusStructures(c *fiber.Ctx) error {
	var result struct {
		Nodes []*model.SyllabusStructure `json:"nodes"`
		Error any                        `json:"error"`
	}
	result.Nodes = []*model.SyllabusStructure{}

	rows, err := h.db.QueryxContext(c.Context(), `
		WITH RECURSIVE t AS (
		  SELECT *
		  FROM syllabus_structures
		  WHERE prev_id IS NULL
		  UNION ALL
		  SELECT s.*
		  FROM syllabus_structures s
		  INNER JOIN t ON s.prev_id = t.id
		)
		SELECT id, prev_id, title FROM t
		UNION
		SELECT id, prev_id, title
		FROM syllabus_structures
		WHERE prev_id = -1
	`)
	if err != nil {
		log.Error().Err(err).Msg("syllabus.syllabusStructures")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}
	defer rows.Close()

	for rows.Next() {
		var node model.SyllabusStructure
		if err := rows.StructScan(&node); err != nil {
			log.Error().Err(err).Msg("syllabus.syllabusStructures")
			result.Error = constant.RespInternalServerError
			return c.Status(fiber.StatusInternalServerError).JSON(result)
		}

		result.Nodes = append(result.Nodes, &node)
	}
	c.Set("X-Total-Count", strconv.Itoa(len(result.Nodes)))

	return c.Status(fiber.StatusOK).JSON(result)
}

func (h *Handler) saveSyllabusStructure(c *fiber.Ctx) error {
	var result struct {
		Success bool `json:"success"`
		Error   any  `json:"error"`
	}

	var body struct {
		PrevID *int   `json:"prevId"`
		Title  string `json:"title"`
	}
	if err := c.BodyParser(&body); err != nil {
		log.Error().Err(err).Msg("syllabus.saveSyllabusStructure")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}

	if structureID, _ := c.ParamsInt("structureId"); structureID != 0 {
		_, err := h.db.ExecContext(c.Context(), `
			UPDATE syllabus_structures
			SET prev_id = ?, title = ?
			WHERE id = ?
		`, body.PrevID, body.Title, structureID)
		if err != nil {
			log.Error().Err(err).Msg("syllabus.saveSyllabusStructure")
			result.Error = constant.RespInternalServerError
			return c.Status(fiber.StatusInternalServerError).JSON(result)
		}
	} else {
		tx := h.db.MustBeginTx(c.Context(), nil)

		// even if this record isn't inserted due to the conflict, the sequence will still increment.
		_, err := tx.ExecContext(c.Context(), `
			INSERT INTO syllabus_structures (prev_id, title)
			VALUES (-1, 'Assignment')
			ON CONFLICT (title) DO NOTHING
		`, body.PrevID, body.Title)
		if err != nil {
			tx.Rollback()
			if err, ok := err.(sqlite3.Error); ok && err.ExtendedCode == sqlite3.ErrConstraintUnique {
				result.Error = fiber.Map{"code": "TITLE_SHOULD_BE_UNIQUE"}
				return c.Status(fiber.StatusConflict).JSON(result)
			}
			log.Error().Err(err).Msg("syllabus.saveSyllabusStructure")
			result.Error = constant.RespInternalServerError
			return c.Status(fiber.StatusInternalServerError).JSON(result)
		}

		_, err = tx.ExecContext(c.Context(), `
			INSERT INTO syllabus_structures (prev_id, title)
			VALUES (?, ?)
		`, body.PrevID, body.Title)
		if err != nil {
			tx.Rollback()
			if err, ok := err.(sqlite3.Error); ok && err.ExtendedCode == sqlite3.ErrConstraintUnique {
				result.Error = fiber.Map{"code": "TITLE_SHOULD_BE_UNIQUE"}
				return c.Status(fiber.StatusConflict).JSON(result)
			}
			log.Error().Err(err).Msg("syllabus.saveSyllabusStructure")
			result.Error = constant.RespInternalServerError
			return c.Status(fiber.StatusInternalServerError).JSON(result)
		}

		tx.Commit()
	}

	result.Success = true
	return c.Status(fiber.StatusOK).JSON(result)
}

func (h *Handler) deleteSyllabusStructure(c *fiber.Ctx) error {
	var result struct {
		Success bool `json:"success"`
		Error   any  `json:"error"`
	}

	structureID := c.Params("structureId")
	if structureID == "all" {
		_, err := h.db.ExecContext(c.Context(), `DELETE FROM syllabus_structures`)
		if err != nil {
			log.Error().Err(err).Msg("syllabus.deleteSyllabusStructure")
			result.Error = constant.RespInternalServerError
			return c.Status(fiber.StatusInternalServerError).JSON(result)
		}
		result.Success = true
	} else if v, err := strconv.Atoi(structureID); err == nil {
		_, err := h.db.ExecContext(c.Context(), `
			WITH t AS (
			  SELECT id FROM syllabus_structures
			)
			DELETE FROM syllabus_structures
			WHERE id IN (
			  SELECT
			    CASE WHEN (SELECT COUNT(*) FROM t) = 2
			      THEN id
			      ELSE ?
			    END
			  FROM t
			  UNION ALL
			  SELECT ? WHERE (SELECT COUNT(*) FROM t) != 2
			)
		`, v, v)
		if err != nil {
			log.Error().Err(err).Msg("syllabus.deleteSyllabusStructure")
			result.Error = constant.RespInternalServerError
			return c.Status(fiber.StatusInternalServerError).JSON(result)
		}
		result.Success = true
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func (h *Handler) syllabuses(c *fiber.Ctx) error {
	var result struct {
		Nodes []*model.Syllabus `json:"nodes"`
		Error any               `json:"error"`
	}
	result.Nodes = []*model.Syllabus{}

	rows, err := h.db.QueryxContext(c.Context(), `SELECT id, parent_id, structure_id, title FROM syllabuses`)
	if err != nil {
		log.Error().Err(err).Msg("syllabus.syllabuses")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}
	defer rows.Close()

	for rows.Next() {
		var node model.Syllabus
		if err := rows.StructScan(&node); err != nil {
			log.Error().Err(err).Msg("syllabus.syllabuses")
			result.Error = constant.RespInternalServerError
			return c.Status(fiber.StatusInternalServerError).JSON(result)
		}
		result.Nodes = append(result.Nodes, &node)
	}
	c.Set("X-Total-Count", strconv.Itoa(len(result.Nodes)))

	return c.Status(fiber.StatusOK).JSON(result)
}

func (h *Handler) syllabus(c *fiber.Ctx) error {
	type Syllabus struct {
		model.BaseSyllabus
		Parents []*model.BaseSyllabus `json:"parents"`
	}

	var result struct {
		Syllabus *Syllabus `json:"syllabus"`
		Error    any       `json:"error"`
	}

	syllabusID, _ := c.ParamsInt("syllabusId")

	rows, err := h.db.QueryxContext(c.Context(), `
		WITH RECURSIVE t AS (
		  SELECT *
		  FROM syllabuses
		  WHERE id = ?
		  UNION ALL
		  SELECT s.*
		  FROM syllabuses s
		  INNER JOIN t ON s.id = t.parent_id
		)
		SELECT id, title FROM t
	`, syllabusID)
	if err != nil {
		log.Error().Err(err).Msg("syllabus.syllabus")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}
	defer rows.Close()

	for rows.Next() {
		var node model.BaseSyllabus
		if err := rows.StructScan(&node); err != nil {
			log.Error().Err(err).Msg("syllabus.syllabus")
			result.Error = constant.RespInternalServerError
			return c.Status(fiber.StatusInternalServerError).JSON(result)
		}

		if node.ID == syllabusID {
			result.Syllabus = &Syllabus{node, nil}
		} else {
			result.Syllabus.Parents = append(result.Syllabus.Parents, &node)
		}
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func (h *Handler) saveSyllabus(c *fiber.Ctx) error {
	var result struct {
		Success bool `json:"success"`
		Error   any  `json:"error"`
	}

	var body struct {
		ParentID    *int   `json:"parentId"`
		StructureID int    `json:"structureId"`
		Title       string `json:"title"`
	}
	if err := c.BodyParser(&body); err != nil {
		log.Error().Err(err).Msg("syllabus.saveSyllabus")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}

	if syllabusID, _ := c.ParamsInt("syllabusID"); syllabusID != 0 {
		_, err := h.db.ExecContext(c.Context(), `
			UPDATE syllabuses
			SET parent_id = ?, title = ?
			WHERE id = ?
		`, body.ParentID, body.Title, syllabusID)
		if err != nil {
			log.Error().Err(err).Msg("syllabus.saveSyllabus")
			result.Error = constant.RespInternalServerError
			return c.Status(fiber.StatusInternalServerError).JSON(result)
		}
	} else {
		_, err := h.db.ExecContext(c.Context(), `
			INSERT INTO syllabuses (parent_id, structure_id, title)
			VALUES (?, ?, ?)
		`, body.ParentID, body.StructureID, body.Title)
		if err != nil {
			log.Error().Err(err).Msg("syllabus.saveSyllabus")
			result.Error = constant.RespInternalServerError
			return c.Status(fiber.StatusInternalServerError).JSON(result)
		}
	}

	result.Success = true
	return c.Status(fiber.StatusOK).JSON(result)
}

func (h *Handler) deleteSyllabus(c *fiber.Ctx) error {
	var result struct {
		Success bool `json:"success"`
		Error   any  `json:"error"`
	}

	syllabusID := c.Params("syllabusId")
	if syllabusID == "all" {
		_, err := h.db.ExecContext(c.Context(), `DELETE FROM syllabuses`)
		if err != nil {
			log.Error().Err(err).Msg("syllabus.deleteSyllabus")
			result.Error = constant.RespInternalServerError
			return c.Status(fiber.StatusInternalServerError).JSON(result)
		}
		result.Success = true
	} else if v, err := strconv.Atoi(syllabusID); err == nil {
		_, err := h.db.ExecContext(c.Context(), `DELETE FROM syllabuses WHERE id = ?`, v)
		if err != nil {
			log.Error().Err(err).Msg("syllabus.deleteSyllabus")
			result.Error = constant.RespInternalServerError
			return c.Status(fiber.StatusInternalServerError).JSON(result)
		}
		result.Success = true
	}

	return c.Status(fiber.StatusOK).JSON(result)
}
