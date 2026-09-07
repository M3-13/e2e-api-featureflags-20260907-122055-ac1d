package api

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"featureflagservice/internal/store"
)

func TestCreateFlagCreated(t *testing.T) {
	st := store.New()
	handler := CreateFlag(st)

	body := `{"key":"feature_x","enabled":true,"description":"hello","rollout_percent":50}`
	req := httptest.NewRequest(http.MethodPost, "/flags", strings.NewReader(body))
	rec := httptest.NewRecorder()

	handler(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusCreated, rec.Body.String())
	}

	got, ok := st.Get("feature_x")
	if !ok {
		t.Fatal("flag was not stored")
	}
	if got.Enabled != true || got.Description != "hello" || got.RolloutPercent != 50 {
		t.Fatalf("unexpected stored flag: %+v", got)
	}
}

func TestCreateFlagDefaultRollout(t *testing.T) {
	st := store.New()
	handler := CreateFlag(st)

	body := `{"key":"feature_y","enabled":false}`
	req := httptest.NewRequest(http.MethodPost, "/flags", strings.NewReader(body))
	rec := httptest.NewRecorder()

	handler(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusCreated, rec.Body.String())
	}
	got, ok := st.Get("feature_y")
	if !ok {
		t.Fatal("flag was not stored")
	}
	if got.RolloutPercent != 100 {
		t.Fatalf("rollout = %d, want default 100", got.RolloutPercent)
	}
}

func TestCreateFlagDuplicate(t *testing.T) {
	st := store.New()
	handler := CreateFlag(st)

	body := `{"key":"dup","enabled":true}`
	first := httptest.NewRecorder()
	handler(first, httptest.NewRequest(http.MethodPost, "/flags", strings.NewReader(body)))
	if first.Code != http.StatusCreated {
		t.Fatalf("first create status = %d, want %d", first.Code, http.StatusCreated)
	}

	second := httptest.NewRecorder()
	handler(second, httptest.NewRequest(http.MethodPost, "/flags", strings.NewReader(body)))

	if second.Code != http.StatusConflict {
		t.Fatalf("status = %d, want %d; body=%s", second.Code, http.StatusConflict, second.Body.String())
	}
	if !strings.Contains(second.Body.String(), `"error"`) {
		t.Fatalf("expected JSON error object, got %s", second.Body.String())
	}
}

func TestCreateFlagInvalidJSON(t *testing.T) {
	st := store.New()
	handler := CreateFlag(st)

	req := httptest.NewRequest(http.MethodPost, "/flags", strings.NewReader(`{not json`))
	rec := httptest.NewRecorder()

	handler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d; body=%s", rec.Code, http.StatusBadRequest, rec.Body.String())
	}
}

func TestCreateFlagEmptyKey(t *testing.T) {
	st := store.New()
	handler := CreateFlag(st)

	req := httptest.NewRequest(http.MethodPost, "/flags", strings.NewReader(`{"key":"","enabled":true}`))
	rec := httptest.NewRecorder()

	handler(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestCreateFlagInvalidKeyChars(t *testing.T) {
	st := store.New()
	handler := CreateFlag(st)

	for _, key := range []string{"bad key", "bad/key", "ümlaut"} {
		body := `{"key":` + toJSONString(key) + `,"enabled":true}`
		rec := httptest.NewRecorder()
		handler(rec, httptest.NewRequest(http.MethodPost, "/flags", strings.NewReader(body)))
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("key %q: status = %d, want %d", key, rec.Code, http.StatusBadRequest)
		}
	}
}

func TestCreateFlagRolloutOutOfRange(t *testing.T) {
	st := store.New()
	handler := CreateFlag(st)

	for _, n := range []int{-1, 101} {
		body := `{"key":"range","enabled":true,"rollout_percent":` + itoa(n) + `}`
		rec := httptest.NewRecorder()
		handler(rec, httptest.NewRequest(http.MethodPost, "/flags", strings.NewReader(body)))
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("rollout %d: status = %d, want %d", n, rec.Code, http.StatusBadRequest)
		}
	}
}

func TestCreateFlagBodyTooLarge(t *testing.T) {
	st := store.New()
	handler := CreateFlag(st)

	big := `{"key":"k","description":"` + strings.Repeat("a", maxBodyBytes) + `"}`
	req := httptest.NewRequest(http.MethodPost, "/flags", strings.NewReader(big))
	rec := httptest.NewRecorder()

	handler(rec, req)

	if rec.Code != http.StatusRequestEntityTooLarge && rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 413 or 400", rec.Code)
	}
}

func toJSONString(s string) string {
	return `"` + strings.ReplaceAll(s, `"`, `\"`) + `"`
}

func itoa(n int) string {
	return strconv.Itoa(n)
}
