package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "caf",
	Short: "CodeAuthorityFirst - AI session tracking and decision provenance",
	Long: `CodeAuthorityFirst (caf) tracks AI coding sessions, captures decisions, 
and provides provenance for every line of code written by AI agents.

It works with any AI agent (opencode, Claude Code, Cursor, etc.) and provides
a universal pre-commit hook, TUI dashboard, and provenance tracking.`,
	Version: "v0.4.1-mvp.1",
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}