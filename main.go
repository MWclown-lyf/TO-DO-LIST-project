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

func (tl *TaskList) UpdateTaskField(id int, field string, value interface{}) bool {
	for i, t := range tl.Tasks {
		if t.ID == id {
			switch field {
			case "name":
				tl.Tasks[i].Name = value.(string)
			case "category":
				tl.Tasks[i].Category = value.(string)
			case "description":
				tl.Tasks[i].Description = value.(string)
			case "status":
				tl.Tasks[i].Status = value.(Status)
			}
			return true
		}
	}
	return false
}

func (tl *TaskList) FindTask(id int) *Task {
	for _, t := range tl.Tasks {
		if t.ID == id {
			return &t
		}
	}
	return nil
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

func (tl *TaskList) ListByStatus(status Status) {
	fmt.Printf("\n--- %s ---\n", status)
	for _, t := range tl.Tasks {
		if t.Status == status {
			fmt.Printf("ID: %d | Name: %s | Category: %s\n", t.ID, t.Name, t.Category)
			fmt.Printf("Description: %s\n", t.Description)
			fmt.Println("---")
		}
	}
}

func main() {
	const filename = "tasks.json"
	taskList := &TaskList{}
	taskList.Load(filename)
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n1. Add Task  2. Update Task  3. View Tasks  4. Save and Exit")
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
			fmt.Print("Enter task ID to update: ")
			var id int
			fmt.Scanf("%d\n", &id)

			task := taskList.FindTask(id)
			if task == nil {
				fmt.Println("Task not found!")
				continue
			}

			fmt.Printf("Current task info:\n")
			fmt.Printf("Name: %s\nCategory: %s\nDescription: %s\nStatus: %s\n",
				task.Name, task.Category, task.Description, task.Status)

		updateLoop:
			for {
				fmt.Println("\nWhat would you like to update?")
				fmt.Println("1. Name  2. Category  3. Description  4. Status  5. Finish updating")
				fmt.Print("Choose field: ")
				fieldInput, _ := reader.ReadString('\n')
				fieldInput = strings.TrimSpace(fieldInput)

				switch fieldInput {
				case "1":
					fmt.Printf("Current name: %s\n", task.Name)
					fmt.Print("Enter new name (press Enter to keep current): ")
					newName, _ := reader.ReadString('\n')
					newName = strings.TrimSpace(newName)
					if newName != "" {
						taskList.UpdateTaskField(id, "name", newName)
						task.Name = newName
						fmt.Println("Name updated!")
					}

				case "2":
					fmt.Printf("Current category: %s\n", task.Category)
					fmt.Print("Enter new category (press Enter to keep current): ")
					newCategory, _ := reader.ReadString('\n')
					newCategory = strings.TrimSpace(newCategory)
					if newCategory != "" {
						taskList.UpdateTaskField(id, "category", newCategory)
						task.Category = newCategory
						fmt.Println("Category updated!")
					}

				case "3":
					fmt.Printf("Current description: %s\n", task.Description)
					fmt.Print("Enter new description (press Enter to keep current): ")
					newDescription, _ := reader.ReadString('\n')
					newDescription = strings.TrimSpace(newDescription)
					if newDescription != "" {
						taskList.UpdateTaskField(id, "description", newDescription)
						task.Description = newDescription
						fmt.Println("Description updated!")
					}

				case "4":
					fmt.Printf("Current status: %s\n", task.Status)
					fmt.Print("Enter new status (to_do/doing/done, press Enter to keep current): ")
					statusStr, _ := reader.ReadString('\n')
					statusStr = strings.TrimSpace(statusStr)

					if statusStr != "" {
						var status Status
						switch statusStr {
						case "to_do":
							status = ToDo
						case "doing":
							status = Doing
						case "done":
							status = Done
						default:
							fmt.Println("Invalid status!")
							continue
						}
						taskList.UpdateTaskField(id, "status", status)
						task.Status = status
						fmt.Println("Status updated!")
					}

				case "5":
					fmt.Println("Task update completed!")
					break updateLoop

				default:
					fmt.Println("Invalid choice, please try again.")
				}
			}

		case "3":
			taskList.ListByStatus(ToDo)
			taskList.ListByStatus(Doing)
			taskList.ListByStatus(Done)

		case "4":
			taskList.Save(filename)
			fmt.Println("Saved successfully, exiting program.")
			return

		default:
			fmt.Println("Invalid choice, please try again.")
		}
	}
}
