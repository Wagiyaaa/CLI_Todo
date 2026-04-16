package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Task struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	Done bool   `json:"done"`
}

func marshal_json(tasks *[]Task) error {
	jsonTask, json_encode_err := json.Marshal(&tasks)
	if json_encode_err != nil {
		fmt.Println(json_encode_err)
		return json_encode_err
	}
	json_write_to_file_err := os.WriteFile("tasks.json", jsonTask, 0644)
	if json_write_to_file_err != nil {
		fmt.Println(json_write_to_file_err)
		return json_write_to_file_err
	}
	return nil
}

func unmarshal_json(tasks *[]Task) error {
	json_data, read_file_err := os.ReadFile("tasks.json")

	if read_file_err != nil {
		return read_file_err
	}
	unmarshal_err := json.Unmarshal(json_data, &tasks)
	if unmarshal_err != nil {
		return unmarshal_err
	}
	return nil
}

func list_tasks() error {
	var tasks []Task
	unmarshal_json_err := unmarshal_json(&tasks)
	if unmarshal_json_err != nil {
		return unmarshal_json_err
	}
	for _, task := range tasks {
		fmt.Printf("Task ID: %d %v Done: %t\n", task.ID, task.Text, task.Done)
	}
	return nil
}

func add_task(text string) error {
	var tasks []Task

	unmarshal_err := unmarshal_json(&tasks)
	if unmarshal_err != nil {
		return unmarshal_err
	}

	var id int
	if len(tasks) == 0 {
		id = 1
	} else {
		id = tasks[len(tasks)-1].ID + 1
	}
	task := Task{
		ID:   id,
		Text: text,
		Done: false,
	}

	tasks = append(tasks, task)

	err := marshal_json(&tasks)
	if err != nil {
		return err
	}

	return nil
}

func initialize_file() error {
	if _, err := os.Stat("tasks.json"); err == nil {
		return err
	} else if errors.Is(err, os.ErrNotExist) {
		err := os.WriteFile("tasks.json", []byte("[]"), 0644)
		if err != nil {
			fmt.Println(err)
			return err
		}
		return err
	} else {
		fmt.Println(err)
		return err
	}
}

func main() {
	args := os.Args
	if len(args) > 1 {
		err := initialize_file()
		if err != nil {
			fmt.Println(err)
			return
		}
		if args[1] == "add" && len(args) == 3 {
			err := add_task(args[2])
			if err != nil {
				fmt.Println(err)
				return
			}
			return
		}
		if args[1] == "list" && len(args) == 2 {
			err := list_tasks()
			if err != nil {
				fmt.Println(err)
				return
			}
			return
		}
	}
	fmt.Println("No command specified...")
}
