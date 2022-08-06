package gowalker

import (
	"testing"
)

func TestExtractFunctionName(t *testing.T) {
	if extractFunctionName("foo(bar)") != "foo" {
		t.Error("could not extract function name")
	}
}

func TestExtractParameters(t *testing.T) {
	if extractParameters("foo(bar)")[0] != "bar" {
		t.Error("failed at extracting one parameter")
	}
	if extractParameters("foo(bar,dawg)")[1] != "dawg" {
		t.Error("failed at extracting second parameter")
	}
	if extractParameters("foo(bar,\"dawg\")")[1] != "\"dawg\"" {
		t.Error("failed at extracting second parameter with quotes")
	}
}
