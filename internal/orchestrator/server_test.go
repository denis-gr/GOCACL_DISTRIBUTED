// Package orchestrator содержит тесты для пакета orchestrator.
package orchestrator

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCalculateHandler(t *testing.T) {
	router := NewRouter()
	reqBody, _ := json.Marshal(CalculateRequest{Expression: "2+2"})
	req, _ := http.NewRequest("POST", "/api/v1/calculate", bytes.NewBuffer(reqBody))
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
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		var createRes CalculateResponse
		err := json.NewDecoder(rr.Body).Decode(&createRes)
		if err != nil {
			t.Errorf("error decoding response: %v", err)
		}

		// Используем полученный ID для запроса
		req, _ = http.NewRequest("GET", "/api/v1/expressions/"+createRes.ID, nil)
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
	router := NewRouter()
	req, _ := http.NewRequest("GET", "/internal/task", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %v, got %v", http.StatusOK, rr.Code)
	}
	var res TaskResponse
	err := json.NewDecoder(rr.Body).Decode(&res)
	if err != nil {
		t.Errorf("error decoding response: %v", err)
	}
}

func TestPostTaskResultHandler(t *testing.T) {
	router := NewRouter()
	// Сначала получаем задачу, чтобы получить ее ID
	req, _ := http.NewRequest("GET", "/internal/task", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	var taskRes TaskResponse
	_ = json.NewDecoder(rr.Body).Decode(&taskRes)

	// Используем полученный ID для отправки результата
	reqBody, _ := json.Marshal(TaskResultRequest{ID: taskRes.Task.ID, Result: 4})
	req, _ = http.NewRequest("POST", "/internal/task", bytes.NewBuffer(reqBody))
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %v, got %v", http.StatusOK, rr.Code)
	}
}

func TestGetTasksHandler(t *testing.T) {
	router := NewRouter()
	req, _ := http.NewRequest("GET", "/internal/tasks", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %v, got %v", http.StatusOK, rr.Code)
	}
	var res TaskFullResponse
	err := json.NewDecoder(rr.Body).Decode(&res)
	if err != nil {
		t.Errorf("error decoding response: %v", err)
	}
}

func TestCalculateHandlerInvalidRequest(t *testing.T) {
	router := NewRouter()
	req, _ := http.NewRequest("POST", "/api/v1/calculate", bytes.NewBuffer([]byte("invalid")))
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected status %v, got %v", http.StatusUnprocessableEntity, rr.Code)
	}
}

func TestPostTaskResultHandlerInvalidRequest(t *testing.T) {
	router := NewRouter()
	req, _ := http.NewRequest("POST", "/internal/task", bytes.NewBuffer([]byte("invalid")))
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected status %v, got %v", http.StatusUnprocessableEntity, rr.Code)
	}
}
