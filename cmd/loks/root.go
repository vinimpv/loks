package loks

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "loks",
	Short: "Loks CLI tool",
}

// Execute starts the CLI.
func Execute() {
	rootCmd.Execute()
}
