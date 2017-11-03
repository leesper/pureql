package ql

import (
	"fmt"
	"strconv"
)

// Token is the set of lexical tokens of GraphQL
type Token int

// Token types defined
const (
	// Special tokens
	ILLEGAL Token = iota
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

var tokens = map[Token]string{
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

var keywords map[string]Token

func init() {
	keywords = make(map[string]Token)
	for i := keywordBeg + 1; i < keywordEnd; i++ {
		keywords[tokens[i]] = i
	}
}

// Lookup maps an identifier to its keyword token or NAME (if not a keyword).
func Lookup(ident string) Token {
	if tok, isKeyword := keywords[ident]; isKeyword {
		return tok
	}
	return NAME
}

// IsPunct returns true for tokens corresponding to punctuators;
// it returns false otherwise.
func (tok Token) IsPunct() bool {
	return punctBeg < tok && tok < punctEnd
}

// IsLiteral returns true for tokens corresponding to names or scalar types;
// it returns false otherwise.
func (tok Token) IsLiteral() bool {
	return literalBeg < tok && tok < literalEnd
}

// IsKeyword returns true for tokens corresponding to keywords;
// it returns false otherwise.
func (tok Token) IsKeyword() bool {
	return keywordBeg < tok && tok < keywordEnd
}

// String returns the string corresponding to the token tok.
func (tok Token) String() string {
	s := ""
	if 0 <= tok && tok < Token(len(tokens)) {
		s = tokens[tok]
	}
	if s == "" {
		s = fmt.Sprintf("token(%s)", strconv.Itoa(int(tok)))
	}
	return s
}
