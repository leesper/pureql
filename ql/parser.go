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

// Parser converts tokens into AST. This is an LL(1) grammar parser based on
// Dr. Terence Parr's "Language Implementation Patterns".
type Parser struct {
	lexer     *Lexer
	lookAhead Token
}

// NewParser returns a new parser equipped with lexer.
func NewParser(l *Lexer) *Parser {
	return &Parser{
		lexer:     l,
		lookAhead: TokenEOF,
	}
}

// Parse parses the tokens, it returns error if something goes wrong.
func (p *Parser) Parse() error {
	if p.lexer == nil {
		return errors.New("lexer nil")
	}
	var err error
	p.lookAhead = p.lexer.Read()
	for p.lookAhead != TokenEOF {
		switch p.lookAhead.Text {
		case tokens[QUERY], tokens[MUTATION], tokens[LBRACE]:
			err = p.document()
			if err != nil {
				return err
			}
		case tokens[INTERFACE]:
			err = p.interfaceDefinition()
			if err != nil {
				return err
			}
		case tokens[SCALAR]:
			err = p.scalarDefinition()
			if err != nil {
				return err
			}
		case tokens[INPUT]:
			err = p.inputObjectDefinition()
			if err != nil {
				return err
			}
		case tokens[TYPE]:
			err = p.typeDefinition()
			if err != nil {
				return err
			}
		case tokens[SCHEMA]:
			err = p.schemaDefinition()
			if err != nil {
				return err
			}
		default:
			expecting := fmt.Sprintf("%s/%s/%s/%s/%s/%s/%s/%s",
				tokens[QUERY], tokens[MUTATION], tokens[LBRACE], tokens[INTERFACE],
				tokens[SCALAR], tokens[INPUT], tokens[TYPE], tokens[SCHEMA])
			return ErrBadParse{p.lexer.Line(), expecting, p.lookAhead}
		}
		p.lookAhead = p.lexer.Read()
	}
	return nil
}

func (p *Parser) interfaceDefinition() error {
	return errors.New("not implemented")
}

func (p *Parser) scalarDefinition() error {
	return errors.New("not implemented")
}

func (p *Parser) inputObjectDefinition() error {
	return errors.New("not implemented")
}

func (p *Parser) typeDefinition() error {
	return errors.New("not implemented")
}

func (p *Parser) schemaDefinition() error {
	return errors.New("not implemented")
}

func (p *Parser) document() error {
	return errors.New("not implemented")
}

func (p *Parser) definition() error {
	return errors.New("not implemented")
}

func (p *Parser) OperationDefinition() error {
	return errors.New("not implemented")
}

func (p *Parser) FragmentDefinition() error {
	return errors.New("not implemented")
}

func (p *Parser) SelectionSet() error {
	return errors.New("not implemented")
}

func (p *Parser) Selection() error {
	return errors.New("not implemented")
}

func (p *Parser) Field() error {
	return errors.New("not implemented")
}

func (p *Parser) FragmentSpread() error {
	return errors.New("not implemented")
}

func (p *Parser) InlineFragment() error {
	return errors.New("not implemented")
}

func (p *Parser) Alias() error {
	return errors.New("not implemented")
}

func (p *Parser) Arguments() error {
	return errors.New("not implemented")
}

func (p *Parser) Directives() error {
	return errors.New("not implemented")
}

func (p *Parser) Directive() error {
	return errors.New("not implemented")
}

func (p *Parser) Argument() error {
	return errors.New("not implemented")
}

func (p *Parser) TypeCondition() error {
	return errors.New("not implemented")
}

func (p *Parser) Value() error {
	return errors.New("not implemented")
}

func (p *Parser) Variable() error {
	return errors.New("not implemented")
}

func (p *Parser) BooleanValue() error {
	return errors.New("not implemented")
}

func (p *Parser) NullValue() error {
	return errors.New("not implemented")
}

func (p *Parser) EnumValue() error {
	return errors.New("not implemented")
}

func (p *Parser) ListValue() error {
	return errors.New("not implemented")
}

func (p *Parser) ObjectValue() error {
	return errors.New("not implemented")
}

func (p *Parser) ObjectField() error {
	return errors.New("not implemented")
}

func (p *Parser) VariableDefinitions() error {
	return errors.New("not implemented")
}

func (p *Parser) DefaultValue() error {
	return errors.New("not implemented")
}

func (p *Parser) Type() error {
	return errors.New("not implemented")
}

func (p *Parser) NamedType() error {
	return errors.New("not implemented")
}

func (p *Parser) ListType() error {
	return errors.New("not implemented")
}

func (p *Parser) NonNullType() error {
	return errors.New("not implemented")
}

func (p *Parser) match(k Kind) error {
	return errors.New("not implemented")
}
