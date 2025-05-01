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

func (c *Conn) onPong(appData string) error {
	c.Conn.SetReadDeadline(time.Now().Add(websocketTimeout))
	return nil
}

func (c *Conn) IsClosed(err error) bool {
	return err != nil && (websocket.IsUnexpectedCloseError(err) || websocket.IsCloseError(err))
}

func (c *Conn) SetWriteDeadline() {
	c.Conn.SetWriteDeadline(time.Now().Add(websocketTimeout))
}

func (c *Conn) Ok(data map[string]any) error {
	c.SetWriteDeadline()
	return c.Conn.WriteJSON(protocol.SuccessResponse{
		Version: c.protocolVersion,
		Data:    data,
	})
}

func (c *Conn) Err(err error) error {
	c.SetWriteDeadline()
	return c.Conn.WriteJSON(protocol.ErrorResponse{
		Version: c.protocolVersion,
		Errno:   protocol.ErrnoOf(err),
	})
}

func InitConn(conn *websocket.Conn, version int) *Conn {
	c := &Conn{
		protocolVersion: version,
		Conn:            conn,
		Clientch:        make(chan []byte),
	}

	c.Conn.SetReadLimit(1024)
	c.Conn.SetReadDeadline(time.Now().Add(websocketTimeout))
	c.Conn.SetPongHandler(c.onPong)

	return c
}
