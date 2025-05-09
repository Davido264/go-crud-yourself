package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch logs or events for a specific server",
}

var watchLogsCmd = &cobra.Command{
	Use:   "logs <server-alias>",
	Short: "Watch logs for a server",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		server := args[0]
		fmt.Printf("Watching logs for server '%s' at middleware %s...\n", server, host)
	},
}

var watchEventsCmd = &cobra.Command{
	Use:   "events <server-alias>",
	Short: "Watch events for a server",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		server := args[0]
		fmt.Printf("Watching events for server '%s' at middleware %s...\n", server, host)
	},
}

func init() {
	watchCmd.AddCommand(watchLogsCmd)
	watchCmd.AddCommand(watchEventsCmd)
}
