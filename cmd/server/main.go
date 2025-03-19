package main

import (
	"context"
	"fmt"
	"mcp-tool-kit/internal/tools"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/server"
)

// EnableCORS is a middleware that enables CORS
func EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// For preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}

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
	fmt.Printf("Loading tools: %s\n", configTools)

	// Initialize tools and track their status
	toolStatus := make(map[string]bool)

	for _, tool := range configTools {
		if tool == "jira" {
			fmt.Println("Initializing Jira tool...")
			tools.NewJiraTool(mcpServer)
			toolStatus["jira"] = true
		}
		
		if tool == "sql-server" {
			fmt.Println("Initializing SQL Server tool...")
			sqlTool := tools.NewSQLServerTool(mcpServer)
			toolStatus["sql-server"] = sqlTool != nil
			if sqlTool != nil {
				fmt.Println("SQL Server tool initialized successfully.")
			} else {
				fmt.Println("Failed to initialize SQL Server tool.")
			}
		}
	}
	
	fmt.Println("MCP server configured successfully.")

	// Create a new ServeMux for easier CORS handling
	mux := http.NewServeMux()

	// Add HTTP endpoints
	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		status := struct {
			Status     string            `json:"status"`
			Tools      map[string]bool   `json:"tools"`
			ServerInfo map[string]string `json:"server_info"`
		}{
			Status: "running",
			Tools:  toolStatus,
			ServerInfo: map[string]string{
				"name":    "mcp-tool-kit",
				"version": "1.0.0",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"%s","tools":%v,"server_info":{"name":"%s","version":"%s"}}`,
			status.Status, toolStatus, status.ServerInfo["name"], status.ServerInfo["version"])
	})

	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"ok","message":"MCP server is running"}`)
	})
	
	// Add SSE endpoint for cursor connection
	mux.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		// Set SSE headers
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		
		// Send initial connection message
		fmt.Fprintf(w, "event: connected\ndata: {\"status\": \"connected\"}\n\n")
		w.(http.Flusher).Flush()
		
		// Keep connection alive with context cancellation
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()
		
		// Create a ticker for heartbeats
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		
		// Listen for client disconnection
		go func() {
			<-ctx.Done()
			ticker.Stop()
		}()
		
		// Send heartbeats
		for {
			select {
			case <-ticker.C:
				// Send heartbeat
				fmt.Fprintf(w, "event: heartbeat\ndata: {\"timestamp\": %d}\n\n", time.Now().Unix())
				w.(http.Flusher).Flush()
			case <-ctx.Done():
				// Client disconnected
				return
			}
		}
	})

	// Apply CORS middleware to all routes
	corsHandler := EnableCORS(mux)

	// Start HTTP server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	fmt.Printf("Starting HTTP server on port %s...\n", port)
	fmt.Printf("Available endpoints:\n")
	fmt.Printf("  - http://localhost:%s/ping\n", port)
	fmt.Printf("  - http://localhost:%s/status\n", port)
	fmt.Printf("  - http://localhost:%s/events (SSE endpoint for Cursor)\n", port)
	
	if err := http.ListenAndServe(":"+port, corsHandler); err != nil {
		fmt.Printf("Error starting HTTP server: %v\n", err)
	}
} 