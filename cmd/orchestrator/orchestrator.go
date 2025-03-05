package main

import (
	"log"
	"net/http"
	"os"

	"github.com/denis/repos/GOCACL_DISTRIBUTED/internal/orchestrator"
)

func main() {
	http.HandleFunc("/api/v1/calculate", orchestrator.StartCalculationHandler)
	http.HandleFunc("/api/v1/expressions", orchestrator.GetExpressionsListHandler)
	http.HandleFunc("/api/v1/expressions/", orchestrator.GetExpressionHandler)
	http.HandleFunc("/internal/task", orchestrator.GetTaskHandler)
	http.HandleFunc("/internal/task", orchestrator.SaveTaskResultHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
