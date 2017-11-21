package parser

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/leesper/pureql/gql/token"
)

type lexer struct {
	input     *bufio.Reader
	lookAhead rune
	line      int
}

func newLexer(source string) *lexer {
	source = strings.TrimRight(source, "\n\t\r ")
	reader := bufio.NewReader(strings.NewReader(source))
	l := &lexer{
		input: reader,
		line:  1,
	}

	l.consume()

	return l
}

func (l *lexer) consume() {
	r, _, err := l.input.ReadRune()
	if err != nil {
		if err == io.EOF {
			r = rune(token.EOF)
		} else {
			r = rune(token.ILLEGAL)
		}
	}

	l.lookAhead = r
}

func (l *lexer) match(r rune) error {
	if l.lookAhead == r {
		l.consume()
		return nil
	}

	return fmt.Errorf("expecting %s, found %s",
		strconv.QuoteRune(r), strconv.QuoteRune(l.lookAhead))
}

// Line returns the line number of current token.
func (l *lexer) Line() int {
	return l.line
}

// Read consumes and returns a token.
func (l *lexer) Read() token.Token {
	for l.lookAhead != rune(token.EOF) {
		switch l.lookAhead {
		case '#':
			l.readComment()
		case '\uFEFF', '\u0009', '\u0020', '\u000A', '\u000D', ',': // ignored
			if l.lookAhead == '\u000A' { // new line
				l.line++
			}
			l.consume()
			continue
		case '!':
			l.consume()
			return token.Token{Kind: token.BANG, Text: "!"}
		case '$':
			l.consume()
			return token.Token{Kind: token.DOLLAR, Text: "$"}
		case '(':
			l.consume()
			return token.Token{Kind: token.LPAREN, Text: "("}
		case ')':
			l.consume()
			return token.Token{Kind: token.RPAREN, Text: ")"}
		case ':':
			l.consume()
			return token.Token{Kind: token.COLON, Text: ":"}
		case '=':
			l.consume()
			return token.Token{Kind: token.EQL, Text: "="}
		case '@':
			l.consume()
			return token.Token{Kind: token.AT, Text: "@"}
		case '[':
			l.consume()
			return token.Token{Kind: token.LBRACK, Text: "["}
		case ']':
			l.consume()
			return token.Token{Kind: token.RBRACK, Text: "]"}
		case '{':
			l.consume()
			return token.Token{Kind: token.LBRACE, Text: "{"}
		case '|':
			l.consume()
			return token.Token{Kind: token.PIPE, Text: "|"}
		case '}':
			l.consume()
			return token.Token{Kind: token.RBRACE, Text: "}"}
		case '.': // ...
			return l.readSpread()
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-':
			return l.readNumber()
		case 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o',
			'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', 'A', 'B', 'C', 'D', 'E',
			'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U',
			'V', 'W', 'X', 'Y', 'Z', '_':
			return l.readName()
		case '"':
			return l.readString()
		default:
			return token.Token{Kind: token.ILLEGAL, Text: string(l.lookAhead)}
		}
	}
	return token.TokenEOF
}

func (l *lexer) readSpread() token.Token {
	// ...
	var b bytes.Buffer
	var err error
	for i := 0; i < 3; i++ {
		if err = l.match('.'); err != nil {
			b.WriteRune(l.lookAhead)
			return token.Token{Kind: token.ILLEGAL, Text: b.String()}
		}
		b.WriteRune('.')
	}

	return token.Token{Kind: token.SPREAD, Text: b.String()}
}

func (l *lexer) readComment() {
	l.consume()
	for l.lookAhead != rune(token.EOF) &&
		(l.lookAhead >= '\u0020' || l.lookAhead == '\u0009') { // SourceCharacter but not LineTerminator
		l.consume()
	}
}

// IntValue : '-'? IntegerPart
// FloatValue : '-'? IntegerPart ('.' Digit*)? ('e'|'E') '-'? Digit+
// IntegerPart : '0' | NonZeroDigit Digit*
func (l *lexer) readNumber() token.Token {
	var b bytes.Buffer
	if l.lookAhead == '-' {
		b.WriteRune('-')
		l.match('-')
	}

	// IntegerPart
	if l.lookAhead == '0' {
		b.WriteRune('0')
		l.match('0')
		if '0' <= l.lookAhead && l.lookAhead <= '9' {
			b.WriteRune(l.lookAhead)
			l.consume()
			return token.Token{Kind: token.ILLEGAL, Text: b.String()}
		}
	} else if '1' <= l.lookAhead && l.lookAhead <= '9' {
		b.WriteRune(l.lookAhead)
		l.consume()
		for '0' <= l.lookAhead && l.lookAhead <= '9' {
			b.WriteRune(l.lookAhead)
			l.consume()
		}
	} else {
		b.WriteRune(l.lookAhead)
		l.consume()
		return token.Token{Kind: token.ILLEGAL, Text: b.String()}
	}

	var isFloat bool
	if l.lookAhead == '.' {
		isFloat = true
		b.WriteRune('.')
		l.match('.')

		if '0' <= l.lookAhead && l.lookAhead <= '9' {
			b.WriteRune(l.lookAhead)
			l.consume()
		} else {
			b.WriteRune(l.lookAhead)
			l.consume()
			return token.Token{Kind: token.ILLEGAL, Text: b.String()}
		}

		for '0' <= l.lookAhead && l.lookAhead <= '9' {
			b.WriteRune(l.lookAhead)
			l.consume()
		}
	}

	if l.lookAhead == 'e' || l.lookAhead == 'E' {
		isFloat = true
		b.WriteRune(l.lookAhead)
		l.consume()

		if l.lookAhead == '-' || l.lookAhead == '+' {
			b.WriteRune(l.lookAhead)
			l.consume()
		}

		if '0' <= l.lookAhead && l.lookAhead <= '9' {
			b.WriteRune(l.lookAhead)
			l.consume()
		} else {
			b.WriteRune(l.lookAhead)
			l.consume()
			return token.Token{Kind: token.ILLEGAL, Text: b.String()}
		}

		for '0' <= l.lookAhead && l.lookAhead <= '9' {
			b.WriteRune(l.lookAhead)
			l.consume()
		}
	}

	if isFloat {
		return token.Token{Kind: token.FLOAT, Text: b.String()}
	}

	return token.Token{Kind: token.INT, Text: b.String()}
}

func (l *lexer) readName() token.Token {
	var b bytes.Buffer
	b.WriteRune(l.lookAhead)
	l.consume()
	for l.lookAhead == '_' ||
		('0' <= l.lookAhead && l.lookAhead <= '9') ||
		('a' <= l.lookAhead && l.lookAhead <= 'z') ||
		('A' <= l.lookAhead && l.lookAhead <= 'Z') {
		b.WriteRune(l.lookAhead)
		l.consume()
	}

	return token.Token{Kind: token.NAME, Text: b.String()}
}

// '"' ([\u0009\u0020-\uFFFF]|EscapedUnicode|EscapedChar)* '"'
// EscapedUnicode: \u [0-9A-Fa-f]{4}
// EscapedChar: \" \\ \/ \b \f \n \r \t
func (l *lexer) readString() token.Token {
	var b bytes.Buffer
	b.WriteRune('"')
	l.consume()

	for l.lookAhead != rune(token.EOF) && l.lookAhead != '"' && l.lookAhead != '\u000A' && l.lookAhead != '\u000D' {

		// SourceCharacter
		if l.lookAhead < '\u0020' && l.lookAhead != '\u0009' {
			b.WriteRune(l.lookAhead)
			l.consume()
			return token.Token{Kind: token.ILLEGAL, Text: b.String()}
		}

		if l.lookAhead == '\\' { // Escaped Char and Unicode
			l.consume()

			switch l.lookAhead {
			case '"':
				b.WriteRune('"')
				l.consume()
			case '\\':
				b.WriteRune('\\')
				l.consume()
			case '/':
				b.WriteRune('/')
				l.consume()
			case 'b':
				b.WriteRune('\b')
				l.consume()
			case 'f':
				b.WriteRune('\f')
				l.consume()
			case 'n':
				b.WriteRune('\n')
				l.consume()
			case 'r':
				b.WriteRune('\r')
				l.consume()
			case 't':
				b.WriteRune('\t')
				l.consume()
			case 'u':
				l.consume()

				hex1 := l.lookAhead
				l.consume()
				hex2 := l.lookAhead
				l.consume()
				hex3 := l.lookAhead
				l.consume()
				hex4 := l.lookAhead
				l.consume()

				quote := fmt.Sprintf(`'\u%s'`, string([]rune{hex1, hex2, hex3, hex4}))
				ucode, err := strconv.Unquote(quote)
				if err != nil {
					b.WriteString(quote)
					return token.Token{Kind: token.ILLEGAL, Text: b.String()}
				}
				b.WriteRune([]rune(ucode)[0])
			default:
				b.WriteRune(l.lookAhead)
				l.consume()
				return token.Token{Kind: token.ILLEGAL, Text: b.String()}
			}
		} else {
			b.WriteRune(l.lookAhead)
			l.consume()
		}
	}

	if l.lookAhead != '"' {
		b.WriteRune(l.lookAhead)
		l.consume()
		return token.Token{Kind: token.ILLEGAL, Text: b.String()}
	}

	b.WriteRune('"')
	l.consume()
	strVal := b.String()
	return token.Token{Kind: token.STRING, Text: strVal[1 : len(strVal)-1]}
}