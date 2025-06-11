# Blue Prince MCP - Architect Notes

![Blue Prince MCP Architect Notes logo](static/blue_prince_mcp_logo.png)

This repository contains the code for an MCP (Model Context Protocol) server designed to act as a dedicated note-taking and brainstorming assistant for playing the game [Blue Prince](https://store.steampowered.com/app/1569580/Blue_Prince/).

This MCP server exposes tools and resources for managing local notes (stored as .md files) that allow users to write notes, lookup information from their notes, and brainstorm with a companion MCP client as they play through the video game Blue Prince. This is designed to help players make connections and recall things they've seen and experienced while avoiding spoilers from online resources.

**⚠️ IMPORTANT: SPOILER-AWARE USAGE**

This MCP server is designed to preserve your Blue Prince gameplay experience. When used with an MCP client (e.g. Claude Desktop):
- **Primary Source**: The Client prioritizes information from your notes
- **Filtered External Access**: The Client MAY use external Blue Prince knowledge, but ONLY for content you've already discovered
- **Spoiler Protection**: External information is filtered to show only what you've documented in your notes
- **Consent Required**: The Client will ask permission before sharing potentially spoiling external information
- **Smart Filtering**: Automatic spoiler prevention rules are provided as an MCP resource
- **Discovery Preservation**: Focus remains on your documented experiences and discoveries

## Features

- **MCP Server:** Implements the MCP protocol to expose note-taking capabilities as tools and resources.
- **Local Vault Storage:** Stores notes as markdown files in a structured local directory (compatible with Obsidian).
- **Structured Notes:** Organizes notes in predefined categories (`people`, `puzzles`, `rooms`, `items`, `lore`, `general`) with intelligent metadata extraction.
- **Resource System:** Exposes all vault files as MCP resources for direct access by AI clients (excludes `.obsidian/` directories).
- **Spoiler-Aware Protection System:** Smart filtering that preserves discovery while enabling helpful context:
  - Dynamic spoiler prevention rules automatically exposed as an MCP resource
  - Client-side enforcement through tool descriptions and server metadata
  - Server-side validation of all content creation with discovery preservation
  - Automatic filtering of external information based on user's documented discoveries
  - Consent-based sharing of potentially spoiling external information
  - Built-in content validation to prevent premature investigation prompts
- **Complete CRUD Operations:** 
  - ✅ `list_notes` - Lists all notes in the vault
  - ✅ `create_note` - Creates structured notes with intelligent categorization and spoiler prevention
  - ✅ `read_note` - Reads complete note content including metadata
  - ✅ `update_note` - Updates existing notes with new content
  - 📋 `delete_note` - Planned for future implementation
- **CLI Testing Tools:** Comprehensive command-line interface for manual testing and debugging.
- **Setup Utility:** Go program to initialize vault directory structure and configuration, as well as OAuth with Google Drive for screenshot syncs.
- **Flexible Configuration:** Supports both file-based config and environment variable overrides.

## Usage Guide
```
"Write a new note. I'm in the corridor. There is a painting of a tiger and a cupcake stand(?). Three windows. Two benches and hats."
```

The MCP Client will then intelligently format the note and tag it appropraitely.

In the future you can ask
```
"Where have I seen windows before? Can you list all rooms that have windows in them?"
```

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
    bin/setup
    ```

    Alternatively, you can specify a custom path for your vault:

    ```bash
    bin/setup /path/to/your/custom/vault
    ```

    The setup utility will ensure the required subdirectories (`notes/people`, `notes/puzzles`, `notes/rooms`, `notes/items`, `notes/lore`, `notes/general`) exist within the vault, along with `meta/` and `screenshots/` directories.

3.  **Configure Google Drive Integration (Optional):**
    To enable automatic screenshot syncing from Google Drive, run the drive setup command:

    ```bash
    bin/setup drive "YourFolderName"
    ```

    This will:
    - Open your browser for Google Drive OAuth authentication
    - Create or access the specified folder in your Google Drive
    - Configure permissions for viewing, editing, creating, and downloading files
    - Save authentication tokens locally for future use

    **Requirements:**
    - The folder name must be specified (cannot be root directory)
    - Google Drive permissions include: view, list, edit, create directories, download files
    - All authentication data is stored locally on your machine. See our [Privacy Policy](privacy-policy.html) for more details on data handling

4.  **Review `config.yaml`:**
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

### Building the Setup Utility
```bash
go build -o ./bin/setup ./cmd/setup/main.go
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
        "GOOGLE_DRIVE_SCREENSHOT_FOLDER": "Blue Prince",
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
│   │   ├── tools/              # MCP Tool implementations
│   │   └── resources/          # MCP Resources
│   ├── models/
│   │   ├── notes/              # Note structure and schemas
│   │   ├── vault/              # Obsidian Vault constants and structure
│   │   └── storage/            # Storage interface abstractions
│   ├── storage/                # Storage implementations
│   │   └── drive/              # Google Drive implementation
│   └── utils/                  # Common utilities (logging, file ops, security)
├── cmd/
│   ├── server/main.go          # Main MCP server application
│   ├── setup/                  # Setup utilities
│   ├── tools/                  # CLI for running MCP Server Tools locally
│   └── config/                 # Configuration management
├── docs/                       # Documentation for GitHub Pages
└── bin/                        # Built binaries
```

## Current Status & Roadmap

### ✅ Completed
- **Core MCP Framework:**
  - MCP server framework with stdio transport
  - Resource system exposing all vault files to AI clients  
  - Structured note schema with metadata and categories
  - Complete CRUD operations: `list_notes`, `create_note`, `read_note`, `update_note`, `delete_note`
  - Vault directory structure and setup utility
- **Google Drive Integration:**
  - OAuth2 authentication flow with automatic browser opening
  - Full Google Drive API permissions (view, list, edit, create, download)
  - Secure local token storage in `~/.blueprince_mcp/`
  - Automated folder creation and access verification
  - Privacy-focused design with no third-party data transmission
  - **Refactored Architecture:**
    - Modular storage interface with `runtime/storage/drive` backend
    - Centralized Google Drive utilities and path management
    - Shared credential loading and token management
    - Clean separation between setup and runtime operations
- **Code Quality & Testing:**
  - **Comprehensive unit test coverage** across all major components
  - Storage utilities testing (path management, token handling, configs)
  - File operations testing (security validation, directory management)
  - Authentication flow testing (OAuth setup, error handling)
  - Mock implementations for external dependencies
  - Benchmark tests for performance validation
- **Security & Reliability:**
  - Multi-layered spoiler prevention system
  - Path security and traversal prevention
  - Input validation and error handling
  - Configuration management with environment variable support

### 📋 Planned
- **Enhanced Screenshot Integration:**
  - Intelligently interpret screenshots to create notes with tags and descriptions of images
  - Embed notes with smart links to related images
  - Serve images back to MCP Client
  - ✅ **Google Drive sync foundation complete** - ready for screenshot synchronization from Steam Deck -> Google Drive -> local vault
  - Automatic download and sync of files from configured Google Drive folder

## Contributing

Contributions are welcome! Please feel free to open issues or submit pull requests.

## License

This project is licensed under the MIT License - see the LICENSE file for details.