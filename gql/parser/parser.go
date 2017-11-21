package parser

import (
	"errors"
	"fmt"

	"github.com/leesper/pureql/gql/token"
)

// ErrBadParse for invalid parse.
type ErrBadParse struct {
	line   int
	expect string
	found  token.Token
}

func (e ErrBadParse) Error() string {
	return fmt.Sprintf("line %d: expecting %s, found %s", e.line, e.expect, e.found)
}

// ParseDocument returns ast.Document.
func ParseDocument(document string) error {
	return newParser(newLexer(document), 2).parseDocument()
}

// ParseSchema returns ast.Schema.
func ParseSchema(schema string) error {
	return newParser(newLexer(schema), 2).parseSchema()
}

// Parser converts GraphQL source into AST.
type parser struct {
	input      *lexer
	lookAheads []token.Token
	curr       int
}

func newParser(l *lexer, k int) *parser {
	if l == nil || k <= 1 {
		return nil
	}

	p := &parser{
		input:      l,
		lookAheads: make([]token.Token, k),
	}

	for i := 0; i < k; i++ {
		p.consume()
	}

	return p
}

func (p *parser) parseDocument() error {
	if p == nil {
		return errors.New("parser nil")
	}

	err := p.definition()
	if err != nil {
		return err
	}

	for p.lookAhead(1) != token.TokenEOF {
		err = p.definition()
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) definition() error {
	if p.lookAhead(1).Text == token.TokenString(token.FRAGMENT) {
		return p.fragmentDefinition()
	}
	return p.operationDefinition()
}

func (p *parser) parseSchema() error {
	if p == nil {
		return errors.New("parser nil")
	}

	err := p.schema()
	if err != nil {
		return err
	}

	for p.lookAhead(1) != token.TokenEOF {
		err = p.schema()
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) schema() error {
	switch p.lookAhead(1).Text {
	case token.TokenString(token.INTERFACE):
		return p.interfaceDefinition()
	case token.TokenString(token.SCALAR):
		return p.scalarDefinition()
	case token.TokenString(token.INPUT):
		return p.inputObjectDefinition()
	case token.TokenString(token.TYPE):
		return p.typeDefinition()
	case token.TokenString(token.EXTEND):
		return p.extendDefinition()
	case token.TokenString(token.DIRECTIVE):
		return p.directiveDefinition()
	case token.TokenString(token.SCHEMA):
		return p.schemaDefinition()
	case token.TokenString(token.ENUM):
		return p.enumDefinition()
	default:
		return p.unionDefinition()
	}
}

func (p *parser) operationDefinition() error {
	if p.lookAhead(1).Kind == token.LBRACE {
		return p.selectionSet()
	}

	// TODO: subscription
	var err error
	if err = p.match(token.QUERY); err != nil {
		if err = p.match(token.MUTATION); err != nil {
			if err = p.match(token.SUBSCRIPTION); err != nil {
				return ErrBadParse{
					line: p.input.Line(),
					expect: fmt.Sprintf("%s or %s or %s",
						token.TokenString(token.QUERY),
						token.TokenString(token.MUTATION),
						token.TokenString(token.SUBSCRIPTION)),
					found: p.lookAhead(1),
				}
			}
		}
	}

	if p.lookAhead(1).Kind == token.NAME {
		p.match(token.NAME)
	}

	if p.lookAhead(1).Kind == token.LPAREN {
		if err = p.variableDefinitions(); err != nil {
			return err
		}
	}

	if p.lookAhead(1).Kind == token.AT {
		if err = p.directives(); err != nil {
			return err
		}
	}

	return p.selectionSet()
}

func (p *parser) variableDefinitions() error {
	err := p.match(token.LPAREN)
	if err != nil {
		return err
	}

	if err = p.variableDefinition(); err != nil {
		return err
	}

	for p.lookAhead(1).Kind != token.RPAREN {
		if err = p.variableDefinition(); err != nil {
			return err
		}
	}

	if err = p.match(token.RPAREN); err != nil {
		return err
	}

	return nil
}

func (p *parser) variableDefinition() error {
	err := p.variable()
	if err != nil {
		return err
	}

	if err = p.match(token.COLON); err != nil {
		return err
	}

	if err = p.types(); err != nil {
		return err
	}

	if p.lookAhead(1).Kind == token.EQL {
		return p.defaultValue()
	}

	return nil
}

func (p *parser) types() error {
	var err error
	if p.lookAhead(1).Kind == token.NAME {
		if err = p.namedType(); err != nil {
			return err
		}
	} else {
		if err = p.listType(); err != nil {
			return err
		}
	}

	// non-null
	if p.lookAhead(1).Kind == token.BANG {
		return p.match(token.BANG)
	}

	return nil
}

func (p *parser) namedType() error {
	return p.match(token.NAME)
}

func (p *parser) listType() error {
	err := p.match(token.LBRACK)
	if err != nil {
		return err
	}

	if err = p.types(); err != nil {
		return err
	}

	if err = p.match(token.RBRACK); err != nil {
		return err
	}

	return nil
}

func (p *parser) defaultValue() error {
	err := p.match(token.EQL)
	if err != nil {
		return err
	}

	if err = p.valueConst(); err != nil {
		return err
	}

	return nil
}

func (p *parser) valueConst() error {
	switch p.lookAhead(1).Kind {
	case token.INT:
		return p.match(token.INT)
	case token.FLOAT:
		return p.match(token.FLOAT)
	case token.STRING:
		return p.match(token.STRING)
	case token.NAME:
		text := p.lookAhead(1).Text
		if text == "true" || text == "false" {
			return p.booleanValue()
		} else if text == "null" {
			return p.nullValue()
		}
		return p.enumValue()
	case token.LBRACK:
		return p.listValueConst()
	case token.LBRACE:
		return p.objectValueConst()
	default:
		expect := fmt.Sprintf("%s/%s/%s/%s/%s/%s/%s",
			token.TokenString(token.INT),
			token.TokenString(token.FLOAT),
			token.TokenString(token.STRING),
			token.TokenString(token.NAME),
			token.TokenString(token.DOLLAR),
			token.TokenString(token.LBRACK),
			token.TokenString(token.LBRACE))
		return ErrBadParse{
			line:   p.input.Line(),
			expect: expect,
			found:  p.lookAhead(1),
		}
	}
}

func (p *parser) listValueConst() error {
	err := p.match(token.LBRACK)
	if err != nil {
		return err
	}

	if p.lookAhead(1).Kind == token.RBRACK {
		return p.match(token.RBRACK)
	}

	if err = p.valueConst(); err != nil {
		return err
	}

	for p.lookAhead(1).Kind != token.RBRACK {
		if err = p.value(); err != nil {
			return err
		}
	}

	return p.match(token.RBRACK)
}

func (p *parser) objectValueConst() error {
	err := p.match(token.LBRACE)
	if err != nil {
		return err
	}

	if p.lookAhead(1).Kind == token.RBRACE {
		return p.match(token.RBRACE)
	}

	if err = p.objectFieldConst(); err != nil {
		return err
	}

	for p.lookAhead(1).Kind != token.RBRACE {
		if err = p.objectFieldConst(); err != nil {
			return err
		}
	}

	return p.match(token.RBRACE)
}

func (p *parser) objectFieldConst() error {
	err := p.match(token.NAME)
	if err != nil {
		return err
	}

	if err = p.match(token.COLON); err != nil {
		return err
	}

	return p.valueConst()
}

func (p *parser) nonNullType() error {
	return errors.New("not implemented")
}

func (p *parser) selectionSet() error {
	err := p.match(token.LBRACE)
	if err != nil {
		return err
	}

	if err = p.selection(); err != nil {
		return err
	}
	for p.lookAhead(1).Kind != token.RBRACE {
		if err = p.selection(); err != nil {
			return err
		}
	}

	if err = p.match(token.RBRACE); err != nil {
		return err
	}
	return nil
}

func (p *parser) selection() error {
	if p.lookAhead(1).Kind == token.SPREAD {
		if p.lookAhead(2).Kind == token.NAME && p.lookAhead(2).Text != token.TokenString(token.ON) {
			return p.fragmentSpread()
		}
		return p.inlineFragment()
	}

	return p.field()
}

func (p *parser) field() error {
	var err error
	if p.lookAhead(1).Kind == token.NAME && p.lookAhead(2).Kind == token.COLON {
		if err = p.alias(); err != nil {
			return err
		}
	}

	if err = p.match(token.NAME); err != nil {
		return err
	}

	if p.lookAhead(1).Kind == token.LPAREN {
		if err = p.arguments(); err != nil {
			return err
		}
	}

	if p.lookAhead(1).Kind == token.AT {
		if err = p.directives(); err != nil {
			return err
		}
	}

	if p.lookAhead(1).Kind == token.LBRACE {
		if err = p.selectionSet(); err != nil {
			return err
		}
	}

	return nil
}

func (p *parser) fragmentSpread() error {
	err := p.match(token.SPREAD)
	if err != nil {
		return err
	}

	if p.lookAhead(1).Text == token.TokenString(token.ON) {
		return ErrBadParse{
			line:   p.input.Line(),
			expect: "NAME but not *on*",
			found:  p.lookAhead(1),
		}
	}

	if err = p.match(token.NAME); err != nil {
		return err
	}

	if p.lookAhead(1).Kind == token.AT {
		return p.directives()
	}

	return nil
}

func (p *parser) inlineFragment() error {
	err := p.match(token.SPREAD)
	if err != nil {
		return err
	}

	if p.lookAhead(1).Text == token.TokenString(token.ON) {
		if err = p.typeCondition(); err != nil {
			return err
		}
	}

	if p.lookAhead(1).Kind == token.AT {
		if err = p.directives(); err != nil {
			return err
		}
	}

	return p.selectionSet()
}

func (p *parser) typeCondition() error {
	err := p.match(token.ON)
	if err != nil {
		return err
	}
	return p.namedType()
}

func (p *parser) alias() error {
	err := p.match(token.NAME)
	if err != nil {
		return err
	}

	if err = p.match(token.COLON); err != nil {
		return err
	}

	return nil
}

func (p *parser) arguments() error {
	err := p.match(token.LPAREN)
	if err != nil {
		return err
	}

	if err = p.argument(); err != nil {
		return err
	}

	for p.lookAhead(1).Kind != token.RPAREN {
		if err = p.argument(); err != nil {
			return err
		}
	}

	if err = p.match(token.RPAREN); err != nil {
		return err
	}
	return nil
}

func (p *parser) argument() error {
	err := p.match(token.NAME)
	if err != nil {
		return err
	}

	if err = p.match(token.COLON); err != nil {
		return err
	}

	return p.value()
}

func (p *parser) value() error {
	switch p.lookAhead(1).Kind {
	case token.INT:
		return p.match(token.INT)
	case token.FLOAT:
		return p.match(token.FLOAT)
	case token.STRING:
		return p.match(token.STRING)
	case token.NAME:
		text := p.lookAhead(1).Text
		if text == "true" || text == "false" {
			return p.booleanValue()
		} else if text == "null" {
			return p.nullValue()
		}
		return p.enumValue()
	case token.DOLLAR:
		return p.variable()
	case token.LBRACK:
		return p.listValue()
	case token.LBRACE:
		return p.objectValue()
	default:
		expect := fmt.Sprintf("%s/%s/%s/%s/%s/%s/%s",
			token.TokenString(token.INT),
			token.TokenString(token.FLOAT),
			token.TokenString(token.STRING),
			token.TokenString(token.NAME),
			token.TokenString(token.DOLLAR),
			token.TokenString(token.LBRACK),
			token.TokenString(token.LBRACE))
		return ErrBadParse{
			line:   p.input.Line(),
			expect: expect,
			found:  p.lookAhead(1),
		}
	}
}

func (p *parser) variable() error {
	err := p.match(token.DOLLAR)
	if err != nil {
		return err
	}
	if err = p.match(token.NAME); err != nil {
		return err
	}
	return nil
}

// FIXME: check whether its true/false/null
func (p *parser) booleanValue() error {
	return p.match(token.NAME)
}

func (p *parser) nullValue() error {
	return p.match(token.NAME)
}

func (p *parser) enumValue() error {
	err := p.match(token.NAME)
	if err != nil {
		return err
	}

	if p.lookAhead(1).Kind == token.AT {
		if err = p.directives(); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) listValue() error {
	err := p.match(token.LBRACK)
	if err != nil {
		return err
	}

	if p.lookAhead(1).Kind == token.RBRACK {
		return p.match(token.RBRACK)
	}

	if err = p.value(); err != nil {
		return err
	}

	for p.lookAhead(1).Kind != token.RBRACK {
		if err = p.value(); err != nil {
			return err
		}
	}

	return p.match(token.RBRACK)
}

func (p *parser) objectValue() error {
	err := p.match(token.LBRACE)
	if err != nil {
		return err
	}

	if p.lookAhead(1).Kind == token.RBRACE {
		return p.match(token.RBRACE)
	}

	if err = p.objectField(); err != nil {
		return err
	}

	for p.lookAhead(1).Kind != token.RBRACE {
		if err = p.objectField(); err != nil {
			return err
		}
	}

	return p.match(token.RBRACE)
}

func (p *parser) objectField() error {
	err := p.match(token.NAME)
	if err != nil {
		return err
	}

	if err = p.match(token.COLON); err != nil {
		return err
	}

	return p.value()
}

func (p *parser) directives() error {
	err := p.directive()
	if err != nil {
		return err
	}

	for p.lookAhead(1).Kind == token.AT {
		if err = p.directive(); err != nil {
			return err
		}
	}

	return nil
}

func (p *parser) directive() error {
	err := p.match(token.AT)
	if err != nil {
		return err
	}

	if err = p.match(token.NAME); err != nil {
		return err
	}

	if p.lookAhead(1).Kind == token.LPAREN {
		return p.arguments()
	}

	return nil
}

func (p *parser) fragmentDefinition() error {
	err := p.match(token.FRAGMENT)
	if err != nil {
		return err
	}

	if p.lookAhead(1).Text == token.TokenString(token.ON) {
		return ErrBadParse{
			line:   p.input.Line(),
			expect: "NAME but not *on*",
			found:  p.lookAhead(1),
		}
	}

	if err = p.match(token.NAME); err != nil {
		return err
	}

	if err = p.typeCondition(); err != nil {
		return err
	}

	if p.lookAhead(1).Kind == token.AT {
		if err = p.directives(); err != nil {
			return err
		}
	}

	return p.selectionSet()
}

func (p *parser) lookAhead(i int) token.Token {
	return p.lookAheads[(p.curr+i-1)%len(p.lookAheads)]
}

func (p *parser) match(k token.Kind) error {
	// fmt.Println("DEBUG tok", p.lookAhead(1))
	if token.IsKeyword(k) {
		if p.lookAhead(1).Kind == token.NAME && p.lookAhead(1).Text == token.TokenString(k) {
			p.consume()
			return nil
		}
	} else {
		if p.lookAhead(1).Kind == k {
			p.consume()
			return nil
		}
	}

	return ErrBadParse{
		line:   p.input.Line(),
		expect: token.TokenString(k),
		found:  p.lookAhead(1),
	}
}

func (p *parser) consume() {
	p.lookAheads[p.curr] = p.input.Read()
	p.curr = (p.curr + 1) % len(p.lookAheads)
}

func (p *parser) interfaceDefinition() error {
	err := p.match(token.INTERFACE)
	if err != nil {
		return err
	}

	if err = p.match(token.NAME); err != nil {
		return err
	}

	if p.lookAhead(1).Kind == token.AT {
		if err = p.directives(); err != nil {
			return err
		}
	}

	if err = p.match(token.LBRACE); err != nil {
		return err
	}

	if err = p.fieldDefinition(); err != nil {
		return err
	}

	for p.lookAhead(1).Kind != token.RBRACE {
		if err = p.fieldDefinition(); err != nil {
			return err
		}
	}

	return p.match(token.RBRACE)
}

func (p *parser) fieldDefinition() error {
	err := p.match(token.NAME)
	if err != nil {
		return err
	}

	if p.lookAhead(1).Kind == token.LPAREN {
		if err = p.argumentsDefinition(); err != nil {
			return err
		}
	}

	if err = p.match(token.COLON); err != nil {
		return err
	}

	if err = p.types(); err != nil {
		return err
	}

	if p.lookAhead(1).Kind == token.AT {
		if err = p.directives(); err != nil {
			return err
		}
	}

	return nil
}

func (p *parser) argumentsDefinition() error {
	err := p.match(token.LPAREN)
	if err != nil {
		return err
	}

	if err = p.inputValueDefinition(); err != nil {
		return err
	}

	for p.lookAhead(1).Kind != token.RPAREN {
		if err = p.inputValueDefinition(); err != nil {
			return err
		}
	}

	return p.match(token.RPAREN)
}

func (p *parser) inputValueDefinition() error {
	err := p.match(token.NAME)
	if err != nil {
		return err
	}

	if err = p.match(token.COLON); err != nil {
		return err
	}

	if err = p.types(); err != nil {
		return err
	}

	if p.lookAhead(1).Kind == token.EQL {
		if err = p.defaultValue(); err != nil {
			return err
		}
	}

	if p.lookAhead(1).Kind == token.AT {
		if err = p.directives(); err != nil {
			return err
		}
	}

	return nil
}

func (p *parser) scalarDefinition() error {
	err := p.match(token.SCALAR)
	if err != nil {
		return err
	}

	if err = p.match(token.NAME); err != nil {
		return err
	}

	if p.lookAhead(1).Kind == token.AT {
		if err = p.directives(); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) inputObjectDefinition() error {
	err := p.match(token.INPUT)
	if err != nil {
		return err
	}

	if err = p.match(token.NAME); err != nil {
		return err
	}

	if p.lookAhead(1).Kind == token.AT {
		if err = p.directives(); err != nil {
			return err
		}
	}

	if err = p.match(token.LBRACE); err != nil {
		return err
	}

	if err = p.inputValueDefinition(); err != nil {
		return err
	}

	for p.lookAhead(1).Kind != token.RBRACE {
		if err = p.inputValueDefinition(); err != nil {
			return err
		}
	}

	return p.match(token.RBRACE)
}

func (p *parser) typeDefinition() error {
	err := p.match(token.TYPE)
	if err != nil {
		return err
	}

	if err = p.match(token.NAME); err != nil {
		return err
	}

	if p.lookAhead(1).Text == token.TokenString(token.IMPLEMENTS) {
		if err = p.implementsInterfaces(); err != nil {
			return err
		}
	}

	if p.lookAhead(1).Kind == token.AT {
		if err = p.directives(); err != nil {
			return err
		}
	}

	if err = p.match(token.LBRACE); err != nil {
		return err
	}

	// if err = p.fieldDefinition(); err != nil {
	// 	return err
	// }

	for p.lookAhead(1).Kind != token.RBRACE {
		if err = p.fieldDefinition(); err != nil {
			return err
		}
	}
	return p.match(token.RBRACE)
}

func (p *parser) implementsInterfaces() error {
	err := p.match(token.IMPLEMENTS)
	if err != nil {
		return err
	}

	if err = p.match(token.NAME); err != nil {
		return err
	}

	for p.lookAhead(1).Kind == token.NAME {
		if err = p.match(token.NAME); err != nil {
			return err
		}
	}

	return nil
}

func (p *parser) extendDefinition() error {
	err := p.match(token.EXTEND)
	if err != nil {
		return err
	}

	return p.typeDefinition()
}

func (p *parser) directiveDefinition() error {
	err := p.match(token.DIRECTIVE)
	if err != nil {
		return err
	}

	if err = p.match(token.AT); err != nil {
		return err
	}

	if err = p.match(token.NAME); err != nil {
		return err
	}

	if p.lookAhead(1).Kind == token.LPAREN {
		if err = p.argumentsDefinition(); err != nil {
			return err
		}
	}

	if err = p.match(token.ON); err != nil {
		return err
	}

	return p.directiveLocations()
}

func (p *parser) directiveLocations() error {
	err := p.directiveLocation()
	if err != nil {
		return err
	}

	for p.lookAhead(1).Kind == token.PIPE {
		if err = p.match(token.PIPE); err != nil {
			return err
		}
		if err = p.directiveLocation(); err != nil {
			return err
		}
	}

	return nil
}

func (p *parser) directiveLocation() error {
	return p.match(token.NAME)
}

func (p *parser) schemaDefinition() error {
	err := p.match(token.SCHEMA)
	if err != nil {
		return err
	}

	if p.lookAhead(1).Kind == token.AT {
		if err = p.directives(); err != nil {
			return err
		}
	}

	if err = p.match(token.LBRACE); err != nil {
		return err
	}

	if err = p.operationTypeDefinition(); err != nil {
		return err
	}

	for p.lookAhead(1).Kind != token.RBRACE {
		if err = p.operationTypeDefinition(); err != nil {
			return err
		}
	}

	return p.match(token.RBRACE)
}

func (p *parser) operationTypeDefinition() error {
	var err error
	if err = p.match(token.QUERY); err != nil {
		if err = p.match(token.MUTATION); err != nil {
			if err = p.match(token.SUBSCRIPTION); err != nil {
				return ErrBadParse{
					line: p.input.Line(),
					expect: fmt.Sprintf("%s or %s or %s",
						token.TokenString(token.QUERY),
						token.TokenString(token.MUTATION),
						token.TokenString(token.SUBSCRIPTION)),
					found: p.lookAhead(1),
				}
			}
		}
	}

	if err := p.match(token.COLON); err != nil {
		return err
	}

	return p.match(token.NAME)
}

func (p *parser) enumDefinition() error {
	err := p.match(token.ENUM)
	if err != nil {
		return err
	}

	if err = p.match(token.NAME); err != nil {
		return err
	}

	if p.lookAhead(1).Kind == token.AT {
		if err = p.directives(); err != nil {
			return err
		}
	}

	if err = p.match(token.LBRACE); err != nil {
		return err
	}

	if err = p.enumValue(); err != nil {
		return err
	}

	for p.lookAhead(1).Kind != token.RBRACE {
		if err = p.enumValue(); err != nil {
			return err
		}
	}

	return p.match(token.RBRACE)
}

func (p *parser) unionDefinition() error {
	err := p.match(token.UNION)
	if err != nil {
		return err
	}

	if err = p.match(token.NAME); err != nil {
		return err
	}

	if p.lookAhead(1).Kind == token.AT {
		if err = p.directives(); err != nil {
			return err
		}
	}

	if err = p.match(token.EQL); err != nil {
		return err
	}

	return p.unionMembers()
}

func (p *parser) unionMembers() error {
	err := p.match(token.NAME)
	if err != nil {
		return err
	}

	for p.lookAhead(1).Kind == token.PIPE {
		if err = p.match(token.PIPE); err != nil {
			return err
		}

		if err = p.match(token.NAME); err != nil {
			return err
		}
	}

	return nil
}
