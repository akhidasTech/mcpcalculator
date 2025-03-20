package main

import (
	"log"
	"github.com/akhidasTech/mcpcalculator/mcp"
)

// Add adds two numbers
func Add(a int, b int) int {
	return a + b
}

// GetGreeting returns a personalized greeting
func GetGreeting(name string) string {
	return "Hello, " + name + "!"
}

func main() {
	// Create a new MCP server
	server := mcp.NewMCPServer("Demo")

	// Register the Add tool
	server.Tool()(Add)

	// Register the greeting resource
	server.Resource("greeting/{name}")(GetGreeting)

	// Start the server
	log.Fatal(server.Start("8080"))
}