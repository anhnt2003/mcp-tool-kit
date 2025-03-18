package tools

import (
	"context"

	"github.com/mark3labs/mcp-go/server"
)

// SQLServer defines the interface for SQL Server operations
type SQLServer interface {
	ExecuteQuery(ctx context.Context, query string) ([]map[string]interface{}, error)
	GetTables(ctx context.Context) ([]string, error)
	GetTableSchema(ctx context.Context, tableName string) ([]map[string]interface{}, error)
	GetSchema(ctx context.Context) ([]string, error)
}

// NewSQLServerTool creates a new instance of SQLServerTool
func NewSQLServerTool(server *server.MCPServer) {
	
}



