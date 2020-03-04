// Package scanner implements a lexical scanner for a proto file
//
package scanner

import "rogchap.com/protoparser/internal/token"

// Scanner is the data structure for a lexer
type Scanner struct {
	src []byte

	// scanning state
	ch       rune // current char
	offset   int  // char offset
	rdOffset int  // reading offset (position after the current char)
}

const bom = 0xFEFF // byte order mark, only permitted as very first character

// Init initiates a Scanner
func (s *Scanner) Init(src []byte) {
	s.src = src
	s.ch = ' '
	s.offset = 0
	s.rdOffset = 0

	s.next()
	if s.ch == bom {
		s.next()
	}
}

// read the next Unicode from the source
// < 0 means end-of-file.
func (s *Scanner) next() {
	if s.rdOffset < len(s.src) {
		s.offset = s.rdOffset
		s.ch = rune(s.src[s.rdOffset])
		s.rdOffset += 1
	} else {
		s.offset = len(s.src)
		s.ch = -1 // eof
	}
}

func (s *Scanner) skipWhitespace() {
	for s.ch == ' ' || s.ch == '\t' || s.ch == '\n' || s.ch == '\r' {
		s.next()
	}
}

// Scan will scan the next rune and consume any literals
func (s *Scanner) Scan() (tok token.Token, lit string) {
	s.skipWhitespace()

	switch ch := s.ch; {
	case isLetter(ch):
		lit = s.scanIdentifier()
		tok = token.Lookup(lit)
	case isDigit(ch):
	default:
		s.next() // always make progress
		switch ch {
		case '=':
			tok = token.ASSIGN
		case '"', '\'':
			tok = token.STRING
			lit = s.scanString()
		case -1:
			tok = token.EOF
		default:
			tok = token.ILLEGAL
			lit = string(ch)
		}
	}
	return
}

func isLetter(ch rune) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z'
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

func (s *Scanner) scanIdentifier() string {
	offs := s.offset
	for isLetter(s.ch) || isDigit(s.ch) || s.ch == '_' {
		s.next()
	}
	return string(s.src[offs:s.offset])
}

func (s *Scanner) scanString() string {
	offs := s.offset - 1 // opening quote already consumed
	quote := rune(s.src[offs])
	for {
		ch := s.ch
		if ch == '\n' || ch < 0 {
			// TODO: deal with errors: string literal not terminated
			break
		}
		s.next()
		if ch == quote {
			break
		}
	}
	return string(s.src[offs:s.offset])
}
