package parser

import (
	"go/token"
)

// Node is the interface for all AST node types.
type Node interface {
	Pos() token.Pos // position of first character belong to the node
	End() token.Pos // position of first character immediately after the node
}

// Query Document related-------------------------------------------------------

// Definition is the interface for all definition node types:
// OperationDefinition, FragmentDefinition
type Definition interface {
	Node
	defnNode()
}

// Selection is the interface for all selection node types:
// Field, FragmentSpread, InlineFragment
type Selection interface {
	Node
	selectNode()
}

// Value is the interface for all value node types:
// Variable, LiteralValue, NameValue, ListValue, ObjectValue
type Value interface {
	Node
	valueNode()
}

// Type is the interface for all type node types:
// NamedType, ListType, NonNullType
type Type interface {
	Node
	typeNode()
}

// Document node
type Document struct {
	Defs []Definition
}

// Pos returns position of first character belong to the node
func (d *Document) Pos() token.Pos {
	return d.Defs[0].Pos()
}

// End returns position of first character immediately after the node
func (d *Document) End() token.Pos {
	return d.Defs[len(d.Defs)-1].End()
}

// Definition related-----------------------------------------------------------

// OperationDefinition node
type OperationDefinition struct {
	OperType Token
	OperPos  token.Pos
	Name     Token
	NamePos  token.Pos
	VarDefns *VariableDefinitions
	Directs  *Directives
	SelSet   *SelectionSet
}

// Pos returns position of first character belong to the node
func (o *OperationDefinition) Pos() token.Pos {
	return o.OperPos
}

// End returns position of first character immediately after the node
func (o *OperationDefinition) End() token.Pos {
	return o.SelSet.End()
}

func (o *OperationDefinition) defnNode() {}

// FragmentDefinition node
type FragmentDefinition struct {
	Fragment token.Pos
	Name     Token
	NamePos  token.Pos
	TypeCond *TypeCondition
	Directs  *Directives
	SelSet   *SelectionSet
}

// Pos returns position of first character belong to the node
func (f *FragmentDefinition) Pos() token.Pos {
	return f.Fragment
}

// End returns position of first character immediately after the node
func (f *FragmentDefinition) End() token.Pos {
	return f.SelSet.End()
}

func (f *FragmentDefinition) defnNode() {}

// Selection related------------------------------------------------------------

// Field node
type Field struct {
	Als     *Alias
	Name    Token
	NamePos token.Pos
	Args    *Arguments
	Directs *Directives
	SelSet  *SelectionSet
}

// Pos returns position of first character belong to the node
func (f *Field) Pos() token.Pos {
	if f.Als != nil {
		return f.Als.NamePos
	}
	return f.NamePos
}

// End returns position of first character immediately after the node
func (f *Field) End() token.Pos {
	switch {
	case f.SelSet != nil:
		return f.SelSet.End()
	case f.Directs != nil:
		return f.Directs.End()
	case f.Args != nil:
		return f.Args.End()
	default:
		return token.Pos(int(f.NamePos) + 1)
	}
}

func (f *Field) selectNode() {}

// FragmentSpread node
type FragmentSpread struct {
	Spread  token.Pos
	Name    Token
	NamePos token.Pos
	Directs *Directives
}

// Pos returns position of first character belong to the node
func (f *FragmentSpread) Pos() token.Pos {
	return f.Spread
}

// End returns position of first character immediately after the node
func (f *FragmentSpread) End() token.Pos {
	if f.Directs != nil {
		return f.Directs.End()
	}
	return token.Pos(int(f.NamePos) + len(f.Name.Text))
}

func (f *FragmentSpread) selectNode() {}

// InlineFragment node
type InlineFragment struct {
	Spread   token.Pos
	TypeCond *TypeCondition
	Directs  *Directives
	SelSet   *SelectionSet
}

// Pos returns position of first character belong to the node
func (i *InlineFragment) Pos() token.Pos {
	return i.Spread
}

// End returns position of first character immediately after the node
func (i *InlineFragment) End() token.Pos {
	return i.SelSet.End()
}

func (i *InlineFragment) selectNode() {}

// Value related----------------------------------------------------------------

// Variable node
type Variable struct {
	Dollar  token.Pos
	Name    Token
	NamePos token.Pos
}

// Pos returns position of first character belong to the node
func (v *Variable) Pos() token.Pos {
	return v.Dollar
}

// End returns position of first character immediately after the node
func (v *Variable) End() token.Pos {
	return token.Pos(int(v.NamePos) + len(v.Name.Text))
}

func (v *Variable) valueNode() {}

// LiteralValue node for IntValue, FloatValue and StringValue
type LiteralValue struct {
	Val    Token
	ValPos token.Pos
}

// Pos returns position of first character belong to the node
func (l *LiteralValue) Pos() token.Pos {
	return l.ValPos
}

// End returns position of first character immediately after the node
func (l *LiteralValue) End() token.Pos {
	extra := 0 // if STRING add offset for two duoble-quotes
	if l.Val.Kind == STRING {
		extra = 2
	}
	return token.Pos(int(l.ValPos) + len(l.Val.Text) + extra)
}

func (l *LiteralValue) valueNode() {}

// NameValue node for BooleanValue and NullValue
type NameValue struct {
	Val    Token
	ValPos token.Pos
}

// Pos returns position of first character belong to the node
func (n *NameValue) Pos() token.Pos {
	return n.ValPos
}

// End returns position of first character immediately after the node
func (n *NameValue) End() token.Pos {
	return token.Pos(int(n.ValPos) + len(n.Val.Text))
}

func (n *NameValue) valueNode() {}

// ListValue node
type ListValue struct {
	Lbrack token.Pos
	Vals   []Value
	Rbrack token.Pos
}

// Pos returns position of first character belong to the node
func (l *ListValue) Pos() token.Pos {
	return l.Lbrack
}

// End returns position of first character immediately after the node
func (l *ListValue) End() token.Pos {
	return token.Pos(int(l.Rbrack) + 1)
}

func (l *ListValue) valueNode() {}

// ObjectValue node
type ObjectValue struct {
	Lbrace    token.Pos
	ObjFields []*ObjectField
	Rbrace    token.Pos
}

// Pos returns position of first character belong to the node
func (o *ObjectValue) Pos() token.Pos {
	return o.Lbrace
}

// End returns position of first character immediately after the node
func (o *ObjectValue) End() token.Pos {
	return token.Pos(int(o.Rbrace) + 1)
}

func (o *ObjectValue) valueNode() {}

// Embedded node types----------------------------------------------------------

// SelectionSet node
type SelectionSet struct {
	Lbrace token.Pos
	Sels   []Selection
	Rbrace token.Pos
}

// Pos returns position of first character belong to the node
func (s *SelectionSet) Pos() token.Pos {
	return s.Lbrace
}

// End returns position of first character immediately after the node
func (s *SelectionSet) End() token.Pos {
	return token.Pos(int(s.Rbrace) + 1)
}

// Alias node
type Alias struct {
	Name    Token
	NamePos token.Pos
	Colon   token.Pos
}

// Pos returns position of first character belong to the node
func (a *Alias) Pos() token.Pos {
	return a.NamePos
}

// End returns position of first character immediately after the node
func (a *Alias) End() token.Pos {
	return token.Pos(int(a.Colon) + 1)
}

// Arguments node
type Arguments struct {
	Lparen token.Pos
	Args   []*Argument
	Rparen token.Pos
}

// Pos returns position of first character belong to the node
func (a *Arguments) Pos() token.Pos {
	return a.Lparen
}

// End returns position of first character immediately after the node
func (a *Arguments) End() token.Pos {
	return token.Pos(int(a.Rparen) + 1)
}

// Argument node
type Argument struct {
	Name    Token
	NamePos token.Pos
	Colon   token.Pos
	Val     Value
}

// Pos returns position of first character belong to the node
func (a *Argument) Pos() token.Pos {
	return a.NamePos
}

// End returns position of first character immediately after the node
func (a *Argument) End() token.Pos {
	return a.Val.End()
}

// TypeCondition node
type TypeCondition struct {
	On       token.Pos
	NamedTyp *NamedType
}

// Pos returns position of first character belong to the node
func (t *TypeCondition) Pos() token.Pos {
	return t.On
}

// End returns position of first character immediately after the node
func (t *TypeCondition) End() token.Pos {
	return t.NamedTyp.End()
}

// ObjectField node
type ObjectField struct {
	Name    Token
	NamePos token.Pos
	Colon   token.Pos
	Val     Value
}

// Pos returns position of first character belong to the node
func (o *ObjectField) Pos() token.Pos {
	return o.NamePos
}

// End returns position of first character immediately after the node
func (o *ObjectField) End() token.Pos {
	return o.Val.End()
}

// VariableDefinitions node
type VariableDefinitions struct {
	Lparen   token.Pos
	VarDefns []*VariableDefinition
	Rparen   token.Pos
}

// Pos returns position of first character belong to the node
func (v *VariableDefinitions) Pos() token.Pos {
	return v.Lparen
}

// End returns position of first character immediately after the node
func (v *VariableDefinitions) End() token.Pos {
	return token.Pos(int(v.Rparen) + 1)
}

// VariableDefinition node
type VariableDefinition struct {
	Var     *Variable
	Colon   token.Pos
	Typ     Type
	DeflVal *DefaultValue
}

// Pos returns position of first character belong to the node
func (v *VariableDefinition) Pos() token.Pos {
	return v.Var.Pos()
}

// End returns position of first character immediately after the node
func (v *VariableDefinition) End() token.Pos {
	if v.DeflVal != nil {
		return v.DeflVal.End()
	}
	return v.Typ.End()
}

// DefaultValue node
type DefaultValue struct {
	Eq  token.Pos
	Val Value
}

// Pos returns position of first character belong to the node
func (d *DefaultValue) Pos() token.Pos {
	return d.Eq
}

// End returns position of first character immediately after the node
func (d *DefaultValue) End() token.Pos {
	return d.Val.End()
}

// NamedType node
type NamedType struct {
	Name    Token
	NamePos token.Pos
	NonNull bool
	BangPos token.Pos
}

// Pos returns position of first character belong to the node
func (n *NamedType) Pos() token.Pos {
	return n.NamePos
}

// End returns position of first character immediately after the node
func (n *NamedType) End() token.Pos {
	if n.NonNull {
		return token.Pos(int(n.BangPos) + 1)
	}
	return token.Pos(int(n.NamePos) + len(n.Name.Text))
}

func (n *NamedType) typeNode() {}

// ListType node
type ListType struct {
	Lbrack  token.Pos
	Typ     Type
	Rbrack  token.Pos
	NonNull bool
	BangPos token.Pos
}

// Pos returns position of first character belong to the node
func (l *ListType) Pos() token.Pos {
	return l.Lbrack
}

// End returns position of first character immediately after the node
func (l *ListType) End() token.Pos {
	if l.NonNull {
		return token.Pos(int(l.BangPos) + 1)
	}
	return token.Pos(int(l.Rbrack) + 1)
}

func (l *ListType) typeNode() {}

// Directives node
type Directives struct {
	Directs []*Directive
}

// Pos returns position of first character belong to the node
func (d *Directives) Pos() token.Pos {
	return d.Directs[0].Pos()
}

// End returns position of first character immediately after the node
func (d *Directives) End() token.Pos {
	return d.Directs[len(d.Directs)-1].End()
}

// Directive node
type Directive struct {
	At      token.Pos
	Name    Token
	NamePos token.Pos
	Args    *Arguments
}

// Pos returns position of first character belong to the node
func (d *Directive) Pos() token.Pos {
	return d.At
}

// End returns position of first character immediately after the node
func (d *Directive) End() token.Pos {
	if d.Args != nil {
		return d.Args.End()
	}
	return token.Pos(int(d.NamePos) + len(d.Name.Text))
}

// Schema related---------------------------------------------------------------

// Schema node contains all definition nodes
type Schema struct {
	Interfaces   []*InterfaceDefinition
	Scalars      []*ScalarDefinition
	InputObjects []*InputObjectDefinition
	Types        []*TypeDefinition
	Extends      []*ExtendDefinition
	Directives   []*DirectiveDefinition
	Schemas      []*SchemaDefinition
	Enums        []*EnumDefinition
	Unions       []*UnionDefinition
	pos, end     token.Pos
}

// Pos returns position of first character belong to the node
func (s *Schema) Pos() token.Pos {
	return s.pos
}

// End returns position of first character immediately after the node
func (s *Schema) End() token.Pos {
	return token.Pos(int(s.end) + 1)
}

// InterfaceDefinition node
type InterfaceDefinition struct {
	Interface  token.Pos
	Name       Token
	NamePos    token.Pos
	Directs    *Directives
	Lbrace     token.Pos
	FieldDefns []*FieldDefinition
	Rbrace     token.Pos
}

// Pos returns position of first character belong to the node
func (i *InterfaceDefinition) Pos() token.Pos {
	return i.Interface
}

// End returns position of first character immediately after the node
func (i *InterfaceDefinition) End() token.Pos {
	return token.Pos(int(i.Rbrace) + 1)
}

// FieldDefinition node
type FieldDefinition struct {
	Name     Token
	NamePos  token.Pos
	ArgDefns *ArgumentsDefinition
	Colon    token.Pos
	Typ      Type
	Directs  *Directives
}

// Pos returns position of first character belong to the node
func (f *FieldDefinition) Pos() token.Pos {
	return f.NamePos
}

// End returns position of first character immediately after the node
func (f *FieldDefinition) End() token.Pos {
	if f.Directs != nil {
		return f.Directs.End()
	}
	return f.Typ.End()
}

// ArgumentsDefinition node
type ArgumentsDefinition struct {
	Lparen        token.Pos
	InputValDefns []*InputValueDefinition
	Rparen        token.Pos
}

// Pos returns position of first character belong to the node
func (a *ArgumentsDefinition) Pos() token.Pos {
	return a.Lparen
}

// End returns position of first character immediately after the node
func (a *ArgumentsDefinition) End() token.Pos {
	return token.Pos(int(a.Rparen) + 1)
}

// InputValueDefinition node
type InputValueDefinition struct {
	Name    Token
	NamePos token.Pos
	Colon   token.Pos
	Typ     Type
	DeflVal *DefaultValue
	Directs *Directives
}

// Pos returns position of first character belong to the node
func (i *InputValueDefinition) Pos() token.Pos {
	return i.NamePos
}

// End returns position of first character immediately after the node
func (i *InputValueDefinition) End() token.Pos {
	if i.Directs != nil {
		return i.Directs.End()
	}
	if i.DeflVal != nil {
		return i.DeflVal.End()
	}
	return i.Typ.End()
}

// ScalarDefinition node
type ScalarDefinition struct {
	Scalar  token.Pos
	Name    Token
	NamePos token.Pos
	Directs *Directives
}

// Pos returns position of first character belong to the node
func (s *ScalarDefinition) Pos() token.Pos {
	return s.Scalar
}

// End returns position of first character immediately after the node
func (s *ScalarDefinition) End() token.Pos {
	if s.Directs != nil {
		return s.Directs.End()
	}
	return token.Pos(int(s.NamePos) + 1)
}

// InputObjectDefinition node
type InputObjectDefinition struct {
	Input         token.Pos
	Name          Token
	NamePos       token.Pos
	Directs       *Directives
	Lbrace        token.Pos
	InputValDefns []*InputValueDefinition
	Rbrace        token.Pos
}

// Pos returns position of first character belong to the node
func (i *InputObjectDefinition) Pos() token.Pos {
	return i.Input
}

// End returns position of first character immediately after the node
func (i *InputObjectDefinition) End() token.Pos {
	return token.Pos(int(i.Rbrace) + 1)
}

// TypeDefinition node
type TypeDefinition struct {
	Typ        token.Pos
	Name       Token
	NamePos    token.Pos
	Implements *ImplementsInterfaces
	Directs    *Directives
	Lbrace     token.Pos
	FieldDefns []*FieldDefinition
	Rbrace     token.Pos
}

// Pos returns position of first character belong to the node
func (t *TypeDefinition) Pos() token.Pos {
	return t.Typ
}

// End returns position of first character immediately after the node
func (t *TypeDefinition) End() token.Pos {
	return token.Pos(int(t.Rbrace) + 1)
}

// ImplementsInterfaces node
type ImplementsInterfaces struct {
	Implements token.Pos
	NamedTyps  []*NamedType
}

// Pos returns position of first character belong to the node
func (i *ImplementsInterfaces) Pos() token.Pos {
	return i.Implements
}

// End returns position of first character immediately after the node
func (i *ImplementsInterfaces) End() token.Pos {
	return i.NamedTyps[len(i.NamedTyps)-1].End()
}

// ExtendDefinition node
type ExtendDefinition struct {
	Extend  token.Pos
	TypDefn *TypeDefinition
}

// Pos returns position of first character belong to the node
func (e *ExtendDefinition) Pos() token.Pos {
	return e.Extend
}

// End returns position of first character immediately after the node
func (e *ExtendDefinition) End() token.Pos {
	return e.TypDefn.End()
}

//DirectiveDefinition node
type DirectiveDefinition struct {
	Direct  token.Pos
	At      token.Pos
	Name    Token
	NamePos token.Pos
	Args    *ArgumentsDefinition
	On      token.Pos
	Locs    *DirectiveLocations
}

// Pos returns position of first character belong to the node
func (d *DirectiveDefinition) Pos() token.Pos {
	return d.Direct
}

// End returns position of first character immediately after the node
func (d *DirectiveDefinition) End() token.Pos {
	return d.Locs.End()
}

// DirectiveLocations node
type DirectiveLocations struct {
	Name    Token
	NamePos token.Pos
	Locs    []*DirectiveLocation
}

// Pos returns position of first character belong to the node
func (d *DirectiveLocations) Pos() token.Pos {
	return d.NamePos
}

// End returns position of first character immediately after the node
func (d *DirectiveLocations) End() token.Pos {
	if d.Locs != nil {
		return d.Locs[len(d.Locs)-1].End()
	}
	return token.Pos(int(d.NamePos) + len(d.Name.Text))
}

// DirectiveLocation node
type DirectiveLocation struct {
	Pipe    token.Pos
	Name    Token
	NamePos token.Pos
}

// Pos returns position of first character belong to the node
func (d *DirectiveLocation) Pos() token.Pos {
	return d.Pipe
}

// End returns position of first character immediately after the node
func (d *DirectiveLocation) End() token.Pos {
	return token.Pos(int(d.NamePos) + len(d.Name.Text))
}

// SchemaDefinition node
type SchemaDefinition struct {
	Schema    token.Pos
	Directs   *Directives
	Lbrace    token.Pos
	OperDefns []*OperationTypeDefinition
	Rbrace    token.Pos
}

// Pos returns position of first character belong to the node
func (s *SchemaDefinition) Pos() token.Pos {
	return s.Schema
}

// End returns position of first character immediately after the node
func (s *SchemaDefinition) End() token.Pos {
	return token.Pos(int(s.Rbrace) + 1)
}

// OperationTypeDefinition node
type OperationTypeDefinition struct {
	OperType Token
	OperPos  token.Pos
	Colon    token.Pos
	NamedTyp *NamedType
}

// Pos returns position of first character belong to the node
func (o *OperationTypeDefinition) Pos() token.Pos {
	return o.OperPos
}

// End returns position of first character immediately after the node
func (o *OperationTypeDefinition) End() token.Pos {
	return o.NamedTyp.End()
}

// EnumDefinition node
type EnumDefinition struct {
	Enum     token.Pos
	Name     Token
	NamePos  token.Pos
	Directs  *Directives
	Lbrace   token.Pos
	EnumVals []*EnumValue
	Rbrace   token.Pos
}

// Pos returns position of first character belong to the node
func (e *EnumDefinition) Pos() token.Pos {
	return e.Enum
}

// End returns position of first character immediately after the node
func (e *EnumDefinition) End() token.Pos {
	return token.Pos(int(e.Rbrace) + 1)
}

// EnumValue node
type EnumValue struct {
	Name    Token
	NamePos token.Pos
	Directs *Directives
}

// Pos returns position of first character belong to the node
func (e *EnumValue) Pos() token.Pos {
	return e.NamePos
}

// End returns position of first character immediately after the node
func (e *EnumValue) End() token.Pos {
	if e.Directs != nil {
		return e.Directs.End()
	}
	return token.Pos(int(e.NamePos) + len(e.Name.Text))
}

// UnionDefinition node
type UnionDefinition struct {
	Union   token.Pos
	Name    Token
	NamePos token.Pos
	Directs *Directives
	Eq      token.Pos
	Members *UnionMembers
}

// Pos returns position of first character belong to the node
func (u *UnionDefinition) Pos() token.Pos {
	return u.Union
}

// End returns position of first character immediately after the node
func (u *UnionDefinition) End() token.Pos {
	return u.Members.End()
}

// UnionMembers node
type UnionMembers struct {
	NamedTyp *NamedType
	Members  []*UnionMember
}

// Pos returns position of first character belong to the node
func (u *UnionMembers) Pos() token.Pos {
	return u.NamedTyp.Pos()
}

// End returns position of first character immediately after the node
func (u *UnionMembers) End() token.Pos {
	if u.Members != nil {
		return u.Members[len(u.Members)-1].End()
	}
	return u.NamedTyp.End()
}

// UnionMember node
type UnionMember struct {
	Pipe     token.Pos
	NamedTyp *NamedType
}

// Pos returns position of first character belong to the node
func (u *UnionMember) Pos() token.Pos {
	return u.Pipe
}

// End returns position of first character immediately after the node
func (u *UnionMember) End() token.Pos {
	return u.NamedTyp.End()
}
