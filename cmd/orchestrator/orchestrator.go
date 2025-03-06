package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/denis-gr/GOCACL_DISTRIBUTED/internal/orchestrator"
)

func main() {
	addr := os.Getenv("ADDR")

	if addr == "" {
		addr = ":8080"
	}

	router := orchestrator.NewRouter()
	fmt.Println("Starting server at", addr)
	err := http.ListenAndServe(addr, router)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
