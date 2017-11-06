package ql

import (
	"fmt"
	"strconv"
)

// Token is the set of lexical tokens of GraphQL.
type Token struct {
	Kind Kind
	Text string
}

// Kind represents the Token kind.
type Kind int

// TokenEOF is a special token for end-of-file.
var (
	TokenEOF = Token{EOF, tokens[EOF]}
)

// Token types defined
const (
	// Special tokens
	ILLEGAL Kind = iota
	EOF
	COMMENT

	// Punctuators
	punctBeg
	BANG     // !
	DOLLAR   // $
	LPAREN   // (
	RPAREN   // )
	ELLIPSIS // ...
	COLON    // :
	EQL      // =
	AT       // @
	LBRACK   // [
	RBRACK   // ]
	LBRACE   // {
	PIPE     // |
	RBRACE   // }
	punctEnd

	// Names and scalar type literals
	literalBeg
	NAME   // query
	INT    // 12345
	FLOAT  // 123.45
	STRING // "abc"
	literalEnd

	// Reserved keywords
	keywordBeg
	QUERY
	MUTATION
	FRAGMENT
	ON
	TYPE
	INTERFACE
	IMPLEMENTS
	UNION
	INPUT
	ENUM
	keywordEnd
)

var tokens = map[Kind]string{
	ILLEGAL:    "ILLEGAL",
	EOF:        "EOF",
	COMMENT:    "COMMENT",
	BANG:       "!",
	DOLLAR:     "$",
	LPAREN:     "(",
	RPAREN:     ")",
	ELLIPSIS:   "...",
	COLON:      ":",
	EQL:        "=",
	AT:         "@",
	LBRACK:     "[",
	RBRACK:     "]",
	LBRACE:     "{",
	PIPE:       "|",
	RBRACE:     "}",
	NAME:       "NAME",
	INT:        "INT",
	FLOAT:      "FLOAT",
	STRING:     "STRING",
	QUERY:      "query",
	MUTATION:   "mutation",
	FRAGMENT:   "fragment",
	ON:         "on",
	TYPE:       "type",
	INTERFACE:  "interface",
	IMPLEMENTS: "implements",
	UNION:      "union",
	INPUT:      "input",
	ENUM:       "enum",
}

var keywords map[string]Kind
var puncts map[string]Kind

func init() {
	keywords = make(map[string]Kind)
	for i := keywordBeg + 1; i < keywordEnd; i++ {
		keywords[tokens[i]] = i
	}
	puncts = make(map[string]Kind)
	for i := punctBeg + 1; i < punctEnd; i++ {
		puncts[tokens[i]] = i
	}
}

// Lookup maps an identifier to its keyword kind or NAME (if not a keyword).
func Lookup(ident string) Kind {
	if tok, isKeyword := keywords[ident]; isKeyword {
		return tok
	}
	return NAME
}

// IsPunct returns true for tokens corresponding to punctuators;
// it returns false otherwise.
func (tok Token) IsPunct() bool {
	return punctBeg < tok.Kind && tok.Kind < punctEnd
}

// IsLiteral returns true for tokens corresponding to names or scalar types;
// it returns false otherwise.
func (tok Token) IsLiteral() bool {
	return literalBeg < tok.Kind && tok.Kind < literalEnd
}

// IsKeyword returns true for tokens corresponding to keywords;
// it returns false otherwise.
func (tok Token) IsKeyword() bool {
	return keywordBeg < tok.Kind && tok.Kind < keywordEnd
}

// String returns the string corresponding to the token tok.
func (tok Token) String() string {
	s := ""
	if 0 <= tok.Kind && tok.Kind < Kind(len(tokens)) {
		s = tokens[tok.Kind]
	}
	if s == "" {
		s = fmt.Sprintf("token(%s)", strconv.Itoa(int(tok.Kind)))
	}
	return s
}
