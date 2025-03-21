package main

import (
	"log"
	"mcp-tool-kit/internal/tools"
	"os"
	"strings"

	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Create a new MCP server instance
	mcpServer := server.NewMCPServer(
		"mcp-tool-kit",
		"1.0.0",
		server.WithLogging(),
		server.WithPromptCapabilities(true),
		server.WithResourceCapabilities(true, true),
	)

	// Load config tools from environment variable
	configTools := strings.Split(os.Getenv("CONFIG_TOOLS"), ",")

	for _, tool := range configTools {
		if tool == "jira" {
			tools.NewJiraTool(mcpServer)
		}
		
		if tool == "sql-server" {
			tools.NewSQLServerTool(mcpServer)
		}
	}

	sseServer := server.NewSSEServer(
		mcpServer, 
		server.WithBaseURL("http://localhost:8080"),
		server.WithMessageEndpoint("/message"),
		server.WithSSEEndpoint("/sse"),
	)
	
	err := sseServer.Start(":8080")
	if err != nil {
		log.Fatalf("Failed to start SSE server: %v", err)
	}

	log.Println("SSE server started on http://localhost:8080")
} 
