package parser

import (
	"testing"

	"github.com/leesper/pureql/gql/token"
)

func TestComment(t *testing.T) {
	comments := `#this is comment
	# this is comment with space
	#				this is comment with tabs			`
	lexer := newLexer(comments)
	if tok := lexer.Read(); tok != token.TokenEOF {
		t.Errorf("returned: %v, expected: %v", tok, token.TokenEOF)
	}
	if lexer.Line() != 3 {
		t.Errorf("returned line: %d, expected: %d", lexer.Line(), 3)
	}
}

func TestLexesPunctuators(t *testing.T) {
	lexer := newLexer("! $ ( ) ... : = @ [ ] { | }")
	tok := lexer.Read()
	expected := token.Token{Kind: token.BANG, Text: "!"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	tok = lexer.Read()
	expected = token.Token{Kind: token.DOLLAR, Text: "$"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	tok = lexer.Read()
	expected = token.Token{Kind: token.LPAREN, Text: "("}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	tok = lexer.Read()
	expected = token.Token{Kind: token.RPAREN, Text: ")"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	tok = lexer.Read()
	expected = token.Token{Kind: token.SPREAD, Text: "..."}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	tok = lexer.Read()
	expected = token.Token{Kind: token.COLON, Text: ":"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	tok = lexer.Read()
	expected = token.Token{Kind: token.EQL, Text: "="}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	tok = lexer.Read()
	expected = token.Token{Kind: token.AT, Text: "@"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	tok = lexer.Read()
	expected = token.Token{Kind: token.LBRACK, Text: "["}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	tok = lexer.Read()
	expected = token.Token{Kind: token.RBRACK, Text: "]"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	tok = lexer.Read()
	expected = token.Token{Kind: token.LBRACE, Text: "{"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	tok = lexer.Read()
	expected = token.Token{Kind: token.PIPE, Text: "|"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	tok = lexer.Read()
	expected = token.Token{Kind: token.RBRACE, Text: "}"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}
}

func TestInvalidPunctuators(t *testing.T) {
	lexer := newLexer("..")
	tok := lexer.Read()
	expected := token.Token{Kind: token.ILLEGAL, Text: ".." + string(rune(token.EOF))}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer("?")
	tok = lexer.Read()
	expected = token.Token{Kind: token.ILLEGAL, Text: "?"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer("\u203B")
	tok = lexer.Read()
	expected = token.Token{Kind: token.ILLEGAL, Text: "\u203B"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer("\u203b")
	tok = lexer.Read()
	expected = token.Token{Kind: token.ILLEGAL, Text: "\u203b"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer("ф")
	tok = lexer.Read()
	expected = token.Token{Kind: token.ILLEGAL, Text: "ф"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}
}

func TestInvalidFloat(t *testing.T) {
	lexer := newLexer(".234")
	tok := lexer.Read()
	expected := token.Token{Kind: token.ILLEGAL, Text: ".2"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}
	lexer = newLexer("..2")
	tok = lexer.Read()
	expected = token.Token{Kind: token.ILLEGAL, Text: "..2"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}
}

func TestUncommonControlChar(t *testing.T) {
	lexer := newLexer("\u0007")
	tok := lexer.Read()
	expected := token.Token{Kind: token.ILLEGAL, Text: "\u0007"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}
}

func TestBOMHeader(t *testing.T) {
	lexer := newLexer("\ufeff foo")
	tok := lexer.Read()
	expected := token.Token{Kind: token.NAME, Text: "foo"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}
}

func TestSkipWhiteSpace(t *testing.T) {
	lexer := newLexer(`
		foo
`)
	tok := lexer.Read()
	expected := token.Token{Kind: token.NAME, Text: "foo"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}
}

func TestSkipComments(t *testing.T) {
	lexer := newLexer(`
	#comment
	foo#comment
`)
	tok := lexer.Read()
	expected := token.Token{Kind: token.NAME, Text: "foo"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}
}

func TestSkipCommas(t *testing.T) {
	lexer := newLexer(",,,query,,,")
	tok := lexer.Read()
	expected := token.Token{Kind: token.NAME, Text: "query"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}
}

func TestLexesStrings(t *testing.T) {
	lexer := newLexer(`"simple"`)
	tok := lexer.Read()
	expected := token.Token{Kind: token.STRING, Text: "simple"}
	if tok != expected {
		t.Errorf("returnd: %v, expected: %v", tok, expected)
	}
	if tok.String() != `<'simple', STRING>` {
		t.Errorf("returned: %s, expected: %s", tok.String(), `<'simple', STRING>`)
	}

	lexer = newLexer(`" white space "`)
	tok = lexer.Read()
	expected = token.Token{Kind: token.STRING, Text: " white space "}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer("\"quote \\\"\"")
	tok = lexer.Read()
	expected = token.Token{Kind: token.STRING, Text: `quote "`}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer("\"escaped \\n\\r\\b\\t\\f\"")
	tok = lexer.Read()
	expected = token.Token{Kind: token.STRING, Text: "escaped \n\r\b\t\f"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer("\"slashes \\\\ \\/\"")
	tok = lexer.Read()
	expected = token.Token{Kind: token.STRING, Text: "slashes \\ /"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer("\"unicode \\u1234\\u5678\\u90AB\\uCDEF\"")
	tok = lexer.Read()
	expected = token.Token{Kind: token.STRING, Text: "unicode \u1234\u5678\u90AB\uCDEF"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer("\"unicode фы世界\"")
	tok = lexer.Read()
	expected = token.Token{Kind: token.STRING, Text: "unicode фы世界"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer("\"фы世界\"")
	tok = lexer.Read()
	expected = token.Token{Kind: token.STRING, Text: "фы世界"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer("\"Has a фы世界 multi-byte character.\"")
	tok = lexer.Read()
	expected = token.Token{Kind: token.STRING, Text: "Has a фы世界 multi-byte character."}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}
}

func TestInvalidStrings(t *testing.T) {
	lexer := newLexer("\"")
	tok := lexer.Read()
	expected := token.Token{Kind: token.ILLEGAL, Text: "\"" + string(rune(token.EOF))}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer("\"no end quote")
	tok = lexer.Read()
	expected = token.Token{Kind: token.ILLEGAL, Text: "\"no end quote" + string(rune(token.EOF))}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer("\"contains unescaped \u0007 control char\"")
	tok = lexer.Read()
	expected = token.Token{Kind: token.ILLEGAL, Text: "\"contains unescaped \u0007"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer("\"null-byte is not \u0000 end of file\"")
	tok = lexer.Read()
	expected = token.Token{Kind: token.ILLEGAL, Text: "\"null-byte is not \u0000"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer("\"multi\nline\"")
	tok = lexer.Read()
	expected = token.Token{Kind: token.ILLEGAL, Text: "\"multi\n"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer("\"multi\rline\"")
	tok = lexer.Read()
	expected = token.Token{Kind: token.ILLEGAL, Text: "\"multi\r"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer("\"bad \\z esc\"")
	tok = lexer.Read()
	expected = token.Token{Kind: token.ILLEGAL, Text: "\"bad z"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer("\"bad \\x esc\"")
	tok = lexer.Read()
	expected = token.Token{Kind: token.ILLEGAL, Text: "\"bad x"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer("\"bad \\u1 esc\"")
	tok = lexer.Read()
	expected = token.Token{Kind: token.ILLEGAL, Text: "\"bad '\\u1 es'"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer("\"bad \\u0XX1 esc\"")
	tok = lexer.Read()
	expected = token.Token{Kind: token.ILLEGAL, Text: "\"bad '\\u0XX1'"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer("\"bad \\uXXXX esc\"")
	tok = lexer.Read()
	expected = token.Token{Kind: token.ILLEGAL, Text: "\"bad '\\uXXXX'"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer("\"bad \\uFXXX esc\"")
	tok = lexer.Read()
	expected = token.Token{Kind: token.ILLEGAL, Text: "\"bad '\\uFXXX'"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer("\"bad \\uXXXF esc\"")
	tok = lexer.Read()
	expected = token.Token{Kind: token.ILLEGAL, Text: "\"bad '\\uXXXF'"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer("\"bad \\u123")
	tok = lexer.Read()
	expected = token.Token{Kind: token.ILLEGAL, Text: "\"bad '\\u123" + string(rune(token.EOF)) + "'"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer("\"bфы世ыы𠱸d \\uXXXF esc\"")
	tok = lexer.Read()
	expected = token.Token{Kind: token.ILLEGAL, Text: "\"bфы世ыы𠱸d '\\uXXXF'"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}
}

func TestLexesNumbers(t *testing.T) {
	lexer := newLexer("4")
	tok := lexer.Read()
	expected := token.Token{Kind: token.INT, Text: "4"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer("4.123")
	tok = lexer.Read()
	expected = token.Token{Kind: token.FLOAT, Text: "4.123"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer("-4")
	tok = lexer.Read()
	expected = token.Token{Kind: token.INT, Text: "-4"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer("9")
	tok = lexer.Read()
	expected = token.Token{Kind: token.INT, Text: "9"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer("0")
	tok = lexer.Read()
	expected = token.Token{Kind: token.INT, Text: "0"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer("-4.123")
	tok = lexer.Read()
	expected = token.Token{Kind: token.FLOAT, Text: "-4.123"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer("0.123")
	tok = lexer.Read()
	expected = token.Token{Kind: token.FLOAT, Text: "0.123"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer("123e4")
	tok = lexer.Read()
	expected = token.Token{Kind: token.FLOAT, Text: "123e4"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer("123E4")
	tok = lexer.Read()
	expected = token.Token{Kind: token.FLOAT, Text: "123E4"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer("123e-4")
	tok = lexer.Read()
	expected = token.Token{Kind: token.FLOAT, Text: "123e-4"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer("123e+4")
	tok = lexer.Read()
	expected = token.Token{Kind: token.FLOAT, Text: "123e+4"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer("-1.123e4")
	tok = lexer.Read()
	expected = token.Token{Kind: token.FLOAT, Text: "-1.123e4"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer("-1.123E4")
	tok = lexer.Read()
	expected = token.Token{Kind: token.FLOAT, Text: "-1.123E4"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer("-1.123e-4")
	tok = lexer.Read()
	expected = token.Token{Kind: token.FLOAT, Text: "-1.123e-4"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer("-1.123e+4")
	tok = lexer.Read()
	expected = token.Token{Kind: token.FLOAT, Text: "-1.123e+4"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer("-1.123e4567")
	tok = lexer.Read()
	expected = token.Token{Kind: token.FLOAT, Text: "-1.123e4567"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}
}

func TestInvalidNumbers(t *testing.T) {
	lexer := newLexer("00")
	tok := lexer.Read()
	expected := token.Token{Kind: token.ILLEGAL, Text: "00"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer("09")
	tok = lexer.Read()
	expected = token.Token{Kind: token.ILLEGAL, Text: "09"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer("+1")
	tok = lexer.Read()
	expected = token.Token{Kind: token.ILLEGAL, Text: "+"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer("1.")
	tok = lexer.Read()
	expected = token.Token{Kind: token.ILLEGAL, Text: "1." + string(rune(token.EOF))}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer(".123")
	tok = lexer.Read()
	expected = token.Token{Kind: token.ILLEGAL, Text: ".1"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer("1.A")
	tok = lexer.Read()
	expected = token.Token{Kind: token.ILLEGAL, Text: "1.A"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer("-A")
	tok = lexer.Read()
	expected = token.Token{Kind: token.ILLEGAL, Text: "-A"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer("1.0e")
	tok = lexer.Read()
	expected = token.Token{Kind: token.ILLEGAL, Text: "1.0e" + string(rune(token.EOF))}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}

	lexer = newLexer("1.0eA")
	tok = lexer.Read()
	expected = token.Token{Kind: token.ILLEGAL, Text: "1.0eA"}
	if tok != expected {
		t.Errorf("returned: %v, expected: %v", tok, expected)
	}
}
