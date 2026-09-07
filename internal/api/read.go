package api

import (
	"net/http"

	"featureflagservice/internal/store"
)

// ListFlags handles GET /flags. It returns all flags as a JSON array; an empty
// store yields an empty (non-null) array.
func ListFlags(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, st.List())
	}
}

// GetFlag handles GET /flags/{key}. It returns the flag with the given key, a
// 404 JSON error object when the key is unknown, and a 400 JSON error object
// when the key is invalid (e.g. longer than 128 characters).
func GetFlag(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.PathValue("key")
		if err := validateKey(key); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		flag, ok := st.Get(key)
		if !ok {
			writeError(w, http.StatusNotFound, "flag not found")
			return
		}
		writeJSON(w, http.StatusOK, flag)
	}
}
