package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

var filename string = "database.json"

type Task struct {
	ID          int
	Description string
	Status      string
	CreatedAt   string
	UpdatedAt   string
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Получена ошибка:", r)
		}
	}()

	if len(os.Args) > 1 {
		if os.Args[1] == "add" && len(os.Args) == 3 {
			addTask(os.Args[2])
		} else if os.Args[1] == "update" && len(os.Args) == 4 {
			id, _ := strconv.Atoi(os.Args[2])
			updateTask(id, os.Args[3])
		} else if os.Args[1] == "delete" && len(os.Args) == 3 {
			id, _ := strconv.Atoi(os.Args[2])
			deleteTask(id)
		} else if os.Args[1] == "list" && len(os.Args) == 2 {
			getList("all")
		} else if os.Args[1] == "list" && len(os.Args) == 3 {
			getList(os.Args[2])
		} else if os.Args[1] == "print" && len(os.Args) == 3 {
			id, _ := strconv.Atoi(os.Args[2])
			printTask(id)
		} else if os.Args[1] == "mark" && len(os.Args) == 4 {
			id, _ := strconv.Atoi(os.Args[3])
			markTask(os.Args[2], id)
		} else {
			panic("Недостаточно аргументов!")
		}
	}
}

func File(filename string) *os.File {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		file, err := os.Create(filename)
		if err != nil {
			panic("Ошибка при создании файла!")
		}
		data, err := json.Marshal([]Task{})
		if err != nil {
			panic("Ошибка при работе с json!")
		}
		file.Write(data)
		file.Close()
	}
	file, err := os.OpenFile(filename, os.O_RDWR, 0666)
	if err != nil {
		panic("Ошибка при открытии файла!")
	}
	return file
}

func getTaskList(file *os.File) []Task {
	buf := make([]byte, 1000)
	n, err := file.Read(buf)
	if err != nil {
		panic("Ошибка чтения файла!")
	}
	var taskList []Task
	err = json.Unmarshal(buf[:n], &taskList)
	if err != nil {
		panic("Ошибка при работе с json!")
	}
	return taskList
}

func uploadTaskList(file *os.File, taskList []Task) {
	data_json, err := json.MarshalIndent(taskList, "", " ")
	if err != nil {
		panic("Ошибка при работе с json!")
	}
	file, err = os.Create(filename)
	_, err = file.Write(data_json)
	if err != nil {
		panic("Ошибка при изменении файла!")
	}
}

func addTask(desc string) {
	openfile := File(filename)
	taskList := getTaskList(openfile)
	date := time.Now().Format("02.01.2006")
	task := Task{len(taskList) + 1, desc, "todo", date, date}
	taskList = append(taskList, task)
	uploadTaskList(openfile, taskList)
	fmt.Println("Задача добавлена!")
	defer openfile.Close()
}

func updateTask(id int, desc string) {
	openfile := File(filename)
	taskList := getTaskList(openfile)
	if id > len(taskList) {
		panic("Нет задачи с таким id!")
	}
	date := time.Now().Format("02.01.2006")
	oldTask := taskList[id-1]
	taskList[id-1] = Task{id, desc, oldTask.Status, oldTask.CreatedAt, date}
	uploadTaskList(openfile, taskList)
	fmt.Println("Задача обновлена!")
	defer openfile.Close()
}

func deleteTask(id int) {
	openfile := File(filename)
	taskList := getTaskList(openfile)
	newTaskList := make([]Task, len(taskList)-1)
	for _, task := range taskList {
		if task.ID < id {
			newTaskList = append(newTaskList, task)
		}
		if task.ID > id {
			task.ID -= 1
			newTaskList = append(newTaskList, task)
		}
	}
	newTaskList = newTaskList[int(len(newTaskList)/2):]
	uploadTaskList(openfile, newTaskList)
	fmt.Println("Задача удалена!")
	defer openfile.Close()
}

func formatPrint(task Task) {
	status := ""
	switch task.Status {
	case "todo":
		status = "Нужно сделать"
	case "in-progress":
		status = "В процессе выполнения"
	case "done":
		status = "Сделано"
	}
	fmt.Printf("id: %d\nОписание: %s\nСтатус: %s\nСоздана: %s\nОбновлена: %s\n", task.ID, task.Description, status, task.CreatedAt, task.UpdatedAt)
	fmt.Println("__________________________")
	fmt.Println()
}

func getList(filter string) {
	openfile := File(filename)
	taskList := getTaskList(openfile)
	for _, task := range taskList {
		if filter == "all" {
			formatPrint(task)
		} else {
			if task.Status == filter {
				formatPrint(task)
			}
		}
	}
}

func printTask(id int) {
	openfile := File(filename)
	taskList := getTaskList(openfile)
	if id > len(taskList) {
		panic("Нет задачи с таким id!")
	}
	for _, task := range taskList {
		if task.ID == id {
			formatPrint(task)
			break
		}
	}
	defer openfile.Close()
}

func markTask(status string, id int) {
	openfile := File(filename)
	taskList := getTaskList(openfile)
	if id > len(taskList) {
		panic("Нет задачи с таким id!")
	}
	date := time.Now().Format("02.01.2006")
	oldTask := taskList[id-1]
	taskList[id-1] = Task{id, oldTask.Description, status, oldTask.CreatedAt, date}
	uploadTaskList(openfile, taskList)
	fmt.Println("Задача обновлена!")
	defer openfile.Close()
}
