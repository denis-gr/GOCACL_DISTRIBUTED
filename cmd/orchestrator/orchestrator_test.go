package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/denis-gr/GOCACL_DISTRIBUTED/internal/orchestrator"
)

func TestMain(t *testing.T) {
	os.Setenv("ADDR", ":8081")

	router := orchestrator.NewRouter()

	ts := httptest.NewServer(router)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/v1/expressions")
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

	resp, err := http.Get(ts.URL + "/api/v1/expressions")
	if err != nil {
		t.Fatalf("Failed to make GET request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
}

func TestServerStart(t *testing.T) {
	os.Setenv("ADDR", ":8082")

	go func() {
		main()
	}()

	time.Sleep(1 * time.Second)

	resp, err := http.Get("http://localhost:8082/api/v1/expressions")
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

	resp, err := http.Get("http://localhost:8080/api/v1/expressions")
	if err != nil {
		t.Fatalf("Failed to make GET request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}
}
