package cluster

import (
	"log"

	"github.com/Davido264/go-crud-yourself/lib/assert"
	"github.com/gorilla/websocket"
)

const mtag = "[MANAGER]"

type ManagerNode struct {
	C    *Conn
	Send <-chan []byte
}

func (m *ManagerNode) Disconnect() {
	assert.Assert(m.C != nil)

	log.Printf("%v Closing connection\n", mtag)
	err := m.C.Conn.Close()
	if err != nil {
		log.Printf("%v Error while closing connection: %v\n", mtag, err)
	}

	close(m.C.Clientch)
	m.C = nil
}

func (m *ManagerNode) Serve() {
	assert.Assert(m.C != nil)

	defer m.Disconnect()

	for msg := range m.C.Clientch {
		m.C.SetWriteDeadline()
		err := m.C.Conn.WriteMessage(websocket.TextMessage, msg)

		if err != nil {
			log.Printf("%v Error on websocket connection: %v\n", mtag, err)
			if !m.C.IsClosed(err) {
				break
			}
		}
	}
}
