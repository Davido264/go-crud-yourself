package commands

import (
	"log"

	connection "github.com/Davido264/go-crud-yourself/pkg/cmdline/utils"
	"github.com/spf13/cobra"
)

var StatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check the status of the middleware",
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := connection.Status()

		if err != nil {
			log.Fatalf("Error checking status: %v\n", err)
		}

		var body []byte
		_, _ = resp.Request.Body.Read(body)
		log.Println(string(body))
	},
}
