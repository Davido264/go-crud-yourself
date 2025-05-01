package event

const (
	EServerJoin eventt = iota
	EServerUnjoin
	EServerMsg
)

type eventt int

type Event struct {
	Type   eventt `json:"type"`
	Server string `json:"server"`
}
