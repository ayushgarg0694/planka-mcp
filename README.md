# Planka MCP Server (Go)

A Model Context Protocol (MCP) server implementation in Go that provides an interface between Large Language Models (LLMs) and [Planka](https://github.com/plankanban/planka), an open-source Kanban board application.

This implementation is inspired by:
- [plankapy](https://github.com/hwelch-fle/plankapy) - Python API for Planka
- [kanban-mcp](https://github.com/bradrisse/kanban-mcp) - Node.js MCP server for Planka

## Features

### 📋 Project & Board Management
- List and view projects
- Create and manage boards
- Organize work across multiple projects

### 📝 List & Card Operations
- Create and manage lists within boards
- Create, update, and delete cards
- Move cards between lists
- Add descriptions, due dates, and labels

### ⏱️ Time Tracking
- Start, stop, and reset stopwatches
- Track time spent on individual tasks
- Analyze time usage patterns

### ✅ Task Management
- Create and manage tasks within cards
- Mark tasks as complete or incomplete
- Organize tasks with positions

### 💬 Comment Management
- Add comments to cards for discussion
- View comment history
- Delete comments

## Installation

### Prerequisites

- Go 1.21 or later
- Access to a Planka instance
- Either a Planka API token OR username/password credentials

### Building

```bash
git clone https://github.com/ayushgarg/mcp-planka.git
cd mcp-planka
go build -o mcp-planka
```

## Configuration

The server requires the following environment variables:

- `PLANKA_URL`: The base URL of your Planka instance (e.g., `https://planka.example.com`)

**Authentication (choose one method):**

### Option 1: API Token (Recommended)

- `PLANKA_TOKEN`: Your Planka API authentication token

**Getting a Planka API Token:**

1. Log in to your Planka instance
2. Navigate to your user settings
3. Generate an API token
4. Use this token as the `PLANKA_TOKEN` environment variable

### Option 2: Username/Password

- `PLANKA_USERNAME`: Your Planka username
- `PLANKA_PASSWORD`: Your Planka password

**Note:** The server will automatically authenticate using username/password if `PLANKA_TOKEN` is not provided. The token will be obtained automatically during login.

## Usage

The server supports two modes of operation:

### Stdio Mode (Default)

The server communicates via stdio using the MCP protocol. This is the default mode and is typically used by MCP clients like Cursor.

#### Using API Token

```bash
export PLANKA_URL="https://planka.example.com"
export PLANKA_TOKEN="your-api-token-here"
./mcp-planka
```

#### Using Username/Password

```bash
export PLANKA_URL="https://planka.example.com"
export PLANKA_USERNAME="your-username"
export PLANKA_PASSWORD="your-password"
./mcp-planka
```

### HTTP Server Mode

The server can also run as an HTTP server that accepts JSON-RPC 2.0 requests over HTTP. This is useful for web clients or remote access.

#### Starting the HTTP Server

```bash
export PLANKA_URL="https://planka.example.com"
export PLANKA_USERNAME="your-username"
export PLANKA_PASSWORD="your-password"

# Start HTTP server on default port 8080
./mcp-planka --http

# Start HTTP server on custom port
./mcp-planka --http --http-port 3000

# Start HTTP server on specific address and port
./mcp-planka --http --http-addr 127.0.0.1 --http-port 8080
```

#### Command-Line Flags

- `--http` - Enable HTTP server mode (default: false, uses stdio)
- `--http-port` - HTTP server port (default: 8080)
- `--http-addr` - HTTP server bind address (default: "0.0.0.0")

#### HTTP Endpoints

**POST /mcp** or **POST /** - Main JSON-RPC endpoint
- Accepts JSON-RPC 2.0 requests in the request body
- Returns JSON-RPC 2.0 responses
- Example request:
  ```json
  {
    "jsonrpc": "2.0",
    "method": "tools/list",
    "id": 1
  }
  ```

**GET /health** - Health check endpoint
- Returns server status
- Example response:
  ```json
  {
    "status": "ok",
    "service": "planka-mcp"
  }
  ```

#### Example HTTP Usage

```bash
# Initialize the server
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "initialize",
    "params": {},
    "id": 1
  }'

# List available tools
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "tools/list",
    "id": 2
  }'

# Call a tool
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "tools/call",
    "params": {
      "name": "get_projects",
      "arguments": {}
    },
    "id": 3
  }'
```

#### HTTP Mode Test Script

A comprehensive test script is provided to test all HTTP endpoints:

```bash
# Make sure the server is running in HTTP mode first
./mcp-planka --http --http-port 8080

# In another terminal, run the test script
./test_http.sh

# Or specify a custom base URL
BASE_URL=http://localhost:9000 ./test_http.sh
```

The test script (`test_http.sh`) includes:
- Health check endpoint
- Initialize and initialized notification
- List all available tools
- Get projects, project details, and boards
- CORS preflight requests
- Error handling (invalid methods, malformed JSON)

**Prerequisites:** The script requires `jq` for JSON parsing:
```bash
# Ubuntu/Debian
sudo apt-get install jq

# macOS
brew install jq
```

**Example Output:**
```
╔══════════════════════════════════════════════════════════════════════════════╗
║                    Planka MCP HTTP Mode Test Suite                          ║
╚══════════════════════════════════════════════════════════════════════════════╝
Testing server at: http://localhost:8080
MCP endpoint: http://localhost:8080/mcp

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Test: Health Check (GET /health)
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
✓ Health check passed
{
  "status": "ok",
  "service": "planka-mcp"
}
...
```

### Running Tests

To test the Planka API connection and verify all endpoints are working:

1. Create a test file (or use the provided `test.go`):

```bash
# The test file should be in a separate directory or temporarily rename main.go
# For example, create test/test.go with your test code
```

2. Run the test (temporarily rename `main.go` to avoid conflicts):

```bash
# Option 1: Temporarily rename main.go
mv main.go main.go.bak
go run test.go
mv main.go.bak main.go

# Option 2: Create test in separate package (requires exporting internal packages)
```

3. Or create a simple test script:

```bash
# Create a test script that uses the planka client
cat > test_connection.sh << 'EOF'
#!/bin/bash
export PLANKA_URL="https://planka.example.com"
export PLANKA_USERNAME="your-username"
export PLANKA_PASSWORD="your-password"

# Test authentication
go run -exec 'mv main.go main.go.bak && go run test.go && mv main.go.bak main.go'
EOF

chmod +x test_connection.sh
./test_connection.sh
```

**Example Test Output:**

```
======================
Testing Planka MCP Server
======================

[Test 1] Authenticating...
✓ Successfully authenticated!

[Test 2] Getting current user...
✓ Logged in as: username (email@example.com)

[Test 3] Getting all projects...
✓ Found 6 project(s):
  1. Project 1 (ID: ...)
  2. Project 2 (ID: ...)
  ...
```

### MCP Client Configuration

#### Cursor Configuration (Stdio Mode)

For Cursor, add the following to your MCP configuration:

#### Using API Token

```json
{
  "mcpServers": {
    "planka": {
      "command": "/path/to/mcp-planka",
      "env": {
        "PLANKA_URL": "https://planka.example.com",
        "PLANKA_TOKEN": "your-api-token-here"
      }
    }
  }
}
```

#### Using Username/Password

```json
{
  "mcpServers": {
    "planka": {
      "command": "/path/to/mcp-planka",
      "env": {
        "PLANKA_URL": "https://planka.example.com",
        "PLANKA_USERNAME": "your-username",
        "PLANKA_PASSWORD": "your-password"
      }
    }
  }
}
```

**Note:** For security reasons, prefer using API tokens in production environments. Username/password authentication is convenient for development and testing.

#### HTTP Server Configuration

When running in HTTP mode, you can configure clients to connect via HTTP:

```bash
# Start the HTTP server
./mcp-planka --http --http-port 8080

# Clients can then make HTTP requests to http://localhost:8080/mcp
```

For web applications, the server includes CORS headers to allow cross-origin requests.

## Available Tools

The server provides the following MCP tools:

### Projects
- `get_projects` - Get all projects
- `get_project` - Get a project by ID
- `create_project` - Create a new project

### Boards
- `get_boards` - Get all boards for a project
- `get_board` - Get a board by ID
- `create_board` - Create a new board

### Lists
- `get_lists` - Get all lists for a board
- `get_list` - Get a list by ID
- `create_list` - Create a new list

### Cards
- `get_cards` - Get all cards for a list
- `get_card` - Get a card by ID
- `create_card` - Create a new card
- `update_card` - Update a card
- `delete_card` - Delete a card
- `move_card` - Move a card to a different list

### Tasks
- `get_tasks` - Get all tasks for a card
- `create_task` - Create a new task
- `update_task` - Update a task
- `delete_task` - Delete a task

### Comments
- `get_comments` - Get all comments for a card
- `create_comment` - Create a new comment
- `delete_comment` - Delete a comment

### Time Tracking
- `get_stopwatch` - Get the stopwatch for a card
- `start_stopwatch` - Start the stopwatch for a card
- `stop_stopwatch` - Stop the stopwatch for a card
- `reset_stopwatch` - Reset the stopwatch for a card

## Development

### Project Structure

```
mcp-planka/
├── main.go                 # Entry point
├── test.go                 # Integration test file (optional)
├── internal/
│   ├── planka/            # Planka API client
│   │   ├── client.go      # HTTP client implementation
│   │   ├── models.go      # Data models
│   │   └── api.go         # API methods
│   └── mcp/               # MCP server implementation
│       ├── server.go      # MCP protocol handling
│       └── tools.go       # Tool definitions and handlers
├── go.mod
└── README.md
```

### Local Development Setup

1. **Clone the repository:**
```bash
git clone https://github.com/ayushgarg/mcp-planka.git
cd mcp-planka
```

2. **Install dependencies:**
```bash
go mod download
```

3. **Build the server:**
```bash
go build -o mcp-planka .
```

4. **Set up environment variables:**
```bash
export PLANKA_URL="https://planka.example.com"
export PLANKA_USERNAME="your-username"
export PLANKA_PASSWORD="your-password"
# OR
export PLANKA_TOKEN="your-api-token"
```

5. **Run the server:**
```bash
./mcp-planka
```

6. **Test the connection:**
```bash
# Create test.go with your test code, then:
mv main.go main.go.bak
go run test.go
mv main.go.bak main.go
```

### Running Tests

#### Unit Tests

```bash
go test ./...
```

#### Integration Tests

To test the actual Planka API connection:

1. Set up your test environment variables:
```bash
export PLANKA_URL="https://planka.example.com"
export PLANKA_USERNAME="your-username"
export PLANKA_PASSWORD="your-password"
```

2. Create a test file (e.g., `test.go`) with your test code

3. Run the test (temporarily rename `main.go`):
```bash
mv main.go main.go.bak
go run test.go
mv main.go.bak main.go
```

**Note:** The test file should import the internal packages and test the Planka client functionality. Make sure to test with a non-production Planka instance if possible.

### Building for Different Platforms

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o mcp-planka-linux

# macOS
GOOS=darwin GOARCH=amd64 go build -o mcp-planka-macos

# Windows
GOOS=windows GOARCH=amd64 go build -o mcp-planka.exe
```

## Security Considerations

- **Input Validation**: All inputs are validated before being sent to the Planka API
- **Token Security**: Never commit your API token to version control
- **HTTPS**: Always use HTTPS for your Planka instance URL
- **Error Handling**: The server includes comprehensive error handling to prevent information leakage

## Limitations

- This implementation is based on Planka API v1.x. Some features may not work with Planka v2.x instances.
- The server supports both Bearer token and username/password authentication methods.
- Some endpoints may return HTML instead of JSON if the API structure differs between Planka versions.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is open-source. Please check the LICENSE file for details.

## Acknowledgments

- [Planka](https://github.com/plankanban/planka) - The Kanban board application
- [plankapy](https://github.com/hwelch-fle/plankapy) - Python API reference
- [kanban-mcp](https://github.com/bradrisse/kanban-mcp) - Node.js MCP server reference

