package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"featureflagservice/internal/store"
)

func TestHealthz(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	healthz(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if got := rec.Body.String(); got != `{"status":"ok"}` {
		t.Fatalf("body = %q, want %q", got, `{"status":"ok"}`)
	}
}

// TestRoutesAreRegistered proves each handler is wired onto the router under
// the agreed path and verb. It does NOT assert what the stub handlers answer
// today (501) — that answer changes when the feature tickets land. Instead it
// probes each path with a verb that is NOT registered for it: a Go 1.22
// ServeMux answers 405 Method Not Allowed (with an Allow header naming the
// registered verbs) when a path is registered but the method is not, and 404
// when no pattern matches at all. Both facts are stable across the stub→real
// transition, so this test stays green before and after the feature tickets.
func TestRoutesAreRegistered(t *testing.T) {
	srv := httptest.NewServer(newHandler(store.New()))
	defer srv.Close()

	tests := []struct {
		name   string
		method string // the verb this route is registered under
		path   string // a concrete path matching the registered pattern
		probe  string // a verb NOT registered for that path
	}{
		{"POST /flags", http.MethodPost, "/flags", http.MethodPatch},
		{"GET /flags", http.MethodGet, "/flags", http.MethodPatch},
		{"GET /flags/{key}", http.MethodGet, "/flags/some-key", http.MethodPatch},
		{"PUT /flags/{key}", http.MethodPut, "/flags/some-key", http.MethodPatch},
		{"DELETE /flags/{key}", http.MethodDelete, "/flags/some-key", http.MethodPatch},
		{"GET /flags/{key}/evaluate", http.MethodGet, "/flags/some-key/evaluate", http.MethodPost},
		{"GET /healthz", http.MethodGet, "/healthz", http.MethodPost},
	}

	for _, tt := range tests {
		req, err := http.NewRequest(tt.probe, srv.URL+tt.path, nil)
		if err != nil {
			t.Fatalf("%s: %v", tt.name, err)
		}
		res, err := srv.Client().Do(req)
		if err != nil {
			t.Fatalf("%s: %v", tt.name, err)
		}
		res.Body.Close()

		if res.StatusCode != http.StatusMethodNotAllowed {
			t.Errorf("%s: probing %s %s => %d, want %d (route not registered)",
				tt.name, tt.probe, tt.path, res.StatusCode, http.StatusMethodNotAllowed)
			continue
		}
		allow := res.Header.Get("Allow")
		if !strings.Contains(allow, tt.method) {
			t.Errorf("%s: Allow header %q does not name %s", tt.name, allow, tt.method)
		}
	}
}
