package main

import (
	"mcp-tool-kit/internal/tools"
	"os"
	"strings"

	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Create a new MCP server instance
	server := server.NewMCPServer(
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
			tools.NewJiraTool(server)
		}
		
		if tool == "sql-server" {
			tools.NewSQLServerTool(server)
		}
	}
} 