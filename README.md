# MCP Calculator

A Go implementation of an MCP server with calculator and greeting functionality. This implementation provides a framework similar to Python's FastMCP, allowing you to easily create tools and resources.

## Features

- Custom MCP server implementation
- Tool registration system
- Dynamic resource routing
- Reflection-based parameter handling
- JSON response formatting

## Requirements

- Go 1.21 or higher

## Installation

```bash
git clone https://github.com/akhidasTech/mcpcalculator.git
cd mcpcalculator
go mod download
```

## Project Structure

```
mcpcalculator/
├── mcp/
│   └── server.go    # MCP server implementation
├── main.go          # Example usage
├── go.mod          # Go module file
└── README.md       # Documentation
```

## Usage

### Creating an MCP Server

```go
import "github.com/akhidasTech/mcpcalculator/mcp"

// Create a new MCP server
server := mcp.NewMCPServer("Demo")
```

### Registering Tools

Tools are functions that can be called via HTTP POST requests:

```go
// Define a tool function
func Add(a int, b int) int {
    return a + b
}

// Register the tool
server.Tool()(Add)
```

To call a tool:
```bash
curl -X POST http://localhost:8080/tool/Add -d '{"a": 5, "b": 3}'
```

### Registering Resources

Resources are dynamic endpoints that can handle parameters:

```go
// Define a resource function
func GetGreeting(name string) string {
    return "Hello, " + name + "!"
}

// Register the resource
server.Resource("greeting/{name}")(GetGreeting)
```

To call a resource:
```bash
curl http://localhost:8080/greeting/John
```

### Starting the Server

```go
server.Start("8080")
```

## Response Format

All endpoints return JSON responses in the following format:

```json
{
    "result": <result_value>,
    "error": "error message" // Only present if there's an error
}
```

## Example

The included example demonstrates:
- Adding two numbers using a tool
- Getting a personalized greeting using a resource

To run the example:

```bash
go run main.go
```

Then try:
- Tool: `curl -X POST http://localhost:8080/tool/Add -d '{"a": 5, "b": 3}'`
- Resource: `curl http://localhost:8080/greeting/John`