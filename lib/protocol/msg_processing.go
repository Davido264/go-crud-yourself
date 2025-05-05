package protocol

import (
	"time"

	"github.com/Davido264/go-crud-yourself/lib/errs"
)

const (
	fieldId            = "id"
	fieldLastTimeStamp = "lastTimeStamp"
)

func ValidateMsg(protocolVersion int, msg Msg) error {
	if msg.Version != protocolVersion {
		return errs.New(errs.ErrnoInvalidProtocolVersion)
	}

	if msg.ClientId == "" {
		return errs.New(errs.ErrnoInvalidField)
	}

	if msg.Action == ActionDel && (msg.Args == nil || msg.Args[fieldId] == nil) {
		return errs.New(errs.ErrnoInvalidField)
	}

	if msg.Action == ActionPut && msg.Args == nil {
		return errs.New(errs.ErrnoInvalidField)
	}

	if msg.Entity == EntityStatus && msg.Action != ActionGet {
		return errs.New(errs.ErrnoNotAllowed)
	}

	if msg.Action == ActionGet && (msg.Args == nil || msg.Args[fieldLastTimeStamp] == nil) {

	}

	return nil
}

func ShouldPropagate(msg Msg) bool {
	return msg.Action != ActionGet
}

func ProcessMsg(msg Msg) []byte {
	lastTimeStamp := time.Now().UTC().UnixMilli()

	if msg.Action == ActionGet {
		// get entity

		return Ok(msg.Version, map[string]any{
			fieldLastTimeStamp: lastTimeStamp,
		})
	}


	// store entity

	return Ok(msg.Version, map[string]any{
		fieldLastTimeStamp: lastTimeStamp,
	})
}
