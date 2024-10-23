package scorecard

import (
	"context"
	"sync"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"github.com/vmihailenco/taskq/v3"
	"github.com/vmihailenco/taskq/v3/memqueue"
)

type GeneratorInterface interface {
	Start()
	Enqueue(ctx context.Context, programID, userID, scorecardID int)
}

type Generator struct {
	db *sqlx.DB

	queue taskq.Queue
	task  *taskq.Task
}

func NewGenerator(db *sqlx.DB) GeneratorInterface {
	return &Generator{db: db}
}

func (g *Generator) Start() {
	factory := memqueue.NewFactory()

	g.queue = factory.RegisterQueue(&taskq.QueueOptions{
		Name:         "scorecard",
		MaxNumWorker: 1,
	})

	g.task = taskq.RegisterTask(&taskq.TaskOptions{
		Name:    "generate",
		Handler: g.generate,
	})
}

func (g *Generator) Enqueue(ctx context.Context, programID, userID, scorecardID int) {
	msg := g.task.WithArgs(ctx, programID, userID, scorecardID)
	if err := g.queue.Add(msg); err != nil {
		log.Error().Err(err).Msg("scorecard.Generator.Enqueue")
	}
}

func (g *Generator) generate(programID, userID, scorecardID int) error {
	type ScorecardStructure struct {
		ID         int
		ParentID   *int `db:"parent_id"`
		SyllabusID *int `db:"syllabus_id"`
	}

	var rootIds []int
	var structures []*ScorecardStructure
	assignments := make(map[int]float64)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		rows, err := g.db.Queryx(`
			SELECT id, parent_id, syllabus_id
			FROM scorecard_structures
			WHERE program_id = ?
		`, programID)
		if err != nil {
			log.Error().Err(err).Msg("scorecard.Generator.generate")
			return
		}
		defer rows.Close()

		for rows.Next() {
			var node ScorecardStructure
			if err := rows.StructScan(&node); err != nil {
				log.Error().Err(err).Msg("scorecard.Generator.generate")
				continue
			}
			if node.ParentID == nil {
				rootIds = append(rootIds, node.ID)
			}
			structures = append(structures, &node)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		rows, err := g.db.Query(`
			SELECT syllabus_id, score
			FROM user_scores
			WHERE user_id = ?
		`, userID)
		if err != nil {
			log.Error().Err(err).Msg("scorecard.Generator.generate")
			return
		}
		defer rows.Close()

		for rows.Next() {
			var syllabusID int
			var score float64
			if err := rows.Scan(&syllabusID, &score); err != nil {
				log.Error().Err(err).Msg("scorecard.Generator.generate")
				continue
			}
			assignments[syllabusID] = score
		}
	}()

	wg.Wait()

	if len(structures) == 0 || len(assignments) == 0 {
		return nil
	}

	nodes := make([]*Node, len(structures))
	for i, structure := range structures {
		var score float64
		if structure.SyllabusID != nil {
			score = assignments[*structure.SyllabusID]
		}

		nodes[i] = &Node{
			ID:       structure.ID,
			ParentID: structure.ParentID,
			Score:    score,
		}
	}

	reducer := NewReducer()
	reducer.SetNodes(nodes)
	reducer.Reduce()

	tx := g.db.MustBegin()

	var score float64

	roots := reducer.GetRoots()
	for _, node := range roots {
		score += node.Score
	}
	score /= float64(len(roots))

	if scorecardID == 0 {
		err := tx.QueryRow(`
			INSERT INTO scorecards (program_id, user_id, score)
			VALUES (?, ?, ?)
			RETURNING id
		`, programID, userID, score).Scan(&scorecardID)
		if err != nil {
			tx.Rollback()
			log.Error().Err(err).Msg("scorecard.Generator.generate")
			return nil
		}
	} else {
		_, err := tx.Exec(`
			UPDATE scorecards
			SET score = ?, is_outdated = FALSE
			WHERE id = ?
		`, score, scorecardID)
		if err != nil {
			tx.Rollback()
			log.Error().Err(err).Msg("scorecard.Generator.generate")
			return nil
		}
	}

	qb := sq.Insert("scorecard_items").Columns("scorecard_id", "structure_id", "score").
		Suffix("ON CONFLICT (scorecard_id, structure_id) DO UPDATE SET score = EXCLUDED.score")
	for _, rootID := range rootIds {
		for _, node := range reducer.GetByParentID(rootID) {
			qb = qb.Values(scorecardID, node.ID, node.Score)
		}
	}

	if _, err := qb.RunWith(tx).Exec(); err != nil {
		tx.Rollback()
		log.Error().Err(err).Msg("scorecard.Generator.generate")
		return nil
	}

	tx.Commit()

	return nil
}
