package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"kubectl-go-mcp-server/pkg/kubectl"
	"kubectl-go-mcp-server/pkg/types"
)

type Server struct {
	kubectlConfig string
	server        *server.MCPServer
	tools         *Tools
	workDir       string
}

func NewServer(kubectlConfig, workDir string) (*Server, error) {
	s := &Server{
		kubectlConfig: kubectlConfig,
		workDir:       workDir, server: server.NewMCPServer(
			"kubectl-go-mcp-server",
			"1.0.0",
			server.WithToolCapabilities(true),
		),
		tools: NewTools(),
	}

	kubectlTool := &kubectl.KubectlTool{}
	s.tools.RegisterTool(kubectlTool)

	for _, tool := range s.tools.AllTools() {
		toolDefn := tool.FunctionDefinition()
		toolInputSchema, err := toolDefn.Parameters.ToRawSchema()
		if err != nil {
			return nil, fmt.Errorf("converting tool schema to json.RawMessage: %w", err)
		}

		s.server.AddTool(mcp.NewToolWithRawSchema(
			toolDefn.Name,
			toolDefn.Description,
			toolInputSchema,
		), s.handleToolCall)
	}

	return s, nil
}

func (s *Server) Serve(ctx context.Context) error {
	return server.ServeStdio(s.server)
}

func (s *Server) GetKubectlConfig() string {
	return s.kubectlConfig
}

func (s *Server) GetWorkDir() string {
	return s.workDir
}

func (s *Server) GetTools() *Tools {
	return s.tools
}

func (s *Server) handleToolCall(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name := request.Params.Name

	argMap, ok := request.Params.Arguments.(map[string]interface{})
	if !ok {
		return mcp.NewToolResultError("Invalid arguments format: expected a map"), nil
	}

	commandVal, ok := argMap["command"]
	if !ok {
		return mcp.NewToolResultError("Missing required parameter: command"), nil
	}
	command, ok := commandVal.(string)
	if !ok {
		return mcp.NewToolResultError("Parameter 'command' must be a string"), nil
	}

	var modifiesResource string
	if modVal, ok := argMap["modifies_resource"]; ok {
		if modStr, ok := modVal.(string); ok {
			modifiesResource = modStr
		}
	}

	log.Printf("Received tool call: tool=%s, command=%s, modifies_resource=%s", name, command, modifiesResource)

	if name != "kubectl" {
		log.Printf("SECURITY WARNING: Attempt to use non-kubectl tool: %s", name)
		return mcp.NewToolResultError(fmt.Sprintf("Only kubectl tool is allowed, tool %s is not permitted", name)), nil
	}

	ctx = context.WithValue(ctx, types.KubeconfigKey, s.kubectlConfig)
	ctx = context.WithValue(ctx, types.WorkdirKey, s.workDir)

	tool := s.tools.Lookup(name)
	if tool == nil {
		return mcp.NewToolResultError(fmt.Sprintf("Tool %s not found", name)), nil
	}

	args := map[string]any{
		"command": command,
	}

	if modifiesResource != "" {
		args["modifies_resource"] = modifiesResource
	}

	output, err := tool.Run(ctx, args)
	if err != nil {
		log.Printf("Error running tool call: %v", err)
		return mcp.NewToolResultError(fmt.Sprintf("Error running tool: %v", err)), nil
	}

	result, err := ToolResultToMap(output)
	if err != nil {
		log.Printf("Error converting tool call output to result: %v", err)
		return mcp.NewToolResultError(fmt.Sprintf("Error processing result: %v", err)), nil
	}

	log.Printf("Tool call output: tool=%s, result=%v", name, result)

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("%v", result),
			},
		},
	}, nil
}

type Tools struct {
	tools map[string]types.Tool
}

func NewTools() *Tools {
	return &Tools{
		tools: make(map[string]types.Tool),
	}
}

func (t *Tools) RegisterTool(tool types.Tool) {
	if _, exists := t.tools[tool.Name()]; exists {
		panic("tool already registered: " + tool.Name())
	}
	t.tools[tool.Name()] = tool
}

func (t *Tools) Lookup(name string) types.Tool {
	return t.tools[name]
}

func (t *Tools) AllTools() []types.Tool {
	var tools []types.Tool
	for _, tool := range t.tools {
		tools = append(tools, tool)
	}
	return tools
}

func (t *Tools) Count() int {
	return len(t.tools)
}

func (t *Tools) HasTool(name string) bool {
	_, exists := t.tools[name]
	return exists
}

func ToolResultToMap(result any) (map[string]interface{}, error) {
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("marshaling result: %w", err)
	}

	var resultMap map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &resultMap); err != nil {
		return map[string]interface{}{
			"output": fmt.Sprintf("%v", result),
		}, nil
	}

	return resultMap, nil
}
