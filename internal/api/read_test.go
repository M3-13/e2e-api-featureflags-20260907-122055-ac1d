package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"featureflagservice/internal/store"
)

func TestListFlags_Empty(t *testing.T) {
	st := store.New()
	req := httptest.NewRequest(http.MethodGet, "/flags", nil)
	rec := httptest.NewRecorder()

	ListFlags(st)(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("Content-Type = %q, want application/json", ct)
	}
	// The body must be the literal empty array, never null.
	if got := strings.TrimSpace(rec.Body.String()); got != "[]" {
		t.Fatalf("body = %q, want []", got)
	}
}

func TestListFlags_Multiple(t *testing.T) {
	st := store.New()
	st.Create(store.Flag{Key: "alpha", Enabled: true, Description: "a", RolloutPercent: 100})
	st.Create(store.Flag{Key: "beta", Enabled: false, Description: "b", RolloutPercent: 0})

	req := httptest.NewRequest(http.MethodGet, "/flags", nil)
	rec := httptest.NewRecorder()
	ListFlags(st)(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	var flags []store.Flag
	if err := json.Unmarshal(rec.Body.Bytes(), &flags); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(flags) != 2 {
		t.Fatalf("len(flags) = %d, want 2", len(flags))
	}
}

func TestGetFlag_Found(t *testing.T) {
	st := store.New()
	want := store.Flag{Key: "alpha", Enabled: true, Description: "a", RolloutPercent: 50}
	st.Create(want)

	req := httptest.NewRequest(http.MethodGet, "/flags/alpha", nil)
	req.SetPathValue("key", "alpha")
	rec := httptest.NewRecorder()
	GetFlag(st)(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	var got store.Flag
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if got.Key != want.Key || got.Enabled != want.Enabled || got.RolloutPercent != want.RolloutPercent {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestGetFlag_NotFound(t *testing.T) {
	st := store.New()

	req := httptest.NewRequest(http.MethodGet, "/flags/missing", nil)
	req.SetPathValue("key", "missing")
	rec := httptest.NewRecorder()
	GetFlag(st)(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if _, ok := body["error"]; !ok {
		t.Fatalf("body = %v, want an \"error\" field", body)
	}
}

func TestGetFlag_KeyTooLong(t *testing.T) {
	st := store.New()
	longKey := strings.Repeat("a", 129)

	req := httptest.NewRequest(http.MethodGet, "/flags/"+longKey, nil)
	req.SetPathValue("key", longKey)
	rec := httptest.NewRecorder()
	GetFlag(st)(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if _, ok := body["error"]; !ok {
		t.Fatalf("body = %v, want an \"error\" field", body)
	}
}
