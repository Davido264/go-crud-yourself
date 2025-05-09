package cluster

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Davido264/go-crud-yourself/lib/event"
	"github.com/Davido264/go-crud-yourself/lib/logger"
	"github.com/Davido264/go-crud-yourself/lib/protocol"
	"github.com/Davido264/go-crud-yourself/lib/queue"
	"github.com/gorilla/websocket"
)

const clustermtag = "[CLUSTER MANAGER]"

type Cluster struct {
	protocolVersion int
	chsz            int

	managers []ManagerNode
	servers  map[string]*Server

	notifych chan protocol.Msg
	msgqueue queue.MsgQueue

	eventch chan event.Event
	logch   chan []byte
}

func (c *Cluster) ConnectServer(id string, conn *websocket.Conn) {
	server := c.servers[id]
	server.C = InitConn(conn, c.protocolVersion, c.chsz, c.notifych, c.eventch)

	go server.ListenAndServe()
}

func (c *Cluster) AddManager(conn *websocket.Conn) {
	nManager := ManagerNode{
		C: InitConn(conn, c.protocolVersion, c.chsz, c.notifych, c.eventch),
	}

	c.managers = append(c.managers, nManager)
	go nManager.Serve()
}

func (c *Cluster) ListenNotifications() {
	for msg := range c.notifych {

		if msg.Action == protocol.ActionGet {
			logger.Printf("%v Procesing request: %v\n", clustermtag, msg)

			data := make(map[string]any)
			switch msg.Entity {
			case protocol.EntityStatus:
				data[protocol.FieldLastTimeStamp] = c.msgqueue.LastTimeStamp()
			default:
				data[protocol.FieldLastTimeStamp] = c.msgqueue.LastTimeStamp()
				actions, err := c.msgqueue.PopSince(data[protocol.FieldLastTimeStamp], msg.Entity)
				if err != nil {
					logger.Printf("%v Error: %v\n", clustermtag, err)
					c.servers[msg.ClientId].C.Clientch <- protocol.Err(c.protocolVersion, err)
					continue
				}
				data["actions"] = actions
			}

			if c.servers[msg.ClientId].C != nil {
				c.servers[msg.ClientId].C.Clientch <- protocol.Ok(c.protocolVersion, data)
			}

			continue
		}

		logger.Printf("%v Propagating message: %v\n", clustermtag, msg)
		c.NotifiyServers(msg)
	}
}

func (c *Cluster) ListenEvents() {
	for _ = range c.eventch {
		// jsonEv, err := json.Marshal(ev)
		// assert.AssertErrNotNil(err)

		// c.NotifyManagers(jsonEv)
	}
}

func (c *Cluster) ListenLogs() {
	for log := range c.logch {
		c.NotifyManagers(log)
	}
}

func (c *Cluster) NotifiyServers(msg protocol.Msg) {
	logger.Printf("%v Notifying servers\n", clustermtag)
	encmsg, err := json.Marshal(msg)
	if err != nil {
		logger.Panic(err)
	}

	missingCount := 0
	for i := range c.servers {
		if c.servers[i].Identifier == msg.ClientId {
			continue
		}

		if c.servers[i].Address != "" {
			msg.ClientId = c.servers[i].Identifier
		}

		if c.servers[i].C == nil {
			logger.Printf("%v Server %v is not connected. Skiping...\n", clustermtag, c.servers[i].DisplayName())
			missingCount++
			continue
		}
		c.servers[i].C.Clientch <- encmsg
	}

	if missingCount != 0 {
		c.msgqueue.Add(msg, missingCount)
	}
}

func (c *Cluster) NotifyManagers(ev []byte) {
	logger.LocalOnlyPrintf("%v Notifying connected manager clients\n", clustermtag)

	for i := range c.managers {
		c.managers[i].C.Clientch <- ev
	}
}

func (c *Cluster) JoinClusters() {
	for i := range c.servers {
		if c.servers[i].Address == "" {
			continue
		}

		logger.Printf("%v Joining cluster %v\n", clustermtag, c.servers[i].Alias)

		u := url.URL{
			Scheme:   "ws",
			Host:     c.servers[i].Address,
			Path:     "/service-integration",
			RawQuery: fmt.Sprintf("clientId=%s", c.servers[i].Identifier),
		}

		conn, res, err := websocket.DefaultDialer.Dial(u.String(), nil)

		if err != nil {
			logger.Printf("%v Error joining cluster %v: %v\n", clustermtag, c.servers[i].Alias, err)
			continue
		}

		if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusSwitchingProtocols {
			logger.Printf("%v Error joining cluster %v: HTTP Status %v\n", clustermtag, c.servers[i].Alias, res.Status)
			continue
		}

		c.servers[i].C = InitConn(conn, c.protocolVersion, c.chsz, c.notifych, c.eventch)
	}
}

func (c *Cluster) GetServer(id string) *Server {
	if server, ok := c.servers[id]; ok {
		return server
	}
	return nil
}

func (c *Cluster) ServerList() []Server {
	servers := make([]Server, 0, len(c.servers))
	for _, server := range c.servers {
		servers = append(servers, *server)
	}
	return servers
}

func (c *Cluster) Status() map[string]any {
	var connected int
	for i := range c.servers {
		if c.servers[i].C != nil {
			connected++
		}
	}

	return map[string]any{
		"onlineServers": connected,
		"lastTimeStamp": c.msgqueue.LastTimeStamp(),
		"queueSize":     c.msgqueue.Len(),
	}
}

func NewCluster(cfg ClusterConfig) *Cluster {
	m := make(map[string]*Server)

	for i := range cfg.Servers {
		id := cfg.Servers[i].Identifier
		m[id] = &cfg.Servers[i]
	}

	c := &Cluster{
		servers:         m,
		protocolVersion: cfg.ProtocolVersion,
		msgqueue:        queue.New(),
		managers:        []ManagerNode{},
		notifych:        make(chan protocol.Msg, cfg.ChannelSize),
		eventch:         make(chan event.Event, cfg.ChannelSize),
		logch:           make(chan []byte, cfg.ChannelSize),
	}

	logger.RegisterLogger(c.logch)
	return c
}
