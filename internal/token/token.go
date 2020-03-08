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
	COMMA     // ,

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

	file_options_beg
	JAVA_PACKAGE
	JAVA_OUTER_CLASSNAME
	JAVA_MULTIPLE_FILES
	JAVA_STRING_CHECK_UTF8
	OPTIMIZE_FOR
	GO_PACKAGE
	CC_GENERIC_SERVICES
	JAVA_GENERIC_SERVICES
	PY_GENERIC_SERVICES
	PHP_GENERIC_SERVICES
	DEPRECATED
	CC_ENABLE_ARENAS
	OBJC_CLASS_PREFIX
	CSHARP_NAMESPACE
	SWIFT_PREFIX
	PHP_CLASS_PREFIX
	PHP_NAMESPACE
	PHP_METADATA_NAMESPACE
	RUBY_PACKAGE
	file_options_end
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
	COMMA:     ",",

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

	JAVA_PACKAGE:           "java_package",
	JAVA_OUTER_CLASSNAME:   "java_outer_classname",
	JAVA_MULTIPLE_FILES:    "java_multiple_files",
	JAVA_STRING_CHECK_UTF8: "java_string_check_utf8",
	OPTIMIZE_FOR:           "optimize_for",
	GO_PACKAGE:             "go_package",
	CC_GENERIC_SERVICES:    "cc_generic_services",
	JAVA_GENERIC_SERVICES:  "java_generic_services",
	PY_GENERIC_SERVICES:    "py_generic_services",
	PHP_GENERIC_SERVICES:   "php_generic_services",
	DEPRECATED:             "deprecated",
	CC_ENABLE_ARENAS:       "cc_enable_arenas",
	OBJC_CLASS_PREFIX:      "objc_class_prefix",
	CSHARP_NAMESPACE:       "csharp_namespace",
	SWIFT_PREFIX:           "swift_prefix",
	PHP_CLASS_PREFIX:       "php_class_prefix",
	PHP_NAMESPACE:          "php_namespace",
	PHP_METADATA_NAMESPACE: "php_metadata_namespace",
	RUBY_PACKAGE:           "ruby_package",
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
var fileOpts map[string]Token

func init() {
	keywords = make(map[string]Token)
	for i := keyword_beg + 1; i < keyword_end; i++ {
		keywords[tokens[i]] = i
	}
	fileOpts = make(map[string]Token)
	for i := file_options_beg + 1; i < file_options_end; i++ {
		fileOpts[tokens[i]] = i
	}
}

// Lookup maps an identifier to its keyword token or IDENT (if not a keyword).
func Lookup(ident string) Token {
	if tok, is_keyword := keywords[ident]; is_keyword {
		return tok
	}
	return IDENT
}

func LookupFileOption(ident string) Token {
	if tok, ok := fileOpts[ident]; ok {
		return tok
	}
	return IDENT
}
