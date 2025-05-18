// Package main содержит точку входа для сервера orchestrator.
package main

import (
	"fmt"
	"os"

	"github.com/denis-gr/GOCACL_DISTRIBUTED/internal/orchestrator"
)

func main() {
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = "localhost:8080"
	}
	grpcAddr := os.Getenv("GRPC_ADDR")
	if grpcAddr == "" {
		grpcAddr = "localhost:8092"
	}
	err := orchestrator.StartServer(addr, grpcAddr)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
