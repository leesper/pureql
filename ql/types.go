package ql

import (
	"fmt"
	"strings"
)

// Type interface for all types.
type Type interface {
	Type() string
}

// Schema is the entry point of GraphQL service.
type Schema struct {
	Qry *Object
	Mut *Object
}

// Scalar represents primitive value.
type Scalar struct {
	Name string
}

// Type returns basic type info.
func (scalar *Scalar) Type() string {
	return fmt.Sprintf("scalar %s", scalar.Name)
}

// Enum represents limited enumerable values.
type Enum struct {
	Name string
}

// Type returns basic type info.
func (enum *Enum) Type() string {
	return fmt.Sprintf("enum %s", enum.Name)
}

// Object defines a set of fields of another type in the type system.
type Object struct {
	Name   string
	Ifaces []*Interface
	Fields []*Field
}

// Type returns basic type info.
func (obj *Object) Type() string {
	var fieldInfos []string
	for _, f := range obj.Fields {
		fieldInfos = append(fieldInfos, f.Typ.Type())
	}
	return fmt.Sprintf("object %s { %s }", obj.Name, strings.Join(fieldInfos, " "))
}

// Interface defines an abstract type for Object to implement.
type Interface struct {
	Name   string
	Fields []*Field
}

// Type returns basic type info.
func (iface *Interface) Type() string {
	var fieldInfos []string
	for _, f := range iface.Fields {
		fieldInfos = append(fieldInfos, f.Typ.Type())
	}
	return fmt.Sprintf("interface %s { %s }", iface.Name, strings.Join(fieldInfos, " "))
}

// Union defines a list of possible Object types.
type Union struct {
	Name string
	Typs []Type
}

// Type returns basic type info.
func (union *Union) Type() string {
	var typeInfos []string
	for _, t := range union.Typs {
		typeInfos = append(typeInfos, t.Type())
	}
	return fmt.Sprintf("union %s %s", union.Name, strings.Join(typeInfos, "|"))
}

// List is a list of other types.
type List struct {
	OfType Type
}

// Type returns basic type info.
func (lt *List) Type() string {
	return fmt.Sprintf("[ %s ]", lt.OfType.Type())
}

// NonNull is a non-null wrapper of other types.
type NonNull struct {
	OfType Type
}

// Type returns basic type info.
func (nn *NonNull) Type() string {
	return fmt.Sprintf("%s!", nn.OfType.Type())
}

// InputObject is a struct for complex input.
type InputObject struct {
	Name   string
	Fields []*Field
}

// Type returns basic type info.
func (io *InputObject) Type() string {
	var fieldInfos []string
	for _, f := range io.Fields {
		fieldInfos = append(fieldInfos, f.Typ.Type())
	}
	return fmt.Sprintf("input %s { %s }", io.Name, strings.Join(fieldInfos, " "))
}

// Directive represents directives the execution engine supports.
type Directive struct{}
