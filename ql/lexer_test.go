package ql

import "testing"

func TestComment(t *testing.T) {
	comments := []string{
		`#this is comment`,
		`# this is comment with space`,
		`#	this is comment with tab`,
		`#		this is comment with tabs		`,
	}
	for _, comment := range comments {
		lexer := NewLexer(comment)
		if tok := lexer.Read(); tok != TokenEOF {
			t.Errorf("returned: %v, expected: %v", tok, TokenEOF)
		}
	}
}

func TestPunct(t *testing.T) {
	lexer := NewLexer(`!	$	(	)	...	:	=	@	[	]	{	|	} # ^ &`)
	var tokens []Token
	tok := lexer.Read()
	for tok != TokenEOF {
		tokens = append(tokens, tok)
		tok = lexer.Read()
	}

	expectedLen := 13
	if len(tokens) != expectedLen {
		t.Errorf("returned len: %d, expected: %d", len(tokens), expectedLen)
	}

	if tokens[4].Kind != SPREAD {
		t.Errorf("returned kind: %d, expected: %d", tokens[4].Kind, SPREAD)
	}
}

func TestIllegalPunct(t *testing.T) {
	lexer := NewLexer(`^`)
	tok := lexer.Read()
	if tok.Kind != ILLEGAL {
		t.Errorf("returned kind: %d, expected: %d", tok.Kind, ILLEGAL)
	}
}

func TestInvalidFloat(t *testing.T) {
	lexer := NewLexer(".234")
	tok := lexer.Read()
	if tok.Kind != ILLEGAL {
		t.Errorf("returned kind: %d, expected: %d", tok.Kind, ILLEGAL)
	}
	lexer = NewLexer("..2")
	tok = lexer.Read()
	if tok.Kind != ILLEGAL {
		t.Errorf("returned kind: %d, expected: %d", tok.Kind, ILLEGAL)
	}
}

func TestUncommonControlChar(t *testing.T) {
	lexer := NewLexer("\u0007")
	tok := lexer.Read()
	if tok.Kind != ILLEGAL {
		t.Errorf("returned kind: %d, expected: %d", tok.Kind, ILLEGAL)
	}
}

func TestBOMHeader(t *testing.T) {
	lexer := NewLexer("\ufeff foo")
	tok := lexer.Read()
	if tok.Kind != NAME {
		t.Errorf("returned kind: %d, expected: %d", tok.Kind, NAME)
	}
}

func TestSkipWhiteSpace(t *testing.T) {
	lexer := NewLexer(`
		foo
`)
	tok := lexer.Read()
	if tok.Kind != NAME {
		t.Errorf("returned kind: %d, expected: %d", tok.Kind, NAME)
	}
}

func TestSkipComments(t *testing.T) {
	lexer := NewLexer(`
	#comment
	foo#comment
`)
	tok := lexer.Read()
	if tok.Kind != NAME {
		t.Errorf("returned kind: %d, expected: %d", tok.Kind, NAME)
	}
}

func TestSkipCommas(t *testing.T) {
	lexer := NewLexer(",,,foo,,,")
	tok := lexer.Read()
	if tok.Kind != NAME {
		t.Errorf("returned kind: %d, expected: %d", tok.Kind, NAME)
	}
}

func TestLexesStrings(t *testing.T) {
	lexer := NewLexer(`"simple"`)
	tok := lexer.Read()
	expected := Token{STRING, "simple"}
	if tok != expected {
		t.Errorf("returnd: %v, expected: %v", tok, expected)
	}
	if tok.String() != `<'simple', STRING>` {
		t.Errorf("returned: %s, expected: %s", tok.String(), `<'simple', STRING>`)
	}

	lexer = NewLexer(`" white space "`)
	tok = lexer.Read()
	expected = Token{STRING, " white space "}
	if tok != expected {
		t.Errorf("returnd: %v, expected: %v", tok, expected)
	}

	lexer = NewLexer("\"quote \\\"\"")
	tok = lexer.Read()
	expected = Token{STRING, `quote "`}
	if tok != expected {
		t.Errorf("returnd: %v, expected: %v", tok, expected)
	}

	lexer = NewLexer("\"escaped \\n\\r\\b\\t\\f\"")
	tok = lexer.Read()
	expected = Token{STRING, "escaped \n\r\b\t\f"}
	if tok != expected {
		t.Errorf("returnd: %v, expected: %v", tok, expected)
	}

	lexer = NewLexer("\"slashes \\\\ \\/\"")
	tok = lexer.Read()
	expected = Token{STRING, "slashes \\ /"}
	if tok != expected {
		t.Errorf("returnd: %v, expected: %v", tok, expected)
	}

	lexer = NewLexer("\"unicode \\u1234\\u5678\\u90AB\\uCDEF\"")
	tok = lexer.Read()
	expected = Token{STRING, "unicode \u1234\u5678\u90AB\uCDEF"}
	if tok != expected {
		t.Errorf("returnd: %v, expected: %v", tok, expected)
	}

	lexer = NewLexer("\"unicode фы世界\"")
	tok = lexer.Read()
	expected = Token{STRING, "unicode фы世界"}
	if tok != expected {
		t.Errorf("returnd: %v, expected: %v", tok, expected)
	}

	lexer = NewLexer("\"фы世界\"")
	tok = lexer.Read()
	expected = Token{STRING, "фы世界"}
	if tok != expected {
		t.Errorf("returnd: %v, expected: %v", tok, expected)
	}

	lexer = NewLexer("\"Has a фы世界 multi-byte character.\"")
	tok = lexer.Read()
	expected = Token{STRING, "Has a фы世界 multi-byte character."}
	if tok != expected {
		t.Errorf("returnd: %v, expected: %v", tok, expected)
	}
}

// func TestLexicalTokens(t *testing.T) {
// 	lexer := NewLexerWithSource(`a = 2, 47 3.14159 1e50 -6.0221413e23 "Golang\n\r\t" # and this is a comment`)
// 	var tokens []Token
// 	tok := lexer.Read()
// 	for tok != TokenEOF {
// 		tokens = append(tokens, tok)
// 		tok = lexer.Read()
// 	}
//
// 	expectedLen := 8
// 	if len(tokens) != expectedLen {
// 		t.Errorf("returned len: %d, expected: %d", len(tokens), expectedLen)
// 	}
// }
//

//
// func TestLexerKeyword(t *testing.T) {
// 	lexer := NewLexerWithSource(`a _b ILoveGo fooBarBAZ 1_a 2b`)
// 	for i := 0; i < 4; i++ {
// 		tok := lexer.Peek(i)
// 		if tok.Kind != NAME {
// 			t.Errorf("returned kind: %d, expected: %d", tok.Kind, NAME)
// 		}
// 	}
// }
//
// func TestParseMutation(t *testing.T) {
// 	source := `
// mutation {
//   likeStory(storyID: 12345) {
//     story {
//       likeCount
//     }
//   }
// }
// 	`
// 	lexer := NewLexerWithSource(source)
// 	var tokens []Token
//
// 	tok := lexer.Read()
// 	for tok != TokenEOF {
// 		tokens = append(tokens, tok)
// 		tok = lexer.Read()
// 	}
//
// 	if !tokens[0].IsKeyword() {
// 		t.Errorf("token %v should be a keyword", tokens[0])
// 	}
//
// 	if !tokens[6].IsLiteral() {
// 		t.Errorf("token %v should be a literal", tokens[6])
// 	}
//
// 	if !tokens[1].IsPunct() {
// 		t.Errorf("token %v should be a punctuator", tokens[1])
// 	}
//
// 	expectedLen := 15
// 	if len(tokens) != expectedLen {
// 		t.Errorf("returned len: %d, expected: %d", len(tokens), expectedLen)
// 	}
//
// 	expected := []Token{
// 		Token{MUTATION, "mutation"},
// 		Token{LBRACE, "{"},
// 		Token{NAME, "likeStory"},
// 		Token{LPAREN, "("},
// 		Token{NAME, "storyID"},
// 		Token{COLON, ":"},
// 		Token{INT, "12345"},
// 		Token{RPAREN, ")"},
// 		Token{LBRACE, "{"},
// 		Token{NAME, "story"},
// 		Token{LBRACE, "{"},
// 		Token{NAME, "likeCount"},
// 		Token{RBRACE, "}"},
// 		Token{RBRACE, "}"},
// 		Token{RBRACE, "}"},
// 	}
//
// 	if !reflect.DeepEqual(tokens, expected) {
// 		t.Errorf("lexer scann error, returned: %v, expected: %v", tokens, expected)
// 	}
// }
//
// func TestParseQueryWithFragments(t *testing.T) {
// 	source := `
// # this is a query with nested fragments
// query withNestedFragments {
//   user(id: 4) {
//     friends(first: 10) {
//       ...friendFields
//     }
//     mutualFriends(first: 10) {
//       ...friendFields
//     }
//   }
// }
//
// # outer fragment
// fragment friendFields on User {
//   id
//   name
//   ...standardProfilePic
// }
//
// # inner fragment
// fragment standardProfilePic on User {
//   profilePic(size: 50)
// }
// `
// 	lexer := NewLexerWithSource(source)
// 	var tokens []Token
//
// 	tok := lexer.Read()
// 	for tok != TokenEOF {
// 		tokens = append(tokens, tok)
// 		tok = lexer.Read()
// 	}
//
// 	expectedLen := 54
// 	if len(tokens) != expectedLen {
// 		t.Errorf("returned len: %d, expected: %d", len(tokens), expectedLen)
// 	}
//
// 	expected := []Token{
// 		Token{QUERY, "query"},
// 		Token{NAME, "withNestedFragments"},
// 		Token{LBRACE, "{"},
// 		Token{NAME, "user"},
// 		Token{LPAREN, "("},
// 		Token{NAME, "id"},
// 		Token{COLON, ":"},
// 		Token{INT, "4"},
// 		Token{RPAREN, ")"},
// 		Token{LBRACE, "{"},
// 		Token{NAME, "friends"},
// 		Token{LPAREN, "("},
// 		Token{NAME, "first"},
// 		Token{COLON, ":"},
// 		Token{INT, "10"},
// 		Token{RPAREN, ")"},
// 		Token{LBRACE, "{"},
// 		Token{ELLIPSIS, "..."},
// 		Token{NAME, "friendFields"},
// 		Token{RBRACE, "}"},
// 		Token{NAME, "mutualFriends"},
// 		Token{LPAREN, "("},
// 		Token{NAME, "first"},
// 		Token{COLON, ":"},
// 		Token{INT, "10"},
// 		Token{RPAREN, ")"},
// 		Token{LBRACE, "{"},
// 		Token{ELLIPSIS, "..."},
// 		Token{NAME, "friendFields"},
// 		Token{RBRACE, "}"},
// 		Token{RBRACE, "}"},
// 		Token{RBRACE, "}"},
// 		Token{FRAGMENT, "fragment"},
// 		Token{NAME, "friendFields"},
// 		Token{ON, "on"},
// 		Token{NAME, "User"},
// 		Token{LBRACE, "{"},
// 		Token{NAME, "id"},
// 		Token{NAME, "name"},
// 		Token{ELLIPSIS, "..."},
// 		Token{NAME, "standardProfilePic"},
// 		Token{RBRACE, "}"},
// 		Token{FRAGMENT, "fragment"},
// 		Token{NAME, "standardProfilePic"},
// 		Token{ON, "on"},
// 		Token{NAME, "User"},
// 		Token{LBRACE, "{"},
// 		Token{NAME, "profilePic"},
// 		Token{LPAREN, "("},
// 		Token{NAME, "size"},
// 		Token{COLON, ":"},
// 		Token{INT, "50"},
// 		Token{RPAREN, ")"},
// 		Token{RBRACE, "}"},
// 	}
//
// 	if !reflect.DeepEqual(tokens, expected) {
// 		t.Errorf("lexer scann error, returned: %v, expected: %v", tokens, expected)
// 	}
// }
