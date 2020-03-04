# Proto Parser - a Protocol Buffer proto file parser for Go

**NOTE: This is a new project that is WIP, and is not ready for production use**

## Motivation

There are a few parsers that are already available:

* https://github.com/emicklei/proto
* https://github.com/yoheimuta/go-protoparser

However, both these parsers have there own defined AST and APIs for managing proto files. The aim of protoparser is to
provide a parser that uses the standard [`descriptorpb.FileDescriptorProto`](https://pkg.go.dev/google.golang.org/protobuf/types/descriptorpb?tab=doc#FileDescriptorProto) as the AST and provide a very basic API.

These makes it easier to use the standard protocol buffer tooling for Go (especially APIv2)
For example using the `protodesc` package for converting `FileDescriptorProto` messages to
`protoreflect.FileDescriptor` values. See
[google.golang.org/protobuf](https://pkg.go.dev/google.golang.org/protobuf?tab=overview) for details.

## Design Spec

The parser should follow the [Protocol Buffers Version 3 Language
Specification](https://developers.google.com/protocol-buffers/docs/reference/proto3-spec)

## Usage

```go
package main

import (
    "fmt"

    "rogchap.com/protoparser"
)

func main() {
    src := `
    syntax = "proto3";
    package foo.bar;

    message Foo {
        string name = 1;
    }
 
    message Bar {
        string id = 1;
        Foo foo = 2;
    }
    `

    // Parse the source
    fd, err := protoparser.ParseFile("", src)
    if err != nil {
        fmt.Println(err)
        return
    }

    // Print the all the messages and their fields
    for _, m := fd.MessageType {
        fmt.Println(m.Name)
        for _, f := m.Field {
            fmt.Println(" - "+f.Name)
        }
        fmt.Println()
    }
}

