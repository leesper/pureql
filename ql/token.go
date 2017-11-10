package ql

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
	COMMENT

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

	// Names and scalar type literals
	literalBeg
	NAME   // query
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
	EOF:          "EOF",
	COMMENT:      "COMMENT",
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

func illegalToken(v string) Token {
	return Token{ILLEGAL, v}
}

func lookupName(name string) Kind {
	kind, ok := keywords[name]
	if ok {
		return kind
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
	return fmt.Sprintf("<'%s', %s>", tok.Text, tokens[tok.Kind])
}
