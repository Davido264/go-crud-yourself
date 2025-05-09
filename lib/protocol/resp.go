package protocol

import (
	"encoding/json"

	"github.com/Davido264/go-crud-yourself/lib/errs"
	"github.com/Davido264/go-crud-yourself/lib/logger"
)

const (
	errnoInvalidFieldStr           string = "ERRNO_INVALID_FIELD"
	errnoInvalidFormatStr          string = "ERRNO_INVALID_FORMAT"
	errnoNotAllowedStr             string = "ERRNO_NOT_ALLOWED"
	errnoInvalidArgsStr            string = "ERRNO_INVALID_ARGS"
	errnoUnknownStr                string = "ERRNO_UNKNOWN"
	errnoInvalidProtocolVersionStr string = "ERRNO_INVALID_PROTOCOL_VERSION"
	errnoMissingDataStr            string = "ERRNO_MISSING_DATA_FOR_QUERY"
)

type Response interface {
	SuccessResponse | ErrorResponse
}

type SuccessResponse struct {
	Version int            `json:"version"`
	Data    map[string]any `json:"data"`
}

type ErrorResponse struct {
	Version int    `json:"version"`
	Errno   string `json:"errno"`
}

func ErrnoOf(err error) string {
	if _, ok := err.(*json.SyntaxError); ok {
		return errnoInvalidFormatStr
	}

	switch {
	case errs.Is(err, errs.ErrnoInvalidField):
		return errnoInvalidFieldStr
	case errs.Is(err, errs.ErrnoInvalidFormat):
		return errnoInvalidFormatStr
	case errs.Is(err, errs.ErrnoNotAllowed):
		return errnoNotAllowedStr
	case errs.Is(err, errs.ErrnoInvalidArgs):
		return errnoInvalidArgsStr
	case errs.Is(err, errs.ErrnoInvalidProtocolVersion):
		return errnoInvalidProtocolVersionStr
	case errs.Is(err, errs.ErrnoMissingData):
		return errnoMissingDataStr
	}

	return errnoUnknownStr
}

func Ok(protocolVersion int, data map[string]any) []byte {
	resp, err := json.Marshal(SuccessResponse{
		Version: protocolVersion,
		Data:    data,
	})

	if err != nil {
		logger.Panic(err)
	}

	return resp
}

func Err(protocolVersion int, err error) []byte {
	resp, err := json.Marshal(ErrorResponse{
		Version: protocolVersion,
		Errno:   ErrnoOf(err),
	})

	if err != nil {
		logger.Panic(err)
	}

	return resp
}
