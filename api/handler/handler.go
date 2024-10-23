package handler

import (
	"github.com/brantem/scorecard/middleware"
	"github.com/brantem/scorecard/scorecard"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type Handler struct {
	db *sqlx.DB

	generator scorecard.GeneratorInterface
}

func New(db *sqlx.DB, generator scorecard.GeneratorInterface) *Handler {
	return &Handler{db, generator}
}

func (h *Handler) Register(r *fiber.App, m middleware.MiddlewareInterface) {
	v1 := r.Group("/v1")

	programs := v1.Group("/programs")
	programs.Get("/", h.programs)
	programs.Put("/:programId<int>?", h.saveProgram)

	programID := programs.Group("/:programId<int>", m.Program)
	programID.Get("/", h.program)
	programID.Delete("/", h.deleteProgram)

	users := programID.Group("/users")
	{
		users.Get("/", h.users)
		users.Put("/:userId<int>?", h.saveUser)

		userID := users.Group("/:userId<int>", m.User)
		userID.Get("/", h.user)
		userID.Get("/scores", h.userScores)
		userID.Delete("/", h.deleteUser)
	}

	syllabuses := programID.Group("/syllabuses")
	{
		structures := syllabuses.Group("/structures")
		structures.Get("/", h.syllabusStructures)
		structures.Put("/:structureId<int>?", m.SyllabusStructure, h.saveSyllabusStructure)
		structures.Delete("/:structureId<int>", m.SyllabusStructure, h.deleteSyllabusStructure)

		syllabuses.Get("/", h.syllabuses)
		syllabuses.Put("/:syllabusId<int>?", m.Syllabus, h.saveSyllabus)

		syllabusID := syllabuses.Group("/:syllabusId<int>", m.Syllabus)
		syllabusID.Get("/", h.syllabus)
		syllabusID.Delete("/", h.deleteSyllabus)

		scores := syllabusID.Group("/scores")
		scores.Get("/", h.syllabusScores)
		scores.Put("/:userId<int>", m.User, h.saveScore)
	}

	scorecards := programID.Group("/scorecards")
	{
		structures := scorecards.Group("/structures")
		structures.Get("/", h.scorecardStructures)
		structures.Post("/copy/:syllabusId<int>", m.Syllabus, h.copySyllabusesIntoStructures)
		structures.Put("/:structureId<int>?", m.ScorecardStructure, h.saveScorecardStructure)
		structures.Delete("/:structureId<int>", m.ScorecardStructure, h.deleteScorecardStructure)

		scorecards.Get("/", h.scorecards)
		scorecards.Post("/generate/:scorecardId<int>?", m.Scorecard, h.generateScorecards)
		scorecards.Get("/:scorecardId", m.Scorecard, h.scorecard)
	}
}
