package scorecard

import (
	"context"

	"github.com/brantem/scorecard/scorecard"
)

type Generator struct {
	EnqueueProgramID   []int
	EnqueueUserID      []int
	EnqueueScorecardID []int
}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) Start() {}

func (g *Generator) Stats() *scorecard.GeneratorStats {
	return &scorecard.GeneratorStats{}
}

func (g *Generator) IsInQueue(scorecardID int) bool {
	return false
}

func (g *Generator) Enqueue(ctx context.Context, programID, userID, scorecardID int) {
	g.EnqueueProgramID = append(g.EnqueueProgramID, programID)
	g.EnqueueUserID = append(g.EnqueueUserID, userID)
	g.EnqueueScorecardID = append(g.EnqueueScorecardID, scorecardID)
}
