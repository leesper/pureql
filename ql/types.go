package ql

import "fmt"

// Type interface for GraphQL eight types.
type Type interface {
	typ()
}

// InputCoercer for input coercion.
type InputCoercer interface {
	InputCoerce(value interface{}) interface{}
}

// ResultCoercer for output coercion
type ResultCoercer interface {
	ResultCoerce(value interface{}) interface{}
}

// Scalar defines for all scalar types: Int, Float, String, Boolean, ID and
// custom-defined.
type Scalar struct {
	result CoerceFunc
	input  CoerceFunc
}

func (s *Scalar) typ() {}

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

// CoerceFunc is the function type for serialize and deserialize.
type CoerceFunc func(value interface{}) interface{}

type Object struct {
	name    string
	ifaces  []*Interface
	directs []*Directive
	fields  []*Field
}

func (o *Object) typ() {}

func (o *Object) Validate() error {
	err := ruleAtLeastOneField(o)
	if err != nil {
		return err
	}
	err = ruleUniqueWithinObject(o)
	if err != nil {
		return err
	}
	err = ruleSuperSetIface(o)
	if err != nil {
		return err
	}
	return nil
}

func ruleAtLeastOneField(o *Object) error {
	if len(o.fields) <= 0 {
		return fmt.Errorf("object %s has no fields", o.name)
	}
	return nil
}

func ruleUniqueWithinObject(o *Object) error {
	set := map[string]bool{}
	for _, f := range o.fields {
		if !set[f.name] {
			set[f.name] = true
		} else {
			return fmt.Errorf("object %s has multiple fields named %s", o.name, f.name)
		}
	}
	return nil
}

func ruleSuperSetIface(o *Object) error {
	for _, iface := range o.ifaces {
		for _, sf := range iface.fields {
			found := false
			for _, f := range o.fields {
				if f.name == sf.name {
					if !isSubType(f, sf) {
						return fmt.Errorf("object %s is not a sub-type of %s", o.name, iface.name)
					}
					if !checkArg(f, sf) {
						return fmt.Errorf("object %s not includes argument defined in %s", o.name, iface.name)
					}
					found = true
				}
			}
			if !found {
				return fmt.Errorf("field %s not found", sf.name)
			}
		}
	}
	return nil
}

func isSubType(field *Field, super *Field) bool {}
func checkArg(field *Field, super *Field) bool  {}

type Field struct {
	name    string
	typ     Type
	directs []*Directive
	args    []*ArgumentDefinition
}

type Directive struct {
	name string
	args []*Argument
}

type Argument struct {
	name string
	val  Value
}

type ArgumentDefinition struct {
	name    string
	typ     Type
	defl    Value
	directs []*Directive
}

type Interface struct {
	name    string
	directs []*Directive
	fields  []*Field
}

func (i *Interface) typ() {}

// // Enum represents GraphQL enums.
// type Enum struct {
// 	valueSet map[int]EnumValue
// 	nameSet  map[string]EnumValue // for O(1) look-up
// }
//
// // NewEnum returns an Enum of possible values, vals should be non-nil.
// func NewEnum(vals []string) *Enum {
// 	if len(vals) == 0 {
// 		return nil
// 	}
// 	enum := &Enum{
// 		valueSet: map[int]EnumValue{},
// 		nameSet:  map[string]EnumValue{},
// 	}
// 	for idx, val := range vals {
// 		eval := EnumValue{name: val, value: idx}
// 		enum.valueSet[idx] = eval
// 		enum.nameSet[val] = eval
// 	}
// 	return enum
// }
//
// // ResultCoerce performs result coercion for Enum.
// func (e *Enum) ResultCoerce(value interface{}) interface{} {
// 	switch value := value.(type) {
// 	case int:
// 		return e.valueSet[value].name
// 	case *int:
// 		return e.valueSet[*value].name
// 	default:
// 		return fmt.Errorf("cannot coerce %T", value)
// 	}
// }
//
// // InputCoerce performs input coercion for Enum.
// func (e *Enum) InputCoerce(value interface{}) interface{} {
// 	switch value := value.(type) {
// 	case int:
// 		return e.valueSet[value]
// 	case *int:
// 		return e.valueSet[*value]
// 	default:
// 		return fmt.Errorf("cannot coerce %T", value)
// 	}
// }
