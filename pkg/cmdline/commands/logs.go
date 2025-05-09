package commands

import (
	"github.com/spf13/cobra"
)

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Fetch recent logs from the middleware",
	Run: func(cmd *cobra.Command, args []string) {
	},
}
