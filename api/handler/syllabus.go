package handler

import (
	"slices"
	"strconv"

	sq "github.com/Masterminds/squirrel"
	"github.com/brantem/scorecard/constant"
	"github.com/brantem/scorecard/model"
	"github.com/gofiber/fiber/v2"
	"github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
)

func (h *Handler) syllabusStructures(c *fiber.Ctx) error {
	var result struct {
		Nodes []*model.SyllabusStructure `json:"nodes"`
		Error any                        `json:"error"`
	}
	result.Nodes = []*model.SyllabusStructure{}

	programID, _ := c.ParamsInt("programId")

	rows, err := h.db.QueryxContext(c.UserContext(), `
		WITH RECURSIVE t AS (
		  SELECT *
		  FROM syllabus_structures
		  WHERE program_id = ?
		    AND prev_id IS NULL
		  UNION ALL
		  SELECT s.*
		  FROM syllabus_structures s
		  INNER JOIN t ON s.prev_id = t.id
		)
		SELECT id, prev_id, title FROM t
		UNION
		SELECT id, prev_id, title
		FROM syllabus_structures
		WHERE program_id = ?
		  AND prev_id = -1
	`, programID, programID)
	if err != nil {
		log.Error().Err(err).Msg("syllabus.syllabusStructures")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}
	defer rows.Close()

	for rows.Next() {
		var node model.SyllabusStructure
		if err := rows.StructScan(&node); err != nil {
			log.Error().Err(err).Msg("syllabus.syllabusStructures")
			result.Error = constant.RespInternalServerError
			return c.Status(fiber.StatusInternalServerError).JSON(result)
		}

		result.Nodes = append(result.Nodes, &node)
	}
	c.Set("X-Total-Count", strconv.Itoa(len(result.Nodes)))

	return c.Status(fiber.StatusOK).JSON(result)
}

func (h *Handler) saveSyllabusStructure(c *fiber.Ctx) error {
	var result struct {
		Success bool `json:"success"`
		Error   any  `json:"error"`
	}

	programID, _ := c.ParamsInt("programId")

	var body struct {
		PrevID *int   `json:"prevId"`
		Title  string `json:"title"`
	}
	if err := c.BodyParser(&body); err != nil {
		log.Error().Err(err).Msg("syllabus.saveSyllabusStructure")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}

	if structureID, _ := c.ParamsInt("structureId"); structureID != 0 {
		// TODO: Update prev_id
		_, err := h.db.ExecContext(c.UserContext(), `
			UPDATE syllabus_structures
			SET title = ?
			WHERE id = ?
		`, body.Title, structureID)
		if err != nil {
			if err, ok := err.(sqlite3.Error); ok && err.ExtendedCode == sqlite3.ErrConstraintUnique {
				result.Error = fiber.Map{"code": "TITLE_SHOULD_BE_UNIQUE"}
				return c.Status(fiber.StatusConflict).JSON(result)
			}
			log.Error().Err(err).Msg("syllabus.saveSyllabusStructure")
			result.Error = constant.RespInternalServerError
			return c.Status(fiber.StatusInternalServerError).JSON(result)
		}
	} else {
		tx := h.db.MustBeginTx(c.UserContext(), nil)

		if body.PrevID == nil {
			_, err := tx.ExecContext(c.UserContext(), `
				INSERT INTO syllabus_structures (program_id, prev_id, title)
				VALUES (?, -1, 'Assignment')
				ON CONFLICT (program_id, title) DO NOTHING
			`, programID)
			if err != nil {
				tx.Rollback()
				log.Error().Err(err).Msg("syllabus.saveSyllabusStructure")
				result.Error = constant.RespInternalServerError
				return c.Status(fiber.StatusInternalServerError).JSON(result)
			}
		}

		_, err := tx.ExecContext(c.UserContext(), `
			INSERT INTO syllabus_structures (program_id, prev_id, title)
			VALUES (?, ?, ?)
		`, programID, body.PrevID, body.Title)
		if err != nil {
			tx.Rollback()
			if err, ok := err.(sqlite3.Error); ok && err.ExtendedCode == sqlite3.ErrConstraintUnique {
				result.Error = fiber.Map{"code": "TITLE_SHOULD_BE_UNIQUE"}
				return c.Status(fiber.StatusConflict).JSON(result)
			}
			log.Error().Err(err).Msg("syllabus.saveSyllabusStructure")
			result.Error = constant.RespInternalServerError
			return c.Status(fiber.StatusInternalServerError).JSON(result)
		}

		tx.Commit()
	}

	result.Success = true
	return c.Status(fiber.StatusOK).JSON(result)
}

func (h *Handler) deleteSyllabusStructure(c *fiber.Ctx) error {
	var result struct {
		Success bool `json:"success"`
		Error   any  `json:"error"`
	}

	programID, _ := c.ParamsInt("programId")
	structureID, _ := c.ParamsInt("structureId", -1)

	tx := h.db.MustBeginTx(c.UserContext(), nil)

	// If structureID == 0, all syllabus structures will be deleted
	_, err := tx.ExecContext(c.UserContext(), `
		DELETE FROM syllabus_structures
		WHERE (? = 0 OR id = ?)
	`, structureID, structureID)
	if err != nil {
		tx.Rollback()
		log.Error().Err(err).Msg("syllabus.deleteSyllabusStructure")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}

	if structureID == 0 {
		_, err = tx.ExecContext(c.UserContext(), `DELETE FROM scorecard_structures WHERE program_id = ?`, programID)
		if err != nil {
			tx.Rollback()
			log.Error().Err(err).Msg("syllabus.deleteSyllabusStructure")
			result.Error = constant.RespInternalServerError
			return c.Status(fiber.StatusInternalServerError).JSON(result)
		}

		_, err = tx.ExecContext(c.UserContext(), `DELETE FROM scorecards WHERE program_id = ?`, programID)
		if err != nil {
			tx.Rollback()
			log.Error().Err(err).Msg("syllabus.deleteSyllabusStructure")
			result.Error = constant.RespInternalServerError
			return c.Status(fiber.StatusInternalServerError).JSON(result)
		}
	}

	tx.Commit()

	result.Success = true
	return c.Status(fiber.StatusOK).JSON(result)
}

func (h *Handler) syllabuses(c *fiber.Ctx) error {
	var result struct {
		Nodes []*model.Syllabus `json:"nodes"`
		Error any               `json:"error"`
	}
	result.Nodes = []*model.Syllabus{}

	rows, err := h.db.QueryxContext(c.UserContext(), `
		SELECT s.id, s.parent_id, s.structure_id, s.title
		FROM syllabuses s
		JOIN syllabus_structures ss ON ss.id = s.structure_id
		WHERE ss.program_id = ?
		ORDER BY s.rowid
	`, c.Params("programId"))
	if err != nil {
		log.Error().Err(err).Msg("syllabus.syllabuses")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}
	defer rows.Close()

	for rows.Next() {
		var node model.Syllabus
		if err := rows.StructScan(&node); err != nil {
			log.Error().Err(err).Msg("syllabus.syllabuses")
			result.Error = constant.RespInternalServerError
			return c.Status(fiber.StatusInternalServerError).JSON(result)
		}
		result.Nodes = append(result.Nodes, &node)
	}
	c.Set("X-Total-Count", strconv.Itoa(len(result.Nodes)))

	return c.Status(fiber.StatusOK).JSON(result)
}

func (h *Handler) syllabus(c *fiber.Ctx) error {
	type Syllabus struct {
		model.BaseSyllabus
		Parents      []*model.BaseSyllabus `json:"parents"`
		IsAssignment bool                  `json:"isAssignment" db:"is_assignment"`
	}

	var result struct {
		Syllabus *Syllabus `json:"syllabus"`
		Error    any       `json:"error"`
	}

	syllabusID, _ := c.ParamsInt("syllabusId")

	rows, err := h.db.QueryxContext(c.UserContext(), `
		WITH RECURSIVE t AS (
		  SELECT s.*
		  FROM syllabuses s
		  JOIN syllabus_structures ss ON ss.id = s.structure_id
		  WHERE s.id = ?
		    AND ss.program_id = ?
		  UNION ALL
		  SELECT s.*
		  FROM syllabuses s
		  INNER JOIN t ON s.id = t.parent_id
		)
		SELECT t.id, t.title, COALESCE(ss.prev_id, 0) = -1 AS is_assignment
		FROM t
		JOIN syllabus_structures ss ON ss.id = t.structure_id
	`, syllabusID, c.Params("programId"))
	if err != nil {
		log.Error().Err(err).Msg("syllabus.syllabus")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}
	defer rows.Close()

	for rows.Next() {
		var node Syllabus
		if err := rows.StructScan(&node); err != nil {
			log.Error().Err(err).Msg("syllabus.syllabus")
			result.Error = constant.RespInternalServerError
			return c.Status(fiber.StatusInternalServerError).JSON(result)
		}

		if node.ID == syllabusID {
			result.Syllabus = &node
		} else {
			result.Syllabus.Parents = append(result.Syllabus.Parents, &node.BaseSyllabus)
		}
	}

	if len(result.Syllabus.Parents) > 1 {
		slices.Reverse(result.Syllabus.Parents)
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

func (h *Handler) saveSyllabus(c *fiber.Ctx) error {
	var result struct {
		Success bool `json:"success"`
		Error   any  `json:"error"`
	}

	var body struct {
		ParentID    *int   `json:"parentId"`
		StructureID int    `json:"structureId"`
		Title       string `json:"title"`
	}
	if err := c.BodyParser(&body); err != nil {
		log.Error().Err(err).Msg("syllabus.saveSyllabus")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}

	if syllabusID, _ := c.ParamsInt("syllabusID"); syllabusID != 0 {
		_, err := h.db.ExecContext(c.UserContext(), `
			UPDATE syllabuses
			SET parent_id = ?, title = ?
			WHERE id = ?
		`, body.ParentID, body.Title, syllabusID)
		if err != nil {
			if err, ok := err.(sqlite3.Error); ok && err.ExtendedCode == sqlite3.ErrConstraintUnique {
				result.Error = fiber.Map{"code": "TITLE_SHOULD_BE_UNIQUE"}
				return c.Status(fiber.StatusConflict).JSON(result)
			}
			log.Error().Err(err).Msg("syllabus.saveSyllabus")
			result.Error = constant.RespInternalServerError
			return c.Status(fiber.StatusInternalServerError).JSON(result)
		}
	} else {
		_, err := h.db.ExecContext(c.UserContext(), `
			INSERT INTO syllabuses (parent_id, structure_id, title)
			VALUES (?, ?, ?)
		`, body.ParentID, body.StructureID, body.Title)
		if err != nil {
			if err, ok := err.(sqlite3.Error); ok && err.ExtendedCode == sqlite3.ErrConstraintUnique {
				result.Error = fiber.Map{"code": "TITLE_SHOULD_BE_UNIQUE"}
				return c.Status(fiber.StatusConflict).JSON(result)
			}
			log.Error().Err(err).Msg("syllabus.saveSyllabus")
			result.Error = constant.RespInternalServerError
			return c.Status(fiber.StatusInternalServerError).JSON(result)
		}
	}

	result.Success = true
	return c.Status(fiber.StatusOK).JSON(result)
}

func (h *Handler) deleteSyllabus(c *fiber.Ctx) error {
	var result struct {
		Success bool `json:"success"`
		Error   any  `json:"error"`
	}

	programID, _ := c.ParamsInt("programId")
	syllabusID, _ := c.ParamsInt("syllabusId")

	tx := h.db.MustBeginTx(c.UserContext(), nil)

	// If structureID == 0, all syllabuses will be deleted
	_, err := tx.ExecContext(c.UserContext(), `
		DELETE FROM syllabuses
		WHERE (? = 0 OR id = ?)
	`, syllabusID, syllabusID)
	if err != nil {
		tx.Rollback()
		log.Error().Err(err).Msg("syllabus.deleteSyllabus")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}

	if syllabusID == 0 {
		_, err = tx.ExecContext(c.UserContext(), `DELETE FROM scorecard_structures WHERE program_id = ?`, programID)
		if err != nil {
			tx.Rollback()
			log.Error().Err(err).Msg("syllabus.deleteSyllabus")
			result.Error = constant.RespInternalServerError
			return c.Status(fiber.StatusInternalServerError).JSON(result)
		}

		_, err = tx.ExecContext(c.UserContext(), `DELETE FROM scorecards WHERE program_id = ?`, programID)
		if err != nil {
			tx.Rollback()
			log.Error().Err(err).Msg("syllabus.deleteSyllabus")
			result.Error = constant.RespInternalServerError
			return c.Status(fiber.StatusInternalServerError).JSON(result)
		}
	} else {
		_, err = tx.ExecContext(c.UserContext(), `
			UPDATE scorecards
			SET is_outdated = TRUE
			WHERE program_id = ?
			  AND EXISTS (
			    SELECT ss.id
			    FROM scorecard_structures ss
			    JOIN syllabuses s ON s.id = ss.syllabus_id
			    WHERE ss.program_id = ?
			      AND s.id = ?
			    LIMIT 1
			  )
		`, programID, programID, syllabusID)
		if err != nil {
			tx.Rollback()
			log.Error().Err(err).Msg("syllabus.deleteSyllabus")
			result.Error = constant.RespInternalServerError
			return c.Status(fiber.StatusInternalServerError).JSON(result)
		}
	}

	tx.Commit()

	result.Success = true
	return c.Status(fiber.StatusOK).JSON(result)
}

// ?limit int
// ?offset int

func (h *Handler) syllabusScores(c *fiber.Ctx) error {
	type Node struct {
		UserID int         `json:"-" db:"user_id"`
		User   *model.User `json:"user" db:"-"`
		Score  *float64    `json:"score"`
	}

	var result struct {
		Nodes []*Node `json:"nodes"`
		Error any     `json:"error"`
	}
	result.Nodes = []*Node{}

	qb := sq.Select().From("users u").
		LeftJoin("user_scores us ON us.user_id = u.id AND us.syllabus_id = ?", c.Params("syllabusId"))

	var totalCount int
	if err := qb.Column("COUNT(u.id)").RunWith(h.db).QueryRowContext(c.UserContext()).Scan(&totalCount); err != nil {
		log.Error().Err(err).Msg("syllabus.syllabusScores")
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

	query, args, err := qb.Columns("u.id AS user_id", "us.score").OrderBy("u.rowid ASC").ToSql()
	if err != nil {
		log.Error().Err(err).Msg("syllabus.syllabusScores")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}

	rows, err := h.db.QueryxContext(c.UserContext(), query, args...)
	if err != nil {
		log.Error().Err(err).Msg("syllabus.syllabusScores")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}
	defer rows.Close()

	var userIds []int
	for rows.Next() {
		var node Node
		if err := rows.StructScan(&node); err != nil {
			log.Error().Err(err).Msg("syllabus.syllabusScores")
			result.Error = constant.RespInternalServerError
			return c.Status(fiber.StatusInternalServerError).JSON(result)
		}
		result.Nodes = append(result.Nodes, &node)
		userIds = append(userIds, node.UserID)
	}

	if users, err := h.getUsers(c.UserContext(), userIds); err == nil {
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
		log.Error().Err(err).Msg("syllabus.saveScore")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}

	tx := h.db.MustBeginTx(c.UserContext(), nil)

	_, err := tx.ExecContext(c.UserContext(), `
		INSERT INTO user_scores (user_id, syllabus_id, score)
		VALUES (?, ?, ?)
		ON CONFLICT (user_id, syllabus_id)
		DO UPDATE SET score = EXCLUDED.score
	`, userID, c.Params("syllabusId"), body.Score)
	if err != nil {
		tx.Rollback()
		log.Error().Err(err).Msg("syllabus.saveScore")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}

	_, err = tx.ExecContext(c.UserContext(), `UPDATE scorecards SET is_outdated = TRUE WHERE user_id = ?`, userID)
	if err != nil {
		tx.Rollback()
		log.Error().Err(err).Msg("syllabus.saveScore")
		result.Error = constant.RespInternalServerError
		return c.Status(fiber.StatusInternalServerError).JSON(result)
	}

	tx.Commit()

	result.Success = true
	return c.Status(fiber.StatusOK).JSON(result)
}
