package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"
	"strings"

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
	date, err := time.Parse(time.RFC3339Nano, round.Date)
	if err != nil {
		return 0, fmt.Errorf("parse round date: %w", err)
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	roundID, err := r.insertRound(ctx, tx, round.Username, round.CourseName, date)
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
		SELECT id, username, course_name, date
		FROM rounds
		WHERE id = $1
	`, id).Scan(&round.Id, &round.Username, &round.CourseName, &date)
	if errors.Is(err, pgx.ErrNoRows) {
		return rounds.Round{}, rounds.ErrRoundNotFound
	}
	if err != nil {
		return rounds.Round{}, fmt.Errorf("get round: %w", err)
	}

	round.Date = date.UTC().Format(time.RFC3339Nano)
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

func (r *RoundRepository) GetUserRounds(
	ctx context.Context,
	username string,
) ([]rounds.Round, error) {
	var usernameLowercase := strings.toLower(username)

	rows, err := r.pool.Query(ctx, `
        SELECT r.id, r.username, r.course_name, r.date,
               h.hole_number, h.par, h.strokes
        FROM rounds r
        LEFT JOIN hole_scores h ON h.round_id = r.id
        WHERE r.username = $1
        ORDER BY r.date DESC, r.id DESC, h.hole_number
    `, usernameLowercase)
	if err != nil {
		return nil, fmt.Errorf("get user rounds: %w", err)
	}
	defer rows.Close()

	userRounds := make([]rounds.Round, 0)
	for rows.Next() {
		var round rounds.Round
		var date time.Time
		var holeNumber, par, strokes *int
		if err := rows.Scan(&round.Id, &round.Username, &round.CourseName, &date,
			&holeNumber, &par, &strokes); err != nil {
			return nil, fmt.Errorf("scan user round: %w", err)
		}

		// Joined rows repeat the round fields for each hole.
		last := len(userRounds) - 1
		if last < 0 || userRounds[last].Id != round.Id {
			round.Date = date.UTC().Format(time.RFC3339Nano)
			round.Holes = make([]rounds.HoleScore, 0)
			userRounds = append(userRounds, round)
			last++
		}
		if holeNumber != nil {
			userRounds[last].Holes = append(userRounds[last].Holes, rounds.HoleScore{
				HoleNumber: *holeNumber, Par: *par, Strokes: *strokes,
			})
		}
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read user rounds: %w", err)
	}
	return userRounds, nil
}

func (r *RoundRepository) insertRound(
	ctx context.Context,
	tx pgx.Tx,
	username string,
	courseName string,
	date time.Time,
) (int64, error) {
	var roundID int64
	var usernameLowercase := strings.toLower(username)
	var courseNameLowercase := strings.toLower(courseName)

	err := tx.QueryRow(ctx, `
		INSERT INTO rounds (username, course_name, date)
		VALUES ($1, $2, $3)
		RETURNING id
	`, usernameLowercase, courseNameLowercase, date).Scan(&roundID)
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
