package parser

import (
	"errors"
	"fmt"
	"strings"
)

// ErrBadParse for invalid parse.
type ErrBadParse struct {
	line   int
	expect string
	found  Token
}

func (e ErrBadParse) Error() string {
	return fmt.Sprintf("line %d: expecting %s, found %s", e.line, e.expect, e.found)
}

// ParseDocument returns ast.Document.
func ParseDocument(document []byte) error {
	document = []byte(strings.TrimRight(string(document), "\n\t\r "))
	return newParser(newLexer(document), 2).parseDocument()
}

// ParseSchema returns ast.Schema.
func ParseSchema(schema []byte) error {
	schema = []byte(strings.TrimRight(string(schema), "\n\t\r "))
	return newParser(newLexer(schema), 2).parseSchema()
}

// Parser converts GraphQL source into AST.
type parser struct {
	input      *lexer
	lookAheads []Token
	curr       int
}

func newParser(l *lexer, k int) *parser {
	if l == nil || k <= 1 {
		return nil
	}

	p := &parser{
		input:      l,
		lookAheads: make([]Token, k),
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

	for p.lookAhead(1) != TokenEOF {
		err = p.definition()
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) definition() error {
	if p.lookAhead(1).Text == Stringify(FRAGMENT) {
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

	for p.lookAhead(1) != TokenEOF {
		err = p.schema()
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) schema() error {
	switch p.lookAhead(1).Text {
	case Stringify(INTERFACE):
		return p.interfaceDefinition()
	case Stringify(SCALAR):
		return p.scalarDefinition()
	case Stringify(INPUT):
		return p.inputObjectDefinition()
	case Stringify(TYPE):
		return p.typeDefinition()
	case Stringify(EXTEND):
		return p.extendDefinition()
	case Stringify(DIRECTIVE):
		return p.directiveDefinition()
	case Stringify(SCHEMA):
		return p.schemaDefinition()
	case Stringify(ENUM):
		return p.enumDefinition()
	default:
		return p.unionDefinition()
	}
}

func (p *parser) operationDefinition() error {
	if p.lookAhead(1).Kind == LBRACE {
		return p.selectionSet()
	}

	// TODO: subscription
	var err error
	if err = p.match(QUERY); err != nil {
		if err = p.match(MUTATION); err != nil {
			if err = p.match(SUBSCRIPTION); err != nil {
				return ErrBadParse{
					line: p.input.Line(),
					expect: fmt.Sprintf("%s or %s or %s",
						Stringify(QUERY),
						Stringify(MUTATION),
						Stringify(SUBSCRIPTION)),
					found: p.lookAhead(1),
				}
			}
		}
	}

	if p.lookAhead(1).Kind == NAME {
		p.match(NAME)
	}

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

func (p *parser) variableDefinitions() error {
	err := p.match(LPAREN)
	if err != nil {
		return err
	}

	if err = p.variableDefinition(); err != nil {
		return err
	}

	for p.lookAhead(1).Kind != RPAREN {
		if err = p.variableDefinition(); err != nil {
			return err
		}
	}

	if err = p.match(RPAREN); err != nil {
		return err
	}

	return nil
}

func (p *parser) variableDefinition() error {
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

func (p *parser) types() error {
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

func (p *parser) namedType() error {
	return p.match(NAME)
}

func (p *parser) listType() error {
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

func (p *parser) defaultValue() error {
	err := p.match(EQL)
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
			Stringify(INT),
			Stringify(FLOAT),
			Stringify(STRING),
			Stringify(NAME),
			Stringify(DOLLAR),
			Stringify(LBRACK),
			Stringify(LBRACE))
		return ErrBadParse{
			line:   p.input.Line(),
			expect: expect,
			found:  p.lookAhead(1),
		}
	}
}

func (p *parser) listValueConst() error {
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

func (p *parser) objectValueConst() error {
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

func (p *parser) objectFieldConst() error {
	err := p.match(NAME)
	if err != nil {
		return err
	}

	if err = p.match(COLON); err != nil {
		return err
	}

	return p.valueConst()
}

func (p *parser) nonNullType() error {
	return errors.New("not implemented")
}

func (p *parser) selectionSet() error {
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

func (p *parser) selection() error {
	if p.lookAhead(1).Kind == SPREAD {
		if p.lookAhead(2).Kind == NAME && p.lookAhead(2).Text != Stringify(ON) {
			return p.fragmentSpread()
		}
		return p.inlineFragment()
	}

	return p.field()
}

func (p *parser) field() error {
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

func (p *parser) fragmentSpread() error {
	err := p.match(SPREAD)
	if err != nil {
		return err
	}

	if p.lookAhead(1).Text == Stringify(ON) {
		return ErrBadParse{
			line:   p.input.Line(),
			expect: "NAME but not *on*",
			found:  p.lookAhead(1),
		}
	}

	if err = p.match(NAME); err != nil {
		return err
	}

	if p.lookAhead(1).Kind == AT {
		return p.directives()
	}

	return nil
}

func (p *parser) inlineFragment() error {
	err := p.match(SPREAD)
	if err != nil {
		return err
	}

	if p.lookAhead(1).Text == Stringify(ON) {
		if err = p.typeCondition(); err != nil {
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

func (p *parser) typeCondition() error {
	err := p.match(ON)
	if err != nil {
		return err
	}
	return p.namedType()
}

func (p *parser) alias() error {
	err := p.match(NAME)
	if err != nil {
		return err
	}

	if err = p.match(COLON); err != nil {
		return err
	}

	return nil
}

func (p *parser) arguments() error {
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

func (p *parser) argument() error {
	err := p.match(NAME)
	if err != nil {
		return err
	}

	if err = p.match(COLON); err != nil {
		return err
	}

	return p.value()
}

func (p *parser) value() error {
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
			Stringify(INT),
			Stringify(FLOAT),
			Stringify(STRING),
			Stringify(NAME),
			Stringify(DOLLAR),
			Stringify(LBRACK),
			Stringify(LBRACE))
		return ErrBadParse{
			line:   p.input.Line(),
			expect: expect,
			found:  p.lookAhead(1),
		}
	}
}

func (p *parser) variable() error {
	err := p.match(DOLLAR)
	if err != nil {
		return err
	}
	if err = p.match(NAME); err != nil {
		return err
	}
	return nil
}

// FIXME: check whether its true/false/null
func (p *parser) booleanValue() error {
	return p.match(NAME)
}

func (p *parser) nullValue() error {
	return p.match(NAME)
}

func (p *parser) enumValue() error {
	err := p.match(NAME)
	if err != nil {
		return err
	}

	if p.lookAhead(1).Kind == AT {
		if err = p.directives(); err != nil {
			return err
		}
	}
	return nil
}

func (p *parser) listValue() error {
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

func (p *parser) objectValue() error {
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

func (p *parser) objectField() error {
	err := p.match(NAME)
	if err != nil {
		return err
	}

	if err = p.match(COLON); err != nil {
		return err
	}

	return p.value()
}

func (p *parser) directives() error {
	err := p.directive()
	if err != nil {
		return err
	}

	for p.lookAhead(1).Kind == AT {
		if err = p.directive(); err != nil {
			return err
		}
	}

	return nil
}

func (p *parser) directive() error {
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

func (p *parser) fragmentDefinition() error {
	err := p.match(FRAGMENT)
	if err != nil {
		return err
	}

	if p.lookAhead(1).Text == Stringify(ON) {
		return ErrBadParse{
			line:   p.input.Line(),
			expect: "NAME but not *on*",
			found:  p.lookAhead(1),
		}
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

func (p *parser) lookAhead(i int) Token {
	return p.lookAheads[(p.curr+i-1)%len(p.lookAheads)]
}

func (p *parser) match(k Kind) error {
	// fmt.Println("DEBUG tok", p.lookAhead(1))
	if IsKeyword(k) {
		if p.lookAhead(1).Kind == NAME && p.lookAhead(1).Text == Stringify(k) {
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
		expect: Stringify(k),
		found:  p.lookAhead(1),
	}
}

func (p *parser) consume() {
	p.lookAheads[p.curr] = p.input.Read()
	p.curr = (p.curr + 1) % len(p.lookAheads)
}

func (p *parser) interfaceDefinition() error {
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

func (p *parser) fieldDefinition() error {
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

func (p *parser) argumentsDefinition() error {
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

func (p *parser) inputValueDefinition() error {
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

func (p *parser) scalarDefinition() error {
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

func (p *parser) inputObjectDefinition() error {
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

func (p *parser) typeDefinition() error {
	err := p.match(TYPE)
	if err != nil {
		return err
	}

	if err = p.match(NAME); err != nil {
		return err
	}

	if p.lookAhead(1).Text == Stringify(IMPLEMENTS) {
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

	// if err = p.fieldDefinition(); err != nil {
	// 	return err
	// }

	for p.lookAhead(1).Kind != RBRACE {
		if err = p.fieldDefinition(); err != nil {
			return err
		}
	}
	return p.match(RBRACE)
}

func (p *parser) implementsInterfaces() error {
	err := p.match(IMPLEMENTS)
	if err != nil {
		return err
	}

	if err = p.match(NAME); err != nil {
		return err
	}

	for p.lookAhead(1).Kind == NAME {
		if err = p.match(NAME); err != nil {
			return err
		}
	}

	return nil
}

func (p *parser) extendDefinition() error {
	err := p.match(EXTEND)
	if err != nil {
		return err
	}

	return p.typeDefinition()
}

func (p *parser) directiveDefinition() error {
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

func (p *parser) directiveLocations() error {
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

func (p *parser) directiveLocation() error {
	return p.match(NAME)
}

func (p *parser) schemaDefinition() error {
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

func (p *parser) operationTypeDefinition() error {
	var err error
	if err = p.match(QUERY); err != nil {
		if err = p.match(MUTATION); err != nil {
			if err = p.match(SUBSCRIPTION); err != nil {
				return ErrBadParse{
					line: p.input.Line(),
					expect: fmt.Sprintf("%s or %s or %s",
						Stringify(QUERY),
						Stringify(MUTATION),
						Stringify(SUBSCRIPTION)),
					found: p.lookAhead(1),
				}
			}
		}
	}

	if err := p.match(COLON); err != nil {
		return err
	}

	return p.match(NAME)
}

func (p *parser) enumDefinition() error {
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

	for p.lookAhead(1).Kind != RBRACE {
		if err = p.enumValue(); err != nil {
			return err
		}
	}

	return p.match(RBRACE)
}

func (p *parser) unionDefinition() error {
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

func (p *parser) unionMembers() error {
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
