package ast

import (
	"errors"
	"fmt"
	"go/token"
)

// converts GraphQL source into AST.
type parser struct {
	input        *lexer
	lookAheads   []Token // LL(2), look two tokens ahead
	tokenOffsets []int   // corresponding offset of two tokens
	curr         int
}

func newParser(source []byte, filename string, fset *token.FileSet) *parser {
	if source == nil {
		return nil
	}

	f := fset.AddFile(filename, -1, len(source))
	l := newLexer(source, f)

	p := &parser{
		input:        l,
		lookAheads:   make([]Token, 2),
		tokenOffsets: make([]int, 2),
	}

	for i := 0; i < 2; i++ {
		p.consume()
	}

	return p
}

func (p *parser) parseDocument() (*Document, error) {
	if p == nil {
		return nil, errors.New("parser nil")
	}

	document := &Document{}

	defn, err := p.definition()
	if err != nil {
		return document, err
	}
	document.Defs = append(document.Defs, defn)

	for p.lookAhead(1) != TokenEOF {
		defn, err = p.definition()
		if err != nil {
			return document, err
		}
		document.Defs = append(document.Defs, defn)
	}
	return document, nil
}

func (p *parser) definition() (Definition, error) {
	if p.lookAhead(1).Text == Stringify(FRAGMENT) {
		return p.fragmentDefinition()
	}
	return p.operationDefinition()
}

func (p *parser) parseSchema() (*Schema, error) {
	if p == nil {
		return nil, errors.New("parser nil")
	}

	schema, err := p.schema()
	if err != nil {
		return nil, err
	}

	return schema, nil
}

func (p *parser) schema() (*Schema, error) {
	s := &Schema{}
	var node, first, last Node
	var err error
	isFirst := true
	for p.lookAhead(1) != TokenEOF {
		switch p.lookAhead(1).Text {
		case Stringify(INTERFACE):
			node, err = p.interfaceDefinition()
			if err != nil {
				return nil, err
			}
			s.Interfaces = append(s.Interfaces, node.(*InterfaceDefinition))
		case Stringify(SCALAR):
			node, err = p.scalarDefinition()
			if err != nil {
				return nil, err
			}
			s.Scalars = append(s.Scalars, node.(*ScalarDefinition))
		case Stringify(INPUT):
			node, err = p.inputObjectDefinition()
			if err != nil {
				return nil, err
			}
			s.InputObjects = append(s.InputObjects, node.(*InputObjectDefinition))
		case Stringify(TYPE):
			node, err = p.typeDefinition()
			if err != nil {
				return nil, err
			}
			s.Types = append(s.Types, node.(*TypeDefinition))
		case Stringify(EXTEND):
			node, err = p.extendDefinition()
			if err != nil {
				return nil, err
			}
			s.Extends = append(s.Extends, node.(*ExtendDefinition))
		case Stringify(DIRECTIVE):
			node, err = p.directiveDefinition()
			if err != nil {
				return nil, err
			}
			s.Directives = append(s.Directives, node.(*DirectiveDefinition))
		case Stringify(SCHEMA):
			node, err = p.schemaDefinition()
			if err != nil {
				return nil, err
			}
			s.Schemas = append(s.Schemas, node.(*SchemaDefinition))
		case Stringify(ENUM):
			node, err = p.enumDefinition()
			if err != nil {
				return nil, err
			}
			s.Enums = append(s.Enums, node.(*EnumDefinition))
		default:
			node, err = p.unionDefinition()
			if err != nil {
				return nil, err
			}
			s.Unions = append(s.Unions, node.(*UnionDefinition))
		}

		// keep recording the last node seen
		last = node

		if isFirst {
			isFirst = false
			first = node
		}
	}

	s.pos, s.end = first.Pos(), last.End()
	return s, nil
}

func (p *parser) operationDefinition() (*OperationDefinition, error) {
	var err error

	operDefn := &OperationDefinition{}

	if p.lookAhead(1).Kind == LBRACE {
		operDefn.SelSet, err = p.selectionSet()
		return operDefn, err
	}

	operDefn.OperType = p.lookAhead(1)
	operDefn.OperPos = p.input.pos(p.tokenOffset(1))
	if err = p.match(QUERY); err != nil {
		if err = p.match(MUTATION); err != nil {
			if err = p.match(SUBSCRIPTION); err != nil {
				expect := fmt.Sprintf("%s or %s or %s",
					Stringify(QUERY),
					Stringify(MUTATION),
					Stringify(SUBSCRIPTION))
				return nil, p.parseError(expect)
			}
		}
	}

	if p.lookAhead(1).Kind == NAME {
		operDefn.Name = p.lookAhead(1)
		operDefn.NamePos = p.input.pos(p.tokenOffset(1))
		p.match(NAME)
	}

	if p.lookAhead(1).Kind == LPAREN {
		operDefn.VarDefns, err = p.variableDefinitions()
		if err != nil {
			return nil, err
		}
	}

	if p.lookAhead(1).Kind == AT {
		operDefn.Directs, err = p.directives()
		if err != nil {
			return nil, err
		}
	}

	operDefn.SelSet, err = p.selectionSet()
	return operDefn, err
}

func (p *parser) variableDefinitions() (*VariableDefinitions, error) {
	varDefns := &VariableDefinitions{
		Lparen: p.input.pos(p.tokenOffset(1)),
	}
	err := p.match(LPAREN)
	if err != nil {
		return varDefns, err
	}

	var varDefn *VariableDefinition
	varDefn, err = p.variableDefinition()
	if err != nil {
		return varDefns, err
	}
	varDefns.VarDefns = append(varDefns.VarDefns, varDefn)

	for p.lookAhead(1).Kind != RPAREN {
		varDefn, err = p.variableDefinition()
		if err != nil {
			return varDefns, err
		}
		varDefns.VarDefns = append(varDefns.VarDefns, varDefn)
	}

	varDefns.Rparen = p.input.pos(p.tokenOffset(1))
	if err = p.match(RPAREN); err != nil {
		return varDefns, err
	}

	return varDefns, nil
}

func (p *parser) variableDefinition() (*VariableDefinition, error) {
	var err error

	varDefn := &VariableDefinition{}
	varDefn.Var, err = p.variable()
	if err != nil {
		return varDefn, err
	}

	varDefn.Colon = p.input.pos(p.tokenOffset(1))
	if err = p.match(COLON); err != nil {
		return varDefn, err
	}

	varDefn.Typ, err = p.types()
	if err != nil {
		return varDefn, err
	}

	if p.lookAhead(1).Kind == EQL {
		varDefn.DeflVal, err = p.defaultValue()
		if err != nil {
			return varDefn, err
		}
	}

	return varDefn, nil
}

func (p *parser) types() (Type, error) {
	var typ Type
	var err error
	var isNamedType bool
	if p.lookAhead(1).Kind == NAME {
		isNamedType = true
		typ, err = p.namedType()
		if err != nil {
			return nil, err
		}
	} else {
		typ, err = p.listType()
		if err != nil {
			return nil, err
		}
	}

	// non-null
	if p.lookAhead(1).Kind == BANG {
		bangPos := p.input.pos(p.tokenOffset(1))
		p.match(BANG)

		if isNamedType {
			typ.(*NamedType).NonNull = true
			typ.(*NamedType).BangPos = bangPos
		} else {
			typ.(*ListType).NonNull = true
			typ.(*ListType).BangPos = bangPos
		}
	}

	return typ, nil
}

func (p *parser) namedType() (*NamedType, error) {
	namedTyp := &NamedType{
		Name:    p.lookAhead(1),
		NamePos: p.input.pos(p.tokenOffset(1)),
	}
	return namedTyp, p.match(NAME)
}

func (p *parser) listType() (*ListType, error) {
	listTyp := &ListType{
		Lbrack: p.input.pos(p.tokenOffset(1)),
	}
	err := p.match(LBRACK)
	if err != nil {
		return nil, err
	}

	var typ Type
	typ, err = p.types()
	if err != nil {
		return nil, err
	}
	listTyp.Typ = typ

	listTyp.Rbrack = p.input.pos(p.tokenOffset(1))
	err = p.match(RBRACK)
	if err != nil {
		return nil, err
	}

	return listTyp, nil
}

func (p *parser) defaultValue() (*DefaultValue, error) {
	deflVal := &DefaultValue{
		Eq: p.input.pos(p.tokenOffset(1)),
	}

	err := p.match(EQL)
	if err != nil {
		return deflVal, err
	}

	deflVal.Val, err = p.valueConst()
	if err != nil {
		return deflVal, err
	}

	return deflVal, nil
}

func (p *parser) valueConst() (Value, error) {
	switch p.lookAhead(1).Kind {
	case INT, FLOAT, STRING:
		val := &LiteralValue{
			Val:    p.lookAhead(1),
			ValPos: p.input.pos(p.tokenOffset(1)),
		}
		return val, p.match(p.lookAhead(1).Kind)
	case NAME:
		return p.nameValue()
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
		return nil, p.parseError(expect)
	}
}

func (p *parser) listValueConst() (*ListValue, error) {
	listVal := &ListValue{}

	listVal.Lbrack = p.input.pos(p.tokenOffset(1))
	err := p.match(LBRACK)
	if err != nil {
		return nil, err
	}

	if p.lookAhead(1).Kind == RBRACK {
		listVal.Rbrack = p.input.pos(p.tokenOffset(1))
		return listVal, p.match(RBRACK)
	}

	var val Value
	val, err = p.valueConst()
	if err != nil {
		return nil, err
	}
	listVal.Vals = append(listVal.Vals, val)

	for p.lookAhead(1).Kind != RBRACK {
		val, err = p.valueConst()
		if err != nil {
			return nil, err
		}
		listVal.Vals = append(listVal.Vals, val)
	}

	listVal.Rbrack = p.input.pos(p.tokenOffset(1))
	return listVal, p.match(RBRACK)
}

func (p *parser) objectValueConst() (*ObjectValue, error) {
	objVal := &ObjectValue{}

	objVal.Lbrace = p.input.pos(p.tokenOffset(1))
	err := p.match(LBRACE)
	if err != nil {
		return nil, err
	}

	if p.lookAhead(1).Kind == RBRACE {
		objVal.Rbrace = p.input.pos(p.tokenOffset(1))
		return objVal, p.match(RBRACE)
	}

	var objField *ObjectField
	objField, err = p.objectFieldConst()
	if err != nil {
		return nil, err
	}
	objVal.ObjFields = append(objVal.ObjFields, objField)

	for p.lookAhead(1).Kind != RBRACE {
		objField, err = p.objectFieldConst()
		if err != nil {
			return nil, err
		}
		objVal.ObjFields = append(objVal.ObjFields, objField)
	}

	objVal.Rbrace = p.input.pos(p.tokenOffset(1))
	return objVal, p.match(RBRACE)
}

func (p *parser) objectFieldConst() (*ObjectField, error) {
	objField := &ObjectField{}

	objField.Name = p.lookAhead(1)
	objField.NamePos = p.input.pos(p.tokenOffset(1))
	err := p.match(NAME)
	if err != nil {
		return nil, err
	}

	objField.Colon = p.input.pos(p.tokenOffset(1))
	err = p.match(COLON)
	if err != nil {
		return nil, err
	}

	objField.Val, err = p.valueConst()
	if err != nil {
		return nil, err
	}

	return objField, nil
}

func (p *parser) selectionSet() (*SelectionSet, error) {
	selSet := &SelectionSet{}

	selSet.Lbrace = p.input.pos(p.tokenOffset(1))
	err := p.match(LBRACE)
	if err != nil {
		return nil, err
	}

	var sel Selection
	sel, err = p.selection()
	if err != nil {
		return nil, err
	}
	selSet.Sels = append(selSet.Sels, sel)

	for p.lookAhead(1).Kind != RBRACE {
		sel, err = p.selection()
		if err != nil {
			return nil, err
		}
		selSet.Sels = append(selSet.Sels, sel)
	}

	selSet.Rbrace = p.input.pos(p.tokenOffset(1))
	if err = p.match(RBRACE); err != nil {
		return nil, err
	}
	return selSet, nil
}

func (p *parser) selection() (Selection, error) {
	if p.lookAhead(1).Kind == SPREAD {
		if p.lookAhead(2).Kind == NAME && p.lookAhead(2).Text != Stringify(ON) {
			return p.fragmentSpread()
		}
		return p.inlineFragment()
	}

	return p.field()
}

func (p *parser) field() (*Field, error) {
	field := &Field{}

	var err error
	if p.lookAhead(1).Kind == NAME && p.lookAhead(2).Kind == COLON {
		field.Als, err = p.alias()
		if err != nil {
			return nil, err
		}
	}

	field.Name = p.lookAhead(1)
	field.NamePos = p.input.pos(p.tokenOffset(1))
	if err = p.match(NAME); err != nil {
		return nil, err
	}

	if p.lookAhead(1).Kind == LPAREN {
		field.Args, err = p.arguments()
		if err != nil {
			return nil, err
		}
	}

	if p.lookAhead(1).Kind == AT {
		field.Directs, err = p.directives()
		if err != nil {
			return nil, err
		}
	}

	if p.lookAhead(1).Kind == LBRACE {
		field.SelSet, err = p.selectionSet()
		if err != nil {
			return nil, err
		}
	}

	return field, nil
}

func (p *parser) fragmentSpread() (*FragmentSpread, error) {
	frag := &FragmentSpread{}

	frag.Spread = p.input.pos(p.tokenOffset(1))
	err := p.match(SPREAD)
	if err != nil {
		return nil, err
	}

	if p.lookAhead(1).Text == Stringify(ON) {
		return nil, p.parseError("NAME but not *on*")
	}

	frag.Name = p.lookAhead(1)
	frag.NamePos = p.input.pos(p.tokenOffset(1))
	if err = p.match(NAME); err != nil {
		return nil, err
	}

	if p.lookAhead(1).Kind == AT {
		frag.Directs, err = p.directives()
		if err != nil {
			return nil, err
		}
	}

	return frag, nil
}

func (p *parser) inlineFragment() (*InlineFragment, error) {
	frag := &InlineFragment{}

	frag.Spread = p.input.pos(p.tokenOffset(1))
	err := p.match(SPREAD)
	if err != nil {
		return nil, err
	}

	if p.lookAhead(1).Text == Stringify(ON) {
		frag.TypeCond, err = p.typeCondition()
		if err != nil {
			return nil, err
		}
	}

	if p.lookAhead(1).Kind == AT {
		frag.Directs, err = p.directives()
		if err != nil {
			return nil, err
		}
	}

	frag.SelSet, err = p.selectionSet()
	if err != nil {
		return nil, err
	}

	return frag, nil
}

func (p *parser) typeCondition() (*TypeCondition, error) {
	typCond := &TypeCondition{}

	typCond.On = p.input.pos(p.tokenOffset(1))
	err := p.match(ON)
	if err != nil {
		return nil, err
	}

	typCond.NamedTyp, err = p.namedType()
	if err != nil {
		return nil, err
	}

	return typCond, nil
}

func (p *parser) alias() (*Alias, error) {
	a := &Alias{}
	a.Name = p.lookAhead(1)
	a.NamePos = p.input.pos(p.tokenOffset(1))
	err := p.match(NAME)
	if err != nil {
		return nil, err
	}

	a.Colon = p.input.pos(p.tokenOffset(1))
	if err = p.match(COLON); err != nil {
		return nil, err
	}

	return a, nil
}

func (p *parser) arguments() (*Arguments, error) {
	args := &Arguments{}

	args.Lparen = p.input.pos(p.tokenOffset(1))
	err := p.match(LPAREN)
	if err != nil {
		return nil, err
	}

	var arg *Argument
	arg, err = p.argument()
	if err != nil {
		return nil, err
	}
	args.Args = append(args.Args, arg)

	for p.lookAhead(1).Kind != RPAREN {
		arg, err = p.argument()
		if err != nil {
			return nil, err
		}
		args.Args = append(args.Args, arg)
	}

	args.Rparen = p.input.pos(p.tokenOffset(1))
	if err = p.match(RPAREN); err != nil {
		return nil, err
	}
	return args, nil
}

func (p *parser) argument() (*Argument, error) {
	arg := &Argument{}

	arg.Name = p.lookAhead(1)
	arg.NamePos = p.input.pos(p.tokenOffset(1))
	err := p.match(NAME)
	if err != nil {
		return nil, err
	}

	arg.Colon = p.input.pos(p.tokenOffset(1))
	if err = p.match(COLON); err != nil {
		return nil, err
	}

	arg.Val, err = p.value()
	if err != nil {
		return nil, err
	}

	return arg, nil
}

func (p *parser) value() (Value, error) {
	switch p.lookAhead(1).Kind {
	case INT, FLOAT, STRING:
		val := &LiteralValue{
			Val:    p.lookAhead(1),
			ValPos: p.input.pos(p.tokenOffset(1)),
		}
		return val, p.match(p.lookAhead(1).Kind)
	case NAME:
		return p.nameValue()
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
		return nil, p.parseError(expect)
	}
}

func (p *parser) variable() (*Variable, error) {
	variable := &Variable{
		Dollar: p.input.pos(p.tokenOffset(1)),
	}

	err := p.match(DOLLAR)
	if err != nil {
		return variable, err
	}

	variable.Name = p.lookAhead(1)
	variable.NamePos = p.input.pos(p.tokenOffset(1))
	if err = p.match(NAME); err != nil {
		return variable, err
	}

	return variable, nil
}

func (p *parser) nameValue() (*NameValue, error) {
	val := &NameValue{
		Val:    p.lookAhead(1),
		ValPos: p.input.pos(p.tokenOffset(1)),
	}
	return val, p.match(NAME)
}

func (p *parser) enumValue() (*EnumValue, error) {
	val := &EnumValue{
		Name:    p.lookAhead(1),
		NamePos: p.input.pos(p.tokenOffset(1)),
	}

	err := p.match(NAME)
	if err != nil {
		return nil, err
	}

	if p.lookAhead(1).Kind == AT {
		directs, err := p.directives()
		if err != nil {
			return nil, err
		}
		val.Directs = directs
	}
	return val, nil
}

func (p *parser) listValue() (*ListValue, error) {
	listVal := &ListValue{}

	listVal.Lbrack = p.input.pos(p.tokenOffset(1))
	err := p.match(LBRACK)
	if err != nil {
		return nil, err
	}

	if p.lookAhead(1).Kind == RBRACK {
		listVal.Rbrack = p.input.pos(p.tokenOffset(1))
		return listVal, p.match(RBRACK)
	}

	var val Value
	val, err = p.value()
	if err != nil {
		return nil, err
	}
	listVal.Vals = append(listVal.Vals, val)

	for p.lookAhead(1).Kind != RBRACK {
		val, err = p.value()
		if err != nil {
			return nil, err
		}
		listVal.Vals = append(listVal.Vals, val)
	}

	listVal.Rbrack = p.input.pos(p.tokenOffset(1))
	return listVal, p.match(RBRACK)
}

func (p *parser) objectValue() (*ObjectValue, error) {
	objVal := &ObjectValue{}

	objVal.Lbrace = p.input.pos(p.tokenOffset(1))
	err := p.match(LBRACE)
	if err != nil {
		return nil, err
	}

	if p.lookAhead(1).Kind == RBRACE {
		objVal.Rbrace = p.input.pos(p.tokenOffset(1))
		return objVal, p.match(RBRACE)
	}

	var objField *ObjectField
	objField, err = p.objectField()
	if err != nil {
		return nil, err
	}
	objVal.ObjFields = append(objVal.ObjFields, objField)

	for p.lookAhead(1).Kind != RBRACE {
		objField, err = p.objectField()
		if err != nil {
			return nil, err
		}
		objVal.ObjFields = append(objVal.ObjFields, objField)
	}

	objVal.Rbrace = p.input.pos(p.tokenOffset(1))
	return objVal, p.match(RBRACE)
}

func (p *parser) objectField() (*ObjectField, error) {
	objField := &ObjectField{}

	objField.Name = p.lookAhead(1)
	objField.NamePos = p.input.pos(p.tokenOffset(1))
	err := p.match(NAME)
	if err != nil {
		return nil, err
	}

	objField.Colon = p.input.pos(p.tokenOffset(1))
	if err = p.match(COLON); err != nil {
		return nil, err
	}

	objField.Val, err = p.value()
	if err != nil {
		return nil, err
	}

	return objField, nil
}

func (p *parser) directives() (*Directives, error) {
	directs := &Directives{}

	direct, err := p.directive()
	if err != nil {
		return nil, err
	}
	directs.Directs = append(directs.Directs, direct)

	for p.lookAhead(1).Kind == AT {
		direct, err = p.directive()
		if err != nil {
			return nil, err
		}
		directs.Directs = append(directs.Directs, direct)
	}

	return directs, nil
}

func (p *parser) directive() (*Directive, error) {
	direct := &Directive{}

	direct.At = p.input.pos(p.tokenOffset(1))
	err := p.match(AT)
	if err != nil {
		return nil, err
	}

	direct.Name = p.lookAhead(1)
	direct.NamePos = p.input.pos(p.tokenOffset(1))
	if err = p.match(NAME); err != nil {
		return nil, err
	}

	if p.lookAhead(1).Kind == LPAREN {
		direct.Args, err = p.arguments()
		if err != nil {
			return nil, err
		}
	}

	return direct, nil
}

func (p *parser) fragmentDefinition() (*FragmentDefinition, error) {
	frag := &FragmentDefinition{}

	frag.Fragment = p.input.pos(p.tokenOffset(1))
	err := p.match(FRAGMENT)
	if err != nil {
		return nil, err
	}

	if p.lookAhead(1).Text == Stringify(ON) {
		return nil, p.parseError("NAME but not *on*")
	}

	frag.Name = p.lookAhead(1)
	frag.NamePos = p.input.pos(p.tokenOffset(1))
	if err = p.match(NAME); err != nil {
		return nil, err
	}

	frag.TypeCond, err = p.typeCondition()
	if err != nil {
		return nil, err
	}

	if p.lookAhead(1).Kind == AT {
		frag.Directs, err = p.directives()
		if err != nil {
			return nil, err
		}
	}

	frag.SelSet, err = p.selectionSet()
	if err != nil {
		return nil, err
	}

	return frag, nil
}

func (p *parser) lookAhead(i int) Token {
	return p.lookAheads[(p.curr+i-1)%len(p.lookAheads)]
}

func (p *parser) tokenOffset(i int) int {
	return p.tokenOffsets[(p.curr+i-1)%len(p.tokenOffsets)]
}

func (p *parser) match(k Kind) error {
	// fmt.Println("DEBUG tok", p.lookAhead(1))
	if IsReserved(k) {
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

	return p.parseError(Stringify(k))
}

func (p *parser) consume() {
	tok, offs := p.input.read()

	// record the token and position of its first character
	p.lookAheads[p.curr] = tok
	p.tokenOffsets[p.curr] = offs - 1 // minus one to start from zero

	p.curr = (p.curr + 1) % len(p.lookAheads)
}

func (p *parser) parseError(expect string) error {
	return ErrBadParse{
		pos:    p.input.positionFor(p.tokenOffset(1)),
		expect: expect,
		found:  p.lookAhead(1).Text,
	}
}

func (p *parser) interfaceDefinition() (*InterfaceDefinition, error) {
	inter := &InterfaceDefinition{}

	inter.Interface = p.input.pos(p.tokenOffset(1))
	err := p.match(INTERFACE)
	if err != nil {
		return nil, err
	}

	inter.Name = p.lookAhead(1)
	inter.NamePos = p.input.pos(p.tokenOffset(1))
	if err = p.match(NAME); err != nil {
		return nil, err
	}

	if p.lookAhead(1).Kind == AT {
		inter.Directs, err = p.directives()
		if err != nil {
			return nil, err
		}
	}

	inter.Lbrace = p.input.pos(p.tokenOffset(1))
	if err = p.match(LBRACE); err != nil {
		return nil, err
	}

	var fieldDefn *FieldDefinition
	fieldDefn, err = p.fieldDefinition()
	if err != nil {
		return nil, err
	}
	inter.FieldDefns = append(inter.FieldDefns, fieldDefn)

	for p.lookAhead(1).Kind != RBRACE {
		fieldDefn, err = p.fieldDefinition()
		if err != nil {
			return nil, err
		}
		inter.FieldDefns = append(inter.FieldDefns, fieldDefn)
	}

	inter.Rbrace = p.input.pos(p.tokenOffset(1))
	return inter, p.match(RBRACE)
}

func (p *parser) fieldDefinition() (*FieldDefinition, error) {
	fieldDefn := &FieldDefinition{}

	fieldDefn.Name = p.lookAhead(1)
	fieldDefn.NamePos = p.input.pos(p.tokenOffset(1))
	err := p.match(NAME)
	if err != nil {
		return nil, err
	}

	if p.lookAhead(1).Kind == LPAREN {
		fieldDefn.ArgDefns, err = p.argumentsDefinition()
		if err != nil {
			return nil, err
		}
	}

	fieldDefn.Colon = p.input.pos(p.tokenOffset(1))
	if err = p.match(COLON); err != nil {
		return nil, err
	}

	fieldDefn.Typ, err = p.types()
	if err != nil {
		return nil, err
	}

	if p.lookAhead(1).Kind == AT {
		fieldDefn.Directs, err = p.directives()
		if err != nil {
			return nil, err
		}
	}

	return fieldDefn, nil
}

func (p *parser) argumentsDefinition() (*ArgumentsDefinition, error) {
	argsDefn := &ArgumentsDefinition{}

	argsDefn.Lparen = p.input.pos(p.tokenOffset(1))
	err := p.match(LPAREN)
	if err != nil {
		return nil, err
	}

	var input *InputValueDefinition
	input, err = p.inputValueDefinition()
	if err != nil {
		return nil, err
	}
	argsDefn.InputValDefns = append(argsDefn.InputValDefns, input)

	for p.lookAhead(1).Kind != RPAREN {
		input, err = p.inputValueDefinition()
		if err != nil {
			return nil, err
		}
		argsDefn.InputValDefns = append(argsDefn.InputValDefns, input)
	}

	argsDefn.Rparen = p.input.pos(p.tokenOffset(1))
	return argsDefn, p.match(RPAREN)
}

func (p *parser) inputValueDefinition() (*InputValueDefinition, error) {
	input := &InputValueDefinition{}

	input.Name = p.lookAhead(1)
	input.NamePos = p.input.pos(p.tokenOffset(1))
	err := p.match(NAME)
	if err != nil {
		return nil, err
	}

	input.Colon = p.input.pos(p.tokenOffset(1))
	if err = p.match(COLON); err != nil {
		return nil, err
	}

	input.Typ, err = p.types()
	if err != nil {
		return nil, err
	}

	if p.lookAhead(1).Kind == EQL {
		input.DeflVal, err = p.defaultValue()
		if err != nil {
			return nil, err
		}
	}

	if p.lookAhead(1).Kind == AT {
		input.Directs, err = p.directives()
		if err != nil {
			return nil, err
		}
	}

	return input, nil
}

func (p *parser) scalarDefinition() (*ScalarDefinition, error) {
	scalar := &ScalarDefinition{}

	scalar.Scalar = p.input.pos(p.tokenOffset(1))
	err := p.match(SCALAR)
	if err != nil {
		return nil, err
	}

	scalar.Name = p.lookAhead(1)
	scalar.NamePos = p.input.pos(p.tokenOffset(1))
	if err = p.match(NAME); err != nil {
		return nil, err
	}

	if p.lookAhead(1).Kind == AT {
		scalar.Directs, err = p.directives()
		if err != nil {
			return nil, err
		}
	}
	return scalar, nil
}

func (p *parser) inputObjectDefinition() (*InputObjectDefinition, error) {
	input := &InputObjectDefinition{}

	input.Input = p.input.pos(p.tokenOffset(1))
	err := p.match(INPUT)
	if err != nil {
		return nil, err
	}

	input.Name = p.lookAhead(1)
	input.NamePos = p.input.pos(p.tokenOffset(1))
	if err = p.match(NAME); err != nil {
		return nil, err
	}

	if p.lookAhead(1).Kind == AT {
		input.Directs, err = p.directives()
		if err != nil {
			return nil, err
		}
	}

	input.Lbrace = p.input.pos(p.tokenOffset(1))
	if err = p.match(LBRACE); err != nil {
		return nil, err
	}

	var vDefn *InputValueDefinition
	vDefn, err = p.inputValueDefinition()
	if err != nil {
		return nil, err
	}
	input.InputValDefns = append(input.InputValDefns, vDefn)

	for p.lookAhead(1).Kind != RBRACE {
		vDefn, err = p.inputValueDefinition()
		if err != nil {
			return nil, err
		}
		input.InputValDefns = append(input.InputValDefns, vDefn)
	}

	input.Rbrace = p.input.pos(p.tokenOffset(1))
	return input, p.match(RBRACE)
}

func (p *parser) typeDefinition() (*TypeDefinition, error) {
	typDefn := &TypeDefinition{}

	typDefn.Typ = p.input.pos(p.tokenOffset(1))
	err := p.match(TYPE)
	if err != nil {
		return nil, err
	}

	typDefn.Name = p.lookAhead(1)
	typDefn.NamePos = p.input.pos(p.tokenOffset(1))
	if err = p.match(NAME); err != nil {
		return nil, err
	}

	if p.lookAhead(1).Text == Stringify(IMPLEMENTS) {
		typDefn.Implements, err = p.implementsInterfaces()
		if err != nil {
			return nil, err
		}
	}

	if p.lookAhead(1).Kind == AT {
		typDefn.Directs, err = p.directives()
		if err != nil {
			return nil, err
		}
	}

	typDefn.Lbrace = p.input.pos(p.tokenOffset(1))
	if err = p.match(LBRACE); err != nil {
		return nil, err
	}

	// if err = p.fieldDefinition(); err != nil {
	// 	return err
	// }

	var fieldDefn *FieldDefinition
	for p.lookAhead(1).Kind != RBRACE {
		fieldDefn, err = p.fieldDefinition()
		if err != nil {
			return nil, err
		}
		typDefn.FieldDefns = append(typDefn.FieldDefns, fieldDefn)
	}

	typDefn.Rbrace = p.input.pos(p.tokenOffset(1))
	return typDefn, p.match(RBRACE)
}

func (p *parser) implementsInterfaces() (*ImplementsInterfaces, error) {
	implement := &ImplementsInterfaces{}

	implement.Implements = p.input.pos(p.tokenOffset(1))
	err := p.match(IMPLEMENTS)
	if err != nil {
		return nil, err
	}

	var namedTyp *NamedType
	namedTyp, err = p.namedType()
	if err != nil {
		return nil, err
	}
	implement.NamedTyps = append(implement.NamedTyps, namedTyp)

	for p.lookAhead(1).Kind == NAME {
		namedTyp, err = p.namedType()
		if err != nil {
			return nil, err
		}
		implement.NamedTyps = append(implement.NamedTyps, namedTyp)
	}

	return implement, nil
}

func (p *parser) extendDefinition() (*ExtendDefinition, error) {
	e := &ExtendDefinition{}

	e.Extend = p.input.pos(p.tokenOffset(1))
	err := p.match(EXTEND)
	if err != nil {
		return nil, err
	}

	e.TypDefn, err = p.typeDefinition()
	if err != nil {
		return nil, err
	}

	return e, nil
}

func (p *parser) directiveDefinition() (*DirectiveDefinition, error) {
	d := &DirectiveDefinition{}

	d.Direct = p.input.pos(p.tokenOffset(1))
	err := p.match(DIRECTIVE)
	if err != nil {
		return nil, err
	}

	d.At = p.input.pos(p.tokenOffset(1))
	if err = p.match(AT); err != nil {
		return nil, err
	}

	d.Name = p.lookAhead(1)
	d.NamePos = p.input.pos(p.tokenOffset(1))
	if err = p.match(NAME); err != nil {
		return nil, err
	}

	if p.lookAhead(1).Kind == LPAREN {
		d.Args, err = p.argumentsDefinition()
		if err != nil {
			return nil, err
		}
	}

	d.On = p.input.pos(p.tokenOffset(1))
	if err = p.match(ON); err != nil {
		return nil, err
	}

	d.Locs, err = p.directiveLocations()
	if err != nil {
		return nil, err
	}

	return d, nil
}

func (p *parser) directiveLocations() (*DirectiveLocations, error) {
	loc := &DirectiveLocations{}

	loc.Name = p.lookAhead(1)
	loc.NamePos = p.input.pos(p.tokenOffset(1))
	err := p.match(NAME)
	if err != nil {
		return nil, err
	}

	var l *DirectiveLocation
	for p.lookAhead(1).Kind == PIPE {
		l, err = p.directiveLocation()
		if err != nil {
			return nil, err
		}
		loc.Locs = append(loc.Locs, l)
	}

	return loc, nil
}

func (p *parser) directiveLocation() (*DirectiveLocation, error) {
	var err error
	loc := &DirectiveLocation{}
	loc.Pipe = p.input.pos(p.tokenOffset(1))
	if err = p.match(PIPE); err != nil {
		return nil, err
	}

	loc.Name = p.lookAhead(1)
	loc.NamePos = p.input.pos(p.tokenOffset(1))
	if err = p.match(NAME); err != nil {
		return nil, err
	}

	return loc, nil
}

func (p *parser) schemaDefinition() (*SchemaDefinition, error) {
	schemaDefn := &SchemaDefinition{}

	schemaDefn.Schema = p.input.pos(p.tokenOffset(1))
	err := p.match(SCHEMA)
	if err != nil {
		return nil, err
	}

	if p.lookAhead(1).Kind == AT {
		schemaDefn.Directs, err = p.directives()
		if err != nil {
			return nil, err
		}
	}

	schemaDefn.Lbrace = p.input.pos(p.tokenOffset(1))
	if err = p.match(LBRACE); err != nil {
		return nil, err
	}

	var oper *OperationTypeDefinition
	oper, err = p.operationTypeDefinition()
	if err != nil {
		return nil, err
	}
	schemaDefn.OperDefns = append(schemaDefn.OperDefns, oper)

	for p.lookAhead(1).Kind != RBRACE {
		oper, err = p.operationTypeDefinition()
		if err != nil {
			return nil, err
		}
		schemaDefn.OperDefns = append(schemaDefn.OperDefns, oper)
	}

	schemaDefn.Rbrace = p.input.pos(p.tokenOffset(1))
	return schemaDefn, p.match(RBRACE)
}

func (p *parser) operationTypeDefinition() (*OperationTypeDefinition, error) {
	var err error
	oper := &OperationTypeDefinition{}

	oper.OperType = p.lookAhead(1)
	oper.OperPos = p.input.pos(p.tokenOffset(1))
	if err = p.match(QUERY); err != nil {
		if err = p.match(MUTATION); err != nil {
			if err = p.match(SUBSCRIPTION); err != nil {
				expect := fmt.Sprintf("%s or %s or %s",
					Stringify(QUERY),
					Stringify(MUTATION),
					Stringify(SUBSCRIPTION))
				return nil, p.parseError(expect)
			}
		}
	}

	oper.Colon = p.input.pos(p.tokenOffset(1))
	if err = p.match(COLON); err != nil {
		return nil, err
	}

	oper.NamedTyp, err = p.namedType()
	if err != nil {
		return nil, err
	}
	return oper, nil
}

func (p *parser) enumDefinition() (*EnumDefinition, error) {
	e := &EnumDefinition{}

	e.Enum = p.input.pos(p.tokenOffset(1))
	err := p.match(ENUM)
	if err != nil {
		return nil, err
	}

	e.Name = p.lookAhead(1)
	e.NamePos = p.input.pos(p.tokenOffset(1))
	if err = p.match(NAME); err != nil {
		return nil, err
	}

	if p.lookAhead(1).Kind == AT {
		e.Directs, err = p.directives()
		if err != nil {
			return nil, err
		}
	}

	e.Lbrace = p.input.pos(p.tokenOffset(1))
	if err = p.match(LBRACE); err != nil {
		return nil, err
	}

	var ev *EnumValue
	ev, err = p.enumValue()
	if err != nil {
		return nil, err
	}
	e.EnumVals = append(e.EnumVals, ev)

	for p.lookAhead(1).Kind != RBRACE {
		ev, err = p.enumValue()
		if err != nil {
			return nil, err
		}
		e.EnumVals = append(e.EnumVals, ev)
	}

	e.Rbrace = p.input.pos(p.tokenOffset(1))
	return e, p.match(RBRACE)
}

func (p *parser) unionDefinition() (*UnionDefinition, error) {
	defn := &UnionDefinition{}

	defn.Union = p.input.pos(p.tokenOffset(1))
	err := p.match(UNION)
	if err != nil {
		return nil, err
	}

	defn.Name = p.lookAhead(1)
	defn.NamePos = p.input.pos(p.tokenOffset(1))
	if err = p.match(NAME); err != nil {
		return nil, err
	}

	if p.lookAhead(1).Kind == AT {
		defn.Directs, err = p.directives()
		if err != nil {
			return nil, err
		}
	}

	defn.Eq = p.input.pos(p.tokenOffset(1))
	if err = p.match(EQL); err != nil {
		return nil, err
	}

	defn.Members, err = p.unionMembers()
	if err != nil {
		return nil, err
	}
	return defn, nil
}

func (p *parser) unionMembers() (*UnionMembers, error) {
	um := &UnionMembers{}

	namedTyp, err := p.namedType()
	if err != nil {
		return nil, err
	}

	um.NamedTyp = namedTyp

	var u *UnionMember
	for p.lookAhead(1).Kind == PIPE {
		u, err = p.unionMember()
		if err != nil {
			return nil, err
		}
		um.Members = append(um.Members, u)
	}

	return um, nil
}

func (p *parser) unionMember() (*UnionMember, error) {
	var err error
	u := &UnionMember{}

	u.Pipe = p.input.pos(p.tokenOffset(1))
	if err = p.match(PIPE); err != nil {
		return nil, err
	}

	u.NamedTyp, err = p.namedType()
	if err != nil {
		return nil, err
	}

	return u, nil
}
