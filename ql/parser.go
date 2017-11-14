package ql

import (
	"errors"
	"fmt"
)

// ErrBadToken for invalid token.
type ErrBadToken struct {
	line int
	tok  Token
}

func (e ErrBadToken) Error() string {
	return fmt.Sprintf("line %d: bad token %s", e.line, e.tok)
}

// ErrBadParse for invalid parse.
type ErrBadParse struct {
	line   int
	expect string
	found  Token
}

func (e ErrBadParse) Error() string {
	return fmt.Sprintf("line %d: expecting %s, found %s", e.line, e.expect, e.found)
}

// ParseDocument returns ast.Document using DefaultParser.
func ParseDocument(document string) error {
	return errors.New("not implemented")
}

// ParseSchema returns ast.Schema using DefaultParser.
func ParseSchema(schema string) error {
	return errors.New("not implemented")
}

// ParseOperation returns ast.Operation using DefaultParser.
func ParseOperation(oper string) error {
	return errors.New("not implemented")
}

// Parser converts GraphQL source into AST.
type Parser struct {
	lexer     *Lexer
	lookAhead Token
}

// NewParser returns a new Parser.
func NewParser(l *Lexer) *Parser {
	return &Parser{lexer: l}
}

// ParseDocument returns ast.Document.
func (p *Parser) ParseDocument() error {
	if p.lexer == nil {
		return errors.New("lexer nil")
	}
	p.lookAhead = p.lexer.Read()

	err := p.definition()
	if err != nil {
		return err
	}

	for p.lookAhead != TokenEOF {
		err = p.definition()
		if err != nil {
			return err
		}
	}
	return nil
}

// ParseSchema returns ast.Schema.
func (p *Parser) ParseSchema(schema string) error {
	return errors.New("not implemented")
}

// ParseOperation returns ast.Operation.
func (p *Parser) ParseOperation(oper string) error {
	return errors.New("not implemented")
}

func (p *Parser) definition() error {
	switch p.lookAhead.Kind {
	case QUERY, MUTATION, LBRACE:
		return p.operationDefinition()
	case FRAGMENT:
		return p.fragmentDefinition()
	default:
		expect := fmt.Sprintf("%s/%s/%s", tokens[QUERY], tokens[MUTATION], tokens[LBRACE])
		return ErrBadParse{
			line:   p.lexer.Line(),
			expect: expect,
			found:  p.lookAhead,
		}
	}
}

func (p *Parser) operationDefinition() error {
	if p.lookAhead.Kind == LBRACE {
		return p.selectionSet()
	}

	if !p.match(QUERY) {
		if !p.match(MUTATION) {
			expect := fmt.Sprintf("%s/%s", tokens[QUERY], tokens[MUTATION])
			return ErrBadParse{
				line:   p.lexer.Line(),
				expect: expect,
				found:  p.lookAhead,
			}
		}
	}

	p.match(NAME)

	var err error
	if p.lookAhead.Kind == LPAREN {
		if err = p.variableDefinitions(); err != nil {
			return err
		}
	}

	for p.lookAhead.Kind == AT {
		if err = p.directives(); err != nil {
			return err
		}
	}

	return p.selectionSet()
}

func (p *Parser) variableDefinitions() error {
	return errors.New("not implemented")
}

func (p *Parser) directives() error {
	return errors.New("not implemented")
}

func (p *Parser) directive() error {
	return errors.New("not implemented")
}

func (p *Parser) match(k Kind) bool {
	if p.lookAhead.Kind == k {
		p.consume()
		return true
	}
	return false
}

func (p *Parser) consume() {
	p.lookAhead = p.lexer.Read()
}

func (p *Parser) selectionSet() error {
	return errors.New("not implemented")
}

func (p *Parser) fragmentDefinition() error {
	return errors.New("not implemented")
}

// Parse parses the tokens, it returns error if something goes wrong.
func Parse() error {
	// var err error
	// p.lookAhead = p.lexer.Read()
	// for p.lookAhead != TokenEOF {
	// 	switch p.lookAhead.Text {
	// 	case tokens[QUERY], tokens[MUTATION], tokens[LBRACE]:
	// 		err = p.document()
	// 		if err != nil {
	// 			return err
	// 		}
	// 	case tokens[INTERFACE]:
	// 		err = p.interfaceDefinition()
	// 		if err != nil {
	// 			return err
	// 		}
	// 	case tokens[SCALAR]:
	// 		err = p.scalarDefinition()
	// 		if err != nil {
	// 			return err
	// 		}
	// 	case tokens[INPUT]:
	// 		err = p.inputObjectDefinition()
	// 		if err != nil {
	// 			return err
	// 		}
	// 	case tokens[TYPE]:
	// 		err = p.typeDefinition()
	// 		if err != nil {
	// 			return err
	// 		}
	// 	case tokens[SCHEMA]:
	// 		err = p.schemaDefinition()
	// 		if err != nil {
	// 			return err
	// 		}
	// 	default:
	// 		expecting := fmt.Sprintf("%s/%s/%s/%s/%s/%s/%s/%s",
	// 			tokens[QUERY], tokens[MUTATION], tokens[LBRACE], tokens[INTERFACE],
	// 			tokens[SCALAR], tokens[INPUT], tokens[TYPE], tokens[SCHEMA])
	// 		return ErrBadParse{p.lexer.Line(), expecting, p.lookAhead}
	// 	}
	// 	p.lookAhead = p.lexer.Read()
	// }
	return nil
}

func interfaceDefinition() error {
	return errors.New("not implemented")
}

func scalarDefinition() error {
	return errors.New("not implemented")
}

func inputObjectDefinition() error {
	return errors.New("not implemented")
}

func typeDefinition() error {
	return errors.New("not implemented")
}

func schemaDefinition() error {
	return errors.New("not implemented")
}

func document() error {
	return errors.New("not implemented")
}

func selection() error {
	return errors.New("not implemented")
}

func field() error {
	return errors.New("not implemented")
}

func fragmentSpread() error {
	return errors.New("not implemented")
}

func inlineFragment() error {
	return errors.New("not implemented")
}

func alias() error {
	return errors.New("not implemented")
}

func arguments() error {
	return errors.New("not implemented")
}

func argument() error {
	return errors.New("not implemented")
}

func typeCondition() error {
	return errors.New("not implemented")
}

func value() error {
	return errors.New("not implemented")
}

func variable() error {
	return errors.New("not implemented")
}

func booleanValue() error {
	return errors.New("not implemented")
}

func nullValue() error {
	return errors.New("not implemented")
}

func enumValue() error {
	return errors.New("not implemented")
}

func listValue() error {
	return errors.New("not implemented")
}

func objectValue() error {
	return errors.New("not implemented")
}

func objectField() error {
	return errors.New("not implemented")
}

func defaultValue() error {
	return errors.New("not implemented")
}

func types() error {
	return errors.New("not implemented")
}

func namedType() error {
	return errors.New("not implemented")
}

func listType() error {
	return errors.New("not implemented")
}

func nonNullType() error {
	return errors.New("not implemented")
}
