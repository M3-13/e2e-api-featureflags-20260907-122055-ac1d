package middleware

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"
)

func captureLog(t *testing.T) (*bytes.Buffer, func()) {
	t.Helper()
	var buf bytes.Buffer
	prev := log.Writer()
	log.SetOutput(&buf)
	return &buf, func() { log.SetOutput(prev) }
}

func TestLoggingLogsMethodPathStatusDuration(t *testing.T) {
	buf, restore := captureLog(t)
	defer restore()

	handler := Logging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte("ok"))
	}))

	req := httptest.NewRequest(http.MethodPost, "/flags/foo/evaluate", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusCreated)
	}

	out := buf.String()
	// One line: POST /flags/foo/evaluate 201 <duration>.
	re := regexp.MustCompile(`POST /flags/foo/evaluate 201 \S+`)
	if !re.MatchString(out) {
		t.Fatalf("log line %q does not match method/path/status/duration", out)
	}
}

func TestLoggingDefaultStatusIs200(t *testing.T) {
	buf, restore := captureLog(t)
	defer restore()

	handler := Logging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	}))

	req := httptest.NewRequest(http.MethodGet, "/flags", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	out := buf.String()
	re := regexp.MustCompile(`GET /flags 200 \S+`)
	if !re.MatchString(out) {
		t.Fatalf("log line %q does not show default 200 status", out)
	}
}

func TestLoggingOmitsQueryStringAndUser(t *testing.T) {
	buf, restore := captureLog(t)
	defer restore()

	handler := Logging(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/flags/myflag/evaluate?user=alice&foo=bar", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	out := buf.String()

	if strings.Contains(out, "alice") {
		t.Fatalf("log leaks user value: %q", out)
	}
	if strings.Contains(out, "user=") {
		t.Fatalf("log leaks user query parameter: %q", out)
	}
	if strings.Contains(out, "foo=bar") || strings.Contains(out, "?") {
		t.Fatalf("log leaks query string: %q", out)
	}
	if !strings.Contains(out, "/flags/myflag/evaluate") {
		t.Fatalf("log missing clean path: %q", out)
	}
}
