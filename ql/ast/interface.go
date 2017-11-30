package ast

import (
	"errors"
	"fmt"
	"go/token"
)

// ErrBadParse for invalid parse.
type ErrBadParse struct {
	pos    token.Position
	expect string
	found  string
}

func (e ErrBadParse) Error() string {
	return fmt.Sprintf("%s: expecting %s, found '%s'", e.pos, e.expect, e.found)
}

// ParseDocument returns ast.Document.
func ParseDocument(document []byte, filename string, fset *token.FileSet) (*Document, error) {
	if fset == nil {
		return nil, errors.New("no token.FileSet provided (fset == nil)")
	}
	return newParser(document, filename, fset).parseDocument()
}

// ParseSchema returns ast.Schema.
func ParseSchema(schema []byte, filename string, fset *token.FileSet) (*Schema, error) {
	if fset == nil {
		return nil, errors.New("no token.FileSet provided (fset == nil)")
	}
	return newParser(schema, filename, fset).parseSchema()
}
