package ql

import (
	"bufio"
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
			return Token{ELLIPSIS, "..."}
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
	return illegalToken(string(l.current))
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

// TODO
func (l *Lexer) readNumber() Token {
	return TokenEOF
}

// TODO
func (l *Lexer) readName() Token {
	return TokenEOF
}

// TODO
func (l *Lexer) readString() Token {
	return TokenEOF
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
