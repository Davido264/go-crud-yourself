package cluster

import (
	"log"

	"github.com/Davido264/go-crud-yourself/lib/assert"
	"github.com/Davido264/go-crud-yourself/lib/event"
	"github.com/Davido264/go-crud-yourself/lib/protocol"
	"github.com/gorilla/websocket"
)

type Server struct {
	Identifier string   `json:"-"`
	Alias      string   `json:"alias"`
	Addr       []string `json:"address"`
	C          *Conn    `json:"-"`
}

func (s *Server) DisplayName() string {
	if s.Alias != "" {
		return s.Alias
	}

	return s.Identifier
}

func (s *Server) Disconnect() {
	assert.Assert(s.C != nil)

	log.Printf("[%v] Closing connection\n", s.DisplayName())
	err := s.C.Conn.Close()
	if err != nil {
		log.Printf("[%v] Error closing connection: %v\n", s.DisplayName(), err)
	}

	s.C.Eventch <- event.Event{
		Type:   event.EServerUnjoin,
		Server: s.Identifier,
	}

	close(s.C.Clientch)
	s.C = nil
}

func (s *Server) Listen() {
	assert.Assert(s.C != nil)

	defer s.Disconnect()

	log.Printf("[%v] Listening\n", s.DisplayName())

	s.C.Eventch <- event.Event{
		Type:   event.EServerJoin,
		Server: s.Identifier,
	}

	for {
		var msg protocol.Msg
		err := s.C.Conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("[%v] Error on websocket connection: %v\n", s.DisplayName(), err)
			if !s.C.IsClosed(err) {
				break
			}

			err = s.C.Err(err)
			if err != nil {
				log.Printf("[%v] Error sending response: %v", s.DisplayName(), err)
				if !s.C.IsClosed(err) {
					break
				}
			}
		}

		msg.ClientId = s.Identifier
		s.C.Notifch <- msg

		s.C.Eventch <- event.Event{
			Type:   event.EServerMsg,
			Server: s.Identifier,
		}

		err = s.C.Ok(nil)
		if err != nil {
			log.Printf("[%v] Error on websocket connection: %v\n", s.DisplayName(), err)
			if !s.C.IsClosed(err) {
				break
			}
		}
	}
}

func (s *Server) Serve() {
	assert.Assert(s.C != nil)

	defer s.Disconnect()

	for msg := range s.C.Clientch {
		s.C.SetWriteDeadline()
		err := s.C.Conn.WriteMessage(websocket.TextMessage, msg)

		if err != nil {
			log.Printf("[%v] Error on websocket connection: %v\n", s.DisplayName(), err)
			if !s.C.IsClosed(err) {
				break
			}
		}
	}
}
