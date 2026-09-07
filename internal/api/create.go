package api

import (
	"net/http"

	"featureflagservice/internal/store"
)

// CreateFlag handles POST /flags. Implemented by a later feature ticket.
func CreateFlag(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusNotImplemented, "POST /flags is not implemented yet")
	}
}
