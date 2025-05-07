package queue

import (
	"slices"
	"time"

	"github.com/Davido264/go-crud-yourself/lib/errs"
	"github.com/Davido264/go-crud-yourself/lib/protocol"
)

type RCMsg struct {
	Msg   protocol.Msg
	count int
}

type MsgQueue struct {
	lastTimeStamp int64
	queue         []RCMsg
}

func (q *MsgQueue) LastTimeStamp() int64 {
	return q.lastTimeStamp
}

func (q *MsgQueue) Add(msg protocol.Msg, count int) {
	q.lastTimeStamp = msg.LastTimeStamp
	q.queue = append(q.queue, RCMsg{Msg: msg, count: count})
}

func (q *MsgQueue) Pop() *protocol.Msg {
	if len(q.queue) == 0 {
		return nil
	}
	msg := q.queue[0]
	msg.count--
	if msg.count == 0 {
		q.queue = q.queue[1:]
	}

	if len(q.queue) == 0 {
		q.lastTimeStamp = time.Now().UTC().UnixMilli()
	}

	return &msg.Msg
}

func (q *MsgQueue) PopSince(timestamp *int64) ([]protocol.Msg, error) {
	if len(q.queue) == 0 {
		return []protocol.Msg{}, nil
	}

	var idx int
	if timestamp != nil {
		idx = slices.IndexFunc(q.queue, func(m RCMsg) bool {
			return m.Msg.LastTimeStamp <= *timestamp
		})

		if idx == -1 {
			return nil, errs.New(errs.ErrnoMissingData)
		}
	} else {
		idx = len(q.queue) - 1
	}

	var result []protocol.Msg
	for i := range q.queue[:idx] {
		msg := q.queue[i]
		msg.count--
		result = append(result, protocol.Msg{
			Version:  msg.Msg.Version,
			Entity:   msg.Msg.Entity,
			Action:   msg.Msg.Action,
			ClientId: msg.Msg.ClientId,
			Args:     msg.Msg.Args,
		})
	}



	return result, nil
}
