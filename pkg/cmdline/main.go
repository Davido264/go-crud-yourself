package main

import (
	"fmt"
	"os"

	"github.com/Davido264/go-crud-yourself/pkg/cmdline/commands"
	connection "github.com/Davido264/go-crud-yourself/pkg/cmdline/utils"
	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use:   "middlewarectl",
		Short: "Command line tool to interact with middleware",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if connection.Host == "" {
				connection.Host = os.Getenv("MIDDLEWARE_HOST")
				if connection.Host == "" {
					return fmt.Errorf("host not specified via --host flag or MIDDLEWARE_HOST environment variable")
				}
			}
			return nil
		},
	}

	cmd.AddCommand(commands.ServersCmd)
	cmd.AddCommand(commands.StatusCmd)
	cmd.AddCommand(commands.WatchCmd)

	cmd.Execute()
}
