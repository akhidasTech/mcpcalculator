package main

import (
	"log"
	"github.com/akhidasTech/mcpcalculator/mcp/protocol"
)

// AddParams represents parameters for the Add tool
type AddParams struct {
	A int `json:"a"`
	B int `json:"b"`
}

// Add adds two numbers
func Add(params interface{}) (interface{}, error) {
	p, ok := params.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid parameters")
	}

	a, ok := p["a"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid parameter 'a'")
	}

	b, ok := p["b"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid parameter 'b'")
	}

	return int(a) + int(b), nil
}

// GreetingParams represents parameters for the Greeting resource
type GreetingParams struct {
	Name string `json:"name"`
}

// GetGreeting returns a personalized greeting
func GetGreeting(params interface{}) (interface{}, error) {
	p, ok := params.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid parameters")
	}

	name, ok := p["name"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid parameter 'name'")
	}

	return fmt.Sprintf("Hello, %s!", name), nil
}

func main() {
	// Create a new MCP server
	server := protocol.NewServer("Demo")

	// Register the Add tool
	server.RegisterTool(
		"add",
		Add,
		"Add two numbers",
		AddParams{},
	)

	// Register the greeting resource
	server.RegisterResource(
		"greeting/{name}",
		GetGreeting,
		"Get a personalized greeting",
	)

	// Start the server
	log.Fatal(server.Start("8080"))
}
