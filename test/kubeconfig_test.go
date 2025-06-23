package test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"kubectl-go-mcp-server/internal/config"
)

func TestGetDefaultKubeconfigPath(t *testing.T) {
	// Save original env vars
	originalKubeconfig := os.Getenv("KUBECONFIG")
	originalHome := os.Getenv("HOME")
	originalUserProfile := os.Getenv("USERPROFILE")

	defer func() {
		// Restore original env vars
		if originalKubeconfig != "" {
			_ = os.Setenv("KUBECONFIG", originalKubeconfig)
		} else {
			_ = os.Unsetenv("KUBECONFIG")
		}
		if originalHome != "" {
			_ = os.Setenv("HOME", originalHome)
		} else {
			_ = os.Unsetenv("HOME")
		}
		if originalUserProfile != "" {
			_ = os.Setenv("USERPROFILE", originalUserProfile)
		} else {
			_ = os.Unsetenv("USERPROFILE")
		}
	}()
	tests := []struct {
		name        string
		kubeconfig  string
		home        string
		userProfile string
		expected    string
		description string
	}{
		{
			name:        "HOME env var set",
			kubeconfig:  "/custom/kubeconfig", // This should be ignored by GetDefaultKubeconfigPath
			home:        "/home/user",
			expected:    "/home/user/.kube/config",
			description: "Should use HOME/.kube/config (ignores KUBECONFIG)",
		},
		{
			name:        "Multiple KUBECONFIG paths ignored",
			kubeconfig:  "/first/config:/second/config", // This should be ignored by GetDefaultKubeconfigPath
			home:        "/home/user",
			expected:    "/home/user/.kube/config",
			description: "Should use HOME/.kube/config (ignores KUBECONFIG)",
		},
		{
			name:        "HOME env var fallback",
			kubeconfig:  "",
			home:        "/home/user",
			expected:    "/home/user/.kube/config",
			description: "Should use HOME/.kube/config when KUBECONFIG not set",
		},
	}

	if runtime.GOOS == "windows" {
		// Add Windows-specific test
		tests = append(tests, struct {
			name        string
			kubeconfig  string
			home        string
			userProfile string
			expected    string
			description string
		}{
			name:        "Windows USERPROFILE fallback",
			kubeconfig:  "",
			home:        "",
			userProfile: "C:\\Users\\user",
			expected:    "C:\\Users\\user\\.kube\\config",
			description: "Should use USERPROFILE on Windows when HOME not available",
		})

		// Add Windows path separator test
		tests = append(tests, struct {
			name        string
			kubeconfig  string
			home        string
			userProfile string
			expected    string
			description string
		}{
			name:        "Windows KUBECONFIG multiple paths",
			kubeconfig:  "C:\\first\\config;C:\\second\\config",
			home:        "C:\\Users\\user",
			expected:    "C:\\first\\config",
			description: "Should handle Windows path separator in KUBECONFIG",
		})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment
			if tt.kubeconfig != "" {
				_ = os.Setenv("KUBECONFIG", tt.kubeconfig)
			} else {
				_ = os.Unsetenv("KUBECONFIG")
			}

			if tt.home != "" {
				_ = os.Setenv("HOME", tt.home)
			} else {
				_ = os.Unsetenv("HOME")
			}

			if tt.userProfile != "" {
				_ = os.Setenv("USERPROFILE", tt.userProfile)
			} else {
				_ = os.Unsetenv("USERPROFILE")
			}

			result := config.GetDefaultKubeconfigPath()

			// Normalize paths for comparison
			expectedNorm := filepath.Clean(tt.expected)
			resultNorm := filepath.Clean(result)

			if resultNorm != expectedNorm {
				t.Errorf("%s: got %q, want %q", tt.description, result, tt.expected)
			}
		})
	}
}

// TestExpandPath is commented out because it tests an unexported function
// The expandPath function is tested indirectly through ValidateKubeconfigPath
// func TestExpandPath(t *testing.T) { ... }

func TestValidateKubeconfigPath(t *testing.T) {
	// Save original HOME
	originalHome := os.Getenv("HOME")
	defer func() {
		if originalHome != "" {
			_ = os.Setenv("HOME", originalHome)
		} else {
			_ = os.Unsetenv("HOME")
		}
	}()

	tests := []struct {
		name     string
		input    string
		home     string
		expected string
		hasError bool
	}{
		{
			name:     "Empty path returns default",
			input:    "",
			home:     "/home/user",
			expected: "/home/user/.kube/config",
			hasError: false,
		},
		{
			name:     "Absolute path unchanged",
			input:    "/custom/kubeconfig",
			home:     "/home/user",
			expected: "/custom/kubeconfig",
			hasError: false,
		},
		{
			name:     "Tilde expansion",
			input:    "~/.kube/custom-config",
			home:     "/home/user",
			expected: "/home/user/.kube/custom-config",
			hasError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.home != "" {
				_ = os.Setenv("HOME", tt.home)
			}

			result, err := config.ValidateKubeconfigPath(tt.input)

			if tt.hasError && err == nil {
				t.Errorf("expected error but got none")
			}

			if !tt.hasError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !tt.hasError {
				// Normalize paths for comparison
				expectedNorm := filepath.Clean(tt.expected)
				resultNorm := filepath.Clean(result)

				if resultNorm != expectedNorm {
					t.Errorf("got %q, want %q", result, tt.expected)
				}
			}
		})
	}
}

// TestGetKubeconfigPaths is commented out because GetKubeconfigPaths function doesn't exist
// This function would parse KUBECONFIG environment variable for multiple paths
// func TestGetKubeconfigPaths(t *testing.T) { ... }

func TestConfigGetKubeconfigPath(t *testing.T) {
	tests := []struct {
		name       string
		configPath string
		home       string
		expected   string
	}{
		{
			name:       "Custom path in config",
			configPath: "/custom/config",
			home:       "/home/user",
			expected:   "/custom/config",
		},
		{
			name:       "Empty config path uses default",
			configPath: "",
			home:       "/home/user",
			expected:   "/home/user/.kube/config",
		},
		{
			name:       "Tilde expansion in config path",
			configPath: "~/.kube/custom",
			home:       "/home/user",
			expected:   "/home/user/.kube/custom",
		},
	}

	// Save original HOME
	originalHome := os.Getenv("HOME")
	originalKubeconfig := os.Getenv("KUBECONFIG")
	defer func() {
		if originalHome != "" {
			_ = os.Setenv("HOME", originalHome)
		} else {
			_ = os.Unsetenv("HOME")
		}
		if originalKubeconfig != "" {
			_ = os.Setenv("KUBECONFIG", originalKubeconfig)
		} else {
			_ = os.Unsetenv("KUBECONFIG")
		}
	}()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear KUBECONFIG to test default behavior
			_ = os.Unsetenv("KUBECONFIG")
			if tt.home != "" {
				_ = os.Setenv("HOME", tt.home)
			}

			cfg := &config.Config{
				Kubeconfig: config.KubeconfigSettings{
					Path: tt.configPath,
				},
			}

			result := cfg.GetKubeconfigPath()

			// Normalize paths for comparison
			expectedNorm := filepath.Clean(tt.expected)
			resultNorm := filepath.Clean(result)

			if resultNorm != expectedNorm {
				t.Errorf("got %q, want %q", result, tt.expected)
			}
		})
	}
}
