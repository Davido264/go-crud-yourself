package queue

import (
	"slices"
	"sync"
	"time"

	"github.com/Davido264/go-crud-yourself/lib/errs"
	"github.com/Davido264/go-crud-yourself/lib/protocol"
)

type RCMsg struct {
	Msg   protocol.Msg
	count int
}

type MsgQueue struct {
	Mutex         sync.Mutex
	lastTimeStamp int64
	queue         []RCMsg
}

func (q *MsgQueue) LastTimeStamp() int64 {
	return q.lastTimeStamp
}

func (q *MsgQueue) Add(msg protocol.Msg, count int) {
	q.Mutex.Lock()
	defer q.Mutex.Unlock()

	q.lastTimeStamp = msg.LastTimeStamp
	q.queue = append(q.queue, RCMsg{Msg: msg, count: count})
}

func (q *MsgQueue) Pop() *protocol.Msg {
	q.Mutex.Lock()
	defer q.Mutex.Unlock()

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

func (q *MsgQueue) PopSince(timestamp any, entity string) ([]protocol.Msg, error) {
	q.Mutex.Lock()
	defer q.Mutex.Unlock()

	if len(q.queue) == 0 {
		return []protocol.Msg{}, nil
	}

	var idx int
	if timestamp != nil {
		ts, ok := timestamp.(int64)

		if !ok {
			return nil, errs.New(errs.ErrnoInvalidArgs)
		}

		idx = slices.IndexFunc(q.queue, func(m RCMsg) bool {
			return m.Msg.LastTimeStamp <= ts
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
		if msg.Msg.Entity != entity {
			continue
		}
		msg.count--
		result = append(result, protocol.Msg{
			Version:  msg.Msg.Version,
			Entity:   msg.Msg.Entity,
			Action:   msg.Msg.Action,
			ClientId: msg.Msg.ClientId,
			Args:     msg.Msg.Args,
		})
	}

	q.queue = slices.DeleteFunc(q.queue, func(m RCMsg) bool {
		return m.count == 0
	})

	return result, nil
}

func New() MsgQueue {
	return MsgQueue{
		Mutex:         sync.Mutex{},
		lastTimeStamp: time.Now().UTC().UnixMilli(),
		queue:         []RCMsg{},
	}
}
