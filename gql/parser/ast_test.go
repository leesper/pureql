package parser

import (
	"fmt"
	"go/token"
	"reflect"
	"testing"
)

func assertEqual(t *testing.T, expected, found interface{}) {
	if !reflect.DeepEqual(expected, found) {
		t.Errorf("expected %#v, found %#v", expected, found)
	}
}

func assertTrue(t *testing.T, v bool) {
	if !v {
		t.Error("should return true")
	}
}

func assertFalse(t *testing.T, v bool) {
	if v {
		t.Error("should return false")
	}
}

func TestDirectives(t *testing.T) {
	directives := "@include(if: $withFriends) @onInputObjectType @onType"
	fset := token.NewFileSet()
	p := newParser([]byte(directives), "", fset)
	d, err := p.directives()
	if err != nil {
		t.Error("unexpected error")
	}

	if len(d.Directs) != 3 {
		t.Error("should have 3 directives")
	}

	assertEqual(t, "include", d.Directs[0].Name.Text)
	assertEqual(t, "onInputObjectType", d.Directs[1].Name.Text)
	assertEqual(t, "onType", d.Directs[2].Name.Text)

	pos := fset.Position(d.Pos())
	expected := "1:1"
	assertEqual(t, expected, pos.String())

	end := fset.Position(d.End())
	expected = fmt.Sprintf("1:%d", len(directives)+1)
	assertEqual(t, expected, end.String())
}

func TestTypes(t *testing.T) {
	types := `Type1
Type2!
[ Type3 ]
[ Type4 ] !
[ [ Type5! ] ]
`
	fset := token.NewFileSet()
	p := newParser([]byte(types), "", fset)
	var typ Type
	var err error

	typ, err = p.types()
	if err != nil {
		t.Error("unexpected error", err)
	}
	nt, ok := typ.(*NamedType)
	if !ok {
		t.Error("should be *NamedType")
	}
	assertEqual(t, "1:1", fset.Position(nt.Pos()).String())
	assertEqual(t, "1:6", fset.Position(nt.End()).String())
	assertFalse(t, nt.NonNull)

	typ, err = p.types()
	if err != nil {
		t.Error("unexpected error", err)
	}
	nt, ok = typ.(*NamedType)
	assertTrue(t, ok)
	assertEqual(t, "2:1", fset.Position(nt.Pos()).String())
	assertEqual(t, "2:7", fset.Position(nt.End()).String())
	assertTrue(t, nt.NonNull)

	typ, err = p.types()
	if err != nil {
		t.Error("unexpected error", err)
	}
	var lt *ListType
	lt, ok = typ.(*ListType)
	assertTrue(t, ok)
	assertEqual(t, "3:1", fset.Position(lt.Pos()).String())
	assertEqual(t, "3:10", fset.Position(lt.End()).String())
	assertFalse(t, lt.NonNull)

	typ, err = p.types()
	if err != nil {
		t.Error("unexpected error", err)
	}
	lt, ok = typ.(*ListType)
	assertTrue(t, ok)
	assertEqual(t, "4:1", fset.Position(lt.Pos()).String())
	assertEqual(t, "4:12", fset.Position(lt.End()).String())
	assertTrue(t, lt.NonNull)

	typ, err = p.types()
	if err != nil {
		t.Error("unexpected error", err)
	}
	lt, ok = typ.(*ListType)
	assertTrue(t, ok)
	assertEqual(t, "5:1", fset.Position(lt.Pos()).String())
	assertEqual(t, "5:15", fset.Position(lt.End()).String())
	assertFalse(t, lt.NonNull)
	var inner *ListType
	inner, ok = lt.Typ.(*ListType)
	assertTrue(t, ok)
	assertEqual(t, "5:3", fset.Position(inner.Pos()).String())
	assertEqual(t, "5:13", fset.Position(inner.End()).String())
	assertFalse(t, inner.NonNull)
	nt, ok = inner.Typ.(*NamedType)
	assertTrue(t, ok)
	assertEqual(t, "5:5", fset.Position(nt.Pos()).String())
	assertEqual(t, "5:11", fset.Position(nt.End()).String())
	assertTrue(t, nt.NonNull)
	assertEqual(t, "Type5", nt.Name.Text)

	_, err = p.types()
	if err == nil {
		t.Error("should return err")
	}
}

func TestVariableDefinitions(t *testing.T) {
	varDefns := `($episode: Episode = "JEDI", $withFriends: Boolean!, $ep: Episode! $review: ReviewInput!)`
	fset := token.NewFileSet()
	p := newParser([]byte(varDefns), "", fset)
	defns, err := p.variableDefinitions()
	if err != nil {
		t.Error("unexpected error", err)
	}

	assertEqual(t, "1:1", fset.Position(defns.Pos()).String())
	assertEqual(t, fmt.Sprintf("1:%d", len(varDefns)+1), fset.Position(defns.End()).String())

	if len(defns.VarDefns) != 4 {
		t.Error("should have 4 variables")
	}

	episode := defns.VarDefns[0]
	assertEqual(t, "1:2", fset.Position(episode.Pos()).String())
	assertEqual(t, "1:28", fset.Position(episode.End()).String())
	assertEqual(t, "1:2", fset.Position(episode.Var.Dollar).String())
	assertEqual(t, "episode", episode.Var.Name.Text)
	assertEqual(t, "1:2", fset.Position(episode.Var.Pos()).String())
	assertEqual(t, "1:10", fset.Position(episode.Var.End()).String())
	assertEqual(t, "1:10", fset.Position(episode.Colon).String())
	assertEqual(t, "Episode", episode.Typ.(*NamedType).Name.Text)
	assertFalse(t, episode.Typ.(*NamedType).NonNull)
	assertEqual(t, "1:20", fset.Position(episode.DeflVal.Eq).String())
	assertEqual(t, "JEDI", episode.DeflVal.Val.(*LiteralValue).Val.Text)
	assertEqual(t, "1:20", fset.Position(episode.DeflVal.Pos()).String())
	assertEqual(t, "1:28", fset.Position(episode.DeflVal.End()).String())
	assertEqual(t, "1:22", fset.Position(episode.DeflVal.Val.Pos()).String())
	assertEqual(t, "1:28", fset.Position(episode.DeflVal.Val.End()).String())

	withFriends := defns.VarDefns[1]
	assertEqual(t, "1:30", fset.Position(withFriends.Pos()).String())
	assertEqual(t, "1:52", fset.Position(withFriends.End()).String())
	assertEqual(t, "withFriends", withFriends.Var.Name.Text)
	assertEqual(t, "Boolean", withFriends.Typ.(*NamedType).Name.Text)
	assertTrue(t, withFriends.Typ.(*NamedType).NonNull)
	assertEqual(t, "1:51", fset.Position(withFriends.Typ.(*NamedType).BangPos).String())
	assertTrue(t, withFriends.DeflVal == nil)

	ep := defns.VarDefns[2]
	assertEqual(t, "1:54", fset.Position(ep.Pos()).String())
	assertEqual(t, "1:67", fset.Position(ep.End()).String())

	review := defns.VarDefns[3]
	assertEqual(t, "1:68", fset.Position(review.Pos()).String())
	assertEqual(t, "1:89", fset.Position(review.End()).String())
	assertEqual(t, "ReviewInput", review.Typ.(*NamedType).Name.Text)
	assertEqual(t, "1:88", fset.Position(review.Typ.(*NamedType).BangPos).String())
	assertTrue(t, review.Typ.(*NamedType).NonNull)
}

func TestValues(t *testing.T) {
	values := `$variable
10225406 3.1415926 "Golang" true false null MALE FEMALE
[] [ 123 true 2.3E-10 $gender]
{} { name: "Hello world", score: 1.0 job: $job}
{search: [
      {
        name: "Han Solo",
        height: 1.8
      },
      {
        name: "Leia Organa",
        height: 1.5
      },
      {
        name: "TIE Advanced x1",
        length: 9.2
      }
		]
}`
	fset := token.NewFileSet()
	fname := "value.graphql"
	p := newParser([]byte(values), fname, fset)
	val, err := p.value()
	if err != nil {
		t.Error("unexpected error", err)
	}
	assertEqual(t, fmt.Sprintf("%s:1:1", fname), fset.Position(val.Pos()).String())
	assertEqual(t, fmt.Sprintf("%s:1:10", fname), fset.Position(val.End()).String())
	assertEqual(t, "variable", val.(*Variable).Name.Text)

	val, err = p.value()
	if err != nil {
		t.Error("unexpected error", err)
	}
	assertEqual(t, fmt.Sprintf("%s:2:1", fname), fset.Position(val.Pos()).String())
	assertEqual(t, fmt.Sprintf("%s:2:9", fname), fset.Position(val.End()).String())
	assertEqual(t, "10225406", val.(*LiteralValue).Val.Text)
	assertEqual(t, INT, val.(*LiteralValue).Val.Kind)

	val, err = p.value()
	if err != nil {
		t.Error("unexpected error", err)
	}
	assertEqual(t, fmt.Sprintf("%s:2:10", fname), fset.Position(val.Pos()).String())
	assertEqual(t, fmt.Sprintf("%s:2:19", fname), fset.Position(val.End()).String())
	assertEqual(t, "3.1415926", val.(*LiteralValue).Val.Text)
	assertEqual(t, FLOAT, val.(*LiteralValue).Val.Kind)

	val, err = p.value()
	if err != nil {
		t.Error("unexpected error", err)
	}
	assertEqual(t, fmt.Sprintf("%s:2:20", fname), fset.Position(val.Pos()).String())
	assertEqual(t, fmt.Sprintf("%s:2:28", fname), fset.Position(val.End()).String())
	assertEqual(t, "Golang", val.(*LiteralValue).Val.Text)
	assertEqual(t, STRING, val.(*LiteralValue).Val.Kind)

	val, err = p.value()
	if err != nil {
		t.Error("unexpected error", err)
	}
	assertEqual(t, fmt.Sprintf("%s:2:29", fname), fset.Position(val.Pos()).String())
	assertEqual(t, fmt.Sprintf("%s:2:33", fname), fset.Position(val.End()).String())
	assertEqual(t, "true", val.(*NameValue).Val.Text)
	assertEqual(t, NAME, val.(*NameValue).Val.Kind)

	val, err = p.value()
	if err != nil {
		t.Error("unexpected error", err)
	}
	assertEqual(t, fmt.Sprintf("%s:2:34", fname), fset.Position(val.Pos()).String())
	assertEqual(t, fmt.Sprintf("%s:2:39", fname), fset.Position(val.End()).String())
	assertEqual(t, "false", val.(*NameValue).Val.Text)
	assertEqual(t, NAME, val.(*NameValue).Val.Kind)

	val, err = p.value()
	if err != nil {
		t.Error("unexpected error", err)
	}
	assertEqual(t, fmt.Sprintf("%s:2:40", fname), fset.Position(val.Pos()).String())
	assertEqual(t, fmt.Sprintf("%s:2:44", fname), fset.Position(val.End()).String())
	assertEqual(t, "null", val.(*NameValue).Val.Text)
	assertEqual(t, NAME, val.(*NameValue).Val.Kind)

	val, err = p.value()
	if err != nil {
		t.Error("unexpected error", err)
	}
	assertEqual(t, fmt.Sprintf("%s:2:45", fname), fset.Position(val.Pos()).String())
	assertEqual(t, fmt.Sprintf("%s:2:49", fname), fset.Position(val.End()).String())
	assertEqual(t, "MALE", val.(*NameValue).Val.Text)
	assertEqual(t, NAME, val.(*NameValue).Val.Kind)

	val, err = p.value()
	if err != nil {
		t.Error("unexpected error", err)
	}
	assertEqual(t, fmt.Sprintf("%s:2:50", fname), fset.Position(val.Pos()).String())
	assertEqual(t, fmt.Sprintf("%s:2:56", fname), fset.Position(val.End()).String())
	assertEqual(t, "FEMALE", val.(*NameValue).Val.Text)
	assertEqual(t, NAME, val.(*NameValue).Val.Kind)

	val, err = p.value()
	if err != nil {
		t.Error("unexpected error", err)
	}
	assertEqual(t, fmt.Sprintf("%s:3:1", fname), fset.Position(val.Pos()).String())
	assertEqual(t, fmt.Sprintf("%s:3:3", fname), fset.Position(val.End()).String())
	assertTrue(t, val.(*ListValue).Vals == nil)

	val, err = p.value()
	if err != nil {
		t.Error("unexpected error", err)
	}
	assertEqual(t, fmt.Sprintf("%s:3:4", fname), fset.Position(val.Pos()).String())
	assertEqual(t, fmt.Sprintf("%s:3:31", fname), fset.Position(val.End()).String())
	assertEqual(t, 4, len(val.(*ListValue).Vals))
	assertEqual(t, "123", val.(*ListValue).Vals[0].(*LiteralValue).Val.Text)
	assertEqual(t, INT, val.(*ListValue).Vals[0].(*LiteralValue).Val.Kind)
	assertEqual(t, "true", val.(*ListValue).Vals[1].(*NameValue).Val.Text)
	assertEqual(t, NAME, val.(*ListValue).Vals[1].(*NameValue).Val.Kind)
	assertEqual(t, "2.3E-10", val.(*ListValue).Vals[2].(*LiteralValue).Val.Text)
	assertEqual(t, FLOAT, val.(*ListValue).Vals[2].(*LiteralValue).Val.Kind)
	assertEqual(t, "gender", val.(*ListValue).Vals[3].(*Variable).Name.Text)

	val, err = p.value()
	if err != nil {
		t.Error("unexpected error", err)
	}
	assertEqual(t, fmt.Sprintf("%s:4:1", fname), fset.Position(val.Pos()).String())
	assertEqual(t, fmt.Sprintf("%s:4:3", fname), fset.Position(val.End()).String())
	assertEqual(t, 0, len(val.(*ObjectValue).ObjFields))

	val, err = p.value()
	if err != nil {
		t.Error("unexpected error", err)
	}
	assertEqual(t, fmt.Sprintf("%s:4:4", fname), fset.Position(val.Pos()).String())
	assertEqual(t, fmt.Sprintf("%s:4:48", fname), fset.Position(val.End()).String())
	assertEqual(t, 3, len(val.(*ObjectValue).ObjFields))
	assertEqual(t, "name", val.(*ObjectValue).ObjFields[0].Name.Text)
	assertEqual(t, "Hello world", val.(*ObjectValue).ObjFields[0].Val.(*LiteralValue).Val.Text)
	assertEqual(t, STRING, val.(*ObjectValue).ObjFields[0].Val.(*LiteralValue).Val.Kind)
	assertEqual(t, "score", val.(*ObjectValue).ObjFields[1].Name.Text)
	assertEqual(t, "1.0", val.(*ObjectValue).ObjFields[1].Val.(*LiteralValue).Val.Text)
	assertEqual(t, FLOAT, val.(*ObjectValue).ObjFields[1].Val.(*LiteralValue).Val.Kind)
	assertEqual(t, "job", val.(*ObjectValue).ObjFields[2].Name.Text)
	assertEqual(t, "job", val.(*ObjectValue).ObjFields[2].Val.(*Variable).Name.Text)
	assertEqual(t, fmt.Sprintf("%s:4:38", fname), fset.Position(val.(*ObjectValue).ObjFields[2].Pos()).String())
	assertEqual(t, fmt.Sprintf("%s:4:47", fname), fset.Position(val.(*ObjectValue).ObjFields[2].End()).String())

	val, err = p.value()
	if err != nil {
		t.Error("unexpected error", err)
	}
	assertEqual(t, fmt.Sprintf("%s:5:1", fname), fset.Position(val.Pos()).String())
	assertEqual(t, fmt.Sprintf("%s:19:2", fname), fset.Position(val.End()).String())
	assertEqual(t, 1, len(val.(*ObjectValue).ObjFields))
	assertEqual(t, 3, len(val.(*ObjectValue).ObjFields[0].Val.(*ListValue).Vals))
	listVal := val.(*ObjectValue).ObjFields[0].Val.(*ListValue)
	assertEqual(t, fmt.Sprintf("%s:5:10", fname), fset.Position(listVal.Pos()).String())
	assertEqual(t, fmt.Sprintf("%s:18:4", fname), fset.Position(listVal.End()).String())
	for _, v := range listVal.Vals {
		assertTrue(t, len(v.(*ObjectValue).ObjFields) == 2)
		assertTrue(t, v.(*ObjectValue).ObjFields[0].Val.(*LiteralValue).Val.Kind == STRING)
		assertTrue(t, v.(*ObjectValue).ObjFields[1].Val.(*LiteralValue).Val.Kind == FLOAT)
	}
}

func TestConstValues(t *testing.T) {
	values := `[ [] ] {  } { a: 1, b: 2 }`
	fset := token.NewFileSet()
	p := newParser([]byte(values), "", fset)

	val, err := p.valueConst()
	if err != nil {
		t.Error("unexpected error", err)
	}
	assertTrue(t, len(val.(*ListValue).Vals) == 1)
	inner := val.(*ListValue).Vals[0].(*ListValue)
	assertTrue(t, len(inner.Vals) == 0)

	val, err = p.valueConst()
	if err != nil {
		t.Error("unexpected error", err)
	}
	assertTrue(t, len(val.(*ObjectValue).ObjFields) == 0)

	val, err = p.valueConst()
	if err != nil {
		t.Error("unexpected error", err)
	}
	assertTrue(t, len(val.(*ObjectValue).ObjFields) == 2)
}
func TestFragmentDefinition(t *testing.T) {
	frag := `
fragment comparisonFields on Character {
	name
	appearsIn
	friends {
		name
	}
}
`
	fset := token.NewFileSet()
	p := newParser([]byte(frag), "", fset)

	def, err := p.fragmentDefinition()
	if err != nil {
		t.Error("unexpected error", err)
	}
	assertEqual(t, "2:1", fset.Position(def.Pos()).String())
	assertEqual(t, "8:2", fset.Position(def.End()).String())
	assertTrue(t, def.Name.Text == "comparisonFields")
	assertTrue(t, def.TypeCond != nil)
	assertTrue(t, def.Directs == nil)
	assertTrue(t, len(def.SelSet.Sels) == 3)
	assertEqual(t, "5:2", fset.Position(def.SelSet.Sels[2].(*Field).Pos()).String())
	assertEqual(t, "7:3", fset.Position(def.SelSet.Sels[2].(*Field).End()).String())
	assertTrue(t, def.SelSet.Sels[2].(*Field).SelSet != nil)
	assertTrue(t, len(def.SelSet.Sels[2].(*Field).SelSet.Sels) == 1)
}

func TestOperationDefinition(t *testing.T) {
	oper := `
query HeroForEpisode($ep: Episode!) {
	hero(episode: $ep) {
		name
		... on Droid {
      primaryFunction
		}
    ... on Human {
      height
    }
  }
}`

	fset := token.NewFileSet()
	p := newParser([]byte(oper), "", fset)
	op, err := p.operationDefinition()
	if err != nil {
		t.Error("unexpected error", err)
	}
	assertEqual(t, "2:1", fset.Position(op.Pos()).String())
	assertEqual(t, "12:2", fset.Position(op.End()).String())
	assertEqual(t, "query", op.OperType.Text)
	assertEqual(t, "HeroForEpisode", op.Name.Text)
	assertTrue(t, op.VarDefns != nil)
	assertTrue(t, op.Directs == nil)
	assertTrue(t, op.SelSet != nil)
	assertTrue(t, len(op.SelSet.Sels) == 1)
	innerSelSet := op.SelSet.Sels[0].(*Field).SelSet
	assertTrue(t, len(innerSelSet.Sels) == 3)
	assertEqual(t, "Droid", innerSelSet.Sels[1].(*InlineFragment).TypeCond.NamedTyp.Name.Text)
	assertEqual(t, "Human", innerSelSet.Sels[2].(*InlineFragment).TypeCond.NamedTyp.Name.Text)
	inlineFrag := innerSelSet.Sels[1].(*InlineFragment)
	assertEqual(t, "5:3", fset.Position(inlineFrag.Pos()).String())
	assertEqual(t, "7:4", fset.Position(inlineFrag.End()).String())
}

func TestDocument(t *testing.T) {
	doc := `
query withFragments {
	user(id: 4) {
		friends(first: 10) {
			...friendFields
		}
		mutualFriends(first: 10) {
			...friendFields
		}
	}
}

fragment friendFields on User {
	id
	name
	profilePic(size: 50)
}`

	fset := token.NewFileSet()
	p := newParser([]byte(doc), "", fset)

	document, err := p.parseDocument()
	if err != nil {
		t.Error("unexpected error", err)
	}

	assertEqual(t, "2:1", fset.Position(document.Pos()).String())
	assertEqual(t, "17:2", fset.Position(document.End()).String())
	assertTrue(t, len(document.Defs) == 2)
	query := document.Defs[0].(*OperationDefinition)
	assertEqual(t, "2:21", fset.Position(query.SelSet.Pos()).String())
	spread := query.SelSet.Sels[0].(*Field).SelSet.Sels[0].(*Field).SelSet.Sels[0].(*FragmentSpread)
	assertEqual(t, "5:4", fset.Position(spread.Pos()).String())
	assertEqual(t, "5:19", fset.Position(spread.End()).String())
	frag := document.Defs[1].(*FragmentDefinition)
	assertEqual(t, "13:1", fset.Position(frag.Pos()).String())
	assertEqual(t, "17:2", fset.Position(frag.End()).String())
}

func TestUnionDefinition(t *testing.T) {
	unions := `
union Feed = Story | Article | Advert
union AnnotatedUnion @onUnion = A | B
union AnnotatedUnionTwo @onUnion = A | B`

	fset := token.NewFileSet()
	p := newParser([]byte(unions), "", fset)
	def, err := p.unionDefinition()
	if err != nil {
		t.Error("unexpected error", err)
	}
	assertEqual(t, "2:1", fset.Position(def.Pos()).String())
	assertEqual(t, "2:38", fset.Position(def.End()).String())
	assertTrue(t, def.Directs == nil)
	assertTrue(t, def.Members != nil)
	assertEqual(t, "Story", def.Members.NamedTyp.Name.Text)
	assertEqual(t, "2:14", fset.Position(def.Members.Pos()).String())
	assertEqual(t, "2:38", fset.Position(def.Members.End()).String())

	def, err = p.unionDefinition()
	if err != nil {
		t.Error("unexpected error", err)
	}
	assertEqual(t, "3:1", fset.Position(def.Pos()).String())
	assertEqual(t, "3:38", fset.Position(def.End()).String())
	assertTrue(t, def.Directs != nil)

	def, err = p.unionDefinition()
	if err != nil {
		t.Error("unexpected error", err)
	}
	assertEqual(t, "4:1", fset.Position(def.Pos()).String())
	assertEqual(t, "4:41", fset.Position(def.End()).String())
	assertTrue(t, def.Directs != nil)
	assertEqual(t, "AnnotatedUnionTwo", def.Name.Text)
}

func TestEnumDefinition(t *testing.T) {
	enums := `
enum Site {
	DESKTOP
	MOBILE
}

enum AnnotatedEnum @onEnum {
	ANNOTATED_VALUE @onEnumValue
	OTHER_VALUE
}`

	fset := token.NewFileSet()
	p := newParser([]byte(enums), "", fset)
	def, err := p.enumDefinition()
	if err != nil {
		t.Error("unexpected error", err)
	}
	assertEqual(t, "2:1", fset.Position(def.Pos()).String())
	assertEqual(t, "5:2", fset.Position(def.End()).String())
	assertTrue(t, def.Directs == nil)
	assertTrue(t, len(def.EnumVals) == 2)
	assertEqual(t, "DESKTOP", def.EnumVals[0].Name.Text)
	assertEqual(t, "MOBILE", def.EnumVals[1].Name.Text)

	def, err = p.enumDefinition()
	if err != nil {
		t.Error("unexpected error", err)
	}
	assertEqual(t, "7:1", fset.Position(def.Pos()).String())
	assertEqual(t, "10:2", fset.Position(def.End()).String())
	assertTrue(t, def.Directs != nil)
	assertTrue(t, len(def.EnumVals) == 2)
	assertEqual(t, "ANNOTATED_VALUE", def.EnumVals[0].Name.Text)
	assertEqual(t, "onEnumValue", def.EnumVals[0].Directs.Directs[0].Name.Text)
	assertEqual(t, "OTHER_VALUE", def.EnumVals[1].Name.Text)
}

func TestSchemaDefinition(t *testing.T) {
	schema := `
schema {
	query: QueryType
	mutation: MutationType
}`
	fset := token.NewFileSet()
	p := newParser([]byte(schema), "", fset)
	def, err := p.schemaDefinition()
	if err != nil {
		t.Error("unexpected error", err)
	}
	assertEqual(t, "2:1", fset.Position(def.Pos()).String())
	assertEqual(t, "5:2", fset.Position(def.End()).String())
	assertTrue(t, len(def.OperDefns) == 2)
	assertEqual(t, "3:2", fset.Position(def.OperDefns[0].Pos()).String())
	assertEqual(t, "3:18", fset.Position(def.OperDefns[0].End()).String())
	assertEqual(t, "QueryType", def.OperDefns[0].NamedTyp.Name.Text)
	assertEqual(t, "query", def.OperDefns[0].OperType.Text)
	assertEqual(t, "MutationType", def.OperDefns[1].NamedTyp.Name.Text)
	assertEqual(t, "mutation", def.OperDefns[1].OperType.Text)
}

func TestDirectiveDefinition(t *testing.T) {
	directive := `
directive @include(if: Boolean!)
	on FIELD
	| FRAGMENT_SPREAD
	| INLINE_FRAGMENT

directive @include2(if: Boolean!) on FIELD`
	fset := token.NewFileSet()
	p := newParser([]byte(directive), "", fset)
	def, err := p.directiveDefinition()
	if err != nil {
		t.Error("unexpected error", err)
	}
	assertEqual(t, "2:1", fset.Position(def.Pos()).String())
	assertEqual(t, "5:19", fset.Position(def.End()).String())
	assertEqual(t, "include", def.Name.Text)
	assertTrue(t, def.Args != nil)
	assertTrue(t, def.Locs != nil)
	assertEqual(t, "4:2", fset.Position(def.Locs.Locs[0].Pos()).String())
	assertEqual(t, "3:5", fset.Position(def.Locs.Pos()).String())
	assertEqual(t, "5:19", fset.Position(def.Locs.End()).String())

	def, err = p.directiveDefinition()
	if err != nil {
		t.Error("unexpected error", err)
	}
	assertEqual(t, "7:1", fset.Position(def.Pos()).String())
	assertEqual(t, "7:43", fset.Position(def.End()).String())
	assertEqual(t, "include2", def.Name.Text)
}

func TestExtendDefinition(t *testing.T) {
	extend := `
extend type Foo @onType {
	seven(argument: [String]): Type
}`
	fset := token.NewFileSet()
	p := newParser([]byte(extend), "", fset)
	def, err := p.extendDefinition()
	if err != nil {
		t.Error("unexpected error", err)
	}
	assertEqual(t, "2:1", fset.Position(def.Pos()).String())
	assertEqual(t, "4:2", fset.Position(def.End()).String())
	assertEqual(t, "3:7", fset.Position(def.TypDefn.FieldDefns[0].ArgDefns.Pos()).String())
	assertEqual(t, "3:27", fset.Position(def.TypDefn.FieldDefns[0].ArgDefns.End()).String())
	assertEqual(t, "3:8", fset.Position(def.TypDefn.FieldDefns[0].ArgDefns.InputValDefns[0].Pos()).String())
	assertEqual(t, "3:26", fset.Position(def.TypDefn.FieldDefns[0].ArgDefns.InputValDefns[0].End()).String())
	assertTrue(t, def.TypDefn != nil)
}

func TestTypeDefinition(t *testing.T) {
	d := `type Human implements Character {
	id: ID!
	name: String!
	friends: [Character]
	appearsIn: [Episode]!
	starships: [Starship]
	totalCredits: Int @onField
}`
	fset := token.NewFileSet()
	p := newParser([]byte(d), "", fset)
	def, err := p.typeDefinition()
	if err != nil {
		t.Error("unexpected error", err)
	}
	assertEqual(t, "1:1", fset.Position(def.Pos()).String())
	assertEqual(t, "8:2", fset.Position(def.End()).String())
	assertEqual(t, "Human", def.Name.Text)
	assertTrue(t, def.Implements != nil)
	assertEqual(t, "1:12", fset.Position(def.Implements.Pos()).String())
	assertEqual(t, "1:32", fset.Position(def.Implements.End()).String())
	assertTrue(t, len(def.Implements.NamedTyps) == 1)
	assertEqual(t, "Character", def.Implements.NamedTyps[0].Name.Text)
	assertTrue(t, def.Directs == nil)
	assertTrue(t, len(def.FieldDefns) == 6)

	assertEqual(t, "2:2", fset.Position(def.FieldDefns[0].Pos()).String())
	assertEqual(t, "2:9", fset.Position(def.FieldDefns[0].End()).String())
	assertEqual(t, "3:2", fset.Position(def.FieldDefns[1].Pos()).String())
	assertEqual(t, "3:15", fset.Position(def.FieldDefns[1].End()).String())
	assertEqual(t, "4:2", fset.Position(def.FieldDefns[2].Pos()).String())
	assertEqual(t, "4:22", fset.Position(def.FieldDefns[2].End()).String())
	assertEqual(t, "5:2", fset.Position(def.FieldDefns[3].Pos()).String())
	assertEqual(t, "5:23", fset.Position(def.FieldDefns[3].End()).String())
	assertEqual(t, "6:2", fset.Position(def.FieldDefns[4].Pos()).String())
	assertEqual(t, "6:23", fset.Position(def.FieldDefns[4].End()).String())
	assertEqual(t, "7:2", fset.Position(def.FieldDefns[5].Pos()).String())
	assertEqual(t, "7:28", fset.Position(def.FieldDefns[5].End()).String())
}

func TestInputObjectDefinition(t *testing.T) {
	definition := `
input AnnotatedInput @onInputObjectType {
	annotatedField: Type @onField
}`

	fset := token.NewFileSet()
	p := newParser([]byte(definition), "", fset)
	def, err := p.inputObjectDefinition()
	if err != nil {
		t.Error("unexpected error", err)
	}
	assertEqual(t, "2:1", fset.Position(def.Pos()).String())
	assertEqual(t, "4:2", fset.Position(def.End()).String())
	assertTrue(t, def.Directs != nil)
	assertTrue(t, len(def.InputValDefns) == 1)
	assertEqual(t, "3:2", fset.Position(def.InputValDefns[0].Pos()).String())
	assertEqual(t, "3:31", fset.Position(def.InputValDefns[0].End()).String())
}

func TestScalarDefinition(t *testing.T) {
	scalar := `scalar CustomScalar @onScalar
scalar Date`
	fset := token.NewFileSet()
	p := newParser([]byte(scalar), "", fset)
	def, err := p.scalarDefinition()
	if err != nil {
		t.Error("unexpected error", err)
	}
	assertEqual(t, "1:1", fset.Position(def.Pos()).String())
	assertEqual(t, "1:30", fset.Position(def.End()).String())

	def, err = p.scalarDefinition()
	if err != nil {
		t.Error("unexpected error", err)
	}
	assertEqual(t, "2:1", fset.Position(def.Pos()).String())
	assertEqual(t, "2:12", fset.Position(def.End()).String())
}

func TestInterfaceDefinition(t *testing.T) {
	iface := `
interface Bar {
	one: Type
	four(argument: String = "string"): String
}

interface AnnotatedInterface @onInterface {
	annotatedField(arg: Type @onArg): Type @onField
}`
	fset := token.NewFileSet()
	p := newParser([]byte(iface), "", fset)
	def, err := p.interfaceDefinition()
	if err != nil {
		t.Error("unexpected error", err)
	}
	assertEqual(t, "2:1", fset.Position(def.Pos()).String())
	assertEqual(t, "5:2", fset.Position(def.End()).String())

	def, err = p.interfaceDefinition()
	if err != nil {
		t.Error("unexpected error", err)
	}
	assertEqual(t, "7:1", fset.Position(def.Pos()).String())
	assertEqual(t, "9:2", fset.Position(def.End()).String())
}

func TestSchema(t *testing.T) {
	schemas := `
interface Entity {
	id: ID!
	name: String
}

type User implements Entity {
	id: ID!
	name: String
	age: Int
	balance: Float
	is_active: Boolean
	friends: [User]!
}

type Root {
	me: User
	users(limit: Int = 10): [User]!
	friends(forUser: ID!, limit: Int = 5): [User]!
}

schema {
	query: Root
}`

	fset := token.NewFileSet()
	p := newParser([]byte(schemas), "", fset)
	def, err := p.schema()
	if err != nil {
		t.Error("unexpected error", err)
	}
	assertEqual(t, "2:1", fset.Position(def.Pos()).String())
	assertEqual(t, "24:2", fset.Position(def.End()).String())
	assertTrue(t, len(def.Interfaces) == 1)
	assertTrue(t, len(def.Types) == 2)
	assertTrue(t, len(def.Schemas) == 1)
}
