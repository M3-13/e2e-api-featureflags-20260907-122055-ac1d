package api

import (
	"hash/fnv"
	"net/http"

	"featureflagservice/internal/store"
)

// evaluateResponse is the JSON body returned by a successful evaluation.
type evaluateResponse struct {
	Key            string `json:"key"`
	User           string `json:"user"`
	Enabled        bool   `json:"enabled"`
	RolloutPercent int    `json:"rollout_percent"`
	Evaluated      bool   `json:"evaluated"`
}

// EvaluateFlag handles GET /flags/{key}/evaluate?user={id}.
//
// The decision is deterministic: the same key and user always map to the same
// evaluated result for a given flag, via a stable FNV-1a hash modulo 100.
func EvaluateFlag(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := r.URL.Query().Get("user")
		if user == "" {
			writeError(w, http.StatusBadRequest, "missing user query parameter")
			return
		}
		if len(user) > 128 {
			writeError(w, http.StatusBadRequest, "user must not exceed 128 characters")
			return
		}

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

		evaluated := false
		if flag.Enabled {
			h := fnv.New64a()
			_, _ = h.Write([]byte(key))
			_, _ = h.Write([]byte(":"))
			_, _ = h.Write([]byte(user))
			evaluated = (h.Sum64() % 100) < uint64(flag.RolloutPercent)
		}

		writeJSON(w, http.StatusOK, evaluateResponse{
			Key:            flag.Key,
			User:           user,
			Enabled:        flag.Enabled,
			RolloutPercent: flag.RolloutPercent,
			Evaluated:      evaluated,
		})
	}
}
