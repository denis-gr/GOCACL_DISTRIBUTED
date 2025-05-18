package main

import (
	"os"
	"strconv"
	"testing"
)

func TestMain(t *testing.T) {
	os.Setenv("COMPUTING_POWER", "2")
	os.Setenv("DELAY_MS", "500")
	os.Setenv("TASK_URL", "localhost:8092")

	go main()

	if computingPower, err := strconv.Atoi(os.Getenv("COMPUTING_POWER")); err != nil || computingPower != 2 {
		t.Errorf("Expected COMPUTING_POWER to be 2, got %d", computingPower)
	}
	if delayMs, err := strconv.ParseInt(os.Getenv("DELAY_MS"), 10, 64); err != nil || delayMs != 500 {
		t.Errorf("Expected DELAY_MS to be 500, got %d", delayMs)
	}
	if url := os.Getenv("TASK_URL"); url != "localhost:8092" {
		t.Errorf("Expected TASK_URL to be 'localhost:8092', got %s", url)
	}
}
