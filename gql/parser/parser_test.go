package parser

import (
	"fmt"
	"testing"
)

func TestInvalidLexerOrK(t *testing.T) {
	if newParser(nil, "") != nil {
		t.Error("should return nil")
	}

	if newParser([]byte("foobar"), "") == nil {
		t.Error("should not return nil")
	}

	if err := newParser(nil, "").parseDocument(); err == nil {
		t.Error("should return error")
	}

	if err := newParser(nil, "").parseSchema(); err == nil {
		t.Error("should return error")
	}
}

func TestBadToken(t *testing.T) {
	err := ParseDocument([]byte(`
query ф {
	me {
		id
	}
}`))

	if err.Error() != "2:7: expecting {, found 'ф'" {
		t.Error(err)
	}
}

func TestParseQuery(t *testing.T) {
	err := ParseDocument([]byte(`
query _ {
	me {
		id
	}
}`))

	if err != nil {
		t.Error(err)
	}
}

func TestQueryShorthand(t *testing.T) {
	err := ParseDocument([]byte("{ field }"))
	if err != nil {
		t.Error(err)
	}
}

func TestParseInvalids(t *testing.T) {
	err := ParseDocument([]byte("{"))
	if err == nil {
		t.Errorf("expecting error, found nil")
	}
	expecting := "-: expecting NAME, found '<EOF>'"
	if err.Error() != expecting {
		t.Errorf("expecting %s, found %s", expecting, err)
	}

	err = ParseDocument([]byte(`
{ ...MissingOn }
fragment MissingOn Type
`))
	if err == nil {
		t.Error("expecting error, found nil")
	}
	expecting = "3:19: expecting on, found 'Type'"
	if err.Error() != expecting {
		t.Errorf("expecting %s, found %s", expecting, err)
	}

	err = ParseDocument([]byte(`{ field: {} }`))
	if err == nil {
		t.Error("expecting error, found nil")
	}
	expecting = "1:10: expecting NAME, found '{'"
	if err.Error() != expecting {
		t.Errorf("expecting %s, found %s", expecting, err)
	}

	err = ParseDocument([]byte(`notAnOper Foo { field }`))
	if err == nil {
		t.Error("expecting error, found nil")
	}
	expecting = "1:1: expecting query or mutation or subscription, found 'notAnOper'"
	if err.Error() != expecting {
		t.Errorf("expecting %s, found %s", expecting, err)
	}

	err = ParseDocument([]byte(`...`))
	if err == nil {
		t.Error("expecting error, found nil")
	}
	expecting = "-: expecting query or mutation or subscription, found '...'"
	if err.Error() != expecting {
		t.Errorf("expecting %s, found %s", expecting, err)
	}

	err = ParseDocument([]byte("query"))
	if err == nil {
		t.Error("expecting error, found nil")
	}
	expecting = "1:4: expecting {, found '<EOF>'"
	if err.Error() != expecting {
		t.Errorf("expecting %s, found %s", expecting, err)
	}
}

func TestParseVariableInline(t *testing.T) {
	err := ParseDocument([]byte(`{ field(complex: { a: { b: [ $var ] } }) }`))
	if err != nil {
		t.Error("unexpected error", err)
	}
}

func TestParseDefaultValues(t *testing.T) {
	err := ParseDocument([]byte(`query Foo($x: Complex = { a: { b: [ true ] } }) { field }`))
	if err != nil {
		t.Error("unexpected error", err)
	}
}

func TestParseInvalidUseOn(t *testing.T) {
	err := ParseDocument([]byte(`fragment on on on { on }`))
	if err == nil {
		t.Error("should return error")
	}
	expecting := "1:10: expecting NAME but not *on*, found 'on'"
	if err.Error() != expecting {
		t.Errorf("expecting %s, found %s", expecting, err)
	}
}

func TestParseInvalidSpreadOfOn(t *testing.T) {
	err := ParseDocument([]byte(`{ ...on }`))
	if err == nil {
		t.Error("should return error")
	}
	expecting := "1:8: expecting NAME, found '}'"
	if err.Error() != expecting {
		t.Errorf("expecting %s, found %s", expecting, err)
	}
}

func TestParseNullMisuse(t *testing.T) {
	err := ParseDocument([]byte(`{ fieldWithNullableStringInput(input: null) }`))
	if err != nil {
		t.Error("unexpected error", err)
	}
}

func TestParseQueryWithComment(t *testing.T) {
	err := ParseDocument([]byte(`
# This comment has a \u0A0A multi-byte character.
{
	field(arg: "Has a \u0A0A multi-byte character.")
}`))
	if err != nil {
		t.Error("unexpected error", err)
	}
}

func TestParseQueryWithUnicode(t *testing.T) {
	err := ParseDocument([]byte(`
# This comment has a фы世界 multi-byte character.
{ field(arg: "Has a фы世界 multi-byte character.") }`))
	if err != nil {
		t.Error("unexpected error", err)
	}
}

func TestParseQueryFile(t *testing.T) {
	err := ParseDocument([]byte(`# Copyright (c) 2015-present, Facebook, Inc.
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
}`))
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
		err := ParseDocument([]byte(query))
		if err != nil {
			t.Error("unexpected error", err, query)
		}
	}
}

func TestParseSubscription(t *testing.T) {
	err := ParseDocument([]byte(`
subscription Foo {
  subscriptionField
}`))
	if err != nil {
		t.Error("unexpected error", err)
	}
}

func TestParseUnnamedSubscription(t *testing.T) {
	err := ParseDocument([]byte(`
subscription {
  subscriptionField
}`))
	if err != nil {
		t.Error("unexpected error", err)
	}
}

func TestParseMutation(t *testing.T) {
	err := ParseDocument([]byte(`
mutation Foo {
  mutationField
}`))
	if err != nil {
		t.Error("unexpected error", err)
	}
}

func TestParseUnnamedMutation(t *testing.T) {
	err := ParseDocument([]byte(`
mutation {
  mutationField
}`))
	if err != nil {
		t.Error("unexpected error", err)
	}
}

func TestInvalidSchema(t *testing.T) {
	err := ParseSchema([]byte(`
type Character {
	name: String!
	appearsIn: [Episode]!
}

notSchema {
	foo: bar
}
`))
	if err == nil {
		t.Error("should return error")
	}

	err = ParseSchema([]byte(`
notSchema {
	foo: bar
}
`))

	if err == nil {
		t.Error("should return error")
	}
}

func TestParseSchemaFile(t *testing.T) {
	err := ParseSchema([]byte(`
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
  | INLINE_FRAGMENT`))

	if err != nil {
		t.Error("unexpected error", err)
	}
}
