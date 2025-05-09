package commands

import (
	"fmt"
	"log"
	"time"

	connection "github.com/Davido264/go-crud-yourself/pkg/cmdline/utils"
	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
)

var WatchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch logs or events for a specific server",
	Run: func(cmd *cobra.Command, args []string) {
		watch()
	},
}

func watch() {
	conn, err := connection.Websocket()

	if err != nil {
		log.Fatalf("Error connecting to events: %v\n", err)
	}

	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	for {
		t, msg, err := conn.ReadMessage()

		if err != nil {
			log.Fatalf("Erroor reading message: %v\n", err)
		}

		if t != websocket.TextMessage {
			continue
		}
		fmt.Println(string(msg))
	}
}
