package handler

import (
	"strconv"

	sq "github.com/Masterminds/squirrel"
	"github.com/brantem/scorecard/constant"
	"github.com/brantem/scorecard/model"
	"github.com/gofiber/fiber/v2"
	"github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
)

// ?limit int
// ?offset int

func (h *Handler) programs(c *fiber.Ctx) error {
	var result struct {
		Nodes []*model.Program `json:"nodes"`
		Error any              `json:"error"`
	}
	result.Nodes = []*model.Program{}

	qb := sq.Select().From("programs")

	var totalCount int
	if err := qb.Column("COUNT(id)").RunWith(h.db).QueryRowContext(c.UserContext()).Scan(&totalCount); err != nil {
		log.Error().Err(err).Msg("program.programs")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}
	c.Set("X-Total-Count", strconv.Itoa(totalCount))

	if c.Method() == fiber.MethodHead {
		return c.SendStatus(fiber.StatusOK)
	}

	if totalCount == 0 {
		return c.Status(fiber.StatusOK).JSON(result)
	}

	if v := c.QueryInt("limit"); v > 0 {
		qb = qb.Limit(uint64(v))
	}

	if v := c.QueryInt("offset"); v > 0 {
		qb = qb.Offset(uint64(v))
	}

	query, args, err := qb.Columns("id", "title").OrderBy("rowid ASC").ToSql()
	if err != nil {
		log.Error().Err(err).Msg("program.programs")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}

	rows, err := h.db.QueryxContext(c.UserContext(), query, args...)
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

	return c.Status(fiber.StatusOK).JSON(result)
}

func (h *Handler) program(c *fiber.Ctx) error {
	var result struct {
		Program *model.Program `json:"program"`
		Error   any            `json:"error"`
	}

	var program model.Program
	err := h.db.QueryRowxContext(c.UserContext(), `
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
		_, err := h.db.ExecContext(c.UserContext(), `UPDATE programs SET title = ? WHERE id = ?`, body.Title, programID)
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
		_, err := h.db.ExecContext(c.UserContext(), `
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

func (h *Handler) deleteProgram(c *fiber.Ctx) error {
	var result struct {
		Success bool `json:"success"`
		Error   any  `json:"error"`
	}

	if _, err := h.db.ExecContext(c.UserContext(), `DELETE FROM programs WHERE id = ?`, c.Params("programId")); err != nil {
		log.Error().Err(err).Msg("program.deleteProgram")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}

	result.Success = true
	return c.Status(fiber.StatusOK).JSON(result)
}
