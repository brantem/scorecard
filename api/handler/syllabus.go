package handler

import (
	"context"
	"strconv"

	"github.com/brantem/scorecard/constant"
	"github.com/brantem/scorecard/model"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
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
		  SELECT
		    id,
		    prev_id,
		    title
		  FROM syllabus_structures
		  WHERE prev_id IS NULL
		  UNION ALL
		  SELECT
		    s.id,
		    s.prev_id,
		    s.title
		  FROM syllabus_structures s
		  INNER JOIN t ON s.prev_id = t.id
		)
		SELECT id, title FROM t
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
		PrevID *string `json:"prevId"`
		Title  string  `json:"title"`
	}
	if err := c.BodyParser(&body); err != nil {
		log.Error().Err(err).Msg("syllabus.saveSyllabusStructure")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}

	if structureID := c.Params("structureId"); structureID != "" {
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
		_, err := h.db.ExecContext(c.Context(), `
			INSERT INTO syllabus_structures (prev_id, title)
			VALUES (?, ?)
		`, body.PrevID, body.Title)
		if err != nil {
			log.Error().Err(err).Msg("syllabus.saveSyllabusStructure")
			result.Error = constant.RespInternalServerError
			return c.Status(fiber.StatusInternalServerError).JSON(result)
		}
	}

	result.Success = true
	return c.Status(fiber.StatusOK).JSON(result)
}

func (h *Handler) deleteSyllabusStructure(c *fiber.Ctx) error {
	var result struct {
		Success bool `json:"success"`
		Error   any  `json:"error"`
	}

	if _, err := h.db.ExecContext(c.Context(), `DELETE FROM syllabus_structures WHERE id = ?`, c.Params("structureId")); err != nil {
		log.Error().Err(err).Msg("syllabus.deleteSyllabusStructure")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}

	result.Success = true
	return c.Status(fiber.StatusOK).JSON(result)
}

func (h *Handler) getSyllabuses(ctx context.Context, ids []int) (map[int]*model.Syllabus, error) {
	if len(ids) == 0 {
		return make(map[int]*model.Syllabus), nil
	}

	query, args, err := sqlx.In(`SELECT id, title FROM syllabuses WHERE id IN (?)`, ids)
	if err != nil {
		log.Error().Err(err).Msg("syllabus.syllabuses")
		return nil, constant.ErrInternalServerError
	}

	rows, err := h.db.QueryxContext(ctx, query, args...)
	if err != nil {
		log.Error().Err(err).Msg("syllabus.syllabuses")
		return nil, constant.ErrInternalServerError
	}
	defer rows.Close()

	m := make(map[int]*model.Syllabus, len(ids))
	for rows.Next() {
		var syllabus model.Syllabus
		if err := rows.StructScan(&syllabus); err != nil {
			log.Error().Err(err).Msg("syllabus.syllabuses")
			return nil, constant.ErrInternalServerError
		}
		m[syllabus.ID] = &syllabus
	}

	return m, nil
}

func (h *Handler) syllabuses(c *fiber.Ctx) error {
	var result struct {
		Nodes []*model.Syllabus `json:"nodes"`
		Error any               `json:"error"`
	}
	result.Nodes = []*model.Syllabus{}

	rows, err := h.db.QueryxContext(c.Context(), `SELECT id, structure_id, title FROM syllabuses`)
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

func (h *Handler) saveSyllabus(c *fiber.Ctx) error {
	var result struct {
		Success bool `json:"success"`
		Error   any  `json:"error"`
	}

	var body struct {
		StructureID int    `json:"structureId"`
		Title       string `json:"title"`
	}
	if err := c.BodyParser(&body); err != nil {
		log.Error().Err(err).Msg("syllabus.saveSyllabus")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}

	if syllabusID := c.Params("syllabusID"); syllabusID != "" {
		_, err := h.db.ExecContext(c.Context(), `
			UPDATE syllabuses
			SET title = ?
			WHERE id = ?
		`, body.Title, syllabusID)
		if err != nil {
			log.Error().Err(err).Msg("syllabus.saveSyllabus")
			result.Error = constant.RespInternalServerError
			return c.Status(fiber.StatusInternalServerError).JSON(result)
		}
	} else {
		_, err := h.db.ExecContext(c.Context(), `
			INSERT INTO syllabuses (structure_id, title)
			VALUES (?, ?)
		`, body.StructureID, body.Title)
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

	_, err := h.db.ExecContext(c.Context(), `DELETE FROM syllabuses WHERE id = ?`, c.Params("syllabusId"))
	if err != nil {
		log.Error().Err(err).Msg("syllabus.deleteSyllabus")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}

	result.Success = true
	return c.Status(fiber.StatusOK).JSON(result)
}
