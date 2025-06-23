package test

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"kubectl-go-mcp-server/internal/mcp"
	"kubectl-go-mcp-server/pkg/types"
)

// MockTool for testing
type MockTool struct {
	name        string
	description string
	funcDef     *types.FunctionDefinition
	runResult   any
	runError    error
	interactive bool
	modifies    string
}

func (m *MockTool) Name() string {
	return m.name
}

func (m *MockTool) Description() string {
	return m.description
}

func (m *MockTool) FunctionDefinition() *types.FunctionDefinition {
	return m.funcDef
}

func (m *MockTool) Run(ctx context.Context, args map[string]any) (any, error) {
	return m.runResult, m.runError
}

func (m *MockTool) IsInteractive(args map[string]any) (bool, error) {
	return m.interactive, nil
}

func (m *MockTool) CheckModifiesResource(args map[string]any) string {
	return m.modifies
}

func TestNewTools(t *testing.T) {
	tools := mcp.NewTools()
	if tools == nil {
		t.Fatal("mcp.NewTools() should not return nil")
	}

	if tools.Count() != 0 {
		t.Errorf("Expected empty tools map, got %d tools", tools.Count())
	}
}

func TestTools_RegisterTool(t *testing.T) {
	tools := mcp.NewTools()

	mockTool := &MockTool{
		name:        "test-tool",
		description: "A test tool",
		funcDef: &types.FunctionDefinition{
			Name:        "test-tool",
			Description: "A test tool",
		},
	}
	t.Run("Register new tool", func(t *testing.T) {
		tools.RegisterTool(mockTool)

		if tools.Count() != 1 {
			t.Errorf("Expected 1 tool, got %d", tools.Count())
		}

		if !tools.HasTool("test-tool") {
			t.Error("Tool should be registered")
		}

		retrieved := tools.Lookup("test-tool")
		if retrieved != mockTool {
			t.Error("Retrieved tool does not match registered tool")
		}
	})

	t.Run("Register duplicate tool panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic when registering duplicate tool")
			}
		}()

		tools.RegisterTool(mockTool)
	})
}

func TestTools_Lookup(t *testing.T) {
	tools := mcp.NewTools()

	mockTool := &MockTool{
		name: "test-tool",
	}

	tools.RegisterTool(mockTool)

	t.Run("Lookup existing tool", func(t *testing.T) {
		result := tools.Lookup("test-tool")
		if result != mockTool {
			t.Error("Lookup should return the registered tool")
		}
	})

	t.Run("Lookup non-existing tool", func(t *testing.T) {
		result := tools.Lookup("non-existing")
		if result != nil {
			t.Error("Lookup should return nil for non-existing tool")
		}
	})
}

func TestTools_AllTools(t *testing.T) {
	tools := mcp.NewTools()

	mockTool1 := &MockTool{name: "tool1"}
	mockTool2 := &MockTool{name: "tool2"}

	tools.RegisterTool(mockTool1)
	tools.RegisterTool(mockTool2)

	allTools := tools.AllTools()

	if len(allTools) != 2 {
		t.Errorf("Expected 2 tools, got %d", len(allTools))
	}

	// Check that both tools are present (order doesn't matter)
	foundTool1, foundTool2 := false, false
	for _, tool := range allTools {
		if tool == mockTool1 {
			foundTool1 = true
		}
		if tool == mockTool2 {
			foundTool2 = true
		}
	}

	if !foundTool1 {
		t.Error("tool1 not found in AllTools() result")
	}
	if !foundTool2 {
		t.Error("tool2 not found in AllTools() result")
	}
}

func TestToolResultToMap(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected map[string]interface{}
		wantErr  bool
	}{{
		name: "ExecResult struct",
		input: &types.ExecResult{
			Command:  "kubectl get pods",
			Stdout:   "pod1   Running",
			ExitCode: 1, // Use non-zero to avoid omitempty
		},
		// Only check for the fields that are actually set
		expected: map[string]interface{}{
			"command":   "kubectl get pods",
			"stdout":    "pod1   Running",
			"exit_code": float64(1), // JSON numbers are float64
		},
		wantErr: false,
	},
		{
			name: "Simple map",
			input: map[string]interface{}{
				"key1": "value1",
				"key2": 42,
			},
			expected: map[string]interface{}{
				"key1": "value1",
				"key2": float64(42), // JSON numbers are float64
			},
			wantErr: false,
		},
		{
			name:  "String value",
			input: "simple string",
			expected: map[string]interface{}{
				"output": "simple string",
			},
			wantErr: false,
		},
		{
			name:  "Integer value",
			input: 42,
			expected: map[string]interface{}{
				"output": "42",
			},
			wantErr: false,
		}, {
			name:     "Nil value",
			input:    nil,
			expected: map[string]interface{}{}, // nil marshals to empty map
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := mcp.ToolResultToMap(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("mcp.ToolResultToMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr { // For ExecResult, check the specific fields rather than exact match
				if tt.name == "ExecResult struct" {
					if command, ok := result["command"].(string); !ok || command != "kubectl get pods" {
						t.Errorf("Expected command 'kubectl get pods', got %v", result["command"])
					}
					if stdout, ok := result["stdout"].(string); !ok || stdout != "pod1   Running" {
						t.Errorf("Expected stdout 'pod1   Running', got %v", result["stdout"])
					}
					if exitCode, ok := result["exit_code"].(float64); !ok || exitCode != 1 {
						t.Errorf("Expected exit_code 1, got %v", result["exit_code"])
					}
				} else if tt.name == "Nil value" {
					// For nil, just check that we get an empty map
					if len(result) != 0 {
						t.Errorf("Expected empty map for nil, got %v", result)
					}
				} else if !reflect.DeepEqual(result, tt.expected) {
					t.Errorf("mcp.ToolResultToMap() = %v, want %v", result, tt.expected)
				}
			}
		})
	}
}

func TestNewServer(t *testing.T) {
	t.Run("Valid server creation", func(t *testing.T) {
		server, err := mcp.NewServer("/path/to/kubeconfig", "/tmp/workdir")

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if server == nil {
			t.Fatal("mcp.NewServer should not return nil")
		}

		if server.GetKubectlConfig() != "/path/to/kubeconfig" {
			t.Errorf("Expected kubectlConfig '/path/to/kubeconfig', got %q", server.GetKubectlConfig())
		}

		if server.GetWorkDir() != "/tmp/workdir" {
			t.Errorf("Expected workDir '/tmp/workdir', got %q", server.GetWorkDir())
		}

		if server.GetTools() == nil {
			t.Error("Server tools should be initialized")
		}
		// Check that kubectl tool is registered
		kubectlTool := server.GetTools().Lookup("kubectl")
		if kubectlTool == nil {
			t.Error("kubectl tool should be registered")
		}
	})
}

// Test edge cases and error conditions
func TestToolResultToMap_EdgeCases(t *testing.T) {
	t.Run("Complex nested structure", func(t *testing.T) {
		complexInput := map[string]interface{}{
			"nested": map[string]interface{}{
				"inner": "value",
				"array": []interface{}{1, 2, 3},
			},
			"bool": true,
		}

		result, err := mcp.ToolResultToMap(complexInput)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		// Verify the structure is preserved
		if nested, ok := result["nested"].(map[string]interface{}); ok {
			if inner, ok := nested["inner"].(string); !ok || inner != "value" {
				t.Error("Nested string value not preserved")
			}
			if array, ok := nested["array"].([]interface{}); !ok || len(array) != 3 {
				t.Error("Nested array not preserved")
			}
		} else {
			t.Error("Nested map not preserved")
		}

		if boolVal, ok := result["bool"].(bool); !ok || !boolVal {
			t.Error("Boolean value not preserved")
		}
	})
	t.Run("Channel type (unmarshalable)", func(t *testing.T) {
		// Channels cannot be marshaled to JSON
		ch := make(chan int)
		defer close(ch)

		result, err := mcp.ToolResultToMap(ch)

		// Should error when marshaling fails
		if err == nil {
			t.Error("Expected error when marshaling channel")
		}

		if result != nil {
			t.Errorf("Expected nil result when error occurs, got %v", result)
		}
	})
}

func TestServer_ErrorConditions(t *testing.T) {
	t.Run("Tool with invalid schema", func(t *testing.T) {
		// Create a mock tool that would cause schema conversion to fail
		// This is tricky to test directly without modifying the Schema.ToRawSchema method
		// For now, we'll test with a valid schema to ensure the happy path works

		mockTool := &MockTool{
			name:        "invalid-tool",
			description: "A tool with valid schema for testing",
			funcDef: &types.FunctionDefinition{
				Name:        "invalid-tool",
				Description: "A tool with valid schema",
				Parameters: &types.Schema{
					Type: types.TypeObject,
					Properties: map[string]*types.Schema{
						"param": {
							Type:        types.TypeString,
							Description: "A parameter",
						},
					},
				},
			},
		}

		// This should work fine
		tools := mcp.NewTools()
		tools.RegisterTool(mockTool)

		server, err := mcp.NewServer("/path/to/kubeconfig", "/tmp/workdir")
		if err != nil {
			t.Errorf("Unexpected error creating server: %v", err)
		}

		if server == nil {
			t.Error("Server should not be nil")
		}
	})
}

// Benchmark tests
func BenchmarkToolResultToMap(b *testing.B) {
	execResult := &types.ExecResult{
		Command:    "kubectl get pods",
		Stdout:     "pod1   Running\npod2   Pending",
		Stderr:     "",
		ExitCode:   0,
		StreamType: "text",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := mcp.ToolResultToMap(execResult)
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
	}
}

func BenchmarkTools_Lookup(b *testing.B) {
	tools := mcp.NewTools()

	// Register multiple tools
	for i := 0; i < 100; i++ {
		mockTool := &MockTool{
			name: fmt.Sprintf("tool-%d", i),
		}
		tools.RegisterTool(mockTool)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Lookup a tool that exists
		tools.Lookup("tool-50")
	}
}

// Integration-style test that tests multiple components together
func TestToolsIntegration(t *testing.T) {
	tools := mcp.NewTools()

	// Create a mock tool that simulates realistic behavior
	mockTool := &MockTool{
		name:        "integration-tool",
		description: "A tool for integration testing",
		funcDef: &types.FunctionDefinition{
			Name:        "integration-tool",
			Description: "A tool for integration testing",
			Parameters: &types.Schema{
				Type: types.TypeObject,
				Properties: map[string]*types.Schema{
					"command": {
						Type:        types.TypeString,
						Description: "Command to execute",
					},
				},
				Required: []string{"command"},
			},
		},
		runResult: &types.ExecResult{
			Command:  "test command",
			Stdout:   "test output",
			ExitCode: 0,
		},
		runError:    nil,
		interactive: false,
		modifies:    "no",
	}

	// Register the tool
	tools.RegisterTool(mockTool)

	// Lookup the tool
	tool := tools.Lookup("integration-tool")
	if tool == nil {
		t.Fatal("Tool lookup failed")
	}

	// Test the tool methods
	if tool.Name() != "integration-tool" {
		t.Error("Tool name mismatch")
	}

	// Test function definition
	funcDef := tool.FunctionDefinition()
	if funcDef == nil {
		t.Fatal("Function definition is nil")
	}

	// Test schema conversion
	_, err := funcDef.Parameters.ToRawSchema()
	if err != nil {
		t.Errorf("Schema conversion failed: %v", err)
	}

	// Test tool execution
	ctx := context.Background()
	result, err := tool.Run(ctx, map[string]any{"command": "test"})
	if err != nil {
		t.Errorf("Tool execution failed: %v", err)
	}

	// Test result conversion
	resultMap, err := mcp.ToolResultToMap(result)
	if err != nil {
		t.Errorf("Result conversion failed: %v", err)
	}

	if command, ok := resultMap["command"].(string); !ok || command != "test command" {
		t.Error("Result conversion did not preserve command")
	}
}

// Test security validations in the MCP server
func TestMCPServerSecurity(t *testing.T) {
	t.Run("Only kubectl tool is allowed", func(t *testing.T) {
		server, err := mcp.NewServer("/path/to/kubeconfig", "/tmp/workdir")
		if err != nil {
			t.Fatalf("Unexpected error creating server: %v", err)
		}

		// Verify only kubectl tool is registered
		tools := server.GetTools().AllTools()
		if len(tools) != 1 {
			t.Errorf("Expected exactly 1 tool, got %d", len(tools))
		}

		kubectlTool := server.GetTools().Lookup("kubectl")
		if kubectlTool == nil {
			t.Error("kubectl tool should be registered")
		}

		// Verify other tools are not registered
		if bashTool := server.GetTools().Lookup("bash"); bashTool != nil {
			t.Error("bash tool should not be registered")
		}
		if shellTool := server.GetTools().Lookup("shell"); shellTool != nil {
			t.Error("shell tool should not be registered")
		}
		if execTool := server.GetTools().Lookup("exec"); execTool != nil {
			t.Error("exec tool should not be registered")
		}
	})

	t.Run("kubectl tool validates commands properly", func(t *testing.T) {
		server, err := mcp.NewServer("/path/to/kubeconfig", "/tmp/workdir")
		if err != nil {
			t.Fatalf("Unexpected error creating server: %v", err)
		}

		kubectlTool := server.GetTools().Lookup("kubectl")
		if kubectlTool == nil {
			t.Fatal("kubectl tool should be registered")
		}

		// Create context with required values
		ctx := context.WithValue(context.Background(), types.KubeconfigKey, "/path/to/kubeconfig")
		ctx = context.WithValue(ctx, types.WorkdirKey, "/tmp/workdir")

		// Test valid kubectl command
		validArgs := map[string]any{
			"command": "kubectl get pods",
		}
		result, err := kubectlTool.Run(ctx, validArgs)
		if err != nil {
			t.Errorf("Valid kubectl command should not return error: %v", err)
		}
		if result == nil {
			t.Error("Valid kubectl command should return result")
		}

		// Test invalid command injection
		invalidArgs := map[string]any{
			"command": "kubectl get pods; rm -rf /",
		}
		result, err = kubectlTool.Run(ctx, invalidArgs)
		if err != nil {
			t.Errorf("Command injection should be handled gracefully: %v", err)
		}
		if execResult, ok := result.(*types.ExecResult); ok {
			if !strings.Contains(execResult.Error, "Security violation") {
				t.Errorf("Expected security violation error, got: %s", execResult.Error)
			}
		} else {
			t.Error("Expected ExecResult with security violation error")
		}

		// Test non-kubectl command
		nonKubectlArgs := map[string]any{
			"command": "ls -la",
		}
		result, err = kubectlTool.Run(ctx, nonKubectlArgs)
		if err != nil {
			t.Errorf("Non-kubectl command should be handled gracefully: %v", err)
		}
		if execResult, ok := result.(*types.ExecResult); ok {
			if !strings.Contains(execResult.Error, "Security violation") {
				t.Errorf("Expected security violation error for non-kubectl command, got: %s", execResult.Error)
			}
		} else {
			t.Error("Expected ExecResult with security violation error for non-kubectl command")
		}
	})

	t.Run("Interactive commands are properly blocked", func(t *testing.T) {
		server, err := mcp.NewServer("/path/to/kubeconfig", "/tmp/workdir")
		if err != nil {
			t.Fatalf("Unexpected error creating server: %v", err)
		}

		kubectlTool := server.GetTools().Lookup("kubectl")
		if kubectlTool == nil {
			t.Fatal("kubectl tool should be registered")
		}

		// Create context with required values
		ctx := context.WithValue(context.Background(), types.KubeconfigKey, "/path/to/kubeconfig")
		ctx = context.WithValue(ctx, types.WorkdirKey, "/tmp/workdir")

		// Test interactive exec command
		interactiveArgs := map[string]any{
			"command": "kubectl exec -it mypod -- /bin/bash",
		}
		result, err := kubectlTool.Run(ctx, interactiveArgs)
		if err != nil {
			t.Errorf("Interactive command should be handled gracefully: %v", err)
		}
		if execResult, ok := result.(*types.ExecResult); ok {
			if !strings.Contains(execResult.Error, "interactive mode not supported") {
				t.Errorf("Expected interactive mode error, got: %s", execResult.Error)
			}
		} else {
			t.Error("Expected ExecResult with interactive mode error")
		}

		// Test port-forward command
		portForwardArgs := map[string]any{
			"command": "kubectl port-forward pod/mypod 8080:80",
		}
		result, err = kubectlTool.Run(ctx, portForwardArgs)
		if err != nil {
			t.Errorf("Port-forward command should be handled gracefully: %v", err)
		}
		if execResult, ok := result.(*types.ExecResult); ok {
			if !strings.Contains(execResult.Error, "interactive mode not supported") {
				t.Errorf("Expected interactive mode error for port-forward, got: %s", execResult.Error)
			}
		} else {
			t.Error("Expected ExecResult with interactive mode error for port-forward")
		}
	})
}
