package middleware

import (
	"github.com/brantem/scorecard/constant"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (m *Middleware) SyllabusStructure(c *fiber.Ctx) error {
	var result struct {
		Error any `json:"error"`
	}

	structureID, _ := c.ParamsInt("structureId")
	if structureID < 1 {
		return c.Next()
	}

	// In a real production app, this should be cached

	var isExists bool
	err := m.db.QueryRowContext(c.Context(), `SELECT EXISTS (
	  SELECT id
	  FROM syllabus_structures
	  WHERE id = ?
	    AND program_id = ?
	)`, structureID, c.Params("programId")).Scan(&isExists)
	if err != nil {
		log.Error().Err(err).Msg("middleware.SyllabusStructure")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}

	if !isExists {
		result.Error = constant.RespNotFound
		return c.Status(fiber.StatusNotFound).JSON(result)
	}

	return c.Next()
}

func (m *Middleware) Syllabus(c *fiber.Ctx) error {
	var result struct {
		Error any `json:"error"`
	}

	syllabusID, _ := c.ParamsInt("syllabusId")
	if syllabusID < 1 {
		return c.Next()
	}

	// In a real production app, this should be cached

	var isExists bool
	err := m.db.QueryRowContext(c.Context(), `SELECT EXISTS (
	  SELECT s.id
	  FROM syllabuses s
	  JOIN syllabus_structures ss ON ss.id = s.structure_id
	  WHERE s.id = ?
	    AND ss.program_id = ?
	)`, syllabusID, c.Params("programId")).Scan(&isExists)
	if err != nil {
		log.Error().Err(err).Msg("middleware.Syllabus")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}

	if !isExists {
		result.Error = constant.RespNotFound
		return c.Status(fiber.StatusNotFound).JSON(result)
	}

	return c.Next()
}
