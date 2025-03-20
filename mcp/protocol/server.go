package protocol

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

// Server represents an MCP server
type Server struct {
	Info      ServerInfo
	handlers  map[string]HandlerFunc
	resources map[string]ResourceHandler
	tools     map[string]ToolHandler
	mutex     sync.RWMutex
}

// HandlerFunc represents a request handler function
type HandlerFunc func(params interface{}) (interface{}, error)

// ResourceHandler represents a resource handler
type ResourceHandler struct {
	Pattern     string
	Handler     interface{}
	Description string
}

// ToolHandler represents a tool handler
type ToolHandler struct {
	Name        string
	Handler     interface{}
	Description string
	Schema      interface{}
}

// NewServer creates a new MCP server
func NewServer(name string) *Server {
	return &Server{
		Info: ServerInfo{
			Name:    name,
			Version: "1.0.0",
			Capabilities: []Capability{
				{
					Name:        "resources",
					Version:     "1.0.0",
					Description: "Resource access capability",
				},
				{
					Name:        "tools",
					Version:     "1.0.0",
					Description: "Tool execution capability",
				},
			},
		},
		handlers:  make(map[string]HandlerFunc),
		resources: make(map[string]ResourceHandler),
		tools:     make(map[string]ToolHandler),
	}
}

// RegisterHandler registers a request handler
func (s *Server) RegisterHandler(method string, handler HandlerFunc) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.handlers[method] = handler
}

// RegisterResource registers a resource handler
func (s *Server) RegisterResource(pattern string, handler interface{}, description string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.resources[pattern] = ResourceHandler{
		Pattern:     pattern,
		Handler:     handler,
		Description: description,
	}
}

// RegisterTool registers a tool handler
func (s *Server) RegisterTool(name string, handler interface{}, description string, schema interface{}) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.tools[name] = ToolHandler{
		Name:        name,
		Handler:     handler,
		Description: description,
		Schema:      schema,
	}
}

// ServeHTTP implements the http.Handler interface
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.writeError(w, &req, -32700, "Parse error", nil)
		return
	}

	if req.JSONRPC != "2.0" {
		s.writeError(w, &req, -32600, "Invalid Request", nil)
		return
	}

	s.handleRequest(w, &req)
}

// handleRequest processes an incoming request
func (s *Server) handleRequest(w http.ResponseWriter, req *Request) {
	s.mutex.RLock()
	handler, exists := s.handlers[req.Method]
	s.mutex.RUnlock()

	if !exists {
		s.writeError(w, req, -32601, "Method not found", nil)
		return
	}

	result, err := handler(req.Params)
	if err != nil {
		s.writeError(w, req, -32603, err.Error(), nil)
		return
	}

	s.writeResponse(w, req, result)
}

// writeResponse writes a successful response
func (s *Server) writeResponse(w http.ResponseWriter, req *Request, result interface{}) {
	resp := Response{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  result,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// writeError writes an error response
func (s *Server) writeError(w http.ResponseWriter, req *Request, code int, message string, data interface{}) {
	resp := Response{
		JSONRPC: "2.0",
		ID:      req.ID,
		Error: &Error{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding error response: %v", err)
	}
}

// Start starts the MCP server
func (s *Server) Start(port string) error {
	addr := fmt.Sprintf(":%s", port)
	log.Printf("Starting %s MCP server on %s", s.Info.Name, addr)
	return http.ListenAndServe(addr, s)
}
