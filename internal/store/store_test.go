package store

import (
	"fmt"
	"sync"
	"testing"
)

func TestCreateGetUpdateDelete(t *testing.T) {
	s := New()

	if err := s.Create(Flag{Key: "a", Enabled: true, RolloutPercent: 100}); err != nil {
		t.Fatalf("Create returned unexpected error: %v", err)
	}
	if err := s.Create(Flag{Key: "a", Enabled: false}); err != ErrKeyExists {
		t.Fatalf("Create duplicate key: got %v, want ErrKeyExists", err)
	}

	f, ok := s.Get("a")
	if !ok || f.Key != "a" || !f.Enabled {
		t.Fatalf("Get returned unexpected result: %+v, ok=%v", f, ok)
	}
	if _, ok := s.Get("missing"); ok {
		t.Fatal("Get for missing key should return ok=false")
	}

	enabled := false
	desc := "updated"
	rollout := 42
	f, ok = s.Update("a", Patch{Enabled: &enabled, Description: &desc, RolloutPercent: &rollout})
	if !ok || f.Enabled || f.Description != "updated" || f.RolloutPercent != 42 {
		t.Fatalf("Update returned unexpected result: %+v, ok=%v", f, ok)
	}
	if _, ok := s.Update("missing", Patch{Enabled: &enabled}); ok {
		t.Fatal("Update for missing key should return ok=false")
	}

	if !s.Delete("a") {
		t.Fatal("Delete of existing key should return true")
	}
	if s.Delete("a") {
		t.Fatal("Delete of missing key should return false")
	}
	if _, ok := s.Get("a"); ok {
		t.Fatal("Get after Delete should return ok=false")
	}
}

func TestListNeverNil(t *testing.T) {
	s := New()
	if l := s.List(); l == nil {
		t.Fatal("List on empty store should not be nil")
	}
	s.Create(Flag{Key: "x"})
	if l := s.List(); len(l) != 1 {
		t.Fatalf("List length = %d, want 1", len(l))
	}
}

func TestConcurrentCreateGetUpdateDelete(t *testing.T) {
	s := New()
	const workers = 32
	const perWorker = 25

	var wg sync.WaitGroup
	// Concurrent creates of distinct keys.
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(base int) {
			defer wg.Done()
			for j := 0; j < perWorker; j++ {
				key := fmt.Sprintf("key-%d-%d", base, j)
				_ = s.Create(Flag{Key: key, Enabled: true, RolloutPercent: 100})
			}
		}(i)
	}
	wg.Wait()

	// Concurrent reads.
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(base int) {
			defer wg.Done()
			for j := 0; j < perWorker; j++ {
				key := fmt.Sprintf("key-%d-%d", base, j)
				_, _ = s.Get(key)
				_ = s.List()
			}
		}(i)
	}

	// Concurrent updates of existing keys.
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(base int) {
			defer wg.Done()
			for j := 0; j < perWorker; j++ {
				key := fmt.Sprintf("key-%d-%d", base, j)
				enabled := false
				_, _ = s.Update(key, Patch{Enabled: &enabled})
			}
		}(i)
	}

	// Concurrent deletes of a shared set of keys.
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(base int) {
			defer wg.Done()
			for j := 0; j < perWorker; j++ {
				key := fmt.Sprintf("key-%d-%d", base, j)
				_ = s.Delete(key)
			}
		}(i)
	}

	wg.Wait()
}
