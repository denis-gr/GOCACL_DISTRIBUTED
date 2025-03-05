package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Task struct {
	Id            string    `json:"id"`
	Arg1          float64   `json:"arg1"`
	Arg2          float64   `json:"arg2"`
	Operation     string    `json:"operation"`
	OperationTime time.Time `json:"operation_time"`
}

func main() {
	computingPower, err := strconv.Atoi(os.Getenv("COMPUTING_POWER"))
	if err != nil {
		log.Fatalf("Invalid COMPUTING_POWER: %v", err)
	}

	for i := 0; i < computingPower; i++ {
		go worker()
	}

	select {}
}

func worker() {
	for {
		task, err := getTask()
		if err != nil {
			log.Printf("Error getting task: %v", err)
			time.Sleep(1 * time.Second)
			continue
		}

		result := performTask(task)

		err = sendResult(task.Id, result)
		if err != nil {
			log.Printf("Error sending result: %v", err)
		}
	}
}

func getTask() (*Task, error) {
	resp, err := http.Get("http://localhost:8080/internal/task")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}

	var task struct {
		Task Task `json:"task"`
	}
	err = json.NewDecoder(resp.Body).Decode(&task)
	if err != nil {
		return nil, err
	}

	return &task.Task, nil
}

func performTask(task *Task) float64 {
	time.Sleep(time.Duration(task.OperationTime) * time.Millisecond)

	switch task.Operation {
	case "+":
		return task.Arg1 + task.Arg2
	case "-":
		return task.Arg1 - task.Arg2
	case "*":
		return task.Arg1 * task.Arg2
	case "/":
		return task.Arg1 / task.Arg2
	default:
		return 0
	}
}

func sendResult(taskId string, result float64) error {
	data := struct {
		Id     string  `json:"id"`
		Result float64 `json:"result"`
	}{
		Id:     taskId,
		Result: result,
	}

	body, err := json.Marshal(data)
	if err != nil {
		return err
	}

	resp, err := http.Post("http://localhost:8080/internal/task", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}

	return nil
}
