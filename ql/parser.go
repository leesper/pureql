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
	input      *Lexer
	lookAheads []Token
	curr       int
}

// NewParser returns a new Parser, it returns nil if lexer nil or k <= 1.
func NewParser(l *Lexer, k int) *Parser {
	if l == nil || k <= 1 {
		return nil
	}

	p := &Parser{
		input:      l,
		lookAheads: make([]Token, k),
	}

	for i := 0; i < k; i++ {
		p.consume()
	}

	return p
}

// ParseDocument returns ast.Document.
func (p *Parser) ParseDocument() error {
	if p == nil {
		return errors.New("parser nil")
	}

	err := p.definition()
	if err != nil {
		return err
	}

	for p.lookAhead(1) != TokenEOF {
		err = p.definition()
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *Parser) definition() error {
	if p.lookAhead(1).Kind == FRAGMENT {
		return p.fragmentDefinition()
	}
	return p.operationDefinition()
}

// ParseSchema returns ast.Schema.
func (p *Parser) ParseSchema(schema string) error {
	return errors.New("not implemented")
}

// ParseOperation returns ast.Operation.
func (p *Parser) ParseOperation(oper string) error {
	return errors.New("not implemented")
}

func (p *Parser) operationDefinition() error {
	if p.lookAhead(1).Kind == LBRACE {
		return p.selectionSet()
	}

	// TODO: subscription
	if p.lookAhead(1).Kind == QUERY {
		p.match(QUERY)
	} else if p.lookAhead(1).Kind == MUTATION {
		p.match(MUTATION)
	} else {
		return ErrBadParse{
			line:   p.input.Line(),
			expect: fmt.Sprintf("%s or %s", tokens[QUERY], tokens[MUTATION]),
			found:  p.lookAhead(1),
		}
	}

	if p.lookAhead(1).Kind == NAME {
		p.match(NAME)
	}

	var err error
	if p.lookAhead(1).Kind == LPAREN {
		if err = p.variableDefinitions(); err != nil {
			return err
		}
	}

	if p.lookAhead(1).Kind == AT {
		if err = p.directives(); err != nil {
			return err
		}
	}

	return p.selectionSet()
}

func (p *Parser) variableDefinitions() error {
	err := p.match(LPAREN)
	if err != nil {
		return err
	}

	if err = p.variableDefinition(); err != nil {
		return err
	}

	if err = p.match(RPAREN); err != nil {
		return err
	}

	return nil
}

func (p *Parser) variableDefinition() error {
	err := p.variable()
	if err != nil {
		return err
	}

	if err = p.match(COLON); err != nil {
		return err
	}

	if err = p.types(); err != nil {
		return err
	}

	if p.lookAhead(1).Kind == EQL {
		return p.defaultValue()
	}

	return nil
}

func (p *Parser) types() error {
	var err error
	if p.lookAhead(1).Kind == NAME {
		if err = p.namedType(); err != nil {
			return err
		}
	} else {
		if err = p.listType(); err != nil {
			return err
		}
	}

	// non-null
	if p.lookAhead(1).Kind == BANG {
		return p.match(BANG)
	}

	return nil
}

func (p *Parser) namedType() error {
	return p.match(NAME)
}

func (p *Parser) listType() error {
	err := p.match(LBRACK)
	if err != nil {
		return err
	}

	if err = p.types(); err != nil {
		return err
	}

	if err = p.match(RBRACK); err != nil {
		return err
	}

	return nil
}

func (p *Parser) defaultValue() error {
	err := p.match(EQL)
	if err != nil {
		return err
	}

	if err = p.valueConst(); err != nil {
		return err
	}

	return nil
}

func (p *Parser) valueConst() error {
	switch p.lookAhead(1).Kind {
	case INT:
		return p.match(INT)
	case FLOAT:
		return p.match(FLOAT)
	case STRING:
		return p.match(STRING)
	case NAME:
		text := p.lookAhead(1).Text
		if text == "true" || text == "false" {
			return p.booleanValue()
		} else if text == "null" {
			return p.nullValue()
		}
		return p.enumValue()
	case LBRACK:
		return p.listValueConst()
	case LBRACE:
		return p.objectValueConst()
	default:
		expect := fmt.Sprintf("%s/%s/%s/%s/%s/%s/%s",
			tokens[INT], tokens[FLOAT], tokens[STRING], tokens[NAME],
			tokens[DOLLAR], tokens[LBRACK], tokens[LBRACE])
		return ErrBadParse{
			line:   p.input.Line(),
			expect: expect,
			found:  p.lookAhead(1),
		}
	}
}

func (p *Parser) listValueConst() error {
	err := p.match(LBRACK)
	if err != nil {
		return err
	}

	if p.lookAhead(1).Kind == RBRACK {
		return p.match(RBRACK)
	}

	if err = p.valueConst(); err != nil {
		return err
	}

	for p.lookAhead(1).Kind != RBRACK {
		if err = p.value(); err != nil {
			return err
		}
	}

	return p.match(RBRACK)
}

func (p *Parser) objectValueConst() error {
	err := p.match(LBRACE)
	if err != nil {
		return err
	}

	if p.lookAhead(1).Kind == RBRACE {
		return p.match(RBRACE)
	}

	if err = p.objectFieldConst(); err != nil {
		return err
	}

	for p.lookAhead(1).Kind != RBRACE {
		if err = p.objectFieldConst(); err != nil {
			return err
		}
	}

	return p.match(RBRACE)
}

func (p *Parser) objectFieldConst() error {
	err := p.match(NAME)
	if err != nil {
		return err
	}

	if err = p.match(COLON); err != nil {
		return err
	}

	return p.valueConst()
}

func (p *Parser) nonNullType() error {
	return errors.New("not implemented")
}

func (p *Parser) selectionSet() error {
	err := p.match(LBRACE)
	if err != nil {
		return err
	}

	if err = p.selection(); err != nil {
		return err
	}
	for p.lookAhead(1).Kind != RBRACE {
		if err = p.selection(); err != nil {
			return err
		}
	}

	if err = p.match(RBRACE); err != nil {
		return err
	}
	return nil
}

func (p *Parser) selection() error {
	if p.lookAhead(1).Kind == SPREAD {
		if p.lookAhead(2).Kind == ON {
			return p.inlineFragment()
		}
		return p.fragmentSpread()
	}

	return p.field()
}

func (p *Parser) field() error {
	var err error
	if p.lookAhead(1).Kind == NAME && p.lookAhead(2).Kind == COLON {
		if err = p.alias(); err != nil {
			return err
		}
	}

	if err = p.match(NAME); err != nil {
		return err
	}

	if p.lookAhead(1).Kind == LPAREN {
		if err = p.arguments(); err != nil {
			return err
		}
	}

	if p.lookAhead(1).Kind == AT {
		if err = p.directives(); err != nil {
			return err
		}
	}

	if p.lookAhead(1).Kind == LBRACE {
		if err = p.selectionSet(); err != nil {
			return err
		}
	}

	return nil
}

func (p *Parser) fragmentSpread() error {
	err := p.match(SPREAD)
	if err != nil {
		return err
	}

	if err = p.match(NAME); err != nil {
		return err
	}

	if p.lookAhead(1).Kind == AT {
		return p.directives()
	}

	return nil
}

func (p *Parser) inlineFragment() error {
	err := p.match(SPREAD)
	if err != nil {
		return err
	}

	if p.lookAhead(1).Kind == ON {
		if err = p.typeCondition(); err != nil {
			return err
		}
	}

	if p.lookAhead(1).Kind == AT {
		return p.directives()
	}

	return p.selectionSet()
}

func (p *Parser) typeCondition() error {
	err := p.match(ON)
	if err != nil {
		return err
	}
	return p.namedType()
}

func (p *Parser) alias() error {
	err := p.match(NAME)
	if err != nil {
		return err
	}

	if err = p.match(COLON); err != nil {
		return err
	}

	return nil
}

func (p *Parser) arguments() error {
	err := p.match(LPAREN)
	if err != nil {
		return err
	}

	if err = p.argument(); err != nil {
		return err
	}

	for p.lookAhead(1).Kind != RPAREN {
		if err = p.argument(); err != nil {
			return err
		}
	}

	if err = p.match(RPAREN); err != nil {
		return err
	}
	return nil
}

func (p *Parser) argument() error {
	err := p.match(NAME)
	if err != nil {
		return err
	}

	if err = p.match(COLON); err != nil {
		return err
	}

	return p.value()
}

func (p *Parser) value() error {
	switch p.lookAhead(1).Kind {
	case INT:
		return p.match(INT)
	case FLOAT:
		return p.match(FLOAT)
	case STRING:
		return p.match(STRING)
	case NAME:
		text := p.lookAhead(1).Text
		if text == "true" || text == "false" {
			return p.booleanValue()
		} else if text == "null" {
			return p.nullValue()
		}
		return p.enumValue()
	case DOLLAR:
		return p.variable()
	case LBRACK:
		return p.listValue()
	case LBRACE:
		return p.objectValue()
	default:
		expect := fmt.Sprintf("%s/%s/%s/%s/%s/%s/%s",
			tokens[INT], tokens[FLOAT], tokens[STRING], tokens[NAME],
			tokens[DOLLAR], tokens[LBRACK], tokens[LBRACE])
		return ErrBadParse{
			line:   p.input.Line(),
			expect: expect,
			found:  p.lookAhead(1),
		}
	}
}

func (p *Parser) variable() error {
	err := p.match(DOLLAR)
	if err != nil {
		return err
	}
	if err = p.match(NAME); err != nil {
		return err
	}
	return nil
}

func (p *Parser) booleanValue() error {
	return p.match(NAME)
}

func (p *Parser) nullValue() error {
	return p.match(NAME)
}

func (p *Parser) enumValue() error {
	return p.match(NAME)
}

func (p *Parser) listValue() error {
	err := p.match(LBRACK)
	if err != nil {
		return err
	}

	if p.lookAhead(1).Kind == RBRACK {
		return p.match(RBRACK)
	}

	if err = p.value(); err != nil {
		return err
	}

	for p.lookAhead(1).Kind != RBRACK {
		if err = p.value(); err != nil {
			return err
		}
	}

	return p.match(RBRACK)
}

func (p *Parser) objectValue() error {
	err := p.match(LBRACE)
	if err != nil {
		return err
	}

	if p.lookAhead(1).Kind == RBRACE {
		return p.match(RBRACE)
	}

	if err = p.objectField(); err != nil {
		return err
	}

	for p.lookAhead(1).Kind != RBRACE {
		if err = p.objectField(); err != nil {
			return err
		}
	}

	return p.match(RBRACE)
}

func (p *Parser) objectField() error {
	err := p.match(NAME)
	if err != nil {
		return err
	}

	if err = p.match(COLON); err != nil {
		return err
	}

	return p.value()
}

func (p *Parser) directives() error {
	err := p.directive()
	if err != nil {
		return err
	}

	for p.lookAhead(1) != TokenEOF {
		if err = p.directive(); err != nil {
			return err
		}
	}

	return nil
}

func (p *Parser) directive() error {
	err := p.match(AT)
	if err != nil {
		return err
	}

	if err = p.match(NAME); err != nil {
		return err
	}

	if p.lookAhead(1).Kind == LPAREN {
		return p.arguments()
	}

	return nil
}

func (p *Parser) fragmentDefinition() error {
	err := p.match(FRAGMENT)
	if err != nil {
		return err
	}

	if err = p.match(NAME); err != nil {
		return err
	}

	if err = p.typeCondition(); err != nil {
		return err
	}

	if p.lookAhead(1).Kind == AT {
		if err = p.directives(); err != nil {
			return err
		}
	}

	return p.selectionSet()
}

func (p *Parser) lookAhead(i int) Token {
	return p.lookAheads[(p.curr+i-1)%len(p.lookAheads)]
}

func (p *Parser) match(k Kind) error {
	if p.lookAhead(1).Kind == k {
		p.consume()
		return nil
	}
	return ErrBadParse{
		line:   p.input.Line(),
		expect: tokens[k],
		found:  p.lookAhead(1),
	}
}

func (p *Parser) consume() {
	p.lookAheads[p.curr] = p.input.Read()
	p.curr = (p.curr + 1) % len(p.lookAheads)
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

func (p *Parser) interfaceDefinition() error {
	err := p.match(INTERFACE)
	if err != nil {
		return err
	}

	if err = p.match(NAME); err != nil {
		return err
	}

	if p.lookAhead(1).Kind == AT {
		if err = p.directives(); err != nil {
			return err
		}
	}

	if err = p.match(LBRACE); err != nil {
		return err
	}

	if err = p.fieldDefinition(); err != nil {
		return err
	}

	for p.lookAhead(1).Kind != RBRACE {
		if err = p.fieldDefinition(); err != nil {
			return err
		}
	}

	return p.match(RBRACE)
}

func (p *Parser) fieldDefinition() error {
	err := p.match(NAME)
	if err != nil {
		return err
	}

	if p.lookAhead(1).Kind == LPAREN {
		if err = p.argumentsDefinition(); err != nil {
			return err
		}
	}

	if err = p.match(COLON); err != nil {
		return err
	}

	if err = p.types(); err != nil {
		return err
	}

	if p.lookAhead(1).Kind == AT {
		if err = p.directives(); err != nil {
			return err
		}
	}

	return nil
}

func (p *Parser) argumentsDefinition() error {
	err := p.match(LPAREN)
	if err != nil {
		return err
	}

	if err = p.inputValueDefinition(); err != nil {
		return err
	}

	for p.lookAhead(1).Kind != RPAREN {
		if err = p.inputValueDefinition(); err != nil {
			return err
		}
	}

	return p.match(RPAREN)
}

func (p *Parser) inputValueDefinition() error {
	err := p.match(NAME)
	if err != nil {
		return err
	}

	if err = p.match(COLON); err != nil {
		return err
	}

	if err = p.types(); err != nil {
		return err
	}

	if p.lookAhead(1).Kind == EQL {
		if err = p.defaultValue(); err != nil {
			return err
		}
	}

	if p.lookAhead(1).Kind == AT {
		if err = p.directives(); err != nil {
			return err
		}
	}

	return nil
}

func (p *Parser) scalarDefinition() error {
	err := p.match(SCALAR)
	if err != nil {
		return err
	}

	if err = p.match(NAME); err != nil {
		return err
	}

	if p.lookAhead(1).Kind == AT {
		if err = p.directives(); err != nil {
			return err
		}
	}
	return nil
}

func (p *Parser) inputObjectDefinition() error {
	err := p.match(INPUT)
	if err != nil {
		return err
	}

	if err = p.match(NAME); err != nil {
		return err
	}

	if p.lookAhead(1).Kind == AT {
		if err = p.directives(); err != nil {
			return err
		}
	}

	if err = p.match(LBRACE); err != nil {
		return err
	}

	if err = p.inputValueDefinition(); err != nil {
		return err
	}

	for p.lookAhead(1).Kind != RBRACE {
		if err = p.inputValueDefinition(); err != nil {
			return err
		}
	}

	return p.match(RBRACE)
}

func (p *Parser) typeDefinition() error {
	err := p.match(TYPE)
	if err != nil {
		return err
	}

	if err = p.match(NAME); err != nil {
		return err
	}

	if p.lookAhead(1).Kind == IMPLEMENTS {
		if err = p.implementsInterfaces(); err != nil {
			return err
		}
	}

	if p.lookAhead(1).Kind == AT {
		if err = p.directives(); err != nil {
			return err
		}
	}

	if err = p.match(LBRACE); err != nil {
		return err
	}

	if err = p.fieldDefinition(); err != nil {
		return err
	}

	for p.lookAhead(1).Kind != RBRACE {
		if err = p.fieldDefinition(); err != nil {
			return err
		}
	}
	return p.match(RBRACE)
}

func (p *Parser) implementsInterfaces() error {
	err := p.match(IMPLEMENTS)
	if err != nil {
		return err
	}

	if err = p.match(NAME); err != nil {
		return err
	}

	for p.lookAhead(1) != TokenEOF {
		if err = p.match(NAME); err != nil {
			return err
		}
	}

	return nil
}

func (p *Parser) typeExtend() error {
	err := p.match(EXTEND)
	if err != nil {
		return err
	}

	return p.typeDefinition()
}

func (p *Parser) directiveDefinition() error {
	err := p.match(DIRECTIVE)
	if err != nil {
		return err
	}

	if err = p.match(AT); err != nil {
		return err
	}

	if err = p.match(NAME); err != nil {
		return err
	}

	if p.lookAhead(1).Kind == LPAREN {
		if err = p.argumentsDefinition(); err != nil {
			return err
		}
	}

	if err = p.match(ON); err != nil {
		return err
	}

	return p.directiveLocations()
}

func (p *Parser) directiveLocations() error {
	err := p.directiveLocation()
	if err != nil {
		return err
	}

	for p.lookAhead(1).Kind == PIPE {
		if err = p.match(PIPE); err != nil {
			return err
		}
		if err = p.directiveLocation(); err != nil {
			return err
		}
	}

	return nil
}

func (p *Parser) directiveLocation() error {
	return p.match(NAME)
}

func (p *Parser) schemaDefinition() error {
	err := p.match(SCHEMA)
	if err != nil {
		return err
	}

	if p.lookAhead(1).Kind == AT {
		if err = p.directives(); err != nil {
			return err
		}
	}

	if err = p.match(LBRACE); err != nil {
		return err
	}

	if err = p.operationTypeDefinition(); err != nil {
		return err
	}

	for p.lookAhead(1).Kind != RBRACE {
		if err = p.operationTypeDefinition(); err != nil {
			return err
		}
	}

	return p.match(RBRACE)
}

func (p *Parser) operationTypeDefinition() error {
	if p.lookAhead(1).Kind == QUERY {
		p.match(QUERY)
	} else if p.lookAhead(1).Kind == MUTATION {
		p.match(MUTATION)
	} else {
		return ErrBadParse{
			line:   p.input.Line(),
			expect: fmt.Sprintf("%s or %s", tokens[QUERY], tokens[MUTATION]),
			found:  p.lookAhead(1),
		}
	}

	if err := p.match(COLON); err != nil {
		return err
	}

	return p.match(NAME)
}

func (p *Parser) enumType() error {
	err := p.match(ENUM)
	if err != nil {
		return err
	}

	if err = p.match(NAME); err != nil {
		return err
	}

	if p.lookAhead(1).Kind == AT {
		if err = p.directives(); err != nil {
			return err
		}
	}

	if err = p.match(LBRACE); err != nil {
		return err
	}

	if err = p.enumValue(); err != nil {
		return err
	}

	for p.lookAhead(1).Kind == RBRACE {
		if err = p.enumValue(); err != nil {
			return err
		}
	}

	return p.match(RBRACE)
}

func (p *Parser) unionDefinition() error {
	err := p.match(UNION)
	if err != nil {
		return err
	}

	if err = p.match(NAME); err != nil {
		return err
	}

	if p.lookAhead(1).Kind == AT {
		if err = p.directives(); err != nil {
			return err
		}
	}

	if err = p.match(EQL); err != nil {
		return err
	}

	return p.unionMembers()
}

func (p *Parser) unionMembers() error {
	err := p.match(NAME)
	if err != nil {
		return err
	}

	for p.lookAhead(1).Kind == PIPE {
		if err = p.match(PIPE); err != nil {
			return err
		}

		if err = p.match(NAME); err != nil {
			return err
		}
	}

	return nil
}
