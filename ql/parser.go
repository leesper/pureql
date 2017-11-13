package ql

import "errors"

type Parser struct {
	lexer     *Lexer
	lookAhead Token
}

func NewParser(l Lexer) *Parser {
	return nil
}

func (p *Parser) Parse() error {
	return errors.New("not implemented")
}

func (p *Parser) Document() error {
	return errors.New("not implemented")
}

func (p *Parser) Definition() error {
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

func (p *Parser) InterfaceDefinition() error {
	return errors.New("not implemented")
}

func (p *Parser) ScalarDefinition() error {
	return errors.New("not implemented")
}

func (p *Parser) InputObjectDefinition() error {
	return errors.New("not implemented")
}

func (p *Parser) TypeDefinition() error {
	return errors.New("not implemented")
}

func (p *Parser) SchemaDefinition() error {
	return errors.New("not implemented")
}

func (p *Parser) match(k Kind) error {
	return errors.New("not implemented")
}
