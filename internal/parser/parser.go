package parser

import (
	"google.golang.org/protobuf/types/descriptorpb"

	"rogchap.com/protoparser/internal/scanner"
	"rogchap.com/protoparser/internal/token"
)

type parser struct {
	scanner scanner.Scanner

	tok token.Token // last read token
	lit string      // token literal

	// TODO: Deal with errors
}

func (p *parser) init(src []byte) {
	p.scanner.Init(src)
	p.next()
}

func (p *parser) next() {
	p.tok, p.lit = p.scanner.Scan()
}

func (p *parser) expect(tok token.Token) {
	if p.tok != tok {
		// TODO: deal with errors
	}
	p.next()
}

func (p *parser) parseSyntax() string {
	p.next()
	p.expect(token.ASSIGN)
	s = string(p.lit[1 : len(p.lit)-1]) // strip the quotes
	p.next()
	return s
}

func (p *parser) parseFile() *descriptorpb.FileDescriptorProto {

	// syntax must be the first non-empty, non-comment line of the file.
	// defaults to proto2 if not defined.
	syntax = "proto2"
	if p.tok == tok.SYNTAX {
		syntax = p.parseSyntax()
	}

	for p.tok != token.EOF {
		// TODO: parse rest of file
		p.next()
	}

	return &descriptorpb.FileDescriptorProto{
		Name:             nil,
		Package:          nil,
		Dependency:       nil,
		PublicDependency: nil,
		MessageType:      nil,
		EnumType:         nil,
		Service:          nil,
		Extension:        nil,
		Options:          nil,
		Syntax:           &syntax,
	}
}
