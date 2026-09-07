package api

import (
	"errors"
	"net/http"

	"featureflagservice/internal/store"
)

// updateRequest is the JSON body accepted by PUT /flags/{key}. Every field is a
// pointer so the handler can tell "absent" from "explicitly zero".
type updateRequest struct {
	Enabled        *bool   `json:"enabled"`
	Description    *string `json:"description"`
	RolloutPercent *int    `json:"rollout_percent"`
}

// UpdateFlag handles PUT /flags/{key}. It applies a partial update: only the
// fields present in the body are changed. It answers 200 with the updated flag,
// 404 when the key is unknown, and 400 for an invalid key, an invalid
// rollout_percent, or an unreadable body (413 when the body exceeds 1 MiB).
func UpdateFlag(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.PathValue("key")

		if err := validateKey(key); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		var req updateRequest
		if err := decodeJSONBody(w, r, &req); err != nil {
			if errors.Is(err, errBodyTooLarge) {
				writeError(w, http.StatusRequestEntityTooLarge, "request body too large")
				return
			}
			writeError(w, http.StatusBadRequest, "invalid JSON body")
			return
		}

		if req.RolloutPercent != nil && !validateRollout(*req.RolloutPercent) {
			writeError(w, http.StatusBadRequest, "rollout_percent must be between 0 and 100")
			return
		}

		flag, ok := st.Update(key, store.Patch{
			Enabled:        req.Enabled,
			Description:    req.Description,
			RolloutPercent: req.RolloutPercent,
		})
		if !ok {
			writeError(w, http.StatusNotFound, "flag not found")
			return
		}

		writeJSON(w, http.StatusOK, flag)
	}
}
