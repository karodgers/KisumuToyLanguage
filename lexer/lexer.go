package lexer

import (
	"unicode"
)

type TokenType string

const (
	TokenKeyword  TokenType = "KEYWORD"
	TokenOperator TokenType = "OPERATOR"
	TokenIdent    TokenType = "IDENT"
	TokenNumber   TokenType = "NUMBER"
	TokenString   TokenType = "STRING"
	TokenEOF      TokenType = "EOF"
	TokenError    TokenType = "ERROR"
)

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
	line         int
	column       int
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input, line: 1, column: 0}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
	l.column++
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) NextToken() Token {
	l.skipWhitespace()

	tok := Token{Line: l.line, Column: l.column}

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			tok.Literal = "=="
			tok.Type = TokenOperator
			l.readChar()
		} else {
			tok = l.newToken(TokenOperator, string(l.ch))
		}
	case '>':
		tok = l.newToken(TokenOperator, string(l.ch))
	case '<':
		tok = l.newToken(TokenOperator, string(l.ch))
	case '(':
		tok = l.newToken(TokenOperator, string(l.ch))
	case ')':
		tok = l.newToken(TokenOperator, string(l.ch))
	case '{':
		tok = l.newToken(TokenOperator, string(l.ch))
	case '}':
		tok = l.newToken(TokenOperator, string(l.ch))
	case '"':
		tok.Type = TokenString
		tok.Literal = l.readString()
	case 0:
		tok.Literal = ""
		tok.Type = TokenEOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = l.identifierType(tok.Literal)
			return tok
		} else if unicode.IsDigit(rune(l.ch)) {
			tok.Literal = l.readNumber()
			tok.Type = TokenNumber
			return tok
		} else {
			tok = l.newToken(TokenError, string(l.ch))
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || unicode.IsDigit(rune(l.ch)) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for unicode.IsDigit(rune(l.ch)) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		if l.ch == '\n' {
			l.line++
			l.column = 0
		}
		l.readChar()
	}
}

func (l *Lexer) newToken(tokenType TokenType, literal string) Token {
	return Token{Type: tokenType, Literal: literal, Line: l.line, Column: l.column}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *Lexer) identifierType(ident string) TokenType {
	keywords := map[string]bool{
		"declare": true, "displayln": true, "if": true, "case": true, "otherwise": true,
	}
	if keywords[ident] {
		return TokenKeyword
	}
	return TokenIdent
}
