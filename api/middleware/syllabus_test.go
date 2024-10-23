package middleware

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/brantem/scorecard/testutil/db"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestSyllabusStructure(t *testing.T) {
	assert := assert.New(t)

	t.Run("structureId < 0", func(t *testing.T) {
		m := Middleware{}

		app := fiber.New()
		app.Get("/:structureId", m.SyllabusStructure, func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})

		req := httptest.NewRequest("GET", "/-1", nil)

		resp, _ := app.Test(req, -1)
		assert.Equal(fiber.StatusNotFound, resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		assert.Equal(`{"error":{"code":"NOT_FOUND"}}`, string(body))
	})

	t.Run("structureId == 0", func(t *testing.T) {
		m := Middleware{}

		app := fiber.New()
		app.Get("/:structureId", m.SyllabusStructure, func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})

		req := httptest.NewRequest("GET", "/0", nil)

		resp, _ := app.Test(req, -1)
		assert.Equal(fiber.StatusOK, resp.StatusCode)
	})

	t.Run("not found", func(t *testing.T) {
		db, mock := db.New()
		m := Middleware{db}

		mock.ExpectQuery("SELECT .+ FROM syllabus_structures").
			WithArgs(1, "1").
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		app := fiber.New()
		app.Get("/:programId/:structureId", m.SyllabusStructure, func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})

		req := httptest.NewRequest("GET", "/1/1", nil)

		resp, _ := app.Test(req, -1)
		assert.Equal(fiber.StatusNotFound, resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		assert.Equal(`{"error":{"code":"NOT_FOUND"}}`, string(body))
	})

	t.Run("success", func(t *testing.T) {
		db, mock := db.New()
		m := Middleware{db}

		mock.ExpectQuery("SELECT .+ FROM syllabus_structures").
			WithArgs(1, "1").
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		app := fiber.New()
		app.Get("/:programId/:structureId", m.SyllabusStructure, func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})

		req := httptest.NewRequest("GET", "/1/1", nil)

		resp, _ := app.Test(req, -1)
		assert.Equal(fiber.StatusOK, resp.StatusCode)
	})
}

func TestSyllabus(t *testing.T) {
	assert := assert.New(t)

	t.Run("syllabusId < 0", func(t *testing.T) {
		m := Middleware{}

		app := fiber.New()
		app.Get("/:syllabusId", m.Syllabus, func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})

		req := httptest.NewRequest("GET", "/-1", nil)

		resp, _ := app.Test(req, -1)
		assert.Equal(fiber.StatusNotFound, resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		assert.Equal(`{"error":{"code":"NOT_FOUND"}}`, string(body))
	})

	t.Run("syllabusId == 0", func(t *testing.T) {
		m := Middleware{}

		app := fiber.New()
		app.Get("/:syllabusId", m.Syllabus, func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})

		req := httptest.NewRequest("GET", "/0", nil)

		resp, _ := app.Test(req, -1)
		assert.Equal(fiber.StatusOK, resp.StatusCode)
	})

	t.Run("not found", func(t *testing.T) {
		db, mock := db.New()
		m := Middleware{db}

		mock.ExpectQuery("SELECT .+ FROM syllabuses").
			WithArgs(1, "1").
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		app := fiber.New()
		app.Get("/:programId/:syllabusId", m.Syllabus, func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})

		req := httptest.NewRequest("GET", "/1/1", nil)

		resp, _ := app.Test(req, -1)
		assert.Equal(fiber.StatusNotFound, resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		assert.Equal(`{"error":{"code":"NOT_FOUND"}}`, string(body))
	})

	t.Run("success", func(t *testing.T) {
		db, mock := db.New()
		m := Middleware{db}

		mock.ExpectQuery("SELECT .+ FROM syllabuses").
			WithArgs(1, "1").
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		app := fiber.New()
		app.Get("/:programId/:syllabusId", m.Syllabus, func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusOK)
		})

		req := httptest.NewRequest("GET", "/1/1", nil)

		resp, _ := app.Test(req, -1)
		assert.Equal(fiber.StatusOK, resp.StatusCode)
	})
}
