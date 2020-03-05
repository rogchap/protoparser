// Package token defines the constants representing the lexical tokens of the Protocol Buffer Version 3 Language Spec
//

package token

import "strconv"

// Token is the set of lexical tokens for a proto file
type Token int

// The list of lexical tokens
const (
	// Special Tokens
	ILLEGAL Token = iota
	EOF
	COMMENT

	// Identifiers and basic literals
	IDENT  // MessageName
	INT    // 1234
	FLOAT  // 1234.12
	BOOL   // true
	STRING // "abc"

	ASSIGN // =

	LPAREN // (
	RPAREN // )
	LBRACE // {
	RBRACE // }
	LBRACK // [
	RBRACK // ]
	LANGLE // <
	RANGLE // >

	SEMICOLON // ;
	DOT       // .

	keyword_beg
	SYNTAX
	IMPORT
	WEAK
	PUBLIC
	PACKAGE
	OPTION
	REPEATED
	ONEOF
	MAP
	RESERVED
	TO
	MAX
	ENUM
	MESSAGE
	SERVICE
	RPC
	STREAM
	RETURNS
	keyword_end
)

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",
	COMMENT: "COMMENT",

	IDENT:  "IDENT",
	INT:    "INT",
	FLOAT:  "FLOAT",
	BOOL:   "BOOL",
	STRING: "STRING",

	ASSIGN: "=",

	LPAREN: "(",
	RPAREN: ")",
	LBRACE: "{",
	RBRACE: "}",
	LBRACK: "[",
	RBRACK: "]",
	LANGLE: "<",
	RANGLE: ">",

	SEMICOLON: ";",
	DOT:       ".",

	SYNTAX:   "syntax",
	IMPORT:   "import",
	WEAK:     "weak",
	PUBLIC:   "public",
	PACKAGE:  "package",
	OPTION:   "option",
	REPEATED: "repeated",
	ONEOF:    "oneof",
	RESERVED: "reserved",
	MAP:      "map",
	TO:       "to",  // only used as a reserved range so maybe should not be a keyword
	MAX:      "max", // only used as a reserved range so maybe should not be a keyword
	ENUM:     "enum",
	MESSAGE:  "message",
	SERVICE:  "service",
	RPC:      "rpc",
	STREAM:   "stream",
	RETURNS:  "returns",
}

// String returns the string corresponding to the token tok.
// For operators, delimiters, and keywords the string is the actual
// token character sequence (e.g., for the token IMPORT, the string is
// "import"). For all other tokens the string corresponds to the token
// constant name (e.g. for the token IDENT, the string is "IDENT").
func (tok Token) String() string {
	s := ""
	if 0 <= tok && tok < Token(len(tokens)) {
		s = tokens[tok]
	}
	// helps with debugging
	if s == "" {
		s = "token(" + strconv.Itoa(int(tok)) + ")"
	}
	return s
}

var keywords map[string]Token

func init() {
	keywords = make(map[string]Token)
	for i := keyword_beg + 1; i < keyword_end; i++ {
		keywords[tokens[i]] = i
	}
}

// Lookup maps an identifier to its keyword token or IDENT (if not a keyword).
func Lookup(ident string) Token {
	if tok, is_keyword := keywords[ident]; is_keyword {
		return tok
	}
	return IDENT
}
