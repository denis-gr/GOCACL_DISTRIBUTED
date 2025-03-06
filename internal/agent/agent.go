package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/denis-gr/GOCACL_DISTRIBUTED/internal/orchestrator"
)

func Worker(delay_ms int64, url string) {
	for {
		nextRun := time.Now().Add(time.Duration(delay_ms) * time.Millisecond)
		task := getTask(url)
		if task != nil {
			result := performTask(task)
			err := sendResult(result, url)
			if err != nil {
				log.Println("Error sending result:", err)
			}
		}
		time.Sleep(time.Until(nextRun))
	}
}

func getTask(url string) *orchestrator.Task {
	resp, err := http.Get(url)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	var task orchestrator.TaskResponse
	err = json.NewDecoder(resp.Body).Decode(&task)
	if err != nil {
		return nil
	}

	return &(task.Task)
}

func performTask(task *orchestrator.Task) *orchestrator.TaskResultRequest {
	wait := time.Now().Add(time.Duration(task.OperationTime) * time.Millisecond)

	var result float64
	switch task.Operation {
	case "+":
		result = task.Arg1 + task.Arg2
	case "-":
		result = task.Arg1 - task.Arg2
	case "*":
		result = task.Arg1 * task.Arg2
	case "/":
		if task.Arg2 != 0 {
			result = task.Arg1 / task.Arg2
		} else {
			result = 0
		}
	default:
		result = 0
	}

	time.Sleep(time.Until(wait))

	return &orchestrator.TaskResultRequest{
		ID:     task.ID,
		Result: result,
	}
}

func sendResult(result *orchestrator.TaskResultRequest, url string) error {
	data, err := json.Marshal(result)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
