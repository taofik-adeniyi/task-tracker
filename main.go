package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strconv"
	"strings"
	"time"
)

type TaskStatus string

const (
	Todo       TaskStatus = "todo"
	Done       TaskStatus = "done"
	InProgress TaskStatus = "in-progress"
	filePath   string     = "tasks.json"
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

var commands = []string{"add", "update", "delete", "mark-in-progress", "mark-done", "list", "list done", "list todo", "list in-progress"}

// Define a map of commands to their corresponding functions
var commandMap = map[string]func(){
	"add":              addTask,
	"update":           updateTask,
	"delete":           deleteTask,
	"mark-in-progress": markTaskInProgress,
	"mark-done":        markTaskDone,
	"list":             listTasks,
	"list done":        listDoneTasks,
	"list in-progress": listInprogressTasks,
	"list todo":        listTodoTasks,
	"-V":               checkVersion,
	"--version":        checkVersion,
}

type Task struct {
	Id          int        `json:"id"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (t Task) taskInfo() {
	// fmt.Print
	// fmt.Printf("%#v\n", t)

	fmt.Printf("Task ID: %v \n", t.Id)
	fmt.Printf("Task Description: %v \n", t.Description)
	fmt.Printf("Task Status: %v \n", t.Status)
	fmt.Printf("Task Created on: %v \n", t.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Task Updated on: %v \n", t.CreatedAt.Format("2006-01-02 15:04:05"))
}

func checkVersion() {
	println("v0.0.1")
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
	// fmt.Println("the command passed is: ", command)
	commandCheck(command)
}

func commandCheck(command string) {

	// fn, exists := commandMap[command]
	// if exists {
	// 	fn()
	// } else {
	// 	help()
	// }

	if strings.Contains(command, "add") {
		addTask()
	} else if strings.Contains(command, "update") {
		updateTask()
	} else if strings.Contains(command, "delete") {
		deleteTask()
	} else if strings.Contains(command, "mark-in-progress") {
		markTaskInProgress()
	} else if strings.Contains(command, "mark-done") {
		markTaskDone()
	} else if strings.Contains(command, "list") {
		listTasks()
	} else if strings.Contains(command, "list done") {
		listDoneTasks()
	} else if strings.Contains(command, "list todo") {
		listTodoTasks()
	} else if strings.Contains(command, "list in-progress") {
		listInprogressTasks()
	} else {
		help()
	}
}

func checkIfJsonFileExists() (*os.File, error) {
	osFile, err := os.Open(filePath)

	if err != nil {
		fmt.Printf("the file path is not available %s \n", err)
		return nil, err
	} else {
		defer osFile.Close()
		return osFile, nil
	}
}
func createFile() (*os.File, error) {
	osFile, err := os.Create(filePath)
	if err != nil {
		return nil, err
	} else {
		defer osFile.Close()
		return osFile, nil
	}
}

type DefaultFileStruct struct {
	Tasks []Task `json:"tasks"`
}

func createFileIfNotExist() (*os.File, error) {
	osFile, err := os.Open(filePath)

	if err != nil {
		fmt.Printf("the file path is not available %s \n", err)
		createdFile, err := os.Create(filePath)
		if err != nil {
			fmt.Printf("Error creating file [%v] that does not exist error is %v \n", filePath, err)
			return nil, err
		}
		defer createdFile.Close()

		osFile = createdFile
		fmt.Println("Name of created file", createdFile.Name())

		encoder := json.NewEncoder(osFile)
		if err := encoder.Encode(DefaultFileStruct{Tasks: []Task{{
			Id:          1,
			Description: "description",
			Status:      Todo,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}}}); err != nil {
			fmt.Printf("Error writing task to file: %v\n", err)
			return nil, err
		}
		fmt.Println("Task written to file successfully")

	} else {
		defer osFile.Close()
		fmt.Println("Name of default file", osFile.Name())
		return osFile, nil
	}
	return osFile, nil
}
func addTask() {
	fmt.Println("Add Task ")
	fmt.Println("e.g: task-tracker add 'Buy groceries'")
	// task-cli add "Buy groceries"

	args := os.Args
	argsLen := len(args)
	if argsLen != 3 {
		fmt.Printf("Incorrect command %v\n", args[0:])
		fmt.Printf("The correct command should look like (task-tracker add 'task one')")
		return
	}
	taskDescription := args[2]
	fmt.Println("task description", taskDescription)
	//  checkIfJsonFileExists()
	// createFile()
	file, err := createFileIfNotExist()
	if err != nil {
		fmt.Println("err check", err.Error())
		return
	}
	fmt.Println("check return file name", file.Name())

	fsys := os.DirFS(".")

	reader, err := fs.ReadFile(fsys, filePath)
	if err != nil {
		// if err == io.EOF {
		// 	fmt.Println("End of file", err.Error())
		// 	return
		// }
		fmt.Println("error reading file", err.Error())
		return
	}
	var fileContent DefaultFileStruct
	// fmt.Println("file content from reader", reader)
	fmt.Println("file content from reader", string(reader))
	err = json.Unmarshal(reader, &fileContent)
	if err != nil {
		fmt.Println("error converting reader file to the default struct", err.Error())
	}
	fmt.Println("before file content struct", fileContent)
	newTask := Task{
		Id:          len(fileContent.Tasks) + 1,
		Description: taskDescription,
		Status:      Todo,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	// newTask.taskInfo()
	fileContent.Tasks = append(fileContent.Tasks, newTask)
	defaultTask := fileContent.Tasks[0]
	addedTask := fileContent.Tasks[1]
	fmt.Println("after file content struct>>")
	defaultTask.taskInfo()
	addedTask.taskInfo()

	jsonData, err := json.Marshal(fileContent)
	if err != nil {
		fmt.Println("Error marshaling struct to JSON:", err)
		return
	}
	// Print the JSON byte slice
	fmt.Println("JSON data:", string(jsonData))
	err = os.WriteFile(filePath, jsonData, 0644)

	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("Struct saved to file as JSON successfully!")

}
func updateTask() {
	fmt.Println("Update Task")
	fmt.Println("e.g: task-tracker update 1 'Buy groceries and cook dinner'")
	args := os.Args
	if len(args) != 4 {
		fmt.Printf("Incorrect command %v\n", args[0:])
		fmt.Printf("The correct command should look like (task-tracker update 1 'description to update to')")
		return
	}

	taskId := args[2]
	taskDescription := args[3]
	taskIdInt, err := strconv.Atoi(taskId)
	if err != nil {
		fmt.Println("Error converting string to int", err.Error())
		return
	}
	var fileContent DefaultFileStruct
	fsys := os.DirFS(".")
	content, err := fs.ReadFile(fsys, filePath)
	if err != nil {
		fmt.Println("err reading file", err.Error())
		return
	}
	err = json.Unmarshal(content, &fileContent)
	if err != nil {
		fmt.Println("error turning bytes into struct", err.Error())
		return
	}

	for i, _ := range fileContent.Tasks {
		if fileContent.Tasks[i].Id == taskIdInt {
			fileContent.Tasks[i].Description = taskDescription
			fileContent.Tasks[i].UpdatedAt = time.Now()
		}
	}
	fmt.Println("file contents", fileContent)

	conte, err := json.Marshal(fileContent)
	if err != nil {
		fmt.Println("err converting struct to bytes", err.Error())
		return
	}
	err = os.WriteFile(filePath, conte, 0644)
	if err != nil {
		fmt.Println("Error writing to file", err.Error())
	}
}
func deleteTask() {
	// task-cli delete 1
	fmt.Println("e.g: task-tracker delete 1")
	fmt.Println("deleting task func")
	args := os.Args
	argsLen := len(args)
	if argsLen != 3 {
		fmt.Printf("Incorrect command %v\n", args[0:])
		fmt.Printf("The correct command should look like (task-tracker delete 1)")
		return
	}

	// command := args[1]
	taskId := args[2]

	taskIdInt, err := strconv.Atoi(taskId)
	if err != nil {
		fmt.Println("Invalid task ID: must be a number")
		return
	}

	fsfs := os.DirFS(".")
	// file, err := os.Open(filePath)
	content, err := fs.ReadFile(fsfs, filePath)
	if err != nil {
		fmt.Println("Unable to read file")
		return
	}
	var fileContent DefaultFileStruct
	err = json.Unmarshal(content, &fileContent)
	if err != nil {
		fmt.Println("Unable to parse the file content to struct", err.Error())
		return
	}

	var newFileContent DefaultFileStruct
	for _, value := range fileContent.Tasks {
		if value.Id == taskIdInt {
			continue
		} else {
			newFileContent.Tasks = append(newFileContent.Tasks, value)
		}
	}
	fmt.Println("newFileContent", newFileContent)

	contentToSave, err := json.Marshal(newFileContent)
	if err != nil {
		fmt.Printf("Unable to marshall the struct %v \n, err: %v", newFileContent, err.Error())
	}

	err = os.WriteFile(filePath, contentToSave, 0644)
	if err != nil {
		fmt.Printf("unable to save file into %v, err: %v\n", filePath, err.Error())
		return
	}
	// if command == "" || taskId == "" {
	// 	fmt.Println("invalid command or no taskId")
	// 	return
	// }

	fmt.Println("Delete tasks")
}
func listTasks() {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file", err.Error())
		return
	}
	defer file.Close()

	fsys := os.DirFS(".")
	var tasks DefaultFileStruct
	bytes, err := fs.ReadFile(fsys, filePath)
	if err != nil {
		fmt.Printf("Error reading file: %v, error %v \n", filePath, err)
	}
	// fmt.Println(string(bytes))
	err = json.Unmarshal(bytes, &tasks)
	if err != nil {
		fmt.Println("err converting bytes to struct", err.Error())
	}

	fmt.Println("List of tasks")
	fmt.Println("")
	for key, value := range tasks.Tasks {
		fmt.Printf("Task %v \n", key+1)
		value.taskInfo()
		fmt.Println("")
		fmt.Println("#######################")
		fmt.Println("")
	}
}
func markTaskInProgress() {
	fmt.Println("Mark a task as in-progress or done")
}
func markTaskDone() {
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
