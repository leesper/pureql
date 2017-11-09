package ql

import (
	"bufio"
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// Lexer converts the GraphQL source text into tokens.
type Lexer struct {
	scanner *bufio.Scanner
	current rune
}

// NewLexerWithSource returns a new Lexer parsing source.
func NewLexerWithSource(source string) *Lexer {
	scanner := bufio.NewScanner(strings.NewReader(source))
	return newLexerWithScanner(scanner)
}

func newLexerWithScanner(scanner *bufio.Scanner) *Lexer {
	scanner.Split(bufio.ScanRunes)
	lexer := &Lexer{
		scanner: scanner,
	}
	lexer.consume()
	return lexer
}

// Read consumes and returns a token.
func (l *Lexer) Read() Token {
	for l.current != rune(EOF) {
		switch l.current {
		case '#':
			l.readComment()
		case '\uFEFF', '\u0009', '\u0020', '\u000A', '\u000D', ',': // ignored
			l.consume()
			continue
		case '!':
			l.consume()
			return Token{BANG, "!"}
		case '$':
			l.consume()
			return Token{DOLLAR, "$"}
		case '(':
			l.consume()
			return Token{LPAREN, "("}
		case ')':
			l.consume()
			return Token{RPAREN, ")"}
		case ':':
			l.consume()
			return Token{COLON, ":"}
		case '=':
			l.consume()
			return Token{EQL, "="}
		case '@':
			l.consume()
			return Token{AT, "@"}
		case '[':
			l.consume()
			return Token{LBRACK, "["}
		case ']':
			l.consume()
			return Token{RBRACK, "]"}
		case '{':
			l.consume()
			return Token{LBRACE, "{"}
		case '|':
			l.consume()
			return Token{PIPE, "|"}
		case '}':
			l.consume()
			return Token{RBRACE, "}"}
		case '.': // ...
			l.consume()
			if l.current != '.' {
				return illegalToken(string(l.current))
			}
			l.consume()
			if l.current != '.' {
				return illegalToken(string(l.current))
			}
			l.consume()
			return Token{SPREAD, "..."}
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '-':
			return l.readNumber()
		case 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o',
			'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', 'A', 'B', 'C', 'D', 'E',
			'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U',
			'V', 'W', 'X', 'Y', 'Z', '_':
			return l.readName()
		case '"':
			return l.readString()
		}
	}
	return TokenEOF
}

func (l *Lexer) consume() {
	if l.scanner.Scan() {
		l.current = []rune(l.scanner.Text())[0]
	} else {
		l.current = rune(EOF)
	}
}

func (l *Lexer) readComment() {
	l.consume()
	for l.current != rune(EOF) &&
		(l.current > '\u001F' || l.current == '\u0009') { // SourceCharacter but not LineTerminator
		l.consume()
	}
}

// IntValue : '-'? IntegerPart
// FloatValue : '-'? IntegerPart ('.' Digit*)? ('e'|'E') '-'? Digit+
// IntegerPart : '0' | NonZeroDigit Digit*
func (l *Lexer) readNumber() Token {
	var b bytes.Buffer
	if l.current == '-' {
		b.WriteRune('-')
		l.consume()
	}

	// IntegerPart
	if l.current == '0' {
		b.WriteRune('0')
		l.consume()
	} else if '1' <= l.current && l.current <= '9' {
		b.WriteRune(l.current)
		l.consume()
		for '0' <= l.current && l.current <= '9' {
			b.WriteRune(l.current)
			l.consume()
		}
	} else {
		b.WriteRune(l.current)
		l.consume()
		return illegalToken(b.String())
	}

	var isFloat bool
	if l.current == '.' {
		isFloat = true
		b.WriteRune('.')
		l.consume()
		for '0' <= l.current && l.current <= '9' {
			b.WriteRune(l.current)
			l.consume()
		}
	}

	if l.current == 'e' || l.current == 'E' {
		isFloat = true
		b.WriteRune(l.current)
		l.consume()
	} else {
		b.WriteRune(l.current)
		l.consume()
		return illegalToken(b.String())
	}

	if l.current == '-' {
		b.WriteRune('-')
		l.consume()
	}

	if '0' <= l.current && l.current <= '9' {
		b.WriteRune(l.current)
		l.consume()
	} else {
		b.WriteRune(l.current)
		l.consume()
		return illegalToken(b.String())
	}

	for '0' <= l.current && l.current <= '9' {
		b.WriteRune(l.current)
		l.consume()
	}

	if isFloat {
		return Token{FLOAT, b.String()}
	}

	return Token{INT, b.String()}
}

// TODO
func (l *Lexer) readName() Token {
	return TokenEOF
}

// '"' ([\u0009\u0020-\uFFFF]|EscapedUnicode|EscapedChar)* '"'
// EscapedUnicode: \u [0-9A-Fa-f]{4}
// EscapedChar: \" \\ \/ \b \f \n \r \t
func (l *Lexer) readString() Token {
	var b bytes.Buffer
	b.WriteRune('"')
	l.consume()

	for l.current != rune(EOF) && l.current != '"' && l.current != '\u000A' && l.current != '\u000D' {

		// SourceCharacter
		if l.current == '\u0009' || ('\u0020' <= l.current && l.current <= '\uFFFF') {
			b.WriteRune(l.current)
			l.consume()
		} else if l.current == '\\' { // Escaped Char and Unicode
			l.consume()

			switch l.current {
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

				hex1 := l.current
				l.consume()
				hex2 := l.current
				l.consume()
				hex3 := l.current
				l.consume()
				hex4 := l.current
				l.consume()

				quote := fmt.Sprintf(`'\u%s'`, string([]rune{hex1, hex2, hex3, hex4}))
				ucode, err := strconv.Unquote(quote)
				if err != nil {
					b.WriteString(quote)
					return illegalToken(b.String())
				}
				b.WriteRune([]rune(ucode)[0])
			default:
				b.WriteRune(l.current)
				l.consume()
				return illegalToken(b.String())
			}
		}
	}

	if l.current != '"' {
		b.WriteRune(l.current)
		l.consume()
		return illegalToken(b.String())
	}

	b.WriteRune('"')
	strVal := b.String()
	return Token{STRING, strVal[1 : len(strVal)-1]}
}

// // Peek returns a token in ith position from current.
// func (l *Lexer) Peek(i int) Token {
// 	tok := TokenEOF
// 	if l.more(i) {
// 		tok = l.tokenBuf[i]
// 	}
// 	return tok
// }
//
// func (l *Lexer) more(i int) bool {
// 	for i >= len(l.tokenBuf) {
// 		if l.hasMore {
// 			l.readLine()
// 		} else {
// 			return false
// 		}
// 	}
// 	return true
// }
//
// func (l *Lexer) readLine() {
// 	l.hasMore = l.scanner.Scan()
// 	if l.hasMore {
// 		line := l.scanner.Text()
// 		tokens := l.tokenize(line)
// 		l.tokenBuf = append(l.tokenBuf, tokens...)
// 	}
// }
//
// func (l *Lexer) tokenize(line string) []Token {
// 	var tokens []Token
// 	return tokens
// }
