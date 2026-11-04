package rounds

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

type roundStoreFunc func(context.Context, Round) (int64, error)

func (f roundStoreFunc) SaveRound(ctx context.Context, round Round) (int64, error) {
	return f(ctx, round)
}

func (f roundStoreFunc) GetRound(context.Context, int64) (Round, error) {
	return Round{}, errors.New("unexpected GetRound call")
}

func TestSaveRound(t *testing.T) {
	round := Round{
		Username:   "JesperNi",
		Id:         7, // A supplied ID must not replace the database-generated ID.
		CourseName: "Visättra",
		Date:       "2026-10-19T15:42:00+02:00",
		Holes:      []HoleScore{{HoleNumber: 1, Par: 3, Strokes: 2}},
	}

	t.Run("passes round and context to storage", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		calls := 0
		service := NewService(roundStoreFunc(func(gotCtx context.Context, got Round) (int64, error) {
			calls++
			if gotCtx != ctx {
				t.Error("storage received a different context")
			}
			if !reflect.DeepEqual(got, round) {
				t.Errorf("storage received %+v, want %+v", got, round)
			}
			return 42, nil
		}))
		id, err := service.SaveRound(ctx, round)
		if err != nil {
			t.Fatal(err)
		}
		if id != 42 {
			t.Fatalf("saved ID = %d, want database ID 42", id)
		}
		if calls != 1 {
			t.Fatalf("storage called %d times, want 1", calls)
		}
	})

	t.Run("rejects invalid round before storage", func(t *testing.T) {
		service := NewService(roundStoreFunc(func(context.Context, Round) (int64, error) {
			t.Fatal("storage called for invalid round")
			return 42, nil
		}))
		if id, err := service.SaveRound(context.Background(), Round{}); err == nil || id != 0 {
			t.Fatal("expected validation error")
		}
	})

	t.Run("preserves storage errors", func(t *testing.T) {
		storageErr := errors.New("database unavailable")
		service := NewService(roundStoreFunc(func(context.Context, Round) (int64, error) {
			return 99, storageErr
		}))
		id, err := service.SaveRound(context.Background(), round)
		if id != 0 {
			t.Fatalf("failed save returned ID %d, want 0", id)
		}
		if !errors.Is(err, storageErr) {
			t.Fatalf("SaveRound() error = %v, want wrapped storage error", err)
		}
		if err.Error() != "save round: database unavailable" {
			t.Fatalf("unexpected error message: %v", err)
		}
	})
}

// A domain test runs the calculation without starting HTTP or PostgreSQL.
func TestTotalStrokes(t *testing.T) {
	round := Round{
		Username:   "JesperNi",
		CourseName: "Visättra",
		Holes: []HoleScore{
			{HoleNumber: 1, Par: 3, Strokes: 2},
			{HoleNumber: 2, Par: 3, Strokes: 3},
		},
	}
	if got := TotalStrokes(round); got != 5 {
		t.Fatalf("TotalStrokes() = %d, want 5", got)
	}
}

func TestTotalPar(t *testing.T) {
	round := Round{
		Username:   "JesperNi",
		CourseName: "Rudan",
		Holes: []HoleScore{
			{HoleNumber: 1, Par: 4, Strokes: 3},
			{HoleNumber: 2, Par: 4, Strokes: 4},
		},
	}
	if got := TotalPar(round); got != 8 {
		t.Fatalf("TotalPar() = %d, want 8", got)
	}
}

func TestRoundEvenPar(t *testing.T) {
	round := Round{
		Username:   "JesperNi",
		CourseName: "Visättra",
		Holes: []HoleScore{
			{HoleNumber: 1, Par: 3, Strokes: 2},
			{HoleNumber: 2, Par: 3, Strokes: 4},
		},
	}
	if got := ScoreRelativeToPar(round); got != 0 {
		t.Fatalf("ScoreRelativeToPar() = %d, want 0", got)
	}
}

func TestRoundUnderPar(t *testing.T) {
	round := Round{
		Username:   "JesperNi",
		CourseName: "Visättra",
		Holes: []HoleScore{
			{HoleNumber: 1, Par: 3, Strokes: 2},
			{HoleNumber: 2, Par: 3, Strokes: 3},
		},
	}
	if got := ScoreRelativeToPar(round); got != -1 {
		t.Fatalf("ScoreRelativeToPar() = %d, want -1", got)
	}
}

func TestRoundOverPar(t *testing.T) {
	round := Round{
		Username:   "JesperNi",
		CourseName: "Visättra",
		Holes: []HoleScore{
			{HoleNumber: 1, Par: 3, Strokes: 4},
			{HoleNumber: 2, Par: 3, Strokes: 3},
		},
	}
	if got := ScoreRelativeToPar(round); got != 1 {
		t.Fatalf("ScoreRelativeToPar() = %d, want 1", got)
	}
}
