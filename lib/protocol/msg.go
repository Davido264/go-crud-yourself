package protocol

import (
	"encoding/json"
	"strings"

	"github.com/Davido264/go-crud-yourself/lib/errs"
)

type MsgAction int
type MsgEntity int

const (
	ActionGet MsgAction = iota
	ActionDel
	ActionPut

	actionGetStr string = "get"
	actionDelStr string = "del"
	actionPutStr string = "put"
)

const (
	EntityStudent MsgEntity = iota
	EntityTeacher
	EntityAssigment
	EntityCycle
	EntityMatr
	EntityRNote
	EntityStatus

	entityStudentStr      string = "estudiante"
	entityTeacherStr      string = "profesor"
	entityAssigmentStr    string = "asignatura"
	entityTeacherCycleStr string = "ciclo"
	entityMatrStr         string = "matricula"
	entityRNoteStr        string = "rnota"
	entityStatusStr       string = "status"
)

type Msg struct {
	Version  int            `json:"version"`
	ClientId string         `json:"clientId"`
	Action   MsgAction      `json:"action"`
	Entity   MsgEntity      `json:"entity"`
	Args     map[string]any `json:"args"`
}

func (m *MsgAction) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	switch strings.ToLower(s) {
	case actionGetStr:
		*m = ActionGet
	case actionDelStr:
		*m = ActionDel
	case actionPutStr:
		*m = ActionPut
	default:
		return errs.New(errs.ErrnoInvalidField)
	}

	return nil
}

func (m MsgAction) MarshalJSON() ([]byte, error) {
	var s string
	switch m {
	case ActionGet:
		s = actionGetStr
	case ActionDel:
		s = actionDelStr
	case ActionPut:
		s = actionPutStr
	}

	return json.Marshal(s)
}

func (m *MsgEntity) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	switch strings.ToLower(s) {
	case entityStudentStr:
		*m = EntityStudent
	case entityTeacherStr:
		*m = EntityTeacher
	case entityAssigmentStr:
		*m = EntityAssigment
	case entityTeacherCycleStr:
		*m = EntityCycle
	case entityMatrStr:
		*m = EntityMatr
	case entityRNoteStr:
		*m = EntityRNote
	case entityStatusStr:
		*m = EntityStatus
	default:
		return errs.New(errs.ErrnoInvalidField)
	}

	return nil
}

func (m MsgEntity) MarshalJSON() ([]byte, error) {
	var s string
	switch m {
	case EntityStudent:
		s = entityStudentStr
	case EntityTeacher:
		s = entityTeacherStr
	case EntityAssigment:
		s = entityAssigmentStr
	case EntityCycle:
		s = entityTeacherCycleStr
	case EntityMatr:
		s = entityMatrStr
	case EntityRNote:
		s = entityRNoteStr
	case EntityStatus:
		s = entityStatusStr
	}

	return json.Marshal(s)
}

type TimedMsg struct {
	Msg Msg
	LastTimeStamp int64
}

