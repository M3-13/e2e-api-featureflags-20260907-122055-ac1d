package store

import (
	"errors"
	"sync"
)

// ErrKeyExists is returned by Create when a flag with the same key already exists.
var ErrKeyExists = errors.New("key already exists")

// Flag is a feature-flag record. Only these fields are ever stored; no
// per-user data is retained beyond the lifetime of a single request.
type Flag struct {
	Key            string `json:"key"`
	Enabled        bool   `json:"enabled"`
	Description    string `json:"description,omitempty"`
	RolloutPercent int    `json:"rollout_percent"`
}

// Patch describes a partial update to a Flag. Only non-nil fields are applied.
type Patch struct {
	Enabled        *bool
	Description    *string
	RolloutPercent *int
}

// Store is a thread-safe in-memory feature-flag store.
type Store struct {
	mu    sync.RWMutex
	flags map[string]Flag
}

// New returns an empty, ready-to-use Store.
func New() *Store {
	return &Store{flags: make(map[string]Flag)}
}

// Create inserts a new flag. It returns ErrKeyExists if the key is taken.
func (s *Store) Create(f Flag) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.flags[f.Key]; ok {
		return ErrKeyExists
	}
	s.flags[f.Key] = f
	return nil
}

// List returns all flags. The result is never nil; an empty store yields an
// empty (non-nil) slice.
func (s *Store) List() []Flag {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Flag, 0, len(s.flags))
	for _, f := range s.flags {
		out = append(out, f)
	}
	return out
}

// Get returns the flag with the given key and whether it was found.
func (s *Store) Get(key string) (Flag, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	f, ok := s.flags[key]
	return f, ok
}

// Update applies the non-nil fields of p to the flag with the given key. It
// returns the updated flag and whether the key existed.
func (s *Store) Update(key string, p Patch) (Flag, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	f, ok := s.flags[key]
	if !ok {
		return Flag{}, false
	}
	if p.Enabled != nil {
		f.Enabled = *p.Enabled
	}
	if p.Description != nil {
		f.Description = *p.Description
	}
	if p.RolloutPercent != nil {
		f.RolloutPercent = *p.RolloutPercent
	}
	s.flags[key] = f
	return f, true
}

// Delete removes the flag with the given key and reports whether it existed.
func (s *Store) Delete(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.flags[key]; !ok {
		return false
	}
	delete(s.flags, key)
	return true
}
