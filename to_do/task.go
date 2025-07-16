package main

import (
	"time"
)

type Status string

const (
	ToDo  Status = "to_do"
	Doing Status = "doing"
	Done  Status = "done"
)

type Task struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Category    string     `json:"category"`
	Description string     `json:"description"`
	Status      Status     `json:"status"`
	Deadline    time.Time  `json:"deadline"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}
