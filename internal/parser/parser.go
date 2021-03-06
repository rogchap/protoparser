package parser

import (
	"strconv"
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

func (p *parser) parseBoolLit() bool {
	b, ok := strconv.ParseBool(p.lit)
	if !ok {
		//TODO: deal with error
		return false
	}
	return b
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

func (p *parser) parserFieldType() (typ descriptorpb.FieldDescriptorProto_Type, name string) {
	key := "TYPE_" + strings.ToUpper(p.lit)
	typ = descriptorpb.FieldDescriptorProto_Type(descriptorpb.FieldDescriptorProto_Type_value[key])
	if typ == 0 {
		// TODO: Need to get type name for Message or Enum
		// will need some sort of lookup
	}
	p.next()
	return
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

	typ, typName = p.parserFieldType()

	if p.tok != token.IDENT {
		//TODO: deal with unexpected token
	}
	name = p.lit
	jsonName = jsonCamelCase(name)
	p.next()

	p.expect(token.ASSIGN)

	if p.tok != token.INT {
		// TODO: deal with error
	}

	i, _ := strconv.Atoi(p.lit)
	number = int32(i)
	p.next()

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

func (p *parser) parseMapField() (*descriptorpb.FieldDescriptorProto, *descriptorpb.DescriptorProto) {
	// mapField = "map" "<" keyType "," type ">" mapName "=" fieldNumber [ "[" fieldOptions "]" ] ";"
	// keyType = "int32" | "int64" | "uint32" | "uint64" | "sint32" | "sint64" |
	//           "fixed32" | "fixed64" | "sfixed32" | "sfixed64" | "bool" | "string"
	p.next()

	var (
		name, entryName, typName string
		lbl                      = descriptorpb.FieldDescriptorProto_LABEL_REPEATED
		typ                      = descriptorpb.FieldDescriptorProto_TYPE_MESSAGE
		entryFields              []*descriptorpb.FieldDescriptorProto
		entryOpts                = &descriptorpb.MessageOptions{MapEntry: boolPtr(true)}
	)

	p.expect(token.LANGLE)
	ktyp, _ := p.parserFieldType()
	if ktyp == 0 {
		// TODO: deal with unexpected key type
	}
	p.expect(token.COMMA)
	vtyp, vtypName := p.parserFieldType()
	p.expect(token.RANGLE)

	entryLbl := descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL
	var n, n2 int32 = 1, 2
	entryFields = append(entryFields, &descriptorpb.FieldDescriptorProto{
		Name:     strPtr("key"),
		JsonName: strPtr("key"),
		Label:    &entryLbl,
		Number:   &n,
		Type:     &ktyp,
	}, &descriptorpb.FieldDescriptorProto{
		Name:     strPtr("value"),
		JsonName: strPtr("value"),
		Label:    &entryLbl,
		Number:   &n2,
		Type:     &vtyp,
		TypeName: strPtr(vtypName),
	})

	if p.tok != token.IDENT {
		//TODO: deal with error
	}
	name = p.lit
	entryName = mapEntryName(p.lit)
	p.next()

	//TODO: Get the full type name for the MapEntry message that we are adding

	p.expect(token.ASSIGN)

	if p.tok != token.INT {
		//TODO: deal with unexpected token
	}
	i, _ := strconv.Atoi(p.lit)
	number := int32(i)

	// TODO: deal with options
	p._skipTo(token.SEMICOLON)

	return &descriptorpb.FieldDescriptorProto{
			Name:     strPtr(name),
			JsonName: strPtr(jsonCamelCase(name)),
			Label:    &lbl,
			Number:   &number,
			Type:     &typ,
			TypeName: strPtr(typName),
		},
		&descriptorpb.DescriptorProto{
			Name:    strPtr(entryName),
			Field:   entryFields,
			Options: entryOpts,
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
			mf, mn := p.parseMapField()
			fields = append(fields, mf)
			nested = append(nested, mn)
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

func (p *parser) parseEnumValue() *descriptorpb.EnumValueDescriptorProto {
	name := p.lit
	p.next()
	p.expect(token.ASSIGN)
	if p.tok != token.INT {
		// TODO: deal with error
	}
	i, _ := strconv.Atoi(p.lit)
	number := int32(i)

	//TODO: parse options
	p._skipTo(token.SEMICOLON)

	return &descriptorpb.EnumValueDescriptorProto{
		Name:   strPtr(name),
		Number: &number,
	}
}

func (p *parser) parseEnum() *descriptorpb.EnumDescriptorProto {
	// enum = "enum" enumName enumBody
	// enumBody = "{" { option | enumField | emptyStatement } "}"
	// enumField = ident "=" intLit [ "[" enumValueOption { ","  enumValueOption } "]" ]";"
	// enumValueOption = optionName "=" constant
	p.next()

	var (
		name string
		vals []*descriptorpb.EnumValueDescriptorProto
		opts *descriptorpb.EnumOptions
	)

	if p.tok != token.IDENT {
		//TODO deal with unexpected token
	}
	name = p.lit
	p.next()
	p.expect(token.LBRACE)

	for p.tok != token.RBRACE && p.tok != token.EOF {
		switch p.tok {
		case token.OPTION:
			// TODO: parse options
			p._skipTo(token.SEMICOLON)
		case token.SEMICOLON:
			p.next()
		case token.IDENT:
			//TODO: parse vals
			vals = append(vals, p.parseEnumValue())
		default:
			// TODO: deal with unexpected token
			p.next()
		}
	}

	p._skipTo(token.RBRACE)

	return &descriptorpb.EnumDescriptorProto{
		Name:    strPtr(name),
		Value:   vals,
		Options: opts,
	}
}

func (p *parser) parseFileOption(opt *descriptorpb.FileOptions) {
	p.next()
	if p.tok != token.IDENT {
		//TODO: deal with errors
	}
	ident := p.lit
	// TODO: deal with fullIdent
	p.next()

	p.expect(token.ASSIGN)

	switch token.LookupFileOption(ident) {
	case token.JAVA_PACKAGE:
		opt.JavaPackage = strPtr(p.parseStrLit())
	case token.JAVA_OUTER_CLASSNAME:
		opt.JavaOuterClassname = strPtr(p.parseStrLit())
	case token.JAVA_MULTIPLE_FILES:
		opt.JavaMultipleFiles = boolPtr(p.parseBoolLit())
	case token.JAVA_STRING_CHECK_UTF8:
		opt.JavaStringCheckUtf8 = boolPtr(p.parseBoolLit())
	case token.OPTIMIZE_FOR:
		// TODO: parse enum values
	case token.GO_PACKAGE:
		opt.GoPackage = strPtr(p.parseStrLit())
	case token.CC_GENERIC_SERVICES:
		opt.CcGenericServices = boolPtr(p.parseBoolLit())
	case token.JAVA_GENERIC_SERVICES:
		opt.JavaGenericServices = boolPtr(p.parseBoolLit())
	case token.PY_GENERIC_SERVICES:
		opt.PyGenericServices = boolPtr(p.parseBoolLit())
	case token.PHP_GENERIC_SERVICES:
		opt.PhpGenericServices = boolPtr(p.parseBoolLit())
	case token.DEPRECATED:
		opt.Deprecated = boolPtr(p.parseBoolLit())
	case token.CC_ENABLE_ARENAS:
		opt.CcEnableArenas = boolPtr(p.parseBoolLit())
	case token.OBJC_CLASS_PREFIX:
		opt.ObjcClassPrefix = strPtr(p.parseStrLit())
	case token.CSHARP_NAMESPACE:
		opt.CsharpNamespace = strPtr(p.parseStrLit())
	case token.SWIFT_PREFIX:
		opt.SwiftPrefix = strPtr(p.parseStrLit())
	case token.PHP_CLASS_PREFIX:
		opt.PhpClassPrefix = strPtr(p.parseStrLit())
	case token.PHP_NAMESPACE:
		opt.PhpClassPrefix = strPtr(p.parseStrLit())
	case token.PHP_METADATA_NAMESPACE:
		opt.PhpMetadataNamespace = strPtr(p.parseStrLit())
	case token.RUBY_PACKAGE:
		opt.RubyPackage = strPtr(p.parseStrLit())
	default:
		//TODO: deal with custom option
	}

	p._skipTo(token.SEMICOLON)
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
		case token.OPTION:
			if opt == nil {
				opt = &descriptorpb.FileOptions{}
			}
			p.parseFileOption(opt)
		case token.MESSAGE:
			msgs = append(msgs, p.parseMessage())
		case token.ENUM:
			enums = append(enums, p.parseEnum())
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

func boolPtr(b bool) *bool {
	return &b
}
