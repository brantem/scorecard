package handler

import (
	"context"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/brantem/scorecard/model"
	"github.com/brantem/scorecard/testutil/db"
	"github.com/brantem/scorecard/testutil/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func Test_getUsers(t *testing.T) {
	assert := assert.New(t)

	user1 := model.User{ID: 1, Name: "User 1"}
	user2 := model.User{ID: 2, Name: "User 2"}

	t.Run("empty", func(t *testing.T) {
		h := New(nil, nil)

		m, err := h.getUsers(context.Background(), []int{})
		assert.Empty(m)
		assert.Nil(err)
	})

	t.Run("success", func(t *testing.T) {
		db, mock := db.New()
		h := New(db, nil)

		mock.ExpectQuery(`SELECT .+ FROM users WHERE id IN \(\?, \?\)`).
			WithArgs(user1.ID, user2.ID).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "name"}).
					AddRow(user1.ID, user1.Name).
					AddRow(user2.ID, user2.Name),
			)

		m, err := h.getUsers(context.Background(), []int{user1.ID, user2.ID})
		assert.Equal(map[int]*model.User{1: &user1, 2: &user2}, m)
		assert.Nil(err)
	})
}

func Test_users(t *testing.T) {
	assert := assert.New(t)

	t.Run("HEAD", func(t *testing.T) {
		db, mock := db.New()
		h := New(db, nil)

		mock.ExpectQuery(`SELECT COUNT\(id\) FROM users`).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

		app := fiber.New()
		h.Register(app, middleware.New())

		req := httptest.NewRequest("HEAD", "/v1/programs/1/users", nil)

		resp, _ := app.Test(req)
		assert.Nil(mock.ExpectationsWereMet())
		assert.Equal(fiber.StatusOK, resp.StatusCode)
		assert.Equal("1", resp.Header.Get("X-Total-Count"))
	})

	t.Run("empty", func(t *testing.T) {
		db, mock := db.New()
		h := New(db, nil)

		mock.ExpectQuery(`SELECT COUNT\(id\) FROM users`).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

		app := fiber.New()
		h.Register(app, middleware.New())

		req := httptest.NewRequest("GET", "/v1/programs/1/users", nil)

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

		mock.ExpectQuery(`SELECT COUNT\(id\) FROM users`).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

		mock.ExpectQuery("SELECT .+ FROM users .+ LIMIT 1 OFFSET 1").
			WithArgs("1").
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(2, "User 2"))

		app := fiber.New()
		h.Register(app, middleware.New())

		req := httptest.NewRequest("GET", "/v1/programs/1/users?limit=1&offset=1", nil)

		resp, _ := app.Test(req)
		assert.Nil(mock.ExpectationsWereMet())
		assert.Equal(fiber.StatusOK, resp.StatusCode)
		assert.Equal("2", resp.Header.Get("X-Total-Count"))
		body, _ := io.ReadAll(resp.Body)
		assert.Equal(`{"nodes":[{"id":2,"name":"User 2"}],"error":null}`, string(body))
	})
}

func Test_user(t *testing.T) {
	db, mock := db.New()
	h := New(db, nil)

	mock.ExpectQuery("SELECT .+ FROM users").
		WithArgs("1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "User 1"))

	app := fiber.New()
	h.Register(app, middleware.New())

	req := httptest.NewRequest("GET", "/v1/programs/1/users/1", nil)

	resp, _ := app.Test(req)
	assert.Nil(t, mock.ExpectationsWereMet())
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, `{"user":{"id":1,"name":"User 1"},"error":null}`, string(body))
}

func Test_userScores(t *testing.T) {
	db, mock := db.New()
	h := New(db, nil)

	mock.ExpectQuery("SELECT .+ FROM syllabus_structures").
		WithArgs("1", "1").
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "parent_id", "title", "score", "is_assignment"}).
				AddRow(1, nil, "Syllabus 1", 0, false).
				AddRow(2, 1, "Syllabus 2", 100, true),
		)

	app := fiber.New()
	h.Register(app, middleware.New())

	req := httptest.NewRequest("GET", "/v1/programs/1/users/1/scores", nil)

	resp, _ := app.Test(req)
	assert.Nil(t, mock.ExpectationsWereMet())
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	assert.Equal(t, "1", resp.Header.Get("X-Total-Count"))
	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, `{"nodes":[{"syllabus":{"id":2,"title":"Syllabus 2","parents":[{"id":1,"title":"Syllabus 1"}]},"score":100}],"error":null}`, string(body))
}

func Test_saveUser(t *testing.T) {
	assert := assert.New(t)

	t.Run("insert not unique", func(t *testing.T) {
		db, mock := db.New()
		h := New(db, nil)

		mock.ExpectExec("INSERT INTO users").
			WithArgs("1", "User 1").
			WillReturnError(sqlite3.Error{ExtendedCode: sqlite3.ErrConstraintUnique})

		app := fiber.New()
		h.Register(app, middleware.New())

		req := httptest.NewRequest("PUT", "/v1/programs/1/users", strings.NewReader(`{"name":"User 1"}`))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Nil(mock.ExpectationsWereMet())
		assert.Equal(fiber.StatusConflict, resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		assert.Equal(`{"success":false,"error":{"code":"NAME_SHOULD_BE_UNIQUE"}}`, string(body))
	})

	t.Run("insert", func(t *testing.T) {
		db, mock := db.New()
		h := New(db, nil)

		mock.ExpectExec("INSERT INTO users").
			WithArgs("1", "User 1").
			WillReturnResult(sqlmock.NewResult(1, 1))

		app := fiber.New()
		h.Register(app, middleware.New())

		req := httptest.NewRequest("PUT", "/v1/programs/1/users", strings.NewReader(`{"name":"User 1"}`))
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

		mock.ExpectExec("UPDATE users").
			WithArgs("User 1a", 1).
			WillReturnError(sqlite3.Error{ExtendedCode: sqlite3.ErrConstraintUnique})

		app := fiber.New()
		h.Register(app, middleware.New())

		req := httptest.NewRequest("PUT", "/v1/programs/1/users/1", strings.NewReader(`{"name":"User 1a"}`))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Nil(mock.ExpectationsWereMet())
		assert.Equal(fiber.StatusConflict, resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		assert.Equal(`{"success":false,"error":{"code":"NAME_SHOULD_BE_UNIQUE"}}`, string(body))
	})

	t.Run("update", func(t *testing.T) {
		db, mock := db.New()
		h := New(db, nil)

		mock.ExpectExec("UPDATE users").
			WithArgs("User 2a", 2).
			WillReturnResult(sqlmock.NewResult(1, 1))

		app := fiber.New()
		h.Register(app, middleware.New())

		req := httptest.NewRequest("PUT", "/v1/programs/1/users/2", strings.NewReader(`{"name":"User 2a"}`))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Nil(mock.ExpectationsWereMet())
		assert.Equal(fiber.StatusOK, resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		assert.Equal(`{"success":true,"error":null}`, string(body))
	})
}

func Test_deleteUser(t *testing.T) {
	db, mock := db.New()
	h := New(db, nil)

	mock.ExpectExec("DELETE FROM users").
		WithArgs("2").
		WillReturnResult(sqlmock.NewResult(0, 1))

	app := fiber.New()
	h.Register(app, middleware.New())

	req := httptest.NewRequest("DELETE", "/v1/programs/1/users/2", nil)

	resp, _ := app.Test(req)
	assert.Nil(t, mock.ExpectationsWereMet())
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
}
