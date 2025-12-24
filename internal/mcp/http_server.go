package mcp

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

// sessionState tracks initialization state per session
type sessionState struct {
	initialized bool
	mu          sync.RWMutex
}

// HTTP server with session management
type httpServer struct {
	server  *Server
	sessions map[string]*sessionState
	mu      sync.RWMutex
}

// StartHTTP starts the MCP server in HTTP mode
func (s *Server) StartHTTP(addr string, port int) error {
	httpSrv := &httpServer{
		server:   s,
		sessions: make(map[string]*sessionState),
	}

	mux := http.NewServeMux()
	
	// Main MCP JSON-RPC endpoint
	mux.HandleFunc("/mcp", httpSrv.handleMCPRequest)
	mux.HandleFunc("/", httpSrv.handleMCPRequest) // Also support root path
	
	// Health check endpoint
	mux.HandleFunc("/health", httpSrv.handleHealth)
	
	serverAddr := fmt.Sprintf("%s:%d", addr, port)
	log.Printf("HTTP server listening on %s", serverAddr)
	log.Printf("MCP endpoint: http://%s/mcp", serverAddr)
	
	return http.ListenAndServe(serverAddr, httpSrv.corsMiddleware(mux))
}

// corsMiddleware adds CORS headers to responses
func (h *httpServer) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Content-Type", "application/json")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

// handleHealth handles health check requests
func (h *httpServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	response := map[string]interface{}{
		"status": "ok",
		"service": "planka-mcp",
	}
	
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// handleMCPRequest handles MCP JSON-RPC requests over HTTP
func (h *httpServer) handleMCPRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed. Use POST for JSON-RPC requests.", http.StatusMethodNotAllowed)
		return
	}

	// Get or create session based on a simple identifier
	// For stateless HTTP, we could use a session token, but for simplicity,
	// we'll track by a combination of factors or make initialization per-request
	sessionID := h.getSessionID(r)
	session := h.getOrCreateSession(sessionID)

	// Decode JSON-RPC request
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.sendHTTPError(w, nil, fmt.Errorf("failed to decode request: %w", err), http.StatusBadRequest)
		return
	}

	method, _ := request["method"].(string)
	id, _ := request["id"]

	// Handle initialization
	if method == "initialize" {
		session.mu.Lock()
		session.initialized = true
		session.mu.Unlock()
		
		response := h.server.buildInitializeResponse(id)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Handle initialized notification
	if method == "notifications/initialized" {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"jsonrpc": "2.0",
			"result":  nil,
			"id":      id,
		})
		return
	}

	// Check if initialized (for HTTP, we're more lenient - allow requests without explicit init)
	session.mu.RLock()
	initialized := session.initialized
	session.mu.RUnlock()

	// For HTTP mode, we can auto-initialize if not done
	if !initialized && method != "initialize" {
		// Auto-initialize for HTTP mode
		session.mu.Lock()
		session.initialized = true
		session.mu.Unlock()
	}

	// Handle the request
	response, err := h.server.handleMCPRequest(request)
	if err != nil {
		h.sendHTTPError(w, request, err, http.StatusOK) // JSON-RPC errors still return 200
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode response: %v", err)
	}
}

// getSessionID generates a session ID from the request
// For simplicity, we can use IP + User-Agent, or a session token if provided
func (h *httpServer) getSessionID(r *http.Request) string {
	// Check for session token in header
	if token := r.Header.Get("X-Session-Token"); token != "" {
		return token
	}
	
	// Fallback to IP + User-Agent combination
	return r.RemoteAddr + r.Header.Get("User-Agent")
}

// getOrCreateSession gets or creates a session state
func (h *httpServer) getOrCreateSession(sessionID string) *sessionState {
	h.mu.RLock()
	session, exists := h.sessions[sessionID]
	h.mu.RUnlock()
	
	if exists {
		return session
	}
	
	h.mu.Lock()
	defer h.mu.Unlock()
	
	// Double-check after acquiring write lock
	if session, exists := h.sessions[sessionID]; exists {
		return session
	}
	
	session = &sessionState{initialized: false}
	h.sessions[sessionID] = session
	return session
}

// sendHTTPError sends a JSON-RPC error response
func (h *httpServer) sendHTTPError(w http.ResponseWriter, request map[string]interface{}, err error, statusCode int) {
	id := interface{}(nil)
	if request != nil {
		id, _ = request["id"]
	}
	
	response := h.server.buildErrorResponse(id, err)
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

