package main

import (
	"bufio"
	"strings"
	"testing"
	"time"
)

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		dur      time.Duration
		expected string
	}{
		{time.Minute * 5, "5 minutes"},
		{time.Hour * 2, "2 hours"},
		{time.Hour*2 + time.Minute*30, "2 hours 30 minutes"},
		{time.Hour * 48, "2 days"},
		{time.Hour*49 + time.Minute*10, "2 days 1 hours"},
		{-time.Hour * 2, "2 hours"},
		{time.Hour * 24, "24 hours"},
	}

	for _, tt := range tests {
		got := formatDuration(tt.dur)
		if got != tt.expected {
			t.Errorf("formatDuration(%v) = %q, want %q", tt.dur, got, tt.expected)
		}
	}
}

func TestParseTime(t *testing.T) {
	now := time.Now()
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"2024-12-31 23:59", "2024-12-31 23:59"},
		{"2024-12-31", "2024-12-31 23:59"},
		{"12-31 23:59", "%d-12-31 23:59"},
		{"invalid", "invalid"},
	}

	for _, tt := range tests {
		got, err := parseTime(tt.input)
		if tt.input == "invalid" {
			if err == nil {
				t.Errorf("parseTime(%q) should fail", tt.input)
			}
			continue
		}
		if err != nil {
			t.Errorf("parseTime(%q) unexpected error: %v", tt.input, err)
			continue
		}

		if tt.input == "" {
			if got.Before(now) {
				t.Errorf("parseTime(\"\") got %v, want after now", got)
			}
		}
	}
}

func TestParseTime_InvalidFormats(t *testing.T) {
	invalidInputs := []string{
		"2024/12/31",
		"31-12-2024",
		"2024-13-01",
		"2024-12-32",
		"notadate",
		"2024-02-30",
	}
	for _, input := range invalidInputs {
		_, err := parseTime(input)
		if err == nil {
			t.Errorf("parseTime(%q) should return error", input)
		}
	}
}

func TestReadInput(t *testing.T) {
	input := "hello world"
	reader := bufio.NewReader(strings.NewReader(input))
	got := readInput(reader, "Please input：")
	if got != "hello world" {
		t.Errorf("readInput() = %q, want %q", got, "hello world")
	}
}

func TestReadInt(t *testing.T) {
	input := "42"
	reader := bufio.NewReader(strings.NewReader(input))
	got, err := readInt(reader, "Please input：")
	if err != nil {
		t.Fatalf("readInt() error: %v", err)
	}
	if got != 42 {
		t.Errorf("readInt() = %d, want 42", got)
	}

	// 非法输入
	reader = bufio.NewReader(strings.NewReader("abc\n"))
	_, err = readInt(reader, "Please input：")
	if err == nil {
		t.Errorf("readInt() should fail on non-integer input")
	}
}
func TestFormatDuration_EdgeCases(t *testing.T) {
	tests := []struct {
		dur      time.Duration
		expected string
	}{
		{0, "0 minutes"},
		{-time.Minute * 5, "5 minutes"},
		{-time.Hour*2 - time.Minute*30, "2 hours 30 minutes"},
		{time.Hour*24*3 + time.Hour*5, "3 days 5 hours"},
		{time.Hour*24*2 + time.Minute*59, "2 days 0 hours"},
	}
	for _, tt := range tests {
		got := formatDuration(tt.dur)
		if got != tt.expected {
			t.Errorf("formatDuration(%v) = %q, want %q", tt.dur, got, tt.expected)
		}
	}
}

func TestReadInput_TrimSpaces(t *testing.T) {
	input := "   test with spaces   \n"
	reader := bufio.NewReader(strings.NewReader(input))
	got := readInput(reader, "Prompt: ")
	if got != "test with spaces" {
		t.Errorf("readInput() = %q, want %q", got, "test with spaces")
	}
}

func TestReadInt_NegativeAndZero(t *testing.T) {
	reader := bufio.NewReader(strings.NewReader("-7\n"))
	got, err := readInt(reader, "Prompt: ")
	if err != nil {
		t.Fatalf("readInt() error: %v", err)
	}
	if got != -7 {
		t.Errorf("readInt() = %d, want -7", got)
	}

	reader = bufio.NewReader(strings.NewReader("0\n"))
	got, err = readInt(reader, "Prompt: ")
	if err != nil {
		t.Fatalf("readInt() error: %v", err)
	}
	if got != 0 {
		t.Errorf("readInt() = %d, want 0", got)
	}
}
