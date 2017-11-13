package ql

type Parser struct {
	lexer *Lexer
}

func (p *Parser) Document() error { return nil }

func (p *Parser) OperationDefinition() error { return nil }

func (p *Parser) FragmentDefinition() error { return nil }

func (p *Parser) SelectionSet() error { return nil }

func (p *Parser) Selection() error { return nil }

func (p *Parser) Field() error { return nil }

func (p *Parser) FragmentSpread() error { return nil }

func (p *Parser) InlineFragment() error { return nil }

func (p *Parser) Alias() error { return nil }

func (p *Parser) Arguments() error { return nil }

func (p *Parser) Directives() error { return nil }

func (p *Parser) Directive() error { return nil }

func (p *Parser) Argument() error { return nil }

func (p *Parser) TypeCondition() error { return nil }

func (p *Parser) Value() error { return nil }

func (p *Parser) Variable() error { return nil }

func (p *Parser) BooleanValue() error { return nil }

func (p *Parser) NullValue() error { return nil }

func (p *Parser) EnumValue() error { return nil }

func (p *Parser) ListValue() error { return nil }

func (p *Parser) ObjectValue() error { return nil }

func (p *Parser) ObjectField() error { return nil }

func (p *Parser) VariableDefinitions() error { return nil }

func (p *Parser) DefaultValue() error { return nil }

func (p *Parser) Type() error { return nil }

func (p *Parser) NamedType() error { return nil }

func (p *Parser) ListType() error { return nil }

func (p *Parser) NonNullType() error { return nil }

func (p *Parser) match(k Kind) error { return nil }
