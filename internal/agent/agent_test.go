package agent

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/denis-gr/GOCACL_DISTRIBUTED/internal/orchestrator"
)

func TestPerformTask(t *testing.T) {
	tests := []struct {
		task     orchestrator.Task
		expected float64
	}{
		{orchestrator.Task{Operation: "+", Arg1: 1, Arg2: 1, OperationTime: 100}, 2},
		{orchestrator.Task{Operation: "-", Arg1: 2, Arg2: 1, OperationTime: 100}, 1},
		{orchestrator.Task{Operation: "*", Arg1: 2, Arg2: 2, OperationTime: 100}, 4},
		{orchestrator.Task{Operation: "/", Arg1: 4, Arg2: 2, OperationTime: 100}, 2},
		{orchestrator.Task{Operation: "/", Arg1: 4, Arg2: 0, OperationTime: 100}, 0},
	}

	for _, test := range tests {
		start := time.Now()
		result := performTask(&test.task)
		duration := time.Since(start)

		if result.Result != test.expected {
			t.Errorf("expected %f, got %f", test.expected, result.Result)
		}

		if duration < time.Duration(test.task.OperationTime)*time.Millisecond {
			t.Errorf("task completed too quickly, expected at least %dms, got %dms", test.task.OperationTime, duration.Milliseconds())
		}
	}

	// Несуществующая операция
	task := orchestrator.Task{Operation: "invalid", Arg1: 1, Arg2: 1, OperationTime: 100}
	result := performTask(&task)
	if result.Result != 0 {
		t.Errorf("expected 0, got %f", result.Result)
	}

	// Тесты с нулевым ожианием
	task = orchestrator.Task{Operation: "+", Arg1: 1, Arg2: 1, OperationTime: 0}
	start := time.Now()
	result = performTask(&task)
	duration := time.Since(start)

	if result.Result != 2 {
		t.Errorf("expected 2, got %f", result.Result)
	}

	if duration > time.Duration(task.OperationTime)*time.Millisecond {
		t.Errorf("task took too long, expected less than %dms, got %dms", task.OperationTime, duration.Milliseconds())
	}
}

func TestGetTask(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"task": {"ID": "1", "Operation": "+", "Arg1": 1, "Arg2": 1, "OperationTime": 100}}`))
	}))
	defer server.Close()

	task := getTask(server.URL)
	if task == nil {
		t.Fatal("expected task, got nil")
	}

	if task.Operation != "+" || task.Arg1 != 1 || task.Arg2 != 1 {
		t.Errorf("unexpected task: %+v", task)
	}

	// Тесты с неверными URL
	task = getTask("http://invalid-url")
	if task != nil {
		t.Fatal("expected nil, got task")
	}

	// Тесты с неверным JSON
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`invalid json`))
	}))
	defer server.Close()

	task = getTask(server.URL)
	if task != nil {
		t.Fatal("expected nil, got task")
	}

	// Тесты с пустым URL
	task = getTask("")
	if task != nil {
		t.Fatal("expected nil, got task")
	}
}

func TestSendResult(t *testing.T) {
	var received orchestrator.TaskResultRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatal(err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	result := &orchestrator.TaskResultRequest{ID: "1", Result: 2}
	err := sendResult(result, server.URL)
	if err != nil {
		t.Fatal(err)
	}

	if received.ID != result.ID || received.Result != result.Result {
		t.Errorf("unexpected result: %+v", received)
	}

	// Тесты с неверным URL
	err = sendResult(&orchestrator.TaskResultRequest{ID: "1", Result: 2}, "http://invalid-url")
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	// Тесты с ошибкой сервера
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	err = sendResult(&orchestrator.TaskResultRequest{ID: "1", Result: 2}, server.URL)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	// Тесты с неверным JSON
	err = sendResult(&orchestrator.TaskResultRequest{ID: "1", Result: 2}, "")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestWorker(t *testing.T) {
	taskServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"task": {"ID": "1", "Operation": "+", "Arg1": 1, "Arg2": 1, "OperationTime": 100}}`))
	}))
	defer taskServer.Close()

	resultServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer resultServer.Close()

	go Worker(1000, taskServer.URL)
	time.Sleep(2 * time.Second)

	// Тесты с неверным URL
	go Worker(1000, "http://invalid-url")
	time.Sleep(2 * time.Second)

	// Тесты с нулевой задержкой
	taskServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"task": {"ID": "1", "Operation": "+", "Arg1": 1, "Arg2": 1, "OperationTime": 100}}`))
	}))
	defer taskServer.Close()

	resultServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer resultServer.Close()

	go Worker(0, taskServer.URL)
	time.Sleep(2 * time.Second)

	// Тесты с ненулевой задержкой
	taskServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"task": {"ID": "1", "Operation": "+", "Arg1": 1, "Arg2": 1, "OperationTime": 100}}`))
	}))
	defer taskServer.Close()

	resultServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer resultServer.Close()

	go Worker(1, taskServer.URL)
	time.Sleep(2 * time.Second)
}
