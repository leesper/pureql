package parser

import "fmt"

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

	// Punctuators
	punctBeg
	BANG   // !
	DOLLAR // $
	LPAREN // (
	RPAREN // )
	SPREAD // ...
	COLON  // :
	EQL    // =
	AT     // @
	LBRACK // [
	RBRACK // ]
	LBRACE // {
	PIPE   // |
	RBRACE // }
	punctEnd

	NAME // query

	// scalar type literals
	literalBeg
	INT    // 12345
	FLOAT  // 123.45
	STRING // "abc"
	literalEnd

	// Reserved keywords
	keywordBeg
	SCHEMA
	QUERY
	MUTATION
	SUBSCRIPTION
	FRAGMENT
	ON
	SCALAR
	TYPE
	INTERFACE
	IMPLEMENTS
	EXTEND
	UNION
	INPUT
	ENUM
	DIRECTIVE
	keywordEnd
)

var tokens = map[Kind]string{
	ILLEGAL:      "ILLEGAL",
	EOF:          "<EOF>",
	BANG:         "!",
	DOLLAR:       "$",
	LPAREN:       "(",
	RPAREN:       ")",
	SPREAD:       "...",
	COLON:        ":",
	EQL:          "=",
	AT:           "@",
	LBRACK:       "[",
	RBRACK:       "]",
	LBRACE:       "{",
	PIPE:         "|",
	RBRACE:       "}",
	NAME:         "NAME",
	INT:          "INT",
	FLOAT:        "FLOAT",
	STRING:       "STRING",
	SCHEMA:       "schema",
	QUERY:        "query",
	MUTATION:     "mutation",
	SUBSCRIPTION: "subscription",
	FRAGMENT:     "fragment",
	ON:           "on",
	SCALAR:       "scalar",
	TYPE:         "type",
	INTERFACE:    "interface",
	IMPLEMENTS:   "implements",
	EXTEND:       "extend",
	UNION:        "union",
	INPUT:        "input",
	ENUM:         "enum",
	DIRECTIVE:    "directive",
}

// String returns the string corresponding to the token tok.
func (tok Token) String() string {
	return fmt.Sprintf("<'%s', %s>", tok.Text, tokens[tok.Kind])
}

// IsKeyword returns true if a token Kind is in keyword range.
func IsKeyword(k Kind) bool {
	return keywordBeg < k && k < keywordEnd
}

// Stringify returns the string representation of Kind.
func Stringify(k Kind) string {
	return tokens[k]
}
