package protocol

import (
	"encoding/json"
	"log"

	"github.com/Davido264/go-crud-yourself/lib/errs"
)

const (
	errnoInvalidFieldStr  string = "ERRNO_INVALID_FIELD"
	errnoInvalidFormatStr string = "ERRNO_INVALID_FORMAT"
	errnoNotAllowedStr    string = "ERRNO_NOT_ALLOWED"
	errnoInvalidArgsStr   string = "ERRNO_INVALID_ARGS"
	errnoUnknownStr      string = "ERRNO_UNKNOWN"
	errnoInvalidProtocolVersionStr string = "ERRNO_INVALID_PROTOCOL_VERSION"
	errnoMissingDataStr string = "ERRNO_MISSING_DATA_FOR_QUERY"
)

type Response interface {
	SuccessResponse | ErrorResponse
}

type SuccessResponse struct {
	Version int            `json:"version"`
	Data    map[string]any `json:"data"`
}

type errno int
type ErrorResponse struct {
	Version int   `json:"version"`
	Errno   errno `json:"errno"`
}

func (e errno) MarshalJSON() ([]byte, error) {
	var s string
	switch errs.Errt(e) {
	case errs.ErrnoInvalidField:
		s = errnoInvalidFieldStr
	case errs.ErrnoInvalidFormat:
		s = errnoInvalidFormatStr
	case errs.ErrnoNotAllowed:
		s = errnoNotAllowedStr
	case errs.ErrnoInvalidArgs:
		s = errnoInvalidArgsStr
	case errs.ErrnoInvalidProtocolVersion:
		s = errnoInvalidProtocolVersionStr
	case errs.ErrnoMissingData:
		s = errnoMissingDataStr
	default:
		s = ""
	}

	return json.Marshal(s)
}

func ErrnoOf(err error) errno {
	if _, ok := err.(*json.SyntaxError); ok {
		return errno(errs.ErrnoInvalidFormat)
	}

	switch {
	case errs.Is(err, errs.ErrnoInvalidField):
		return errno(errs.ErrnoInvalidField)
	case errs.Is(err, errs.ErrnoInvalidFormat):
		return errno(errs.ErrnoInvalidFormat)
	case errs.Is(err, errs.ErrnoNotAllowed):
		return errno(errs.ErrnoNotAllowed)
	case errs.Is(err, errs.ErrnoInvalidArgs):
		return errno(errs.ErrnoInvalidArgs)
	case errs.Is(err, errs.ErrnoInvalidProtocolVersion):
		return errno(errs.ErrnoInvalidProtocolVersion)
	case errs.Is(err, errs.ErrnoMissingData):
		return errno(errs.ErrnoMissingData)
	}

	return errno(errs.ErrnoUnknown)
}

func Ok(protocolVersion int, data map[string]any) []byte {
	resp, err := json.Marshal(SuccessResponse{
		Version: protocolVersion,
		Data:    data,
	})

	if err != nil {
		log.Panic(err)
	}

	return resp
}

func Err(protocolVersion int, err error) []byte {
	resp, err := json.Marshal(ErrorResponse{
		Version: protocolVersion,
		Errno:   ErrnoOf(err),
	})

	if err != nil {
		log.Panic(err)
	}

	return resp
}
