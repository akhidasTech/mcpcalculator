# MCP Calculator

A Go implementation of an MCP server with calculator and greeting functionality.

## Features

- Addition endpoint (`/add`)
- Dynamic greeting endpoint (`/greeting/{name}`)

## Requirements

- Go 1.21 or higher
- gorilla/mux package

## Installation

```bash
git clone https://github.com/akhidasTech/mcpcalculator.git
cd mcpcalculator
go mod download
```

## Running the Server

```bash
go run main.go
```

The server will start on port 8080.

## API Endpoints

### Add Numbers

```
GET /add?a={number}&b={number}
```

Example: `http://localhost:8080/add?a=5&b=3`

### Get Greeting

```
GET /greeting/{name}
```

Example: `http://localhost:8080/greeting/John`

## Response Format

All endpoints return JSON responses in the following format:

```json
{
    "result": <result_value>,
    "error": "error message" // Only present if there's an error
}
```