package commands

import (
	"github.com/spf13/cobra"
)

var serversCmd = &cobra.Command{
	Use:   "servers",
	Short: "List all servers connected to the middleware",
	Run: func(cmd *cobra.Command, args []string) {
		//
	},
}
