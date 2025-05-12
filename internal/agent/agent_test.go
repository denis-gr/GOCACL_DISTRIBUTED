// Package agent содержит тесты для пакета agent.
package agent

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	pb "github.com/denis-gr/GOCACL_DISTRIBUTED/internal/gen"
	"google.golang.org/grpc"
)

// Mock для pb.OrchestratorServiceClient
type mockOrchestratorServiceClient struct {
	pb.UnimplementedOrchestratorServiceServer
}

func (m *mockOrchestratorServiceClient) GetTask(ctx context.Context, req *pb.Empty, opts ...grpc.CallOption) (*pb.TaskResponse, error) {
	return &pb.TaskResponse{Task: &pb.Task{Id: "1", Operation: "+", Arg1: 1, Arg2: 1, OperationTime: 100}}, nil
}

func (m *mockOrchestratorServiceClient) SendResult(ctx context.Context, req *pb.TaskResultRequest, opts ...grpc.CallOption) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

func TestPerformTask(t *testing.T) {
	tests := []struct {
		task     pb.Task
		expected float64
	}{
		{pb.Task{Operation: "+", Arg1: 1, Arg2: 1, OperationTime: 100}, 2},
		{pb.Task{Operation: "-", Arg1: 2, Arg2: 1, OperationTime: 100}, 1},
		{pb.Task{Operation: "*", Arg1: 2, Arg2: 2, OperationTime: 100}, 4},
		{pb.Task{Operation: "/", Arg1: 4, Arg2: 2, OperationTime: 100}, 2},
		{pb.Task{Operation: "/", Arg1: 4, Arg2: 0, OperationTime: 100}, 0},
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
	task := pb.Task{Operation: "invalid", Arg1: 1, Arg2: 1, OperationTime: 100}
	result := performTask(&task)
	if result.Result != 0 {
		t.Errorf("expected 0, got %f", result.Result)
	}

	// Тесты с нулевым ожиданием
	task = pb.Task{Operation: "+", Arg1: 1, Arg2: 1, OperationTime: 0}
	start := time.Now()
	result = performTask(&task)
	duration := time.Since(start)

	if result.Result != 2 {
		t.Errorf("expected 2, got %f", result.Result)
	}

	if duration > time.Duration(task.OperationTime)*time.Millisecond+10*time.Millisecond {
		t.Errorf("task took too long, expected less than %dms, got %dms", task.OperationTime, duration.Milliseconds())
	}
}

func TestGetTask(t *testing.T) {
	client := &mockOrchestratorServiceClient{}
	task := getTask(client)
	if task == nil || task.Id != "1" {
		t.Errorf("unexpected task: %+v", task)
	}
}

func TestSendResult(t *testing.T) {
	client := &mockOrchestratorServiceClient{}
	result := &pb.TaskResultRequest{Id: "1", Result: 2}
	err := sendResult(client, result)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestWorker(t *testing.T) {
	taskServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"task": {"ID": "1", "Operation": "+", "Arg1": 1, "Arg2": 1, "OperationTime": 100}}`))
		if err != nil {
			t.Fatal(err)
		}
	}))
	defer taskServer.Close()

	resultServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer resultServer.Close()

	go Worker(1000, taskServer.URL)
	time.Sleep(2 * time.Second)

	// Тесты с неверным URL
	go Worker(1000, "http://invalid-url")
	time.Sleep(2 * time.Second)

	// Тесты с нулевой задержкой
	taskServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"task": {"ID": "1", "Operation": "+", "Arg1": 1, "Arg2": 1, "OperationTime": 100}}`))
		if err != nil {
			t.Fatal(err)
		}
	}))
	defer taskServer.Close()

	resultServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer resultServer.Close()

	go Worker(0, taskServer.URL)
	time.Sleep(2 * time.Second)

	// Тесты с ненулевой задержкой
	taskServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"task": {"ID": "1", "Operation": "+", "Arg1": 1, "Arg2": 1, "OperationTime": 100}}`))
		if err != nil {
			t.Fatal(err)
		}
	}))
	defer taskServer.Close()

	resultServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer resultServer.Close()

	go Worker(1, taskServer.URL)
	time.Sleep(2 * time.Second)
}
