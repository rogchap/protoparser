package parser

import (
	"strings"

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
	s := string(p.lit[1 : len(p.lit)-1]) // strip the quotes
	p.expect(token.SEMICOLON)
	return s
}

func (p *parser) parseFullIdent() string {
	var sb strings.Builder
	if p.tok != token.IDENT {
		//TODO: deal with error
		return ""
	}
	sb.WriteString(p.lit)
	prevTok := p.tok
	p.next()

	for p.tok == token.IDENT || p.tok == token.DOT {
		if p.tok == prevTok {
			//TODO: deal with error: invalid package name
			return ""
		}
		sb.WriteString(p.lit)
		prevTok = p.tok
		p.next()
	}

	return sb.String()
}

func (p *parser) parsePackage() string {
	p.next()
	s := p.parseFullIdent()
	p.expect(token.SEMICOLON)
	return s
}

func (p *parser) parseFile() *descriptorpb.FileDescriptorProto {

	// syntax must be the first non-empty, non-comment line of the file.
	// defaults to proto2 if not defined.
	syntax := "proto2"
	if p.tok == token.SYNTAX {
		syntax = p.parseSyntax()
		// TODO: validatate proto2 or proto3
	}

	var (
		name, pkg string
		deps      []string
		pDeps     []int32
		msgs      []*descriptorpb.DescriptorProto
		enums     []*descriptorpb.EnumDescriptorProto
		srcs      []*descriptorpb.ServiceDescriptorProto
		exts      []*descriptorpb.FieldDescriptorProto
		opt       *descriptorpb.FileOptions
	)

	for p.tok != token.EOF {
		switch p.tok {
		case token.PACKAGE:
			// [RC] if pacakge is declared twice is this a systax err
			// or do we expect the first or last pacakge declaration?
			pkg = p.parsePackage()
		default:
			// TODO: deal with unexpected token error
			p.next()
		}

	}

	return &descriptorpb.FileDescriptorProto{
		Name:             &name,
		Package:          &pkg,
		Dependency:       deps,
		PublicDependency: pDeps,
		MessageType:      msgs,
		EnumType:         enums,
		Service:          srcs,
		Extension:        exts,
		Options:          opt,
		Syntax:           &syntax,
	}
}
