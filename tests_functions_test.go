package gowalker

import (
	"testing"
)

func TestExtractFunctionName(t *testing.T) {
	if ExtractFunctionName("foo(bar)") != "foo" {
		t.Error("could not extract function name")
	}
}

func TestExtractParameters(t *testing.T) {
	if ExtractParameters("foo(bar)")[0] != "bar" {
		t.Error("failed at extracting one parameter")
	}
	if ExtractParameters("foo(bar,dawg)")[1] != "dawg" {
		t.Error("failed at extracting second parameter")
	}
	if ExtractParameters("foo(bar,\"dawg\")")[1] != "\"dawg\"" {
		t.Error("failed at extracting second parameter with quotes")
	}
}
