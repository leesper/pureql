package ql

import (
	"fmt"
	"testing"
)

func TestBadToken(t *testing.T) {
	err := ParseDocument(`
query ф {
	me {
		id
	}
}`)

	if err.Error() != "line 2: expecting {, found <'ф', ILLEGAL>" {
		t.Error(err)
	}
}

func TestParseQuery(t *testing.T) {
	err := ParseDocument(`
query _ {
	me {
		id
	}
}`)

	if err != nil {
		t.Error(err)
	}
}

func TestQueryShorthand(t *testing.T) {
	err := ParseDocument("{ field }")
	if err != nil {
		t.Error(err)
	}
}

func TestParseInvalids(t *testing.T) {
	err := ParseDocument("{")
	if err == nil {
		t.Errorf("expecting error, found nil")
	}
	expecting := "line 1: expecting NAME, found <'EOF', EOF>"
	if err.Error() != expecting {
		t.Errorf("expecting %s, found %s", expecting, err)
	}

	err = ParseDocument(`
{ ...MissingOn }
fragment MissingOn Type
`)
	if err == nil {
		t.Error("expecting error, found nil")
	}
	expecting = "line 3: expecting on, found <'Type', NAME>"
	if err.Error() != expecting {
		t.Errorf("expecting %s, found %s", expecting, err)
	}

	err = ParseDocument(`{ field: {} }`)
	if err == nil {
		t.Error("expecting error, found nil")
	}
	expecting = "line 1: expecting NAME, found <'{', {>"
	if err.Error() != expecting {
		t.Errorf("expecting %s, found %s", expecting, err)
	}

	err = ParseDocument(`notAnOper Foo { field }`)
	if err == nil {
		t.Error("expecting error, found nil")
	}
	expecting = "line 1: expecting query or mutation, found <'notAnOper', NAME>"
	if err.Error() != expecting {
		t.Errorf("expecting %s, found %s", expecting, err)
	}

	err = ParseDocument(`...`)
	if err == nil {
		t.Error("expecting error, found nil")
	}
	expecting = "line 1: expecting query or mutation, found <'...', ...>"
	if err.Error() != expecting {
		t.Errorf("expecting %s, found %s", expecting, err)
	}

	err = ParseDocument("query")
	if err == nil {
		t.Error("expecting error, found nil")
	}
	expecting = "line 1: expecting {, found <'EOF', EOF>"
	if err.Error() != expecting {
		t.Errorf("expecting %s, found %s", expecting, err)
	}
}

func TestParseVariableInline(t *testing.T) {
	err := ParseDocument(`{ field(complex: { a: { b: [ $var ] } }) }`)
	if err != nil {
		t.Error("unexpected error", err)
	}
}

func TestParseDefaultValues(t *testing.T) {
	err := ParseDocument(`query Foo($x: Complex = { a: { b: [ true ] } }) { field }`)
	if err != nil {
		t.Error("unexpected error", err)
	}
}

func TestParseInvalidUseOn(t *testing.T) {
	err := ParseDocument(`fragment on on on { on }`)
	if err == nil {
		t.Error("should return error")
	}
	expecting := "line 1: expecting NAME but not *on*, found <'on', NAME>"
	if err.Error() != expecting {
		t.Errorf("expecting %s, found %s", expecting, err)
	}
}

func TestParseInvalidSpreadOfOn(t *testing.T) {
	err := ParseDocument(`{ ...on }`)
	if err == nil {
		t.Error("should return error")
	}
	expecting := "line 1: expecting NAME, found <'}', }>"
	if err.Error() != expecting {
		t.Errorf("expecting %s, found %s", expecting, err)
	}
}

func TestParseNullMisuse(t *testing.T) {
	err := ParseDocument(`{ fieldWithNullableStringInput(input: null) }`)
	if err != nil {
		t.Error("unexpected error", err)
	}
}

func TestParseQueryWithComment(t *testing.T) {
	err := ParseDocument(`
# This comment has a \u0A0A multi-byte character.
{
	field(arg: "Has a \u0A0A multi-byte character.")
}`)
	if err != nil {
		t.Error("unexpected error", err)
	}
}

func TestParseQueryWithUnicode(t *testing.T) {
	err := ParseDocument(`
# This comment has a фы世界 multi-byte character.
{ field(arg: "Has a фы世界 multi-byte character.") }`)
	if err != nil {
		t.Error("unexpected error", err)
	}
}

func TestParseQueryFile(t *testing.T) {
	err := ParseDocument(`# Copyright (c) 2015-present, Facebook, Inc.
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
}`)
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
		err := ParseDocument(query)
		if err != nil {
			t.Error("unexpected error", err, query)
		}
	}
}

func TestParseSubscription(t *testing.T) {
	err := ParseDocument(`
subscription Foo {
  subscriptionField
}`)
	if err != nil {
		t.Error("unexpected error", err)
	}
}

func TestParseUnnamedSubscription(t *testing.T) {
	err := ParseDocument(`
subscription {
  subscriptionField
}`)
	if err != nil {
		t.Error("unexpected error", err)
	}
}

func TestParseMutation(t *testing.T) {
	err := ParseDocument(`
mutation Foo {
  mutationField
}`)
	if err != nil {
		t.Error("unexpected error", err)
	}
}

func TestParseUnnamedMutation(t *testing.T) {
	err := ParseDocument(`
mutation {
  mutationField
}`)
	if err != nil {
		t.Error("unexpected error", err)
	}
}

func TestParseSchemaFile(t *testing.T) {
	err := ParseSchema(`
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
  | INLINE_FRAGMENT`)

	if err != nil {
		t.Error("unexpected error", err)
	}
}
