package main

import (
	"fmt"
	"keeper/internal/command/agent"
	"log"
)

var (
	version     = "N/A"
	buildTime   = "N/A"
	buildCommit = "N/A"
)

func main() {
	fmt.Printf("Build version: %s\n", version)
	fmt.Printf("Build date: %s\n", buildTime)
	fmt.Printf("Build commit: %s\n", buildCommit)

	if err := agent.Execute(); err != nil {
		log.Fatalf("failed to run: %v", err)
	}
}
