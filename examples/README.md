# kubectl-go-mcp-server Examples

This directory contains practical examples for using kubectl-go-mcp-server with various MCP clients.

## VS Code Setup

### 1. Configure VS Code Settings

Add to your VS Code settings (`.vscode/settings.json` or global settings):

```json
{
    "mcp": {        
        "servers": {
            "kubectl-go-mcp-server": {
                "type": "stdio",
                "command": "/path/to/kubectl-go-mcp-server",
                "args": [],
                "env": {
                    "KUBECONFIG": "/path/to/.kube/config"
                }
            }
        }
    }
}
```

### 2. Installation Options

**Option A: Download Binary**
```bash
# Download from releases page
curl -L https://github.com/your-username/kubectl-go-mcp-server/releases/latest/download/kubectl-go-mcp-server-linux-amd64 -o kubectl-go-mcp-server
chmod +x kubectl-go-mcp-server
```

**Option B: Build from Source**
```bash
git clone https://github.com/your-username/kubectl-go-mcp-server.git
cd kubectl-go-mcp-server
make build
# Binary will be in ./bin/kubectl-go-mcp-server
```

**Option C: Docker**
```bash
docker run --rm -v ~/.kube:/root/.kube kubectl-go-mcp-server
```

## Example Usage

Once configured, you can use natural language in VS Code to interact with your Kubernetes cluster:

### Common Commands
- "Show me all pods in the default namespace"
- "Describe the nginx deployment"
- "Get the logs from the api-server pod"
- "What nodes are in my cluster?"
- "Scale the web deployment to 3 replicas"

### Security Features
- ✅ Only kubectl commands are allowed
- ✅ Interactive commands are blocked automatically
- ✅ Command injection prevention
- ✅ Configurable kubeconfig paths

## Troubleshooting

- **Binary not found**: Update the `command` path in VS Code settings
- **Permission denied**: Ensure the binary is executable (`chmod +x`)
- **Kubeconfig issues**: Verify the `KUBECONFIG` path in the environment settings
