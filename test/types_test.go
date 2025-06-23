package test

import (
	"context"
	"encoding/json"
	"testing"

	"kubectl-go-mcp-server/pkg/types"
)

func TestSchemaType_Constants(t *testing.T) {
	tests := []struct {
		name       string
		schemaType types.SchemaType
		expected   string
	}{
		{"types.TypeObject", types.TypeObject, "object"},
		{"types.TypeArray", types.TypeArray, "array"},
		{"types.TypeString", types.TypeString, "string"},
		{"types.TypeBoolean", types.TypeBoolean, "boolean"},
		{"types.TypeInteger", types.TypeInteger, "integer"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.schemaType) != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, string(tt.schemaType))
			}
		})
	}
}

func TestSchema_ToRawSchema(t *testing.T) {
	tests := []struct {
		name    string
		schema  *types.Schema
		wantErr bool
	}{
		{
			name: "Simple string schema",
			schema: &types.Schema{
				Type:        types.TypeString,
				Description: "A simple string",
			},
			wantErr: false,
		},
		{
			name: "Object schema with properties",
			schema: &types.Schema{
				Type: types.TypeObject,
				Properties: map[string]*types.Schema{
					"name": {
						Type:        types.TypeString,
						Description: "Name field",
					},
					"age": {
						Type:        types.TypeInteger,
						Description: "Age field",
					},
				},
				Required: []string{"name"},
			},
			wantErr: false,
		},
		{
			name: "Array schema",
			schema: &types.Schema{
				Type: types.TypeArray,
				Items: &types.Schema{
					Type: types.TypeString,
				},
			},
			wantErr: false,
		},
		{
			name:    "Empty schema",
			schema:  &types.Schema{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rawSchema, err := tt.schema.ToRawSchema()
			if (err != nil) != tt.wantErr {
				t.Errorf("types.Schema.ToRawSchema() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify the JSON is valid by unmarshaling
				var result map[string]interface{}
				if err := json.Unmarshal(rawSchema, &result); err != nil {
					t.Errorf("Generated JSON is invalid: %v", err)
				}

				// Verify that marshaling back produces the same structure
				backToSchema, err := json.Marshal(tt.schema)
				if err != nil {
					t.Errorf("Failed to marshal original schema: %v", err)
				}

				if string(rawSchema) != string(backToSchema) {
					t.Errorf("ToRawSchema() produced different JSON than direct marshal")
				}
			}
		})
	}
}

func TestExecResult_String(t *testing.T) {
	tests := []struct {
		name     string
		result   *types.ExecResult
		expected string
	}{
		{
			name: "Complete result",
			result: &types.ExecResult{
				Command:    "kubectl get pods",
				Error:      "connection failed",
				Stdout:     "pod1   Running",
				Stderr:     "warning: deprecated",
				ExitCode:   1,
				StreamType: "text",
			},
			expected: `Command: "kubectl get pods"
Error: "connection failed"
Stdout: "pod1   Running"
Stderr: "warning: deprecated"
ExitCode: 1
StreamType: "text"`,
		},
		{
			name:   "Empty result",
			result: &types.ExecResult{},
			expected: `Command: ""
Error: ""
Stdout: ""
Stderr: ""
ExitCode: 0
StreamType: ""`,
		},
		{
			name: "Success result",
			result: &types.ExecResult{
				Command:  "kubectl version",
				Stdout:   "Client Version: v1.28.0",
				ExitCode: 0,
			},
			expected: `Command: "kubectl version"
Error: ""
Stdout: "Client Version: v1.28.0"
Stderr: ""
ExitCode: 0
StreamType: ""`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.result.String()
			if actual != tt.expected {
				t.Errorf("types.ExecResult.String() = %q, want %q", actual, tt.expected)
			}
		})
	}
}

func TestFunctionDefinition(t *testing.T) {
	t.Run("Complete function definition", func(t *testing.T) {
		funcDef := &types.FunctionDefinition{
			Name:        "test-tool",
			Description: "A test tool",
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
		}

		if funcDef.Name != "test-tool" {
			t.Errorf("Expected name 'test-tool', got %q", funcDef.Name)
		}
		if funcDef.Description != "A test tool" {
			t.Errorf("Expected description 'A test tool', got %q", funcDef.Description)
		}
		if funcDef.Parameters == nil {
			t.Error("Expected parameters to be set")
		}
		if funcDef.Parameters.Type != types.TypeObject {
			t.Errorf("Expected parameters type 'object', got %q", funcDef.Parameters.Type)
		}
	})

	t.Run("Minimal function definition", func(t *testing.T) {
		funcDef := &types.FunctionDefinition{
			Name: "minimal",
		}

		if funcDef.Name != "minimal" {
			t.Errorf("Expected name 'minimal', got %q", funcDef.Name)
		}
		if funcDef.Description != "" {
			t.Errorf("Expected empty description, got %q", funcDef.Description)
		}
		if funcDef.Parameters != nil {
			t.Error("Expected parameters to be nil")
		}
	})
}

func TestMockTool(t *testing.T) {
	t.Run("MockTool implements types.Tool interface", func(t *testing.T) {
		funcDef := &types.FunctionDefinition{
			Name:        "mock-tool",
			Description: "A mock tool for testing",
		}

		mock := &MockTool{
			name:        "mock-tool",
			description: "A mock tool for testing",
			funcDef:     funcDef,
			runResult:   "test result",
			runError:    nil,
			interactive: false,
			modifies:    "no",
		}

		// Test types.Tool interface methods
		var tool types.Tool = mock // This line will fail to compile if MockTool doesn't implement types.Tool

		if tool.Name() != "mock-tool" {
			t.Errorf("Expected name 'mock-tool', got %q", tool.Name())
		}
		if tool.Description() != "A mock tool for testing" {
			t.Errorf("Expected description 'A mock tool for testing', got %q", tool.Description())
		}

		def := tool.FunctionDefinition()
		if def != funcDef {
			t.Error("Expected same function definition")
		}

		result, err := tool.Run(context.Background(), map[string]any{})
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if result != "test result" {
			t.Errorf("Expected 'test result', got %v", result)
		}

		interactive, err := tool.IsInteractive(map[string]any{})
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if interactive != false {
			t.Errorf("Expected false, got %v", interactive)
		}

		modifies := tool.CheckModifiesResource(map[string]any{})
		if modifies != "no" {
			t.Errorf("Expected 'no', got %q", modifies)
		}
	})
}

func TestExecResultJSONSerialization(t *testing.T) {
	t.Run("JSON marshaling and unmarshaling", func(t *testing.T) {
		original := &types.ExecResult{
			Command:    "kubectl get pods",
			Error:      "some error",
			Stdout:     "output data",
			Stderr:     "error data",
			ExitCode:   1,
			StreamType: "text",
		}

		// Marshal to JSON
		data, err := json.Marshal(original)
		if err != nil {
			t.Fatalf("Failed to marshal types.ExecResult: %v", err)
		}

		// Unmarshal from JSON
		var unmarshaled types.ExecResult
		if err := json.Unmarshal(data, &unmarshaled); err != nil {
			t.Fatalf("Failed to unmarshal types.ExecResult: %v", err)
		}

		// Compare fields
		if unmarshaled.Command != original.Command {
			t.Errorf("Command mismatch: expected %q, got %q", original.Command, unmarshaled.Command)
		}
		if unmarshaled.Error != original.Error {
			t.Errorf("Error mismatch: expected %q, got %q", original.Error, unmarshaled.Error)
		}
		if unmarshaled.Stdout != original.Stdout {
			t.Errorf("Stdout mismatch: expected %q, got %q", original.Stdout, unmarshaled.Stdout)
		}
		if unmarshaled.Stderr != original.Stderr {
			t.Errorf("Stderr mismatch: expected %q, got %q", original.Stderr, unmarshaled.Stderr)
		}
		if unmarshaled.ExitCode != original.ExitCode {
			t.Errorf("ExitCode mismatch: expected %d, got %d", original.ExitCode, unmarshaled.ExitCode)
		}
		if unmarshaled.StreamType != original.StreamType {
			t.Errorf("StreamType mismatch: expected %q, got %q", original.StreamType, unmarshaled.StreamType)
		}
	})
}
