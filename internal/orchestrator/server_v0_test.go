// Package orchestrator содержит тесты для пакета orchestrator.
package orchestrator

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCalculateHandlerV0(t *testing.T) {
	router := NewRouter()
	reqBody, _ := json.Marshal(CalculateRequest{Expression: "2+2"})
	req, _ := http.NewRequest("POST", "/api/v0/calculate", bytes.NewBuffer(reqBody))
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

func TestGetExpressionsHandlerV0(t *testing.T) {
	router := NewRouter()
	req, _ := http.NewRequest("GET", "/api/v0/expressions", nil)
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

func TestGetExpressionByIDHandlerV0(t *testing.T) {
	router := NewRouter()
	expressions := []string{"2+2", "5-3", "4*3", "8/2", "8/0"}

	for _, expr := range expressions {
		// Сначала создаем выражение, чтобы получить его ID
		reqBody, _ := json.Marshal(CalculateRequest{Expression: expr})
		req, _ := http.NewRequest("POST", "/api/v0/calculate", bytes.NewBuffer(reqBody))
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		var createRes CalculateResponse
		err := json.NewDecoder(rr.Body).Decode(&createRes)
		if err != nil {
			t.Errorf("error decoding response: %v", err)
		}

		// Используем полученный ID для запроса
		req, _ = http.NewRequest("GET", "/api/v0/expressions/"+createRes.ID, nil)
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

func TestGetTaskHandlerV0(t *testing.T) {
	router := NewRouter()
	req, _ := http.NewRequest("GET", "/api/v0/task", nil)
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

func TestPostTaskResultHandlerV0(t *testing.T) {
	router := NewRouter()
	// Сначала получаем задачу, чтобы получить ее ID
	req, _ := http.NewRequest("GET", "/api/v0/task", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	var taskRes TaskResponse
	_ = json.NewDecoder(rr.Body).Decode(&taskRes)

	// Используем полученный ID для отправки результата
	reqBody, _ := json.Marshal(TaskResultRequest{ID: taskRes.Task.ID, Result: 4})
	req, _ = http.NewRequest("POST", "/api/v0/task", bytes.NewBuffer(reqBody))
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %v, got %v", http.StatusOK, rr.Code)
	}
}

func TestGetTasksHandlerV0(t *testing.T) {
	router := NewRouter()
	req, _ := http.NewRequest("GET", "/api/v0/tasks", nil)
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

func TestCalculateHandlerInvalidRequestV0(t *testing.T) {
	router := NewRouter()
	req, _ := http.NewRequest("POST", "/api/v0/calculate", bytes.NewBuffer([]byte("invalid")))
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected status %v, got %v", http.StatusUnprocessableEntity, rr.Code)
	}
}

func TestPostTaskResultHandlerInvalidRequestV0(t *testing.T) {
	router := NewRouter()
	req, _ := http.NewRequest("POST", "/api/v0/task", bytes.NewBuffer([]byte("invalid")))
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected status %v, got %v", http.StatusUnprocessableEntity, rr.Code)
	}
}
