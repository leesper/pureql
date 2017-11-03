package ql

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var (
	regexComment = `(?P<comment>#.*)`
	regexPunct   = `(?P<punct>!|\$|\(|\)|\.{3}|:|=|@|\[|\]|\{|\}|\|)`
	regexName    = `(?P<name>[_A-Za-z][_0-9A-Za-z]*)`
	regexInt     = `(?P<int>-?(?:0|[1-9][0-9]*))`
	regexFloat   = `(?P<float>-?(?:0|[1-9][0-9]*)(?:\.[0-9]+)?(?:(?:E|e)(?:\+|-)?[0-9]+)?)`
	regexString  = `(?P<string>\"(?:[^"\\\n\r]|(?:\\(?:u[0-9a-fA-F]{4}|["\\/bfnrt])))*\")`
	regexPat     = fmt.Sprintf(`\s*%s|%s|%s|%s|%s|%s`,
		regexComment,
		regexFloat,
		regexInt,
		regexPunct,
		regexName,
		regexString)
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
	tok := EOF
	if l.more(0) {
		tok = l.tokenBuf[0]
		l.tokenBuf = l.tokenBuf[1:]
	}
	return tok
}

// Peek returns a token in ith position from current.
func (l *Lexer) Peek(i int) Token {
	tok := EOF
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
		// TODO: parse line into tokens, append them in tokenBuf
		fmt.Println("line", line)
	}
	// example:
	// fmt.Println(regexPat)
	// re := regexp.MustCompile(regexPat)
	// allMatches := re.FindAllStringSubmatch(`computer = "net"`, -1)
	// names := re.SubexpNames()
	// for _, match := range allMatches {
	// 	group := map[string]string{}
	// 	for i := 0; i < len(names); i++ {
	// 		if names[i] != "" && match[i] != "" {
	// 			group[names[i]] = match[i]
	// 		}
	// 	}
	// 	fmt.Printf("group: %#v\n", group)
	// }
}
