package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jesrnsu/jesrnsu.github.io/internal/rounds"
)

// RoundRepository will store and retrieve rounds using the shared pool.
// It does not own the pool; application startup and shutdown manage its lifetime.
type RoundRepository struct {
	pool *pgxpool.Pool
}

func NewRoundRepository(pool *pgxpool.Pool) *RoundRepository {
	return &RoundRepository{pool: pool}
}

func (r *RoundRepository) SaveRound(
	ctx context.Context,
	round rounds.Round,
) (int64, error) {
	date, err := time.Parse("2006-01-02", round.Date)
	if err != nil {
		return 0, fmt.Errorf("parse round date: %w", err)
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	roundID, err := r.insertRound(ctx, tx, round.CourseName, date)
	if err != nil {
		return 0, err
	}

	if err := r.insertHoleScores(ctx, tx, roundID, round.Holes); err != nil {
		return 0, err
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("commit round: %w", err)
	}

	return roundID, nil
}

func (r *RoundRepository) GetRound(
	ctx context.Context,
	id int64,
) (rounds.Round, error) {
	var round rounds.Round
	var date time.Time

	err := r.pool.QueryRow(ctx, `
		SELECT id, course_name, date
		FROM rounds
		WHERE id = $1
	`, id).Scan(&round.Id, &round.CourseName, &date)
	if errors.Is(err, pgx.ErrNoRows) {
		return rounds.Round{}, rounds.ErrRoundNotFound
	}
	if err != nil {
		return rounds.Round{}, fmt.Errorf("get round: %w", err)
	}

	round.Date = date.Format("2006-01-02")
	round.Holes = make([]rounds.HoleScore, 0)

	rows, err := r.pool.Query(ctx, `
		SELECT hole_number, par, strokes
		FROM hole_scores
		WHERE round_id = $1
		ORDER BY hole_number
	`, id)
	if err != nil {
		return rounds.Round{}, fmt.Errorf("get hole scores: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var hole rounds.HoleScore
		if err := rows.Scan(
			&hole.HoleNumber,
			&hole.Par,
			&hole.Strokes,
		); err != nil {
			return rounds.Round{}, fmt.Errorf("scan hole score: %w", err)
		}

		round.Holes = append(round.Holes, hole)
	}
	if err := rows.Err(); err != nil {
		return rounds.Round{}, fmt.Errorf("read hole scores: %w", err)
	}

	return round, nil
}

func (r *RoundRepository) insertRound(
	ctx context.Context,
	tx pgx.Tx,
	courseName string,
	date time.Time,
) (int64, error) {
	var roundID int64

	err := tx.QueryRow(ctx, `
		INSERT INTO rounds (course_name, date)
		VALUES ($1, $2)
		RETURNING id
	`, courseName, date).Scan(&roundID)
	if err != nil {
		return 0, fmt.Errorf("insert round: %w", err)
	}

	return roundID, nil
}

func (r *RoundRepository) insertHoleScores(
	ctx context.Context,
	tx pgx.Tx,
	roundID int64,
	holes []rounds.HoleScore,
) error {
	for _, hole := range holes {
		_, err := tx.Exec(ctx, `
			INSERT INTO hole_scores
				(round_id, hole_number, par, strokes)
			VALUES ($1, $2, $3, $4)
		`, roundID, hole.HoleNumber, hole.Par, hole.Strokes)
		if err != nil {
			return fmt.Errorf("insert hole %d: %w", hole.HoleNumber, err)
		}
	}

	return nil
}
