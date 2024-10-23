package handler

import (
	"strconv"

	"github.com/brantem/scorecard/constant"
	"github.com/brantem/scorecard/model"
	"github.com/gofiber/fiber/v2"
	"github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
)

func (h *Handler) programs(c *fiber.Ctx) error {
	var result struct {
		Nodes []*model.Program `json:"nodes"`
		Error any              `json:"error"`
	}
	result.Nodes = []*model.Program{}

	rows, err := h.db.QueryxContext(c.Context(), `SELECT id, title FROM programs`)
	if err != nil {
		log.Error().Err(err).Msg("program.programs")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}
	defer rows.Close()

	for rows.Next() {
		var node model.Program
		if err := rows.StructScan(&node); err != nil {
			log.Error().Err(err).Msg("program.programs")
			result.Error = constant.RespInternalServerError
			return c.Status(fiber.StatusInternalServerError).JSON(result)
		}
		result.Nodes = append(result.Nodes, &node)
	}
	c.Set("X-Total-Count", strconv.Itoa(len(result.Nodes)))

	return c.Status(fiber.StatusOK).JSON(result)
}

func (h *Handler) program(c *fiber.Ctx) error {
	var result struct {
		Program *model.Program `json:"program"`
		Error   any            `json:"error"`
	}

	var program model.Program
	err := h.db.QueryRowxContext(c.Context(), `
		SELECT id, title
		FROM programs
		WHERE id = ?
	`, c.Params("programId")).StructScan(&program)
	if err != nil {
		log.Error().Err(err).Msg("program.program")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}
	result.Program = &program

	return c.Status(fiber.StatusOK).JSON(result)
}

func (h *Handler) saveProgram(c *fiber.Ctx) error {
	var result struct {
		Success bool `json:"success"`
		Error   any  `json:"error"`
	}

	var body struct {
		Title string `json:"title"`
	}
	if err := c.BodyParser(&body); err != nil {
		log.Error().Err(err).Msg("program.saveProgram")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}

	if programID, _ := c.ParamsInt("programId"); programID != 0 {
		_, err := h.db.ExecContext(c.Context(), `UPDATE programs SET title = ? WHERE id = ?`, body.Title, programID)
		if err != nil {
			if err, ok := err.(sqlite3.Error); ok && err.ExtendedCode == sqlite3.ErrConstraintUnique {
				result.Error = fiber.Map{"code": "TITLE_SHOULD_BE_UNIQUE"}
				return c.Status(fiber.StatusConflict).JSON(result)
			}
			log.Error().Err(err).Msg("program.saveProgram")
			result.Error = constant.RespInternalServerError
			return c.Status(fiber.StatusInternalServerError).JSON(result)
		}
	} else {
		_, err := h.db.ExecContext(c.Context(), `
			INSERT INTO programs (title)
			VALUES (?)
		`, body.Title)
		if err != nil {
			if err, ok := err.(sqlite3.Error); ok && err.ExtendedCode == sqlite3.ErrConstraintUnique {
				result.Error = fiber.Map{"code": "TITLE_SHOULD_BE_UNIQUE"}
				return c.Status(fiber.StatusConflict).JSON(result)
			}
			log.Error().Err(err).Msg("program.saveProgram")
			result.Error = constant.RespInternalServerError
			return c.Status(fiber.StatusInternalServerError).JSON(result)
		}
	}

	result.Success = true
	return c.Status(fiber.StatusOK).JSON(result)
}
