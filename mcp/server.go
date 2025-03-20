package mcp

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"strings"
)

// Tool represents a registered tool function
type Tool struct {
	Name        string
	Handler     interface{}
	Description string
}

// Resource represents a registered resource endpoint
type Resource struct {
	Pattern     string
	Handler     interface{}
	Description string
	Params      []string
}

// MCPServer represents the main MCP server
type MCPServer struct {
	Name      string
	Tools     map[string]Tool
	Resources map[string]Resource
}

// NewMCPServer creates a new MCP server instance
func NewMCPServer(name string) *MCPServer {
	return &MCPServer{
		Name:      name,
		Tools:     make(map[string]Tool),
		Resources: make(map[string]Resource),
	}
}

// Tool registers a new tool function
func (s *MCPServer) Tool() func(interface{}) {
	return func(handler interface{}) {
		handlerType := reflect.TypeOf(handler)
		if handlerType.Kind() != reflect.Func {
			log.Fatal("Tool handler must be a function")
		}

		// Get function name
		handlerValue := reflect.ValueOf(handler)
		name := runtime.FuncForPC(handlerValue.Pointer()).Name()
		name = strings.TrimSuffix(strings.TrimPrefix(name, "main."), "-fm")

		// Get function description from comments if available
		description := fmt.Sprintf("Tool: %s", name)

		s.Tools[name] = Tool{
			Name:        name,
			Handler:     handler,
			Description: description,
		}
	}
}

// Resource registers a new resource endpoint
func (s *MCPServer) Resource(pattern string) func(interface{}) {
	return func(handler interface{}) {
		handlerType := reflect.TypeOf(handler)
		if handlerType.Kind() != reflect.Func {
			log.Fatal("Resource handler must be a function")
		}

		// Extract parameters from pattern
		params := extractParams(pattern)
		
		// Convert pattern to regex for matching
		regexPattern := patternToRegex(pattern)

		// Get function name and description
		handlerValue := reflect.ValueOf(handler)
		name := runtime.FuncForPC(handlerValue.Pointer()).Name()
		name = strings.TrimSuffix(strings.TrimPrefix(name, "main."), "-fm")
		description := fmt.Sprintf("Resource: %s", pattern)

		s.Resources[pattern] = Resource{
			Pattern:     regexPattern,
			Handler:     handler,
			Description: description,
			Params:      params,
		}
	}
}

// extractParams extracts parameter names from a pattern
func extractParams(pattern string) []string {
	var params []string
	re := regexp.MustCompile(`\{([^}]+)\}`)
	matches := re.FindAllStringSubmatch(pattern, -1)
	for _, match := range matches {
		params = append(params, match[1])
	}
	return params
}

// patternToRegex converts a pattern to a regex pattern
func patternToRegex(pattern string) string {
	re := regexp.MustCompile(`\{([^}]+)\}`)
	return re.ReplaceAllString(pattern, `([^/]+)`)
}

// ServeHTTP implements the http.Handler interface
func (s *MCPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	// Check resources
	for pattern, resource := range s.Resources {
		re := regexp.MustCompile("^" + resource.Pattern + "$")
		matches := re.FindStringSubmatch(path)
		
		if matches != nil {
			// Extract parameters
			params := matches[1:]
			
			// Create args for the handler
			args := make([]reflect.Value, len(params))
			handlerType := reflect.TypeOf(resource.Handler)
			
			for i, param := range params {
				paramType := handlerType.In(i)
				args[i] = reflect.ValueOf(param).Convert(paramType)
			}
			
			// Call handler
			result := reflect.ValueOf(resource.Handler).Call(args)
			
			// Send response
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"result": result[0].Interface(),
			})
			return
		}
	}

	// Handle tool calls
	if r.Method == http.MethodPost && strings.HasPrefix(path, "/tool/") {
		toolName := strings.TrimPrefix(path, "/tool/")
		if tool, exists := s.Tools[toolName]; exists {
			var params map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
				http.Error(w, "Invalid request body", http.StatusBadRequest)
				return
			}

			// Convert parameters and call tool
			handlerType := reflect.TypeOf(tool.Handler)
			args := make([]reflect.Value, handlerType.NumIn())
			
			for i := 0; i < handlerType.NumIn(); i++ {
				paramName := handlerType.In(i).Name()
				if val, ok := params[paramName]; ok {
					args[i] = reflect.ValueOf(val).Convert(handlerType.In(i))
				}
			}

			result := reflect.ValueOf(tool.Handler).Call(args)
			
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"result": result[0].Interface(),
			})
			return
		}
	}

	http.NotFound(w, r)
}

// Start starts the MCP server
func (s *MCPServer) Start(port string) error {
	addr := fmt.Sprintf(":%s", port)
	log.Printf("Starting %s MCP server on %s", s.Name, addr)
	return http.ListenAndServe(addr, s)
}