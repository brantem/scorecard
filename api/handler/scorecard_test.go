package handler

import (
	"context"
	"fmt"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/brantem/scorecard/model"
	"github.com/brantem/scorecard/testutil/db"
	"github.com/brantem/scorecard/testutil/middleware"
	"github.com/brantem/scorecard/testutil/scorecard"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func Test_getScorecardStructureSyllabuses(t *testing.T) {
	assert := assert.New(t)

	syllabus1 := model.ScorecardStructureSyllabus{
		BaseSyllabus: model.BaseSyllabus{
			ID:    1,
			Title: "Syllabus 1",
		},
		IsAssignment: false,
	}
	syllabus2 := model.ScorecardStructureSyllabus{
		BaseSyllabus: model.BaseSyllabus{
			ID:    2,
			Title: "Syllabus 2",
		},
		IsAssignment: true,
	}

	t.Run("empty", func(t *testing.T) {
		h := New(nil, nil)

		m, err := h.getScorecardStructureSyllabuses(context.Background(), []int{})
		assert.Empty(m)
		assert.Nil(err)
	})

	t.Run("success", func(t *testing.T) {
		db, mock := db.New()
		h := New(db, nil)

		mock.ExpectQuery("SELECT .+ FROM syllabuses .+ WHERE s.id IN (?)").
			WithArgs(syllabus1.ID, syllabus2.ID).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "title", "is_assignment"}).
					AddRow(syllabus1.ID, syllabus1.Title, syllabus1.IsAssignment).
					AddRow(syllabus2.ID, syllabus2.Title, syllabus2.IsAssignment),
			)

		m, err := h.getScorecardStructureSyllabuses(context.Background(), []int{syllabus1.ID, syllabus2.ID})
		assert.Equal(map[int]*model.ScorecardStructureSyllabus{1: &syllabus1, 2: &syllabus2}, m)
		assert.Nil(err)
	})
}

func Test_scorecardStructures(t *testing.T) {
	db, mock := db.New()
	h := New(db, nil)

	mock.ExpectQuery("SELECT .+ FROM scorecard_structures").
		WithArgs(1, 1, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "parent_id", "title", "syllabus_id"}).AddRow(1, nil, "Structure 1", 1))

	mock.ExpectQuery("SELECT .+ FROM syllabuses .+ WHERE s.id IN (?)").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "is_assignment"}).AddRow(1, "Syllabus 1", false))

	app := fiber.New()
	h.Register(app, middleware.New())

	req := httptest.NewRequest("GET", "/v1/programs/1/scorecards/structures?depth=1", nil)

	resp, _ := app.Test(req)
	assert.Nil(t, mock.ExpectationsWereMet())
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	assert.Equal(t, "1", resp.Header.Get("X-Total-Count"))
	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, `{"nodes":[{"id":1,"parentId":null,"title":"Structure 1","syllabus":{"id":1,"title":"Syllabus 1","isAssignment":false}}],"error":null}`, string(body))
}

func Test_copySyllabusesIntoStructures(t *testing.T) {
	syllabusIds := [2]int{2, 3}

	db, mock := db.New()
	h := New(db, nil)

	mock.ExpectQuery(`.+ SELECT \* FROM syllabuses`).
		WithArgs(0, 0).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "parent_id"}).
				AddRow(syllabusIds[0], 1).
				AddRow(syllabusIds[1], syllabusIds[0]),
		)

	mock.ExpectBegin()

	structureIds := [2]int{1, 2}
	mock.ExpectQuery("INSERT INTO scorecard_structures").
		WithArgs(1, syllabusIds[0], syllabusIds[1]).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "syllabus_id"}).
				AddRow(structureIds[0], syllabusIds[0]).
				AddRow(structureIds[1], syllabusIds[1]),
		)

	mock.ExpectExec(fmt.Sprintf(
		`WITH .+ \(%d, NULL\), \(%d, %d\) .+ UPDATE scorecard_structures`,
		structureIds[0],
		structureIds[1], structureIds[0],
	)).
		WillReturnResult(sqlmock.NewResult(0, 2))

	mock.ExpectExec("UPDATE scorecards").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectCommit()

	app := fiber.New()
	h.Register(app, middleware.New())

	req := httptest.NewRequest("POST", "/v1/programs/1/scorecards/structures/copy/0", strings.NewReader(`{"parentId":null}`))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	assert.Nil(t, mock.ExpectationsWereMet())
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, `{"success":true,"error":null}`, string(body))
}

func Test_saveScorecardStructure(t *testing.T) {
	assert := assert.New(t)

	t.Run("insert", func(t *testing.T) {
		db, mock := db.New()
		h := New(db, nil)

		mock.ExpectExec("INSERT INTO scorecard_structures").
			WithArgs(1, nil, "Structure 1").
			WillReturnResult(sqlmock.NewResult(0, 1))

		app := fiber.New()
		h.Register(app, middleware.New())

		req := httptest.NewRequest("PUT", "/v1/programs/1/scorecards/structures", strings.NewReader(`{"prevId":null,"title":"Structure 1"}`))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Nil(mock.ExpectationsWereMet())
		assert.Equal(fiber.StatusOK, resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		assert.Equal(`{"success":true,"error":null}`, string(body))
	})

	t.Run("update", func(t *testing.T) {
		db, mock := db.New()
		h := New(db, nil)

		mock.ExpectExec("UPDATE scorecard_structures").
			WithArgs(nil, "Structure 2a", 2).
			WillReturnResult(sqlmock.NewResult(0, 1))

		app := fiber.New()
		h.Register(app, middleware.New())

		req := httptest.NewRequest("PUT", "/v1/programs/1/scorecards/structures/2", strings.NewReader(`{"parentId":null,"title":"Structure 2a"}`))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Nil(mock.ExpectationsWereMet())
		assert.Equal(fiber.StatusOK, resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		assert.Equal(`{"success":true,"error":null}`, string(body))
	})
}

func Test_deleteScorecardStructure(t *testing.T) {
	assert := assert.New(t)

	t.Run("structureId == 0", func(t *testing.T) {
		db, mock := db.New()
		h := New(db, nil)

		mock.ExpectBegin()

		mock.ExpectExec("DELETE FROM scorecard_structures").
			WithArgs(0, 0).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectExec("DELETE FROM scorecards").
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectCommit()

		app := fiber.New()
		h.Register(app, middleware.New())

		req := httptest.NewRequest("DELETE", "/v1/programs/1/scorecards/structures/0", nil)

		resp, _ := app.Test(req)
		assert.Nil(mock.ExpectationsWereMet())
		assert.Equal(fiber.StatusOK, resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		assert.Equal(`{"success":true,"error":null}`, string(body))
	})

	t.Run("structureId != 0", func(t *testing.T) {
		db, mock := db.New()
		h := New(db, nil)

		mock.ExpectBegin()

		mock.ExpectExec("DELETE FROM scorecard_structures").
			WithArgs(2, 2).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectExec("UPDATE scorecards").
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectCommit()

		app := fiber.New()
		h.Register(app, middleware.New())

		req := httptest.NewRequest("DELETE", "/v1/programs/1/scorecards/structures/2", nil)

		resp, _ := app.Test(req)
		assert.Nil(mock.ExpectationsWereMet())
		assert.Equal(fiber.StatusOK, resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		assert.Equal(`{"success":true,"error":null}`, string(body))
	})
}

func Test_generateScorecards(t *testing.T) {
	assert := assert.New(t)

	t.Run("structureId != 0", func(t *testing.T) {
		db, mock := db.New()
		generator := scorecard.NewGenerator()
		h := New(db, generator)

		mock.ExpectQuery("SELECT .+ FROM scorecards").
			WithArgs(1, 2).
			WillReturnRows(sqlmock.NewRows([]string{"user_id"}).AddRow(3))

		app := fiber.New()
		h.Register(app, middleware.New())

		req := httptest.NewRequest("POST", "/v1/programs/1/scorecards/generate/2", nil)

		resp, _ := app.Test(req)
		assert.Nil(mock.ExpectationsWereMet())
		assert.Equal([]int{1}, generator.EnqueueProgramID)
		assert.Equal([]int{3}, generator.EnqueueUserID)
		assert.Equal([]int{2}, generator.EnqueueScorecardID)
		assert.Equal(fiber.StatusOK, resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		assert.Equal(`{"success":true,"error":null}`, string(body))
	})

	t.Run("structureId == 0", func(t *testing.T) {
		db, mock := db.New()
		generator := scorecard.NewGenerator()
		h := New(db, generator)

		mock.ExpectQuery("SELECT .+ FROM syllabus_structures").
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"user_id", "scorecard_id"}).AddRow(1, 1).AddRow(2, 0))

		app := fiber.New()
		h.Register(app, middleware.New())

		req := httptest.NewRequest("POST", "/v1/programs/1/scorecards/generate", nil)

		resp, _ := app.Test(req)
		assert.Nil(mock.ExpectationsWereMet())
		assert.Equal([]int{1, 1}, generator.EnqueueProgramID)
		assert.Equal([]int{1, 2}, generator.EnqueueUserID)
		assert.Equal([]int{1, 0}, generator.EnqueueScorecardID)
		assert.Equal(fiber.StatusOK, resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		assert.Equal(`{"success":true,"error":null}`, string(body))
	})
}

func Test_scorecards(t *testing.T) {
	db, mock := db.New()
	h := New(db, nil)

	mock.ExpectQuery("SELECT .+ FROM scorecards").
		WithArgs("1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "score", "is_outdated", "generated_at"}).AddRow(1, 1, 100, false, "2024-01-01 00:00:00"))

	mock.ExpectQuery("SELECT .+ FROM users WHERE id IN (?)").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "User 1"))

	app := fiber.New()
	h.Register(app, middleware.New())

	req := httptest.NewRequest("GET", "/v1/programs/1/scorecards", nil)

	resp, _ := app.Test(req)
	assert.Nil(t, mock.ExpectationsWereMet())
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	assert.Equal(t, "1", resp.Header.Get("X-Total-Count"))
	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, `{"nodes":[{"id":1,"user":{"id":1,"name":"User 1"},"score":100,"items":null,"isOutdated":false,"generatedAt":"2024-01-01T00:00:00Z"}],"error":null}`, string(body))
}

func Test_scorecard(t *testing.T) {
	db, mock := db.New()
	h := New(db, nil)

	mock.MatchExpectationsInOrder(false)

	mock.ExpectQuery("SELECT .+ FROM scorecards").
		WithArgs("1", 2).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "score", "is_outdated", "generated_at"}).AddRow(2, 1, 100, false, "2024-01-01 00:00:00"))

	mock.ExpectQuery("SELECT .+ FROM users WHERE id IN (?)").
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "User 1"))

	mock.ExpectQuery("SELECT .+ FROM scorecard_items").
		WithArgs(2).
		WillReturnRows(sqlmock.NewRows([]string{"id", "structure_id", "score"}).AddRow(1, 1, 100))

	app := fiber.New()
	h.Register(app, middleware.New())

	req := httptest.NewRequest("GET", "/v1/programs/1/scorecards/2", nil)

	resp, _ := app.Test(req)
	assert.Nil(t, mock.ExpectationsWereMet())
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, `{"scorecard":{"id":2,"user":{"id":1,"name":"User 1"},"score":100,"items":[{"id":1,"structureId":1,"score":100}],"isOutdated":false,"generatedAt":"2024-01-01T00:00:00Z"},"error":null}`, string(body))
}
