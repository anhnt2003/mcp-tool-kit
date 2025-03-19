#!/bin/bash

# Load environment variables from .env file
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# Set default port if not specified
if [ -z "$PORT" ]; then
    export PORT=8080
fi

# Start SQL Server container
echo "Starting SQL Server container..."
./start-sqlserver.sh

# Run the application
echo "Starting the application..."
go run cmd/server/main.go 