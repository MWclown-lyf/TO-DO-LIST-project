# To-Do List Project

A comprehensive command-line to-do list application written in Go, featuring task management, deadline tracking, and priority reminders.

## 📁 Project Structure

```
TO-DO-LIST-project/
├── main.go      # Main application entry point and user interface
├── task.go      # Task data structures and type definitions
├── tasklist.go  # Core business logic for task management
├── display.go   # Display and output formatting functions
├── utils.go     # Utility functions for input handling and time parsing
├── go.mod       # Go module definition
└── README.md    # Project documentation
```

## 📋 File Descriptions

### `main.go`
- **Purpose**: Application entry point and interactive user interface
- **Functions**: 
  - Main program loop with menu system
  - User input handling for all operations
  - Interactive task creation and updates
- **Key Features**: Command-line interface with 6 main options

### `task.go`
- **Purpose**: Core data structures and type definitions
- **Contains**:
  - `Task` struct with ID, name, category, description, status, deadline, and timestamps
  - `Status` enum (to_do, doing, done)
  - Task-related constants and type definitions

### `tasklist.go`
- **Purpose**: Core business logic and task management operations
- **Functions**:
  - `AddTask()`: Create new tasks
  - `UpdateTask()`: Modify existing tasks
  - `FindTask()`: Search for specific tasks
  - `GetTasksByCondition()`: Filter tasks by custom criteria
  - `Save()/Load()`: JSON file persistence

### `display.go`
- **Purpose**: User interface and output formatting
- **Functions**:
  - `Display()`: Show tasks organized by status with visual indicators
  - `ShowReminders()`: Display urgent and overdue task alerts
  - `ShowStats()`: Present completion statistics and metrics
  - `formatDuration()`: Convert time durations to readable format

### `utils.go`
- **Purpose**: Utility functions for common operations
- **Functions**:
  - `readInput()`: Handle user input with prompts
  - `readInt()`: Parse integer input with error handling
  - `parseTime()`: Parse various date/time formats with timezone support

## 🚀 Installation

### Prerequisites
- Go 1.16 or later installed on your system

### Quick Start
```bash
# Clone the repository
git clone https://github.com/liangyifan/to-do-list.git
cd to-do-list

# Run the application
go run .
```

### Build Executable
```bash
# Build for your platform
go build -o todo .

# Run the executable
./todo        # Linux/macOS
todo.exe      # Windows
```

## 📖 Usage Guide

### Main Menu Options
When you run the application, you'll see:
```
1. Add Task      # Create a new task
2. Update Task   # Modify existing task properties
3. View Tasks    # Display all tasks organized by status
4. Reminders     # Show urgent and overdue tasks
5. Statistics    # View completion metrics
6. Exit          # Save and quit the application
```

### 1. Adding Tasks
- Enter task name, category, and description
- Set deadline using formats:
  - `YYYY-MM-DD HH:MM` (e.g., 2024-12-25 15:30)
  - `YYYY-MM-DD` (defaults to end of day)
  - `MM-DD HH:MM` (current year)
  - `MM-DD` (current year, end of day)

### 2. Updating Tasks
- Enter task ID to modify
- Choose field to update:
  - Name, Category, Description
  - Status (to_do → doing → done)
  - Deadline

### 3. Task Display Features
- **Visual Indicators**:
  - 📋 To-do tasks
  - 🔄 In-progress tasks
  - ✅ Completed tasks
  - ❌ Overdue tasks
  - ⚠️ Urgent tasks (due within 24 hours)

### 4. Reminders System
- **Overdue Tasks**: Past deadline with duration
- **Urgent Tasks**: Due within 24 hours
- **Smart Notifications**: Automatic alerts on startup

### 5. Statistics Tracking
- Total tasks created
- Completion rate
- On-time completion percentage
- Task duration analysis

## ⏰ Time Management Features

### Deadline Formats Supported
```bash
2024-12-25 15:30    # Full date and time
2024-12-25          # Date only (23:59 assumed)
12-25 15:30         # Month-day with time (current year)
12-25               # Month-day only (current year, 23:59)
```

### Duration Display
- Minutes: "45 minutes"
- Hours: "2 hours 30 minutes"
- Days: "3 days 5 hours" (48+ hours only)

## 💾 Data Persistence

- Tasks are automatically saved to `tasks.json`
- Data persists between application sessions
- JSON format allows easy backup and sharing

## 🎯 Key Features

- ✅ **Intuitive CLI Interface**: Easy-to-use menu system
- ⏰ **Smart Time Tracking**: Automatic deadline monitoring
- 📊 **Progress Analytics**: Detailed completion statistics
- 🔄 **Status Management**: Three-stage workflow (to_do → doing → done)
- 💾 **Persistent Storage**: JSON-based data saving
- 🌍 **Timezone Support**: Local timezone handling
- 📱 **Cross-Platform**: Runs on Windows, macOS, and Linux

## 🛠️ Development

### Project Architecture
- **Modular Design**: Separated concerns across multiple files
- **Clean Code**: Well-documented functions and clear naming
- **Error Handling**: Robust input validation and error recovery
- **Extensible**: Easy to add new features and functionality

### Example Workflow
```bash
go run .                    # Start application
# 1. Add Task → Create "Complete project"
# 2. Update Task → Set status to "doing"
# 3. View Tasks → See progress
# 4. Reminders → Check urgent items
# 5. Statistics → Review performance
# 6. Exit → Save and quit
```

## 📝 License

This project is open source and available under the MIT License.

## 🤝 Contributing

Contributions, issues, and feature requests are welcome! Feel free to check the issues page.

---

**Happy Task Managing! 🎉**
