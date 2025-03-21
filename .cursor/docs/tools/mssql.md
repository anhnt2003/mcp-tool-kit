# SQL Server Tool

The SQL Server tool provides functionality to interact with Microsoft SQL Server databases. It allows executing queries, retrieving schema information, and exploring database objects.

## Overview

The tool implements a common database interface and provides several capabilities:

- Executing SQL queries
- Retrieving lists of database tables
- Getting schema information for specific tables
- Retrieving database schema information

## Configuration

The SQL Server tool requires the following environment variables to be set:

| Variable | Description |
|----------|-------------|
| `SQL_SERVER` | SQL Server hostname or IP address |
| `SQL_PORT` | SQL Server port (typically 1433) |
| `SQL_USER` | Username for SQL Server authentication |
| `SQL_PASSWORD` | Password for SQL Server authentication |
| `SQL_DATABASE` | Default database name |

## Connection

The tool establishes a connection to SQL Server using the provided credentials. The connection string is formatted as:

```
sqlserver://[username]:[password]@[server]:[port]?database=[database]&encrypt=disable&TrustServerCertificate=true
```

For security purposes, SSL encryption is disabled by default, and the tool trusts the server certificate.

## Available Functions

### ExecuteQuery

Executes a SQL query and returns the results as a collection of row maps.

```go
ExecuteQuery(ctx context.Context, query string) ([]map[string]any, error)
```

**Parameters:**
- `ctx`: Context for query execution
- `query`: SQL query to execute

**Returns:**
- Array of maps, where each map represents a row with column names as keys
- Error if the query execution fails

### GetTables

Retrieves a list of all tables in the database.

```go
GetTables(ctx context.Context) ([]string, error)
```

**Parameters:**
- `ctx`: Context for query execution

**Returns:**
- Array of table names
- Error if the operation fails

### GetTableSchema

Retrieves the schema information for a specific table.

```go
GetTableSchema(ctx context.Context, tableName string) ([]map[string]any, error)
```

**Parameters:**
- `ctx`: Context for query execution
- `tableName`: Name of the table to get schema for

**Returns:**
- Array of maps containing column information (name, data type, length, nullability, default value)
- Error if the operation fails

### GetSchema

Retrieves a list of all schemas in the database.

```go
GetSchema(ctx context.Context) ([]string, error)
```

**Parameters:**
- `ctx`: Context for query execution

**Returns:**
- Array of schema names
- Error if the operation fails

## MCP Tools

When initialized with an MCP server, the SQL Server tool registers the following tools:

### sql_execute_query

Executes a SQL query and returns the results in a formatted table.

**Parameters:**
- `query`: The SQL query to execute (required)

**Example:**
```
sql_execute_query(query="SELECT TOP 10 * FROM Customers")
```

### sql_get_tables

Returns a list of all tables in the database.

**Parameters:**
- None

**Example:**
```
sql_get_tables()
```

### sql_get_table_schema

Returns the schema of a specific table.

**Parameters:**
- `table_name`: The name of the table to get the schema for (required)

**Example:**
```
sql_get_table_schema(table_name="Customers")
```

### sql_get_schemas

Returns a list of all schemas in the database.

**Parameters:**
- None

**Example:**
```
sql_get_schemas()
```

## Implementation Details

The SQL Server tool internally uses Go's standard `database/sql` package with the Microsoft SQL Server driver. Results from queries are transformed into maps for easier consumption by other tools and services.

For handling binary data returned from the database, the tool automatically converts byte arrays to strings when appropriate.

## Error Handling

All functions return detailed error messages if operations fail. Errors are wrapped with context information to help diagnose issues.

Common error scenarios include:
- Connection failures
- Authentication issues
- Query syntax errors
- Permission problems
