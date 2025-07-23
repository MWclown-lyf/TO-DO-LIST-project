package main

import (
	"os"
	"testing"
	"time"
)

func TestAddTask(t *testing.T) {
	tl := &TaskList{}
	deadline := time.Now().Add(24 * time.Hour)
	tl.AddTask("Task1", "Work", "Description1", deadline)

	if len(tl.Tasks) != 1 {
		t.Fatalf("Expected 1 task, got %d", len(tl.Tasks))
	}
	if tl.Tasks[0].Name != "Task1" {
		t.Errorf("Task name is incorrect")
	}
}

func TestUpdateTask(t *testing.T) {
	tl := &TaskList{}
	deadline := time.Now().Add(24 * time.Hour)
	tl.AddTask("Task1", "Work", "Description1", deadline)

	updates := map[string]interface{}{
		"name":        "NewTaskName",
		"category":    "Life",
		"description": "NewDescription",
		"status":      Done,
	}
	ok := tl.UpdateTask(1, updates)
	if !ok {
		t.Fatalf("Failed to update task")
	}
	task := tl.Tasks[0]
	if task.Name != "NewTaskName" || task.Category != "Life" || task.Description != "NewDescription" || task.Status != Done {
		t.Errorf("Task update is incorrect")
	}
	if task.CompletedAt == nil {
		t.Errorf("Completion time not set")
	}
}

func TestFindTask(t *testing.T) {
	tl := &TaskList{}
	tl.AddTask("Task1", "Work", "Description1", time.Now())
	task := tl.FindTask(1)
	if task == nil {
		t.Fatalf("Task not found")
	}
	if task.Name != "Task1" {
		t.Errorf("Task name is incorrect")
	}
}

func TestGetTasksByCondition(t *testing.T) {
	tl := &TaskList{}
	tl.AddTask("Task1", "Work", "Description1", time.Now())
	tl.AddTask("Task2", "Life", "Description2", time.Now())
	tasks := tl.GetTasksByCondition(func(task Task) bool {
		return task.Category == "Life"
	})
	if len(tasks) != 1 || tasks[0].Name != "Task2" {
		t.Errorf("Filtering tasks by condition failed")
	}
}

func TestSaveAndLoad(t *testing.T) {
	tl := &TaskList{}
	tl.AddTask("Task1", "Work", "Description1", time.Now())
	filename := "test_tasks.json"
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			t.Fatalf("Failed to clean up test file: %v", err)
		}
	}(filename)

	if err := tl.Save(filename); err != nil {
		t.Fatalf("Failed to save: %v", err)
	}

	tl2 := &TaskList{}
	if err := tl2.Load(filename); err != nil {
		t.Fatalf("Failed to load: %v", err)
	}
	if len(tl2.Tasks) != 1 || tl2.Tasks[0].Name != "Task1" {
		t.Errorf("Loaded content is incorrect")
	}
}
