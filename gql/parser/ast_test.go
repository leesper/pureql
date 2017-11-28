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
func TestConstValues(t *testing.T)             {}
func TestFragmentDefinition(t *testing.T)      {}
func TestSelectionSet(t *testing.T)            {}
func TestOperationDefinition(t *testing.T)     {}
func TestDefinition(t *testing.T)              {}
func TestDocument(t *testing.T)                {}
func TestUnionDefinition(t *testing.T)         {}
func TestEnumDefinition(t *testing.T)          {}
func TestOperationTypeDefinition(t *testing.T) {}
func TestSchemaDefinition(t *testing.T)        {}
func TestDirectiveDefinition(t *testing.T)     {}
func TestExtendDefinition(t *testing.T)        {}
func TestTypeDefinition(t *testing.T)          {}
func TestInputObjectDefinition(t *testing.T)   {}
func TestScalarDefinition(t *testing.T)        {}
func TestInputValueDefinition(t *testing.T)    {}
func TestArgumentsDefinition(t *testing.T)     {}
func TestFieldDefinition(t *testing.T)         {}
func TestInterfaceDefinition(t *testing.T)     {}
func TestSchema(t *testing.T)                  {}
func TestVisitor(t *testing.T)                 {}
func TestInspect(t *testing.T)                 {}
