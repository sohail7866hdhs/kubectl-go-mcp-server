# Copilot Instructions for kubectl-go-mcp-server

<!-- Use this file to provide workspace-specific custom instructions to Copilot. For more details, visit https://code.visualstudio.com/docs/copilot/copilot-customization#_use-a-githubcopilotinstructionsmd-file -->

This is a Go project that implements a Model Context Protocol (MCP) server for Kubernetes operations using kubectl commands.

## Project Context

You can find more info and examples at https://modelcontextprotocol.io/llms-full.txt

## Key Guidelines

1. **Go Best Practices**: Follow standard Go conventions and idioms
2. **Error Handling**: Use proper error wrapping and meaningful error messages
3. **Testing**: Write comprehensive tests for all new functionality
4. **Security**: Be cautious with kubectl operations and validate all inputs
5. **Documentation**: Keep documentation up-to-date and provide clear examples

## Project Structure

- `cmd/`: Main application entry point
- `internal/mcp/`: MCP server implementation following the MCP specification
- `pkg/kubectl/`: kubectl command wrapper with validation and concurrency control
- `pkg/types/`: shared types and data structures
- `test/`: Integration tests and test utilities

## Code Style

- Use `make fmt` before committing
- Run `make lint` to check for issues
- Ensure all tests pass with `make test`
- Follow the existing patterns for error handling and logging

## Security Considerations

- All kubectl commands are validated before execution
- Destructive operations are disabled by default
- Input sanitization is required for all user-provided data
- Proper timeout and concurrency controls are in place

## MCP Protocol

This server implements the Model Context Protocol specification and provides:

- Tool execution capabilities for kubectl commands
- Resource listing for cluster information
- Proper JSON-RPC 2.0 communication
- Error handling and validation
