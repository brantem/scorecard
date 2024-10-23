package handler

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/brantem/scorecard/testutil/db"
	"github.com/brantem/scorecard/testutil/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func Test_syllabusStructures(t *testing.T) {
	db, mock := db.New()
	h := New(db, nil)

	mock.ExpectQuery("SELECT .+ FROM syllabus_structures").
		WithArgs(1, 1).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "prev_id", "title"}).
				AddRow(2, nil, "Structure 1").
				AddRow(1, -1, "Assignment"),
		)

	app := fiber.New()
	h.Register(app, middleware.New())

	req := httptest.NewRequest("GET", "/v1/programs/1/syllabuses/structures", nil)

	resp, _ := app.Test(req)
	assert.Nil(t, mock.ExpectationsWereMet())
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	assert.Equal(t, "2", resp.Header.Get("X-Total-Count"))
	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, `{"nodes":[{"id":2,"prevId":null,"title":"Structure 1"},{"id":1,"prevId":-1,"title":"Assignment"}],"error":null}`, string(body))
}

func Test_saveSyllabusStructure(t *testing.T) {
	assert := assert.New(t)

	t.Run("insert not unique", func(t *testing.T) {
		db, mock := db.New()
		h := New(db, nil)

		mock.ExpectBegin()

		mock.ExpectExec("INSERT INTO syllabus_structures .+ 'Assignment'").
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectExec("INSERT INTO syllabus_structures").
			WithArgs(1, nil, "Structure 1").
			WillReturnError(sqlite3.Error{ExtendedCode: sqlite3.ErrConstraintUnique})

		mock.ExpectRollback()

		app := fiber.New()
		h.Register(app, middleware.New())

		req := httptest.NewRequest("PUT", "/v1/programs/1/syllabuses/structures", strings.NewReader(`{"prevId":null,"title":"Structure 1"}`))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Nil(mock.ExpectationsWereMet())
		assert.Equal(fiber.StatusConflict, resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		assert.Equal(`{"success":false,"error":{"code":"TITLE_SHOULD_BE_UNIQUE"}}`, string(body))
	})

	t.Run("insert", func(t *testing.T) {
		db, mock := db.New()
		h := New(db, nil)

		mock.ExpectBegin()

		mock.ExpectExec("INSERT INTO syllabus_structures .+ 'Assignment'").
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectExec("INSERT INTO syllabus_structures").
			WithArgs(1, nil, "Structure 1").
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectCommit()

		app := fiber.New()
		h.Register(app, middleware.New())

		req := httptest.NewRequest("PUT", "/v1/programs/1/syllabuses/structures", strings.NewReader(`{"prevId":null,"title":"Structure 1"}`))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Nil(mock.ExpectationsWereMet())
		assert.Equal(fiber.StatusOK, resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		assert.Equal(`{"success":true,"error":null}`, string(body))
	})

	t.Run("update not unique", func(t *testing.T) {
		db, mock := db.New()
		h := New(db, nil)

		mock.ExpectExec("UPDATE syllabus_structures").
			WithArgs("Structure 2a", 2).
			WillReturnError(sqlite3.Error{ExtendedCode: sqlite3.ErrConstraintUnique})

		app := fiber.New()
		h.Register(app, middleware.New())

		req := httptest.NewRequest("PUT", "/v1/programs/1/syllabuses/structures/2", strings.NewReader(`{"title":"Structure 2a"}`))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Nil(mock.ExpectationsWereMet())
		assert.Equal(fiber.StatusConflict, resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		assert.Equal(`{"success":false,"error":{"code":"TITLE_SHOULD_BE_UNIQUE"}}`, string(body))
	})

	t.Run("update", func(t *testing.T) {
		db, mock := db.New()
		h := New(db, nil)

		mock.ExpectExec("UPDATE syllabus_structures").
			WithArgs("Structure 2a", 2).
			WillReturnResult(sqlmock.NewResult(0, 1))

		app := fiber.New()
		h.Register(app, middleware.New())

		req := httptest.NewRequest("PUT", "/v1/programs/1/syllabuses/structures/2", strings.NewReader(`{"title":"Structure 2a"}`))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Nil(mock.ExpectationsWereMet())
		assert.Equal(fiber.StatusOK, resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		assert.Equal(`{"success":true,"error":null}`, string(body))
	})
}

func Test_deleteSyllabusStructure(t *testing.T) {
	assert := assert.New(t)

	t.Run("structureId != 0", func(t *testing.T) {
		db, mock := db.New()
		h := New(db, nil)

		mock.ExpectBegin()

		mock.ExpectExec("DELETE FROM syllabus_structures").
			WithArgs(1, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectCommit()

		app := fiber.New()
		h.Register(app, middleware.New())

		req := httptest.NewRequest("DELETE", "/v1/programs/1/syllabuses/structures/1", nil)

		resp, _ := app.Test(req)
		assert.Nil(mock.ExpectationsWereMet())
		assert.Equal(fiber.StatusOK, resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		assert.Equal(`{"success":true,"error":null}`, string(body))
	})

	t.Run("structureId == 0", func(t *testing.T) {
		db, mock := db.New()
		h := New(db, nil)

		mock.ExpectBegin()

		mock.ExpectExec("DELETE FROM syllabus_structures").
			WithArgs(0, 0).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectExec("DELETE FROM scorecard_structures").
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectExec("DELETE FROM scorecards").
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectCommit()

		app := fiber.New()
		h.Register(app, middleware.New())

		req := httptest.NewRequest("DELETE", "/v1/programs/1/syllabuses/structures/0", nil)

		resp, _ := app.Test(req)
		assert.Nil(mock.ExpectationsWereMet())
		assert.Equal(fiber.StatusOK, resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		assert.Equal(`{"success":true,"error":null}`, string(body))
	})
}

func Test_syllabuses(t *testing.T) {
	db, mock := db.New()
	h := New(db, nil)

	mock.ExpectQuery("SELECT .+ FROM syllabuses").
		WithArgs("1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "parent_id", "structure_id", "title"}).AddRow(1, nil, 1, "Syllabus 1"))

	app := fiber.New()
	h.Register(app, middleware.New())

	req := httptest.NewRequest("GET", "/v1/programs/1/syllabuses", nil)

	resp, _ := app.Test(req)
	assert.Nil(t, mock.ExpectationsWereMet())
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	assert.Equal(t, "1", resp.Header.Get("X-Total-Count"))
	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, `{"nodes":[{"id":1,"title":"Syllabus 1","parentId":null,"structureId":1}],"error":null}`, string(body))
}

func Test_syllabus(t *testing.T) {
	db, mock := db.New()
	h := New(db, nil)

	mock.ExpectQuery("SELECT .+ FROM syllabuses").
		WithArgs(3, "1").
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "title", "is_assignment"}).
				AddRow(3, "Syllabus 3", true).
				AddRow(2, "Syllabus 2", false).
				AddRow(1, "Syllabus 1", false),
		)

	app := fiber.New()
	h.Register(app, middleware.New())

	req := httptest.NewRequest("GET", "/v1/programs/1/syllabuses/3", nil)

	resp, _ := app.Test(req)
	assert.Nil(t, mock.ExpectationsWereMet())
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, `{"syllabus":{"id":3,"title":"Syllabus 3","parents":[{"id":1,"title":"Syllabus 1"},{"id":2,"title":"Syllabus 2"}],"isAssignment":true},"error":null}`, string(body))
}

func Test_saveSyllabus(t *testing.T) {
	assert := assert.New(t)

	t.Run("insert not unique", func(t *testing.T) {
		db, mock := db.New()
		h := New(db, nil)

		mock.ExpectExec("INSERT INTO syllabuses").
			WithArgs(nil, 1, "Syllabus 1").
			WillReturnError(sqlite3.Error{ExtendedCode: sqlite3.ErrConstraintUnique})

		app := fiber.New()
		h.Register(app, middleware.New())

		req := httptest.NewRequest("PUT", "/v1/programs/1/syllabuses", strings.NewReader(`{"parentId":null,"structureId":1,"title":"Syllabus 1"}`))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Nil(mock.ExpectationsWereMet())
		assert.Equal(fiber.StatusConflict, resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		assert.Equal(`{"success":false,"error":{"code":"TITLE_SHOULD_BE_UNIQUE"}}`, string(body))
	})

	t.Run("insert", func(t *testing.T) {
		db, mock := db.New()
		h := New(db, nil)

		mock.ExpectExec("INSERT INTO syllabuses").
			WithArgs(nil, 1, "Syllabus 1").
			WillReturnResult(sqlmock.NewResult(1, 1))

		app := fiber.New()
		h.Register(app, middleware.New())

		req := httptest.NewRequest("PUT", "/v1/programs/1/syllabuses", strings.NewReader(`{"parentId":null,"structureId":1,"title":"Syllabus 1"}`))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Nil(mock.ExpectationsWereMet())
		assert.Equal(fiber.StatusOK, resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		assert.Equal(`{"success":true,"error":null}`, string(body))
	})

	t.Run("update not unique", func(t *testing.T) {
		db, mock := db.New()
		h := New(db, nil)

		mock.ExpectExec("UPDATE syllabuses").
			WithArgs(1, "Syllabus 2a", 2).
			WillReturnError(sqlite3.Error{ExtendedCode: sqlite3.ErrConstraintUnique})

		app := fiber.New()
		h.Register(app, middleware.New())

		req := httptest.NewRequest("PUT", "/v1/programs/1/syllabuses/2", strings.NewReader(`{"parentId":1,"title":"Syllabus 2a"}`))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Nil(mock.ExpectationsWereMet())
		assert.Equal(fiber.StatusConflict, resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		assert.Equal(`{"success":false,"error":{"code":"TITLE_SHOULD_BE_UNIQUE"}}`, string(body))
	})

	t.Run("update", func(t *testing.T) {
		db, mock := db.New()
		h := New(db, nil)

		mock.ExpectExec("UPDATE syllabuses").
			WithArgs(1, "Syllabus 2a", 2).
			WillReturnResult(sqlmock.NewResult(1, 1))

		app := fiber.New()
		h.Register(app, middleware.New())

		req := httptest.NewRequest("PUT", "/v1/programs/1/syllabuses/2", strings.NewReader(`{"parentId":1,"title":"Syllabus 2a"}`))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Nil(mock.ExpectationsWereMet())
		assert.Equal(fiber.StatusOK, resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		assert.Equal(`{"success":true,"error":null}`, string(body))
	})
}

func Test_deleteSyllabus(t *testing.T) {
	assert := assert.New(t)

	t.Run("syllabusId == 0", func(t *testing.T) {
		db, mock := db.New()
		h := New(db, nil)

		mock.ExpectBegin()

		mock.ExpectExec("DELETE FROM syllabuses").
			WithArgs(0, 0).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectExec("DELETE FROM scorecard_structures").
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectExec("DELETE FROM scorecards").
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectCommit()

		app := fiber.New()
		h.Register(app, middleware.New())

		req := httptest.NewRequest("DELETE", "/v1/programs/1/syllabuses/0", nil)

		resp, _ := app.Test(req)
		assert.Nil(mock.ExpectationsWereMet())
		assert.Equal(fiber.StatusOK, resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		assert.Equal(`{"success":true,"error":null}`, string(body))
	})

	t.Run("syllabusId != 0", func(t *testing.T) {
		db, mock := db.New()
		h := New(db, nil)

		mock.ExpectBegin()

		mock.ExpectExec("DELETE FROM syllabuses").
			WithArgs(2, 2).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectExec("UPDATE scorecards").
			WithArgs(1, 1, 2).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectCommit()

		app := fiber.New()
		h.Register(app, middleware.New())

		req := httptest.NewRequest("DELETE", "/v1/programs/1/syllabuses/2", nil)

		resp, _ := app.Test(req)
		assert.Nil(mock.ExpectationsWereMet())
		assert.Equal(fiber.StatusOK, resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		assert.Equal(`{"success":true,"error":null}`, string(body))
	})
}

func Test_syllabusScores(t *testing.T) {
	db, mock := db.New()
	h := New(db, nil)

	mock.ExpectQuery("SELECT .+ FROM users .+ JOIN user_scores").
		WithArgs("2").
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "score"}).AddRow(1, nil))

	mock.ExpectQuery(`SELECT .+ FROM users WHERE id IN \(\?\)`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "User 1"))

	app := fiber.New()
	h.Register(app, middleware.New())

	req := httptest.NewRequest("GET", "/v1/programs/1/syllabuses/2/scores", nil)

	resp, _ := app.Test(req)
	assert.Nil(t, mock.ExpectationsWereMet())
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	assert.Equal(t, "1", resp.Header.Get("X-Total-Count"))
	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, `{"nodes":[{"user":{"id":1,"name":"User 1"},"score":null}],"error":null}`, string(body))
}

func Test_saveScore(t *testing.T) {
	db, mock := db.New()
	h := New(db, nil)

	mock.ExpectBegin()

	mock.ExpectExec("INSERT INTO user_scores").
		WithArgs(3, "2", float64(100)).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectExec("UPDATE scorecards").
		WithArgs(3).
		WillReturnResult(sqlmock.NewResult(0, 1))

	mock.ExpectCommit()

	app := fiber.New()
	h.Register(app, middleware.New())

	req := httptest.NewRequest("PUT", "/v1/programs/1/syllabuses/2/scores/3", strings.NewReader(`{"score":100}`))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)
	assert.Nil(t, mock.ExpectationsWereMet())
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, `{"success":true,"error":null}`, string(body))
}
