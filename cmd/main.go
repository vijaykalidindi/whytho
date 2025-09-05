package main

import (
	"log"
	"os"

	"github.com/vinamra28/operator-reviewer/internal/config"
	"github.com/vinamra28/operator-reviewer/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	srv := server.New(cfg)
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting GitLab MR Reviewer Bot on port %s", port)
	if err := srv.Start(":" + port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}