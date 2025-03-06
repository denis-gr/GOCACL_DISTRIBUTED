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
	delay_ms, err := strconv.ParseInt(os.Getenv("DELAY_MS"), 10, 64)
	if err != nil {
		delay_ms = 1000
	}
	url := os.Getenv("TASK_URL")
	if url == "" {
		url = "http://localhost:8080/internal/task"
	}

	fmt.Printf("Starting %d workers with delay %d ms\n, orchestrator url is %s", computingPower, delay_ms, url)

	for i := 0; i < computingPower; i++ {
		go agent.Worker(delay_ms, url)
	}

	select {}
}
