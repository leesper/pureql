package ast

import (
	"go/token"

	"github.com/leesper/pureql/gql/parser"
)

// Node is the interface for all AST node types.
type Node interface {
	Pos() token.Pos // position of first character belong to the node
	End() token.Pos // position of first character immediately after the node
}

// Value related----------------------------------------------------------------
// Value is the interface for all value nodes.
type Value interface {
	Node
	Value()
}

// Variable -> $ Name
type Variable struct {
	Dollar token.Pos
	N      parser.Token
	NPos   token.Pos
}

// LiteralValue -> IntValue | FloatValue | StringValue
type LiteralValue struct {
	Tok    parser.Token
	TokPos token.Pos
}

// NameValue -> BooleanValue | NullValue | EnumValue
type NameValue struct {
	Tok    parser.Token
	TokPos token.Pos
}

// ListValue -> [ Value* ]
type ListValue struct {
	Lbrack token.Pos
	Vals   []Value
	Rbrack token.Pos
}

// ObjectValue -> { ObjectField* }
type ObjectValue struct {
	Lbrace    token.Pos
	objFields []ObjectField
	Rbrace    token.Pos
}

// ObjectField -> Name : Value
type ObjectField struct {
	N     parser.Token
	NPos  token.Pos
	Colon token.Pos
	Val   Value
}
