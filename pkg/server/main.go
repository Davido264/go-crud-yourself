package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Davido264/go-crud-yourself/lib/cluster"
	_ "github.com/Davido264/go-crud-yourself/lib/errs"
	"github.com/gorilla/websocket"
)

var configPaths = []string{"config.json", "/etc/config.json"}

func main() {
	var cfgpath string
	for _, file := range configPaths {
		f, err := os.Stat(file)
		if err != nil || !f.Mode().IsRegular() {
			log.Printf("Unable to use %v\n", file)
			continue
		}
		cfgpath = file
		break
	}

	var cfg cluster.ClusterConfig
	if cfgpath != "" {
		cfg = cluster.ReadConfig(cfgpath)
	} else {
		cfg = cluster.DefaultConfig()
	}

	c := cluster.NewCluster(cfg)

	go c.ListenNotifications()
	go c.ListenEvents()

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	http.HandleFunc("GET /service-integration", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			log.Printf("Cannot upgrade connection. %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = c.Connect(r.RemoteAddr, conn)
		if err != nil {
			log.Printf("Error %v\n", err)
			return
		}
	})

	http.HandleFunc("GET /adm", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			log.Printf("Cannot upgrade connection %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		c.AddManager(conn)
	})

	err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", cfg.Port), nil)

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Cannot start service: %v", err)
	}
}
