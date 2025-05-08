package cluster

import (
	"log"
	"sync"
	"time"

	"github.com/Davido264/go-crud-yourself/lib/assert"
	"github.com/Davido264/go-crud-yourself/lib/event"
	"github.com/Davido264/go-crud-yourself/lib/protocol"
	"github.com/gorilla/websocket"
)

type Server struct {
	Identifier string `json:"token"`
	Alias      string `json:"alias"`
	C          *Conn  `json:"-"`
}

func (s *Server) DisplayName() string {
	if s.Alias != "" {
		return s.Alias
	}

	return s.Identifier
}

func (s *Server) Disconnect() {

	if s.C == nil {
		log.Printf("[%v] Already disconnected\n", s.DisplayName())
		return
	}

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

func (s *Server) ListenAndServe() {
	assert.Assert(s.C != nil)

	defer s.Disconnect()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go s.listen(&wg)
	go s.serve(&wg)

	wg.Wait()
}

func (s *Server) listen(wg *sync.WaitGroup) {
	defer wg.Done()

	log.Printf("[%v] Listening\n", s.DisplayName())

	s.C.Eventch <- event.Event{
		Type:   event.EServerJoin,
		Server: s.Identifier,
	}

	for {
		log.Printf("[%v] Received message from client\n", s.DisplayName())
		var msg protocol.Msg
		err := s.C.Conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("[%v] Error on websocket connection: %v\n", s.DisplayName(), err)
			break
		}

		err = protocol.ValidateMsg(s.C.protocolVersion, msg)
		if err != nil {
			log.Printf("[%v] Invalid message: %v\n", s.DisplayName(), err)
			s.C.Clientch <- protocol.Err(s.C.protocolVersion, err)
			continue
		}

		msg.LastTimeStamp = time.Now().UTC().UnixMilli()
		msg.ClientId = s.Identifier

		s.C.Eventch <- event.Event{
			Type:   event.EServerMsg,
			Server: s.Identifier,
		}

		if msg.Action == protocol.ActionGet {
			s.C.Notifch <- msg
			continue
		}

		s.C.Clientch <- protocol.Ok(msg.Version, map[string]any{
			protocol.FieldLastTimeStamp: msg.LastTimeStamp,
		})
	}
}

func (s *Server) serve(wg *sync.WaitGroup) {
	defer wg.Done()

	log.Printf("[%v] Serving\n", s.DisplayName())
	for msg := range s.C.Clientch {
		log.Printf("[%v] Sending message to client\n", s.DisplayName())
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
