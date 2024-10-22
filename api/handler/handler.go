package handler

import (
	"github.com/brantem/scorecard/scorecard"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type Handler struct {
	db *sqlx.DB

	generator *scorecard.Queue
}

func New(db *sqlx.DB, generator *scorecard.Queue) *Handler {
	return &Handler{db, generator}
}

func (h *Handler) Register(r *fiber.App) {
	v1 := r.Group("/v1")

	users := v1.Group("/users")
	users.Get("/", h.users)
	users.Put("/:userId<int>?", h.saveUser)

	userID := users.Group("/:userId<int>")
	userID.Get("/", h.user)
	userID.Get("/scores", h.userScores)
	userID.Delete("/", h.deleteUser)

	syllabuses := v1.Group("/syllabuses")
	{
		structures := syllabuses.Group("/structures")
		structures.Get("/", h.syllabusStructures)
		structures.Put("/:structureId<int>?", h.saveSyllabusStructure)
		structures.Delete("/:structureId", h.deleteSyllabusStructure)

		syllabuses.Get("/", h.syllabuses)
		syllabuses.Put("/:syllabusId<int>?", h.saveSyllabus)
		syllabuses.Delete("/:syllabusId", h.deleteSyllabus)

		syllabusID := syllabuses.Group("/:syllabusId<int>")
		syllabusID.Get("/", h.syllabus)

		scores := syllabusID.Group("/scores")
		scores.Get("/", h.syllabusScores)
		scores.Put("/:userId<int>", h.saveScore)
	}

	scorecards := v1.Group("/scorecards")
	{
		structures := scorecards.Group("/structures")
		structures.Get("/", h.scorecardStructures)
		structures.Post("/copy/:syllabusId<int>", h.copySyllabusesIntoStructures)
		structures.Put("/:structureId<int>?", h.saveScorecardStructure)
		structures.Delete("/:structureId", h.deleteScorecardStructure)

		scorecards.Get("/", h.scorecards)
		scorecards.Post("/generate/:scorecardId<int>?", h.generateScorecards)
		scorecards.Get("/:scorecardId<int>", h.scorecard)
	}
}
