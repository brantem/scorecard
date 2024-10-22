package handler

import (
	"strconv"

	"github.com/brantem/scorecard/constant"
	"github.com/brantem/scorecard/model"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func (h *Handler) scores(c *fiber.Ctx) error {
	var result struct {
		Nodes []*model.Score `json:"nodes"`
		Error any            `json:"error"`
	}
	result.Nodes = []*model.Score{}

	rows, err := h.db.QueryxContext(c.Context(), `
		SELECT u.id AS user_id, us.score
		FROM users u
		LEFT JOIN user_scores us ON us.user_id = u.id AND us.syllabus_id = ?
	`, c.Params("syllabusId"))
	if err != nil {
		log.Error().Err(err).Msg("score.scores")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}
	defer rows.Close()

	var userIds []int
	for rows.Next() {
		var node model.Score
		if err := rows.StructScan(&node); err != nil {
			log.Error().Err(err).Msg("score.scores")
			result.Error = constant.RespInternalServerError
			return c.Status(fiber.StatusInternalServerError).JSON(result)
		}
		result.Nodes = append(result.Nodes, &node)
		userIds = append(userIds, node.UserID)
	}
	c.Set("X-Total-Count", strconv.Itoa(len(result.Nodes)))

	if users, err := h.getUsers(c.Context(), userIds); err == nil {
		for _, node := range result.Nodes {
			node.User = users[node.UserID]
		}
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func (h *Handler) saveScore(c *fiber.Ctx) error {
	var result struct {
		Success bool `json:"success"`
		Error   any  `json:"error"`
	}

	userID, _ := c.ParamsInt("userId")

	var body struct {
		Score float64 `json:"score"`
	}
	if err := c.BodyParser(&body); err != nil {
		log.Error().Err(err).Msg("score.saveScore")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}

	tx := h.db.MustBeginTx(c.Context(), nil)

	_, err := tx.ExecContext(c.UserContext(), `
		INSERT INTO user_scores (user_id, syllabus_id, score)
		VALUES (?, ?, ?)
		ON CONFLICT (user_id, syllabus_id)
		DO UPDATE SET score = EXCLUDED.score
	`, userID, c.Params("syllabusId"), body.Score)
	if err != nil {
		tx.Rollback()
		log.Error().Err(err).Msg("score.saveScore")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}

	_, err = tx.ExecContext(c.Context(), `UPDATE scorecards SET is_outdated = TRUE WHERE user_id = ?`, userID)
	if err != nil {
		tx.Rollback()
		log.Error().Err(err).Msg("score.saveScore")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}

	tx.Commit()

	result.Success = true
	return c.Status(fiber.StatusOK).JSON(result)
}
