package api

import (
	"net/http"

	"featureflagservice/internal/store"
)

// DeleteFlag handles DELETE /flags/{key}. It removes the flag and responds
// with 204, or 404 when the key does not exist. An invalid key (empty or
// longer than 128 characters) yields 400.
func DeleteFlag(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.PathValue("key")
		if err := validateKey(key); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		if !st.Delete(key) {
			writeError(w, http.StatusNotFound, "flag not found")
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
