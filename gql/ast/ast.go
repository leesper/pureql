package ast

import "go/token"

// Node is the interface for all AST node types.
type Node interface {
	Pos() token.Pos // position of first character belong to the node
	End() token.Pos // position of first character immediately after the node
}
