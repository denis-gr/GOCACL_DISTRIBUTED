// Package main содержит тесты для сервера orchestrator.
package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/denis-gr/GOCACL_DISTRIBUTED/internal/orchestrator"
)

func generateTestToken() string {
	token, _ := orchestrator.GenerateJWTToken("t", "t")
	return token
}

func TestMain(t *testing.T) {
	os.Setenv("ADDR", ":8081")

	router := orchestrator.NewRouter()

	ts := httptest.NewServer(router)
	defer ts.Close()

	req, err := http.NewRequest("GET", ts.URL+"/api/v1/expressions", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+generateTestToken())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to make GET request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestMain_NoAddr(t *testing.T) {
	os.Unsetenv("ADDR")

	router := orchestrator.NewRouter()

	ts := httptest.NewServer(router)
	defer ts.Close()

	req, err := http.NewRequest("GET", ts.URL+"/api/v1/expressions", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+generateTestToken())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to make GET request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestServerStart(t *testing.T) {
	httpAddr := ":8180"
	grpcAddr := ":8191"

	go func() {
		err := orchestrator.StartServer(httpAddr, grpcAddr)
		if err != nil {
			t.Fatalf("Failed to start server: %v", err)
		}
	}()

	time.Sleep(1 * time.Second)

	req, err := http.NewRequest("GET", "http://localhost:8180/api/v1/expressions", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+generateTestToken())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to make GET request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestServerStart_NoAddr(t *testing.T) {
	os.Unsetenv("ADDR")

	go func() {
		main()
	}()

	time.Sleep(1 * time.Second)

	req, err := http.NewRequest("GET", "http://localhost:8080/api/v1/expressions", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+generateTestToken())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to make GET request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
}
