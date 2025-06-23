package test

import (
	"context"
	"testing"
	"time"

	"kubectl-go-mcp-server/pkg/kubectl"
	"kubectl-go-mcp-server/pkg/types"
)

// TestKubectlIntegration tests the kubectl tool integration
// This test requires kubectl to be installed and configured
func TestKubectlIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tool := &kubectl.KubectlTool{}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Set context values that the tool expects
	ctx = context.WithValue(ctx, types.KubeconfigKey, "")
	ctx = context.WithValue(ctx, types.WorkdirKey, "/tmp")

	t.Run("GetVersion", func(t *testing.T) {
		args := map[string]any{
			"command": "kubectl version --client",
		}

		result, err := tool.Run(ctx, args)
		if err != nil {
			t.Logf("kubectl version test failed (expected if kubectl not available): %v", err)
			return
		}
		if execResult, ok := result.(*types.ExecResult); ok {
			if execResult.Error != "" {
				t.Logf("kubectl version command failed: %s", execResult.Error)
				return
			}
			t.Logf("kubectl version successful: %s", execResult.Command)
		}
	})

	t.Run("CheckModifiesResource", func(t *testing.T) {
		tests := []struct {
			name     string
			args     map[string]any
			expected string
		}{
			{
				name: "get command does not modify",
				args: map[string]any{
					"command": "kubectl get pods",
				},
				expected: "no",
			},
			{
				name: "delete command modifies",
				args: map[string]any{
					"command": "kubectl delete pod test",
				},
				expected: "yes",
			},
			{
				name: "unknown command",
				args: map[string]any{
					"command": "kubectl unknown-verb",
				},
				expected: "unknown",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := tool.CheckModifiesResource(tt.args)
				if result != tt.expected {
					t.Errorf("CheckModifiesResource() = %s, expected %s", result, tt.expected)
				}
			})
		}
	})

	t.Run("IsInteractive", func(t *testing.T) {
		tests := []struct {
			name     string
			args     map[string]any
			expected bool
		}{
			{
				name: "regular command not interactive",
				args: map[string]any{
					"command": "kubectl get pods",
				},
				expected: false,
			},
			{
				name: "exec with -it is interactive",
				args: map[string]any{
					"command": "kubectl exec -it pod test -- bash",
				},
				expected: true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result, _ := tool.IsInteractive(tt.args)
				if result != tt.expected {
					t.Errorf("IsInteractive() = %v, expected %v", result, tt.expected)
				}
			})
		}
	})
}
