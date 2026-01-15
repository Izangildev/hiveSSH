# HiveSSH - SSH Server Management Tool

A Go-based CLI tool for managing multiple SSH servers and executing commands across server groups.

## Features

- **Server Management**: Add, list, and manage SSH servers with connection details
- **Group Operations**: Create server groups for batch command execution
- **Concurrent Execution**: Run commands on multiple servers simultaneously
- **Configuration Persistence**: Stores server and group data in JSON files

## Quick Start

```bash
# Add a server
hivessh join server-name 192.168.1.10 username 22 "description"

# Create a server group
hivessh group create production

# Add server to group
hivessh group join production server-name

# Execute command on group
hivessh run production "uptime"

# List all servers
hivessh list
```

## Architecture

The application follows a three-tier architecture [1](#0-0) :

1. **CLI Layer**: Cobra-based command interface in `cmd/` directory
2. **Business Logic Layer**: Core operations in `logic/` directory  
3. **Data Persistence Layer**: JSON file storage for configuration

## Configuration Files

- `servers.json`: Stores server definitions with connection details [2](#0-1) 
- `groups.json`: Stores server groups and memberships [3](#0-2) 

## Project Structure

```
hivessh/
├── main.go              # Application entry point
├── cmd/                 # CLI command definitions
│   ├── root.go          # Root command setup
│   ├── join.go          # Server registration
│   ├── run.go           # Command execution
│   ├── list.go          # Server listing
│   └── group/           # Group management commands
├── logic/               # Business logic
│   ├── servers.go       # Server operations
│   ├── groups.go        # Group operations
│   ├── run.go           # SSH execution
│   ├── join.go          # Server registration logic
│   └── list.go          # Listing operations
└── env/                 # Environment configuration
```

## Requirements

- Go 1.19+
- SSH access to target servers
- SSH key authentication

## Installation

```bash
git clone https://github.com/Izangildev/hiveSSH
cd hiveSSH
go build -o hivessh
sudo mv hivessh /usr/local/bin/
```

### Citations

**File:** main.go (L10-14)
```go
func main() {
	logic.LoadGroups(env.GroupsFile)
	logic.LoadServers(env.ServersFile)
	cmd.Execute()
}
```

**File:** logic/servers.go (L71-83)
```go
func SaveServers() {
	data, err := json.MarshalIndent(Servers, "", "  ")
	if err != nil {
		fmt.Printf("[❌] Failed to convert in JSON: %s\n", err)
		return
	}

	err = os.WriteFile(env.ServersFile, data, 0644)
	if err != nil {
		fmt.Printf("[❌] Failed to write servers file: %s\n", err)
		return
	}
}
```

**File:** logic/groups.go (L31-42)
```go
func SaveGroups() {
	data, err := json.MarshalIndent(Groups, "", "  ")
	if err != nil {
		fmt.Printf("[❌] Failed to save groups: %s\n", err)
		return
	}

	err = os.WriteFile(env.GroupsFile, data, 0644)
	if err != nil {
		fmt.Printf("[❌] Failed to write groups file: %s\n", err)
	}
}
```

**File:** logic/group/create.go (L34-38)
```go
func createID() string {
	hash := md5.New()
	hash.Write([]byte(fmt.Sprintf("%d", time.Now().UnixNano())))
	return fmt.Sprintf("%x", hash.Sum(nil))
}
```
