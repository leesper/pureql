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
func TestLexerIllegal(t *testing.T)       {}
func TestLexerQueryDocument(t *testing.T) {}
func TestLexerTypeDefs(t *testing.T)      {}
