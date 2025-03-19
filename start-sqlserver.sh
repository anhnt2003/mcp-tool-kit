#!/bin/bash

# Stop and remove existing container if it exists
if [ $(docker ps -a -q -f name=mcp-sqlserver) ]; then
    echo "Stopping and removing existing SQL Server container..."
    docker stop mcp-sqlserver
    docker rm mcp-sqlserver
fi

# Create and start new SQL Server container
echo "Creating and starting new SQL Server container..."
docker run \
    --name mcp-sqlserver \
    -e "ACCEPT_EULA=Y" \
    -e "MSSQL_SA_PASSWORD=StrongPassword123!" \
    -p 1433:1433 \
    -d \
    --platform linux/amd64 \
    mcr.microsoft.com/mssql/server:2019-latest

# Wait for SQL Server to start
echo "Waiting for SQL Server to start..."
sleep 15  # Give it more time to initialize

# Check if SQL Server is running
docker ps | grep mcp-sqlserver
if [ $? -eq 0 ]; then
    echo "SQL Server is running."
    echo "Container logs:"
    docker logs mcp-sqlserver | tail -n 10
else
    echo "SQL Server failed to start. Check Docker logs with: docker logs mcp-sqlserver"
fi 