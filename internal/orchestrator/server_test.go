// Package orchestrator содержит тесты для пакета orchestrator.
package orchestrator

import (
	"bytes"
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	pb "github.com/denis-gr/GOCACL_DISTRIBUTED/internal/gen"
	"google.golang.org/grpc"
)

func generateTestToken() string {
	token, _ := GenerateJWTToken("t", "t")
	return token
}

var grpcServer *grpc.Server

func startTestGRPCServer() string {
	listener, _ := net.Listen("tcp", "localhost:8092")

	grpcServer = grpc.NewServer()
	pb.RegisterOrchestratorServiceServer(grpcServer, &OrchestratorGRPCServer{})

	go func() {
		grpcServer.Serve(listener)
	}()

	return listener.Addr().String()
}

func stopTestGRPCServer() {
	if grpcServer != nil {
		grpcServer.Stop()
	}
}

func TestMain(m *testing.M) {
	startTestGRPCServer()
	defer stopTestGRPCServer()

	m.Run()
}

func TestCalculateHandler(t *testing.T) {
	router := NewRouter()
	reqBody, _ := json.Marshal(CalculateRequest{Expression: "2+2"})
	req, _ := http.NewRequest("POST", "/api/v1/calculate", bytes.NewBuffer(reqBody))
	req.Header.Set("Authorization", "Bearer "+generateTestToken())
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %v, got %v", http.StatusOK, rr.Code)
	}
	var res CalculateResponse
	err := json.NewDecoder(rr.Body).Decode(&res)
	if err != nil {
		t.Errorf("error decoding response: %v", err)
	}
	if res.ID == "" {
		t.Errorf("expected non-empty ID")
	}
}

func TestGetExpressionsHandler(t *testing.T) {
	router := NewRouter()
	req, _ := http.NewRequest("GET", "/api/v1/expressions", nil)
	req.Header.Set("Authorization", "Bearer "+generateTestToken())
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %v, got %v", http.StatusOK, rr.Code)
	}
	var res ExpressionsResponse
	err := json.NewDecoder(rr.Body).Decode(&res)
	if err != nil {
		t.Errorf("error decoding response: %v", err)
	}
}

func TestGetExpressionByIDHandler(t *testing.T) {
	router := NewRouter()
	expressions := []string{"2+2", "5-3", "4*3", "8/2", "8/0"}

	for _, expr := range expressions {
		// Сначала создаем выражение, чтобы получить его ID
		reqBody, _ := json.Marshal(CalculateRequest{Expression: expr})
		req, _ := http.NewRequest("POST", "/api/v1/calculate", bytes.NewBuffer(reqBody))
		req.Header.Set("Authorization", "Bearer "+generateTestToken())
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		var createRes CalculateResponse
		err := json.NewDecoder(rr.Body).Decode(&createRes)
		if err != nil {
			t.Errorf("error decoding response: %v", err)
		}

		// Используем полученный ID для запроса
		req, _ = http.NewRequest("GET", "/api/v1/expressions/"+createRes.ID, nil)
		req.Header.Set("Authorization", "Bearer "+generateTestToken())
		rr = httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status %v, got %v", http.StatusOK, rr.Code)
		}
		var res ExpressionResponse
		err = json.NewDecoder(rr.Body).Decode(&res)
		if err != nil {
			t.Errorf("error decoding response: %v", err)
		}
	}
}

func TestGetTaskHandler(t *testing.T) {
	conn, err := grpc.Dial("localhost:8092", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := pb.NewOrchestratorServiceClient(conn)
	res, err := client.GetTask(context.Background(), &pb.Empty{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.Task == nil {
		t.Errorf("expected a task, got nil")
	}
}

func TestPostTaskResultHandler(t *testing.T) {
	conn, err := grpc.Dial("localhost:8092", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := pb.NewOrchestratorServiceClient(conn)

	// Получаем задачу
	taskRes, err := client.GetTask(context.Background(), &pb.Empty{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Отправляем результат
	_, err = client.SendResult(context.Background(), &pb.TaskResultRequest{
		Id:     taskRes.Task.Id,
		Result: 4,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestCalculateHandlerInvalidRequest(t *testing.T) {
	router := NewRouter()
	req, _ := http.NewRequest("POST", "/api/v1/calculate", bytes.NewBuffer([]byte("invalid")))
	req.Header.Set("Authorization", "Bearer "+generateTestToken())
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected status %v, got %v", http.StatusUnprocessableEntity, rr.Code)
	}
}

func TestPostTaskResultHandlerInvalidRequest(t *testing.T) {
	conn, err := grpc.Dial("localhost:8092", grpc.WithInsecure())
	if err != nil {
		t.Fatalf("failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := pb.NewOrchestratorServiceClient(conn)

	// Отправляем некорректный результат
	_, err = client.SendResult(context.Background(), &pb.TaskResultRequest{
		Id:     "invalid-id",
		Result: 4,
	})
	if err == nil {
		t.Fatalf("expected an error, got nil")
	}
}

func TestRegisterUserHandler(t *testing.T) {
	router := NewRouter()

	// Создаем корректный запрос
	validUser := UserCreateForm{
		Username: "testuser173",
		Password: "password123",
	}
	validBody, _ := json.Marshal(validUser)
	req, _ := http.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(validBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("expected status %v, got %v", http.StatusCreated, rr.Code)
	}

	// Создаем некорректный запрос
	invalidBody := []byte("invalid-json")
	req, _ = http.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(invalidBody))
	req.Header.Set("Content-Type", "application/json")

	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected status %v, got %v", http.StatusUnprocessableEntity, rr.Code)
	}
}

func TestLoginUserHandler(t *testing.T) {
	router := NewRouter()

	validUser := UserCreateForm{
		Username: "testuser",
		Password: "password123",
	}
	validBody, _ := json.Marshal(validUser)
	req, _ := http.NewRequest("POST", "/api/v1/register", bytes.NewBuffer(validBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	// Создаем корректный запрос
	validLogin := UserLoginForm{
		Username: "testuser",
		Password: "password123",
	}
	validBody, _ = json.Marshal(validLogin)
	req, _ = http.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(validBody))
	req.Header.Set("Content-Type", "application/json")

	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %v, got %v", http.StatusOK, rr.Code)
	}

	// Создаем некорректный запрос
	invalidBody := []byte("invalid-json")
	req, _ = http.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(invalidBody))
	req.Header.Set("Content-Type", "application/json")

	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected status %v, got %v", http.StatusUnprocessableEntity, rr.Code)
	}
}
