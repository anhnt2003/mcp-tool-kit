# MCP Tool Kit

A collection of tools for the Model-Controller-Presenter (MCP) framework.

## Features

- SQL Server integration for database operations
- Jira integration for issue tracking

## Prerequisites

- Go 1.18 or higher
- Docker for running SQL Server

## Setup

### SQL Server Setup for macOS

The project uses SQL Server in a Docker container:

1. Make sure Docker is installed and running on your macOS system
2. Update the credentials in the `.env` file:
   ```
   # SQL Server Configuration
   SQL_SERVER=localhost
   SQL_PORT=1433
   SQL_USER=sa
   SQL_PASSWORD=StrongPassword123!
   SQL_DATABASE=master
   
   # MCP Tools Configuration
   CONFIG_TOOLS=sql-server
   ```

3. Run the application:
   ```
   ./run.sh
   ```

This will:
- Start SQL Server in a Docker container
- Configure the SQL Server connection
- Initialize the MCP server with the SQL Server tools

### SQL Server Tools

The following SQL Server tools are available:

- `sql_execute_query`: Execute a SQL query
- `sql_get_tables`: Get a list of all tables in the database
- `sql_get_table_schema`: Get the schema of a specific table
- `sql_get_schemas`: Get a list of all schemas in the database

### Troubleshooting

If you encounter connection issues:

1. Make sure Docker is running
2. Check if the SQL Server container is running: `docker ps | grep mcp-sqlserver`
3. View container logs: `docker logs mcp-sqlserver`
4. Verify the connection parameters in the `.env` file

## Development

### Adding New Tools

To add a new tool:

1. Create a new implementation in the `internal/tools` directory
2. Register the tool in `main.go`
3. Update the `CONFIG_TOOLS` environment variable to include the new tool

## Project Structure

```
mcp-tool-kit/
├── main.go          # Main server implementation
├── go.mod           # Go module definition
└── README.md        # Project documentation
```

## License

MIT License 