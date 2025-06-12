package main

import (
	"flag"
	"log"

	"mock-server/config"
	"mock-server/server"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to configuration file")
	port := flag.String("port", "8080", "port to run the server on")
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Validate configuration
	if err := config.ValidateConfig(cfg); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	// Create and start server
	srv := server.NewServer(cfg)
	log.Printf("Starting mock server on port %s", *port)
	if err := srv.Run(":" + *port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
