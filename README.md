# MCP Calculator

A Go implementation of a Model Context Protocol (MCP) server with calculator and greeting functionality. This implementation follows the official MCP specification.

## Features

- Full JSON-RPC 2.0 implementation
- MCP protocol support
- Tool registration and execution
- Resource handling
- Server capability negotiation
- Error handling according to spec

## Project Structure

```
mcpcalculator/
├── mcp/
│   └── protocol/
│       ├── types.go    # Protocol types and structures
│       ├── server.go   # MCP server implementation
│       └── client.go   # MCP client implementation
├── main.go            # Example usage
├── go.mod             # Go module file
└── README.md          # Documentation
```

## Requirements

- Go 1.21 or higher

## Installation

```bash
git clone https://github.com/akhidasTech/mcpcalculator.git
cd mcpcalculator
go mod download
```

## Usage

### Starting the Server

```bash
go run main.go
```

The server will start on port 8080.

### Making Requests

Tools and resources follow the JSON-RPC 2.0 specification:

#### Add Tool

```bash
curl -X POST http://localhost:8080 \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "add",
    "params": {"a": 5, "b": 3}
  }'
```

#### Greeting Resource

```bash
curl -X POST http://localhost:8080 \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "greeting",
    "params": {"name": "John"}
  }'
```

### Response Format

All responses follow the JSON-RPC 2.0 format:

```json
{
    "jsonrpc": "2.0",
    "id": 1,
    "result": <result_value>
}
```

Or for errors:

```json
{
    "jsonrpc": "2.0",
    "id": 1,
    "error": {
        "code": <error_code>,
        "message": "error message",
        "data": <additional_data>
    }
}
```

## Protocol Implementation

This implementation follows the Model Context Protocol specification:

1. **JSON-RPC 2.0**: All communication uses the JSON-RPC 2.0 protocol
2. **Capabilities**: Server advertises its capabilities during initialization
3. **Tools**: Implements the tool registration and execution protocol
4. **Resources**: Implements the resource access protocol
5. **Error Handling**: Uses standard error codes and formats

## Security

This implementation follows MCP security guidelines:

1. **User Consent**: Tools and resources require explicit invocation
2. **Data Privacy**: No data is shared without explicit requests
3. **Tool Safety**: Tool execution is controlled and validated
4. **Error Handling**: Proper error reporting and handling

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
