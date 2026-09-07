package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"featureflagservice/internal/store"
)

var longString = strings.Repeat("a", 129)

func newEvalServer(st *store.Store) http.Handler {
	return EvaluateFlag(st)
}

func doEval(t *testing.T, h http.Handler, key, user string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/flags/"+key+"/evaluate?user="+user, nil)
	req.SetPathValue("key", key)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func TestEvaluateDeterministic(t *testing.T) {
	st := store.New()
	if err := st.Create(store.Flag{Key: "feature", Enabled: true, RolloutPercent: 50}); err != nil {
		t.Fatal(err)
	}
	h := newEvalServer(st)

	first := doEval(t, h, "feature", "user-1")
	second := doEval(t, h, "feature", "user-1")

	if first.Code != http.StatusOK || second.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d and %d", first.Code, second.Code)
	}

	var a, b evaluateResponse
	if err := json.NewDecoder(first.Body).Decode(&a); err != nil {
		t.Fatal(err)
	}
	if err := json.NewDecoder(second.Body).Decode(&b); err != nil {
		t.Fatal(err)
	}
	if a.Evaluated != b.Evaluated {
		t.Fatalf("expected deterministic result, got %v and %v", a.Evaluated, b.Evaluated)
	}
}

func TestEvaluateDisabledAlwaysFalse(t *testing.T) {
	st := store.New()
	if err := st.Create(store.Flag{Key: "feature", Enabled: false, RolloutPercent: 100}); err != nil {
		t.Fatal(err)
	}
	h := newEvalServer(st)

	for _, user := range []string{"u1", "u2", "u3"} {
		rec := doEval(t, h, "feature", user)
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rec.Code)
		}
		var body evaluateResponse
		if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body.Evaluated {
			t.Fatalf("expected evaluated=false when disabled, got true for user %s", user)
		}
	}
}

func TestEvaluateRolloutPercent100True(t *testing.T) {
	st := store.New()
	if err := st.Create(store.Flag{Key: "feature", Enabled: true, RolloutPercent: 100}); err != nil {
		t.Fatal(err)
	}
	h := newEvalServer(st)

	for _, user := range []string{"u1", "u2", "u3", "u4", "u5"} {
		rec := doEval(t, h, "feature", user)
		var body evaluateResponse
		if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if !body.Evaluated {
			t.Fatalf("expected evaluated=true for rollout 100, got false for user %s", user)
		}
	}
}

func TestEvaluateRolloutPercent0False(t *testing.T) {
	st := store.New()
	if err := st.Create(store.Flag{Key: "feature", Enabled: true, RolloutPercent: 0}); err != nil {
		t.Fatal(err)
	}
	h := newEvalServer(st)

	for _, user := range []string{"u1", "u2", "u3", "u4", "u5"} {
		rec := doEval(t, h, "feature", user)
		var body evaluateResponse
		if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body.Evaluated {
			t.Fatalf("expected evaluated=false for rollout 0, got true for user %s", user)
		}
	}
}

func TestEvaluateUnknownKeyNotFound(t *testing.T) {
	st := store.New()
	h := newEvalServer(st)

	rec := doEval(t, h, "missing", "user-1")
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestEvaluateMissingUserBadRequest(t *testing.T) {
	st := store.New()
	if err := st.Create(store.Flag{Key: "feature", Enabled: true, RolloutPercent: 100}); err != nil {
		t.Fatal(err)
	}
	h := newEvalServer(st)

	req := httptest.NewRequest(http.MethodGet, "/flags/feature/evaluate", nil)
	req.SetPathValue("key", "feature")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestEvaluateUserTooLongBadRequest(t *testing.T) {
	st := store.New()
	if err := st.Create(store.Flag{Key: "feature", Enabled: true, RolloutPercent: 100}); err != nil {
		t.Fatal(err)
	}
	h := newEvalServer(st)

	rec := doEval(t, h, "feature", longString)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestEvaluateKeyTooLongBadRequest(t *testing.T) {
	st := store.New()
	h := newEvalServer(st)

	longKey := strings.Repeat("a", 129)

	rec := doEval(t, h, longKey, "user-1")
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}
