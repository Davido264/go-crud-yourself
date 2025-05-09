package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/Davido264/go-crud-yourself/lib/assert"
	"github.com/Davido264/go-crud-yourself/lib/cluster"
	_ "github.com/Davido264/go-crud-yourself/lib/errs"
	"github.com/Davido264/go-crud-yourself/lib/logger"
	"github.com/gorilla/websocket"
)

var configPaths = []string{"config.json", "/etc/config.json"}

func main() {
	var cfgpath string
	for _, file := range configPaths {
		f, err := os.Stat(file)
		if err != nil || !f.Mode().IsRegular() {
			logger.Printf("Unable to use %v\n", file)
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
		logger.Printf("Service integration request from %v\n", r.RemoteAddr)
		logger.Printf("ClientId %v\n", r.URL.Query().Get("clientId"))
		id := r.URL.Query().Get("clientId")

		srv := c.GetServer(id)

		if srv == nil {
			logger.Printf("No server registered for %v\n", id)
			w.WriteHeader(http.StatusForbidden)
			return
		}

		if srv.IsConnected() {
			logger.Printf("Server %v already connected\n", id)
			w.WriteHeader(http.StatusForbidden)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			logger.Printf("Cannot upgrade connection. %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		c.ConnectServer(id, conn)
	})

	http.HandleFunc("GET /adm", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			logger.Printf("Cannot upgrade connection %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		c.AddManager(conn)
	})

	http.HandleFunc("GET /adm/servers", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")

		servers, err := json.Marshal(c.ServerList())
		assert.AssertErrNotNil(err)
		w.Write(servers)
	})

	err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", cfg.Port), nil)

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Fatalf("Cannot start service: %v", err)
	}
}
