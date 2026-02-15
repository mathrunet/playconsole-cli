package main

import (
	"os"

	"github.com/anthropics/gpc/cmd/gpc/commands"
)

// Version information (set by build)
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	commands.SetVersionInfo(version, commit, date)
	if err := commands.Execute(); err != nil {
		os.Exit(1)
	}
}
