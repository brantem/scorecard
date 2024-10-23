package scorecard

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/brantem/scorecard/testutil/db"
	"github.com/stretchr/testify/assert"
)

func TestQueue_generate(t *testing.T) {
	assert := assert.New(t)

	programID := 1
	userID := 2

	score := float64(100)

	t.Run("empty structures", func(t *testing.T) {
		db, mock := db.New()
		q := NewQueue(db)

		mock.MatchExpectationsInOrder(false)

		mock.ExpectQuery("SELECT .+ FROM scorecard_structures").
			WithArgs(programID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "parent_id", "syllabus_id"}))

		mock.ExpectQuery("SELECT .+ FROM user_scores").
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"syllabus_id", "score"}).AddRow(1, 100))

		assert.Nil(q.generate(programID, userID, 0))
		assert.Nil(mock.ExpectationsWereMet())
	})

	t.Run("empty assignments", func(t *testing.T) {
		db, mock := db.New()
		q := NewQueue(db)

		mock.MatchExpectationsInOrder(false)

		mock.ExpectQuery("SELECT .+ FROM scorecard_structures").
			WithArgs(programID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "parent_id", "syllabus_id"}).AddRow(1, nil, 1))

		mock.ExpectQuery("SELECT .+ FROM user_scores").
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"syllabus_id", "score"}))

		assert.Nil(q.generate(programID, userID, 0))
		assert.Nil(mock.ExpectationsWereMet())
	})

	t.Run("insert", func(t *testing.T) {
		db, mock := db.New()
		q := NewQueue(db)

		mock.MatchExpectationsInOrder(false)

		mock.ExpectQuery("SELECT .+ FROM scorecard_structures").
			WithArgs(programID).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "parent_id", "syllabus_id"}).
					AddRow(1, nil, 1).
					AddRow(2, 1, 2),
			)

		mock.ExpectQuery("SELECT .+ FROM user_scores").
			WithArgs(userID).
			WillReturnRows(
				sqlmock.NewRows([]string{"syllabus_id", "score"}).
					AddRow(1, 0).
					AddRow(2, 100),
			)

		mock.ExpectBegin()

		scorecardID := 1
		mock.ExpectQuery("INSERT INTO scorecards").
			WithArgs(programID, userID, score).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(scorecardID))

		mock.ExpectExec("INSERT INTO scorecard_items").
			WithArgs(scorecardID, 2, score).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		assert.Nil(q.generate(programID, userID, 0))
		assert.Nil(mock.ExpectationsWereMet())
	})

	t.Run("update", func(t *testing.T) {
		db, mock := db.New()
		q := NewQueue(db)

		mock.MatchExpectationsInOrder(false)

		mock.ExpectQuery("SELECT .+ FROM scorecard_structures").
			WithArgs(programID).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "parent_id", "syllabus_id"}).
					AddRow(1, nil, 1).
					AddRow(2, 1, 2),
			)

		mock.ExpectQuery("SELECT .+ FROM user_scores").
			WithArgs(userID).
			WillReturnRows(
				sqlmock.NewRows([]string{"syllabus_id", "score"}).
					AddRow(1, 0).
					AddRow(2, 100),
			)

		mock.ExpectBegin()

		scorecardID := 1
		mock.ExpectExec("UPDATE scorecards").
			WithArgs(score, scorecardID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectExec("INSERT INTO scorecard_items").
			WithArgs(scorecardID, 2, score).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectCommit()

		assert.Nil(q.generate(programID, userID, scorecardID))
		assert.Nil(mock.ExpectationsWereMet())
	})

}
