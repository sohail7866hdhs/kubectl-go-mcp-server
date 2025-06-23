package main

import (
	"kubectl-go-mcp-server/internal/cli"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cli.Main(version, commit, date)
}
