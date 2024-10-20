package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type Handler struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Handler {
	return &Handler{db}
}

func (h *Handler) Register(r *fiber.App) {
	v1 := r.Group("/v1")

	users := v1.Group("/users")
	users.Get("/", h.users)
	users.Put("/:userId<int>?", h.saveUser)

	syllabuses := v1.Group("/syllabuses")
	{
		structures := syllabuses.Group("/structures")
		structures.Get("/", h.syllabusStructures)
		structures.Put("/:structureId<int>?", h.saveSyllabusStructure)
		structures.Delete("/:structureId", h.deleteSyllabusStructure)

		syllabuses.Get("/", h.syllabuses)
		syllabuses.Put("/:syllabusId<int>?", h.saveSyllabus)
		syllabuses.Delete("/:syllabusId", h.deleteSyllabus)
	}

	scores := syllabuses.Group("/:syllabusId<int>/scores")
	scores.Get("/", h.scores)
	scores.Put("/:userId<int>", h.saveScore)

	structures := v1.Group("/structures")
	structures.Get("/", h.structures)
	structures.Put("/:structureId<int>?", h.saveStructure)
	structures.Delete("/:structureId", h.deleteStructure)
}
