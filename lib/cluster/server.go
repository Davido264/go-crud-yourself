package cluster

import (
	"context"
	"sync"
	"time"

	"github.com/Davido264/go-crud-yourself/lib/assert"
	"github.com/Davido264/go-crud-yourself/lib/event"
	"github.com/Davido264/go-crud-yourself/lib/logger"
	"github.com/Davido264/go-crud-yourself/lib/protocol"
	"github.com/gorilla/websocket"
)

type Server struct {
	Identifier string     `json:"token"`
	Alias      string     `json:"alias"`
	Address    string     `json:"address"`
	C          *Conn      `json:"-"`
	Mut        sync.Mutex `json:"-"`
}

func (s *Server) IsConnected() bool {
	s.Mut.Lock()
	defer s.Mut.Unlock()
	return s.C != nil
}

func (s *Server) DisplayName() string {
	if s.Alias != "" {
		return s.Alias
	}

	return s.Identifier
}

func (s *Server) Disconnect() {
	s.Mut.Lock()
	defer s.Mut.Unlock()

	if s.C == nil {
		logger.Printf("[%v] Already disconnected\n", s.DisplayName())
		return
	}

	logger.Printf("[%v] Closing connection\n", s.DisplayName())
	err := s.C.Conn.Close()
	if err != nil && s.C.IsClosedOrBrokenPipe(err) {
		logger.Printf("[%v] Error closing connection: %v\n", s.DisplayName(), err)
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

	ctx, cancel := context.WithCancel(context.Background())

	wg := sync.WaitGroup{}
	wg.Add(2)
	go s.listen(cancel, &wg)
	go s.serve(ctx, cancel, &wg)

	wg.Wait()
}

func (s *Server) listen(cancel context.CancelFunc, wg *sync.WaitGroup) {
	defer wg.Done()

	logger.Printf("[%v] Listening\n", s.DisplayName())

	s.C.Eventch <- event.Event{
		Type:   event.EServerJoin,
		Server: s.Identifier,
	}

	for {
		var msg protocol.Msg
		err := s.C.Conn.ReadJSON(&msg)

		if err != nil {
			if !s.C.IsClosedOrBrokenPipe(err) {
				logger.Printf("[%v] Error on websocket connection: %v\n", s.DisplayName(), err)
			}
			cancel()
			return
		}

		err = protocol.ValidateMsg(s.C.protocolVersion, msg)
		if err != nil {
			logger.Printf("[%v] Invalid message: %v\n", s.DisplayName(), err)
			s.C.Clientch <- protocol.Err(s.C.protocolVersion, err)
			continue
		}

		msg.LastTimeStamp = time.Now().UTC().UnixMilli()
		msg.ClientId = s.Identifier

		s.C.Eventch <- event.Event{
			Type:   event.EServerMsg,
			Server: s.Identifier,
			Msg:    &msg,
		}

		s.C.Notifch <- msg

		if msg.Action != protocol.ActionGet {
			s.C.Clientch <- protocol.Ok(msg.Version, map[string]any{
				protocol.FieldLastTimeStamp: msg.LastTimeStamp,
			})
		}
	}
}

func (s *Server) serve(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup) {
	defer wg.Done()

	logger.Printf("[%v] Serving\n", s.DisplayName())
	for {
		select {
		case msg := <-s.C.Clientch:
			{
				logger.Printf("[%v] Sending message to client\n", s.DisplayName())
				s.C.Mut.Lock()
				s.C.SetWriteDeadline()
				err := s.C.Conn.WriteMessage(websocket.TextMessage, msg)
				s.C.Mut.Unlock()

				if err != nil {
					if !s.C.IsClosedOrBrokenPipe(err) {
						logger.Printf("[%v] Error on websocket connection: %v\n", s.DisplayName(), err)
					}
					cancel()
					return
				}

				logger.Printf("[%v] Message sent to client\n", s.DisplayName())
			}
		case <-ctx.Done():
			return
		}
	}
}
