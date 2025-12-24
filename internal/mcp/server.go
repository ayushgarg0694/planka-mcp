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

// StartStdio starts the MCP server in stdio mode
func (s *Server) StartStdio() error {
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

// buildInitializeResponse builds the response for initialize
func (s *Server) buildInitializeResponse(id interface{}) map[string]interface{} {
	return map[string]interface{}{
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
}

// handleInitialize handles the initialize request (stdio mode)
func (s *Server) handleInitialize(request map[string]interface{}, encoder *json.Encoder, id interface{}) error {
	response := s.buildInitializeResponse(id)
	return encoder.Encode(response)
}

// handleMCPRequest handles an MCP request and returns the response map
// This is the shared request handler used by both stdio and HTTP modes
func (s *Server) handleMCPRequest(request map[string]interface{}) (map[string]interface{}, error) {
	method, ok := request["method"].(string)
	if !ok {
		return nil, fmt.Errorf("missing method in request")
	}

	id, _ := request["id"]

	switch method {
	case "tools/list":
		return s.buildToolsListResponse(id), nil
	case "tools/call":
		return s.buildToolsCallResponse(request, id)
	default:
		return nil, fmt.Errorf("unknown method: %s", method)
	}
}

// handleRequest handles an MCP request (stdio mode)
func (s *Server) handleRequest(request map[string]interface{}, encoder *json.Encoder) error {
	response, err := s.handleMCPRequest(request)
	if err != nil {
		return err
	}
	return encoder.Encode(response)
}

// buildToolsListResponse builds the response for tools/list
func (s *Server) buildToolsListResponse(id interface{}) map[string]interface{} {
	tools := s.getTools()
	return map[string]interface{}{
		"jsonrpc": "2.0",
		"result": map[string]interface{}{
			"tools": tools,
		},
		"id": id,
	}
}

// handleToolsList handles the tools/list request (stdio mode)
func (s *Server) handleToolsList(encoder *json.Encoder, id interface{}) error {
	response := s.buildToolsListResponse(id)
	return encoder.Encode(response)
}

// buildToolsCallResponse builds the response for tools/call
func (s *Server) buildToolsCallResponse(request map[string]interface{}, id interface{}) (map[string]interface{}, error) {
	params, ok := request["params"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("missing params in request")
	}

	toolName, ok := params["name"].(string)
	if !ok {
		return nil, fmt.Errorf("missing name in params")
	}

	arguments, _ := params["arguments"].(map[string]interface{})

	result, err := s.callTool(toolName, arguments)
	if err != nil {
		return nil, fmt.Errorf("tool call failed: %w", err)
	}

	return map[string]interface{}{
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
	}, nil
}

// handleToolsCall handles the tools/call request (stdio mode)
func (s *Server) handleToolsCall(request map[string]interface{}, encoder *json.Encoder, id interface{}) error {
	response, err := s.buildToolsCallResponse(request, id)
	if err != nil {
		return err
	}
	return encoder.Encode(response)
}

// buildErrorResponse builds an error response
func (s *Server) buildErrorResponse(id interface{}, err error) map[string]interface{} {
	return map[string]interface{}{
		"jsonrpc": "2.0",
		"error": map[string]interface{}{
			"code":    -32603,
			"message": err.Error(),
		},
		"id": id,
	}
}

// sendError sends an error response (stdio mode)
func (s *Server) sendError(encoder *json.Encoder, request map[string]interface{}, err error) {
	id, _ := request["id"]
	response := s.buildErrorResponse(id, err)
	encoder.Encode(response)
}


