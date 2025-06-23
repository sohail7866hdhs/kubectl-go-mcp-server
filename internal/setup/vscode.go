package setup

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

// VSCodeSettings represents the structure of VS Code settings.json
type VSCodeSettings struct {
	MCP *MCPConfig `json:"mcp,omitempty"`
	// Keep other settings as raw JSON to preserve existing configuration
	OtherSettings map[string]interface{} `json:"-"`
}

// MCPConfig represents the MCP configuration section
type MCPConfig struct {
	Servers map[string]*MCPServer `json:"servers"`
}

// MCPServer represents a single MCP server configuration
type MCPServer struct {
	Type    string            `json:"type"`
	Command string            `json:"command"`
	Args    []string          `json:"args"`
	Env     map[string]string `json:"env,omitempty"`
}

// SetupVSCode configures VS Code to use kubectl-go-mcp-server as an MCP server
func SetupVSCode() error {
	fmt.Println("🔍 Setting up kubectl-go-mcp-server for VS Code...")

	// 1. Find VS Code settings path
	settingsPath, err := getVSCodeSettingsPath()
	if err != nil {
		return fmt.Errorf("failed to find VS Code settings: %w", err)
	}
	fmt.Printf("✅ Found VS Code settings at: %s\n", settingsPath)

	// 2. Detect binary location
	binaryPath, err := getBinaryPath()
	if err != nil {
		return fmt.Errorf("failed to detect binary path: %w", err)
	}
	fmt.Printf("✅ Detected binary at: %s\n", binaryPath)

	// 3. Read existing settings
	settings, err := readVSCodeSettings(settingsPath)
	if err != nil {
		return fmt.Errorf("failed to read VS Code settings: %w", err)
	}

	// 4. Create backup
	if err := createBackup(settingsPath); err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}
	fmt.Printf("✅ Created backup: %s\n", settingsPath+".backup")

	// 5. Add MCP configuration
	addMCPConfiguration(settings, binaryPath)

	// 6. Write updated settings
	if err := writeVSCodeSettings(settingsPath, settings); err != nil {
		return fmt.Errorf("failed to write VS Code settings: %w", err)
	}
	fmt.Printf("✅ Added kubectl-go-mcp-server to MCP servers\n")

	fmt.Println("🎉 Setup complete! Restart VS Code to use the new MCP server.")
	return nil
}

// getVSCodeSettingsPath returns the platform-specific VS Code settings path
func getVSCodeSettingsPath() (string, error) {
	var settingsPath string

	switch runtime.GOOS {
	case "linux":
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		settingsPath = filepath.Join(home, ".config", "Code", "User", "settings.json")
	case "darwin":
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		settingsPath = filepath.Join(home, "Library", "Application Support", "Code", "User", "settings.json")
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData == "" {
			return "", fmt.Errorf("APPDATA environment variable not set")
		}
		settingsPath = filepath.Join(appData, "Code", "User", "settings.json")
	default:
		return "", fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	// Create the directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(settingsPath), 0755); err != nil {
		return "", fmt.Errorf("failed to create VS Code settings directory: %w", err)
	}

	return settingsPath, nil
}

// getBinaryPath detects the current binary path
func getBinaryPath() (string, error) {
	// First try to get the path of the current executable
	exe, err := os.Executable()
	if err == nil {
		// Resolve any symlinks
		resolved, err := filepath.EvalSymlinks(exe)
		if err == nil {
			return resolved, nil
		}
		return exe, nil
	}

	// Fallback: try to find kubectl-go-mcp-server in PATH
	path, err := exec.LookPath("kubectl-go-mcp-server")
	if err == nil {
		return path, nil
	}

	return "", fmt.Errorf("could not determine binary path: executable detection failed and kubectl-go-mcp-server not found in PATH")
}

// readVSCodeSettings reads and parses the VS Code settings file
func readVSCodeSettings(settingsPath string) (map[string]interface{}, error) {
	// If file doesn't exist, start with empty settings
	if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
		return make(map[string]interface{}), nil
	}

	data, err := os.ReadFile(settingsPath)
	if err != nil {
		return nil, err
	}

	// Handle empty file
	if len(data) == 0 {
		return make(map[string]interface{}), nil
	}

	var settings map[string]interface{}
	if err := json.Unmarshal(data, &settings); err != nil {
		return nil, fmt.Errorf("failed to parse settings.json (invalid JSON): %w", err)
	}

	return settings, nil
}

// createBackup creates a backup of the settings file
func createBackup(settingsPath string) error {
	// Only create backup if the file exists
	if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
		return nil // No backup needed for non-existent file
	}

	backupPath := settingsPath + ".backup"
	
	// If backup already exists, add timestamp
	if _, err := os.Stat(backupPath); err == nil {
		timestamp := time.Now().Format("20060102-150405")
		backupPath = settingsPath + ".backup." + timestamp
	}

	data, err := os.ReadFile(settingsPath)
	if err != nil {
		return err
	}

	return os.WriteFile(backupPath, data, 0644)
}

// addMCPConfiguration adds the MCP server configuration to settings
func addMCPConfiguration(settings map[string]interface{}, binaryPath string) {
	// Get or create MCP section
	var mcpConfig map[string]interface{}
	if mcp, exists := settings["mcp"]; exists {
		if mcpMap, ok := mcp.(map[string]interface{}); ok {
			mcpConfig = mcpMap
		} else {
			mcpConfig = make(map[string]interface{})
		}
	} else {
		mcpConfig = make(map[string]interface{})
	}

	// Get or create servers section
	var servers map[string]interface{}
	if srvs, exists := mcpConfig["servers"]; exists {
		if srvMap, ok := srvs.(map[string]interface{}); ok {
			servers = srvMap
		} else {
			servers = make(map[string]interface{})
		}
	} else {
		servers = make(map[string]interface{})
	}

	// Add kubectl-go-mcp-server configuration
	servers["kubectl-go-mcp-server"] = map[string]interface{}{
		"type":    "stdio",
		"command": binaryPath,
		"args":    []string{},
		"env": map[string]string{
			"KUBECONFIG": getDefaultKubeconfigPath(),
		},
	}

	mcpConfig["servers"] = servers
	settings["mcp"] = mcpConfig
}

// getDefaultKubeconfigPath returns the default kubeconfig path for the platform
func getDefaultKubeconfigPath() string {
	if kubeconfig := os.Getenv("KUBECONFIG"); kubeconfig != "" {
		return kubeconfig
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "~/.kube/config" // fallback
	}

	return filepath.Join(home, ".kube", "config")
}

// writeVSCodeSettings writes the updated settings back to the file
func writeVSCodeSettings(settingsPath string, settings map[string]interface{}) error {
	// Pretty-print the JSON for better readability
	data, err := json.MarshalIndent(settings, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(settingsPath, data, 0644)
}
