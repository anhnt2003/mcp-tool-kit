package main

import (
	"log"
	"mcp-tool-kit/internal/tools"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/mark3labs/mcp-go/server"
)

// initConfig loads and validates environment configurations
func initConfig() error {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found. Using environment variables.")
	}
	return nil
}

// registerTools initializes and registers tools based on configuration
func registerTools(mcpServer *server.MCPServer) {
	configTools := strings.Split(os.Getenv("CONFIG_TOOLS"), ",")
	log.Printf("Initializing configured tools: %v", configTools)

	for _, tool := range configTools {
		tool = strings.TrimSpace(tool)
		if tool == "" {
			continue
		}

		switch tool {
		case "jira":
			tools.NewJiraTool(mcpServer)
		case "sql-server":
			tools.NewSQLServerTool(mcpServer)
		default:
			log.Printf("Unknown tool configuration: %s", tool)
		}
	}
}

func main() {
	// Initialize logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Load configuration
	if err := initConfig(); err != nil {
		log.Fatalf("Failed to initialize configuration: %v", err)
	}

	// Create a new MCP server instance
	mcpServer := server.NewMCPServer(
		"mcp-tool-kit",
		"1.0.0",
		server.WithLogging(),
		server.WithPromptCapabilities(true),
		server.WithResourceCapabilities(true, true),
	)

	// Register tools
	registerTools(mcpServer)

	// Determine the execution mode
	mode := strings.ToLower(os.Getenv("MCP_MODE"))
	if mode == "" {
		mode = "stdio" // Default to stdio if not specified
	}

	switch mode {
	case "stdio":
		// Run in stdio mode
		log.Println("Starting in stdio mode")
		err := server.ServeStdio(mcpServer)
		if err != nil {
			log.Fatalf("Failed to start stdio server: %v", err)
		}
	case "sse":
		// Create a new SSE server instance
		sseServer := server.NewSSEServer(
			mcpServer, 
			server.WithBaseURL("http://localhost:8080"),
			server.WithMessageEndpoint("/message"),
			server.WithSSEEndpoint("/sse"),
		)

		log.Printf("Starting SSE server on http://localhost:8080")
		
		// Start the SSE server
		err := sseServer.Start(":8080")
		if err != nil {
			log.Fatalf("Failed to start SSE server: %v", err)
		}
	default:
		log.Fatalf("Unknown MCP_MODE: %s. Supported modes are 'stdio' and 'sse'", mode)
	}
} 
