package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "project",
	Short: "Placeholder for project command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("project command placeholder")
	},
}

func init() {
	// Add command-specific flags here
}
