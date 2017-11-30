package ast

import (
	"fmt"
	"go/token"
	"testing"
)

func TestInvalidLexerOrK(t *testing.T) {
	if newParser(nil, "", token.NewFileSet()) != nil {
		t.Error("should return nil")
	}

	if newParser([]byte("foobar"), "", token.NewFileSet()) == nil {
		t.Error("should not return nil")
	}

	if _, err := newParser(nil, "", token.NewFileSet()).parseDocument(); err == nil {
		t.Error("should return error")
	}

	if _, err := newParser(nil, "", token.NewFileSet()).parseSchema(); err == nil {
		t.Error("should return error")
	}
}

func TestBadToken(t *testing.T) {
	_, err := ParseDocument([]byte(`
query ф {
	me {
		id
	}
}`), "", token.NewFileSet())

	if err.Error() != "2:8: expecting {, found 'ф'" {
		t.Error(err)
	}
}

func TestParseQuery(t *testing.T) {
	_, err := ParseDocument([]byte(`
query _ {
	me {
		id
	}
}`), "", token.NewFileSet())

	if err != nil {
		t.Error(err)
	}
}

func TestQueryShorthand(t *testing.T) {
	_, err := ParseDocument([]byte("{ field }"), "", token.NewFileSet())
	if err != nil {
		t.Error(err)
	}
}

func TestParseInvalids(t *testing.T) {
	_, err := ParseDocument([]byte("{"), "", token.NewFileSet())
	if err == nil {
		t.Errorf("expecting error, found nil")
	}
	expecting := "1:1: expecting NAME, found '<EOF>'"
	if err.Error() != expecting {
		t.Errorf("expecting %s, found %s", expecting, err)
	}

	_, err = ParseDocument([]byte(`
{ ...MissingOn }
fragment MissingOn Type
`), "", token.NewFileSet())
	if err == nil {
		t.Error("expecting error, found nil")
	}
	expecting = "3:20: expecting on, found 'Type'"
	if err.Error() != expecting {
		t.Errorf("expecting %s, found %s", expecting, err)
	}

	_, err = ParseDocument([]byte(`{ field: {} }`), "", token.NewFileSet())
	if err == nil {
		t.Error("expecting error, found nil")
	}
	expecting = "1:10: expecting NAME, found '{'"
	if err.Error() != expecting {
		t.Errorf("expecting %s, found %s", expecting, err)
	}

	_, err = ParseDocument([]byte(`notAnOper Foo { field }`), "", token.NewFileSet())
	if err == nil {
		t.Error("expecting error, found nil")
	}
	expecting = "1:1: expecting query or mutation or subscription, found 'notAnOper'"
	if err.Error() != expecting {
		t.Errorf("expecting %s, found %s", expecting, err)
	}

	_, err = ParseDocument([]byte(`...`), "", token.NewFileSet())
	if err == nil {
		t.Error("expecting error, found nil")
	}
	expecting = "1:1: expecting query or mutation or subscription, found '...'"
	if err.Error() != expecting {
		t.Errorf("expecting %s, found %s", expecting, err)
	}

	_, err = ParseDocument([]byte("query"), "", token.NewFileSet())
	if err == nil {
		t.Error("expecting error, found nil")
	}
	expecting = "1:5: expecting {, found '<EOF>'"
	if err.Error() != expecting {
		t.Errorf("expecting %s, found %s", expecting, err)
	}
}

func TestParseVariableInline(t *testing.T) {
	_, err := ParseDocument([]byte(`{ field(complex: { a: { b: [ $var ] } }) }`), "", token.NewFileSet())
	if err != nil {
		t.Error("unexpected error", err)
	}
}

func TestParseDefaultValues(t *testing.T) {
	_, err := ParseDocument([]byte(`query Foo($x: Complex = { a: { b: [ true ] } }) { field }`), "", token.NewFileSet())
	if err != nil {
		t.Error("unexpected error", err)
	}
}

func TestParseInvalidUseOn(t *testing.T) {
	_, err := ParseDocument([]byte(`fragment on on on { on }`), "", token.NewFileSet())
	if err == nil {
		t.Error("should return error")
	}
	expecting := "1:10: expecting NAME but not *on*, found 'on'"
	if err.Error() != expecting {
		t.Errorf("expecting %s, found %s", expecting, err)
	}
}

func TestParseInvalidSpreadOfOn(t *testing.T) {
	_, err := ParseDocument([]byte(`{ ...on }`), "", token.NewFileSet())
	if err == nil {
		t.Error("should return error")
	}
	expecting := "1:9: expecting NAME, found '}'"
	if err.Error() != expecting {
		t.Errorf("expecting %s, found %s", expecting, err)
	}
}

func TestParseNullMisuse(t *testing.T) {
	_, err := ParseDocument([]byte(`{ fieldWithNullableStringInput(input: null) }`), "", token.NewFileSet())
	if err != nil {
		t.Error("unexpected error", err)
	}
}

func TestParseQueryWithComment(t *testing.T) {
	_, err := ParseDocument([]byte(`
# This comment has a \u0A0A multi-byte character.
{
	field(arg: "Has a \u0A0A multi-byte character.")
}`), "", token.NewFileSet())
	if err != nil {
		t.Error("unexpected error", err)
	}
}

func TestParseQueryWithUnicode(t *testing.T) {
	_, err := ParseDocument([]byte(`
# This comment has a фы世界 multi-byte character.
{ field(arg: "Has a фы世界 multi-byte character.") }`), "", token.NewFileSet())
	if err != nil {
		t.Error("unexpected error", err)
	}
}

func TestParseQueryFile(t *testing.T) {
	_, err := ParseDocument([]byte(`# Copyright (c) 2015-present, Facebook, Inc.
#
# This source code is licensed under the MIT license found in the
# LICENSE file in the root directory of this source tree.

query queryName($foo: ComplexType, $site: Site = MOBILE) {
	whoever123is: node(id: [123, 456]) {
		id ,
		... on User @defer {
      field2 {
        id ,
        alias: field1(first:10, after:$foo,) @include(if: $foo) {
          id,
          ...frag
        }
      }
    }
    ... @skip(unless: $foo) {
      id
    }
    ... {
      id
    }
  }
}

mutation likeStory {
  like(story: 123) @defer {
    story {
      id
    }
  }
}

subscription StoryLikeSubscription($input: StoryLikeSubscribeInput) {
  storyLikeSubscribe(input: $input) {
    story {
      likers {
        count
      }
      likeSentence {
        text
      }
    }
  }
}

fragment frag on Friend {
  foo(size: $size, bar: $b, obj: {key: "value"})
}

{
  unnamed(truthy: true, falsey: false, nullish: null),
  query
}`), "", token.NewFileSet())
	if err != nil {
		t.Error("unexpected error", err)
	}
}

func TestParseNonKeywordWhereNameAllowed(t *testing.T) {
	nonKeywords := []string{
		"on",
		"fragment",
		"query",
		"mutation",
		"subscription",
		"true",
		"false",
	}
	for _, nonKeyword := range nonKeywords {
		fragmentName := nonKeyword
		if nonKeyword == "on" {
			fragmentName = "alpha"
		}
		query := fmt.Sprintf(`
query %v {
	... %v
	... on %v { field }
}
fragment %v on Type {
	%v(%v: $%v) @%v(%v: $%v)
}`, nonKeyword, fragmentName, nonKeyword, fragmentName, nonKeyword, nonKeyword,
			nonKeyword, nonKeyword, nonKeyword, nonKeyword)
		_, err := ParseDocument([]byte(query), "", token.NewFileSet())
		if err != nil {
			t.Error("unexpected error", err, query)
		}
	}
}

func TestParseSubscription(t *testing.T) {
	_, err := ParseDocument([]byte(`
subscription Foo {
  subscriptionField
}`), "", token.NewFileSet())
	if err != nil {
		t.Error("unexpected error", err)
	}
}

func TestParseUnnamedSubscription(t *testing.T) {
	_, err := ParseDocument([]byte(`
subscription {
  subscriptionField
}`), "", token.NewFileSet())
	if err != nil {
		t.Error("unexpected error", err)
	}
}

func TestParseMutation(t *testing.T) {
	_, err := ParseDocument([]byte(`
mutation Foo {
  mutationField
}`), "", token.NewFileSet())
	if err != nil {
		t.Error("unexpected error", err)
	}
}

func TestParseUnnamedMutation(t *testing.T) {
	_, err := ParseDocument([]byte(`
mutation {
  mutationField
}`), "", token.NewFileSet())
	if err != nil {
		t.Error("unexpected error", err)
	}
}

func TestInvalidSchema(t *testing.T) {
	_, err := ParseSchema([]byte(`
type Character {
	name: String!
	appearsIn: [Episode]!
}

notSchema {
	foo: bar
}
`), "", token.NewFileSet())
	if err == nil {
		t.Error("should return error")
	}

	_, err = ParseSchema([]byte(`
notSchema {
	foo: bar
}
`), "", token.NewFileSet())

	if err == nil {
		t.Error("should return error")
	}
}

func TestParseSchemaFile(t *testing.T) {
	_, err := ParseSchema([]byte(`
# Copyright (c) 2015-present, Facebook, Inc.
#
# This source code is licensed under the MIT license found in the
# LICENSE file in the root directory of this source tree.

schema {
  query: QueryType
  mutation: MutationType
}

type Foo implements Bar {
  one: Type
  two(argument: InputType!): Type
  three(argument: InputType, other: String): Int
  four(argument: String = "string"): String
  five(argument: [String] = ["string", "string"]): String
  six(argument: InputType = {key: "value"}): Type
  seven(argument: Int = null): Type
}

type AnnotatedObject @onObject(arg: "value") {
  annotatedField(arg: Type = "default" @onArg): Type @onField
}

interface Bar {
  one: Type
  four(argument: String = "string"): String
}

interface AnnotatedInterface @onInterface {
  annotatedField(arg: Type @onArg): Type @onField
}

union Feed = Story | Article | Advert

union AnnotatedUnion @onUnion = A | B

union AnnotatedUnionTwo @onUnion = A | B

scalar CustomScalar

scalar AnnotatedScalar @onScalar

enum Site {
  DESKTOP
  MOBILE
}

enum AnnotatedEnum @onEnum {
  ANNOTATED_VALUE @onEnumValue
  OTHER_VALUE
}

input InputType {
  key: String!
  answer: Int = 42
}

input AnnotatedInput @onInputObjectType {
  annotatedField: Type @onField
}

extend type Foo {
  seven(argument: [String]): Type
}

extend type Foo @onType {}

type NoFields {}

directive @skip(if: Boolean!) on FIELD | FRAGMENT_SPREAD | INLINE_FRAGMENT

directive @include(if: Boolean!)
  on FIELD
   | FRAGMENT_SPREAD
   | INLINE_FRAGMENT

directive @include2(if: Boolean!) on
  FIELD
  | FRAGMENT_SPREAD
  | INLINE_FRAGMENT`), "", token.NewFileSet())

	if err != nil {
		t.Error("unexpected error", err)
	}
}
