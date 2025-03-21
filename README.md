# MCP Tool Kit

A comprehensive toolkit for the Model Context Protocol (MCP) that provides integrations with various database systems and development tools.

## Features

- **SQL Server Integration**: Seamless interaction with Microsoft SQL Server databases
  - Execute custom SQL queries
  - Retrieve table and schema information
  - Explore database structure
- **Jira Integration**: Powerful issue tracking capabilities
- **Extensible Architecture**: Easily add new tools and integrations

## Prerequisites

- Go 1.18 or higher
- Docker (for running SQL Server locally)
- SQL Server instance (local or remote)
- Environment variables configured for database access

## Setup

### Environment Configuration

Create a `.env` file in the project root with the following configuration:

```
# SQL Server Configuration
SQL_SERVER=your-server-address
SQL_PORT=1433
SQL_USER=sa
SQL_PASSWORD=YourStrongPassword!
SQL_DATABASE=your-database-name
SQL_MULTIPLE_ACTIVE_RESULT_SETS=True
SQL_TRUST_SERVER_CERTIFICATE=True

# Optional: complete connection string (will be used if provided)
SQL_CONNECTION_STRING=Server=your-server-address;Database=your-database-name;User Id=sa;Password=YourStrongPassword!;MultipleActiveResultSets=True;TrustServerCertificate=True

# MCP Tools Configuration
CONFIG_TOOLS=sql-server
```

### SQL Server Setup

#### Using Docker (macOS/Linux/Windows)

1. Make sure Docker is installed and running on your system
2. Run SQL Server in a Docker container:

```bash
docker run -e "ACCEPT_EULA=Y" -e "SA_PASSWORD=YourStrongPassword!" -p 1433:1433 --name mcp-sqlserver -d mcr.microsoft.com/mssql/server:2019-latest
```

3. Update the credentials in the `.env` file to match your Docker container configuration

#### Using Existing SQL Server

If you already have a SQL Server instance running:

1. Ensure you have the necessary credentials and permissions
2. Update the `.env` file with the connection details

### Running the Application

Execute the application using the provided script:

```bash
./run.sh
```

This will:
- Load environment variables from the `.env` file
- Initialize the MCP server with configured tools
- Start the server and make tools available

## Available Tools

### SQL Server Tools

The MCP Tool Kit provides the following SQL Server tools:

#### sql_execute_query

Executes a SQL query and returns the results in a formatted table.

**Parameters:**
- `query`: The SQL query to execute (required)

**Example:**
```
sql_execute_query(query="SELECT TOP 10 * FROM Customers")
```

#### sql_get_tables

Returns a list of all tables in the database.

**Parameters:**
- None

**Example:**
```
sql_get_tables()
```

#### sql_get_table_schema

Returns the schema of a specific table.

**Parameters:**
- `table_name`: The name of the table to get the schema for (required)

**Example:**
```
sql_get_table_schema(table_name="Customers")
```

#### sql_get_schemas

Returns a list of all schemas in the database.

**Parameters:**
- None

**Example:**
```
sql_get_schemas()
```

### Jira Tools

Documentation for Jira tools will be added soon.

## Troubleshooting

### SQL Server Connection Issues

If you encounter connection issues with SQL Server:

1. Verify that SQL Server is running:
   - If using Docker: `docker ps | grep mcp-sqlserver`
   - If using a local instance: Check SQL Server services

2. Check connection parameters:
   - Ensure environment variables are correctly set
   - Verify network connectivity to the SQL Server instance
   - Check firewall settings

3. View logs for more information:
   - Docker container logs: `docker logs mcp-sqlserver`
   - Application logs in the console output

4. Common solutions:
   - Restart the SQL Server container/service
   - Update connection string parameters
   - Check for conflicting port usage

## Development

### Project Structure

```
mcp-tool-kit/
├── .cursor/             # Documentation directory
│   └── docs/
│       └── tools/       # Tool-specific documentation
├── internal/            # Internal packages
│   ├── interfaces/      # Interface definitions
│   └── tools/           # Tool implementations
├── main.go              # Main server implementation
├── go.mod               # Go module definition
├── go.sum               # Go module checksums
└── README.md            # Project documentation
```

### Adding New Tools

To add a new tool to the MCP Tool Kit:

1. Create a new implementation in the `internal/tools` directory
2. Implement the appropriate interface(s) from `internal/interfaces`
3. Register the tool in `main.go`
4. Update the `CONFIG_TOOLS` environment variable to include the new tool
5. Add documentation in the `.cursor/docs/tools` directory

### Database Implementation Guidelines

When implementing database tools:

1. Follow the `Database` interface defined in `internal/interfaces/database.go`
2. Implement proper error handling and context management
3. Ensure resource cleanup (connection closing, etc.)
4. Add appropriate documentation

## Documentation

Detailed documentation for each tool is available in the `.cursor/docs/tools` directory:

- [SQL Server Tool Documentation](.cursor/docs/tools/mssql.md)
- More tool documentation coming soon
