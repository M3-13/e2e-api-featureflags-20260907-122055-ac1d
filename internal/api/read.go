package api

import (
	"net/http"

	"featureflagservice/internal/store"
)

// ListFlags handles GET /flags. Implemented by a later feature ticket.
func ListFlags(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusNotImplemented, "GET /flags is not implemented yet")
	}
}

// GetFlag handles GET /flags/{key}. Implemented by a later feature ticket.
func GetFlag(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusNotImplemented, "GET /flags/{key} is not implemented yet")
	}
}
