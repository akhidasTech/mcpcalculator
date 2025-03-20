package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
	"encoding/json"
)

// Response represents the JSON response structure
type Response struct {
	Result interface{} `json:"result"`
	Error  string      `json:"error,omitempty"`
}

// MCPServer represents our server structure
type MCPServer struct {
	name   string
	router *mux.Router
}

// NewMCPServer creates a new MCP server instance
func NewMCPServer(name string) *MCPServer {
	return &MCPServer{
		name:   name,
		router: mux.NewRouter(),
	}
}

// Add handles the addition of two numbers
func (s *MCPServer) Add(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	a, err1 := strconv.Atoi(queryParams.Get("a"))
	b, err2 := strconv.Atoi(queryParams.Get("b"))

	if err1 != nil || err2 != nil {
		json.NewEncoder(w).Encode(Response{Error: "Invalid parameters"})
		return
	}

	result := a + b
	json.NewEncoder(w).Encode(Response{Result: result})
}

// GetGreeting handles the greeting endpoint
func (s *MCPServer) GetGreeting(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	greeting := fmt.Sprintf("Hello, %s!", name)
	json.NewEncoder(w).Encode(Response{Result: greeting})
}

// SetupRoutes configures all the routes for our server
func (s *MCPServer) SetupRoutes() {
	s.router.HandleFunc("/add", s.Add).Methods("GET")
	s.router.HandleFunc("/greeting/{name}", s.GetGreeting).Methods("GET")
}

// Start starts the server on the specified port
func (s *MCPServer) Start(port string) error {
	s.SetupRoutes()
	log.Printf("Starting %s server on port %s", s.name, port)
	return http.ListenAndServe(":" + port, s.router)
}

func main() {
	// Create a new MCP server
	mcp := NewMCPServer("Demo")

	// Start the server
	log.Fatal(mcp.Start("8080"))
}