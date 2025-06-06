# Blue Prince MCP Tools CLI

A command-line interface for testing the Blue Prince MCP server tools locally. This CLI connects to the running MCP server and allows you to test create, read, update, and list operations on your notes.

## Prerequisites

1. **MCP Server Running**: Make sure the MCP server is running:
   ```bash
   go run ./cmd/server/main.go
   ```

2. **Configuration**: Ensure `cmd/config/local/config.yaml` contains the correct server settings:
   ```yaml
   server:
     host: localhost
     port: 8001
   obsidian_vault_path: /path/to/your/vault
   ```

## Building the CLI

```bash
go build -o ./bin/blueprince-tools ./cmd/tools/
```

## Available Commands

### 1. List Notes
```bash
# List all notes in the vault
./bin/blueprince-tools list

# With verbose output
./bin/blueprince-tools list --verbose
```

### 2. Read Note
```bash
# Read a specific note
./bin/blueprince-tools read people/simon_jones.md

# With verbose output
./bin/blueprince-tools read rooms/library.md --verbose
```

### 3. Create Note
```bash
# Simple create with defaults
./bin/blueprince-tools create people/new_character.md --content "A mysterious figure"

# Full create with all metadata
./bin/blueprince-tools create rooms/study.md \
  --title "Study Room" \
  --category "rooms" \
  --primary-subject "study" \
  --tags "rooms,books,desk,quiet" \
  --confidence "high" \
  --status "confirmed" \
  --content "# Study Room\n\nA quiet room with a large desk and bookshelves."
```

### 4. Update Note
```bash
# Update existing note (content required)
./bin/blueprince-tools update people/character.md \
  --content "# Updated Character\n\nNew information about this character."

# Update with new metadata
./bin/blueprince-tools update rooms/library.md \
  --title "Grand Library" \
  --tags "rooms,books,secret_passage" \
  --status "confirmed" \
  --content "# Grand Library\n\nLarge library with hidden passages."
```

## Global Flags

- `--config`: Path to config file (default: `cmd/config/local/config.yaml`)
- `--host`: Override server host
- `--port`: Override server port  
- `--verbose`: Enable verbose output for debugging

## Examples

Run the test script to see all commands in action:

```bash
./cmd/tools/test_examples.sh
```

## Metadata Fields

### Categories
- `people` - Characters and NPCs
- `rooms` - Locations and areas  
- `items` - Objects and artifacts
- `lore` - Background information
- `puzzles` - Riddles and challenges
- `general` - Miscellaneous notes

### Confidence Levels
- `high` - Very certain information
- `medium` - Somewhat certain
- `low` - Uncertain or speculative

### Status Values
- `complete` - Fully documented
- `needs_investigation` - Requires more info
- `active_investigation` - Currently being researched
- `theory` - Speculative information
- `confirmed` - Verified information

## Troubleshooting

1. **Connection refused**: Make sure the MCP server is running
2. **Config not found**: Check the config file path
3. **Tool errors**: Check server logs for detailed error messages
4. **Invalid metadata**: Ensure category, confidence, and status values are valid

## Development

The CLI is organized into separate files for maintainability:

- `main.go` - Main CLI setup and root command
- `client.go` - HTTP client for MCP communication
- `list.go` - List notes command
- `read.go` - Read note command  
- `create.go` - Create note command
- `update.go` - Update note command
- `test_examples.sh` - Test script with examples