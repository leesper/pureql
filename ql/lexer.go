package ql

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// TODO: using regular expression is somehow a slow method, maybe using ANTLR or
// doing character-by-character lexer instead.
var (
	regexComment = `(?P<comment>#[\x{0009}\x{0020}-\x{FFFF}]*)`
	regexPunct   = `(?P<punct>!|\$|\(|\)|\.{3}|:|=|@|\[|\]|\{|\}|\|)`
	regexName    = `(?P<name>[_A-Za-z][_0-9A-Za-z]*)`
	regexNumeric = `(?P<numeric>-?(?:0|[1-9][0-9]*)(?:\.[0-9]+)?(?:(?:E|e)(?:\+|-)?[0-9]+)?)`
	regexString  = `(?P<string>\"(?:[^"\\\x{000A}\x{000D}]|(?:\\(?:u[0-9a-fA-F]{4}|["\\/bfnrt])))*\")`
	regexPat     = fmt.Sprintf(`\s*%s|%s|%s|%s|%s`,
		regexComment,
		regexNumeric,
		regexPunct,
		regexName,
		regexString)
	re     = regexp.MustCompile(regexPat)
	groups = re.SubexpNames()[1:] // skip the whole pattern group
)

// Lexer converts the GraphQL source text into tokens.
type Lexer struct {
	tokenBuf []Token
	scanner  *bufio.Scanner
	hasMore  bool
}

// NewLexerWithFile returns a new Lexer parsing from file.
func NewLexerWithFile(file *os.File) *Lexer {
	return &Lexer{
		scanner: bufio.NewScanner(file),
		hasMore: true,
	}
}

// NewLexerWithSource returns a new Lexer parsing source.
func NewLexerWithSource(source string) *Lexer {
	return &Lexer{
		scanner: bufio.NewScanner(strings.NewReader(source)),
		hasMore: true,
	}
}

// Read consumes and returns a token.
func (l *Lexer) Read() Token {
	tok := TokenEOF
	if l.more(0) {
		tok = l.tokenBuf[0]
		l.tokenBuf = l.tokenBuf[1:]
	}
	return tok
}

// Peek returns a token in ith position from current.
func (l *Lexer) Peek(i int) Token {
	tok := TokenEOF
	if l.more(i) {
		tok = l.tokenBuf[i]
	}
	return tok
}

func (l *Lexer) more(i int) bool {
	for i >= len(l.tokenBuf) {
		if l.hasMore {
			l.readLine()
		} else {
			return false
		}
	}
	return true
}

func (l *Lexer) readLine() {
	l.hasMore = l.scanner.Scan()
	if l.hasMore {
		line := l.scanner.Text()
		tokens := l.tokenize(line)
		l.tokenBuf = append(l.tokenBuf, tokens...)
	}
}

func (l *Lexer) tokenize(line string) []Token {
	var tokens []Token
	matches := re.FindAllStringSubmatch(line, -1)
	for _, match := range matches {
		match = match[1:] // skip the whole pattern group
		for idx, val := range match {
			if val != "" {
				group := groups[idx]
				switch group {
				case "numeric":
					if _, err := strconv.Atoi(val); err != nil {
						tokens = append(tokens, Token{FLOAT, val})
					} else {
						tokens = append(tokens, Token{INT, val})
					}
				case "comment":
					// ignore, do nothing
				case "punct":
					tokens = append(tokens, Token{puncts[val], val})
				case "name":
					tokens = append(tokens, Token{Lookup(val), val})
				case "string":
					tokens = append(tokens, Token{STRING, strings.Trim(val, `"`)})
				}
			}
		}
	}
	return tokens
}
