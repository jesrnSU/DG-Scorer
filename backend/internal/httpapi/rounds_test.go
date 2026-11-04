package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/jesrnsu/jesrnsu.github.io/internal/rounds"
)

type roundStoreFunc func(context.Context, rounds.Round) (int64, error)

func (f roundStoreFunc) SaveRound(ctx context.Context, round rounds.Round) (int64, error) {
	return f(ctx, round)
}

func (f roundStoreFunc) GetRound(context.Context, int64) (rounds.Round, error) {
	return rounds.Round{}, errors.New("unexpected GetRound call")
}

func TestSaveRoundEndpoint(t *testing.T) {
	const valid = `{"username":"JesperNi","course_name":"Practice","date":"2026-10-20T15:42:00+02:00","holes":[{"hole_number":1,"par":3,"strokes":4}]}`
	for _, tt := range []struct {
		name        string
		payload     string
		unavailable bool
		storeErr    error
		status      int
		wantCalls   int
	}{
		{name: "saved", payload: valid, status: http.StatusCreated, wantCalls: 1},
		{name: "invalid JSON", payload: `{`, status: http.StatusBadRequest},
		{name: "invalid round", payload: `{}`, status: http.StatusBadRequest},
		{name: "database failure", payload: valid, storeErr: errors.New("private database details"), status: http.StatusInternalServerError, wantCalls: 1},
		{name: "no database configured", payload: valid, unavailable: true, status: http.StatusServiceUnavailable},
	} {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/api/v1/rounds", strings.NewReader(tt.payload))
			request.Header.Set("Content-Type", "application/json")
			calls := 0
			var service *rounds.Service
			if !tt.unavailable {
				service = rounds.NewService(roundStoreFunc(func(ctx context.Context, round rounds.Round) (int64, error) {
					calls++
					if ctx != request.Context() {
						t.Error("request context was not forwarded")
					}
					if round.CourseName != "Practice" || len(round.Holes) != 1 || round.Holes[0].Strokes != 4 {
						t.Errorf("unexpected saved round: %+v", round)
					}
					return 42, tt.storeErr
				}))
			}
			response := httptest.NewRecorder()
			NewRouter(service).ServeHTTP(response, request)
			if response.Code != tt.status {
				t.Fatalf("status = %d, want %d; body = %s", response.Code, tt.status, response.Body.String())
			}
			if response.Code == http.StatusCreated {
				var body struct {
					ID int64 `json:"id"`
				}
				if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
					t.Fatal(err)
				}
				if body.ID != 42 {
					t.Fatalf("response ID = %d, want database ID 42", body.ID)
				}
			}
			if calls != tt.wantCalls {
				t.Errorf("storage calls = %d, want %d", calls, tt.wantCalls)
			}
			if strings.Contains(response.Body.String(), "private database details") {
				t.Error("response exposed storage error")
			}
		})
	}
}

// The example contains multiple rounds; the API accepts one per request.
func exampleRoundPayload(t *testing.T) []byte {
	t.Helper()
	data, err := os.ReadFile("../../examples/round.json")
	if err != nil {
		t.Fatal(err)
	}
	var examples []rounds.Round
	if err := json.Unmarshal(data, &examples); err != nil {
		t.Fatal(err)
	}
	if len(examples) != 2 {
		t.Fatalf("example contains %d rounds, want 2", len(examples))
	}
	for _, round := range examples {
		if err := rounds.ValidateRound(round); err != nil {
			t.Fatalf("invalid example round: %v", err)
		}
		if len(round.Holes) != 18 {
			t.Fatalf("example round contains %d holes, want 18", len(round.Holes))
		}
	}
	payload, err := json.Marshal(examples[0])
	if err != nil {
		t.Fatal(err)
	}
	return payload
}

// Exercise the real router, decoder, service, and response together.
// httptest does this in memory: no listening port or database is needed.
func TestCalculateRoundTotal(t *testing.T) {
	payload := exampleRoundPayload(t)
	request := httptest.NewRequest(http.MethodPost, "/api/v1/rounds/total", strings.NewReader(string(payload)))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	NewRouter(nil).ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
	var body struct {
		TotalStrokes       int `json:"total_strokes"`
		TotalPar           int `json:"total_par"`
		ScoreRelativeToPar int `json:"score_relative_to_par"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.TotalStrokes != 61 {
		t.Fatalf("total_strokes = %d, want 61", body.TotalStrokes)
	}
	if body.TotalPar != 60 {
		t.Fatalf("total_par = %d, want 60", body.TotalPar)
	}
	if body.ScoreRelativeToPar != 1 {
		t.Fatalf("score_relative_to_par = %d, want 1", body.ScoreRelativeToPar)
	}
}

func TestCalculateParTotal(t *testing.T) {
	payload := exampleRoundPayload(t)
	request := httptest.NewRequest(http.MethodPost, "/api/v1/rounds/course_par", strings.NewReader(string(payload)))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	NewRouter(nil).ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
	}
	var body struct {
		TotalPar int `json:"total_par"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.TotalPar != 60 {
		t.Fatalf("total_par = %d, want 60", body.TotalPar)
	}

}

func TestCalculateRoundTotalRejectsInvalidJSON(t *testing.T) {
	for _, payload := range []string{`{`, `{"holes":[{"strokes":"four"}]}`} {
		t.Run(payload, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/api/v1/rounds/total", strings.NewReader(payload))
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()
			NewRouter(nil).ServeHTTP(response, request)
			if response.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want 400", response.Code)
			}
		})
	}
}

func TestScoringEndpointsRejectInvalidRound(t *testing.T) {
	for _, path := range []string{"/api/v1/rounds/total", "/api/v1/rounds/course_par"} {
		t.Run(path, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, path, strings.NewReader(
				`{"username":"JesperNi","course_name":"Practice","date":"2026-10-16T15:42:00+02:00","holes":[{"hole_number":1,"par":3,"strokes":-1}]}`))
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()
			NewRouter(nil).ServeHTTP(response, request)
			if response.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want 400", response.Code)
			}
			var body struct {
				Error string `json:"error"`
			}
			if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
				t.Fatal(err)
			}
			if body.Error != "hole 1: strokes must be positive" {
				t.Fatalf("unexpected error: %q", body.Error)
			}
		})
	}
}

// Check that zero and negative scores are present in the JSON response.
func TestRoundTotalRelativeScores(t *testing.T) {
	for _, tt := range []struct {
		name     string
		strokes  int
		relative int
	}{{"even", 3, 0}, {"under", 2, -1}} {
		t.Run(tt.name, func(t *testing.T) {
			payload, err := json.Marshal(map[string]any{
				"username": "JesperNi", "course_name": "Practice", "date": "2026-10-16T15:42:00+02:00",
				"holes": []map[string]int{{"hole_number": 1, "par": 3, "strokes": tt.strokes}},
			})
			if err != nil {
				t.Fatal(err)
			}
			request := httptest.NewRequest(http.MethodPost, "/api/v1/rounds/total", strings.NewReader(string(payload)))
			request.Header.Set("Content-Type", "application/json")
			response := httptest.NewRecorder()
			NewRouter(nil).ServeHTTP(response, request)
			if response.Code != http.StatusOK {
				t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
			}
			var body map[string]int
			if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
				t.Fatal(err)
			}
			for key, want := range map[string]int{"total_strokes": tt.strokes, "total_par": 3, "score_relative_to_par": tt.relative} {
				if got, ok := body[key]; !ok || got != want {
					t.Errorf("%s = %d (present: %t), want %d", key, got, ok, want)
				}
			}
		})
	}
}

type getRoundStoreFunc func(context.Context, int64) (rounds.Round, error)

func (f getRoundStoreFunc) SaveRound(context.Context, rounds.Round) (int64, error) {
	return 0, errors.New("unexpected SaveRound call")
}

func (f getRoundStoreFunc) GetRound(ctx context.Context, id int64) (rounds.Round, error) {
	return f(ctx, id)
}

func TestGetRoundEndpoint(t *testing.T) {
	want := rounds.Round{
		Username: "JesperNi",
		Id:       42, CourseName: "Practice", Date: "2026-10-29T15:42:00+02:00",
		Holes: []rounds.HoleScore{{HoleNumber: 1, Par: 3, Strokes: 4}},
	}
	for _, tt := range []struct {
		name        string
		id          string
		storeErr    error
		unavailable bool
		status      int
		wantCalls   int
		wantError   string
	}{
		{name: "found", id: "42", status: http.StatusOK, wantCalls: 1},
		{name: "missing", id: "42", storeErr: rounds.ErrRoundNotFound, status: http.StatusNotFound, wantCalls: 1, wantError: "round not found"},
		{name: "storage failure", id: "42", storeErr: errors.New("private database details"), status: http.StatusInternalServerError, wantCalls: 1, wantError: "could not get round"},
		{name: "storage unavailable", id: "42", unavailable: true, status: http.StatusServiceUnavailable, wantError: "round storage is unavailable"},
		{name: "non-numeric", id: "abc", status: http.StatusBadRequest, wantError: "round id must be a positive integer"},
		{name: "zero", id: "0", status: http.StatusBadRequest, wantError: "round id must be a positive integer"},
		{name: "negative", id: "-1", status: http.StatusBadRequest, wantError: "round id must be a positive integer"},
		{name: "overflow", id: "9223372036854775808", status: http.StatusBadRequest, wantError: "round id must be a positive integer"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/api/v1/rounds/"+tt.id, nil)
			calls := 0
			var service *rounds.Service
			if !tt.unavailable {
				service = rounds.NewService(getRoundStoreFunc(func(ctx context.Context, id int64) (rounds.Round, error) {
					calls++
					if ctx != request.Context() || id != want.Id {
						t.Errorf("unexpected storage arguments: context = %v, id = %d", ctx, id)
					}
					return want, tt.storeErr
				}))
			}
			response := httptest.NewRecorder()
			NewRouter(service).ServeHTTP(response, request)
			if response.Code != tt.status {
				t.Fatalf("status = %d, want %d; body = %s", response.Code, tt.status, response.Body.String())
			}
			if calls != tt.wantCalls {
				t.Errorf("storage calls = %d, want %d", calls, tt.wantCalls)
			}
			if tt.status == http.StatusOK {
				var got rounds.Round
				if err := json.Unmarshal(response.Body.Bytes(), &got); err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(got, want) {
					t.Errorf("round = %+v, want %+v", got, want)
				}
			} else {
				var body struct {
					Error string `json:"error"`
				}
				if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
					t.Fatal(err)
				}
				if body.Error != tt.wantError {
					t.Errorf("error = %q, want %q", body.Error, tt.wantError)
				}
			}
		})
	}
}
