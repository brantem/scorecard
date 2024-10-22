package handler

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/brantem/scorecard/constant"
	"github.com/brantem/scorecard/model"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

func (h *Handler) getScorecardStructureSyllabuses(ctx context.Context, ids []int) (map[int]*model.ScorecardStructureSyllabus, error) {
	if len(ids) == 0 {
		return make(map[int]*model.ScorecardStructureSyllabus), nil
	}

	query, args, err := sqlx.In(`
		SELECT s.id, s.title, COALESCE(ss.prev_id, 0) = -1 AS is_assignment
		FROM syllabuses s
		JOIN syllabus_structures ss ON ss.id = s.structure_id
		WHERE s.id IN (?)
	`, ids)
	if err != nil {
		log.Error().Err(err).Msg("syllabus.getScorecardStructureSyllabuses")
		return nil, constant.ErrInternalServerError
	}

	rows, err := h.db.QueryxContext(ctx, query, args...)
	if err != nil {
		log.Error().Err(err).Msg("syllabus.getScorecardStructureSyllabuses")
		return nil, constant.ErrInternalServerError
	}
	defer rows.Close()

	m := make(map[int]*model.ScorecardStructureSyllabus, len(ids))
	for rows.Next() {
		var syllabus model.ScorecardStructureSyllabus
		if err := rows.StructScan(&syllabus); err != nil {
			log.Error().Err(err).Msg("syllabus.getScorecardStructureSyllabuses")
			return nil, constant.ErrInternalServerError
		}
		m[syllabus.ID] = &syllabus
	}

	return m, nil
}

// ?depth int

func (h *Handler) scorecardStructures(c *fiber.Ctx) error {
	var result struct {
		Nodes []*model.ScorecardStructure `json:"nodes"`
		Error any                         `json:"error"`
	}
	result.Nodes = []*model.ScorecardStructure{}

	depth := c.QueryInt("depth")

	rows, err := h.db.QueryxContext(c.Context(), `
		WITH RECURSIVE t AS (
		  SELECT *, 1 AS depth
		  FROM scorecard_structures
		  WHERE parent_id IS NULL
		  UNION ALL
		  SELECT ss.*, t.depth + 1
		  FROM scorecard_structures ss
		  JOIN t ON ss.parent_id = t.id
		  WHERE (? = 0 OR t.depth < ?)
		)
		SELECT id, parent_id, title, syllabus_id FROM t
	`, depth, depth)
	if err != nil {
		log.Error().Err(err).Msg("scorecard.scorecardStructures")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}
	defer rows.Close()

	var syllabusIds []int
	for rows.Next() {
		var node model.ScorecardStructure
		if err := rows.StructScan(&node); err != nil {
			log.Error().Err(err).Msg("scorecard.scorecardStructures")
			result.Error = constant.RespInternalServerError
			return c.Status(fiber.StatusInternalServerError).JSON(result)
		}
		result.Nodes = append(result.Nodes, &node)
		if node.SyllabusID != nil {
			syllabusIds = append(syllabusIds, *node.SyllabusID)
		}
	}
	c.Set("X-Total-Count", strconv.Itoa(len(result.Nodes)))

	if syllabuses, err := h.getScorecardStructureSyllabuses(c.Context(), syllabusIds); err == nil {
		for _, node := range result.Nodes {
			if node.SyllabusID != nil {
				node.Syllabus = syllabuses[*node.SyllabusID]
			}
		}
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func (h *Handler) copySyllabusesIntoStructures(c *fiber.Ctx) error {
	var result struct {
		Success bool `json:"success"`
		Error   any  `json:"error"`
	}

	var body struct {
		ParentID *int `json:"parentId"`
	}
	if err := c.BodyParser(&body); err != nil {
		log.Error().Err(err).Msg("scorecard.copySyllabusesIntoStructures")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}

	targetID, _ := c.ParamsInt("syllabusId")

	syllabuses := make(map[int]*int)
	var syllabusIds []int

	rows, err := h.db.QueryContext(c.Context(), `
		WITH RECURSIVE t AS (
		  SELECT *
		  FROM syllabuses
		  WHERE (? = 0 OR id = ?)
		  UNION ALL
		  SELECT s.*
		  FROM syllabuses s
		  INNER JOIN t ON s.parent_id = t.id
		)
		SELECT id, parent_id FROM t
	`, targetID, targetID)
	if err != nil {
		log.Error().Err(err).Msg("scorecard.copySyllabusesIntoStructures")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}
	defer rows.Close()

	for rows.Next() {
		var syllabusID int
		var parentID *int
		if err := rows.Scan(&syllabusID, &parentID); err != nil {
			log.Error().Err(err).Msg("scorecard.copySyllabusesIntoStructures")
			result.Error = constant.RespInternalServerError
			return c.Status(fiber.StatusInternalServerError).JSON(result)
		}
		syllabuses[syllabusID] = parentID
		syllabusIds = append(syllabusIds, syllabusID)
	}

	if len(syllabuses) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(result)
	}

	query, args, err := sqlx.In(`SELECT title, id FROM syllabuses WHERE id IN (?)`, syllabusIds)
	if err != nil {
		log.Error().Err(err).Msg("scorecard.copySyllabusesIntoStructures")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}

	tx := h.db.MustBeginTx(c.Context(), nil)

	rows, err = tx.QueryContext(c.Context(), fmt.Sprintf(`INSERT INTO scorecard_structures (title, syllabus_id) %s RETURNING id, syllabus_id`, query), args...)
	if err != nil {
		tx.Rollback()
		log.Error().Err(err).Msg("scorecard.copySyllabusesIntoStructures")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}
	defer rows.Close()

	structures := make(map[int]int, len(syllabuses))
	for rows.Next() {
		var structureID, syllabusID int
		if err := rows.Scan(&structureID, &syllabusID); err != nil {
			tx.Rollback()
			log.Error().Err(err).Msg("scorecard.copySyllabusesIntoStructures")
			result.Error = constant.RespInternalServerError
			return c.Status(fiber.StatusInternalServerError).JSON(result)
		}
		structures[syllabusID] = structureID
	}

	values := make([]string, 0, len(structures))
	for syllabusID, structureID := range structures {
		var newID *int
		if syllabusID == targetID {
			newID = body.ParentID
		} else if v := syllabuses[syllabusID]; v != nil {
			_newID := structures[*v]
			newID = &_newID
		}

		if newID != nil {
			values = append(values, fmt.Sprintf("(%d, %d)", structureID, *newID))
		} else {
			values = append(values, fmt.Sprintf("(%d, NULL)", structureID))
		}
	}

	_, err = tx.ExecContext(c.Context(), fmt.Sprintf(`
		WITH t(id, parent_id) AS (
		  VALUES %s
		)
		UPDATE scorecard_structures
		SET parent_id = t.parent_id
		FROM t
		WHERE scorecard_structures.id = t.id
	`, strings.Join(values, ", ")))
	if err != nil {
		tx.Rollback()
		log.Error().Err(err).Msg("scorecard.copySyllabusesIntoStructures")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}

	_, err = tx.ExecContext(c.Context(), `UPDATE scorecards SET is_outdated = TRUE`)
	if err != nil {
		tx.Rollback()
		log.Error().Err(err).Msg("scorecard.copySyllabusesIntoStructures")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}

	tx.Commit()

	result.Success = true
	return c.Status(fiber.StatusOK).JSON(result)
}

func (h *Handler) saveScorecardStructure(c *fiber.Ctx) error {
	var result struct {
		Success bool `json:"success"`
		Error   any  `json:"error"`
	}

	var body struct {
		ParentID *int   `json:"parentId"`
		Title    string `json:"title"`
	}
	if err := c.BodyParser(&body); err != nil {
		log.Error().Err(err).Msg("scorecard.saveScorecardStructure")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}

	structureID, _ := c.ParamsInt("structureId")

	if structureID == 0 {
		_, err := h.db.ExecContext(c.Context(), `
		INSERT INTO scorecard_structures (parent_id, title)
		VALUES (?, ?)
	`, body.ParentID, body.Title)
		if err != nil {
			log.Error().Err(err).Msg("scorecard.saveScorecardStructure")
			result.Error = constant.RespInternalServerError
			return c.Status(fiber.StatusInternalServerError).JSON(result)
		}
	} else if structureID > 0 {
		_, err := h.db.ExecContext(c.Context(), `
			UPDATE scorecard_structures
			SET parent_id = ?, title = ?
			WHERE id = ?
		`, body.ParentID, body.Title, structureID)
		if err != nil {
			log.Error().Err(err).Msg("scorecard.saveScorecardStructure")
			result.Error = constant.RespInternalServerError
			return c.Status(fiber.StatusInternalServerError).JSON(result)
		}
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(result)
	}

	result.Success = true
	return c.Status(fiber.StatusOK).JSON(result)
}

func (h *Handler) deleteScorecardStructure(c *fiber.Ctx) error {
	var result struct {
		Success bool `json:"success"`
		Error   any  `json:"error"`
	}

	tx := h.db.MustBeginTx(c.Context(), nil)

	structureID, _ := c.ParamsInt("structureId")
	_, err := tx.ExecContext(c.Context(), `
		DELETE FROM scorecard_structures
		WHERE (? = 0 OR id = ?)
	`, structureID, structureID)
	if err != nil {
		tx.Rollback()
		log.Error().Err(err).Msg("scorecard.deleteScorecardStructure")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}

	_, err = tx.ExecContext(c.Context(), `UPDATE scorecards SET is_outdated = TRUE`)
	if err != nil {
		tx.Rollback()
		log.Error().Err(err).Msg("scorecard.deleteScorecardStructure")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}

	tx.Commit()

	result.Success = true
	return c.Status(fiber.StatusOK).JSON(result)
}

func (h *Handler) generateScorecards(c *fiber.Ctx) error {
	var result struct {
		Success bool `json:"success"`
		Error   any  `json:"error"`
	}

	if scorecardID, _ := c.ParamsInt("scorecardId"); scorecardID > 0 {
		var userID int
		err := h.db.QueryRowContext(c.Context(), `
			SELECT user_id
			FROM scorecards s
			WHERE id = ?
		`, scorecardID).Scan(&userID)
		if err != nil {
			log.Error().Err(err).Msg("scorecard.generateScorecards")
			result.Error = constant.RespInternalServerError
			return c.Status(fiber.StatusInternalServerError).JSON(result)
		}
		h.generator.Add(c.Context(), userID, scorecardID)
		result.Success = true
	} else {
		rows, err := h.db.QueryContext(c.Context(), `
			SELECT DISTINCT us.user_id, COALESCE(s.id, 0)
			FROM user_scores us
			LEFT JOIN scorecards s ON s.user_id = us.user_id
		`)
		if err != nil {
			log.Error().Err(err).Msg("scorecard.generateScorecards")
			result.Error = constant.RespInternalServerError
			return c.Status(fiber.StatusInternalServerError).JSON(result)
		}
		defer rows.Close()

		for rows.Next() {
			result.Success = true
			var userID, scorecardID int
			if err := rows.Scan(&userID, &scorecardID); err != nil {
				log.Error().Err(err).Msg("scorecard.generateScorecards")
				result.Error = constant.RespInternalServerError
				return c.Status(fiber.StatusInternalServerError).JSON(result)
			}
			h.generator.Add(c.Context(), userID, scorecardID)
		}
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func (h *Handler) scorecards(c *fiber.Ctx) error {
	var result struct {
		Nodes []*model.Scorecard `json:"nodes"`
		Error any                `json:"error"`
	}
	result.Nodes = []*model.Scorecard{}

	rows, err := h.db.QueryxContext(c.Context(), `
		SELECT id, user_id, score, is_outdated, generated_at
		FROM scorecards
	`)
	if err != nil {
		log.Error().Err(err).Msg("scorecard.scorecards")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}
	defer rows.Close()

	var userIds []int
	for rows.Next() {
		var node model.Scorecard
		if err := rows.StructScan(&node); err != nil {
			log.Error().Err(err).Msg("scorecard.scorecards")
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

func (h *Handler) scorecard(c *fiber.Ctx) error {
	var result struct {
		Scorecard *model.Scorecard `json:"scorecard"`
		Error     any              `json:"error"`
	}

	scorecardID, _ := c.ParamsInt("scorecardId")

	scorecard := model.Scorecard{Items: []*model.ScorecardItem{}}
	err := h.db.QueryRowxContext(c.Context(), `
		SELECT id, user_id, score, is_outdated, generated_at
		FROM scorecards
		WHERE id = ?
	`, scorecardID).StructScan(&scorecard)
	if err != nil {
		log.Error().Err(err).Msg("scorecard.scorecard")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}
	result.Scorecard = &scorecard

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		userID := result.Scorecard.UserID
		users, _ := h.getUsers(c.Context(), []int{userID})
		result.Scorecard.User = users[userID]
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		rows, err := h.db.QueryxContext(c.Context(), `
			SELECT id, structure_id, score
			FROM scorecard_items
			WHERE scorecard_id = ?
		`, result.Scorecard.ID)
		if err != nil {
			log.Error().Err(err).Msg("scorecard.scorecard")
			return
		}
		defer rows.Close()

		for rows.Next() {
			var node model.ScorecardItem
			if err := rows.StructScan(&node); err != nil {
				continue
			}
			result.Scorecard.Items = append(result.Scorecard.Items, &node)
		}
	}()

	wg.Wait()

	return c.Status(fiber.StatusOK).JSON(result)
}
