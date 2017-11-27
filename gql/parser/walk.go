package parser

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
			v.Visit(def)
		}
	case *OperationDefinition:
		if n.VarDefns != nil {
			v.Visit(n.VarDefns)
		}
		if n.Directs != nil {
			v.Visit(n.Directs)
		}
		v.Visit(n.SelSet)
	case *SelectionSet:
		for _, s := range n.Sels {
			v.Visit(s)
		}
	case *Field:
		if n.Als != nil {
			v.Visit(n.Als)
		}
		if n.Args != nil {
			v.Visit(n.Args)
		}
		if n.Directs != nil {
			v.Visit(n.Directs)
		}
		if n.SelSet != nil {
			v.Visit(n.SelSet)
		}
	case *Alias:
		// do nothing
	case *Arguments:
		for _, a := range n.Args {
			v.Visit(a)
		}
	case *Argument:
		v.Visit(n.Val)
	case *FragmentSpread:
		if n.Directs != nil {
			v.Visit(n.Directs)
		}
	case *InlineFragment:
		if n.TypeCond != nil {
			v.Visit(n.TypeCond)
		}
		if n.Directs != nil {
			v.Visit(n.Directs)
		}
		v.Visit(n.SelSet)
	case *FragmentDefinition:
		v.Visit(n.TypeCond)
		if n.Directs != nil {
			v.Visit(n.Directs)
		}
		v.Visit(n.SelSet)
	case *TypeCondition:
		v.Visit(n.NamedTyp)
	case *Variable, *LiteralValue, *NameValue:
		// do nothing
	case *ListValue:
		for _, val := range n.Vals {
			v.Visit(val)
		}
	case *ObjectValue:
		for _, obj := range n.ObjFields {
			v.Visit(obj)
		}
	case *ObjectField:
		v.Visit(n.Val)
	case *VariableDefinitions:
		for _, vd := range n.VarDefns {
			v.Visit(vd)
		}
	case *VariableDefinition:
		v.Visit(n.Var)
		v.Visit(n.Typ)
		if n.DeflVal != nil {
			v.Visit(n.DeflVal)
		}
	case *DefaultValue:
		v.Visit(n.Val)
	case *NamedType:
		// do nothing
	case *ListType:
		v.Visit(n.Typ)
	case *Directives:
		for _, d := range n.Directs {
			v.Visit(d)
		}
	case *Directive:
		if n.Args != nil {
			v.Visit(n.Args)
		}
	case *Schema:
		for _, iface := range n.Interfaces {
			v.Visit(iface)
		}
		for _, scalar := range n.Scalars {
			v.Visit(scalar)
		}
		for _, input := range n.InputObjects {
			v.Visit(input)
		}
		for _, typ := range n.Types {
			v.Visit(typ)
		}
		for _, extend := range n.Extends {
			v.Visit(extend)
		}
		for _, direct := range n.Directives {
			v.Visit(direct)
		}
		for _, schema := range n.Schemas {
			v.Visit(schema)
		}
		for _, enum := range n.Enums {
			v.Visit(enum)
		}
		for _, union := range n.Unions {
			v.Visit(union)
		}
	case *InterfaceDefinition:
		if n.Directs != nil {
			v.Visit(n.Directs)
		}
		for _, fd := range n.FieldDefns {
			v.Visit(fd)
		}
	case *FieldDefinition:
		if n.ArgDefns != nil {
			v.Visit(n.ArgDefns)
		}
		v.Visit(n.Typ)
		if n.Directs != nil {
			v.Visit(n.Directs)
		}
	case *ArgumentsDefinition:
		for _, input := range n.InputValDefns {
			v.Visit(input)
		}
	case *InputValueDefinition:
		v.Visit(n.Typ)
		if n.DeflVal != nil {
			v.Visit(n.DeflVal)
		}
		if n.Directs != nil {
			v.Visit(n.Directs)
		}
	case *ScalarDefinition:
		if n.Directs != nil {
			v.Visit(n.Directs)
		}
	case *InputObjectDefinition:
		if n.Directs != nil {
			v.Visit(n.Directs)
		}
		for _, input := range n.InputValDefns {
			v.Visit(input)
		}
	case *TypeDefinition:
		if n.Implements != nil {
			v.Visit(n.Implements)
		}
		if n.Directs != nil {
			v.Visit(n.Directs)
		}
		for _, fd := range n.FieldDefns {
			v.Visit(fd)
		}
	case *ImplementsInterfaces:
		for _, namdTyp := range n.NamedTyps {
			v.Visit(namdTyp)
		}
	case *ExtendDefinition:
		if n.TypDefn != nil {
			v.Visit(n.TypDefn)
		}
	case *DirectiveDefinition:
		if n.Args != nil {
			v.Visit(n.Args)
		}
		if n.Locs != nil {
			v.Visit(n.Locs)
		}
	case *DirectiveLocations:
		for _, l := range n.Locs {
			v.Visit(l)
		}
	case *DirectiveLocation:
		// do nothing
	case *SchemaDefinition:
		if n.Directs != nil {
			v.Visit(n.Directs)
		}
		for _, o := range n.OperDefns {
			v.Visit(o)
		}
	case *OperationTypeDefinition:
		if n.NamedTyp != nil {
			v.Visit(n.NamedTyp)
		}
	case *EnumDefinition:
		if n.Directs != nil {
			v.Visit(n.Directs)
		}
		for _, e := range n.EnumVals {
			v.Visit(e)
		}
	case *EnumValue:
		if n.Directs != nil {
			v.Visit(n.Directs)
		}
	case *UnionDefinition:
		if n.Directs != nil {
			v.Visit(n.Directs)
		}
		if n.Members != nil {
			v.Visit(n.Members)
		}
	case *UnionMembers:
		if n.NamedTyp != nil {
			v.Visit(n.NamedTyp)
		}
		for _, m := range n.Members {
			v.Visit(m)
		}
	case *UnionMember:
		if n.NamedTyp != nil {
			v.Visit(n.NamedTyp)
		}
	}

	v.Visit(nil)
}
