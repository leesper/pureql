package ast

import "fmt"

// Inspect traverses the AST in depth-first order; It starts by calling f(node);
// node must be non-nil. If f returns true, Inspect invokes f recursively for
// each of the non-nil children, followed by a call of f(nil)
func Inspect(node Node, f func(Node) bool) {
	Walk(inspector(f), node)
}

// Visitor interface for visitors traversing AST. The Visit method is called on
// every node encountered by Walk. if the result visitor is non-nil Walk will
// visits each of the children of node with the visitor v, followed by a call of
// v.Visit(nil).
type Visitor interface {
	Visit(node Node) (v Visitor)
}

type inspector func(Node) bool

func (f inspector) Visit(node Node) Visitor {
	if f(node) {
		return f
	}
	return nil
}

// Walk traverses an AST in depth-first order: It starts by calling v.Visit(node);
// node must be non-nil. If the visitor v returned by v.Visit(node) is not nil,
// Walk is invoked recursively with visitor v for each of the non-nil children
// of node, followed by a call of v.Visit(nil).
func Walk(v Visitor, node Node) {
	if v = v.Visit(node); v == nil {
		return
	}

	switch n := node.(type) {
	case *Document:
		for _, def := range n.Defs {
			Walk(v, def)
		}
	case *OperationDefinition:
		if n.VarDefns != nil {
			Walk(v, n.VarDefns)
		}
		if n.Directs != nil {
			Walk(v, n.Directs)
		}
		Walk(v, n.SelSet)
	case *SelectionSet:
		for _, s := range n.Sels {
			Walk(v, s)
		}
	case *Field:
		if n.Als != nil {
			Walk(v, n.Als)
		}
		if n.Args != nil {
			Walk(v, n.Args)
		}
		if n.Directs != nil {
			Walk(v, n.Directs)
		}
		if n.SelSet != nil {
			Walk(v, n.SelSet)
		}
	case *Alias:
		// do nothing
	case *Arguments:
		for _, a := range n.Args {
			Walk(v, a)
		}
	case *Argument:
		Walk(v, n.Val)
	case *FragmentSpread:
		if n.Directs != nil {
			Walk(v, n.Directs)
		}
	case *InlineFragment:
		if n.TypeCond != nil {
			Walk(v, n.TypeCond)
		}
		if n.Directs != nil {
			Walk(v, n.Directs)
		}
		Walk(v, n.SelSet)
	case *FragmentDefinition:
		Walk(v, n.TypeCond)
		if n.Directs != nil {
			Walk(v, n.Directs)
		}
		Walk(v, n.SelSet)
	case *TypeCondition:
		Walk(v, n.NamedTyp)
	case *Variable, *LiteralValue, *NameValue:
		// do nothing
	case *ListValue:
		for _, val := range n.Vals {
			Walk(v, val)
		}
	case *ObjectValue:
		for _, obj := range n.ObjFields {
			Walk(v, obj)
		}
	case *ObjectField:
		Walk(v, n.Val)
	case *VariableDefinitions:
		for _, vd := range n.VarDefns {
			Walk(v, vd)
		}
	case *VariableDefinition:
		Walk(v, n.Var)
		Walk(v, n.Typ)
		if n.DeflVal != nil {
			Walk(v, n.DeflVal)
		}
	case *DefaultValue:
		Walk(v, n.Val)
	case *NamedType:
		// do nothing
	case *ListType:
		Walk(v, n.Typ)
	case *Directives:
		for _, d := range n.Directs {
			Walk(v, d)
		}
	case *Directive:
		if n.Args != nil {
			Walk(v, n.Args)
		}
	case *Schema:
		for _, iface := range n.Interfaces {
			Walk(v, iface)
		}
		for _, scalar := range n.Scalars {
			Walk(v, scalar)
		}
		for _, input := range n.InputObjects {
			Walk(v, input)
		}
		for _, typ := range n.Types {
			Walk(v, typ)
		}
		for _, extend := range n.Extends {
			Walk(v, extend)
		}
		for _, direct := range n.Directives {
			Walk(v, direct)
		}
		for _, schema := range n.Schemas {
			Walk(v, schema)
		}
		for _, enum := range n.Enums {
			Walk(v, enum)
		}
		for _, union := range n.Unions {
			Walk(v, union)
		}
	case *InterfaceDefinition:
		if n.Directs != nil {
			Walk(v, n.Directs)
		}
		for _, fd := range n.FieldDefns {
			Walk(v, fd)
		}
	case *FieldDefinition:
		if n.ArgDefns != nil {
			Walk(v, n.ArgDefns)
		}
		Walk(v, n.Typ)
		if n.Directs != nil {
			Walk(v, n.Directs)
		}
	case *ArgumentsDefinition:
		for _, input := range n.InputValDefns {
			Walk(v, input)
		}
	case *InputValueDefinition:
		Walk(v, n.Typ)
		if n.DeflVal != nil {
			Walk(v, n.DeflVal)
		}
		if n.Directs != nil {
			Walk(v, n.Directs)
		}
	case *ScalarDefinition:
		if n.Directs != nil {
			Walk(v, n.Directs)
		}
	case *InputObjectDefinition:
		if n.Directs != nil {
			Walk(v, n.Directs)
		}
		for _, input := range n.InputValDefns {
			Walk(v, input)
		}
	case *TypeDefinition:
		if n.Implements != nil {
			Walk(v, n.Implements)
		}
		if n.Directs != nil {
			Walk(v, n.Directs)
		}
		for _, fd := range n.FieldDefns {
			Walk(v, fd)
		}
	case *ImplementsInterfaces:
		for _, namdTyp := range n.NamedTyps {
			Walk(v, namdTyp)
		}
	case *ExtendDefinition:
		if n.TypDefn != nil {
			Walk(v, n.TypDefn)
		}
	case *DirectiveDefinition:
		if n.Args != nil {
			Walk(v, n.Args)
		}
		if n.Locs != nil {
			Walk(v, n.Locs)
		}
	case *DirectiveLocations:
		for _, l := range n.Locs {
			Walk(v, l)
		}
	case *DirectiveLocation:
		// do nothing
	case *SchemaDefinition:
		if n.Directs != nil {
			Walk(v, n.Directs)
		}
		for _, o := range n.OperDefns {
			Walk(v, o)
		}
	case *OperationTypeDefinition:
		if n.NamedTyp != nil {
			Walk(v, n.NamedTyp)
		}
	case *EnumDefinition:
		if n.Directs != nil {
			Walk(v, n.Directs)
		}
		for _, e := range n.EnumVals {
			Walk(v, e)
		}
	case *EnumValue:
		if n.Directs != nil {
			Walk(v, n.Directs)
		}
	case *UnionDefinition:
		if n.Directs != nil {
			Walk(v, n.Directs)
		}
		if n.Members != nil {
			Walk(v, n.Members)
		}
	case *UnionMembers:
		if n.NamedTyp != nil {
			Walk(v, n.NamedTyp)
		}
		for _, m := range n.Members {
			Walk(v, m)
		}
	case *UnionMember:
		if n.NamedTyp != nil {
			Walk(v, n.NamedTyp)
		}
	default:
		panic(fmt.Sprintf("parser.Walk: unexpected node type %T", n))
	}

	v.Visit(nil)
}
