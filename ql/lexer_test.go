package ql

import (
	"testing"
)

func TestLexerComment(t *testing.T) {
	comments := []string{
		`#this is comment`,
		`# this is comment with space`,
		`#	this is comment with tab`,
		`#		this is comment with tabs		`,
	}
	for _, comment := range comments {
		lexer := NewLexerWithSource(comment)
		if tok := lexer.Read(); tok != TokenEOF {
			t.Errorf("returned: %v, expected: %v", tok, TokenEOF)
		}
	}
}

func TestLexerPunct(t *testing.T)         {}
func TestLexerScalar(t *testing.T)        {}
func TestLexerKeyword(t *testing.T)       {}
func TestLexerIllegal(t *testing.T)       {}
func TestLexerQueryDocument(t *testing.T) {}
func TestLexerTypeDefs(t *testing.T)      {}
