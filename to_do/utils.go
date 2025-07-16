package main

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
	"time"
)

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
