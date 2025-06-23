package cli

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"kubectl-go-mcp-server/internal/config"
	"kubectl-go-mcp-server/internal/mcp"
)

type Options struct {
	KubeConfigPath string `json:"kubeConfigPath,omitempty"`
}

func (o *Options) BindCLIFlags(f *pflag.FlagSet) error {
	f.StringVar(&o.KubeConfigPath, "kubeconfig", o.KubeConfigPath, "path to kubeconfig file")
	return nil
}

func BuildRootCommand(opt *Options, version, commit, date string) (*cobra.Command, error) {
	rootCmd := &cobra.Command{
		Use:   "kubectl-go-mcp-server",
		Short: "Kubernetes MCP Server - Execute kubectl commands via Model Context Protocol",
		Long:  "kubectl-go-mcp-server is a Model Context Protocol (MCP) server that allows language models to interact with your Kubernetes cluster using kubectl commands safely and securely.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunRootCommand(cmd.Context(), *opt, args)
		},
	}
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print the version number of kubectl-go-mcp-server",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("version: %s\ncommit: %s\ndate: %s\n", version, commit, date)
		},
	})

	if err := opt.BindCLIFlags(rootCmd.Flags()); err != nil {
		return nil, err
	}
	return rootCmd, nil
}

func RunRootCommand(ctx context.Context, opt Options, args []string) error {
	// Only validate and expand the kubeconfig path if one was explicitly provided
	if opt.KubeConfigPath != "" {
		// Validate and expand the provided path
		if expanded, err := config.ValidateKubeconfigPath(opt.KubeConfigPath); err == nil {
			opt.KubeConfigPath = expanded
		}
		// If validation fails, keep the original path (let kubectl handle the error)
	}
	// When no kubeconfig is specified, kubectl will use its default behavior (~/.kube/config)

	if err := StartMCPServer(ctx, opt); err != nil {
		return fmt.Errorf("failed to start MCP server: %w", err)
	}
	return nil
}

func StartMCPServer(ctx context.Context, opt Options) error {
	workDir := filepath.Join(os.TempDir(), "kubectl-go-mcp-server")
	if err := os.MkdirAll(workDir, 0o755); err != nil {
		return fmt.Errorf("error creating work directory: %w", err)
	}

	server, err := mcp.NewServer(opt.KubeConfigPath, workDir)
	if err != nil {
		return fmt.Errorf("creating mcp server: %w", err)
	}
	return server.Serve(ctx)
}

func Main(version, commit, date string) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	opt := &Options{}
	cmd, err := BuildRootCommand(opt, version, commit, date)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error building command: %v\n", err)
		os.Exit(1)
	}

	if err := cmd.ExecuteContext(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
