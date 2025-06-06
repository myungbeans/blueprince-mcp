# Blue Prince Architect Notes (MCP Server)

This repository contains the code for an MCP (Multi-Agent Communication Protocol) server designed to act as a dedicated note-taking and brainstorming assistant specifically tailored for playing the game **Blue Prince**.

This MCP server exposes tools and resources for managing local notes (stored as .md files) that allow users to write notes, lookup information from their notes, and brainstorm with a companion MCP client as they play through the video game Blue Prince. This is designed to help players make connections and recall things they've seen and experienced while avoiding spoilers from online resources.

## Features

- **MCP Server:** Implements the MCP protocol to expose note-taking capabilities as tools and resources.
- **Local Vault Storage:** Stores notes as markdown files in a structured local directory (compatible with Obsidian).
- **Structured Notes:** Organizes notes in predefined categories (`people`, `puzzles`, `rooms`, `items`, `lore`, `general`) with intelligent metadata extraction.
- **Resource System:** Exposes all vault files as MCP resources for direct access by AI clients.
- **Smart Note Creation:** AI-powered note creation with structured metadata and content templates designed to avoid spoilers.
- **Note Management Tools:** 
  - ✅ `list_notes` - Lists all notes in the vault
  - ✅ `create_note` - Creates structured notes with intelligent categorization (schema defined, handler in progress)
  - 🚧 Read, Update, Delete operations (files exist, implementation pending)
- **Setup Utility:** Go program to initialize vault directory structure and configuration.
- **Flexible Configuration:** Supports configs for local server and plugin to Claude Desktop.
- **Structured Logging:** Uses `go.uber.org/zap` for comprehensive logging.

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

## Build
run `go build -o ./bin/blueprince-mcp-server ./cmd/server/main.go`
(TODO: use Earthly to handle build and test cmds?) 

### Running the Server

#### Local Development
Make sure you are in the project root directory.

```bash
go run ./cmd/server/main.go
```

The server will start and listen for MCP connections via stdio transport.

#### Claude Desktop
See https://modelcontextprotocol.io/quickstart/server#testing-your-server-with-claude-for-desktop

## Project Structure

```
blueprince-mcp/
├── cmd/
│   ├── server/main.go          # Main MCP server application
│   ├── setup/main.go           # Vault initialization utility
│   └── config/                 # Configuration management
├── runtime/
│   ├── mcp/
│   │   ├── tools/              # MCP tool implementations
│   │   │   ├── list.go         # ✅ List notes tool
│   │   │   ├── create.go       # 🚧 Create note tool (in progress)
│   │   │   ├── read.go         # 📋 Read note tool (planned)
│   │   │   ├── update.go       # 📋 Update note tool (planned)
│   │   │   └── delete.go       # 📋 Delete note tool (planned)
│   │   └── resources/          # MCP resource system
│   ├── models/
│   │   ├── notes/              # Note structure and schemas
│   │   └── vault/              # Vault constants and structure
│   └── utils/                  # Common utilities (logging, file ops)
└── bin/                        # Built binaries
```

## Current Status & Roadmap

### ✅ Completed
- MCP server framework with stdio transport
- Resource system exposing all vault files to AI clients
- Structured note schema with metadata and categories
- `list_notes` tool implementation
- Vault directory structure and setup utility
- Comprehensive logging and error handling

### 🚧 In Progress
- `create_note` tool implementation (schema complete, handler in development)

### 📋 Planned
- Complete CRUD operations: `read_note`, `update_note`, `delete_note`
- Search and query tools for note discovery
- Note relationship and connection mapping
- Enhanced metadata extraction and categorization
- Template system for different note types
- Integration with game progress tracking

## Contributing

Contributions are welcome! Please feel free to open issues or submit pull requests.

## License

This project is licensed under the MIT License - see the LICENSE file for details.