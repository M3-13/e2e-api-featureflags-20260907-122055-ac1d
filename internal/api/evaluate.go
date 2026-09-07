package api

import (
	"net/http"

	"featureflagservice/internal/store"
)

// EvaluateFlag handles GET /flags/{key}/evaluate. Implemented by a later feature ticket.
func EvaluateFlag(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusNotImplemented, "GET /flags/{key}/evaluate is not implemented yet")
	}
}
