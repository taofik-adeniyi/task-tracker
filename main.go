package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

type TaskStatus string

const (
	Todo       TaskStatus = "todo"
	Done       TaskStatus = "done"
	InProgress TaskStatus = "in-progress"
)

// ValidateTaskStatus checks if a TaskStatus is valid
func ValidateTaskStatus(status TaskStatus) error {
	switch status {
	case Todo, InProgress, Done:
		return nil
	default:
		return errors.New("invalid status: must be 'todo', 'in-progress', or 'done'")
	}
}

var commands = [9]string{"add", "update", "delete", "mark-in-progress", "mark-done", "list", "list done", "list todo", "list in-progress"}

type Task struct {
	Id          string
	Description string
	Status      TaskStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func main() {
	res := os.Args
	fmt.Println("Welcome to Task Tracker")
	if len(res) < 2 {
		help()
		return
	}
	commandsLists := res[1:]
	command := strings.Join(commandsLists, " ")
	fmt.Println("the command passed is: ", command)
	commandCheck(command)
}

func commandCheck(command string) {
	switch command {
	case "add":
		addTask()
	case "update":
		updateTask()
	case "delete":
		deleteTask()
	case "mark-in-progress":
		markTask()
	case "mark-done":
		markTask()
	case "list":
		listTasks()
	case "list done":
		listDoneTasks()
	case "list in-progress":
		listInprogressTasks()
	case "list todo":
		listTodoTasks()
	default:
		help()
	}
}

func addTask() {
	fmt.Println("Add Task ")
}
func updateTask() {
	fmt.Println("Update Task")
}
func deleteTask() {
	fmt.Println("Delete tasks")
}
func listTasks() {}
func markTask() {
	fmt.Println("Mark a task as in-progress or done")
}
func listDoneTasks() {
	fmt.Println("list task marked as done")
}
func listTodoTasks() {
	fmt.Println("list all tasks marked as todo")
}
func listInprogressTasks() {
	fmt.Println("list all tasks marked as in progress")
}

func help() {
	fmt.Println("List of Commands")
	for key, value := range commands {
		fmt.Printf("%v: %v \n", key, value)
	}
	// 	# Adding a new task
	// task-cli add "Buy groceries"
	// # Output: Task added successfully (ID: 1)

	// # Updating and deleting tasks
	// task-cli update 1 "Buy groceries and cook dinner"
	// task-cli delete 1

	// # Marking a task as in progress or done
	// task-cli mark-in-progress 1
	// task-cli mark-done 1

	// # Listing all tasks
	// task-cli list

	// # Listing tasks by status
	// task-cli list done
	// task-cli list todo
	// task-cli list in-progress
}
