package ql

import (
	"fmt"
	"strconv"

	"github.com/leesper/pureql/ql/ast"
)

const (
	query    = "query"
	mutation = "mutation"
)

// Type interface for NamedType and ListType.
type Type interface {
	typeType()
}

// Value interface for NamedValue and ListValue.
type Value interface {
	valueType()
}

// Document is a query document type.
type Document struct {
	opers []*Operation
	frags []*Fragment
}

func newDocument(node *ast.Document) *Document {
	doc := &Document{}
	for _, def := range node.Defs {
		switch def := def.(type) {
		case *ast.OperationDefinition:
			doc.opers = append(doc.opers, newOperation(def))
		case *ast.FragmentDefinition:
			doc.frags = append(doc.frags, newFragment(def))
		}
	}
	return doc
}

// Operation .
type Operation struct {
	operType  string
	name      string
	varDefns  []*VariableDefinition
	directive []*Directive
	selSet    *SelectionSet
}

func newOperation(node *ast.OperationDefinition) *Operation {
	oper := &Operation{
		operType: node.OperType.Text,
		name:     node.Name.Text,
	}
	if node.VarDefns != nil {
		oper.varDefns = newVariableDefinitions(node.VarDefns)
	}
	if node.Directs != nil {
		oper.directive = newDirectives(node.Directs)
	}
	oper.selSet = newSelectionSet(node.SelSet)

	return oper
}

// Fragment is a fragment definition type.
type Fragment struct{}

func newFragment(node *ast.FragmentDefinition) *Fragment {
	return &Fragment{}
}

// Directive is a defined directive.
type Directive struct {
	name string
	args []*Argument
}

func newDirective(node *ast.Directive) *Directive {
	return &Directive{
		name: node.Name.Text,
		args: newArguments(node.Args),
	}
}

func newDirectives(node *ast.Directives) []*Directive {
	var directives []*Directive
	for _, d := range node.Directs {
		directives = append(directives, newDirective(d))
	}
	return directives
}

// Argument represents argument in GraphQL.
type Argument struct {
	name string
	val  Value
}

func newArgument(node *ast.Argument) *Argument {
	return &Argument{
		name: node.Name.Text,
		val:  newValue(node.Val),
	}
}

func newArguments(node *ast.Arguments) []*Argument {
	var args []*Argument
	for _, a := range node.Args {
		args = append(args, newArgument(a))
	}
	return args
}

// VariableDefinition represents variable definitions.
type VariableDefinition struct {
	va   *Variable
	typ  Type
	defl *DefaultValue
}

func newVariableDefinition(node *ast.VariableDefinition) *VariableDefinition {
	return &VariableDefinition{
		va:   newVariable(node.Var),
		typ:  newType(node.Typ),
		defl: newDefaultValue(node.DeflVal),
	}
}

func newVariableDefinitions(node *ast.VariableDefinitions) []*VariableDefinition {
	var defs []*VariableDefinition
	for _, def := range node.VarDefns {
		defs = append(defs, newVariableDefinition(def))
	}
	return defs
}

// Variable represents varibles in variable definition.
type Variable struct {
	name string
}

func (v *Variable) valueType() {}

func newVariable(node *ast.Variable) *Variable {
	return &Variable{
		name: node.Name.Text,
	}
}

// IntValue represents integer values in GraphQL.
type IntValue int

func (i IntValue) valueType() {}

// FloatValue represents floating values in GraphQL.
type FloatValue float64

func (f FloatValue) valueType() {}

// StringValue represents string values in GraphQL.
type StringValue string

func (s StringValue) valueType() {}

// NamedValue represents name value in GraphQL.
type NamedValue string

func (n NamedValue) valueType() {}

// ListValue represents list of values in GraphQL.
type ListValue []Value

func (l ListValue) valueType() {}

func newListValue(node *ast.ListValue) ListValue {
	var vals []Value
	for _, v := range node.Vals {
		vals = append(vals, newValue(v))
	}
	return ListValue(vals)
}

// ObjectValue represents object value in GraphQL.
type ObjectValue map[string]Value

func (o ObjectValue) valueType() {}

func newObjectValue(node *ast.ObjectValue) ObjectValue {
	ov := map[string]Value{}
	for _, f := range node.ObjFields {
		ov[f.Name.Text] = newValue(f.Val)
	}
	return ObjectValue(ov)
}

// DefaultValue represents default value in variable definition.
type DefaultValue struct {
	val Value
}

func newDefaultValue(node *ast.DefaultValue) *DefaultValue {
	return &DefaultValue{
		val: newValue(node.Val),
	}
}

func newType(node ast.Type) Type {
	switch typ := node.(type) {
	case *ast.NamedType:
		return newNamedType(typ)
	case *ast.ListType:
		return newListType(typ)
	default:
		panic(fmt.Sprintf("unexpected type %T", typ))
	}
}

// NamedType represents a named type defined in schema, nonNull indicates whether
// it is nullable or not.
type NamedType struct {
	name    string
	nonNull bool
}

func newNamedType(node *ast.NamedType) *NamedType {
	return &NamedType{
		name:    node.Name.Text,
		nonNull: node.NonNull,
	}
}

func (nt *NamedType) typeType() {}

// ListType represents a list of type defined in schema., nonNull indicates whether
// it is nullable or not.
type ListType struct {
	typ     Type
	nonNull bool
}

func newListType(node *ast.ListType) *ListType {
	return &ListType{
		typ:     newType(node.Typ),
		nonNull: node.NonNull,
	}
}

func (lt *ListType) typeType() {}

func newValue(node ast.Value) Value {
	switch val := node.(type) {
	case *ast.Variable:
		return newVariable(val)
	case *ast.LiteralValue:
		if val.Val.Kind == ast.INT {
			i, _ := strconv.Atoi(val.Val.Text)
			return IntValue(i)
		} else if val.Val.Kind == ast.FLOAT {
			f, _ := strconv.ParseFloat(val.Val.Text, 0)
			return FloatValue(f)
		}
		return StringValue(val.Val.Text)
	case *ast.NameValue:
		return NamedValue(val.Val.Text)
	case *ast.ListValue:
		return newListValue(val)
	case *ast.ObjectValue:
		return newObjectValue(val)
	default:
		panic(fmt.Sprintf("unexpected value type %T", val))
	}
}

// SelectionSet represents selection set defined in GraphQL.
type SelectionSet struct{}

func newSelectionSet(node *ast.SelectionSet) *SelectionSet {
	return &SelectionSet{}
}

// Selection is an interface for Field, FragmentSpread, InlineFragment
type Selection interface{}
