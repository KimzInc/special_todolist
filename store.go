package main

import (
	"context"
	"errors"
)

// task represents data about a task
type Task struct {
	TaskID      string `json:"task_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	IsCompleted bool   `json:"is_completed"`
}

//Sentinel errors. Every store implementation returns these, so the
// handlers never need to know whether the data lives in the memory or Postgres.

var (
	ErrNotFound = errors.New("Task not found")
	ErrDuplicate = errors.New("Task already exists")
)

// TaskStore is the behaviour the handlers need. Anything with these three
// methods safisfies it, with no implements keyword required.

type TaskStore interface {
	List(ctx context.Context) ([]Task, error)
	Get(ctx context.Context, id string) (Task, error)
	Add(ctx context.Context, t Task) error
}