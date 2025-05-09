package commands

import (
	"fmt"
	"io"
	"log"

	connection "github.com/Davido264/go-crud-yourself/pkg/cmdline/utils"
	"github.com/spf13/cobra"
)

var ServersCmd = &cobra.Command{
	Use:   "servers",
	Short: "List all servers connected to the middleware",
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := connection.Servers()

		if err != nil {
			log.Fatalf("Error getting servers: %v\n", err)
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Error reading response body: %v\n", err)
		}

		fmt.Println(string(body))
	},
}
