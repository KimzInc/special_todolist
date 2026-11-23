package main

import (
	"context"
	"sync"
)

// MemoryStore keeps tasks in a slice. The mutex lives inside the struct,
// next to the data it guards, so no handler can forget to lock.
type MemoryStore struct {
	mu sync.RWMutex
	tasks []Task
}

// Compile-time check: fails the build if *MemoryStore stops satisfying TaskStore.
var _ TaskStore = (*MemoryStore)(nil)

// NewMemoryStore returns a store pre-loaded with seed data.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		tasks: []Task{
			{TaskID: "1", Title: "Buy groceries", Description: "Go to the market and buy all weeks groceries", IsCompleted: false},
			{TaskID: "2", Title: "Pancakes", Description: "Prepare pancakes using EX home baking flour", IsCompleted: false},
			{TaskID: "3", Title: "Chapter 7 Interfaces", Description: "Read the whole of chapter 7 of idiomatic Go interfaces", IsCompleted: false},
		},
	}
}

// List returns a copy of all tasks, so callers can't race with later appends.
func (s *MemoryStore) List(ctx context.Context) ([]Task, error){
	s.mu.RLock()

	defer s.mu.RUnlock()

	out := make([]Task, len(s.tasks))
	copy(out, s.tasks)
	return out, nil
}

// Get returns the task with the given ID, or ErrNotFound.
func (s *MemoryStore) Get(ctx context.Context, id string) (Task, error){
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, t := range s.tasks {
		if t.TaskID == id {
			return t, nil
		}
	}
	return Task{}, ErrNotFound
}


// Add appends a task, or returns ErrDuplicate if the ID is already taken.
// The duplicate check and the append happen under one lock, so two
// concurrent requests can't both pass the check.

func (s *MemoryStore) Add(ctx context.Context, t Task) error  {
	s.mu.Lock()

	defer s.mu.Unlock()

	for _, existing := range s.tasks {
		if existing.TaskID == t.TaskID {
			return ErrDuplicate
		}
		
	}
	s.tasks = append(s.tasks, t)
	return nil
}