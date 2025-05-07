package cluster

import (
	"time"

	"github.com/Davido264/go-crud-yourself/lib/event"
	"github.com/Davido264/go-crud-yourself/lib/protocol"
	"github.com/gorilla/websocket"
)

const (
	websocketTimeout = time.Minute
)

type Conn struct {
	protocolVersion int
	Conn            *websocket.Conn
	Clientch        chan []byte
	Notifch         chan<- protocol.Msg
	Eventch         chan<- event.Event
}

func (c *Conn) IsClosed(err error) bool {
	return err != nil && (websocket.IsUnexpectedCloseError(err) || websocket.IsCloseError(err))
}

func (c *Conn) SetWriteDeadline() {
	c.Conn.SetWriteDeadline(time.Now().Add(websocketTimeout))
}

func InitConn(conn *websocket.Conn, version int, notifch chan<- protocol.Msg, eventch chan<- event.Event) *Conn {
	return &Conn{
		protocolVersion: version,
		Conn:            conn,
		Clientch:        make(chan []byte),
		Notifch:         notifch,
		Eventch:         eventch,
	}
}
