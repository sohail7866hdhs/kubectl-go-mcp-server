# kubectl-go-mcp-server

A Model Context Protocol (MCP) server that provides Kubernetes cluster interaction capabilities through kubectl commands. This server enables MCP-compatible clients (like VS Code with Copilot) to execute kubectl commands and retrieve Kubernetes cluster information safely and securely.

## Features

- **Kubernetes Integration**: Execute kubectl commands through MCP interface
- **Interactive Command Protection**: Prevents execution of interactive commands that could hang
- **Resource Modification Detection**: Identifies commands that modify cluster resources
- **Robust Security**: Multiple validation layers to prevent command injection and unsafe operations
- **Configurable Kubeconfig**: Support for custom kubeconfig paths
- **Standard Go Project Layout**: Following Go best practices for maintainability
- **Cobra CLI Integration**: Professional command-line interface with subcommands

## Architecture

kubectl-go-mcp-server acts as a bridge between MCP clients (like VS Code with Copilot) and Kubernetes clusters through kubectl commands:

```
VS Code/Copilot → MCP Client → kubectl-go-mcp-server → kubectl → Kubernetes Cluster
```

### Key Components

- **MCP Server**: Handles JSON-RPC communication and tool registration
- **kubectl Tool**: Validates and executes kubectl commands safely
- **Security Layer**: Prevents interactive commands and command injection

For detailed architecture information, see [docs/architecture.md](docs/architecture.md).

```
pkg/
├── types/          # 🔧 Core interfaces and data structures
│   ├── Tool        # Interface for all MCP tools
│   ├── Schema      # JSON schema definitions
│   └── ExecResult  # Command execution results
│
├── kubectl/        # 🎯 kubectl-specific implementation
│   ├── KubectlTool # Main tool implementation
│   ├── Validation  # Command safety checks
│   └── Execution   # kubectl command runner
│
internal/
├── mcp/           # 🌐 MCP protocol implementation
│   ├── Server     # MCP server and protocol handling
│   ├── Tools      # Tool registry and management
│   └── Protocol   # JSON-RPC message handling
│
└── config/        # ⚙️ Configuration management
    ├── Config     # Application configuration
    └── Defaults   # Default settings
```

### Extension Points

The architecture is designed for extensibility:

1. **New Tools**: Implement the `Tool` interface to add new capabilities
2. **Custom Validation**: Add validation layers for specific use cases
3. **Protocol Extensions**: Extend MCP handling for additional features
4. **Output Formatters**: Add custom result processing

### Performance Considerations

- **Concurrent Safety**: All components are designed for concurrent access
- **Resource Management**: Proper cleanup and resource disposal
- **Timeout Handling**: Configurable timeouts for all operations
- **Memory Efficiency**: Streaming and buffered I/O for large outputs

## Installation

### Prerequisites

- Go 1.23 or later
- kubectl installed and configured
- Access to a Kubernetes cluster

### Build from Source

```bash
# Clone the repository
git clone <repository-url>
cd kubectl-go-mcp-server

# Build the binary
make build

# Or install directly
make install
```

### Download Binary

Download the latest release from the [releases page](releases) for your platform.

## Usage

### Standalone

```bash
# Run with default kubeconfig
./kubectl-go-mcp-server

# Run with custom kubeconfig
./kubectl-go-mcp-server --kubeconfig /path/to/kubeconfig

# Show version
./kubectl-go-mcp-server version
```

### VS Code Integration

For comprehensive VS Code MCP integration guides, see the [examples directory](./examples/):

- **[Docker Integration](./examples/docker-integration.md)**: Complete Docker-based configurations for all platforms
- **[macOS Native](./examples/macos-native.md)**: Native macOS binary installation and configuration
- **[Windows Native](./examples/windows-native.md)**: Native Windows binary installation and configuration
- **[Linux Native](./examples/linux-native.md)**: Native Linux binary installation and configuration

#### Quick Configuration Example

Add the server to your VS Code MCP configuration in `settings.json`:

```json
{
  "mcp": {
    "servers": {
      "kubectl-go-mcp-server": {
        "type": "stdio",
        "command": "/path/to/kubectl-go-mcp-server",
        "args": [],
        "env": {}
      }
    }
  }
}
```

**Note**: The exact configuration varies by platform and installation method. See the platform-specific guides in the `examples/` directory for detailed instructions.

## Available Tools

The MCP server provides the following tool:

### kubectl

Execute kubectl commands with comprehensive validation and safety checks.

**Parameters:**

- `command` (required): The complete kubectl command to execute (including 'kubectl' prefix)
- `modifies_resource` (optional): Indicates if the command modifies resources ("yes", "no", "unknown")

**Example:**

```json
{
  "name": "kubectl",
  "arguments": {
    "command": "kubectl get pods -o json",
    "modifies_resource": "no"
  }
}
```

**Safety Features:**

- **Interactive Command Detection**: Prevents hanging on interactive commands like `kubectl exec -it`, `kubectl edit`, `kubectl port-forward`
- **Resource Modification Tracking**: Automatically detects destructive operations
- **Command Validation**: Ensures only valid kubectl commands are executed

## Security

This server implements multiple security layers including command validation, injection prevention, and interactive command blocking. For detailed security information, see:

- [Security Overview](docs/security.md) - Technical security implementation details
- [Security Policy](SECURITY.md) - Vulnerability reporting and security best practices

## Development

See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup, workflow, and contribution guidelines.

### Quick Start for Developers

```bash
# Install dependencies and build
make deps && make build

# Run tests
make test

# Format and lint code  
make fmt && make lint
```

## Troubleshooting

### Common Issues

- **Server not responding**: Verify kubectl is installed and kubeconfig is accessible
- **Interactive command errors**: Use non-interactive alternatives (see [Security Overview](docs/security.md))
- **Permission denied**: Check kubectl permissions and cluster connectivity

For detailed debugging information, the server logs all tool calls, validation results, and errors.

## Support

- [Create an issue](https://github.com/your-username/kubectl-go-mcp-server/issues) for bug reports or feature requests
- Check [existing issues](https://github.com/your-username/kubectl-go-mcp-server/issues) for known problems
- See [CONTRIBUTING.md](CONTRIBUTING.md) for development questions

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
