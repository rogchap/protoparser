package protoparser

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"

	"google.golang.org/protobuf/types/descriptorpb"
)

func readSource(filename string, src interface{}) ([]byte, error) {
	if src != nil {
		switch s := src.(type) {
		case string:
			return []byte(s), nil
		case []byte:
			return s, nil
		case *bytes.Buffer:
			if s != nil {
				return s.Bytes(), nil
			}
		case io.Reader:
			return ioutil.ReadAll(s)
		}
		return nil, errors.New("protoparser: invalid source")
	}
	return ioutil.ReadFile(filename)
}

func ParseFile(filename string, src interface{}) (*descriptorpb.FileDescriptorProto, error) {
	source, err := readSource(filename, src)
	if err != nil {
		return nil, err
	}

	// TODO: parse the source
	_ = source

	return &descriptorpb.FileDescriptorProto{}, nil
}
