package token_test

import (
	"testing"

	"rogchap.com/protoparser/internal/token"
)

func TestLookup(t *testing.T) {
	t.Parallel()
	var tests = [...]struct {
		given    string
		expected token.Token
	}{
		{"syntax", token.SYNTAX},
		{"import", token.IMPORT},
		{"weak", token.WEAK},
		{"public", token.PUBLIC},
		{"package", token.PACKAGE},
		{"option", token.OPTION},
		{"repeated", token.REPEATED},
		{"oneof", token.ONEOF},
		{"reserved", token.RESERVED},
		{"map", token.MAP},
		{"to", token.TO},
		{"enum", token.ENUM},
		{"message", token.MESSAGE},
		{"service", token.SERVICE},
		{"rpc", token.RPC},
		{"returns", token.RETURNS},
		{"stream", token.STREAM},
		{"my_message", token.IDENT},
		{"semicolon", token.IDENT},
		{"", token.IDENT},
		{" ", token.IDENT},
		{"123", token.IDENT},
		{"你好", token.IDENT},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.given, func(t *testing.T) {
			t.Parallel()
			actual := token.Lookup(tt.given)
			if actual != tt.expected {
				t.Errorf("(%s): expected %s, actual %s", tt.given, tt.expected, actual)
			}

		})
	}
}
