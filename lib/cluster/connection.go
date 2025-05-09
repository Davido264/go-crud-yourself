package cluster

import (
	"errors"
	"sync"
	"syscall"
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
	Mut             sync.Mutex
	Conn            *websocket.Conn
	Clientch        chan []byte
	Notifch         chan<- protocol.Msg
	Eventch         chan<- event.Event
}

func (c *Conn) IsClosedOrBrokenPipe(err error) bool {
	return err != nil && (websocket.IsCloseError(err,
		websocket.CloseGoingAway,
		websocket.CloseAbnormalClosure,
		websocket.CloseNormalClosure) ||
		errors.Is(err, syscall.EPIPE))
}

func (c *Conn) SetWriteDeadline() {
	c.Conn.SetWriteDeadline(time.Now().Add(websocketTimeout))
}

func InitConn(conn *websocket.Conn, version int, chsz int, notifch chan<- protocol.Msg, eventch chan<- event.Event) *Conn {
	return &Conn{
		protocolVersion: version,
		Mut:             sync.Mutex{},
		Conn:            conn,
		Clientch:        make(chan []byte, chsz),
		Notifch:         notifch,
		Eventch:         eventch,
	}
}
