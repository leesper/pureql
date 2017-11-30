package ast

import (
	"fmt"
	"go/token"
)

func ExampleInspect() {
	schemas := `schema {
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

interface Bar {
	one: Type
	four(argument: String = "string"): String
}

union Feed = Story | Article | Advert

scalar CustomScalar

enum Site {
	DESKTOP
	MOBILE
}

input InputType {
	key: String!
	answer: Int = 42
}

extend type Foo {
	seven(argument: [String]): Type
}

directive @include2(if: Boolean!) on
	FIELD
	| FRAGMENT_SPREAD
	| INLINE_FRAGMENT`

	fset := token.NewFileSet()
	p := newParser([]byte(schemas), "", fset)
	def, _ := p.parseSchema()

	Inspect(def, func(n Node) bool {
		switch n.(type) {
		case *Schema:
			fmt.Printf("%s:\t%s\n", fset.Position(n.Pos()), "schemas detected")
		case *SchemaDefinition:
			fmt.Printf("%s:\t%s\n", fset.Position(n.Pos()), "schema detected")
		case *TypeDefinition:
			fmt.Printf("%s:\t%s\n", fset.Position(n.Pos()), "type detected")
		case *InterfaceDefinition:
			fmt.Printf("%s:\t%s\n", fset.Position(n.Pos()), "interface detected")
		case *UnionDefinition:
			fmt.Printf("%s:\t%s\n", fset.Position(n.Pos()), "union detected")
		case *ScalarDefinition:
			fmt.Printf("%s:\t%s\n", fset.Position(n.Pos()), "scalar detected")
		case *EnumDefinition:
			fmt.Printf("%s:\t%s\n", fset.Position(n.Pos()), "enum detected")
		case *InputObjectDefinition:
			fmt.Printf("%s:\t%s\n", fset.Position(n.Pos()), "input object detected")
		case *ExtendDefinition:
			fmt.Printf("%s:\t%s\n", fset.Position(n.Pos()), "extend detected")
		case *DirectiveDefinition:
			fmt.Printf("%s:\t%s\n", fset.Position(n.Pos()), "directive detected")
		}

		return true
	})

	// Output:
	// 1:1:	schemas detected
	// 16:1:	interface detected
	// 23:1:	scalar detected
	// 30:1:	input object detected
	// 6:1:	type detected
	// 35:1:	extend detected
	// 35:8:	type detected
	// 39:1:	directive detected
	// 1:1:	schema detected
	// 25:1:	enum detected
	// 21:1:	union detected
}

func ExampleWalk() {
	document := `
query queryName($foo: ComplexType, $site: Site = MOBILE) {
  whoever123is: node(id: [123, 456]) {
    id ,
    ... on User @defer {
      field2 {
        id,
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
}`

	fset := token.NewFileSet()
	p := newParser([]byte(document), "", fset)
	doc, _ := p.parseDocument()

	Inspect(doc, func(n Node) bool {
		switch n.(type) {
		case *Document:
			fmt.Printf("%s:\t%s\n", fset.Position(n.Pos()), "document detected")
		case *OperationDefinition:
			fmt.Printf("%s:\t%s\n", fset.Position(n.Pos()), "operation definition detected")
		case *FragmentDefinition:
			fmt.Printf("%s:\t%s\n", fset.Position(n.Pos()), "fragment definition detected")

		}

		return true
	})

	// Output:
	// 2:1:	document detected
	// 2:1:	operation definition detected
	// 23:1:	operation definition detected
	// 31:1:	operation definition detected
	// 44:1:	fragment definition detected
	// 48:1:	operation definition detected
}
