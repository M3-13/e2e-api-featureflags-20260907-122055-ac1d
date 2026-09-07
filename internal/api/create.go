package api

import (
	"errors"
	"net/http"

	"featureflagservice/internal/store"
)

// defaultRolloutPercent is applied when the request omits rollout_percent.
const defaultRolloutPercent = 100

// createFlagRequest is the JSON body accepted by POST /flags.
type createFlagRequest struct {
	Key            string `json:"key"`
	Enabled        bool   `json:"enabled"`
	Description    string `json:"description"`
	RolloutPercent *int   `json:"rollout_percent"`
}

// CreateFlag handles POST /flags. It validates the request, stores a new flag
// and responds with 201 and the created flag, or an error object.
func CreateFlag(st *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req createFlagRequest
		if err := decodeJSONBody(w, r, &req); err != nil {
			if errors.Is(err, errBodyTooLarge) {
				writeError(w, http.StatusRequestEntityTooLarge, "request body too large")
				return
			}
			writeError(w, http.StatusBadRequest, "invalid JSON")
			return
		}

		if err := validateKey(req.Key); err != nil {
			writeError(w, http.StatusBadRequest, "invalid key")
			return
		}

		rollout := defaultRolloutPercent
		if req.RolloutPercent != nil {
			rollout = *req.RolloutPercent
		}
		if !validateRollout(rollout) {
			writeError(w, http.StatusBadRequest, "invalid rollout_percent")
			return
		}

		f := store.Flag{
			Key:            req.Key,
			Enabled:        req.Enabled,
			Description:    req.Description,
			RolloutPercent: rollout,
		}

		if err := st.Create(f); err != nil {
			if errors.Is(err, store.ErrKeyExists) {
				writeError(w, http.StatusConflict, "key already exists")
				return
			}
			writeError(w, http.StatusInternalServerError, "could not create flag")
			return
		}

		writeJSON(w, http.StatusCreated, f)
	}
}
