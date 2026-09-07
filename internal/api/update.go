package api

import (
	"net/http"

	"featureflagservice/internal/store"
)

// UpdateFlag handles PUT /flags/{key}. Implemented by a later feature ticket.
func UpdateFlag(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusNotImplemented, "PUT /flags/{key} is not implemented yet")
	}
}
