package main

import (
	"fmt"
	"time"
)

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
				if time.Until(t.Deadline) <= 0 {
					icon = "❌"
				} else if time.Until(t.Deadline) <= 24*time.Hour {
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
				remaining := time.Until(t.Deadline)
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

func (tl *TaskList) ShowReminders() {
	now := time.Now()

	overdue := tl.GetTasksByCondition(func(t Task) bool {
		return t.Status != Done && now.After(t.Deadline)
	})

	urgent := tl.GetTasksByCondition(func(t Task) bool {
		return t.Status != Done && time.Until(t.Deadline) <= 24*time.Hour && time.Until(t.Deadline) > 0
	})

	if len(overdue) > 0 {
		fmt.Printf("\n🚨 OVERDUE TASKS (%d):\n", len(overdue))
		for _, task := range overdue {
			fmt.Printf("❌ ID: %d | %s | Overdue: %s\n",
				task.ID, task.Name, formatDuration(time.Since(task.Deadline)))
		}
	}

	if len(urgent) > 0 {
		fmt.Printf("\n⚠️  URGENT TASKS (%d):\n", len(urgent))
		for _, task := range urgent {
			fmt.Printf("⏰ ID: %d | %s | Due in: %s\n",
				task.ID, task.Name, formatDuration(time.Until(task.Deadline)))
		}
	}

	if len(overdue) == 0 && len(urgent) == 0 {
		fmt.Println("\n✅ No urgent or overdue tasks!")
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
