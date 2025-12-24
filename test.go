package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/ayushgarg/mcp-planka/internal/planka"
)

// RunTests runs all the Planka API connection tests
func RunTests() {
	// Get configuration from environment variables
	baseURL := os.Getenv("PLANKA_URL")
	if baseURL == "" {
		baseURL = "https://planka.cosmicdragon.xyz"
		fmt.Printf("Using default URL: %s\n", baseURL)
	}

	username := os.Getenv("PLANKA_USERNAME")
	password := os.Getenv("PLANKA_PASSWORD")
	token := os.Getenv("PLANKA_TOKEN")

	var client *planka.Client
	var err error

	fmt.Println("=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=")
	fmt.Println("Testing Planka MCP Server")
	fmt.Println("=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=")

	// Test 1: Authentication
	fmt.Println("\n[Test 1] Authenticating...")
	if token != "" {
		client = planka.NewClient(baseURL, token)
		fmt.Println("✓ Using token authentication")
	} else if username != "" && password != "" {
		client, err = planka.NewClientWithPassword(baseURL, username, password)
		if err != nil {
			log.Fatalf("❌ Authentication failed: %v", err)
		}
		fmt.Println("✓ Successfully authenticated with username/password!")
	} else {
		log.Fatal("❌ Either PLANKA_TOKEN or both PLANKA_USERNAME and PLANKA_PASSWORD must be set")
	}

	// Test 2: Get current user
	fmt.Println("\n[Test 2] Getting current user...")
	user, err := client.GetMe()
	if err != nil {
		log.Fatalf("❌ Failed to get user info: %v", err)
	}
	fmt.Printf("✓ Logged in as: %s", user.Username)
	if user.Email != "" {
		fmt.Printf(" (%s)", user.Email)
	}
	fmt.Println()

	// Test 3: Get all projects
	fmt.Println("\n[Test 3] Getting all projects...")
	projects, err := client.GetProjects()
	if err != nil {
		log.Fatalf("❌ Failed to get projects: %v", err)
	}
	fmt.Printf("✓ Found %d project(s):\n", len(projects))
	for i, project := range projects {
		fmt.Printf("  %d. %s (ID: %s)\n", i+1, project.Name, project.ID)
	}

	if len(projects) == 0 {
		fmt.Println("⚠ No projects found, skipping further tests")
		return
	}

	// Test 4: Get a specific project
	fmt.Printf("\n[Test 4] Getting project '%s'...\n", projects[0].Name)
	project, err := client.GetProject(projects[0].ID)
	if err != nil {
		log.Printf("❌ Failed to get project: %v", err)
	} else {
		fmt.Printf("✓ Project retrieved: %s\n", project.Name)
	}

	// Test 5: Get boards for a project
	fmt.Printf("\n[Test 5] Getting boards for project '%s'...\n", projects[0].Name)
	boards, err := client.GetBoards(projects[0].ID)
	if err != nil {
		log.Printf("❌ Failed to get boards: %v", err)
	} else {
		fmt.Printf("✓ Found %d board(s):\n", len(boards))
		for i, board := range boards {
			fmt.Printf("  %d. %s (ID: %s)\n", i+1, board.Name, board.ID)
		}
	}

	if len(boards) == 0 {
		fmt.Println("⚠ No boards found, skipping further tests")
		return
	}

	// Test 6: Get a specific board
	fmt.Printf("\n[Test 6] Getting board '%s'...\n", boards[0].Name)
	board, err := client.GetBoard(boards[0].ID)
	if err != nil {
		log.Printf("❌ Failed to get board: %v", err)
	} else {
		fmt.Printf("✓ Board retrieved: %s\n", board.Name)
	}

	// Test 7: Get lists for a board
	fmt.Printf("\n[Test 7] Getting lists for board '%s'...\n", boards[0].Name)
	lists, err := client.GetLists(boards[0].ID)
	if err != nil {
		log.Printf("❌ Failed to get lists: %v", err)
	} else {
		fmt.Printf("✓ Found %d list(s):\n", len(lists))
		for i, list := range lists {
			fmt.Printf("  %d. %s (ID: %s)\n", i+1, list.Name, list.ID)
		}
	}

	if len(lists) == 0 {
		fmt.Println("⚠ No lists found, skipping further tests")
		return
	}

	// Test 8: Get cards for a list
	fmt.Printf("\n[Test 8] Getting cards for list '%s'...\n", lists[0].Name)
	cards, err := client.GetCards(lists[0].ID)
	if err != nil {
		log.Printf("❌ Failed to get cards: %v", err)
	} else {
		fmt.Printf("✓ Found %d card(s):\n", len(cards))
		if len(cards) > 0 {
			for i, card := range cards {
				fmt.Printf("  %d. %s (ID: %s)\n", i+1, card.Name, card.ID)
				if i >= 4 { // Limit output to first 5 cards
					fmt.Printf("  ... and %d more card(s)\n", len(cards)-5)
					break
				}
			}
		} else {
			fmt.Println("  (No cards in this list)")
		}
	}

	// Try to find a list with cards if the first one is empty
	if len(cards) == 0 && len(lists) > 1 {
		fmt.Printf("\n[Test 8b] Trying to find cards in other lists...\n")
		for _, list := range lists[1:] {
			cards, err = client.GetCards(list.ID)
			if err == nil && len(cards) > 0 {
				fmt.Printf("✓ Found %d card(s) in list '%s':\n", len(cards), list.Name)
				for i, card := range cards {
					fmt.Printf("  %d. %s (ID: %s)\n", i+1, card.Name, card.ID)
					if i >= 4 {
						fmt.Printf("  ... and %d more card(s)\n", len(cards)-5)
						break
					}
				}
				break
			}
		}
	}

	if len(cards) > 0 {
		// Test 9: Get a specific card
		fmt.Printf("\n[Test 9] Getting card '%s'...\n", cards[0].Name)
		card, err := client.GetCard(cards[0].ID)
		if err != nil {
			log.Printf("❌ Failed to get card: %v", err)
		} else {
			fmt.Printf("✓ Card retrieved: %s\n", card.Name)
			if card.Description != "" {
				desc := card.Description
				if len(desc) > 100 {
					desc = desc[:100] + "..."
				}
				fmt.Printf("  Description: %s\n", desc)
			}

			// Test 10: Get tasks for a card
			fmt.Printf("\n[Test 10] Getting tasks for card '%s'...\n", card.Name)
			tasks, err := client.GetTasks(card.ID)
			if err != nil {
				log.Printf("❌ Failed to get tasks: %v", err)
			} else {
				fmt.Printf("✓ Found %d task(s):\n", len(tasks))
				if len(tasks) > 0 {
					for i, task := range tasks {
						status := "incomplete"
						if task.IsCompleted {
							status = "completed"
						}
						fmt.Printf("  %d. %s (%s)\n", i+1, task.Name, status)
					}
				} else {
					fmt.Println("  (No tasks in this card)")
				}
			}

			// Test 11: Get comments for a card
			fmt.Printf("\n[Test 11] Getting comments for card '%s'...\n", card.Name)
			comments, err := client.GetComments(card.ID)
			if err != nil {
				log.Printf("❌ Failed to get comments: %v", err)
			} else {
				fmt.Printf("✓ Found %d comment(s)\n", len(comments))
				if len(comments) > 0 {
					for i, comment := range comments {
						text := comment.Text
						if len(text) > 50 {
							text = text[:50] + "..."
						}
						fmt.Printf("  %d. %s\n", i+1, text)
					}
				}
			}

			// Test 12: Get stopwatch for a card
			fmt.Printf("\n[Test 12] Getting stopwatch for card '%s'...\n", card.Name)
			stopwatch, err := client.GetStopwatch(card.ID)
			if err != nil {
				log.Printf("⚠ Failed to get stopwatch (may not exist): %v", err)
			} else {
				fmt.Printf("✓ Stopwatch retrieved\n")
				if stopwatch.StartedAt != nil {
					fmt.Printf("  Started at: %s\n", stopwatch.StartedAt.Format("2006-01-02 15:04:05"))
				}
				fmt.Printf("  Duration: %d seconds\n", stopwatch.Duration)
			}
		}
	}

	// Test 13: Test JSON marshaling (for MCP responses)
	fmt.Println("\n[Test 13] Testing JSON marshaling...")
	testData := map[string]interface{}{
		"projects": len(projects),
		"boards":    len(boards),
		"lists":     len(lists),
		"cards":     len(cards),
		"user":      user.Username,
	}
	jsonData, err := json.MarshalIndent(testData, "", "  ")
	if err != nil {
		log.Printf("❌ Failed to marshal JSON: %v", err)
	} else {
		fmt.Println("✓ JSON marshaling works:")
		fmt.Println(string(jsonData))
	}

	fmt.Println("\n" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=")
	fmt.Println("✓ All tests completed successfully!")
	fmt.Println("=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=" + "=")
}
