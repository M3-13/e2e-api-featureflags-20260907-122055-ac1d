package api

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"featureflagservice/internal/store"
)

func newUpdateStore(t *testing.T) *store.Store {
	t.Helper()
	st := store.New()
	if err := st.Create(store.Flag{
		Key:            "beta",
		Enabled:        false,
		Description:    "old",
		RolloutPercent: 0,
	}); err != nil {
		t.Fatalf("seed flag: %v", err)
	}
	return st
}

func doUpdate(st *store.Store, key, body string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodPut, "/flags/"+key, strings.NewReader(body))
	req.SetPathValue("key", key)
	rec := httptest.NewRecorder()
	UpdateFlag(st)(rec, req)
	return rec
}

func TestUpdateFlagEnabledOnly(t *testing.T) {
	st := newUpdateStore(t)
	rec := doUpdate(st, "beta", `{"enabled":true}`)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body.String())
	}

	f, ok := st.Get("beta")
	if !ok {
		t.Fatal("flag not found after update")
	}
	if !f.Enabled {
		t.Errorf("enabled = false, want true")
	}
	if f.Description != "old" {
		t.Errorf("description = %q, want unchanged %q", f.Description, "old")
	}
	if f.RolloutPercent != 0 {
		t.Errorf("rollout_percent = %d, want unchanged 0", f.RolloutPercent)
	}
}

func TestUpdateFlagDescriptionOnly(t *testing.T) {
	st := newUpdateStore(t)
	rec := doUpdate(st, "beta", `{"description":"new desc"}`)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body.String())
	}

	f, _ := st.Get("beta")
	if f.Description != "new desc" {
		t.Errorf("description = %q, want %q", f.Description, "new desc")
	}
	if f.Enabled {
		t.Errorf("enabled = true, want unchanged false")
	}
}

func TestUpdateFlagRolloutOnly(t *testing.T) {
	st := newUpdateStore(t)
	rec := doUpdate(st, "beta", `{"rollout_percent":75}`)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200; body = %s", rec.Code, rec.Body.String())
	}

	f, _ := st.Get("beta")
	if f.RolloutPercent != 75 {
		t.Errorf("rollout_percent = %d, want 75", f.RolloutPercent)
	}
}

func TestUpdateFlagUnknownKey(t *testing.T) {
	st := newUpdateStore(t)
	rec := doUpdate(st, "missing", `{"enabled":true}`)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
}

func TestUpdateFlagRolloutOutOfRange(t *testing.T) {
	for _, v := range []int{-1, 101} {
		st := newUpdateStore(t)
		body := `{"rollout_percent":` + strconv.Itoa(v) + `}`
		rec := doUpdate(st, "beta", body)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("rollout_percent=%d: status = %d, want 400", v, rec.Code)
		}

		f, _ := st.Get("beta")
		if f.RolloutPercent != 0 {
			t.Errorf("rollout_percent=%d: flag was modified, rollout = %d", v, f.RolloutPercent)
		}
	}
}

func TestUpdateFlagBodyTooLarge(t *testing.T) {
	st := newUpdateStore(t)
	// A description field large enough to blow past the 1 MiB read limit.
	body := `{"description":"` + strings.Repeat("x", maxBodyBytes) + `"}`
	rec := doUpdate(st, "beta", body)

	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("status = %d, want 413", rec.Code)
	}
}

func TestUpdateFlagKeyTooLong(t *testing.T) {
	st := newUpdateStore(t)
	longKey := strings.Repeat("a", 129)
	rec := doUpdate(st, longKey, `{"enabled":true}`)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", rec.Code)
	}
}
