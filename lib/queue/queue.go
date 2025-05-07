package queue

import (
	"github.com/Davido264/go-crud-yourself/lib/protocol"
)

type MsgQueue struct {
	lastTimeStamp int64
	queue         []protocol.Msg
}

func (q *MsgQueue) Add(msg protocol.Msg) {
	q.lastTimeStamp = msg.LastTimeStamp
	q.queue = append(q.queue, msg)
}

func (q *MsgQueue) Pop() *protocol.Msg {
	if len(q.queue) == 0 {
		return nil
	}
	msg := q.queue[0]
	q.queue = q.queue[1:]
	return &msg
}
