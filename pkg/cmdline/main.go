package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var host string
	cmd := &cobra.Command{
		Use:   "middlewarectl",
		Short: "Command line tool to interact with middleware",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if host == "" {
				host = os.Getenv("MIDDLEWARE_HOST")
				if host == "" {
					return fmt.Errorf("host not specified via --host flag or MIDDLEWARE_HOST environment variable")
				}
			}
			return nil
		},
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "status",
		Short: "Check the status of the middleware",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Checking status of middleware at %s...\n", host)
			// Simulate action here
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "logs",
		Short: "Fetch recent logs from the middleware",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Fetching logs from middleware at %s...\n", host)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "events",
		Short: "Fetch recent events from the middleware",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Fetching events from middleware at %s...\n", host)
		},
	})

	cmd.Execute()
}
