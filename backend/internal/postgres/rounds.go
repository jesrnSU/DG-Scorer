package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"
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
	date, err := time.Parse(time.RFC3339Nano, round.Date)
	if err != nil {
		return 0, fmt.Errorf("parse round date: %w", err)
	}

	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	userID, err := findUserID(ctx, tx, round.Username)
	if err != nil {
		return 0, err
	}

	courseID, err := findOrCreateCourse(ctx, tx, round.CourseName)
	if err != nil {
		return 0, err
	}

	layoutName := round.LayoutName
	if strings.TrimSpace(layoutName) == "" {
		layoutName = "Default"
	}
	if err := ensureCourseLayout(ctx, tx, courseID, layoutName, round.Holes); err != nil {
		return 0, err
	}

	var roundID int64
	err = tx.QueryRow(ctx, `
		INSERT INTO "round" (date, user_id, course_id, layout_name)
		VALUES ($1, $2, $3, $4)
		RETURNING round_id
	`, date, userID, courseID, layoutName).Scan(&roundID)
	if err != nil {
		return 0, fmt.Errorf("insert round: %w", err)
	}

	if err := r.insertHoleScores(ctx, tx, roundID, round.Holes); err != nil {
		return 0, err
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("commit round: %w", err)
	}

	return roundID, nil
}

func (r *RoundRepository) GetUserRounds(
	ctx context.Context,
	username string,
) ([]rounds.Round, error) {
	rows, err := r.pool.Query(ctx, `
        SELECT r.round_id, u.username, c.name, r.layout_name, r.date,
               h.hole_number, h.par, h.strokes
        FROM "round" r
        JOIN users u ON u.id = r.user_id
        JOIN course c ON c.course_id = r.course_id
        LEFT JOIN round_hole_score h ON h.round_id = r.round_id
        WHERE LOWER(u.username) = LOWER($1)
        ORDER BY r.date DESC, r.round_id DESC, h.hole_number
    `, username)
	if err != nil {
		return nil, fmt.Errorf("get user rounds: %w", err)
	}
	defer rows.Close()

	userRounds := make([]rounds.Round, 0)
	for rows.Next() {
		var round rounds.Round
		var date time.Time
		var holeNumber, par, strokes *int
		if err := rows.Scan(&round.ID, &round.Username, &round.CourseName, &round.LayoutName, &date,
			&holeNumber, &par, &strokes); err != nil {
			return nil, fmt.Errorf("scan user round: %w", err)
		}

		// Joined rows repeat the round fields for each hole.
		last := len(userRounds) - 1
		if last < 0 || userRounds[last].ID != round.ID {
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

func (r *RoundRepository) GetAllUsers(
	ctx context.Context,
) ([]rounds.User, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, username FROM users ORDER BY username, id`)
	if err != nil {
		return nil, fmt.Errorf("get all users: %w", err)
	}
	defer rows.Close()
	users := make([]rounds.User, 0)
	for rows.Next() {
		var user rounds.User
		if err := rows.Scan(&user.ID, &user.Username); err != nil {
			return nil, fmt.Errorf("scan user: %w", err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("read users: %w", err)
	}
	return users, nil
}

func findUserID(ctx context.Context, tx pgx.Tx, username string) (int64, error) {
	var userID int64
	err := tx.QueryRow(ctx, `SELECT id FROM users WHERE LOWER(username) = LOWER($1)`, username).Scan(&userID)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, rounds.ErrUsernameNotFound
	}
	if err != nil {
		return 0, fmt.Errorf("find round user: %w", err)
	}
	return userID, nil
}

func findOrCreateCourse(ctx context.Context, tx pgx.Tx, courseName string) (int64, error) {
	var courseID int64
	err := tx.QueryRow(ctx, `
		INSERT INTO course (name)
		VALUES ($1)
		ON CONFLICT (name) DO NOTHING
		RETURNING course_id
	`, courseName).Scan(&courseID)
	if errors.Is(err, pgx.ErrNoRows) {
		// A separate statement can see a conflicting row committed during the insert.
		err = tx.QueryRow(ctx, `SELECT course_id FROM course WHERE name = $1`, courseName).Scan(&courseID)
	}
	if err != nil {
		return 0, fmt.Errorf("find or create course: %w", err)
	}
	return courseID, nil
}

func ensureCourseLayout(ctx context.Context, tx pgx.Tx, courseID int64, layoutName string, holes []rounds.HoleScore) error {
	layoutPar := 0
	for _, hole := range holes {
		layoutPar += hole.Par
	}

	result, err := tx.Exec(ctx, `
		INSERT INTO course_layout (layout_name, course_id, par)
		VALUES ($1, $2, $3)
		ON CONFLICT (layout_name, course_id) DO NOTHING
	`, layoutName, courseID, layoutPar)
	if err != nil {
		return fmt.Errorf("create course layout: %w", err)
	}
	if result.RowsAffected() == 0 {
		return nil
	}

	for _, hole := range holes {
		_, err := tx.Exec(ctx, `
			INSERT INTO course_layout_holes (layout_name, course_id, hole_nr, hole_par)
			VALUES ($1, $2, $3, $4)
		`, layoutName, courseID, hole.HoleNumber, hole.Par)
		if err != nil {
			return fmt.Errorf("insert layout hole %d: %w", hole.HoleNumber, err)
		}
	}
	return nil
}

func (r *RoundRepository) insertHoleScores(
	ctx context.Context,
	tx pgx.Tx,
	roundID int64,
	holes []rounds.HoleScore,
) error {
	for _, hole := range holes {
		_, err := tx.Exec(ctx, `
			INSERT INTO round_hole_score
				(round_id, hole_number, par, strokes)
			VALUES ($1, $2, $3, $4)
		`, roundID, hole.HoleNumber, hole.Par, hole.Strokes)
		if err != nil {
			return fmt.Errorf("insert hole %d: %w", hole.HoleNumber, err)
		}
	}
	return nil
}
