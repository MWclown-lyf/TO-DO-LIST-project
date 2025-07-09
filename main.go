package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Status string

const (
	ToDo  Status = "to_do"
	Doing Status = "doing"
	Done  Status = "done"
)

type Task struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Description string `json:"description"`
	Status      Status `json:"status"`
}

type TaskList struct {
	Tasks []Task `json:"tasks"`
}

func (tl *TaskList) AddTask(name, category, description string) {
	id := 1
	if len(tl.Tasks) > 0 {
		id = tl.Tasks[len(tl.Tasks)-1].ID + 1
	}
	task := Task{
		ID:          id,
		Name:        name,
		Category:    category,
		Description: description,
		Status:      ToDo,
	}
	tl.Tasks = append(tl.Tasks, task)
}

func (tl *TaskList) UpdateTask(id int, name, category, description string, status Status) bool {
	for i, t := range tl.Tasks {
		if t.ID == id {
			tl.Tasks[i].Name = name
			tl.Tasks[i].Category = category
			tl.Tasks[i].Description = description
			tl.Tasks[i].Status = status
			return true
		}
	}
	return false
}

