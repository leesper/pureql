package ql

import "testing"

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

	// parser = NewParser(NewLexer(`{ field: {} }`))
	// err = parser.Parse()
	// if err == nil {
	// 	t.Error("should return error")
	// }
	// switch err.(type) {
	// case ErrBadParse:
	// 	// that's what we want
	// default:
	// 	t.Errorf("expecting ErrBadParse, found %v", err)
	// }
	//
	// parser = NewParser(NewLexer(`notanoperation Foo { field }`))
	// err = parser.Parse()
	// if err == nil {
	// 	t.Error("should return error")
	// }
	// switch err.(type) {
	// case ErrBadParse:
	// 	// that's what we want
	// default:
	// 	t.Errorf("expecting ErrBadParse, found %v", err)
	// }
	//
	// parser = NewParser(NewLexer("..."))
	// err = parser.Parse()
	// if err == nil {
	// 	t.Error("should return error")
	// }
	// switch err.(type) {
	// case ErrBadParse:
	// 	// that's what we want
	// default:
	// 	t.Errorf("expecting ErrBadParse, found %v", err)
	// }
	//
	// parser = NewParser(NewLexer("query"))
	// err = parser.Parse()
	// if err == nil {
	// 	t.Error("should return error")
	// }
	// switch err.(type) {
	// case ErrBadParse:
	// 	// that's what we want
	// default:
	// 	t.Errorf("expecting ErrBadParse, found %v", err)
	// }
}

// func TestParseVariableInline(t *testing.T) {
// 	parser := NewParser(NewLexer(`{ field(complex: { a: { b: [ $var ] } }) }`))
// 	err := parser.Parse()
// 	if err != nil {
// 		t.Error("unexpected error", err)
// 	}
// }
//
// func TestParseDefaultValues(t *testing.T) {
// 	parser := NewParser(NewLexer(`query Foo($x: Complex = { a: { b: [ $var ] } }) { field }`))
// 	err := parser.Parse()
// 	if err != nil {
// 		t.Error("unexpected error", err)
// 	}
// }
//
// func TestParseInvalidUseOn(t *testing.T) {
// 	parser := NewParser(NewLexer(`fragment on on on { on }`))
// 	err := parser.Parse()
// 	if err == nil {
// 		t.Error("should return error")
// 	}
// 	switch err.(type) {
// 	case ErrBadParse:
// 		// that's what we want
// 	default:
// 		t.Errorf("expecting ErrBadParse, found %v", err)
// 	}
// }
//
// func TestParseInvalidSpreadOfOn(t *testing.T) {
// 	parser := NewParser(NewLexer(`{ ...on }`))
// 	err := parser.Parse()
// 	if err == nil {
// 		t.Error("should return error")
// 	}
// 	switch err.(type) {
// 	case ErrBadParse:
// 		// that's what we want
// 	default:
// 		t.Errorf("expecting ErrBadParse, found %v", err)
// 	}
// }
//
// func TestParseNullMisuse(t *testing.T) {
// 	parser := NewParser(NewLexer(`{ fieldWithNullableStringInput(input: null) }`))
// 	err := parser.Parse()
// 	if err == nil {
// 		t.Error("should return error")
// 	}
// 	switch err.(type) {
// 	case ErrBadParse:
// 		// that's what we want
// 	default:
// 		t.Errorf("expecting ErrBadParse, found %v", err)
// 	}
// }
//
// func TestParseQueryWithComment(t *testing.T) {
// 	parser := NewParser(NewLexer(`
// 		# This comment has a \u0A0A multi-byte character.
// 		{ field(arg: "Has a \u0A0A multi-byte character.") }
// 	`))
// 	err := parser.Parse()
// 	if err != nil {
// 		t.Error("unexpected error", err)
// 	}
// }
//
// func TestParseQueryWithUnicode(t *testing.T) {
// 	parser := NewParser(NewLexer(`
// 		# This comment has a фы世界 multi-byte character.
//     { field(arg: "Has a фы世界 multi-byte character.") }
// 	`))
// 	err := parser.Parse()
// 	if err != nil {
// 		t.Error("unexpected error", err)
// 	}
// }
//
// func TestParseQueryFile(t *testing.T) {
// 	parser := NewParser(NewLexer(`# Copyright (c) 2015-present, Facebook, Inc.
// #
// # This source code is licensed under the MIT license found in the
// # LICENSE file in the root directory of this source tree.
//
// query queryName($foo: ComplexType, $site: Site = MOBILE) {
//   whoever123is: node(id: [123, 456]) {
//     id ,
//     ... on User @defer {
//       field2 {
//         id ,
//         alias: field1(first:10, after:$foo,) @include(if: $foo) {
//           id,
//           ...frag
//         }
//       }
//     }
//     ... @skip(unless: $foo) {
//       id
//     }
//     ... {
//       id
//     }
//   }
// }
//
// mutation likeStory {
//   like(story: 123) @defer {
//     story {
//       id
//     }
//   }
// }
//
// subscription StoryLikeSubscription($input: StoryLikeSubscribeInput) {
//   storyLikeSubscribe(input: $input) {
//     story {
//       likers {
//         count
//       }
//       likeSentence {
//         text
//       }
//     }
//   }
// }
//
// fragment frag on Friend {
//   foo(size: $size, bar: $b, obj: {key: "value"})
// }
//
// {
//   unnamed(truthy: true, falsey: false, nullish: null),
//   query
// }`))
// 	err := parser.Parse()
// 	if err != nil {
// 		t.Error("unexpected error", err)
// 	}
// }
//
// func TestParseNonKeywordWhereNameAllowed(t *testing.T) {
// 	nonKeywords := []string{
// 		"on",
// 		"fragment",
// 		"query",
// 		"mutation",
// 		"subscription",
// 		"true",
// 		"false",
// 	}
// 	for _, nonKeyword := range nonKeywords {
// 		fragmentName := nonKeyword
// 		if nonKeyword == "on" {
// 			nonKeyword = "alpha"
// 		}
// 		query := fmt.Sprintf(`
// 		query %v {
// 			... %v
// 			... on %v { field }
// 		}
// 		fragment %v on Type {
// 			%v(%v: $%v) @%v(%v: $%v)
// 		}
// 		`, nonKeyword, fragmentName, nonKeyword, nonKeyword, nonKeyword, nonKeyword,
// 			nonKeyword, nonKeyword, nonKeyword, nonKeyword)
// 		parser := NewParser(NewLexer(query))
// 		err := parser.Parse()
// 		if err != nil {
// 			t.Error("unexpected error", err)
// 		}
// 	}
// }
//
// func TestParseSubscription(t *testing.T) {
// 	parser := NewParser(NewLexer(`
// subscription Foo {
//   subscriptionField
// }`))
// 	err := parser.Parse()
// 	if err != nil {
// 		t.Error("unexpected error", err)
// 	}
// }
//
// func TestParseUnnamedSubscription(t *testing.T) {
// 	parser := NewParser(NewLexer(`
// subscription {
//   subscriptionField
// }`))
// 	err := parser.Parse()
// 	if err != nil {
// 		t.Error("unexpected error", err)
// 	}
// }
//
// func TestParseMutation(t *testing.T) {
// 	parser := NewParser(NewLexer(`
// mutation Foo {
//   mutationField
// }`))
// 	err := parser.Parse()
// 	if err != nil {
// 		t.Error("unexpected error", err)
// 	}
// }
//
// func TestParseUnnamedMutation(t *testing.T) {
// 	parser := NewParser(NewLexer(`
// mutation {
//   mutationField
// }`))
// 	err := parser.Parse()
// 	if err != nil {
// 		t.Error("unexpected error", err)
// 	}
// }
//
// func TestParseSchemaFile(t *testing.T) {
// 	parser := NewParser(NewLexer(`# Copyright (c) 2015-present, Facebook, Inc.
// #
// # This source code is licensed under the MIT license found in the
// # LICENSE file in the root directory of this source tree.
//
// schema {
//   query: QueryType
//   mutation: MutationType
// }
//
// type Foo implements Bar {
//   one: Type
//   two(argument: InputType!): Type
//   three(argument: InputType, other: String): Int
//   four(argument: String = "string"): String
//   five(argument: [String] = ["string", "string"]): String
//   six(argument: InputType = {key: "value"}): Type
//   seven(argument: Int = null): Type
// }
//
// type AnnotatedObject @onObject(arg: "value") {
//   annotatedField(arg: Type = "default" @onArg): Type @onField
// }
//
// interface Bar {
//   one: Type
//   four(argument: String = "string"): String
// }
//
// interface AnnotatedInterface @onInterface {
//   annotatedField(arg: Type @onArg): Type @onField
// }
//
// union Feed = Story | Article | Advert
//
// union AnnotatedUnion @onUnion = A | B
//
// union AnnotatedUnionTwo @onUnion = | A | B
//
// scalar CustomScalar
//
// scalar AnnotatedScalar @onScalar
//
// enum Site {
//   DESKTOP
//   MOBILE
// }
//
// enum AnnotatedEnum @onEnum {
//   ANNOTATED_VALUE @onEnumValue
//   OTHER_VALUE
// }
//
// input InputType {
//   key: String!
//   answer: Int = 42
// }
//
// input AnnotatedInput @onInputObjectType {
//   annotatedField: Type @onField
// }
//
// extend type Foo {
//   seven(argument: [String]): Type
// }
//
// extend type Foo @onType {}
//
// type NoFields {}
//
// directive @skip(if: Boolean!) on FIELD | FRAGMENT_SPREAD | INLINE_FRAGMENT
//
// directive @include(if: Boolean!)
//   on FIELD
//    | FRAGMENT_SPREAD
//    | INLINE_FRAGMENT
//
// directive @include2(if: Boolean!) on
//   | FIELD
//   | FRAGMENT_SPREAD
//   | INLINE_FRAGMENT`))
//
// 	err := parser.Parse()
// 	if err != nil {
// 		t.Error("unexpected error", err)
// 	}
// }
