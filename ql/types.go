package ql

import (
	"fmt"
	"reflect"
)

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

// Fielder is a type which has fields.
type Fielder interface {
	Fields() []*Field
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

// Object represents a list of named fields.
type Object struct {
	name    string
	ifaces  []*Interface
	directs []*Directive
	fields  []*Field
}

func (obj *Object) typ() {}

// Fields returns fields of Object.
func (obj *Object) Fields() []*Field {
	return obj.fields
}

// Validate checks whether Object adhered to a set of rules.
func (obj *Object) Validate() error {
	err := ruleAtLeastOneField(obj)
	if err != nil {
		return err
	}
	err = ruleUniqueWithinObject(obj)
	if err != nil {
		return err
	}
	err = ruleSuperSetIface(obj)
	if err != nil {
		return err
	}
	return nil
}

func ruleAtLeastOneField(f Fielder) error {
	if len(f.Fields()) <= 0 {
		return fmt.Errorf("type %T has no fields", f)
	}
	return nil
}

func ruleUniqueWithinObject(f Fielder) error {
	set := map[string]bool{}
	for _, field := range f.Fields() {
		if !set[field.name] {
			set[field.name] = true
		} else {
			return fmt.Errorf("type %T has multiple fields named %s", f, field.name)
		}
	}
	return nil
}

func ruleSuperSetIface(obj *Object) error {
	for _, iface := range obj.ifaces {
		if err := ruleIncludeField(obj, iface); err != nil {
			return err
		}
	}
	return nil
}

func ruleIncludeField(obj *Object, iface *Interface) error {
	fieldMap := map[string]*Field{}
	for _, f := range obj.fields {
		fieldMap[f.name] = f
	}

	for _, f := range iface.fields {
		if fieldMap[f.name] == nil {
			return fmt.Errorf("object %s has no field %s of interface %s", obj.name, f.name, iface.name)
		}
		err := ruleSubType(fieldMap[f.name].typ, f.typ)
		if err != nil {
			return err
		}
		err = ruleArgSameName(fieldMap[f.name].args, f.args)
		if err != nil {
			return err
		}
	}
	return nil
}

func ruleSubType(typ Type, super Type) error {
	if reflect.DeepEqual(typ, super) {
		return nil
	}

	obj, isObject := typ.(*Object)
	iface, isIface := super.(*Interface)
	if isObject && isIface {
		return ruleIncludeField(obj, iface)
	}

	union, isUnion := super.(*Union)
	if isObject && isUnion {
		for _, typ := range union.typs {
			if err := ruleSubType(obj, typ); err == nil {
				return nil
			}
		}
		return fmt.Errorf("type %T is not a sub-type of %T", typ, super)
	}

	lto, isListObject := typ.(*ListType)
	lti, isListIface := super.(*ListType)
	if isListObject && isListIface {
		return ruleSubType(lto.underline, lti.underline)
	}

	nt, isNonNull := typ.(*NonNullType)
	if isNonNull {
		return ruleSubType(nt.underline, super)
	}

	return nil
}

func ruleArgSameName(args []*ArgumentDefinition, iargs []*ArgumentDefinition) error {
	argMap := map[string]*ArgumentDefinition{}
	for _, a := range args {
		argMap[a.name] = a
	}

	for _, a := range iargs {
		if argMap[a.name] == nil {
			return fmt.Errorf("no argument definition found for %s", a.name)
		}
		if !reflect.DeepEqual(argMap[a.name].typ, a.typ) {
			return fmt.Errorf("expected argument type %T, found %T", a.typ, argMap[a.name].typ)
		}
	}
	return nil
}

// Field represents fields in Object and Interface.
type Field struct {
	name    string
	typ     Type
	directs []*Directive
	args    []*ArgumentDefinition
}

// Directive decorating other entities.
type Directive struct {
	name string
	args []*Argument
}

// Argument defines arguments in other entities.
type Argument struct {
	name string
	val  Value
}

// ArgumentDefinition represents argument definitions in Object and Interface.
type ArgumentDefinition struct {
	name    string
	typ     Type
	defl    Value
	directs []*Directive
}

// Interface represents a list of fields and their arguments.
type Interface struct {
	name    string
	directs []*Directive
	fields  []*Field
}

func (iface *Interface) typ() {}

// Fields returns fields of Interface.
func (iface *Interface) Fields() []*Field {
	return iface.fields
}

// Validate checks whether Interface adhered to a set of rules.
func (iface *Interface) Validate() error {
	err := ruleAtLeastOneField(iface)
	if err != nil {
		return err
	}

	err = ruleUniqueWithinObject(iface)
	if err != nil {
		return err
	}

	return nil
}

// Union represents an object that could be one of a list of Object types.
type Union struct {
	name    string
	directs []*Directive
	typs    []Type
}

func (uni *Union) typ() {}

// Validate checks whether Union adhered to a set of rules.
func (uni *Union) Validate() error {
	err := ruleAtLeastOneType(uni)
	if err != nil {
		return err
	}
	err = ruleAllObject(uni)
	if err != nil {
		return err
	}
	return nil
}

func ruleAllObject(uni *Union) error {
	for _, typ := range uni.typs {
		if _, ok := typ.(*Object); !ok {
			return fmt.Errorf("type %T not an Object", typ)
		}
	}
	return nil
}

func ruleAtLeastOneType(uni *Union) error {
	if len(uni.typs) <= 0 {
		return fmt.Errorf("union must define at least one type")
	}
	return nil
}

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

// ListType is a collection type wrapper.
type ListType struct {
	underline Type
}

func (lt *ListType) typ() {}

// NonNullType is a non-null type wrapper.
type NonNullType struct {
	underline Type
}

func (nt *NonNullType) typ() {}
