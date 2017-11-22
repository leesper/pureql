package parser

import "testing"

func TestComment(t *testing.T) {
	comments := []byte(`#this is comment
	# this is comment with space
	#				this is comment with tabs			`)
	lexer := newLexer(comments, nil)
	if tok := lexer.read(); tok != TokenEOF {
		t.Errorf("returned: %v, expected: %v", tok, TokenEOF)
	}
	if lexer.line() != 3 {
		t.Errorf("returned line: %d, expected: %d", lexer.line(), 3)
	}
}

func TestLexesPunctuators(t *testing.T) {
	lexer := newLexer([]byte("! $ ( ) ... : = @ [ ] { | }"), nil)
	tok := lexer.read()
	expected := Token{Kind: BANG, Text: "!"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	tok = lexer.read()
	expected = Token{Kind: DOLLAR, Text: "$"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	tok = lexer.read()
	expected = Token{Kind: LPAREN, Text: "("}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	tok = lexer.read()
	expected = Token{Kind: RPAREN, Text: ")"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	tok = lexer.read()
	expected = Token{Kind: SPREAD, Text: "..."}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	tok = lexer.read()
	expected = Token{Kind: COLON, Text: ":"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	tok = lexer.read()
	expected = Token{Kind: EQL, Text: "="}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	tok = lexer.read()
	expected = Token{Kind: AT, Text: "@"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	tok = lexer.read()
	expected = Token{Kind: LBRACK, Text: "["}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	tok = lexer.read()
	expected = Token{Kind: RBRACK, Text: "]"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	tok = lexer.read()
	expected = Token{Kind: LBRACE, Text: "{"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	tok = lexer.read()
	expected = Token{Kind: PIPE, Text: "|"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	tok = lexer.read()
	expected = Token{Kind: RBRACE, Text: "}"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}
}

func TestInvalidPunctuators(t *testing.T) {
	lexer := newLexer([]byte(".."), nil)
	tok := lexer.read()
	expected := Token{Kind: ILLEGAL, Text: ".." + string(rune(EOF))}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer([]byte("?"), nil)
	tok = lexer.read()
	expected = Token{Kind: ILLEGAL, Text: "?"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer([]byte("\u203B"), nil)
	tok = lexer.read()
	expected = Token{Kind: ILLEGAL, Text: "\u203B"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer([]byte("\u203b"), nil)
	tok = lexer.read()
	expected = Token{Kind: ILLEGAL, Text: "\u203b"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer([]byte("ф"), nil)
	tok = lexer.read()
	expected = Token{Kind: ILLEGAL, Text: "ф"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}
}

func TestInvalidFloat(t *testing.T) {
	lexer := newLexer([]byte(".234"), nil)
	tok := lexer.read()
	expected := Token{Kind: ILLEGAL, Text: ".2"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}
	lexer = newLexer([]byte("..2"), nil)
	tok = lexer.read()
	expected = Token{Kind: ILLEGAL, Text: "..2"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}
}

func TestUncommonControlChar(t *testing.T) {
	lexer := newLexer([]byte("\u0007"), nil)
	tok := lexer.read()
	expected := Token{Kind: ILLEGAL, Text: "\u0007"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}
}

func TestBOMHeader(t *testing.T) {
	lexer := newLexer([]byte("\ufeff foo"), nil)
	tok := lexer.read()
	expected := Token{Kind: NAME, Text: "foo"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}
}

func TestSkipWhiteSpace(t *testing.T) {
	lexer := newLexer([]byte(`
		foo
`), nil)
	tok := lexer.read()
	expected := Token{Kind: NAME, Text: "foo"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}
}

func TestSkipComments(t *testing.T) {
	lexer := newLexer([]byte(`
	#comment
	foo#comment
`), nil)
	tok := lexer.read()
	expected := Token{Kind: NAME, Text: "foo"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}
}

func TestSkipCommas(t *testing.T) {
	lexer := newLexer([]byte(",,,query,,,"), nil)
	tok := lexer.read()
	expected := Token{Kind: NAME, Text: "query"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}
}

func TestLexesStrings(t *testing.T) {
	lexer := newLexer([]byte(`"simple"`), nil)
	tok := lexer.read()
	expected := Token{Kind: STRING, Text: "simple"}
	if tok != expected {
		t.Errorf("returnd: %v, expected: %v", tok, expected)
	}
	if tok.String() != `<'simple', STRING>` {
		t.Errorf("returned: %s, expected: %s", tok.String(), `<'simple', STRING>`)
	}

	lexer = newLexer([]byte(`" white space "`), nil)
	tok = lexer.read()
	expected = Token{Kind: STRING, Text: " white space "}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer([]byte("\"quote \\\"\""), nil)
	tok = lexer.read()
	expected = Token{Kind: STRING, Text: `quote "`}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer([]byte("\"escaped \\n\\r\\b\\t\\f\""), nil)
	tok = lexer.read()
	expected = Token{Kind: STRING, Text: "escaped \n\r\b\t\f"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer([]byte("\"slashes \\\\ \\/\""), nil)
	tok = lexer.read()
	expected = Token{Kind: STRING, Text: "slashes \\ /"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer([]byte("\"unicode \\u1234\\u5678\\u90AB\\uCDEF\""), nil)
	tok = lexer.read()
	expected = Token{Kind: STRING, Text: "unicode \u1234\u5678\u90AB\uCDEF"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer([]byte("\"unicode фы世界\""), nil)
	tok = lexer.read()
	expected = Token{Kind: STRING, Text: "unicode фы世界"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer([]byte("\"фы世界\""), nil)
	tok = lexer.read()
	expected = Token{Kind: STRING, Text: "фы世界"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer([]byte("\"Has a фы世界 multi-byte character.\""), nil)
	tok = lexer.read()
	expected = Token{Kind: STRING, Text: "Has a фы世界 multi-byte character."}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}
}

func TestInvalidStrings(t *testing.T) {
	lexer := newLexer([]byte("\""), nil)
	tok := lexer.read()
	expected := Token{Kind: ILLEGAL, Text: "\"" + string(rune(EOF))}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer([]byte("\"no end quote"), nil)
	tok = lexer.read()
	expected = Token{Kind: ILLEGAL, Text: "\"no end quote" + string(rune(EOF))}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer([]byte("\"contains unescaped \u0007 control char\""), nil)
	tok = lexer.read()
	expected = Token{Kind: ILLEGAL, Text: "\"contains unescaped \u0007"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer([]byte("\"null-byte is not \u0000 end of file\""), nil)
	tok = lexer.read()
	expected = Token{Kind: ILLEGAL, Text: "\"null-byte is not \u0000"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer([]byte("\"multi\nline\""), nil)
	tok = lexer.read()
	expected = Token{Kind: ILLEGAL, Text: "\"multi\n"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer([]byte("\"multi\rline\""), nil)
	tok = lexer.read()
	expected = Token{Kind: ILLEGAL, Text: "\"multi\r"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer([]byte("\"bad \\z esc\""), nil)
	tok = lexer.read()
	expected = Token{Kind: ILLEGAL, Text: "\"bad z"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer([]byte("\"bad \\x esc\""), nil)
	tok = lexer.read()
	expected = Token{Kind: ILLEGAL, Text: "\"bad x"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer([]byte("\"bad \\u1 esc\""), nil)
	tok = lexer.read()
	expected = Token{Kind: ILLEGAL, Text: "\"bad '\\u1 es'"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer([]byte("\"bad \\u0XX1 esc\""), nil)
	tok = lexer.read()
	expected = Token{Kind: ILLEGAL, Text: "\"bad '\\u0XX1'"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer([]byte("\"bad \\uXXXX esc\""), nil)
	tok = lexer.read()
	expected = Token{Kind: ILLEGAL, Text: "\"bad '\\uXXXX'"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer([]byte("\"bad \\uFXXX esc\""), nil)
	tok = lexer.read()
	expected = Token{Kind: ILLEGAL, Text: "\"bad '\\uFXXX'"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer([]byte("\"bad \\uXXXF esc\""), nil)
	tok = lexer.read()
	expected = Token{Kind: ILLEGAL, Text: "\"bad '\\uXXXF'"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer([]byte("\"bad \\u123"), nil)
	tok = lexer.read()
	expected = Token{Kind: ILLEGAL, Text: "\"bad '\\u123" + string(rune(EOF)) + "'"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer([]byte("\"bфы世ыы𠱸d \\uXXXF esc\""), nil)
	tok = lexer.read()
	expected = Token{Kind: ILLEGAL, Text: "\"bфы世ыы𠱸d '\\uXXXF'"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}
}

func TestLexesNumbers(t *testing.T) {
	lexer := newLexer([]byte("4"), nil)
	tok := lexer.read()
	expected := Token{Kind: INT, Text: "4"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer([]byte("4.123"), nil)
	tok = lexer.read()
	expected = Token{Kind: FLOAT, Text: "4.123"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer([]byte("-4"), nil)
	tok = lexer.read()
	expected = Token{Kind: INT, Text: "-4"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer([]byte("9"), nil)
	tok = lexer.read()
	expected = Token{Kind: INT, Text: "9"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer([]byte("0"), nil)
	tok = lexer.read()
	expected = Token{Kind: INT, Text: "0"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer([]byte("-4.123"), nil)
	tok = lexer.read()
	expected = Token{Kind: FLOAT, Text: "-4.123"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer([]byte("0.123"), nil)
	tok = lexer.read()
	expected = Token{Kind: FLOAT, Text: "0.123"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer([]byte("123e4"), nil)
	tok = lexer.read()
	expected = Token{Kind: FLOAT, Text: "123e4"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer([]byte("123E4"), nil)
	tok = lexer.read()
	expected = Token{Kind: FLOAT, Text: "123E4"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer([]byte("123e-4"), nil)
	tok = lexer.read()
	expected = Token{Kind: FLOAT, Text: "123e-4"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer([]byte("123e+4"), nil)
	tok = lexer.read()
	expected = Token{Kind: FLOAT, Text: "123e+4"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer([]byte("-1.123e4"), nil)
	tok = lexer.read()
	expected = Token{Kind: FLOAT, Text: "-1.123e4"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer([]byte("-1.123E4"), nil)
	tok = lexer.read()
	expected = Token{Kind: FLOAT, Text: "-1.123E4"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer([]byte("-1.123e-4"), nil)
	tok = lexer.read()
	expected = Token{Kind: FLOAT, Text: "-1.123e-4"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer([]byte("-1.123e+4"), nil)
	tok = lexer.read()
	expected = Token{Kind: FLOAT, Text: "-1.123e+4"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer([]byte("-1.123e4567"), nil)
	tok = lexer.read()
	expected = Token{Kind: FLOAT, Text: "-1.123e4567"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}
}

func TestInvalidNumbers(t *testing.T) {
	lexer := newLexer([]byte("00"), nil)
	tok := lexer.read()
	expected := Token{Kind: ILLEGAL, Text: "00"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer([]byte("09"), nil)
	tok = lexer.read()
	expected = Token{Kind: ILLEGAL, Text: "09"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer([]byte("+1"), nil)
	tok = lexer.read()
	expected = Token{Kind: ILLEGAL, Text: "+"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer([]byte("1."), nil)
	tok = lexer.read()
	expected = Token{Kind: ILLEGAL, Text: "1." + string(rune(EOF))}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer([]byte(".123"), nil)
	tok = lexer.read()
	expected = Token{Kind: ILLEGAL, Text: ".1"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer([]byte("1.A"), nil)
	tok = lexer.read()
	expected = Token{Kind: ILLEGAL, Text: "1.A"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer([]byte("-A"), nil)
	tok = lexer.read()
	expected = Token{Kind: ILLEGAL, Text: "-A"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer([]byte("1.0e"), nil)
	tok = lexer.read()
	expected = Token{Kind: ILLEGAL, Text: "1.0e" + string(rune(EOF))}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer([]byte("1.0eA"), nil)
	tok = lexer.read()
	expected = Token{Kind: ILLEGAL, Text: "1.0eA"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}
}
