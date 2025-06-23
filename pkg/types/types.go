package types

import (
	"context"
	"encoding/json"
	"fmt"
)

type contextKey string

const (
	KubeconfigKey contextKey = "kubeconfig"
	WorkdirKey    contextKey = "workdir"
)

type Tool interface {
	Name() string
	Description() string
	FunctionDefinition() *FunctionDefinition
	Run(ctx context.Context, args map[string]any) (any, error)
	IsInteractive(args map[string]any) (bool, error)
	CheckModifiesResource(args map[string]any) string
}

type FunctionDefinition struct {
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	Parameters  *Schema `json:"parameters,omitempty"`
}

type Schema struct {
	Type        SchemaType         `json:"type,omitempty"`
	Properties  map[string]*Schema `json:"properties,omitempty"`
	Items       *Schema            `json:"items,omitempty"`
	Description string             `json:"description,omitempty"`
	Required    []string           `json:"required,omitempty"`
}

type SchemaType string

const (
	TypeObject  SchemaType = "object"
	TypeArray   SchemaType = "array"
	TypeString  SchemaType = "string"
	TypeBoolean SchemaType = "boolean"
	TypeInteger SchemaType = "integer"
)

func (s *Schema) ToRawSchema() (json.RawMessage, error) {
	return json.Marshal(s)
}

type ExecResult struct {
	Command    string `json:"command,omitempty"`
	Error      string `json:"error,omitempty"`
	Stdout     string `json:"stdout,omitempty"`
	Stderr     string `json:"stderr,omitempty"`
	ExitCode   int    `json:"exit_code,omitempty"`
	StreamType string `json:"stream_type,omitempty"`
}

func (e *ExecResult) String() string {
	return fmt.Sprintf("Command: %q\nError: %q\nStdout: %q\nStderr: %q\nExitCode: %d\nStreamType: %q",
		e.Command, e.Error, e.Stdout, e.Stderr, e.ExitCode, e.StreamType)
}
