package orchestrator

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStartCalculationHandler(t *testing.T) {
	reqBody := `{"expression":"2+2"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate", bytes.NewBufferString(reqBody))
	w := httptest.NewRecorder()

	StartCalculationHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("expected status %v, got %v", http.StatusCreated, resp.StatusCode)
	}

	var response StartCalCulationResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Errorf("error decoding response: %v", err)
	}

	if response.Id == "" {
		t.Errorf("expected non-empty id")
	}
}

func TestGetExpressionsListHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/expressions", nil)
	w := httptest.NewRecorder()

	GetExpressionsListHandler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %v, got %v", http.StatusOK, resp.StatusCode)
	}

	var response GetExpressionsListResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Errorf("error decoding response: %v", err)
	}
}
