// Package interfaces provides interfaces for database connections and operations
package interfaces

// Database is a common interface for all database connections
type Database interface {
	// Connect establishes a connection to the database
	Connect() error

	// Disconnect closes the connection to the database
	Disconnect() error

	// Query executes a query and returns results
	// query: SQL query to execute
	// params: Parameters for the query
	Query(query string, params ...any) ([]map[string]any, error)

	// Execute runs a query that doesn't return results (INSERT, UPDATE, DELETE)
	// query: SQL query to execute
	// params: Parameters for the query
	Execute(query string, params ...any) error

	// GetSchema returns database schema information
	GetSchema() (SchemaInfo, error)

	// GetTables returns all table names
	GetTables() ([]string, error)

	// GetTableSchema returns column information for a specific table
	// tableName: Name of the table
	GetTableSchema(tableName string) (TableSchema, error)
}

// SchemaInfo contains database schema information
type SchemaInfo struct {
	// DatabaseName is the name of the database
	DatabaseName string

	// Tables contains information about all tables in the database
	Tables []TableSchema
}

// TableSchema contains table schema information
type TableSchema struct {
	// TableName is the name of the table
	TableName string

	// Columns contains information about all columns in the table
	Columns []ColumnInfo
}

// ColumnInfo contains information about a database column
type ColumnInfo struct {
	// Name is the column name
	Name string

	// Type is the column data type
	Type string

	// Nullable indicates whether the column can contain NULL values
	Nullable bool

	// IsPrimaryKey indicates whether the column is part of the primary key
	IsPrimaryKey bool

	// DefaultValue is the default value for the column (if any)
	DefaultValue interface{}
} 