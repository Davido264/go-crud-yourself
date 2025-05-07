package protocol

import (
	"github.com/Davido264/go-crud-yourself/lib/errs"
)

const (
	fieldId            = "id_"
	fieldIdTeacher     = "id_profesores"
	fieldIdMatricula   = "id_matriculas"
	fieldIdStudents    = "id_estudiantes"
	fieldIdTCA         = "id_profesores_ciclos_asignaturas"
	fieldIdCycles      = "id_ciclos"
	fieldIdAssigments  = "id_asignaturas"
	fieldIdRNotes      = "id_registro_notas"
	fieldName          = "nombre"
	fieldVersion       = "version"
	fieldCycle         = "ciclo"
	fieldAssigment     = "nombre_asignatura"
	fieldNote1         = "nota1"
	fieldNote2         = "nota2"
	fieldSup           = "sup"
	FieldLastTimeStamp = "lastTimeStamp"
)

func ValidateMsg(protocolVersion int, msg Msg) error {
	if msg.Version != protocolVersion {
		return errs.New(errs.ErrnoInvalidProtocolVersion)
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

	if msg.Action == ActionDel && msg.Args[fieldId] == nil {
		return errs.New(errs.ErrnoInvalidArgs)
	}

	if msg.Action == ActionGet && isNaN(msg.Args[FieldLastTimeStamp]) {
		return errs.New(errs.ErrnoInvalidArgs)
	}

	switch msg.Entity {
	case EntityStudent:
		if !validateStudent(msg.Args) {
			return errs.New(errs.ErrnoInvalidArgs)
		}
	case EntityTeacher:
		if !validateTeacher(msg.Args) {
			return errs.New(errs.ErrnoInvalidArgs)
		}
	case EntityAssigment:
		if !validateAssigment(msg.Args) {
			return errs.New(errs.ErrnoInvalidArgs)
		}
	case EntityAsignation:
		if !validateAsignation(msg.Args) {
			return errs.New(errs.ErrnoInvalidArgs)
		}
	case EntityCycle:
		if !validateCycle(msg.Args) {
			return errs.New(errs.ErrnoInvalidArgs)
		}
	case EntityMatr:
		if !validateMatr(msg.Args) {
			return errs.New(errs.ErrnoInvalidArgs)
		}
	case EntityRNote:
		if !validateRnota(msg.Args) {
			return errs.New(errs.ErrnoInvalidArgs)
		}
	}

	return nil
}

func ShouldPropagate(msg Msg) bool {
	return msg.Action != ActionGet
}

func isNaN(a any) bool {
	switch a.(type) {
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64:
		return false
	default:
		return true
	}
}

func validateTeacher(entity map[string]any) bool {
	if entity[fieldId] == nil ||
		entity[fieldIdTeacher] == nil ||
		entity[fieldName] == nil ||
		entity[fieldVersion] == nil {
		return false
	}

	if _, ok := entity[fieldId].(string); !ok {
		return false
	}

	if _, ok := entity[fieldIdTeacher].(string); !ok {
		return false
	}

	if _, ok := entity[fieldName].(string); !ok {
		return false
	}

	if isNaN(entity[fieldVersion]) {
		return false
	}

	return true
}

func validateMatr(entity map[string]any) bool {
	if entity[fieldId] == nil ||
		entity[fieldIdMatricula] == nil ||
		entity[fieldIdStudents] == nil ||
		entity[fieldIdTCA] == nil ||
		entity[fieldVersion] == nil {
		return false
	}

	if _, ok := entity[fieldId].(string); !ok {
		return false
	}

	if _, ok := entity[fieldIdMatricula].(string); !ok {
		return false
	}

	if _, ok := entity[fieldIdStudents].(string); !ok {
		return false
	}

	if _, ok := entity[fieldIdTCA].(string); !ok {
		return false
	}

	if isNaN(entity[fieldVersion]) {
		return false
	}

	return true
}

func validateStudent(entity map[string]any) bool {
	if entity[fieldId] == nil ||
		entity[fieldIdStudents] == nil ||
		entity[fieldName] == nil ||
		entity[fieldVersion] == nil {
		return false
	}

	if _, ok := entity[fieldId].(string); !ok {
		return false
	}

	if _, ok := entity[fieldIdStudents].(string); !ok {
		return false
	}

	if _, ok := entity[fieldName].(string); !ok {
		return false
	}

	if isNaN(entity[fieldVersion]) {
		return false
	}

	return true
}

func validateCycle(entity map[string]any) bool {
	if entity[fieldId] == nil ||
		entity[fieldIdCycles] == nil ||
		entity[fieldCycle] == nil ||
		entity[fieldVersion] == nil {
		return false
	}

	if _, ok := entity[fieldId].(string); !ok {
		return false
	}

	if _, ok := entity[fieldIdCycles].(string); !ok {
		return false
	}

	if _, ok := entity[fieldCycle].(string); !ok {
		return false
	}

	if isNaN(entity[fieldVersion]) {
		return false
	}

	return true

}

func validateAssigment(entity map[string]any) bool {
	if entity[fieldId] == nil ||
		entity[fieldIdAssigments] == nil ||
		entity[fieldAssigment] == nil ||
		entity[fieldVersion] == nil {
		return false
	}

	if _, ok := entity[fieldId].(string); !ok {
		return false
	}

	if _, ok := entity[fieldIdTCA].(string); !ok {
		return false
	}

	if _, ok := entity[fieldAssigment].(string); !ok {
		return false
	}

	if isNaN(entity[fieldVersion]) {
		return false
	}

	return true
}

func validateAsignation(entity map[string]any) bool {
	if entity[fieldId] == nil ||
		entity[fieldIdTCA] == nil ||
		entity[fieldIdTeacher] == nil ||
		entity[fieldIdAssigments] == nil ||
		entity[fieldIdCycles] == nil ||
		entity[fieldVersion] == nil {
		return false
	}

	if _, ok := entity[fieldId].(string); !ok {
		return false
	}

	if _, ok := entity[fieldIdTCA].(string); !ok {
		return false
	}

	if _, ok := entity[fieldIdTeacher].(string); !ok {
		return false
	}

	if _, ok := entity[fieldIdAssigments].(string); !ok {
		return false
	}

	if _, ok := entity[fieldIdCycles].(string); !ok {
		return false
	}

	if isNaN(entity[fieldVersion]) {
		return false
	}

	return true
}

func validateRnota(entity map[string]any) bool {
	if entity[fieldId] == nil ||
		entity[fieldIdRNotes] == nil ||
		entity[fieldIdMatricula] == nil ||
		entity[fieldNote1] == nil ||
		entity[fieldNote2] == nil ||
		entity[fieldSup] == nil ||
		entity[fieldVersion] == nil {
		return false
	}

	if _, ok := entity[fieldId].(string); !ok {
		return false
	}

	if _, ok := entity[fieldIdRNotes].(string); !ok {
		return false
	}

	if _, ok := entity[fieldIdMatricula].(string); !ok {
		return false
	}

	if isNaN(entity[fieldNote1]) {
		return false
	}

	if isNaN(entity[fieldNote2]) {
		return false
	}

	if isNaN(entity[fieldSup]) {
		return false
	}

	if isNaN(entity[fieldVersion]) {
		return false
	}

	return true
}
