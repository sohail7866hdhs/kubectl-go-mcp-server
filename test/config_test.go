package test

import (
	"os"
	"path/filepath"
	"testing"

	"kubectl-go-mcp-server/internal/config"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name       string
		configPath string
		configData string
		wantErr    bool
	}{
		{
			name:       "empty config path returns default",
			configPath: "",
			wantErr:    false,
		},
		{
			name:       "non-existent file returns default",
			configPath: "/non/existent/path",
			wantErr:    false,
		},
		{
			name:       "valid config file",
			configPath: "test-config.json",
			configData: `{"name": "test-server", "debug": true}`,
			wantErr:    false,
		},
		{
			name:       "invalid json",
			configPath: "invalid-config.json",
			configData: `{"name": "test-server"`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test file if configData is provided
			if tt.configData != "" {
				tmpFile := filepath.Join(os.TempDir(), tt.configPath)
				err := os.WriteFile(tmpFile, []byte(tt.configData), 0644)
				if err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
				defer func() {
					if err := os.Remove(tmpFile); err != nil {
						t.Logf("Warning: failed to remove test file: %v", err)
					}
				}()
				tt.configPath = tmpFile
			}

			cfg, err := config.Load(tt.configPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("config.Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && cfg == nil {
				t.Error("config.Load() returned nil config without error")
			}
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := config.DefaultConfig()

	if cfg.Name == "" {
		t.Error("Default config should have a name")
	}

	if cfg.MCP.MaxConcurrentOps <= 0 {
		t.Error("Default config should have positive MaxConcurrentOps")
	}

	if cfg.MCP.OperationTimeout <= 0 {
		t.Error("Default config should have positive OperationTimeout")
	}
}

func TestGetKubeconfigPath(t *testing.T) {
	tests := []struct {
		name           string
		config         *config.Config
		envHome        string
		envUserProfile string
		want           string
	}{
		{
			name: "explicit path",
			config: &config.Config{
				Kubeconfig: config.KubeconfigSettings{
					Path: "/custom/kubeconfig",
				},
			},
			want: "/custom/kubeconfig",
		},
		{
			name: "default with HOME",
			config: &config.Config{
				Kubeconfig: config.KubeconfigSettings{},
			},
			envHome: "/home/user",
			want:    "/home/user/.kube/config",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variables
			if tt.envHome != "" {
				if err := os.Setenv("HOME", tt.envHome); err != nil {
					t.Fatalf("Failed to set HOME: %v", err)
				}
				defer func() {
					if err := os.Unsetenv("HOME"); err != nil {
						t.Logf("Warning: failed to unset HOME: %v", err)
					}
				}()
			}
			if tt.envUserProfile != "" {
				if err := os.Setenv("USERPROFILE", tt.envUserProfile); err != nil {
					t.Fatalf("Failed to set USERPROFILE: %v", err)
				}
				defer func() {
					if err := os.Unsetenv("USERPROFILE"); err != nil {
						t.Logf("Warning: failed to unset USERPROFILE: %v", err)
					}
				}()
			}

			got := tt.config.GetKubeconfigPath()
			if got != tt.want {
				t.Errorf("config.GetKubeconfigPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
