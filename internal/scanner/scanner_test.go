package scanner_test

import (
	"testing"

	"rogchap.com/protoparser/internal/scanner"
	"rogchap.com/protoparser/internal/token"
)

type el struct {
	tok token.Token
	lit string
}

var tokens = [...]el{
	{token.IDENT, "foobar"},
}

const whitespace = "  \t  \n\n\n" // to separate tokens

var source = func() []byte {
	var src []byte
	for _, t := range tokens {
		src = append(src, t.lit...)
		src = append(src, whitespace...)
	}
	return src
}

func TestScan(t *testing.T) {

	var s scanner.Scanner
	s.Init(source())

	for _, e := range tokens {
		tok, lit := s.Scan()

		// check token
		if tok != e.tok {
			t.Errorf("bad token for %q: got %s, expected %s", lit, tok, e.tok)
		}
	}
}
