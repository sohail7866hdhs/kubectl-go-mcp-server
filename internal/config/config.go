package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Version     string `json:"version,omitempty"`
	Debug       bool   `json:"debug,omitempty"`

	Kubeconfig KubeconfigSettings `json:"kubeconfig"`

	MCP MCPSettings `json:"mcp"`
}

type KubeconfigSettings struct {
	Path      string `json:"path,omitempty"`
	Context   string `json:"context,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

type MCPSettings struct {
	MaxConcurrentOps int  `json:"maxConcurrentOps,omitempty"`
	OperationTimeout int  `json:"operationTimeout,omitempty"`
	AllowDestructive bool `json:"allowDestructive,omitempty"`
}

func Load(configPath string) (*Config, error) {
	cfg := DefaultConfig()

	if configPath == "" {
		return cfg, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return cfg, nil
}

func (c *Config) Save(configPath string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func DefaultConfig() *Config {
	return &Config{
		Name:        "kubectl-go-mcp-server",
		Description: "Model Context Protocol server for Kubernetes operations using kubectl",
		Version:     "1.0.0",
		Debug:       false,
		Kubeconfig: KubeconfigSettings{
			Path:      "", // Use default kubeconfig location
			Context:   "", // Use current context
			Namespace: "", // Use default namespace
		},
		MCP: MCPSettings{
			MaxConcurrentOps: 5,
			OperationTimeout: 30,
			AllowDestructive: false,
		},
	}
}

func (c *Config) GetKubeconfigPath() string {
	if c.Kubeconfig.Path != "" {
		if expanded, err := ValidateKubeconfigPath(c.Kubeconfig.Path); err == nil {
			return expanded
		}
		return c.Kubeconfig.Path
	}

	return GetDefaultKubeconfigPath()
}
