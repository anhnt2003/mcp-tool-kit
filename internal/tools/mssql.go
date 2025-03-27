package tools

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// SQLServer defines the interface for SQL Server operations
type SQLServer interface {
	ExecuteQuery(ctx context.Context, query string) ([]map[string]any, error)
	GetTables(ctx context.Context) ([]string, error)
	GetTableSchema(ctx context.Context, tableName string) ([]map[string]any, error)
	GetSchema(ctx context.Context) ([]string, error)
}

// sqlServerImpl implements the SQLServer interface
type sqlServerImpl struct {
	db *sql.DB
}

// ExecuteQuery executes a SQL query and returns the results
func (s *sqlServerImpl) ExecuteQuery(ctx context.Context, query string) ([]map[string]any, error) {
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("error getting column names: %w", err)
	}

	// Prepare result slice
	var results []map[string]any

	// Prepare values for scan
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range columns {
		valuePtrs[i] = &values[i]
	}

	// Iterate through rows
	for rows.Next() {
		err = rows.Scan(valuePtrs...)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}

		// Create a map for this row
		row := make(map[string]any)
		for i, col := range columns {
			var v any
			val := values[i]
			
			// Convert bytes to string if needed
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			
			row[col] = v
		}
		
		results = append(results, row)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return results, nil
}

// GetTables returns a list of all tables in the database
func (s *sqlServerImpl) GetTables(ctx context.Context) ([]string, error) {
	query := `
		SELECT TABLE_NAME 
		FROM INFORMATION_SCHEMA.TABLES 
		WHERE TABLE_TYPE = 'BASE TABLE' 
		ORDER BY TABLE_NAME
	`
	
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error getting tables: %w", err)
	}
	defer rows.Close()
	
	var tables []string
	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, fmt.Errorf("error scanning table name: %w", err)
		}
		tables = append(tables, tableName)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tables: %w", err)
	}
	
	return tables, nil
}

// GetTableSchema returns the schema of a specific table
func (s *sqlServerImpl) GetTableSchema(ctx context.Context, tableName string) ([]map[string]any, error) {
	query := `
		SELECT 
			COLUMN_NAME, 
			DATA_TYPE, 
			CHARACTER_MAXIMUM_LENGTH, 
			IS_NULLABLE, 
			COLUMN_DEFAULT 
		FROM INFORMATION_SCHEMA.COLUMNS 
		WHERE TABLE_NAME = @p1 
		ORDER BY ORDINAL_POSITION
	`
	
	return s.ExecuteQuery(ctx, query)
}

// GetSchema returns a list of all schemas in the database
func (s *sqlServerImpl) GetSchema(ctx context.Context) ([]string, error) {
	query := `
		SELECT SCHEMA_NAME 
		FROM INFORMATION_SCHEMA.SCHEMATA 
		ORDER BY SCHEMA_NAME
	`
	
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error getting schemas: %w", err)
	}
	defer rows.Close()
	
	var schemas []string
	for rows.Next() {
		var schemaName string
		if err := rows.Scan(&schemaName); err != nil {
			return nil, fmt.Errorf("error scanning schema name: %w", err)
		}
		schemas = append(schemas, schemaName)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating schemas: %w", err)
	}
	
	return schemas, nil
}

// createConnection establishes a connection to the SQL Server
func createConnection() (*sql.DB, error) {
	server := os.Getenv("SQL_SERVER")
	port := os.Getenv("SQL_PORT")
	user := os.Getenv("SQL_USER")
	password := os.Getenv("SQL_PASSWORD")
	database := os.Getenv("SQL_DATABASE")
	
	// Connection string format for Microsoft's driver - updated for macOS compatibility
	connectionString := fmt.Sprintf("Server=%s,%s;User ID=%s;Password=%s;Database=%s;Encrypt=disable;TrustServerCertificate=true", 
		server, port, user, password, database)
	
	fmt.Printf("Connecting to SQL Server with connection string: %s\n", 
		strings.Replace(connectionString, password, "********", 1))
	
	db, err := sql.Open("mssql", connectionString)
	if err != nil {
		return nil, fmt.Errorf("error connecting to SQL Server: %w", err)
	}
	
	// Test the connection
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error connecting to SQL Server: %w", err)
	}
	
	return db, nil
}

// NewSQLServerTool creates a new instance of SQLServerTool
func NewSQLServerTool(server *server.MCPServer) SQLServer {
	db, err := createConnection()
	if err != nil {
		// Log the error but continue without the tool
		fmt.Printf("Failed to initialize SQL Server tool: %v\n", err)
		return nil
	}

	fmt.Println("SQL Server tool initialized successfully")
	
	// Create a new SQL Server implementation
	sqlServerTool := &sqlServerImpl{
		db: db,
	}
	
	// Add the SQL Server tools to the MCP server
	if server != nil {
		// Register tool for executing SQL queries
		executeQueryTool := mcp.NewTool("sql_execute_query",
			mcp.WithDescription("Execute a SQL query"),
			mcp.WithString("query",
				mcp.Required(),
				mcp.Description("The SQL query to execute"),
			),
		)
		
		server.AddTool(executeQueryTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			query, ok := request.Params.Arguments["query"].(string)
			if !ok {
				return mcp.NewToolResultError("query must be a string"), nil
			}
			
			results, err := sqlServerTool.ExecuteQuery(ctx, query)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			
			// Format results as text
			var resultText strings.Builder
			resultText.WriteString(fmt.Sprintf("Query executed with %d results:\n\n", len(results)))
			
			if len(results) > 0 {
				// Get column names from the first result
				var columns []string
				for col := range results[0] {
					columns = append(columns, col)
				}
				
				// Print column headers
				for _, col := range columns {
					resultText.WriteString(fmt.Sprintf("%s\t", col))
				}
				resultText.WriteString("\n")
				
				// Print separator
				for range columns {
					resultText.WriteString("----------\t")
				}
				resultText.WriteString("\n")
				
				// Print data rows
				for _, row := range results {
					for _, col := range columns {
						resultText.WriteString(fmt.Sprintf("%v\t", row[col]))
					}
					resultText.WriteString("\n")
				}
			}
			
			return mcp.NewToolResultText(resultText.String()), nil
		})
		
		// Register tool for getting all tables
		getTablesTool := mcp.NewTool("sql_get_tables",
			mcp.WithDescription("Get a list of all tables in the database"),
		)
		
		server.AddTool(getTablesTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			tables, err := sqlServerTool.GetTables(ctx)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			
			// Format tables as text
			var resultText strings.Builder
			resultText.WriteString(fmt.Sprintf("Found %d tables:\n\n", len(tables)))
			
			for i, table := range tables {
				resultText.WriteString(fmt.Sprintf("%d. %s\n", i+1, table))
			}
			
			return mcp.NewToolResultText(resultText.String()), nil
		})
		
		// Register tool for getting table schema
		getTableSchemaTool := mcp.NewTool("sql_get_table_schema",
			mcp.WithDescription("Get the schema of a specific table"),
			mcp.WithString("table_name",
				mcp.Required(),
				mcp.Description("The name of the table to get the schema for"),
			),
		)
		
		server.AddTool(getTableSchemaTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			tableName, ok := request.Params.Arguments["table_name"].(string)
			if !ok {
				return mcp.NewToolResultError("table_name must be a string"), nil
			}
			
			schema, err := sqlServerTool.GetTableSchema(ctx, tableName)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			
			// Format schema as text
			var resultText strings.Builder
			resultText.WriteString(fmt.Sprintf("Schema for table %s:\n\n", tableName))
			
			if len(schema) > 0 {
				// Get column names
				var columns []string
				for col := range schema[0] {
					columns = append(columns, col)
				}
				
				// Print headers
				for _, col := range columns {
					resultText.WriteString(fmt.Sprintf("%s\t", col))
				}
				resultText.WriteString("\n")
				
				// Print separator
				for range columns {
					resultText.WriteString("----------\t")
				}
				resultText.WriteString("\n")
				
				// Print data
				for _, row := range schema {
					for _, col := range columns {
						resultText.WriteString(fmt.Sprintf("%v\t", row[col]))
					}
					resultText.WriteString("\n")
				}
			}
			
			return mcp.NewToolResultText(resultText.String()), nil
		})
		
		// Register tool for getting database schemas
		getSchemasTool := mcp.NewTool("sql_get_schemas",
			mcp.WithDescription("Get a list of all schemas in the database"),
		)
		
		server.AddTool(getSchemasTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			schemas, err := sqlServerTool.GetSchema(ctx)
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			
			// Format schemas as text
			var resultText strings.Builder
			resultText.WriteString(fmt.Sprintf("Found %d schemas:\n\n", len(schemas)))
			
			for i, schema := range schemas {
				resultText.WriteString(fmt.Sprintf("%d. %s\n", i+1, schema))
			}
			
			return mcp.NewToolResultText(resultText.String()), nil
		})
	}
	
	return sqlServerTool
}



