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

#### sql_get_tables

Returns a list of all tables in the database.

#### sql_get_table_schema

Returns the schema of a specific table.

#### sql_get_schemas

Returns a list of all schemas in the database.

### Jira Tools

Documentation for Jira tools will be added soon.