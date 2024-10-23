package middleware

import (
	"github.com/brantem/scorecard/middleware"
	"github.com/gofiber/fiber/v2"
)

type Middleware struct{}

func New() middleware.MiddlewareInterface {
	return &Middleware{}
}

func (m *Middleware) Program(c *fiber.Ctx) error {
	return c.Next()
}

func (m *Middleware) User(c *fiber.Ctx) error {
	return c.Next()
}

func (m *Middleware) SyllabusStructure(c *fiber.Ctx) error {
	return c.Next()
}

func (m *Middleware) Syllabus(c *fiber.Ctx) error {
	return c.Next()
}

func (m *Middleware) ScorecardStructure(c *fiber.Ctx) error {
	return c.Next()
}
