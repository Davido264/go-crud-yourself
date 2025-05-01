package cluster

import (
	"encoding/json"
	"log"
	"slices"

	"github.com/Davido264/go-crud-yourself/lib/assert"
	"github.com/Davido264/go-crud-yourself/lib/errs"
	"github.com/Davido264/go-crud-yourself/lib/event"
	"github.com/Davido264/go-crud-yourself/lib/protocol"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const clustermtag = "[CLUSTER MANAGER]"

type Cluster struct {
	protocolVersion int
	managers        []ManagerNode
	servers         map[string]*Server
	notifych        chan protocol.Msg
	eventch         chan event.Event
}

func (c *Cluster) Connect(addr string, conn *websocket.Conn) error {
	id := c.findRegistered(addr)

	if id == "" {
		return errs.New(errs.ErrnoNoServerRegistered)
	}

	server := c.servers[id]
	server.C = InitConn(conn, c.protocolVersion)

	go server.Listen()
	go server.Serve()
	return nil
}

func (c *Cluster) AddManager(conn *websocket.Conn) {
	nManager := ManagerNode{
		C: InitConn(conn, c.protocolVersion),
	}

	c.managers = append(c.managers, nManager)
	go nManager.Serve()
}

func (c *Cluster) ListenNotifications() {
	for msg := range c.notifych {
		c.NotifiyServers(msg)
	}
}

func (c *Cluster) ListenEvents() {
	for ev := range c.eventch {
		c.NotifyManagers(ev)
	}
}

func (c *Cluster) NotifiyServers(msg protocol.Msg) {
	log.Printf("%v Notifying servers\n", clustermtag)
	encmsg, err := json.Marshal(msg)
	if err != nil {
		log.Panic(err)
	}

	for i := range c.servers {
		if c.servers[i].Identifier == msg.ClientId {
			continue
		}

		c.servers[i].C.Clientch <- encmsg
	}
}

func (c *Cluster) NotifyManagers(ev event.Event) {
	log.Printf("%v Notifying connected manager clients\n", clustermtag)
	encev, err := json.Marshal(ev)
	assert.AssertErrNotNil(err)

	for i := range c.managers {
		c.managers[i].C.Clientch <- encev
	}
}

func (c *Cluster) IsValidMsg(msg protocol.Msg) bool { return true }

func (c *Cluster) findRegistered(addr string) string {
	for i, server := range c.servers {
		if slices.Contains(server.Addr, addr) {
			return i
		}
	}

	return ""
}

func NewCluster(cfg ClusterConfig) *Cluster {
	m := make(map[string]*Server)

	for _, server := range cfg.Servers {
		id := uuid.Must(uuid.NewUUID()).String()
		server.Identifier = id
		m[id] = &server
	}

	return &Cluster{
		servers:         m,
		protocolVersion: cfg.ProtocolVersion,
		managers:        []ManagerNode{},
		notifych:        make(chan protocol.Msg, cfg.ChannelSize),
		eventch:         make(chan event.Event, cfg.ChannelSize),
	}
}
