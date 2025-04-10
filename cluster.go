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

/**
Profesores
{
  "cedula": "0987654321",
  "nombre": "María González"
}

asignaturas
{
  "id": "matematicas",
  "asignatura": "Matemáticas"
}

profesores_ciclo
{
  "id_profesor": "0987654321",
  "id_asignaturas": "matematicas",
  "ciclo": "2025-A"
}

Estudiantes
{
  "cedula": "1234567890",
  "nombre": "Juan Pérez"
}

Matriculas
{
  "id": "Ma1",
  "cedula_estudiante": "1234567890",
  "id_profesores_ciclo": "M2025L",
  "nota1": 8.5,
  "nota2": 1.0
}
*/

func (c *Cluster) Notify(req *http.Request, sender string) {
	req = req.Clone(context.Background())
	req.URL, _ = url.Parse(req.RequestURI)
	req.RequestURI = ""

	if req.URL.Scheme == "" {
		req.URL.Scheme = "http"
	}

	var content string
	if hasJsonBody(req) {
		contentb, err := io.ReadAll(req.Body)
		if err != nil {
			log.Fatalf("Could not read request body: %v\n", err)
		}
		content = string(contentb)
	}

	go func() {
		wg := sync.WaitGroup{}

		for _, server := range c.servers {
			wg.Add(1)

			go func () {
				defer wg.Done()
				server.Mut.Lock()

				if server.Uid == sender || !server.Valid || !server.Online {
					if !server.Online {
						server.Queue = append(server.Queue, Pending{Action: req.Method, Content: content})
					}
					server.Mut.Unlock()
					return
				}

				server.Mut.Unlock()
				log.Printf("Notifying server %s\n", server.Addresses)

				req.Host = server.Addresses[0]
				req.URL.Host = server.Addresses[0]
				req.Header.Set("User-Agent", "middleware")

				resp, err := http.DefaultClient.Do(req)
				checkResponse(&server, resp, err)

				server.Mut.Lock()
				if !server.Online {
					server.Queue = append(server.Queue, Pending{Action: req.Method, Content: content, Endpoint: req.URL.Path})
				}
				server.Mut.Unlock()
			}()
		}
	}()
}

func (c *Cluster) GetServerUUID(addr string) *string {
	for _, server := range c.servers {
		if slices.Contains(server.Addresses, addr) {
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
	for _, server := range c.servers {
		wg.Add(1)
		go checkServer(server, &wg)
	}
	wg.Wait()
}

func checkServer(server Server, wg *sync.WaitGroup) {
	defer wg.Done()

	server.Mut.Lock()
	defer server.Mut.Unlock()

	if server.Online && server.Valid {
		return
	}

	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s/info", server.Addresses[0]), nil)

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		server.Online = false
		return
	}

	if resp.StatusCode != http.StatusOK {
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
		checkResponse(&server, resp, err)
	}
}

func checkResponse(server *Server, resp *http.Response, err error) {
	server.Mut.Lock()
	defer server.Mut.Unlock()

	if err != nil {
		log.Printf("Server %s disconnected", server.Uid)
		server.Online = false
		return
	}

	if (resp.StatusCode >= http.StatusBadRequest && resp.StatusCode < http.StatusInternalServerError) || resp.StatusCode != http.StatusNotImplemented {
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
