package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func loadConfig(fileName string) ([][]string, error) {
	bytes, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	servers := make([][]string, 0)
	for server := range strings.Lines(string(bytes)) {
		servers = append(servers, strings.Split(strings.TrimSpace(server), " "))
	}
	return servers, nil
}

func main() {
	fileName := flag.String("config", "servers", "Path to the config file")
	port := flag.Int("port", 8080, "Port to listen on")
	flag.Parse()

	serverUrls, err := loadConfig(*fileName)
	if err != nil {
		log.Fatal(err)
	}

	cluster, err := NewCluster(serverUrls)
	if err != nil {
		log.Fatalf("Could not create cluster: %v\n", err)
	}

	log.Println("Cluster created")

	cluster.BeginWatchman()

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		method := req.Method
		log.Printf("Received request: %s %s\n", method, req.URL.Path)

		if method == http.MethodGet {
			if !strings.HasSuffix(req.URL.Path, "/info") {
				w.WriteHeader(http.StatusMethodNotAllowed)
			} else {
				w.WriteHeader(http.StatusOK)
			}
			return
		}

		if method == "PING" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("PONG"))
			return
		}

		serverUUID := cluster.GetServerUUID(req.RemoteAddr)
		if serverUUID == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		cluster.Notify(req, *serverUUID)
		w.WriteHeader(http.StatusNoContent)
	})

	http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
}
