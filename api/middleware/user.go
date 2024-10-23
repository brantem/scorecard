package middleware

import (
	"github.com/brantem/scorecard/constant"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (m *Middleware) User(c *fiber.Ctx) error {
	var result struct {
		Error any `json:"error"`
	}

	userID, _ := c.ParamsInt("userId")
	switch {
	case userID < 0:
		result.Error = constant.RespNotFound
		return c.Status(fiber.StatusNotFound).JSON(result)
	case userID == 0:
		return c.Next()
	}

	// In a real production app, this should be cached

	var isExists bool
	err := m.db.QueryRowContext(c.UserContext(), `SELECT EXISTS (
	  SELECT id
	  FROM users
	  WHERE id = ?
	    AND program_id = ?
	)`, userID, c.Params("programId")).Scan(&isExists)
	if err != nil {
		log.Error().Err(err).Msg("middleware.User")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}

	if !isExists {
		result.Error = constant.RespNotFound
		return c.Status(fiber.StatusNotFound).JSON(result)
	}

	return c.Next()
}
