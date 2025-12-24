package main

import (
	"log"
	"os"

	"github.com/ayushgarg/mcp-planka/internal/mcp"
	"github.com/ayushgarg/mcp-planka/internal/planka"
)

func main() {
	// Check if we should run tests instead
	if len(os.Args) > 1 && os.Args[1] == "test" {
		RunTests()
		return
	}

	// Get configuration from environment variables
	plankaURL := os.Getenv("PLANKA_URL")
	if plankaURL == "" {
		log.Fatal("PLANKA_URL environment variable is required")
	}

	var client *planka.Client
	var err error

	// Try token authentication first, then username/password
	plankaToken := os.Getenv("PLANKA_TOKEN")
	if plankaToken != "" {
		client = planka.NewClient(plankaURL, plankaToken)
	} else {
		// Try username/password authentication
		username := os.Getenv("PLANKA_USERNAME")
		password := os.Getenv("PLANKA_PASSWORD")
		if username == "" || password == "" {
			log.Fatal("Either PLANKA_TOKEN or both PLANKA_USERNAME and PLANKA_PASSWORD environment variables are required")
		}
		client, err = planka.NewClientWithPassword(plankaURL, username, password)
		if err != nil {
			log.Fatalf("Failed to authenticate with username/password: %v", err)
		}
		log.Println("Successfully authenticated with username/password")
	}

	// Initialize MCP server
	server := mcp.NewServer(client)

	// Start the MCP server
	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start MCP server: %v", err)
	}
}

