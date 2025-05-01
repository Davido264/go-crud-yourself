package protocol

import (
	"encoding/json"

	"github.com/Davido264/go-crud-yourself/lib/errs"
)

const (
	errnoInvalidFieldStr  string = "ERRNO_INVALID_FIELD"
	errnoInvalidFormatStr string = "ERRNO_INVALID_FORMAT"
	errnoNotAllowedStr    string = "ERRNO_NOT_ALLOWED"
	errnoInvalidArgsStr   string = "ERRNO_INVALID_ARGS"
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

func (e errno) MarshallJSON() ([]byte, error) {
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
	default:
		s = ""
	}

	return json.Marshal(s)
}

func ErrnoOf(err error) errno {
	if _, ok := err.(*json.SyntaxError); ok {
		return errno(errs.ErrnoInvalidFormat)
	}

	return errno(errs.ErrnoUnknown)
}
