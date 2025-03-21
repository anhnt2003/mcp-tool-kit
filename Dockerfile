# Use Go 1.22 as the base image
FROM golang:1.22-alpine AS builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o mcp-tool-kit .

# Use a smaller base image for the final image
FROM alpine:latest

# Install CA certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/mcp-tool-kit .
# Copy .env file
COPY .env .

# Set environment variables
ENV CONFIG_TOOLS="jira,sql-server"

# Expose port 8080
EXPOSE 8080

# Run the application
CMD ["./mcp-tool-kit"] 