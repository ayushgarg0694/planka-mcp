package mcp

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/ayushgarg/mcp-planka/internal/planka"
)

// Server represents an MCP server
type Server struct {
	client *planka.Client
}

// NewServer creates a new MCP server
func NewServer(client *planka.Client) *Server {
	return &Server{
		client: client,
	}
}

// Start starts the MCP server
func (s *Server) Start() error {
	// MCP servers communicate via stdio
	decoder := json.NewDecoder(os.Stdin)
	encoder := json.NewEncoder(os.Stdout)

	// Wait for and handle initialization request
	initialized := false
	for {
		var request map[string]interface{}
		if err := decoder.Decode(&request); err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to decode request: %w", err)
		}

		method, _ := request["method"].(string)
		id, _ := request["id"]

		// Handle initialization
		if method == "initialize" {
			if err := s.handleInitialize(request, encoder, id); err != nil {
				return fmt.Errorf("failed to handle initialize: %w", err)
			}
			initialized = true
			continue
		}

		// Handle initialized notification
		if method == "notifications/initialized" {
			// Client is now ready, continue to normal request handling
			continue
		}

		// Only handle other requests after initialization
		if !initialized {
			return fmt.Errorf("received request before initialization")
		}

		if err := s.handleRequest(request, encoder); err != nil {
			log.Printf("Error handling request: %v", err)
			s.sendError(encoder, request, err)
		}
	}

	return nil
}

// handleInitialize handles the initialize request
func (s *Server) handleInitialize(request map[string]interface{}, encoder *json.Encoder, id interface{}) error {
	response := map[string]interface{}{
		"jsonrpc": "2.0",
		"result": map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities": map[string]interface{}{
				"tools": map[string]interface{}{},
			},
			"serverInfo": map[string]interface{}{
				"name":    "planka-mcp",
				"version": "1.0.0",
			},
		},
		"id": id,
	}
	return encoder.Encode(response)
}

// handleRequest handles an MCP request
func (s *Server) handleRequest(request map[string]interface{}, encoder *json.Encoder) error {
	method, ok := request["method"].(string)
	if !ok {
		return fmt.Errorf("missing method in request")
	}

	id, _ := request["id"]

	switch method {
	case "tools/list":
		return s.handleToolsList(encoder, id)
	case "tools/call":
		return s.handleToolsCall(request, encoder, id)
	default:
		return fmt.Errorf("unknown method: %s", method)
	}
}

// handleToolsList handles the tools/list request
func (s *Server) handleToolsList(encoder *json.Encoder, id interface{}) error {
	tools := s.getTools()
	response := map[string]interface{}{
		"jsonrpc": "2.0",
		"result": map[string]interface{}{
			"tools": tools,
		},
		"id": id,
	}
	return encoder.Encode(response)
}

// handleToolsCall handles the tools/call request
func (s *Server) handleToolsCall(request map[string]interface{}, encoder *json.Encoder, id interface{}) error {
	params, ok := request["params"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("missing params in request")
	}

	toolName, ok := params["name"].(string)
	if !ok {
		return fmt.Errorf("missing name in params")
	}

	arguments, _ := params["arguments"].(map[string]interface{})

	result, err := s.callTool(toolName, arguments)
	if err != nil {
		return fmt.Errorf("tool call failed: %w", err)
	}

	response := map[string]interface{}{
		"jsonrpc": "2.0",
		"result": map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": result,
				},
			},
		},
		"id": id,
	}
	return encoder.Encode(response)
}

// sendError sends an error response
func (s *Server) sendError(encoder *json.Encoder, request map[string]interface{}, err error) {
	id, _ := request["id"]
	response := map[string]interface{}{
		"jsonrpc": "2.0",
		"error": map[string]interface{}{
			"code":    -32603,
			"message": err.Error(),
		},
		"id": id,
	}
	encoder.Encode(response)
}

