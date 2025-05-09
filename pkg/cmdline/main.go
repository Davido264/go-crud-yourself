package main

import (
	"fmt"

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
				return fmt.Errorf("host not specified via --host flag or MIDDLEWARE_HOST environment variable")
			}
			return nil
		},
	}

	cmd.AddCommand(commands.ServersCmd)
	cmd.AddCommand(commands.StatusCmd)
	cmd.AddCommand(commands.WatchCmd)

	cmd.PersistentFlags().StringVar(&connection.Host, "host", "", "host")
	cmd.Execute()
}
