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

func (p *parser) parseStrLit() string {
	if p.tok != token.STRING {
		//TODO: deal with error
		return ""
	}
	return string(p.lit[1 : len(p.lit)-1]) // strip the quotes
}

func (p *parser) parseSyntax() string {
	p.next()
	p.expect(token.ASSIGN)
	s := p.parseStrLit()
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

func (p *parser) parseDependency() (dep string, isPublic, isWeak bool) {
	p.next()
	isPublic = p.tok == token.PUBLIC
	isWeak = p.tok == token.WEAK
	if isPublic || isWeak {
		p.next()
	}
	dep = p.parseStrLit()
	p.expect(token.SEMICOLON)
	return
}

//TODO Remove once finished
func (p *parser) _skipTo(tok token.Token) {
	for p.tok != tok && p.tok != token.EOF {
		p.next()
	}
	p.next()
}

func (p *parser) parseNormalField() *descriptorpb.FieldDescriptorProto {
	var (
		name     string
		number   int32
		label    descriptorpb.FieldDescriptorProto_Label
		typ      descriptorpb.FieldDescriptorProto_Type
		typName  string
		jsonName string
		opt      *descriptorpb.FieldOptions
	)
	switch p.tok {
	case token.REPEATED:
		label = descriptorpb.FieldDescriptorProto_LABEL_REPEATED
		p.next()
	default:
		label = descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL
	}

	// TODO: deal with normal, message, enum type
	p.next()

	if p.tok != token.IDENT {
		//TODO: deal with unexpected token
	}
	name = p.lit
	jsonName = jsonCamelCase(name)

	p.expect(token.ASSIGN)

	p._skipTo(token.SEMICOLON)
	return &descriptorpb.FieldDescriptorProto{
		Name:     strPtr(name),
		Number:   &number,
		Label:    &label,
		Type:     &typ,
		TypeName: strPtr(typName),
		JsonName: strPtr(jsonName),
		Options:  opt,
	}
}

func (p *parser) parseMessage() *descriptorpb.DescriptorProto {
	// message = "message" messageName messageBody
	// messageBody = "{" { field | enum | message | option | oneof | mapField |
	// reserved | emptyStatement } "}"
	p.next()

	var (
		name    string
		fields  []*descriptorpb.FieldDescriptorProto
		exts    []*descriptorpb.FieldDescriptorProto
		nested  []*descriptorpb.DescriptorProto
		enums   []*descriptorpb.EnumDescriptorProto
		extRng  []*descriptorpb.DescriptorProto_ExtensionRange
		oneofs  []*descriptorpb.OneofDescriptorProto
		opt     *descriptorpb.MessageOptions
		resRng  []*descriptorpb.DescriptorProto_ReservedRange
		resName []string
	)

	if p.tok != token.IDENT {
		// TODO deal with error
	}
	name = p.lit
	p.expect(token.LBRACE)
	for p.tok != token.RBRACE && p.tok != token.EOF {
		switch p.tok {
		case token.OPTION:
			//TODO parse options
			p._skipTo(token.SEMICOLON)
		case token.MESSAGE:
			nested = append(nested, p.parseMessage())
		case token.REPEATED, token.IDENT:
			fields = append(fields, p.parseNormalField())
		case token.MAP:
			p._skipTo(token.SEMICOLON)
		default:
			// TODO: deal with unexpected token
			p.next()
		}
	}
	p.expect(token.RBRACE)

	return &descriptorpb.DescriptorProto{
		Name:           strPtr(name),
		Field:          fields,
		Extension:      exts,
		NestedType:     nested,
		EnumType:       enums,
		ExtensionRange: extRng,
		OneofDecl:      oneofs,
		Options:        opt,
		ReservedRange:  resRng,
		ReservedName:   resName,
	}
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
		wDeps     []int32
		msgs      []*descriptorpb.DescriptorProto
		enums     []*descriptorpb.EnumDescriptorProto
		srcs      []*descriptorpb.ServiceDescriptorProto
		exts      []*descriptorpb.FieldDescriptorProto
		opt       *descriptorpb.FileOptions
	)

	for p.tok != token.EOF {
		switch p.tok {
		case token.PACKAGE:
			// [RC] if pacakge is declared twice is this a syntax err
			// or do we expect the first or last pacakge declaration?
			pkg = p.parsePackage()
		case token.IMPORT:
			d, p, w := p.parseDependency()
			deps = append(deps, d)
			if p {
				pDeps = append(pDeps, int32(len(deps)-1))
			}
			if w {
				wDeps = append(wDeps, int32(len(deps)-1))
			}
		case token.MESSAGE:
			println("Outer Message")
			msgs = append(msgs, p.parseMessage())
		default:
			// TODO: deal with unexpected token error
			p.next()
		}

	}

	return &descriptorpb.FileDescriptorProto{
		Name:             strPtr(name),
		Package:          strPtr(pkg),
		Dependency:       deps,
		PublicDependency: pDeps,
		WeakDependency:   wDeps,
		MessageType:      msgs,
		EnumType:         enums,
		Service:          srcs,
		Extension:        exts,
		Options:          opt,
		Syntax:           strPtr(syntax),
	}
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
