package postgres

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jesrnsu/jesrnsu.github.io/internal/rounds"
)

func TestRoundUsernameAndTimestampPersistence(t *testing.T) {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("set TEST_DATABASE_URL to run PostgreSQL integration tests")
	}
	ctx := context.Background()
	admin, err := pgxpool.New(ctx, url)
	if err != nil {
		t.Fatal(err)
	}
	defer admin.Close()
	schema := fmt.Sprintf("round_test_%d", time.Now().UnixNano())
	if _, err := admin.Exec(ctx, "CREATE SCHEMA "+schema); err != nil {
		t.Fatal(err)
	}
	defer admin.Exec(ctx, "DROP SCHEMA "+schema+" CASCADE")
	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		t.Fatal(err)
	}
	config.ConnConfig.RuntimeParams["search_path"] = schema
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()
	apply := func(file string) {
		t.Helper()
		sql, err := os.ReadFile("../../migrations/" + file)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := pool.Exec(ctx, string(sql)); err != nil {
			t.Fatal(err)
		}
	}
	apply("001_create_rounds.sql")
	if _, err := pool.Exec(ctx, "INSERT INTO rounds (course_name, date) VALUES ('Legacy', '2026-10-10')"); err != nil {
		t.Fatal(err)
	}
	apply("002_round_username_timestamp.sql")
	repo := NewRoundRepository(pool)
	legacy, err := repo.GetRound(ctx, 1)
	if err != nil {
		t.Fatal(err)
	}
	if legacy.Username != "unknown" || legacy.Date != "2026-10-10T00:00:00Z" {
		t.Fatalf("legacy round = %+v", legacy)
	}

	var firstID int64
	for _, timestamp := range []string{"2026-10-10T15:42:00+02:00", "2026-10-10T17:42:00+02:00"} {
		input := rounds.Round{Username: "JesperNi", CourseName: "Visättra", Date: timestamp,
			Holes: []rounds.HoleScore{{HoleNumber: 1, Par: 3, Strokes: 4}}}
		id, err := rounds.NewService(repo).SaveRound(ctx, input)
		if err != nil {
			t.Fatal(err)
		}
		if id == firstID {
			t.Fatal("different rounds must have different IDs")
		}
		firstID = id
		got, err := repo.GetRound(ctx, id)
		if err != nil {
			t.Fatal(err)
		}
		expectedTime, _ := time.Parse(time.RFC3339, timestamp)
		actualTime, err := time.Parse(time.RFC3339, got.Date)
		if err != nil || !actualTime.Equal(expectedTime) || got.Username != input.Username || len(got.Holes) != 1 {
			t.Fatalf("round did not retain username, timestamp and holes: %+v", got)
		}
	}
}
