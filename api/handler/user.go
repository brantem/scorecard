package handler

import (
	"context"
	"slices"
	"strconv"

	sq "github.com/Masterminds/squirrel"
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

// ?limit int
// ?offset int

func (h *Handler) users(c *fiber.Ctx) error {
	var result struct {
		Nodes []*model.User `json:"nodes"`
		Error any           `json:"error"`
	}
	result.Nodes = []*model.User{}

	qb := sq.Select().From("users").
		Where("program_id = ?", c.Params("programId"))

	var totalCount int
	if err := qb.Column("COUNT(id)").RunWith(h.db).QueryRowContext(c.UserContext()).Scan(&totalCount); err != nil {
		log.Error().Err(err).Msg("user.users")
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

	query, args, err := qb.Columns("id", "name").OrderBy("rowid ASC").ToSql()
	if err != nil {
		log.Error().Err(err).Msg("user.users")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}

	rows, err := h.db.QueryxContext(c.UserContext(), query, args...)
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

	return c.Status(fiber.StatusOK).JSON(result)
}

func (h *Handler) user(c *fiber.Ctx) error {
	var result struct {
		User  *model.User `json:"user"`
		Error any         `json:"error"`
	}

	var user model.User
	err := h.db.QueryRowxContext(c.UserContext(), `
		SELECT id, name
		FROM users
		WHERE id = ?
	`, c.Params("userId")).StructScan(&user)
	if err != nil {
		log.Error().Err(err).Msg("user.user")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}
	result.User = &user

	return c.Status(fiber.StatusOK).JSON(result)
}

func (h *Handler) userScores(c *fiber.Ctx) error {
	type SyllabusWithParentID struct {
		model.BaseSyllabus
		ParentID *int `json:"-" db:"parent_id"`
	}

	type SyllabusWithParents struct {
		SyllabusWithParentID
		Parents []*SyllabusWithParentID `json:"parents"`
	}

	type Node struct {
		Syllabus SyllabusWithParents `json:"syllabus"`
		Score    *float64            `json:"score"`
	}

	var result struct {
		Nodes []*Node `json:"nodes"`
		Error any     `json:"error"`
	}

	rows, err := h.db.QueryxContext(c.UserContext(), `
		SELECT s.id, s.parent_id, s.title, us.score, COALESCE(ss.prev_id, 0) = -1 AS is_assignment
		FROM syllabus_structures ss
		JOIN syllabuses s ON s.structure_id = ss.id
		LEFT JOIN user_scores us ON us.user_id = ? AND us.syllabus_id = s.id
		WHERE ss.program_id = ?
	`, c.Params("userId"), c.Params("programId"))
	if err != nil {
		log.Error().Err(err).Msg("user.userScores")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}

	m := make(map[int]*SyllabusWithParentID)
	for rows.Next() {
		var row struct {
			SyllabusWithParentID
			Score        *float64 `json:"score"`
			IsAssignment bool     `db:"is_assignment"`
		}
		if err := rows.StructScan(&row); err != nil {
			log.Error().Err(err).Msg("user.users")
			result.Error = constant.RespInternalServerError
			return c.Status(fiber.StatusInternalServerError).JSON(result)
		}

		if row.IsAssignment {
			result.Nodes = append(result.Nodes, &Node{
				Syllabus: SyllabusWithParents{row.SyllabusWithParentID, nil},
				Score:    row.Score,
			})
		} else {
			m[row.ID] = &row.SyllabusWithParentID
		}
	}
	c.Set("X-Total-Count", strconv.Itoa(len(result.Nodes)))

	// fill parents
	for _, node := range result.Nodes {
		parentID := node.Syllabus.ParentID
		for parentID != nil {
			parent := m[*parentID]
			node.Syllabus.Parents = append(node.Syllabus.Parents, parent)
			parentID = parent.ParentID
		}

		if len(node.Syllabus.Parents) > 1 {
			slices.Reverse(node.Syllabus.Parents)
		}
	}

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

	if userID, _ := c.ParamsInt("userId"); userID != 0 {
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
		_, err := h.db.ExecContext(c.UserContext(), `
			INSERT INTO users (program_id, name)
			VALUES (?, ?)
		`, c.Params("programId"), body.Name)
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

func (h *Handler) deleteUser(c *fiber.Ctx) error {
	var result struct {
		Success bool `json:"success"`
		Error   any  `json:"error"`
	}

	if _, err := h.db.ExecContext(c.UserContext(), `DELETE FROM users WHERE id = ?`, c.Params("userId")); err != nil {
		log.Error().Err(err).Msg("user.deleteUser")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}

	result.Success = true
	return c.Status(fiber.StatusOK).JSON(result)
}
