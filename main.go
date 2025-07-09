package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Task struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Description string `json:"description"`
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
	}
	tl.Tasks = append(tl.Tasks, task)
}

func (tl *TaskList) Save(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(tl)
}

func (tl *TaskList) Load(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	return decoder.Decode(tl)
}

func (tl *TaskList) ListTasks() {
	fmt.Printf("\n--- TO-DO TASKS ---\n")
	if len(tl.Tasks) == 0 {
		fmt.Println("No tasks found.")
		return
	}
	for _, t := range tl.Tasks {
		fmt.Printf("ID: %d | Name: %s | Category: %s\n", t.ID, t.Name, t.Category)
		fmt.Printf("Description: %s\n", t.Description)
		fmt.Println("---")
	}
}

func main() {
	const filename = "tasks.json"
	taskList := &TaskList{}
	taskList.Load(filename)
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n1. Add Task  2. View Tasks  3. Save and Exit")
		fmt.Print("Choose operation: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			fmt.Print("Enter task name: ")
			name, _ := reader.ReadString('\n')
			name = strings.TrimSpace(name)

			fmt.Print("Enter task category: ")
			category, _ := reader.ReadString('\n')
			category = strings.TrimSpace(category)

			fmt.Print("Enter task description: ")
			description, _ := reader.ReadString('\n')
			description = strings.TrimSpace(description)

			taskList.AddTask(name, category, description)
			fmt.Println("Task added successfully!")

		case "2":
			taskList.ListTasks()

		case "3":
			taskList.Save(filename)
			fmt.Println("Saved successfully, exiting program.")
			return

		default:
			fmt.Println("Invalid choice, please try again.")
		}
	}
}