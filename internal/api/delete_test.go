package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"featureflagservice/internal/store"
)

func TestDeleteFlag_RemovesFlag(t *testing.T) {
	st := store.New()
	if err := st.Create(store.Flag{Key: "feature-x", Enabled: true}); err != nil {
		t.Fatalf("seed: %v", err)
	}

	req := httptest.NewRequest(http.MethodDelete, "/flags/feature-x", nil)
	req.SetPathValue("key", "feature-x")
	rec := httptest.NewRecorder()
	DeleteFlag(st)(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNoContent)
	}
	if body := rec.Body.String(); body != "" {
		t.Fatalf("body = %q, want empty", body)
	}
	if _, ok := st.Get("feature-x"); ok {
		t.Fatalf("flag still present after delete")
	}
}

func TestDeleteFlag_UnknownKeyReturns404(t *testing.T) {
	st := store.New()

	req := httptest.NewRequest(http.MethodDelete, "/flags/missing", nil)
	req.SetPathValue("key", "missing")
	rec := httptest.NewRecorder()
	DeleteFlag(st)(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
	if body := rec.Body.String(); !strings.Contains(body, `"error"`) {
		t.Fatalf("body = %q, want JSON error object", body)
	}
}

func TestDeleteFlag_KeyTooLongReturns400(t *testing.T) {
	st := store.New()
	key := strings.Repeat("a", 129)

	req := httptest.NewRequest(http.MethodDelete, "/flags/"+key, nil)
	req.SetPathValue("key", key)
	rec := httptest.NewRecorder()
	DeleteFlag(st)(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}
