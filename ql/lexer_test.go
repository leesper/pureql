package ql

import (
	"reflect"
	"testing"
)

func TestComment(t *testing.T) {
	comments := []string{
		`#this is comment`,
		`# this is comment with space`,
		`#	this is comment with tab`,
		`#		this is comment with tabs		`,
	}
	for _, comment := range comments {
		lexer := NewLexerWithSource(comment)
		if tok := lexer.Read(); tok != TokenEOF {
			t.Errorf("returned: %v, expected: %v", tok, TokenEOF)
		}
	}
}

func TestCommas(t *testing.T) {
	lexer := NewLexerWithSource(",,,,, ,, ,\n")
	if tok := lexer.Read(); tok != TokenEOF {
		t.Errorf("returned: %v, expected: %v", tok, TokenEOF)
	}
}

func TestLexicalTokens(t *testing.T) {
	lexer := NewLexerWithSource(`a = 2, 47 3.14159 1e50 -6.0221413e23 "Golang\n\r\t" # and this is a comment`)
	var tokens []Token
	tok := lexer.Read()
	for tok != TokenEOF {
		tokens = append(tokens, tok)
		tok = lexer.Read()
	}

	expectedLen := 8
	if len(tokens) != expectedLen {
		t.Errorf("returned len: %d, expected: %d", len(tokens), expectedLen)
	}
}

func TestPunct(t *testing.T) {
	lexer := NewLexerWithSource(`!	$	(	)	...	:	=	@	[	]	{	|	} # ^ &`)
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

	if tokens[4].Kind != ELLIPSIS {
		t.Errorf("returned kind: %d, expected: %d", tokens[4].Kind, ELLIPSIS)
	}
}

func TestLexerKeyword(t *testing.T) {
	lexer := NewLexerWithSource(`a _b ILoveGo fooBarBAZ 1_a 2b`)
	for i := 0; i < 4; i++ {
		tok := lexer.Peek(i)
		if tok.Kind != NAME {
			t.Errorf("returned kind: %d, expected: %d", tok.Kind, NAME)
		}
	}
}

func TestParseMutation(t *testing.T) {
	source := `
mutation {
  likeStory(storyID: 12345) {
    story {
      likeCount
    }
  }
}
	`
	lexer := NewLexerWithSource(source)
	var tokens []Token

	tok := lexer.Read()
	for tok != TokenEOF {
		tokens = append(tokens, tok)
		tok = lexer.Read()
	}

	if !tokens[0].IsKeyword() {
		t.Errorf("token %v should be a keyword", tokens[0])
	}

	if !tokens[6].IsLiteral() {
		t.Errorf("token %v should be a literal", tokens[6])
	}

	if !tokens[1].IsPunct() {
		t.Errorf("token %v should be a punctuator", tokens[1])
	}

	expectedLen := 15
	if len(tokens) != expectedLen {
		t.Errorf("returned len: %d, expected: %d", len(tokens), expectedLen)
	}

	expected := []Token{
		Token{MUTATION, "mutation"},
		Token{LBRACE, "{"},
		Token{NAME, "likeStory"},
		Token{LPAREN, "("},
		Token{NAME, "storyID"},
		Token{COLON, ":"},
		Token{INT, "12345"},
		Token{RPAREN, ")"},
		Token{LBRACE, "{"},
		Token{NAME, "story"},
		Token{LBRACE, "{"},
		Token{NAME, "likeCount"},
		Token{RBRACE, "}"},
		Token{RBRACE, "}"},
		Token{RBRACE, "}"},
	}

	if !reflect.DeepEqual(tokens, expected) {
		t.Errorf("lexer scann error, returned: %v, expected: %v", tokens, expected)
	}
}

func TestParseQueryWithFragments(t *testing.T) {
	source := `
# this is a query with nested fragments
query withNestedFragments {
  user(id: 4) {
    friends(first: 10) {
      ...friendFields
    }
    mutualFriends(first: 10) {
      ...friendFields
    }
  }
}

# outer fragment
fragment friendFields on User {
  id
  name
  ...standardProfilePic
}

# inner fragment
fragment standardProfilePic on User {
  profilePic(size: 50)
}
`
	lexer := NewLexerWithSource(source)
	var tokens []Token

	tok := lexer.Read()
	for tok != TokenEOF {
		tokens = append(tokens, tok)
		tok = lexer.Read()
	}

	expectedLen := 54
	if len(tokens) != expectedLen {
		t.Errorf("returned len: %d, expected: %d", len(tokens), expectedLen)
	}

	expected := []Token{
		Token{QUERY, "query"},
		Token{NAME, "withNestedFragments"},
		Token{LBRACE, "{"},
		Token{NAME, "user"},
		Token{LPAREN, "("},
		Token{NAME, "id"},
		Token{COLON, ":"},
		Token{INT, "4"},
		Token{RPAREN, ")"},
		Token{LBRACE, "{"},
		Token{NAME, "friends"},
		Token{LPAREN, "("},
		Token{NAME, "first"},
		Token{COLON, ":"},
		Token{INT, "10"},
		Token{RPAREN, ")"},
		Token{LBRACE, "{"},
		Token{ELLIPSIS, "..."},
		Token{NAME, "friendFields"},
		Token{RBRACE, "}"},
		Token{NAME, "mutualFriends"},
		Token{LPAREN, "("},
		Token{NAME, "first"},
		Token{COLON, ":"},
		Token{INT, "10"},
		Token{RPAREN, ")"},
		Token{LBRACE, "{"},
		Token{ELLIPSIS, "..."},
		Token{NAME, "friendFields"},
		Token{RBRACE, "}"},
		Token{RBRACE, "}"},
		Token{RBRACE, "}"},
		Token{FRAGMENT, "fragment"},
		Token{NAME, "friendFields"},
		Token{ON, "on"},
		Token{NAME, "User"},
		Token{LBRACE, "{"},
		Token{NAME, "id"},
		Token{NAME, "name"},
		Token{ELLIPSIS, "..."},
		Token{NAME, "standardProfilePic"},
		Token{RBRACE, "}"},
		Token{FRAGMENT, "fragment"},
		Token{NAME, "standardProfilePic"},
		Token{ON, "on"},
		Token{NAME, "User"},
		Token{LBRACE, "{"},
		Token{NAME, "profilePic"},
		Token{LPAREN, "("},
		Token{NAME, "size"},
		Token{COLON, ":"},
		Token{INT, "50"},
		Token{RPAREN, ")"},
		Token{RBRACE, "}"},
	}

	if !reflect.DeepEqual(tokens, expected) {
		t.Errorf("lexer scann error, returned: %v, expected: %v", tokens, expected)
	}
}
