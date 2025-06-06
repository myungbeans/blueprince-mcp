# Blue Prince Architect Notes (MCP Server)

This repository contains the code for an MCP (Multi-Agent Communication Protocol) server designed to act as a dedicated note-taking and brainstorming assistant specifically tailored for playing the game **Blue Prince**.

This MCP server exposes tools and resources for managing local notes (stored as .md files) that allow users to write notes, lookup information from their notes, and brainstorm with a companion MCP client as they play through the video game Blue Prince. This is designed to help players make connections and recall things they've seen and experienced while avoiding spoilers from online resources.

**⚠️ IMPORTANT: SPOILER-FREE USAGE**
This MCP server is designed to preserve your Blue Prince gameplay experience. When using with an MCP client (e.g. Claude Desktop):
- The Client will ONLY use information from your notes
- The Client cannot and will not reference external Blue Prince information  
- Spoiler prevention rules are automatically provided as an MCP resource
- The Client will have access to explicit spoiler prevention guidelines

## Features

- **MCP Server:** Implements the MCP protocol to expose note-taking capabilities as tools and resources.
- **Local Vault Storage:** Stores notes as markdown files in a structured local directory (compatible with Obsidian).
- **Structured Notes:** Organizes notes in predefined categories (`people`, `puzzles`, `rooms`, `items`, `lore`, `general`) with intelligent metadata extraction.
- **Resource System:** Exposes all vault files as MCP resources for direct access by AI clients (excludes `.obsidian/` directories).
- **Spoiler Prevention System:** Multi-layered protection including:
  - Built-in content validation to prevent investigation prompts
  - Spoiler prevention rules automatically exposed as an MCP resource
  - Client-side enforcement through tool descriptions and server metadata
  - Server-side validation of all content creation
- **Complete CRUD Operations:** 
  - ✅ `list_notes` - Lists all notes in the vault
  - ✅ `create_note` - Creates structured notes with intelligent categorization and spoiler prevention
  - ✅ `read_note` - Reads complete note content including metadata
  - ✅ `update_note` - Updates existing notes with new content
  - 📋 `delete_note` - Planned for future implementation
- **CLI Testing Tools:** Comprehensive command-line interface for manual testing and debugging.
- **Setup Utility:** Go program to initialize vault directory structure and configuration.
- **Flexible Configuration:** Supports both file-based config and environment variable overrides.
- **Structured Logging:** Uses `go.uber.org/zap` for comprehensive logging and debugging.

## Usage Guide
"Write a new note. I'm in the corridor. There is a painting of a tiger and a cupcake stand? Three windows. Two benches and hats."

The MCP Client will then intelligently format the note and tag it appropraitely.

In the future you can ask
"Where have I seen windows before? Can you list all rooms that have windows in them?"

The MCP Client will then scan all of your notes (and only your notes) to look for what _you_ know to be all rooms that have windows in them.

## Getting Started

### Prerequisites

- Go (version 1.20 or higher recommended)
- An MCP client (e.g., a compatible AI agent or a testing tool)
- Git

### Setup

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/myungbeans/blueprince-mcp.git
    cd blueprince-mcp
    ```

2.  **Run the setup utility:**
    This program will create the necessary Obsidian vault directory structure and update your `config.yaml` file with the vault's path.

    By default, it will create the vault at `~/Documents/blueprince_mcp/`.

    ```bash
    go run ./cmd/setup
    ```

    Alternatively, you can specify a custom path for your vault:

    ```bash
    go run ./cmd/setup /path/to/your/custom/vault
    ```

    The setup utility will ensure the required subdirectories (`notes/people`, `notes/puzzles`, `notes/rooms`, `notes/items`, `notes/lore`, `notes/general`) exist within the vault, along with `meta/` and `screenshots/` directories.

3.  **Review `config.yaml`:**
    The setup utility updates `config.yaml` with the `obsidian_vault_path`. You can review this file and adjust other settings like `server.host` or `server.port` if needed.

    ```yaml
    # Example config.yaml
    server:
      host: "localhost"
      port: 8001

    obsidian_vault_path: "/Users/michael.myung/Documents/blueprince_mcp" # This will be set by the setup script
    backup_dir_name: ".obsidian_backup" # Directory name for potential future backups within the vault
    ```

## Build & Usage

### Building the Server
```bash
go build -o ./bin/blueprince-mcp-server ./cmd/server/main.go
```

### Building the CLI Tools
```bash
go build -o ./bin/blueprince-tools ./cmd/tools/
```

### Running the Server

#### Local Development
Make sure you are in the project root directory.

```bash
go run ./cmd/server/main.go
```

The server will start and listen for MCP connections via stdio transport.

#### Environment Configuration
You can override the vault path using an environment variable:

```bash
OBSIDIAN_VAULT_PATH=/path/to/vault go run ./cmd/server/main.go
```

#### Claude Desktop Integration
See the [Claude Desktop instructions for adding custom MCP servers](https://modelcontextprotocol.io/quickstart/server#testing-your-server-with-claude-for-desktop)

Tl;dr 
Add this to your Claude Desktop config:
```json
{
  "mcpServers": {
    "blueprince-notes": {
      "command": "/path/to/blueprince-mcp/bin/blueprince-mcp-server",
      "env": {
        "OBSIDIAN_VAULT_PATH": "/path/to/your/vault"
      }
    }
  }
}
```
**make sure to update your /path/to's**

### Testing with CLI Tools

The project includes a comprehensive CLI for manual testing:

```bash
# List all notes
./bin/blueprince-tools list

# Create a new note
./bin/blueprince-tools create people/character.md \
  --title "Character Name" \
  --content "Character description"

# Read a note
./bin/blueprince-tools read people/character.md

# Update a note
./bin/blueprince-tools update people/character.md \
  --content "Updated character information"

# Use verbose mode for debugging
./bin/blueprince-tools list --verbose
```

See [`cmd/tools/README.md`](cmd/tools/README.md) for detailed CLI documentation and examples.

## Project Structure

```
blueprince-mcp/
├── runtime/
│   ├── mcp/
│   │   ├── tools/              # MCP tool implementations
│   │   │   ├── list.go         # ✅ List notes tool
│   │   │   ├── create.go       # ✅ Create note tool
│   │   │   ├── read.go         # ✅ Read note tool
│   │   │   ├── update.go       # ✅ Update note tool
│   │   │   ├── delete.go       # ✅ Delete note tool
│   │   │   └── register.go     # Tool registration
│   │   └── resources/          # MCP resource system
│   ├── models/
│   │   ├── notes/              # Note structure and schemas
│   │   └── vault/              # Vault constants and structure
│   └── utils/                  # Common utilities (logging, file ops)
├── cmd/
│   ├── server/main.go          # Main MCP server application
│   ├── setup/main.go           # Vault initialization utility
│   ├── tools/                  # CLI testing tools
│   │   ├── main.go             # CLI root command
│   │   ├── client.go           # MCP client implementation
│   │   ├── list.go             # List command
│   │   ├── read.go             # Read command  
│   │   ├── create.go           # Create command
│   │   ├── update.go           # Update command
│   │   └── README.md           # CLI documentation
│   └── config/                 # Configuration management
└── bin/                        # Built binaries
```

## Current Status & Roadmap

### ✅ Completed
- MCP server framework with stdio transport
- Resource system exposing all vault files to AI clients  
- Structured note schema with metadata and categories
- Complete CRUD operations: `list_notes`, `create_note`, `read_note`, `update_note`, `delete_note`
- Vault directory structure and setup utility
- Comprehensive logging and error handling
- Multi-layered spoiler prevention system:
  - Server-side content validation and spoiler detection
  - Spoiler prevention rules exposed as MCP resource
  - Client-side enforcement through tool descriptions
  - Automatic rule delivery to MCP clients
- Path security and traversal prevention
- .obsidian directory filtering for clean vault management
- Utility function abstraction to eliminate code duplication
- Complete CLI testing tools with Cobra framework
- Subprocess communication for reliable MCP testing
- Robust configuration system with environment variable support

### 📋 Planned
- Integration with screenshots
  - Intelligently interpret screenshots to create notes with tags and descriptions of images
  - Embed notes with smart links to related images
  - Serve images back to MCP Client
  - Integrate with Google Drive to sync screenshots from Steam Deck -> Drive -> local

## Contributing

Contributions are welcome! Please feel free to open issues or submit pull requests.

## License

This project is licensed under the MIT License - see the LICENSE file for details.