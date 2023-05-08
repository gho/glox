package main

import (
	"unicode"
	"unicode/utf8"
)

type TokenType byte

const (
	TokenError TokenType = iota
	TokenEOF
	TokenLeftParen
	TokenRightParen
	TokenLeftBrace
	TokenRightBrace
	TokenComma
	TokenDot
	TokenPlus
	TokenMinus
	TokenStar
	TokenSlash
	TokenEqual
	TokenEqualEqual
	TokenBang
	TokenBangEqual
	TokenLess
	TokenLessEqual
	TokenGreater
	TokenGreaterEqual
	TokenSemicolon
	TokenString
	TokenNumber
	TokenIdentifier
	TokenAnd
	TokenClass
	TokenElse
	TokenFalse
	TokenFor
	TokenFun
	TokenIf
	TokenNil
	TokenOr
	TokenPrint
	TokenReturn
	TokenSuper
	TokenTrue
	TokenVar
	TokenWhile
)

type Scanner interface {
	nextToken() Token
}

type scanner struct {
	source  string
	start   int
	current int
	line    int
}

type Token struct {
	typ  TokenType
	line int
	data string
}

func newScanner(source string) Scanner {
	return &scanner{source: source}
}

func (s *scanner) nextToken() Token {
	s.skipWhitespace()
	s.start = s.current

	if s.isEOF() {
		return s.makeToken(TokenEOF)
	}

	r, size := s.currentRune()
	s.current += size

	if isDigit(r) {
		return s.number()
	}

	if isAlpha(r) {
		return s.identifier()
	}

	switch r {
	case '(':
		return s.makeToken(TokenLeftParen)
	case ')':
		return s.makeToken(TokenRightParen)
	case '{':
		return s.makeToken(TokenLeftBrace)
	case '}':
		return s.makeToken(TokenRightBrace)
	case ',':
		return s.makeToken(TokenComma)
	case '.':
		return s.makeToken(TokenDot)
	case '+':
		return s.makeToken(TokenPlus)
	case '-':
		return s.makeToken(TokenMinus)
	case '*':
		return s.makeToken(TokenStar)
	case '/':
		return s.makeToken(TokenSlash)
	case '=':
		if s.match('=') {
			return s.makeToken(TokenEqualEqual)
		} else {
			return s.makeToken(TokenEqual)
		}
	case '!':
		if s.match('=') {
			return s.makeToken(TokenBangEqual)
		} else {
			return s.makeToken(TokenBang)
		}
	case '<':
		if s.match('=') {
			return s.makeToken(TokenLessEqual)
		} else {
			return s.makeToken(TokenLess)
		}
	case '>':
		if s.match('=') {
			return s.makeToken(TokenGreaterEqual)
		} else {
			return s.makeToken(TokenGreater)
		}
	case ';':
		return s.makeToken(TokenSemicolon)
	case '"':
		return s.string()
	}

	return s.makeToken(TokenError)
}

func (s *scanner) string() Token {
	for {
		r, size := s.currentRune()
		if r == '"' || s.isEOF() {
			break
		}

		/* TODO multi-line string
		if r == '\n' {
			s.line++
		}
		*/

		s.current += size
	}

	if s.isEOF() {
		// unterminated string
		return s.makeToken(TokenError)
	}

	// closing quote
	_, size := s.currentRune()
	s.current += size

	return s.makeToken(TokenString)
}

func isDigit(r rune) bool {
	return unicode.IsDigit(r)
}

func (s *scanner) number() Token {
	for {
		r, size := s.currentRune()
		if isDigit(r) {
			s.current += size
			continue
		} else if r == '.' {
			s.current += size
			r, size := s.currentRune()
			for isDigit(r) {
				s.current += size
				r, size = s.currentRune()
			}
		}
		break
	}
	return s.makeToken(TokenNumber)
}

func isAlpha(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}

func (s *scanner) identifier() Token {
	r, size := s.currentRune()
	for isAlpha(r) || isDigit(r) {
		s.current += size
		r, size = s.currentRune()
	}

	token := s.makeToken(TokenIdentifier)

	switch token.data {
	case "and":
		token.typ = TokenAnd
	case "class":
		token.typ = TokenClass
	case "else":
		token.typ = TokenElse
	case "false":
		token.typ = TokenFalse
	case "for":
		token.typ = TokenFor
	case "fun":
		token.typ = TokenFun
	case "if":
		token.typ = TokenIf
	case "nil":
		token.typ = TokenNil
	case "or":
		token.typ = TokenOr
	case "print":
		token.typ = TokenPrint
	case "return":
		token.typ = TokenReturn
	case "super":
		token.typ = TokenSuper
	case "true":
		token.typ = TokenTrue
	case "var":
		token.typ = TokenVar
	case "while":
		token.typ = TokenWhile
	}

	return token
}

func (s *scanner) makeToken(typ TokenType) Token {
	return Token{
		typ:  typ,
		line: s.line + 1,
		data: s.source[s.start:s.current],
	}
}

func (s *scanner) isEOF() bool {
	return s.current >= len(s.source)
}

func (s *scanner) match(expected rune) bool {
	if s.isEOF() {
		return false
	}

	r, size := s.currentRune()
	if r != expected {
		return false
	}

	s.current += size

	return true
}

func (s *scanner) skipWhitespace() {
	for {
		r, size := s.currentRune()
		switch r {
		case ' ', '\r', '\t':
			s.current += size
			continue
		case '\n':
			s.line++
			s.current += size
			continue
		case '/':
			if n, _ := s.runeAt(s.current + size); n == '/' {
				s.skipUntilNewLine()
			}
		}
		break
	}
}

func (s *scanner) skipUntilNewLine() {
	r, size := s.currentRune()
	for r != '\n' && !s.isEOF() {
		s.current += size
		r, size = s.currentRune()
	}
}

func (s *scanner) currentRune() (rune, int) {
	return s.runeAt(s.current)
}

func (s *scanner) runeAt(index int) (rune, int) {
	if index >= utf8.RuneCountInString(s.source) {
		return -1, 0
	}
	return utf8.DecodeRuneInString(s.source[index:])
}

//go:generate stringer -type=TokenType
