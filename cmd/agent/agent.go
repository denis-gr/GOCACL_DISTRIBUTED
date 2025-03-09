// Package main запускает агент, который выполняет задачи, полученные от оркестратора.
package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/denis-gr/GOCACL_DISTRIBUTED/internal/agent"
)

func main() {
	computingPower, err := strconv.Atoi(os.Getenv("COMPUTING_POWER"))
	if err != nil {
		computingPower = 1
	}
	delayMs, err := strconv.ParseInt(os.Getenv("DELAY_MS"), 10, 64)
	if err != nil {
		delayMs = 1000
	}
	url := os.Getenv("TASK_URL")
	if url == "" {
		url = "http://localhost:8080/internal/task"
	}

	fmt.Printf("Starting %d workers with delay %d ms\n, orchestrator url is %s", computingPower, delayMs, url)

	for i := 0; i < computingPower; i++ {
		go agent.Worker(delayMs, url)
	}

	select {}
}
