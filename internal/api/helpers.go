package api

import (
	"encoding/json"
	"errors"
	"net/http"
)

// maxBodyBytes limits how much of a request body the service will read (1 MiB).
const maxBodyBytes = 1 << 20

// errBodyTooLarge is returned by decodeJSONBody when the body exceeds maxBodyBytes.
var errBodyTooLarge = errors.New("request body too large")

// writeJSON encodes v as JSON with the given HTTP status.
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// writeError writes a JSON error object {"error": msg} with the given status.
// Callers must pass a fixed, internal message; never client-supplied data.
func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// validateKey checks that key is non-empty, contains only [A-Za-z0-9_-], and is
// at most 128 characters long.
func validateKey(key string) error {
	if key == "" {
		return errors.New("key must not be empty")
	}
	if len(key) > 128 {
		return errors.New("key must not exceed 128 characters")
	}
	for _, r := range key {
		switch {
		case r >= 'a' && r <= 'z':
		case r >= 'A' && r <= 'Z':
		case r >= '0' && r <= '9':
		case r == '-' || r == '_':
		default:
			return errors.New("key must contain only [A-Za-z0-9_-]")
		}
	}
	return nil
}

// validateRollout reports whether n is within the valid rollout range 0–100.
func validateRollout(n int) bool {
	return n >= 0 && n <= 100
}

// decodeJSONBody decodes the request body into dst, limiting reads to 1 MiB via
// http.MaxBytesReader. It returns errBodyTooLarge when the body exceeds the
// limit, or the underlying decode error otherwise.
func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(dst); err != nil {
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			return errBodyTooLarge
		}
		return err
	}
	return nil
}
