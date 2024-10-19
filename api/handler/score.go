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
		SELECT user_id, score
		FROM user_scores
		WHERE syllabus_id = ?
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

	users, err := h.getUsers(c.Context(), userIds)
	if err == nil {
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

	var body struct {
		Score float64 `json:"score"`
	}
	if err := c.BodyParser(&body); err != nil {
		log.Error().Err(err).Msg("score.saveScore")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}

	_, err := h.db.ExecContext(c.UserContext(), `
		INSERT INTO user_scores (user_id, syllabus_id, score)
		VALUES (?, ?, ?)
		ON CONFLICT (user_id, syllabus_id)
		DO UPDATE SET score = EXCLUDED.score
	`, c.Params("userId"), c.Params("syllabusId"), body.Score)
	if err != nil {
		log.Error().Err(err).Msg("score.saveScore")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}

	result.Success = true
	return c.Status(fiber.StatusOK).JSON(result)
}
