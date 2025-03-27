# MCP Tool Kit

A comprehensive toolkit for the Model Context Protocol (MCP) that provides integrations with various database systems and development tools.

## Features

- **SQL Server Integration**: Seamless interaction with Microsoft SQL Server databases
  - Execute custom SQL queries
  - Retrieve table and schema information
  - Explore database structure
- **Jira Integration**: Powerful issue tracking capabilities
- **Extensible Architecture**: Easily add new tools and integrations
- **Multiple Operation Modes**: Run in stdio or SSE server mode

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

# Optional: complete connection string (will be used if provided)
SQL_CONNECTION_STRING=Server=your-server-address;Database=your-database-name;User Id=sa;Password=YourStrongPassword!;MultipleActiveResultSets=True;TrustServerCertificate=True

# Jira Configuration
JIRA_API_KEY=your-jira-api-key
JIRA_URL=https://your-domain.atlassian.net
JIRA_EMAIL=your-email@domain.com

# MCP Tools Configuration
CONFIG_TOOLS=sql-server,jira

# MCP Mode Configuration
MCP_MODE=stdio  # Options: stdio, sse
```
## Operation Modes

### STDIO Mode (Default)

Standard I/O mode is the default operation mode. In this mode, the tool communicates through standard input and output streams, making it ideal for command-line integration.

To run in STDIO mode: Claude, cursor

```
{
  "mcpServers": {
    "script": {
      "command": "/path/to/your/mcp-script",
      "args": ["-env", "path-to-env-file"]
    }
  }
}
```

### SSE Mode

Server-Sent Events (SSE) mode provides a web-based interface for the toolkit. It starts a web server that allows communication through HTTP requests and SSE for real-time updates.

To run in SSE mode:
```
MCP_MODE=sse go run main.go
```

This starts an SSE server on http://localhost:8080 with the following endpoints:
- `/message` - For sending messages to the server
- `/sse` - For receiving SSE events from the server

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

The MCP Tool Kit provides the following Jira tools:

#### jira_get_issue

Retrieves detailed information about a specific Jira issue, including status, assignee, description, subtasks, and available transitions.

#### jira_search_issue

Searches for Jira issues using JQL (Jira Query Language) and returns details such as summary, status, assignee, and priority.

#### jira_list_sprints

Lists all active and future sprints for a given Jira board.

#### jira_create_issue

Creates a new Jira issue and returns the created issue's key, ID, and URL.

#### jira_update_issue

Updates an existing Jira issue with new details. Only provided fields will be updated.

#### jira_list_statuses

Retrieves all available issue statuses for a specific Jira project.

#### jira_transition_issue

Transitions an issue through its workflow using a valid transition ID.

