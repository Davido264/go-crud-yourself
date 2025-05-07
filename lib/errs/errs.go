package errs

const (
	ErrnoNoServerRegistered Errt = iota
	ErrnoInvalidField
	ErrnoInvalidFormat
	ErrnoUnknown
	ErrnoNotAllowed
	ErrnoInvalidArgs
	ErrnoInvalidProtocolVersion
	ErrnoMissingData
)

type Errt int

type Err struct {
	errno Errt
}

func (e Err) Error() string {
	switch e.errno {
	case ErrnoNoServerRegistered:
		return "No server registered on the config"
	case ErrnoInvalidField:
		return "Invalid field"
	case ErrnoInvalidFormat:
		return "Invalid format"
	case ErrnoNotAllowed:
		return "Not allowed"
	case ErrnoInvalidArgs:
		return "Invalid arguments"
	case ErrnoInvalidProtocolVersion:
		return "Invalid protocol version"
	case ErrnoMissingData:
		return "No data matching query parameters"
	default:
		return "Unknown error"
	}
}

func Is(err error, t Errt) bool {
	if err == nil {
		return false
	}

	e, ok := err.(Err)

	if !ok {
		return false
	}

	return e.errno == t
}

func New(t Errt) error {
	return Err{errno: t}
}
