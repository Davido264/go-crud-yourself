package event

import "github.com/Davido264/go-crud-yourself/lib/protocol"

const (
	EServerJoin eventt = iota
	EServerUnjoin
	EServerMsg
)

type eventt int

type Event struct {
	Type   eventt        `json:"type"`
	Server string        `json:"server"`
	Msg    *protocol.Msg `json:"msg"`
}
