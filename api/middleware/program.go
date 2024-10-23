package middleware

import (
	"github.com/brantem/scorecard/constant"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (m *Middleware) Program(c *fiber.Ctx) error {
	var result struct {
		Error any `json:"error"`
	}

	programID, _ := c.ParamsInt("programId")
	switch {
	case programID < 0:
		result.Error = constant.RespNotFound
		return c.Status(fiber.StatusNotFound).JSON(result)
	case programID == 0:
		return c.Next()
	}

	// In a real production app, this should be cached

	var isExists bool
	err := m.db.QueryRowContext(c.UserContext(), `SELECT EXISTS (
	  SELECT id
	  FROM programs
	  WHERE id = ?
	)`, programID).Scan(&isExists)
	if err != nil {
		log.Error().Err(err).Msg("middleware.Program")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}

	if !isExists {
		result.Error = constant.RespNotFound
		return c.Status(fiber.StatusNotFound).JSON(result)
	}

	return c.Next()
}
