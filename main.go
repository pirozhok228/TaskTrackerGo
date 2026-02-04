package main

import (
	"encoding/json"
	"fmt"
	"os"
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
	if len(os.Args) > 1 {
		if os.Args[1] == "add" {
			addTask(os.Args[2])
		}
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Получена ошибка", r)
		}
	}()
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

func addTask(desc string) {
	file := File(filename)
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
	task := Task{len(taskList) + 1, desc, "todo", "03.02.2026", "03.02.2026"}
	taskList = append(taskList, task)
	data_json, err := json.MarshalIndent(taskList, "", " ")
	if err != nil {
		panic("Ошибка при работе с json!")
	}
	file, err = os.Create(filename)
	_, err = file.Write(data_json)
	if err != nil {
		panic("Ошибка при изменении файла!")
	}
	fmt.Println("Задача добавлена!")
	defer file.Close()
}
