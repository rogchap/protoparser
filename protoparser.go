// Package TODO
//
package protoparser

import (
	"google.golang.org/protobuf/types/descriptorpb"
	"rogchap.com/protoparser/internal/parser"
)

// TODO: Should we provide so wrapper functions to make this easier for the caller?
// for example:
// Parse(filename string)
// ParseReader(r io.Reader)
// ParseBuffer(buf bytes.Buffer)

// ParseFile TODO
func ParseFile(filename string, src interface{}) (*descriptorpb.FileDescriptorProto, error) {
	return parser.ParseFile(filename, src)
}
