package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func main() {
	const filename = "tasks.json"
	taskList := &TaskList{}
	taskList.Load(filename)
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Current time: %s\n", time.Now().Format("2006-01-02 15:04"))
	taskList.ShowReminders()

	for {
		fmt.Println("\n1. Add Task  2. Update Task  3. View Tasks  4. Reminders  5. Statistics  6. Exit")
		choice := readInput(reader, "Choose: ")

		switch choice {
		case "1":
			addTaskInteractive(taskList, reader)
		case "2":
			updateTaskInteractive(taskList, reader)
		case "3":
			taskList.Display()
		case "4":
			taskList.ShowReminders()
		case "5":
			taskList.ShowStats()
		case "6":
			taskList.Save(filename)
			fmt.Println("👋 Goodbye!")
			return
		default:
			fmt.Println("❌ Invalid choice!")
		}
	}
}

func addTaskInteractive(taskList *TaskList, reader *bufio.Reader) {
	name := readInput(reader, "Task name: ")
	category := readInput(reader, "Category: ")
	description := readInput(reader, "Description: ")

	deadlineStr := readInput(reader, "Deadline (YYYY-MM-DD HH:MM): ")
	deadline, err := parseTime(deadlineStr)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	taskList.AddTask(name, category, description, deadline)
	fmt.Println("✅ Task added!")
}

func updateTaskInteractive(taskList *TaskList, reader *bufio.Reader) {
	id, err := readInt(reader, "Task ID: ")
	if err != nil || taskList.FindTask(id) == nil {
		fmt.Println("❌ Task not found!")
		return
	}

	fmt.Println("Fields: 1=Name 2=Category 3=Description 4=Status 5=Deadline")
	field := readInput(reader, "Update field (1-5): ")

	updates := make(map[string]interface{})
	switch field {
	case "1":
		updates["name"] = readInput(reader, "New name: ")
	case "2":
		updates["category"] = readInput(reader, "New category: ")
	case "3":
		updates["description"] = readInput(reader, "New description: ")
	case "4":
		statusStr := readInput(reader, "Status (to_do/doing/done): ")
		updates["status"] = Status(statusStr)
	case "5":
		deadlineStr := readInput(reader, "New deadline: ")
		if deadline, err := parseTime(deadlineStr); err == nil {
			updates["deadline"] = deadline
		}
	}

	if taskList.UpdateTask(id, updates) {
		fmt.Println("✅ Updated!")
	}
}
