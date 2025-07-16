package main

import (
	"encoding/json"
	"os"
	"time"
)

type TaskList struct {
	Tasks []Task `json:"tasks"`
}

func (tl *TaskList) AddTask(name, category, description string, deadline time.Time) {
	task := Task{
		ID:          len(tl.Tasks) + 1,
		Name:        name,
		Category:    category,
		Description: description,
		Status:      ToDo,
		Deadline:    deadline,
		CreatedAt:   time.Now(),
	}
	tl.Tasks = append(tl.Tasks, task)
}

func (tl *TaskList) UpdateTask(id int, updates map[string]interface{}) bool {
	for i := range tl.Tasks {
		if tl.Tasks[i].ID == id {
			for field, value := range updates {
				switch field {
				case "name":
					tl.Tasks[i].Name = value.(string)
				case "category":
					tl.Tasks[i].Category = value.(string)
				case "description":
					tl.Tasks[i].Description = value.(string)
				case "status":
					oldStatus := tl.Tasks[i].Status
					newStatus := value.(Status)
					tl.Tasks[i].Status = newStatus

					if newStatus == Done && oldStatus != Done {
						now := time.Now()
						tl.Tasks[i].CompletedAt = &now
					} else if newStatus != Done && oldStatus == Done {
						tl.Tasks[i].CompletedAt = nil
					}
				case "deadline":
					tl.Tasks[i].Deadline = value.(time.Time)
				}
			}
			return true
		}
	}
	return false
}

func (tl *TaskList) FindTask(id int) *Task {
	for i := range tl.Tasks {
		if tl.Tasks[i].ID == id {
			return &tl.Tasks[i]
		}
	}
	return nil
}

func (tl *TaskList) GetTasksByCondition(condition func(Task) bool) []Task {
	var tasks []Task
	for _, task := range tl.Tasks {
		if condition(task) {
			tasks = append(tasks, task)
		}
	}
	return tasks
}

func (tl *TaskList) Save(filename string) error {
	data, err := json.MarshalIndent(tl, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func (tl *TaskList) Load(filename string) error {
	data, err := os.ReadFile(filename)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	return json.Unmarshal(data, tl)
}
