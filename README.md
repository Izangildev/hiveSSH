# HiveSSH

HiveSSH is a powerful command-line tool written in Go that allows you to manage multiple SSH connections and execute remote commands across your server infrastructure. It provides a simple interface to store SSH targets and execute commands remotely.

## 📋 Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Commands](#commands)
- [Project Structure](#project-structure)
- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Contributing](#contributing)
- [License](#license)

## ✨ Features

- **Server Management**: Add, list, and manage SSH targets
- **Remote Command Execution**: Execute commands on remote servers
- **Connection Status**: Check SSH connectivity status for all servers
- **JSON Storage**: Store server configurations in a simple JSON format
- **SSH Key Authentication**: Secure authentication using SSH private keys
- **Cross-Platform**: Built with Go for cross-platform compatibility

## 🚀 Installation

### Prerequisites

- Go 1.24.4 or later
- SSH private key configured (default: `~/.ssh/id_rsa`)

### Build from Source

```bash
git clone https://github.com/Izangildev/hiveSSH.git
cd hiveSSH
go build
```

## 🏃 Quick Start

1. **Add a server to your inventory:**
   ```bash
   ./hivessh join myserver 192.168.1.100
   ```

2. **List all configured servers:**
   ```bash
   ./hivessh list
   ```

3. **Execute a command on a remote server:**
   ```bash
   ./hivessh run "uptime" --to myserver
   ```

## 📚 Commands

### `hivessh join <name> <ip>`

Adds a new server to the database.

**Arguments:**
- `name`: Unique identifier for the server
- `ip`: IP address of the server

**Example:**
```bash
./hivessh join webserver 192.168.1.50
```

**Features:**
- Validates IP address format
- Prevents duplicate server names or IPs
- Automatically saves to servers.json

### `hivessh list`

Displays all configured servers with their connection status.

**Output Format:**
```
NAME       IP                 SSH STATUS    
────────── ────────────────── ──────────────
webserver  192.168.1.50       reachable     
────────── ────────────────── ──────────────
dbserver   192.168.1.51       unreachable   
────────── ────────────────── ──────────────
```

**Features:**
- Real-time SSH connectivity testing (port 22)
- Formatted table output
- Status indicators: `reachable` or `unreachable`

### `hivessh run <command> --to <target>`

Executes a command on a remote server.

**Arguments:**
- `command`: The command to execute remotely
- `--to`: Target server (name or IP address)

**Example:**
```bash
./hivessh run "df -h" --to webserver
./hivessh run "ls -la" --to 192.168.1.50
```

**Features:**
- Supports both server names and IP addresses as targets
- Displays command output
- Error handling for connection failures
- Uses SSH key authentication

## 🏗️ Project Structure

```
hiveSSH/
├── cmd/                    # Command-line interface definitions
│   ├── root.go            # Root command configuration
│   ├── join.go            # Join command implementation
│   ├── list.go            # List command implementation
│   └── run.go             # Run command implementation
├── logic/                  # Core business logic
│   ├── servers.go         # Server management functions
│   ├── join.go            # Join operation logic
│   ├── list.go            # List operation logic
│   └── run.go             # Remote execution logic
├── env/                    # Environment configuration
│   └── env.go             # Configuration constants
├── main.go                # Application entry point
├── servers.json           # Server inventory storage
├── go.mod                 # Go module definition
├── go.sum                 # Go module checksums
└── README.md             # This documentation
```

## 📁 Detailed File Documentation

### Core Files

#### `main.go`
**Purpose**: Application entry point that initializes the server database and starts the CLI.

**Functions:**
- `main()`: Loads server configuration and executes CLI commands

**Dependencies:** 
- `hivessh/cmd` (CLI framework)
- `hivessh/env` (configuration)
- `hivessh/logic` (core functionality)

#### `servers.json`
**Purpose**: JSON file storing server inventory with name-to-IP mappings.

**Format:**
```json
{
  "server1": "192.168.1.100",
  "server2": "192.168.1.101"
}
```

### Command Layer (`cmd/`)

#### `cmd/root.go`
**Purpose**: Defines the root command and global CLI configuration.

**Key Components:**
- `rootCmd`: Main cobra command configuration
- `Execute()`: Starts the CLI application
- `init()`: Initializes command flags

**Functions:**
- `Execute()`: Entry point for command execution with error handling

#### `cmd/join.go`
**Purpose**: Implements the server registration command.

**Key Components:**
- `joinCmd`: Cobra command for adding servers
- Input validation for server names and IP addresses
- Integration with join logic

**Functions:**
- Command handler: Validates inputs and calls `logic.Join()`
- IP validation using `net.ParseIP()`
- Error handling and user feedback

#### `cmd/list.go`
**Purpose**: Implements the server listing command.

**Key Components:**
- `listCmd`: Cobra command for displaying servers
- Simple interface to list logic

**Functions:**
- Command handler: Calls `logic.List()` to display server inventory

#### `cmd/run.go`
**Purpose**: Implements the remote command execution.

**Key Components:**
- `runCmd`: Cobra command for remote execution
- `target` flag for specifying destination server
- Command validation

**Functions:**
- Command handler: Validates command and target, calls `logic.Run()`
- Flag definition for `--to` parameter

### Logic Layer (`logic/`)

#### `logic/servers.go`
**Purpose**: Core server management functionality and data persistence.

**Global Variables:**
- `servers`: In-memory map storing server inventory

**Functions:**

- `getStatus(ip string) bool`
  - Tests SSH connectivity to a server
  - Uses TCP connection to port 22 with 5-second timeout
  - Returns true if connection successful, false otherwise

- `serverExists(identifier string) (bool, string)`
  - Checks if a server exists by name or IP
  - Returns existence status and identifier type ("name" or "IP")
  - Used for duplicate detection and target resolution

- `existServersFile(serversFile string) bool`
  - Checks if the servers.json file exists
  - Handles file system errors gracefully
  - Returns false if file doesn't exist or access fails

- `SaveServers()`
  - Persists the in-memory server map to servers.json
  - Uses JSON marshaling with indentation for readability
  - Handles file write errors with user feedback

- `LoadServers(serversFile string)`
  - Loads server inventory from servers.json into memory
  - Handles missing files gracefully (empty initialization)
  - Includes JSON parsing error handling
  - Called during application startup

#### `logic/join.go`
**Purpose**: Server registration logic.

**Functions:**

- `Join(serverName, ip string) error`
  - Adds a new server to the inventory
  - Validates uniqueness of both name and IP
  - Updates in-memory map and persists to file
  - Returns error if server or IP already exists

#### `logic/list.go`
**Purpose**: Server listing and status display logic.

**Types:**
- `extractableServer`: Internal struct for server data representation
  - `name string`: Server identifier
  - `ip string`: Server IP address  
  - `status bool`: Connectivity status

**Functions:**

- `List()`
  - Displays formatted table of all servers
  - Tests connectivity for each server using `getStatus()`
  - Formats output with fixed-width columns
  - Shows status as "reachable" or "unreachable"

#### `logic/run.go`
**Purpose**: Remote command execution via SSH.

**Functions:**

- `Run(command, identifier string) error`
  - Executes commands on remote servers
  - Resolves server identifier (name or IP) to actual IP
  - Establishes SSH connection using private key authentication
  - Captures and displays command output (stdout/stderr)
  - Comprehensive error handling for all failure modes

**SSH Implementation Details:**
- Uses `github.com/melbahja/goph` for SSH operations
- Authenticates with private key from `env.Private_key`
- Connects as "root" user (configurable in code)
- Creates new SSH session for each command
- Handles both stdout and stderr streams

### Environment Layer (`env/`)

#### `env/env.go`
**Purpose**: Configuration constants and environment settings.

**Variables:**
- `ServersFile`: Path to server inventory file (default: "servers.json")
- `Private_key`: Path to SSH private key (default: "~/.ssh/id_rsa")

**Features:**
- Uses `os.Getenv("HOME")` for cross-platform home directory resolution
- Configurable paths for different deployment scenarios

## ⚙️ Configuration

### Server Storage
Servers are stored in `servers.json` in the application directory:

```json
{
  "production-web": "10.0.1.100",
  "staging-db": "10.0.2.50",
  "development": "192.168.1.200"
}
```

### SSH Configuration
- **Default private key location**: `~/.ssh/id_rsa`
- **Default user**: `root` (configurable in source)
- **Connection timeout**: 5 seconds for status checks
- **SSH port**: 22 (standard)

### Customization
To modify default paths, edit `env/env.go`:

```go
var (
    ServersFile = "/path/to/custom/servers.json"
    Private_key = "/path/to/custom/key"
)
```

## 📦 Dependencies

HiveSSH uses the following external libraries:

### Core Dependencies
- **`github.com/spf13/cobra v1.9.1`**: CLI framework for command structure
- **`github.com/melbahja/goph v1.4.0`**: SSH client implementation
- **`github.com/pkg/sftp v1.13.5`**: SFTP support (via goph)
- **`golang.org/x/crypto v0.6.0`**: Cryptographic operations

### Utility Dependencies
- **`github.com/pkg/errors v0.9.1`**: Enhanced error handling
- **`github.com/google/uuid v1.2.0`**: UUID generation (indirect)
- **`github.com/go-ping/ping v1.2.0`**: Network ping functionality (indirect)

### Standard Library Usage
- **`encoding/json`**: Server data persistence
- **`net`**: IP address validation and network operations
- **`os`**: File system operations
- **`time`**: Connection timeouts
- **`fmt`**: Output formatting
- **`bytes`**: SSH output handling

## 🔧 Advanced Usage

### Error Handling
HiveSSH provides detailed error messages for common scenarios:

- **Invalid IP addresses**: Validates IPv4/IPv6 format
- **Duplicate servers**: Prevents name/IP conflicts
- **Connection failures**: Reports SSH connectivity issues
- **File permissions**: Handles configuration file access errors
- **Command failures**: Displays stderr output for failed commands

### Security Considerations
- Uses SSH key authentication (no password storage)
- Requires proper SSH key permissions (600)
- Connects to standard SSH port (22)
- No credential storage in configuration files

### Performance
- Connection status checks use 5-second timeouts
- Concurrent connectivity testing for list command
- JSON file I/O for persistent storage
- In-memory caching of server inventory

## 🚀 Future Enhancements

Potential improvements for HiveSSH:

1. **Configuration Management**
   - YAML/TOML configuration support
   - Per-server SSH key configuration
   - Custom SSH ports and users

2. **Advanced Features**
   - Server groups and bulk operations
   - Command history and favorites
   - SSH agent integration
   - Key rotation utilities

3. **User Experience**
   - Interactive server selection
   - Auto-completion support
   - Colored output and themes
   - Progress indicators for long operations

4. **Enterprise Features**
   - LDAP/AD integration
   - Audit logging
   - Role-based access control
   - API interface

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Setup
```bash
git clone https://github.com/Izangildev/hiveSSH.git
cd hiveSSH
go mod download
go build
./hivessh --help
```

### Code Style
- Follow Go conventions and formatting (`go fmt`)
- Add comments for exported functions
- Include error handling for all operations
- Write unit tests for new functionality

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 📞 Support

For questions, issues, or contributions:
- Open an issue on GitHub
- Check existing documentation
- Review the code examples in this README

---

**Happy SSH managing with HiveSSH! 🐝**