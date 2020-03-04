package protoparser_test

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"

	"rogchap.com/protoparser"
)

func TestParseFile(t *testing.T) {
	filename := filepath.Join("testdata", "test.proto")
	pb, err := protoparser.ParseFile(filename, nil)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	var actual, expected interface{}

	raw, _ := json.MarshalIndent(pb, "", " ")
	json.Unmarshal(raw, &actual)

	goldfile := filepath.Join("testdata", "proto.golden")
	graw, _ := ioutil.ReadFile(goldfile)
	json.Unmarshal(graw, &expected)

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("ParseFile() mismatch (-want +got):\n%s", diff)
	}
}
