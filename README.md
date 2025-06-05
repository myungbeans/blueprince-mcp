# Blue Prince Architect Notes (MCP Server)

This repository contains the code for an MCP (Multi-Agent Communication Protocol) server designed to act as a dedicated note-taking and brainstorming assistant specifically tailored for playing the game **Blue Prince**.

By integrating with an Obsidian vault, this server allows an AI agent (or any MCP client) to interact with your game notes, helping you track characters, locations, items, lore, puzzles, and general thoughts as you explore the mysterious world of Blue Prince.

## Features

- **MCP Server:** Implements the MCP protocol to expose note-taking capabilities as tools.
- **Obsidian Vault Integration:** Stores notes directly within a local Obsidian vault directory structure.
- **Structured Notes:** Enforces a predefined directory structure within the vault (`/people`, `/puzzles`, `/rooms`, `/items`, `/lore`, `/general`) to help organize information and facilitate targeted AI interactions.
- **Basic Note Management Tools:** Provides initial tools for listing files (representing notes) within the vault. (Future tools will include Create, Read, Update, and Delete operations).
- **Setup Utility:** A convenient Go program (`cmd/setup`) to initialize the Obsidian vault directory and configure the server's `config.yaml`.
- **Configuration:** Uses a `config.yaml` file for server settings and vault location.
- **Structured Logging:** Utilizes `go.uber.org/zap` for improved logging.

## Getting Started

### Prerequisites

- Go (version 1.20 or higher recommended)
- An MCP client (e.g., a compatible AI agent or a testing tool)
- Git

### Setup

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/your_username/blueprince-mcp.git # Replace with your repo URL
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

    The setup utility will ensure the required subdirectories (`people`, `puzzles`, `rooms`, `items`, `lore`, `general`) exist within the vault.

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

### Running the Server

#### Local
Make sure you are in the project root directory.

```bash
go run ./cmd/server/main.go
```

The server will start and listen for MCP connections on the configured host and port (defaulting to `localhost:8001`).

#### Claude Desktop
See https://modelcontextprotocol.io/quickstart/server#testing-your-server-with-claude-for-desktop

## Project Structure

- `cmd/server`: Contains the main MCP server application.
- `cmd/setup`: Contains the utility program for initial setup.
- `cmd/server/config`: Configuration loading logic.
- `runtime/handlers`: Placeholder and initial implementations for tool handlers (e.g., `list.go`).
- `runtime/utils`: Common utility functions (e.g., file system operations in `files.go`).

## Future Development

- Implement full CRUD (Create, Read, Update, Delete) operations for notes.
- Implement a search tool for notes.
- Add more sophisticated note parsing and structuring capabilities.
- Explore additional tools for brainstorming and game assistance.

## Contributing

Contributions are welcome! Please feel free to open issues or submit pull requests.

## License

This project is licensed under the MIT License - see the LICENSE file for details.