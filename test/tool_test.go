package test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"kubectl-go-mcp-server/pkg/kubectl"
	"kubectl-go-mcp-server/pkg/types"
)

func TestKubectlTool_Name(t *testing.T) {
	tool := &kubectl.KubectlTool{}
	expected := "kubectl"
	if tool.Name() != expected {
		t.Errorf("Expected name %q, got %q", expected, tool.Name())
	}
}

func TestKubectlTool_Description(t *testing.T) {
	tool := &kubectl.KubectlTool{}
	description := tool.Description()

	// Check that description is not empty and contains key information
	if description == "" {
		t.Error("Description should not be empty")
	}

	if !strings.Contains(description, "kubectl") {
		t.Error("Description should mention kubectl")
	}

	if !strings.Contains(description, "Kubernetes") {
		t.Error("Description should mention Kubernetes")
	}

	if !strings.Contains(description, "Interactive commands") {
		t.Error("Description should mention interactive commands limitation")
	}
}

func TestKubectlTool_FunctionDefinition(t *testing.T) {
	tool := &kubectl.KubectlTool{}
	funcDef := tool.FunctionDefinition()

	if funcDef == nil {
		t.Fatal("FunctionDefinition should not be nil")
	}

	if funcDef.Name != "kubectl" {
		t.Errorf("Expected function name 'kubectl', got %q", funcDef.Name)
	}

	if funcDef.Description == "" {
		t.Error("Function description should not be empty")
	}

	if funcDef.Parameters == nil {
		t.Fatal("Parameters should not be nil")
	}

	if funcDef.Parameters.Type != types.TypeObject {
		t.Errorf("Expected parameters type 'object', got %q", funcDef.Parameters.Type)
	}

	// Check required parameters
	expectedRequired := []string{"command"}
	if len(funcDef.Parameters.Required) != len(expectedRequired) {
		t.Errorf("Expected %d required parameters, got %d", len(expectedRequired), len(funcDef.Parameters.Required))
	}

	for i, req := range expectedRequired {
		if i >= len(funcDef.Parameters.Required) || funcDef.Parameters.Required[i] != req {
			t.Errorf("Expected required parameter %q, got %q", req, funcDef.Parameters.Required[i])
		}
	}

	// Check that command parameter exists
	commandProp, exists := funcDef.Parameters.Properties["command"]
	if !exists {
		t.Error("Command property should exist")
	}

	if commandProp.Type != types.TypeString {
		t.Errorf("Expected command type 'string', got %q", commandProp.Type)
	}

	// Check that modifies_resource parameter exists
	modifiesProp, exists := funcDef.Parameters.Properties["modifies_resource"]
	if !exists {
		t.Error("modifies_resource property should exist")
	}

	if modifiesProp.Type != types.TypeString {
		t.Errorf("Expected modifies_resource type 'string', got %q", modifiesProp.Type)
	}
}

func TestKubectlTool_Run(t *testing.T) {
	tool := &kubectl.KubectlTool{}

	t.Run("Missing command parameter", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), types.KubeconfigKey, "")
		ctx = context.WithValue(ctx, types.WorkdirKey, "/tmp")

		result, err := tool.Run(ctx, map[string]any{})
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		execResult, ok := result.(*types.ExecResult)
		if !ok {
			t.Fatalf("Expected *types.ExecResult, got %T", result)
		}

		if !strings.Contains(execResult.Error, "command not provided") {
			t.Errorf("Expected error about missing command, got %q", execResult.Error)
		}
	})

	t.Run("Nil command parameter", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), types.KubeconfigKey, "")
		ctx = context.WithValue(ctx, types.WorkdirKey, "/tmp")

		result, err := tool.Run(ctx, map[string]any{"command": nil})
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		execResult, ok := result.(*types.ExecResult)
		if !ok {
			t.Fatalf("Expected *types.ExecResult, got %T", result)
		}

		if !strings.Contains(execResult.Error, "nil") {
			t.Errorf("Expected error about nil command, got %q", execResult.Error)
		}
	})

	t.Run("Invalid command type", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), types.KubeconfigKey, "")
		ctx = context.WithValue(ctx, types.WorkdirKey, "/tmp")

		result, err := tool.Run(ctx, map[string]any{"command": 123})
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		execResult, ok := result.(*types.ExecResult)
		if !ok {
			t.Fatalf("Expected *types.ExecResult, got %T", result)
		}

		if !strings.Contains(execResult.Error, "must be a string") {
			t.Errorf("Expected error about string type, got %q", execResult.Error)
		}
	})

	t.Run("Interactive command", func(t *testing.T) {
		ctx := context.WithValue(context.Background(), types.KubeconfigKey, "")
		ctx = context.WithValue(ctx, types.WorkdirKey, "/tmp")

		result, err := tool.Run(ctx, map[string]any{"command": "kubectl exec pod-name -it -- /bin/bash"})
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		execResult, ok := result.(*types.ExecResult)
		if !ok {
			t.Fatalf("Expected *types.ExecResult, got %T", result)
		}

		if !strings.Contains(execResult.Error, "interactive mode not supported") {
			t.Errorf("Expected interactive mode error, got %q", execResult.Error)
		}
	})
}

func TestKubectlTool_IsInteractive(t *testing.T) {
	tool := &kubectl.KubectlTool{}

	t.Run("Missing command", func(t *testing.T) {
		interactive, err := tool.IsInteractive(map[string]any{})
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if interactive {
			t.Error("Expected false for missing command")
		}
	})

	t.Run("Nil command", func(t *testing.T) {
		interactive, err := tool.IsInteractive(map[string]any{"command": nil})
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if interactive {
			t.Error("Expected false for nil command")
		}
	})

	t.Run("Invalid command type", func(t *testing.T) {
		interactive, err := tool.IsInteractive(map[string]any{"command": 123})
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if interactive {
			t.Error("Expected false for invalid command type")
		}
	})

	t.Run("Non-interactive command", func(t *testing.T) {
		interactive, err := tool.IsInteractive(map[string]any{"command": "kubectl get pods"})
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if interactive {
			t.Error("Expected false for non-interactive command")
		}
	})

	t.Run("Interactive exec command", func(t *testing.T) {
		interactive, err := tool.IsInteractive(map[string]any{"command": "kubectl exec pod-name -it -- /bin/bash"})
		if err == nil {
			t.Error("Expected error for interactive command")
		}
		if !interactive {
			t.Error("Expected true for interactive command")
		}
	})
}

func TestKubectlTool_CheckModifiesResource(t *testing.T) {
	tool := &kubectl.KubectlTool{}

	t.Run("Invalid command type", func(t *testing.T) {
		result := tool.CheckModifiesResource(map[string]any{"command": 123})
		if result != "unknown" {
			t.Errorf("Expected 'unknown' for invalid command type, got %q", result)
		}
	})

	t.Run("Read-only command", func(t *testing.T) {
		result := tool.CheckModifiesResource(map[string]any{"command": "kubectl get pods"})
		if result != "no" {
			t.Errorf("Expected 'no' for read-only command, got %q", result)
		}
	})

	t.Run("Modifying command", func(t *testing.T) {
		result := tool.CheckModifiesResource(map[string]any{"command": "kubectl delete pod my-pod"})
		if result != "yes" {
			t.Errorf("Expected 'yes' for modifying command, got %q", result)
		}
	})
}

func TestIsInteractiveCommand(t *testing.T) {
	tests := []struct {
		name        string
		command     string
		interactive bool
		expectError bool
	}{
		{"Empty command", "", false, false},
		{"Non-kubectl command", "ls -la", false, false},
		{"Get pods", "kubectl get pods", false, false},
		{"Describe deployment", "kubectl describe deployment app", false, false},
		{"Interactive exec", "kubectl exec pod-name -it -- /bin/bash", true, true},
		{"Port forward", "kubectl port-forward pod-name 8080:80", true, true},
		{"Edit resource", "kubectl edit deployment app", true, true},
		{"Non-interactive exec", "kubectl exec pod-name -- ps aux", false, false},
		{"Logs", "kubectl logs pod-name", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			interactive, err := kubectl.IsInteractiveCommand(tt.command)

			if (err != nil) != tt.expectError {
				t.Errorf("Expected error: %v, got: %v", tt.expectError, err)
			}

			if interactive != tt.interactive {
				t.Errorf("Expected interactive: %v, got: %v", tt.interactive, interactive)
			}
		})
	}
}

func TestModifiesResource(t *testing.T) {
	tests := []struct {
		name     string
		command  string
		expected string
	}{
		{"Empty command", "", "unknown"},
		{"Short command", "kubectl", "unknown"},
		{"Non-kubectl", "docker ps", "unknown"},
		{"Get pods", "kubectl get pods", "no"},
		{"Describe", "kubectl describe deployment app", "no"},
		{"Logs", "kubectl logs pod-name", "no"},
		{"Top", "kubectl top nodes", "no"},
		{"Version", "kubectl version", "no"},
		{"Cluster info", "kubectl cluster-info", "no"},
		{"Config", "kubectl config view", "no"},
		{"Create", "kubectl create deployment app --image=nginx", "yes"},
		{"Apply", "kubectl apply -f deployment.yaml", "yes"},
		{"Delete", "kubectl delete pod my-pod", "yes"},
		{"Patch", "kubectl patch deployment app -p '{\"spec\":{\"replicas\":3}}'", "yes"},
		{"Replace", "kubectl replace -f deployment.yaml", "yes"},
		{"Scale", "kubectl scale deployment app --replicas=3", "yes"},
		{"Rollout", "kubectl rollout restart deployment app", "yes"},
		{"Annotate", "kubectl annotate pods my-pod key=value", "yes"},
		{"Label", "kubectl label pods my-pod key=value", "yes"},
		{"Exec", "kubectl exec pod-name -- ps aux", "no"},
		{"Port forward", "kubectl port-forward pod-name 8080:80", "no"},
		{"Proxy", "kubectl proxy", "no"},
		{"Unknown verb", "kubectl unknown-command", "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := kubectl.ModifiesResource(tt.command)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestLookupBashBin(t *testing.T) {
	result := kubectl.LookupBashBin()
	if result == "" {
		t.Error("LookupBashBin should not return empty string")
	}

	// Should either find bash or return default path
	if !strings.Contains(result, "bash") {
		t.Errorf("Expected result to contain 'bash', got %q", result)
	}
}

func TestExpandShellVar(t *testing.T) {
	// Set up test environment
	if os.Getenv("HOME") == "" {
		t.Setenv("HOME", "/tmp/test-home")
	}

	tests := []struct {
		name     string
		path     string
		expected func(string) bool // function to validate result
		wantErr  bool
	}{
		{
			name: "Home directory expansion",
			path: "~/.kube/config",
			expected: func(result string) bool {
				return !strings.HasPrefix(result, "~") && strings.HasSuffix(result, "/.kube/config")
			},
			wantErr: false,
		},
		{
			name: "Absolute path no expansion",
			path: "/absolute/path",
			expected: func(result string) bool {
				return result == "/absolute/path"
			},
			wantErr: false,
		},
		{
			name: "Relative path no expansion",
			path: "relative/path",
			expected: func(result string) bool {
				return result == "relative/path"
			},
			wantErr: false,
		}, {
			name: "Just tilde",
			path: "~",
			expected: func(result string) bool {
				return result == "~" // expandShellVar only handles ~/path, not bare ~
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := kubectl.ExpandShellVar(tt.path)

			if (err != nil) != tt.wantErr {
				t.Errorf("ExpandShellVar() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !tt.expected(result) {
				t.Errorf("ExpandShellVar() = %q, validation failed", result)
			}
		})
	}
}

func TestRunKubectlCommand(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := filepath.Join(os.TempDir(), "kubectl-test")
	if err := os.MkdirAll(tmpDir, 0755); err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Logf("Warning: failed to remove temp directory: %v", err)
		}
	}()

	t.Run("Interactive command rejection", func(t *testing.T) {
		ctx := context.Background()
		result, err := kubectl.RunKubectlCommand(ctx, "kubectl exec pod-name -it -- /bin/bash", tmpDir, "")

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if result == nil {
			t.Fatal("Expected result, got nil")
		}

		if !strings.Contains(result.Error, "interactive mode not supported") {
			t.Errorf("Expected interactive mode error, got %q", result.Error)
		}
	})

	t.Run("Invalid command", func(t *testing.T) {
		ctx := context.Background()
		result, err := kubectl.RunKubectlCommand(ctx, "kubectl-nonexistent-command", tmpDir, "")

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if result == nil {
			t.Fatal("Expected result, got nil")
		}

		// Command should fail with non-zero exit code
		if result.ExitCode == 0 && result.Error == "" {
			t.Error("Expected command to fail")
		}
	})

	t.Run("Valid non-destructive command", func(t *testing.T) {
		ctx := context.Background()
		// Use kubectl version as it should work even without cluster access
		result, err := kubectl.RunKubectlCommand(ctx, "kubectl version --client", tmpDir, "")

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if result == nil {
			t.Fatal("Expected result, got nil")
		}

		// Command should have some output (either success or client-side error)
		if result.Stdout == "" && result.Error == "" {
			t.Error("Expected some output from kubectl version command")
		}
	})
}
