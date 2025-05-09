package commands

import (
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

		var body []byte
		_, _ = resp.Request.Body.Read(body)
		log.Println(string(body))
	},
}
