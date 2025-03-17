# MCP Tool Kit

A Go-based MCP (Mission Control Protocol) server implementation using the mark3labs/mcp-go library.

## Prerequisites

- Go 1.16 or higher
- Git

## Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/mcp-tool-kit.git
cd mcp-tool-kit
```

2. Install dependencies:
```bash
go mod download
```

## Running the Server

To start the MCP server:

```bash
go run main.go
```

The server will start and listen for connections. To stop the server, press Ctrl+C.

## Project Structure

```
mcp-tool-kit/
├── main.go          # Main server implementation
├── go.mod           # Go module definition
└── README.md        # Project documentation
```

## Features

- Basic MCP server implementation
- Graceful shutdown handling
- Signal handling for clean termination

## License

MIT License 