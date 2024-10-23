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

func Test_programs(t *testing.T) {
	assert := assert.New(t)

	t.Run("HEAD", func(t *testing.T) {
		db, mock := db.New()
		h := New(db, nil)

		mock.ExpectQuery(`SELECT COUNT\(id\) FROM programs`).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		app := fiber.New()
		h.Register(app, middleware.New())

		req := httptest.NewRequest("HEAD", "/v1/programs", nil)

		resp, _ := app.Test(req)
		assert.Nil(mock.ExpectationsWereMet())
		assert.Equal(fiber.StatusOK, resp.StatusCode)
		assert.Equal("1", resp.Header.Get("X-Total-Count"))
	})

	t.Run("empty", func(t *testing.T) {
		db, mock := db.New()
		h := New(db, nil)

		mock.ExpectQuery(`SELECT COUNT\(id\) FROM programs`).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		app := fiber.New()
		h.Register(app, middleware.New())

		req := httptest.NewRequest("GET", "/v1/programs", nil)

		resp, _ := app.Test(req)
		assert.Nil(mock.ExpectationsWereMet())
		assert.Equal(fiber.StatusOK, resp.StatusCode)
		assert.Equal("0", resp.Header.Get("X-Total-Count"))
		body, _ := io.ReadAll(resp.Body)
		assert.Equal(`{"nodes":[],"error":null}`, string(body))
	})

	t.Run("success", func(t *testing.T) {
		db, mock := db.New()
		h := New(db, nil)

		mock.ExpectQuery(`SELECT COUNT\(id\) FROM programs`).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

		mock.ExpectQuery("SELECT .+ FROM programs .+ LIMIT 1 OFFSET 1").
			WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).AddRow(2, "Program 2"))

		app := fiber.New()
		h.Register(app, middleware.New())

		req := httptest.NewRequest("GET", "/v1/programs?limit=1&offset=1", nil)

		resp, _ := app.Test(req)
		assert.Nil(mock.ExpectationsWereMet())
		assert.Equal(fiber.StatusOK, resp.StatusCode)
		assert.Equal("2", resp.Header.Get("X-Total-Count"))
		body, _ := io.ReadAll(resp.Body)
		assert.Equal(`{"nodes":[{"id":2,"title":"Program 2"}],"error":null}`, string(body))
	})
}

func Test_program(t *testing.T) {
	db, mock := db.New()
	h := New(db, nil)

	mock.ExpectQuery("SELECT .+ FROM programs").
		WithArgs("1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).AddRow(1, "Program 1"))

	app := fiber.New()
	h.Register(app, middleware.New())

	req := httptest.NewRequest("GET", "/v1/programs/1", nil)

	resp, _ := app.Test(req)
	assert.Nil(t, mock.ExpectationsWereMet())
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, `{"program":{"id":1,"title":"Program 1"},"error":null}`, string(body))
}

func Test_saveProgram(t *testing.T) {
	assert := assert.New(t)

	t.Run("insert not unique", func(t *testing.T) {
		db, mock := db.New()
		h := New(db, nil)

		mock.ExpectExec("INSERT INTO programs").
			WithArgs("Program 1").
			WillReturnError(sqlite3.Error{ExtendedCode: sqlite3.ErrConstraintUnique})

		app := fiber.New()
		h.Register(app, middleware.New())

		req := httptest.NewRequest("PUT", "/v1/programs", strings.NewReader(`{"title":"Program 1"}`))
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

		mock.ExpectExec("INSERT INTO programs").
			WithArgs("Program 1").
			WillReturnResult(sqlmock.NewResult(1, 1))

		app := fiber.New()
		h.Register(app, middleware.New())

		req := httptest.NewRequest("PUT", "/v1/programs", strings.NewReader(`{"title":"Program 1"}`))
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

		mock.ExpectExec("UPDATE programs").
			WithArgs("Program 1a", 1).
			WillReturnError(sqlite3.Error{ExtendedCode: sqlite3.ErrConstraintUnique})

		app := fiber.New()
		h.Register(app, middleware.New())

		req := httptest.NewRequest("PUT", "/v1/programs/1", strings.NewReader(`{"title":"Program 1a"}`))
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

		mock.ExpectExec("UPDATE programs").
			WithArgs("Program 1a", 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		app := fiber.New()
		h.Register(app, middleware.New())

		req := httptest.NewRequest("PUT", "/v1/programs/1", strings.NewReader(`{"title":"Program 1a"}`))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Nil(mock.ExpectationsWereMet())
		assert.Equal(fiber.StatusOK, resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		assert.Equal(`{"success":true,"error":null}`, string(body))
	})
}

func Test_deleteProgram(t *testing.T) {
	db, mock := db.New()
	h := New(db, nil)

	mock.ExpectExec("DELETE FROM programs").
		WithArgs("1").
		WillReturnResult(sqlmock.NewResult(0, 1))

	app := fiber.New()
	h.Register(app, middleware.New())

	req := httptest.NewRequest("DELETE", "/v1/programs/1", nil)

	resp, _ := app.Test(req)
	assert.Nil(t, mock.ExpectationsWereMet())
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}
