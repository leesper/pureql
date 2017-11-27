package parser

import (
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
func ParseDocument(document []byte) (*Document, error) {
	// document = []byte(strings.TrimRight(string(document), "\n\t\r "))
	return newParser(document, "").parseDocument()
}

// ParseSchema returns ast.Schema.
func ParseSchema(schema []byte) (*Schema, error) {
	// schema = []byte(strings.TrimRight(string(schema), "\n\t\r "))
	return newParser(schema, "").parseSchema()
}
