package rounds

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type RoundStore interface {
	SaveRound(ctx context.Context, round Round) (int64, error)
	GetRound(ctx context.Context, id int64) (Round, error)
}

type Service struct {
	store RoundStore
}

func NewService(store RoundStore) *Service {
	return &Service{store: store}
}

func (s *Service) SaveRound(ctx context.Context, round Round) (int64, error) {
	if err := ValidateRound(round); err != nil {
		return 0, err
	}

	id, err := s.store.SaveRound(ctx, round)
	if err != nil {
		return 0, fmt.Errorf("save round: %w", err)
	}

	return id, nil
}

func (s *Service) GetRound(ctx context.Context, id int64) (Round, error) {
	round, err := s.store.GetRound(ctx, id)
	if err != nil {
		return Round{}, fmt.Errorf("get round: %w", err)
	}

	return round, nil
}

func ValidateRound(round Round) error {
	if strings.TrimSpace(round.Username) == "" {
		return fmt.Errorf("username is required")
	}
	if strings.TrimSpace(round.CourseName) == "" {
		return fmt.Errorf("course_name is required")
	}
	if _, err := time.Parse(time.RFC3339Nano, round.Date); err != nil {
		return fmt.Errorf("date must be a valid RFC 3339 timestamp with a timezone")
	}
	if len(round.Holes) == 0 {
		return fmt.Errorf("round must contain at least one hole")
	}
	seen := make(map[int]bool)
	for _, hole := range round.Holes {
		if hole.HoleNumber <= 0 {
			return fmt.Errorf("hole_number must be positive")
		}
		if seen[hole.HoleNumber] {
			return fmt.Errorf("hole %d appears more than once", hole.HoleNumber)
		}
		seen[hole.HoleNumber] = true
		if hole.Par <= 0 {
			return fmt.Errorf("hole %d: par must be positive", hole.HoleNumber)
		}
		if hole.Strokes <= 0 {
			return fmt.Errorf("hole %d: strokes must be positive", hole.HoleNumber)
		}
	}
	return nil
}

func TotalStrokes(round Round) int {
	total := 0

	for _, hole := range round.Holes {
		total += hole.Strokes
	}

	return total
}

func TotalPar(round Round) int {
	total := 0

	for _, hole := range round.Holes {
		total += hole.Par
	}

	return total
}

func ScoreRelativeToPar(round Round) int {
	return TotalStrokes(round) - TotalPar(round)
}

func ConvertToStringDate(round Round) string {
	return "not implemented"
}
