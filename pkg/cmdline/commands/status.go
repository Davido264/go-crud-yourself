package commands

import (
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check the status of the middleware",
	Run: func(cmd *cobra.Command, args []string) {
		// Simulate action here
	},
}
