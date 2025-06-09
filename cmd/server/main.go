package main

import (
	"fmt"
	"keeper/internal/command/server"
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

	if err := server.Execute(); err != nil {
		log.Fatalf("server exited with error: %v", err)
	}
}
