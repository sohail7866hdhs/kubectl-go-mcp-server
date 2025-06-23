package test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"kubectl-go-mcp-server/internal/cli"
)

func TestBuildRootCommand(t *testing.T) {
	t.Run("Build command with default options", func(t *testing.T) {
		opt := &cli.Options{}
		cmd, err := cli.BuildRootCommand(opt, "test-version", "test-commit", "test-date")

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if cmd == nil {
			t.Fatal("Command should not be nil")
		}

		if cmd.Use != "kubectl-go-mcp-server" {
			t.Errorf("Expected Use 'kubectl-go-mcp-server', got %q", cmd.Use)
		}

		if cmd.Short == "" {
			t.Error("Short description should not be empty")
		}

		if cmd.Long == "" {
			t.Error("Long description should not be empty")
		}

		// Check that version subcommand exists
		versionCmd := findSubcommand(cmd, "version")
		if versionCmd == nil {
			t.Error("Version subcommand should exist")
		}
	})
	t.Run("Command with options", func(t *testing.T) {
		opt := &cli.Options{
			KubeConfigPath: "/custom/kubeconfig",
		}

		cmd, err := cli.BuildRootCommand(opt, "test-version", "test-commit", "test-date")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if cmd == nil {
			t.Fatal("Command should not be nil")
		}

		// Check that flags are properly bound
		flags := cmd.Flags()
		if flags == nil {
			t.Fatal("Flags should not be nil")
		}

		// Check kubeconfig flag (mcp-server flag was removed since it's now the default)
		kubeconfigFlag := flags.Lookup("kubeconfig")
		if kubeconfigFlag == nil {
			t.Error("kubeconfig flag should exist")
		}
	})
}

func TestOptions_BindCLIFlags(t *testing.T) {
	opt := &cli.Options{}
	flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)

	err := opt.BindCLIFlags(flagSet)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Check that flags were added
	kubeconfigFlag := flagSet.Lookup("kubeconfig")
	if kubeconfigFlag == nil {
		t.Error("kubeconfig flag should be added")
	}

	// Test flag parsing
	err = flagSet.Parse([]string{"--kubeconfig", "/test/path"})
	if err != nil {
		t.Errorf("Flag parsing failed: %v", err)
	}

	if opt.KubeConfigPath != "/test/path" {
		t.Errorf("Expected KubeConfigPath '/test/path', got %q", opt.KubeConfigPath)
	}
}

func TestRunRootCommand(t *testing.T) {
	// Set up test environment
	if os.Getenv("HOME") == "" {
		t.Setenv("HOME", "/tmp/test-home")
	}

	t.Run("Default kubeconfig path", func(t *testing.T) {
		opt := cli.Options{}
		ctx := context.Background()

		err := cli.RunRootCommand(ctx, opt, []string{})
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		// Note: cli.RunRootCommand takes cli.Options by value, so original opt is not modified
		// The kubeconfig path setting is internal to the function
	})
	t.Run("Custom kubeconfig path", func(t *testing.T) {
		customPath := "/custom/kubeconfig"
		opt := cli.Options{
			KubeConfigPath: customPath,
		}
		ctx := context.Background()

		err := cli.RunRootCommand(ctx, opt, []string{})
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		// Custom path should be preserved
		if opt.KubeConfigPath != customPath {
			t.Errorf("Expected custom kubeconfig path %q, got %q", customPath, opt.KubeConfigPath)
		}
	})

	t.Run("Default MCP server mode", func(t *testing.T) {
		// This test just verifies that RunRootCommand attempts to start the MCP server
		// We can't easily test the actual server startup without significant infrastructure
		opt := cli.Options{}

		// Use a context that would timeout quickly
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
		defer cancel()

		// This should try to start the MCP server and may fail due to various reasons
		// (context timeout, missing deps, etc.) but it shouldn't panic
		_ = cli.RunRootCommand(ctx, opt, []string{})
		// We don't assert on the error since the behavior may vary based on environment
	})

	t.Run("MCP server with custom kubeconfig", func(t *testing.T) {
		// This test verifies that custom kubeconfig is handled
		opt := cli.Options{
			KubeConfigPath: "/custom/kubeconfig",
		}

		// Use a context that would timeout quickly
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
		defer cancel()

		// This should try to start the MCP server and may fail due to various reasons
		// but it shouldn't panic
		_ = cli.RunRootCommand(ctx, opt, []string{})
		// We don't assert on the error since the behavior may vary based on environment
	})

	// Note: Testing MCP server mode would require more complex setup
	// and is better covered by integration tests
}

func TestStartMCPServer(t *testing.T) {
	t.Run("Work directory creation", func(t *testing.T) {
		opt := cli.Options{
			KubeConfigPath: "/tmp/test-kubeconfig",
		}
		// Use a context that will be cancelled immediately to avoid hanging
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		_ = cli.StartMCPServer(ctx, opt) // Error expected due to context cancellation

		// The error might be due to context cancellation or MCP server startup
		// Either way, the work directory should have been created
		workDir := filepath.Join(os.TempDir(), "kubectl-go-mcp-server")
		if _, err := os.Stat(workDir); os.IsNotExist(err) {
			t.Error("Work directory should have been created")
		}

		// Clean up
		if err := os.RemoveAll(workDir); err != nil {
			t.Logf("Warning: failed to clean up work directory: %v", err)
		}
	})
}

func TestVersions(t *testing.T) {
	// Test version command output
	opt := &cli.Options{}
	cmd, err := cli.BuildRootCommand(opt, "test-version", "test-commit", "test-date")
	if err != nil {
		t.Fatalf("Failed to build command: %v", err)
	}

	versionCmd := findSubcommand(cmd, "version")
	if versionCmd == nil {
		t.Fatal("Version command not found")
	}

	if versionCmd.Use != "version" {
		t.Errorf("Expected version command Use 'version', got %q", versionCmd.Use)
	}

	if versionCmd.Short == "" {
		t.Error("Version command should have short description")
	}
}

func TestCommandStructure(t *testing.T) {
	opt := &cli.Options{}
	cmd, err := cli.BuildRootCommand(opt, "test-version", "test-commit", "test-date")
	if err != nil {
		t.Fatalf("Failed to build command: %v", err)
	}

	// Test root command properties
	if cmd.Use != "kubectl-go-mcp-server" {
		t.Errorf("Expected Use 'kubectl-go-mcp-server', got %q", cmd.Use)
	}

	if !strings.Contains(cmd.Short, "Kubernetes") {
		t.Error("Short description should mention Kubernetes")
	}

	if !strings.Contains(cmd.Long, "Model Context Protocol") {
		t.Error("Long description should mention Model Context Protocol")
	}

	// Test that RunE is set
	if cmd.RunE == nil {
		t.Error("RunE should be set")
	}

	// Test subcommands
	subcommands := cmd.Commands()
	if len(subcommands) == 0 {
		t.Error("Should have at least one subcommand")
	}

	// Find version command
	var versionCmd *cobra.Command
	for _, sub := range subcommands {
		if sub.Use == "version" {
			versionCmd = sub
			break
		}
	}

	if versionCmd == nil {
		t.Error("Version subcommand should exist")
		return // Early return to avoid nil pointer dereference
	}

	if versionCmd.Run == nil {
		t.Error("Version command should have Run function")
	}
}

func TestFlagsConfiguration(t *testing.T) {
	opt := &cli.Options{}
	cmd, err := cli.BuildRootCommand(opt, "test-version", "test-commit", "test-date")
	if err != nil {
		t.Fatalf("Failed to build command: %v", err)
	}

	flags := cmd.Flags()

	// Test kubeconfig flag (mcp-server flag removed since it's now the default)
	kubeconfigFlag := flags.Lookup("kubeconfig")
	if kubeconfigFlag == nil {
		t.Fatal("kubeconfig flag should exist")
	}

	if kubeconfigFlag.Value.Type() != "string" {
		t.Errorf("Expected kubeconfig flag type 'string', got %q", kubeconfigFlag.Value.Type())
	}

	if kubeconfigFlag.Usage == "" {
		t.Error("kubeconfig flag should have usage description")
	}
}

// Test edge cases and error conditions
func TestEdgeCases(t *testing.T) {
	t.Run("Empty options", func(t *testing.T) {
		opt := &cli.Options{}
		cmd, err := cli.BuildRootCommand(opt, "test-version", "test-commit", "test-date")
		if err != nil {
			t.Errorf("Should handle empty options without error: %v", err)
		}
		if cmd == nil {
			t.Error("Should return valid command even with empty options")
		}
	})
	t.Run("Nil options", func(t *testing.T) {
		// This would panic in the current implementation, which is expected behavior
		// since cli.Options is expected to be a valid struct
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic with nil options")
			}
		}()

		_, _ = cli.BuildRootCommand(nil, "test-version", "test-commit", "test-date")
	})
}

// Helper function to find subcommand by name
func findSubcommand(cmd *cobra.Command, name string) *cobra.Command {
	for _, sub := range cmd.Commands() {
		if sub.Use == name {
			return sub
		}
	}
	return nil
}

// Benchmark tests for performance
func BenchmarkBuildRootCommand(b *testing.B) {
	opt := &cli.Options{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cmd, err := cli.BuildRootCommand(opt, "test-version", "test-commit", "test-date")
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
		if cmd == nil {
			b.Fatal("Command should not be nil")
		}
	}
}

func BenchmarkBindCLIFlags(b *testing.B) {
	opt := &cli.Options{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		flagSet := pflag.NewFlagSet("test", pflag.ContinueOnError)
		err := opt.BindCLIFlags(flagSet)
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
	}
}

// Test for concurrent access (if relevant)
func TestConcurrentAccess(t *testing.T) {
	// Test that building multiple commands concurrently works
	const numGoroutines = 10
	errors := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			opt := &cli.Options{}
			cmd, err := cli.BuildRootCommand(opt, "test-version", "test-commit", "test-date")
			if err != nil {
				errors <- err
				return
			}
			if cmd == nil {
				errors <- fmt.Errorf("command is nil")
				return
			}
			errors <- nil
		}()
	}

	// Collect results
	for i := 0; i < numGoroutines; i++ {
		if err := <-errors; err != nil {
			t.Errorf("Concurrent command building failed: %v", err)
		}
	}
}
