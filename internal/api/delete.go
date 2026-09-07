package api

import (
	"net/http"

	"featureflagservice/internal/store"
)

// DeleteFlag handles DELETE /flags/{key}. Implemented by a later feature ticket.
func DeleteFlag(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusNotImplemented, "DELETE /flags/{key} is not implemented yet")
	}
}
