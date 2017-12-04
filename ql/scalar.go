package ql

import (
	"errors"
	"fmt"
	"math"
	"strconv"
)

// Leaf is the interface for all scalars.
type Leaf interface {
	Type
	leaf()
}

// InputField is the interface for all input fields: scalars, enums, input objects
type InputField interface {
	Type
	inputField()
}

// Scalar defines for all scalar types: Int, Float, String, Boolean, ID and
// custom-defined.
type Scalar struct {
	result CoerceFunc
	input  CoerceFunc
}

// NewScalar returns a new Scalar.
func NewScalar(resFunc, inFunc CoerceFunc) *Scalar {
	return &Scalar{
		result: resFunc,
		input:  inFunc,
	}
}

// InputCoerce performs input coercion for Enum.
func (s *Scalar) InputCoerce(value interface{}) interface{} {
	return s.input(value)
}

// ResultCoerce performs result coercion for Enum.
func (s *Scalar) ResultCoerce(value interface{}) interface{} {
	return s.result(value)
}

func (s *Scalar) leaf()       {}
func (s *Scalar) typ()        {}
func (s *Scalar) inputField() {}

// CoerceFunc is the function type for serialize and deserialize.
type CoerceFunc func(value interface{}) interface{}

// Validator is the interface for all types needing validation.
type Validator interface {
	Validate() error
}

// Object represents GraphQL objects.
type Object struct {
	fields []Field
}

// Validate returns an error if Object is invalid.
func (o *Object) Validate() error {
	return errors.New("not implemented")
}

func (o *Object) typ() {}

// InputObject represents GraphQL input objects.
type InputObject struct {
	name   string
	fields []ObjectField
}

// InputCoerce performs input coercion for InputObject.
func (io *InputObject) InputCoerce(value interface{}) interface{} {
	return errors.New("not implemented")
}

// ObjectField represents fields of GraphQL input objects.
type ObjectField struct {
	name  string
	value string // TODO
}

// Interface represents GraphQL interfaces.
type Interface struct {
	fields []Field
}

// Validate returns an error if Interface is invalid.
func (i *Interface) Validate() error {
	return errors.New("not implemented")
}

func (i *Interface) typ() {}

// Field .
type Field struct {
	name string
	args []Argument
	typ  Type
}

// Union represents GraphQL unions.
type Union struct {
	objects []Object
}

// Validate returns an error if Union is invalid.
func (u *Union) Validate() error {
	return errors.New("not implemented")
}

func (u *Union) typ() {}

// Enum represents GraphQL enums.
type Enum struct {
	valueSet map[int]EnumValue
	nameSet  map[string]EnumValue // for O(1) look-up
}

// NewEnum returns an Enum of possible values, vals should be non-nil.
func NewEnum(vals []string) *Enum {
	if len(vals) == 0 {
		return nil
	}
	enum := &Enum{
		valueSet: map[int]EnumValue{},
		nameSet:  map[string]EnumValue{},
	}
	for idx, val := range vals {
		eval := EnumValue{name: val, value: idx}
		enum.valueSet[idx] = eval
		enum.nameSet[val] = eval
	}
	return enum
}

func (e *Enum) typ()        {}
func (e *Enum) inputField() {}

// ResultCoerce performs result coercion for Enum.
func (e *Enum) ResultCoerce(value interface{}) interface{} {
	switch value := value.(type) {
	case int:
		return e.valueSet[value].name
	case *int:
		return e.valueSet[*value].name
	default:
		return fmt.Errorf("cannot coerce %T", value)
	}
}

// InputCoerce performs input coercion for Enum.
func (e *Enum) InputCoerce(value interface{}) interface{} {
	switch value := value.(type) {
	case int:
		return e.valueSet[value]
	case *int:
		return e.valueSet[*value]
	default:
		return fmt.Errorf("cannot coerce %T", value)
	}
}

// EnumValue is the values of some enum type.
type EnumValue struct {
	name  string
	value int
}

// built-in scalars
var (
	Int     = NewScalar(resultInt, inputInt)
	Float   = NewScalar(resultFloat, inputFloat)
	String  = NewScalar(resultString, inputString)
	Boolean = NewScalar(resultBoolean, inputBoolean)
	ID      = NewScalar(resultString, inputID)
)

func inputID(value interface{}) interface{} {
	switch value := value.(type) {
	case int:
		return resultString(value)
	case string:
		return resultString(value)
	case *int:
		return resultString(*value)
	case *string:
		return resultString(*value)
	default:
		return fmt.Errorf("cannot coerce %T", value)
	}
}

func resultBoolean(value interface{}) interface{} {
	switch value := value.(type) {
	case int:
		if value == 0 {
			return false
		}
		return true
	case float32:
		if value != 0 {
			return true
		}
		return false
	case float64:
		if value != 0 {
			return true
		}
		return false
	case string:
		if value == "false" || value == "" {
			return false
		}
		return true
	case bool:
		return value
	case *int:
		return resultBoolean(*value)
	case *float32:
		return resultBoolean(*value)
	case *float64:
		return resultBoolean(*value)
	case *string:
		return resultBoolean(*value)
	case *bool:
		return resultBoolean(*value)
	default:
		return fmt.Errorf("cannot coerce %T", value)
	}
}

func inputBoolean(value interface{}) interface{} {
	switch value := value.(type) {
	case bool:
		return value
	case *bool:
		return *value
	default:
		return fmt.Errorf("cannot coerce %T", value)
	}
}

func resultString(value interface{}) interface{} {
	if val, ok := value.(*string); ok {
		return *val
	}
	return fmt.Sprintf("%v", value)
}

func inputString(value interface{}) interface{} {
	switch value := value.(type) {
	case *string:
		return *value
	case string:
		return value
	default:
		return fmt.Errorf("cannot coerce %T", value)
	}
}

func resultFloat(value interface{}) interface{} {
	switch value := value.(type) {
	case int:
		return float64(value)
	case float32:
		return float64(value)
	case float64:
		return value
	case string:
		val, err := strconv.ParseFloat(value, 0)
		if err != nil {
			return err
		}
		return val
	case bool:
		if value {
			return float64(1)
		}
		return float64(0)
	case *int:
		return float64(*value)
	case *float32:
		return float64(*value)
	case *float64:
		return *value
	case *string:
		return resultFloat(*value)
	case *bool:
		return resultFloat(*value)
	default:
		return fmt.Errorf("cannot coerce %T", value)
	}
}

func inputFloat(value interface{}) interface{} {
	switch value := value.(type) {
	case int:
		return float64(value)
	case float32:
		return float64(value)
	case float64:
		return value
	case *int:
		return float64(*value)
	case *float32:
		return float64(*value)
	case *float64:
		return float64(*value)
	default:
		return fmt.Errorf("cannot coerce %T", value)
	}
}

func resultInt(value interface{}) interface{} {
	switch value := value.(type) {
	case int:
		if value > math.MaxInt32 || value < math.MinInt32 {
			return errors.New("integer out of range(-2^31~2^31)")
		}
		return value
	case float32:
		if value < float32(math.MinInt32) || value > float32(math.MaxInt32) {
			return errors.New("integer out of range(-2^31~2^31)")
		}
		return int(value)
	case float64:
		if value < float64(math.MinInt32) || value > float64(math.MaxInt32) {
			return errors.New("integer out of range(-2^31~2^31)")
		}
		return int(value)
	case string:
		val, err := strconv.ParseFloat(value, 0)
		if err != nil {
			return err
		}
		return resultInt(val)
	case bool:
		if value {
			return 1
		}
		return 0
	case *int:
		return *value
	case *float32:
		return resultInt(*value)
	case *float64:
		return resultInt(*value)
	case *string:
		return resultInt(*value)
	case *bool:
		return resultInt(*value)
	default:
		return fmt.Errorf("cannot coerce %T", value)
	}
}

func inputInt(value interface{}) interface{} {
	switch value := value.(type) {
	case int:
		if value > math.MaxInt32 || value < math.MinInt32 {
			return errors.New("integer out of range(-2^31~2^31)")
		}
		return value
	case *int:
		return *value
	default:
		return fmt.Errorf("cannot coerce %T", value)
	}
}

// Schema is a type-system type parsed from ast.SchemaDefinition.
type Schema struct {
	RootQuery    Object
	RootMutation Object
}

// GroupedFieldSet .
type GroupedFieldSet struct {
	rspKey string
	fields []Field
}

// InputCoercer for input coercion.
type InputCoercer interface {
	InputCoerce(value interface{}) interface{}
}

// ResultCoercer for output coercion
type ResultCoercer interface {
	ResultCoerce(value interface{}) interface{}
}

// // Info for type name and description.
// type Info struct {
// 	name string
// 	desc string
// }
//
// // Meta information for type system.
// type Meta struct {
// 	Info
// 	kind string
// }
//
// // ScalarMeta represents meta information for scalar types.
// type ScalarMeta struct {
// 	Meta
// }
