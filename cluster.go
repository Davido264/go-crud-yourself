package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Pending struct {
	Action   string
	Content  string
	Endpoint string
}

type Server struct {
	Uid       string
	Addresses []string
	Valid     bool
	Online    bool
	Queue     []Pending
	Mut       *sync.RWMutex
}

type Cluster struct {
	servers []Server
}

func (c *Cluster) Notify(req *http.Request, sender string) {
	var content []byte
	if hasJsonBody(req) {
		var err error
		content, err = io.ReadAll(req.Body)
		if err != nil {
			log.Fatalf("Could not read request body: %v\n", err)
		}
	}

	log.Printf("Action: %v", req.Method)
	log.Printf("Body: %v", string(content))
	log.Println("Event Propagation Starts")
	go func() {
		wg := sync.WaitGroup{}

		for i := range c.servers {
			server := c.servers[i]
			wg.Add(1)

			go func() {
				defer wg.Done()

				if server.Uid == sender || !server.Valid || !server.Online {
					if !server.Online {
						log.Printf("Storing queue for %s\n", server.Uid)
						server.Queue = append(server.Queue, Pending{Action: req.Method, Content: string(content)})
					}
					return
				}

				log.Printf("Notifying server %s\n", server.Addresses)

				ctx := context.Background()
				host := normHttp(server.Addresses[0])
				u := fmt.Sprintf("%v%v", host, req.URL.String())
				log.Println(u)
				req2, err := http.NewRequestWithContext(ctx, req.Method, u, strings.NewReader(string(content)))
				if err != nil {
					log.Fatalf("Error creating request: %v", err)
				}

				normalizeRequests(req2, req, content)
				req2.Host = server.Addresses[0]
				req2.URL.Host = server.Addresses[0]
				req2.Header.Set("User-Agent", "middleware")
				req2.Header.Set("X-Middleware-Sent-By", sender)

				resp, err := http.DefaultClient.Do(req2)
				checkResponse(&server, resp, err)

				if !server.Online {
					server.Queue = append(server.Queue, Pending{Action: req.Method, Content: string(content), Endpoint: req.URL.Path})
				}
			}()
			wg.Wait()
		}
	}()
}

func (c *Cluster) GetServerUUID(addr string) *string {
	u, err := url.Parse(normHttp(addr))
	if err != nil {
		log.Fatalln(err.Error())
	}

	for _, server := range c.servers {
		contains := slices.ContainsFunc(server.Addresses, func(value string) bool {
			valurl, err := url.Parse(normHttp(value))
			if err != nil {
				log.Fatalln(err.Error())
			}
			return valurl.Hostname() == u.Hostname()
		})

		if contains {
			return &server.Uid
		}
	}
	return nil
}

func (c *Cluster) BeginWatchman() {
	c.checkServers()
	t := time.NewTicker(time.Duration(30) * time.Second)
	go func() {
		for range t.C {
			c.checkServers()
		}
	}()
}

func (c *Cluster) checkServers() {
	log.Println("Checking server status...")
	wg := sync.WaitGroup{}
	for i := range c.servers {
		wg.Add(1)
		go checkServer(&c.servers[i], &wg)
	}
	wg.Wait()
}

func checkServer(server *Server, wg *sync.WaitGroup) {
	defer wg.Done()

	if server.Online && server.Valid {
		return
	}

	serverhost := normHttp(server.Addresses[0])
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/info", serverhost), nil)

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Printf("Error: %s", err)
		log.Printf("Setting server %s to disconnected state", server.Uid)
		server.Online = false
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Response status: %d", resp.StatusCode)
		log.Printf("Server %s invalid", server.Uid)
		server.Valid = false
		return
	}

	server.Valid = true
	server.Online = true
	log.Printf("Server %s is now online", server.Uid)

	if len(server.Queue) == 0 {
		return
	}

	for _, p := range server.Queue {
		if !server.Online || !server.Valid {
			return
		}

		req, _ := http.NewRequest(p.Action, fmt.Sprintf("http://%s%s", server.Addresses[0], p.Endpoint), strings.NewReader(p.Content))
		resp, err := http.DefaultClient.Do(req)
		checkResponse(server, resp, err)
	}
}

func checkResponse(server *Server, resp *http.Response, err error) {
	if err != nil {
		log.Printf("Error: %v", err)
		log.Printf("Setting server %s to disconnected state", server.Uid)
		server.Online = false
		return
	}

	if (resp.StatusCode >= http.StatusBadRequest && resp.StatusCode < http.StatusInternalServerError) || resp.StatusCode != http.StatusNotImplemented {
		log.Printf("Response status: %d", resp.StatusCode)
		log.Printf("Server %s invalid", server.Uid)
		server.Valid = false
		return
	}

	if resp.StatusCode >= http.StatusInternalServerError {
		log.Printf("Server %s disconnected", server.Uid)
		server.Online = false
		return
	}

}

func hasJsonBody(req *http.Request) bool {
	return req.Header.Get("Content-Type") == "application/json"
}

func normHttp(url string) string {
	if !strings.HasPrefix(url, "http://") || !strings.HasPrefix(url, "https://") {
		return "http://" + url
	}

	return url
}

func NewCluster(serverUrls [][]string) (Cluster, error) {
	servers := make([]Server, len(serverUrls))

	for i, url := range serverUrls {
		servers[i] = Server{
			Uid:       uuid.New().String(),
			Addresses: url,
			Valid:     false,
			Online:    false,
			Queue:     make([]Pending, 0),
			Mut:       &sync.RWMutex{},
		}
	}

	return Cluster{
		servers: servers,
	}, nil
}
