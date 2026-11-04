package rounds

import "testing"

func TestValidateRound(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*Round)
		want   string
	}{
		{"valid", func(r *Round) {}, ""},
		{"missing username", func(r *Round) { r.Username = "" }, "username is required"},
		{"blank username", func(r *Round) { r.Username = "  " }, "username is required"},
		{"date without time", func(r *Round) { r.Date = "2026-10-16" }, "date must be a valid RFC 3339 timestamp with a timezone"},
		{"time without timezone", func(r *Round) { r.Date = "2026-10-16T15:42:00" }, "date must be a valid RFC 3339 timestamp with a timezone"},
		{"UTC timestamp", func(r *Round) { r.Date = "2026-10-16T13:42:00Z" }, ""},
		{"blank course", func(r *Round) { r.CourseName = " \t\n" }, "course_name is required"},
		{"missing date", func(r *Round) { r.Date = "" }, "date must be a valid RFC 3339 timestamp with a timezone"},
		{"wrong date format", func(r *Round) { r.Date = "16/10/2026" }, "date must be a valid RFC 3339 timestamp with a timezone"},
		{"impossible date", func(r *Round) { r.Date = "2026-02-30" }, "date must be a valid RFC 3339 timestamp with a timezone"},
		{"leap day", func(r *Round) { r.Date = "2024-02-29T15:42:00+02:00" }, ""},
		{"missing holes", func(r *Round) { r.Holes = nil }, "round must contain at least one hole"},
		{"empty holes", func(r *Round) { r.Holes = []HoleScore{} }, "round must contain at least one hole"},
		{"zero hole", func(r *Round) { r.Holes[0].HoleNumber = 0 }, "hole_number must be positive"},
		{"negative hole", func(r *Round) { r.Holes[0].HoleNumber = -1 }, "hole_number must be positive"},
		{"duplicate hole", func(r *Round) { r.Holes = append(r.Holes, r.Holes[0]) }, "hole 1 appears more than once"},
		{"zero par", func(r *Round) { r.Holes[0].Par = 0 }, "hole 1: par must be positive"},
		{"negative par", func(r *Round) { r.Holes[0].Par = -1 }, "hole 1: par must be positive"},
		{"zero strokes", func(r *Round) { r.Holes[0].Strokes = 0 }, "hole 1: strokes must be positive"},
		{"negative strokes", func(r *Round) { r.Holes[0].Strokes = -1 }, "hole 1: strokes must be positive"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			round := Round{Username: "JesperNi", CourseName: "Practice course", Date: "2026-10-16T15:42:00+02:00", Holes: []HoleScore{{HoleNumber: 1, Par: 3, Strokes: 4}}}
			tt.mutate(&round)
			err := ValidateRound(round)
			if tt.want == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			} else if err == nil || err.Error() != tt.want {
				t.Fatalf("error = %v, want %q", err, tt.want)
			}
		})
	}
}
