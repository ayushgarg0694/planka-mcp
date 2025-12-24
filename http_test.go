// Package main contains integration tests for the MCP HTTP server.
// These tests communicate with an actual Planka instance via HTTP.
//
// To run these tests:
//   1. Start the MCP server in HTTP mode: ./mcp-planka --http --http-port 8080
//   2. Set environment variables if needed: BASE_URL=http://localhost:8080
//   3. Run tests: go test -v ./http_test.go
//
// The tests are idempotent and will clean up created resources.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

// Test configuration
var (
	baseURL        = getEnv("BASE_URL", "http://localhost:8080")
	mcpEndpoint    = baseURL + "/mcp"
	healthEndpoint = baseURL + "/health"
)

// E2E test resource tracking
type e2eResources struct {
	projectID string
	boardID   string
	listID    string
	cardIDs   []string
	taskIDs   []string
}

var e2eRes e2eResources

// JSON-RPC structures
type jsonRPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
	ID      interface{} `json:"id"`
}

type jsonRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	Result  interface{}     `json:"result,omitempty"`
	Error   *jsonRPCError   `json:"error,omitempty"`
	ID      interface{}     `json:"id,omitempty"`
}

type jsonRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type toolCallResult struct {
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func makeJSONRPCRequest(method string, params interface{}, id interface{}) (*jsonRPCResponse, error) {
	req := jsonRPCRequest{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
		ID:      id,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	resp, err := http.Post(mcpEndpoint, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	var jsonResp jsonRPCResponse
	if err := json.NewDecoder(resp.Body).Decode(&jsonResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &jsonResp, nil
}

func callTool(t *testing.T, toolName string, arguments map[string]interface{}) (map[string]interface{}, error) {
	params := map[string]interface{}{
		"name":      toolName,
		"arguments": arguments,
	}

	resp, err := makeJSONRPCRequest("tools/call", params, time.Now().UnixNano())
	if err != nil {
		return nil, err
	}

	if resp.Error != nil {
		return nil, fmt.Errorf("JSON-RPC error: %d - %s", resp.Error.Code, resp.Error.Message)
	}

	// Extract result content
	resultMap, ok := resp.Result.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid result format")
	}

	content, ok := resultMap["content"].([]interface{})
	if !ok || len(content) == 0 {
		return nil, fmt.Errorf("no content in result")
	}

	contentItem, ok := content[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid content format")
	}

	text, ok := contentItem["text"].(string)
	if !ok {
		return nil, fmt.Errorf("no text in content")
	}

	// Try parsing as array first (for get_projects, get_boards, etc.)
	var arrResult []interface{}
	if err := json.Unmarshal([]byte(text), &arrResult); err == nil {
		return map[string]interface{}{"items": arrResult}, nil
	}

	// Try parsing as object
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(text), &result); err != nil {
		// If parsing fails, it might be a plain string (e.g., delete operations)
		// Return a success indicator with the message
		return map[string]interface{}{
			"success": true,
			"message": text,
		}, nil
	}

	return result, nil
}

func extractID(result map[string]interface{}) string {
	if id, ok := result["id"].(string); ok {
		return id
	}
	return ""
}

// Test 1: Health Check
func TestHealth(t *testing.T) {
	resp, err := http.Get(healthEndpoint)
	if err != nil {
		t.Fatalf("Failed to connect to health endpoint: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}

	var healthResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&healthResp); err != nil {
		t.Fatalf("Failed to decode health response: %v", err)
	}

	if status, ok := healthResp["status"].(string); !ok || status != "ok" {
		t.Fatalf("Expected status 'ok', got %v", healthResp["status"])
	}

	t.Log("✓ Health check passed")
}

// Test 2: Initialize
func TestInitialize(t *testing.T) {
	params := map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"capabilities":    map[string]interface{}{},
		"clientInfo": map[string]interface{}{
			"name":    "test-client",
			"version": "1.0.0",
		},
	}

	resp, err := makeJSONRPCRequest("initialize", params, 1)
	if err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}

	if resp.Error != nil {
		t.Fatalf("Initialize error: %d - %s", resp.Error.Code, resp.Error.Message)
	}

	if resp.Result == nil {
		t.Fatal("Expected result in initialize response")
	}

	t.Log("✓ Initialize successful")
}

// Test 3: Initialized Notification
func TestInitialized(t *testing.T) {
	resp, err := makeJSONRPCRequest("notifications/initialized", nil, 2)
	if err != nil {
		t.Fatalf("Failed to send initialized notification: %v", err)
	}

	if resp.Error != nil {
		t.Fatalf("Initialized notification error: %d - %s", resp.Error.Code, resp.Error.Message)
	}

	t.Log("✓ Initialized notification accepted")
}

// Test 4: List Tools
func TestListTools(t *testing.T) {
	resp, err := makeJSONRPCRequest("tools/list", nil, 3)
	if err != nil {
		t.Fatalf("Failed to list tools: %v", err)
	}

	if resp.Error != nil {
		t.Fatalf("List tools error: %d - %s", resp.Error.Code, resp.Error.Message)
	}

	resultMap, ok := resp.Result.(map[string]interface{})
	if !ok {
		t.Fatal("Invalid result format")
	}

	tools, ok := resultMap["tools"].([]interface{})
	if !ok {
		t.Fatal("No tools in result")
	}

	t.Logf("✓ List tools successful (found %d tools)", len(tools))
}

// Test 5: Get Projects
func TestGetProjects(t *testing.T) {
	result, err := callTool(t, "get_projects", map[string]interface{}{})
	if err != nil {
		t.Fatalf("Failed to get projects: %v", err)
	}

	items, ok := result["items"].([]interface{})
	if !ok {
		// Try direct array
		items = []interface{}{result}
	}

	t.Logf("✓ Get projects successful (found %d projects)", len(items))
}

// Test 6: E2E - Create Project
func TestE2ECreateProject(t *testing.T) {
	projectName := fmt.Sprintf("MCP Test Project - %d", time.Now().Unix())

	// Check if test project already exists
	projects, err := callTool(t, "get_projects", map[string]interface{}{})
	if err == nil {
		if items, ok := projects["items"].([]interface{}); ok {
			for _, item := range items {
				if proj, ok := item.(map[string]interface{}); ok {
					if name, ok := proj["name"].(string); ok && strings.HasPrefix(name, "MCP Test Project") {
						if id, ok := proj["id"].(string); ok {
							e2eRes.projectID = id
							t.Logf("⚠ Using existing test project: %s", id)
							return
						}
					}
				}
			}
		}
	}

	// Create new project
	result, err := callTool(t, "create_project", map[string]interface{}{
		"name": projectName,
	})
	if err != nil {
		t.Fatalf("Failed to create project: %v", err)
	}

	e2eRes.projectID = extractID(result)
	if e2eRes.projectID == "" {
		t.Fatal("Failed to extract project ID")
	}

	t.Logf("✓ Created project: %s (ID: %s)", projectName, e2eRes.projectID)
}

// Test 7: E2E - Create Board
func TestE2ECreateBoard(t *testing.T) {
	if e2eRes.projectID == "" {
		t.Skip("Skipping: no project ID")
	}

	// Check if test board already exists
	boards, err := callTool(t, "get_boards", map[string]interface{}{
		"projectId": e2eRes.projectID,
	})
	if err == nil {
		if items, ok := boards["items"].([]interface{}); ok {
			for _, item := range items {
				if board, ok := item.(map[string]interface{}); ok {
					if name, ok := board["name"].(string); ok && name == "MCP Test Board" {
						if id, ok := board["id"].(string); ok {
							e2eRes.boardID = id
							t.Logf("⚠ Using existing test board: %s", id)
							return
						}
					}
				}
			}
		}
	}

	// Create new board
	result, err := callTool(t, "create_board", map[string]interface{}{
		"projectId": e2eRes.projectID,
		"name":      "MCP Test Board",
	})
	if err != nil {
		t.Fatalf("Failed to create board: %v", err)
	}

	e2eRes.boardID = extractID(result)
	if e2eRes.boardID == "" {
		t.Fatal("Failed to extract board ID")
	}

	t.Logf("✓ Created board: MCP Test Board (ID: %s)", e2eRes.boardID)
}

// Test 8: E2E - Create List
func TestE2ECreateList(t *testing.T) {
	if e2eRes.boardID == "" {
		t.Skip("Skipping: no board ID")
	}

	// Check if test list already exists
	lists, err := callTool(t, "get_lists", map[string]interface{}{
		"boardId": e2eRes.boardID,
	})
	if err == nil {
		if items, ok := lists["items"].([]interface{}); ok {
			for _, item := range items {
				if list, ok := item.(map[string]interface{}); ok {
					if name, ok := list["name"].(string); ok && name == "MCP Test List" {
						if id, ok := list["id"].(string); ok {
							e2eRes.listID = id
							t.Logf("⚠ Using existing test list: %s", id)
							return
						}
					}
				}
			}
		}
	}

	// Create new list
	result, err := callTool(t, "create_list", map[string]interface{}{
		"boardId": e2eRes.boardID,
		"name":     "MCP Test List",
		"position": 1.0,
	})
	if err != nil {
		t.Fatalf("Failed to create list: %v", err)
	}

	e2eRes.listID = extractID(result)
	if e2eRes.listID == "" {
		t.Fatal("Failed to extract list ID")
	}

	t.Logf("✓ Created list: MCP Test List (ID: %s)", e2eRes.listID)
}

// Test 9: E2E - Create Cards
func TestE2ECreateCards(t *testing.T) {
	if e2eRes.listID == "" {
		t.Skip("Skipping: no list ID")
	}

	cardNames := []string{"Test Card 1", "Test Card 2", "Test Card 3"}
	position := 1.0

	for _, cardName := range cardNames {
		result, err := callTool(t, "create_card", map[string]interface{}{
			"listId":   e2eRes.listID,
			"name":     cardName,
			"position": position,
		})
		if err != nil {
			t.Errorf("Failed to create card %s: %v", cardName, err)
			continue
		}

		cardID := extractID(result)
		if cardID == "" {
			t.Errorf("Failed to extract card ID for %s", cardName)
			continue
		}

		e2eRes.cardIDs = append(e2eRes.cardIDs, cardID)
		t.Logf("✓ Created card: %s (ID: %s)", cardName, cardID)
		position++
	}

	if len(e2eRes.cardIDs) == 0 {
		t.Fatal("Failed to create any cards")
	}
}

// Test 10: E2E - Create Tasks
func TestE2ECreateTasks(t *testing.T) {
	if len(e2eRes.cardIDs) == 0 {
		t.Skip("Skipping: no cards available")
	}

	taskNames := []string{"Task 1", "Task 2", "Task 3"}

	for _, cardID := range e2eRes.cardIDs {
		position := 1.0
		for _, taskName := range taskNames {
			result, err := callTool(t, "create_task", map[string]interface{}{
				"cardId":  cardID,
				"name":    taskName,
				"position": position,
			})
			if err != nil {
				t.Errorf("Failed to create task %s for card %s: %v", taskName, cardID, err)
				continue
			}

			taskID := extractID(result)
			if taskID == "" {
				t.Errorf("Failed to extract task ID for %s", taskName)
				continue
			}

			e2eRes.taskIDs = append(e2eRes.taskIDs, taskID)
			t.Logf("✓ Created task: %s for card %s (Task ID: %s)", taskName, cardID, taskID)
			position++
		}
	}

	if len(e2eRes.taskIDs) == 0 {
		t.Fatal("Failed to create any tasks")
	}
}

// Test 11: E2E - Verify Resources
func TestE2EVerifyResources(t *testing.T) {
	// Verify project
	if e2eRes.projectID != "" {
		result, err := callTool(t, "get_project", map[string]interface{}{
			"projectId": e2eRes.projectID,
		})
		if err == nil {
			if name, ok := result["name"].(string); ok {
				t.Logf("✓ Verified project: %s", name)
			}
		}
	}

	// Verify board
	if e2eRes.boardID != "" {
		result, err := callTool(t, "get_board", map[string]interface{}{
			"boardId": e2eRes.boardID,
		})
		if err == nil {
			if name, ok := result["name"].(string); ok {
				t.Logf("✓ Verified board: %s", name)
			}
		}
	}

	// Verify list
	if e2eRes.listID != "" {
		result, err := callTool(t, "get_list", map[string]interface{}{
			"listId": e2eRes.listID,
		})
		if err == nil {
			if name, ok := result["name"].(string); ok {
				t.Logf("✓ Verified list: %s", name)
			}
		}
	}

	// Verify cards
	if e2eRes.listID != "" {
		result, err := callTool(t, "get_cards", map[string]interface{}{
			"listId": e2eRes.listID,
		})
		if err == nil {
			var cardCount int
			if items, ok := result["items"].([]interface{}); ok {
				cardCount = len(items)
			}
			t.Logf("✓ Verified cards: %d cards found", cardCount)
		}
	}

	// Verify tasks
	if len(e2eRes.cardIDs) > 0 && len(e2eRes.taskIDs) > 0 {
		result, err := callTool(t, "get_tasks", map[string]interface{}{
			"cardId": e2eRes.cardIDs[0],
		})
		if err == nil {
			var taskCount int
			if items, ok := result["items"].([]interface{}); ok {
				taskCount = len(items)
			}
			t.Logf("✓ Verified tasks: %d tasks found for first card", taskCount)
		}
	}
}

// Test 12: E2E - Cleanup
func TestE2ECleanup(t *testing.T) {
	// Delete tasks
	for _, taskID := range e2eRes.taskIDs {
		_, err := callTool(t, "delete_task", map[string]interface{}{
			"taskId": taskID,
		})
		if err == nil {
			t.Logf("  ✓ Deleted task: %s", taskID)
		} else {
			t.Logf("  ⚠ Failed to delete task: %s", taskID)
		}
	}

	// Delete cards
	for _, cardID := range e2eRes.cardIDs {
		_, err := callTool(t, "delete_card", map[string]interface{}{
			"cardId": cardID,
		})
		if err == nil {
			t.Logf("  ✓ Deleted card: %s", cardID)
		} else {
			t.Logf("  ⚠ Failed to delete card: %s", cardID)
		}
	}

	// Delete list
	if e2eRes.listID != "" {
		_, err := callTool(t, "delete_list", map[string]interface{}{
			"listId": e2eRes.listID,
		})
		if err == nil {
			t.Logf("  ✓ Deleted list: %s", e2eRes.listID)
		} else {
			t.Logf("  ⚠ Failed to delete list: %s", e2eRes.listID)
		}
	}

	// Delete board
	if e2eRes.boardID != "" {
		_, err := callTool(t, "delete_board", map[string]interface{}{
			"boardId": e2eRes.boardID,
		})
		if err == nil {
			t.Logf("  ✓ Deleted board: %s", e2eRes.boardID)
		} else {
			t.Logf("  ⚠ Failed to delete board: %s", e2eRes.boardID)
		}
	}

	// Delete project
	if e2eRes.projectID != "" {
		_, err := callTool(t, "delete_project", map[string]interface{}{
			"projectId": e2eRes.projectID,
		})
		if err == nil {
			t.Logf("  ✓ Deleted project: %s", e2eRes.projectID)
		} else {
			t.Logf("  ⚠ Failed to delete project: %s", e2eRes.projectID)
		}
	}

	t.Log("✓ Cleanup completed")
}

// Test 13: Invalid Method
func TestInvalidMethod(t *testing.T) {
	resp, err := makeJSONRPCRequest("invalid_method", nil, 9)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}

	if resp.Error == nil {
		t.Fatal("Expected error for invalid method, but got success")
	}

	// Accept any JSON-RPC error code
	if resp.Error.Code < -32768 || resp.Error.Code > -32600 {
		t.Fatalf("Unexpected error code: %d", resp.Error.Code)
	}

	t.Logf("✓ Invalid method correctly rejected (error code: %d)", resp.Error.Code)
}

// Test 14: Invalid JSON
func TestInvalidJSON(t *testing.T) {
	resp, err := http.Post(mcpEndpoint, "application/json", bytes.NewBufferString("invalid json"))
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	var jsonResp jsonRPCResponse
	if err := json.NewDecoder(resp.Body).Decode(&jsonResp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if jsonResp.Error == nil {
		t.Fatal("Expected error for invalid JSON, but got success")
	}

	t.Log("✓ Invalid JSON correctly rejected")
}

// Test 15: CORS Preflight
func TestCORS(t *testing.T) {
	req, err := http.NewRequest("OPTIONS", mcpEndpoint, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Access-Control-Request-Method", "POST")
	req.Header.Set("Access-Control-Request-Headers", "Content-Type")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	corsHeader := resp.Header.Get("Access-Control-Allow-Origin")
	if corsHeader == "" {
		t.Fatal("CORS headers missing")
	}

	t.Log("✓ CORS headers present")
}

