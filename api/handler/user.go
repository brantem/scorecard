package handler

import (
	"context"
	"strconv"

	"github.com/brantem/scorecard/constant"
	"github.com/brantem/scorecard/model"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
)

func (h *Handler) getUsers(ctx context.Context, ids []int) (map[int]*model.User, error) {
	if len(ids) == 0 {
		return make(map[int]*model.User), nil
	}

	query, args, err := sqlx.In(`SELECT id, name FROM users WHERE id IN (?)`, ids)
	if err != nil {
		log.Error().Err(err).Msg("user.users")
		return nil, constant.ErrInternalServerError
	}

	rows, err := h.db.QueryxContext(ctx, query, args...)
	if err != nil {
		log.Error().Err(err).Msg("user.users")
		return nil, constant.ErrInternalServerError
	}
	defer rows.Close()

	m := make(map[int]*model.User, len(ids))
	for rows.Next() {
		var user model.User
		if err := rows.StructScan(&user); err != nil {
			log.Error().Err(err).Msg("user.users")
			return nil, constant.ErrInternalServerError
		}
		m[user.ID] = &user
	}

	return m, nil
}

func (h *Handler) users(c *fiber.Ctx) error {
	var result struct {
		Nodes []*model.User `json:"nodes"`
		Error any           `json:"error"`
	}
	result.Nodes = []*model.User{}

	rows, err := h.db.QueryxContext(c.UserContext(), `SELECT id, name FROM users`)
	if err != nil {
		log.Error().Err(err).Msg("user.users")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}
	defer rows.Close()

	for rows.Next() {
		var node model.User
		if err := rows.StructScan(&node); err != nil {
			log.Error().Err(err).Msg("user.users")
			result.Error = constant.RespInternalServerError
			return c.Status(fiber.StatusInternalServerError).JSON(result)
		}
		result.Nodes = append(result.Nodes, &node)
	}
	c.Set("X-Total-Count", strconv.Itoa(len(result.Nodes)))

	return c.Status(fiber.StatusOK).JSON(result)
}

func (h *Handler) saveUser(c *fiber.Ctx) error {
	var result struct {
		Success bool `json:"success"`
		Error   any  `json:"error"`
	}

	var body struct {
		Name string `json:"name"`
	}
	if err := c.BodyParser(&body); err != nil {
		log.Error().Err(err).Msg("user.saveUser")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}

	if userID := c.Params("userId"); userID != "" {
		_, err := h.db.ExecContext(c.UserContext(), `UPDATE users SET name = ? WHERE id = ?`, body.Name, userID)
		if err != nil {
			if err, ok := err.(sqlite3.Error); ok && err.ExtendedCode == sqlite3.ErrConstraintUnique {
				result.Error = fiber.Map{"code": "NAME_SHOULD_BE_UNIQUE"}
				return c.Status(fiber.StatusConflict).JSON(result)
			}
			log.Error().Err(err).Msg("user.saveUser")
			result.Error = constant.RespInternalServerError
			return c.Status(fiber.StatusInternalServerError).JSON(result)
		}
	} else {
		_, err := h.db.ExecContext(c.UserContext(), `INSERT INTO users (name) VALUES (?)`, body.Name)
		if err != nil {
			if err, ok := err.(sqlite3.Error); ok && err.ExtendedCode == sqlite3.ErrConstraintUnique {
				result.Error = fiber.Map{"code": "NAME_SHOULD_BE_UNIQUE"}
				return c.Status(fiber.StatusConflict).JSON(result)
			}
			log.Error().Err(err).Msg("user.saveUser")
			result.Error = constant.RespInternalServerError
			return c.Status(fiber.StatusInternalServerError).JSON(result)
		}
	}

	result.Success = true
	return c.Status(fiber.StatusOK).JSON(result)
}
