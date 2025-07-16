package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
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

func (tl *TaskList) ShowReminders() {
	now := time.Now()

	overdue := tl.GetTasksByCondition(func(t Task) bool {
		return t.Status != Done && now.After(t.Deadline)
	})

	urgent := tl.GetTasksByCondition(func(t Task) bool {
		return t.Status != Done && t.Deadline.Sub(now) <= 24*time.Hour && t.Deadline.Sub(now) > 0
	})

	if len(overdue) > 0 {
		fmt.Printf("\n🚨 OVERDUE TASKS (%d):\n", len(overdue))
		for _, task := range overdue {
			fmt.Printf("❌ ID: %d | %s | Overdue: %s\n",
				task.ID, task.Name, formatDuration(now.Sub(task.Deadline)))
		}
	}

	if len(urgent) > 0 {
		fmt.Printf("\n⚠️  URGENT TASKS (%d):\n", len(urgent))
		for _, task := range urgent {
			fmt.Printf("⏰ ID: %d | %s | Due in: %s\n",
				task.ID, task.Name, formatDuration(task.Deadline.Sub(now)))
		}
	}

	if len(overdue) == 0 && len(urgent) == 0 {
		fmt.Println("\n✅ No urgent or overdue tasks!")
	}
}

func formatDuration(d time.Duration) string {
	if d < 0 {
		d = -d
	}

	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60

	switch {
	case hours >= 48:
		days := hours / 24
		remainingHours := hours % 24
		if remainingHours == 0 && minutes == 0 {
			return fmt.Sprintf("%d days", days)
		}
		return fmt.Sprintf("%d days %d hours", days, remainingHours)
	case hours > 0:
		if minutes == 0 {
			return fmt.Sprintf("%d hours", hours)
		}
		return fmt.Sprintf("%d hours %d minutes", hours, minutes)
	default:
		return fmt.Sprintf("%d minutes", minutes)
	}
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

func (tl *TaskList) Display() {
	statusOrder := []Status{ToDo, Doing, Done}

	for _, status := range statusOrder {
		fmt.Printf("\n--- %s ---\n", status)
		tasks := tl.GetTasksByCondition(func(t Task) bool { return t.Status == status })

		if len(tasks) == 0 {
			fmt.Println("No tasks")
			continue
		}

		for _, t := range tasks {
			icon := map[Status]string{ToDo: "📋", Doing: "🔄", Done: "✅"}[t.Status]

			if t.Status != Done {
				now := time.Now()
				if now.After(t.Deadline) {
					icon = "❌"
				} else if t.Deadline.Sub(now) <= 24*time.Hour {
					icon = "⚠️"
				}
			}

			fmt.Printf("%s ID: %d | %s | %s\n", icon, t.ID, t.Name, t.Category)
			fmt.Printf("   Description: %s\n", t.Description)
			fmt.Printf("   Deadline: %s\n", t.Deadline.Format("2006-01-02 15:04"))

			if t.Status == Done && t.CompletedAt != nil {
				fmt.Printf("   ✅ Completed: %s\n", t.CompletedAt.Format("2006-01-02 15:04"))
				duration := t.CompletedAt.Sub(t.CreatedAt)
				fmt.Printf("   ⏱️  Duration: %s\n", formatDuration(duration))
			} else if t.Status != Done {
				remaining := t.Deadline.Sub(time.Now())
				if remaining > 0 {
					fmt.Printf("   ⏰ Remaining: %s\n", formatDuration(remaining))
				} else {
					fmt.Printf("   ⚠️ Overdue: %s\n", formatDuration(-remaining))
				}
			}
			fmt.Println("---")
		}
	}
}

func (tl *TaskList) ShowStats() {
	total := len(tl.Tasks)
	completed := tl.GetTasksByCondition(func(t Task) bool { return t.Status == Done })
	onTime := tl.GetTasksByCondition(func(t Task) bool {
		return t.Status == Done && t.CompletedAt != nil && !t.CompletedAt.After(t.Deadline)
	})

	fmt.Printf("\n📊 STATISTICS:\n")
	fmt.Printf("Total: %d | Completed: %d | Pending: %d\n", total, len(completed), total-len(completed))

	if len(completed) > 0 {
		onTimeRate := float64(len(onTime)) / float64(len(completed)) * 100
		fmt.Printf("On-time rate: %.1f%%\n", onTimeRate)
	}
}

func parseTime(input string) (time.Time, error) {
	if input == "" {
		return time.Now().AddDate(1, 0, 0), nil
	}

	location := time.Now().Location()
	formats := []string{"2006-01-02 15:04", "2006-01-02", "01-02 15:04", "01-02"}

	for _, format := range formats {
		if format == "01-02" || format == "01-02 15:04" {
			input = fmt.Sprintf("%d-%s", time.Now().Year(), input)
			format = strings.Replace(format, "01-02", "2006-01-02", 1)
		}

		if t, err := time.ParseInLocation(format, input, location); err == nil {
			if !strings.Contains(format, "15:04") {
				t = t.Add(23*time.Hour + 59*time.Minute)
			}
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid date format")
}

func readInput(reader *bufio.Reader, prompt string) string {
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func readInt(reader *bufio.Reader, prompt string) (int, error) {
	input := readInput(reader, prompt)
	return strconv.Atoi(input)
}

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
			name := readInput(reader, "Task name: ")
			category := readInput(reader, "Category: ")
			description := readInput(reader, "Description: ")

			deadlineStr := readInput(reader, "Deadline (YYYY-MM-DD HH:MM): ")
			deadline, err := parseTime(deadlineStr)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}

			taskList.AddTask(name, category, description, deadline)
			fmt.Println("✅ Task added!")

		case "2":
			id, err := readInt(reader, "Task ID: ")
			if err != nil || taskList.FindTask(id) == nil {
				fmt.Println("❌ Task not found!")
				continue
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
