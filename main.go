package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
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
	Version    string     = "0.0.1"
)

func ValidateTaskStatus(status TaskStatus) error {
	switch status {
	case Todo, InProgress, Done:
		return nil
	default:
		return errors.New("invalid status: must be 'todo', 'in-progress', or 'done'")
	}
}

type Commands struct {
	Command     string
	Description string
}

var comamndsList []Commands = []Commands{
	{
		Command:     "add",
		Description: "add {{task description e.g 'Buy groceries'}}",
	},
	{
		Command:     "update",
		Description: "update id description  e.g update 1 'Buy groceries and cook dinner'",
	},
	{
		Command:     "delete",
		Description: "delete id e.g delete 1",
	},
	{
		Command:     "mark-in-progress",
		Description: "mark-in-progress id e.g mark-in-progress 1",
	},
	{
		Command:     "mark-done",
		Description: "mark-done id e.g mark-done 1",
	},
	{
		Command:     "list",
		Description: "list",
	},
	{
		Command:     "list done",
		Description: "list done",
	},
	{
		Command:     "list todo",
		Description: "list todo",
	},
	{
		Command:     "list in-progress",
		Description: "list in-progress",
	},
	{
		Command:     "-V",
		Description: "list -V",
	},
	{
		Command:     "--version",
		Description: "list --version",
	},
}

type Task struct {
	Id          int        `json:"id"`
	Description string     `json:"description"`
	Status      TaskStatus `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func userInputs() []string {
	args := os.Args
	return args
}

func (t Task) taskInfo() {

	fmt.Printf("Task ID: %v \n", t.Id)
	fmt.Printf("Task Description: %v \n", t.Description)
	fmt.Printf("Task Status: %v \n", t.Status)
	fmt.Printf("Task Created on: %v \n", t.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Task Updated on: %v \n", t.CreatedAt.Format("2006-01-02 15:04:05"))
}

func checkVersion() string {
	version := "task-tracker@" + Version
	return version
}
func displayVersion() {
	fmt.Println("task-tracker@" + Version)
}

func main() {

	_, err := createFileIfNotExist()
	if err != nil {
		fmt.Println("err check", err.Error())
		return
	}

	res := userInputs()
	fmt.Println("")
	if len(res) < 2 {
		help()
		return
	}
	commandsLists := res[1:]
	command := strings.Join(commandsLists, " ")
	commandCheck(command)
}

func commandCheck(command string) {
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
	} else if strings.Contains(command, "list done") {
		listDoneTasks()
	} else if strings.Contains(command, "list todo") {
		listTodoTasks()
	} else if strings.Contains(command, "list in-progress") {
		listInprogressTasks()
	} else if strings.Contains(command, "list") {
		listTasks()
	} else if strings.Contains(command, "help") {
		help()
	} else if strings.Contains(command, "-V") {
		displayVersion()
	} else if strings.Contains(command, "--version") {
		displayVersion()
	} else {
		fmt.Printf("task-tracker: '%v' is not a task-tracker command. See 'task-tracker help'.\n", command)
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
		return osFile, nil
	}
	return osFile, nil
}
func addTask() {
	args := userInputs()
	argsLen := len(args)
	if argsLen != 3 {
		fmt.Printf("Incorrect command %v\n", args[0:])
		fmt.Printf("The correct command should look like (task-tracker add 'task one')")
		return
	}
	taskDescription := args[2]

	fsys := os.DirFS(".")

	reader, err := fs.ReadFile(fsys, filePath)
	if err != nil {
		log.Fatal(err)
		return
	}
	var fileContent DefaultFileStruct
	err = json.Unmarshal(reader, &fileContent)
	if err != nil {
		log.Fatal(err)
	}
	taskId := len(fileContent.Tasks) + 1
	newTask := Task{
		Id:          taskId,
		Description: taskDescription,
		Status:      Todo,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	fileContent.Tasks = append(fileContent.Tasks, newTask)

	jsonData, err := json.Marshal(fileContent)
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(filePath, jsonData, 0644)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Task added successfully (ID: %v)\n", taskId)

}
func updateTask() {
	args := userInputs()
	if len(args) != 4 {
		fmt.Printf("Incorrect command %v\n", args[0:])
		fmt.Printf("The correct command should look like (task-tracker update 1 'description to update to task 2')")
		return
	}

	taskId := args[2]
	taskDescription := args[3]
	taskIdInt, err := strconv.Atoi(taskId)
	if err != nil {
		log.Fatal(err)
		return
	}
	var fileContent DefaultFileStruct
	fsys := os.DirFS(".")
	content, err := fs.ReadFile(fsys, filePath)
	if err != nil {
		log.Fatal(err)
		return
	}
	err = json.Unmarshal(content, &fileContent)
	if err != nil {
		log.Fatal(err)
	}

	for i, _ := range fileContent.Tasks {
		if fileContent.Tasks[i].Id == taskIdInt {
			fileContent.Tasks[i].Description = taskDescription
			fileContent.Tasks[i].UpdatedAt = time.Now()
		}
	}

	conte, err := json.Marshal(fileContent)
	if err != nil {
		log.Fatal(err)
		return
	}
	err = os.WriteFile(filePath, conte, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
func deleteTask() {
	args := userInputs()
	argsLen := len(args)
	if argsLen != 3 {
		errMsg := fmt.Sprintf("Incorrect command %v\nThe correct command should look like (task-tracker delete 1)", args[0:])
		log.Fatal(errMsg)
	}

	taskId := args[2]

	taskIdInt, err := strconv.Atoi(taskId)
	if err != nil {
		fmt.Println("Invalid task ID: must be a number")
		return
	}

	fsfs := os.DirFS(".")
	content, err := fs.ReadFile(fsfs, filePath)
	if err != nil {
		log.Fatal(err)
	}
	var fileContent DefaultFileStruct
	err = json.Unmarshal(content, &fileContent)
	if err != nil {
		log.Fatal(err)
	}

	var newFileContent DefaultFileStruct
	for _, value := range fileContent.Tasks {
		if value.Id == taskIdInt {
			continue
		} else {
			newFileContent.Tasks = append(newFileContent.Tasks, value)
		}
		if value.Id != taskIdInt {
			newFileContent.Tasks = append(newFileContent.Tasks, value)
		}
	}

	contentToSave, err := json.Marshal(newFileContent)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(filePath, contentToSave, 0644)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Task with id of %v deleted", taskId)
}
func listTasks() {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fsys := os.DirFS(".")
	var tasks DefaultFileStruct
	bytes, err := fs.ReadFile(fsys, filePath)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(bytes, &tasks)
	if err != nil {
		log.Fatal(err)
	}

	if len(tasks.Tasks) < 1 {
		fmt.Println("No tasks added yet use the add command to add a task")
		os.Exit(1)
	}

	fmt.Printf("List of tasks")
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
	args := userInputs()
	taskId := args[2]
	taskIdInt, err := strconv.Atoi(taskId)
	if err != nil {
		log.Fatal(err)
	}
	if len(args) != 3 {
		fmt.Printf("Incorrect command %v\n", args[0:])
		fmt.Printf("The correct command should look like (task-tracker mark-in-progress 1)")
		return
	}

	var defaultTaskJson DefaultFileStruct
	fsys := os.DirFS(".")
	fileContent, err := fs.ReadFile(fsys, filePath)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(fileContent, &defaultTaskJson)
	if err != nil {
		log.Fatal(err)
	}
	for index, _ := range defaultTaskJson.Tasks {
		if defaultTaskJson.Tasks[index].Id == taskIdInt {
			defaultTaskJson.Tasks[index].Status = InProgress
		}
	}

	updatedTasks, err := json.Marshal(defaultTaskJson)
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(filePath, updatedTasks, 0644)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Task with id of %v has been marked as %v \n", taskIdInt, InProgress)
}
func markTaskDone() {
	args := userInputs()

	if len(args) != 3 {
		fmt.Printf("Incorrect command %v\n", args[0:])
		fmt.Printf("The correct command should look like (task-tracker mark-done 1)")
		return
	}

	taskId := args[2]
	taskIdInt, err := strconv.Atoi(taskId)
	if err != nil {
		log.Fatal(err)
	}

	var fileDefault DefaultFileStruct
	fsys := os.DirFS(".")
	fileContent, err := fs.ReadFile(fsys, filePath)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(fileContent, &fileDefault)
	if err != nil {
		log.Fatal(err)
	}

	for index, _ := range fileDefault.Tasks {
		if fileDefault.Tasks[index].Id == taskIdInt {
			fileDefault.Tasks[index].Status = Done
		}
	}
	updatedContent, err := json.Marshal(fileDefault)
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(filePath, updatedContent, 0644)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Task with id of %v has been marked as %v\n", taskId, Done)
}
func listDoneTasks() {
	args := userInputs()

	if len(args) != 3 {
		fmt.Printf("Incorrect command %v\n", args[0:])
		fmt.Printf("The correct command should look like (task-tracker list done)")
		return
	}

	var tasksJson DefaultFileStruct
	fsys := os.DirFS(".")
	content, err := fs.ReadFile(fsys, filePath)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = json.Unmarshal(content, &tasksJson)
	if err != nil {
		log.Fatal(err)
		return
	}
	if len(tasksJson.Tasks) < 1 {
		fmt.Println("No tasks added yet use the add command to add a task")
		os.Exit(1)
	}

	var doneTasks []Task
	var doneTasksCount int
	for _, value := range tasksJson.Tasks {
		if value.Status == Done {
			doneTasksCount++
			doneTasks = append(doneTasks, value)
		}
	}
	if doneTasksCount < 1 {
		fmt.Println("No task has been marked as done")
		os.Exit(1)
	}
	fmt.Printf("List of all %v tasks:\n", Todo)
	for _, value := range doneTasks {
		value.taskInfo()
		fmt.Println("____________")
	}
}
func listTodoTasks() {

	args := userInputs()
	if len(args) != 3 {
		fmt.Printf("Incorrect command %v\n", args[0:])
		fmt.Printf("The correct command should look like (task-tracker list todo)")
		return
	}

	var defaultFileStructure DefaultFileStruct
	fsys := os.DirFS(".")
	file_contents, err := fs.ReadFile(fsys, filePath)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(file_contents, &defaultFileStructure)
	if err != nil {
		log.Fatal(err)
	}

	var todoTasks []Task
	var todoTasksCount int
	for _, value := range defaultFileStructure.Tasks {
		if value.Status == Todo {
			todoTasksCount++
			todoTasks = append(todoTasks, value)
		}
	}

	if todoTasksCount < 1 {
		fmt.Println("No task marked as todo")
		os.Exit(1)
	}

	fmt.Printf("List of all %v tasks:\n", Todo)
	for _, value := range todoTasks {
		value.taskInfo()
	}
}
func listInprogressTasks() {

	args := userInputs()
	if len(args) != 3 {
		fmt.Printf("Incorrect command %v\n", args[0:])
		fmt.Printf("The correct command should look like (task-tracker list in-progress)")
		return
	}

	var defaultFileStructure DefaultFileStruct
	fsys := os.DirFS(".")
	file_contents, err := fs.ReadFile(fsys, filePath)
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(file_contents, &defaultFileStructure)
	if err != nil {
		log.Fatal(err)
	}

	var inProgressTasks []Task
	var inProgressCount int
	for _, value := range defaultFileStructure.Tasks {
		if value.Status == InProgress {
			inProgressCount++
			inProgressTasks = append(inProgressTasks, value)
		}
	}
	if inProgressCount < 1 {
		fmt.Println("No task has been marked as in-progress")
		os.Exit(1)
	}

	for _, value := range inProgressTasks {
		value.taskInfo()
	}

}

func help() {
	v := checkVersion()
	fmt.Println("Available Commands for ", v)

	fmt.Println("")
	for key, value := range comamndsList {
		fmt.Printf("%v: %v \n", key+1, value.Description)
	}
}
