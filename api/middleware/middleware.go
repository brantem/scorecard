package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type MiddlewareInterface interface {
	Program(c *fiber.Ctx) error
	User(c *fiber.Ctx) error
	SyllabusStructure(c *fiber.Ctx) error
	Syllabus(c *fiber.Ctx) error
	ScorecardStructure(c *fiber.Ctx) error
	Scorecard(c *fiber.Ctx) error
}

type Middleware struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) MiddlewareInterface {
	return &Middleware{db}
}
