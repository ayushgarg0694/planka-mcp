# Cursor MCP Configuration for Planka

## Configuration File Location

Cursor MCP configuration is typically located at:
- **macOS/Linux**: `~/.cursor/mcp.json` or `~/.config/cursor/mcp.json`
- **Windows**: `%APPDATA%\Cursor\mcp.json` or `C:\Users\<username>\AppData\Roaming\Cursor\mcp.json`

## Full Configuration

Copy the following configuration to your Cursor MCP config file:

### Option 1: Using Username/Password (Recommended for testing)

```json
{
  "mcpServers": {
    "planka": {
      "command": "/path/to/your/mcp-planka",
      "env": {
        "PLANKA_URL": "https://planka.cosmicdragon.xyz",
        "PLANKA_USERNAME": "your-username",
        "PLANKA_PASSWORD": "your-username"
      }
    }
  }
}
```

### Option 2: Using API Token (Recommended for production)

```json
{
  "mcpServers": {
    "planka": {
      "command": "/path/to/your/mcp-planka",
      "env": {
        "PLANKA_URL": "https://planka.cosmicdragon.xyz",
        "PLANKA_TOKEN": "your-api-token-here"
      }
    }
  }
}
```

## Setup Instructions

1. **Build the MCP server** (if not already built):
   ```bash
   cd /path/to/your
   go build -o mcp-planka .
   ```

2. **Make sure the binary is executable**:
   ```bash
   chmod +x mcp-planka
   ```

3. **Update the command path** in the config to match your actual path:
   - Replace `/path/to/your/mcp-planka` with your actual path
   - Or use an absolute path to where you've installed the binary

4. **Add the configuration to Cursor**:
   - Open Cursor settings
   - Navigate to MCP settings
   - Add the configuration above
   - Or manually edit the MCP config file at the location mentioned above

5. **Restart Cursor** to load the new MCP server

## Testing the Configuration

After adding the configuration, you can test it by:

1. Opening Cursor
2. Using the MCP tools in a chat session
3. Try commands like:
   - "List all my Planka projects"
   - "Show me the boards in my Wedding Project"
   - "What cards are in the TODO list?"

## Available MCP Tools

Once configured, you'll have access to these tools:

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

## Troubleshooting

### Server not found
- Make sure the path to `mcp-planka` is correct and absolute
- Verify the binary exists and is executable: `ls -la /path/to/your/mcp-planka`

### Authentication errors
- Check that `PLANKA_URL` is correct
- Verify username/password or token are correct
- Test authentication manually: `./mcp-planka test`

### Connection issues
- Ensure your Planka instance is accessible
- Check network connectivity
- Verify the URL doesn't have a trailing slash

## Security Notes

- **For production**: Use API tokens instead of username/password
- **Never commit** your MCP config file with credentials to version control
- Consider using environment variables or a secrets manager for sensitive data

