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
	if extractParameters("foo(bar.dawg)")[0] != "bar.dawg" {
		t.Error("failed to extract parameter with dot")
	}
}

func TestSizeFunction(t *testing.T) {
	if res, _ := Walk("size()", map[string]string{"foo": "bar", "foo2": "bar2"}, NewFunctions()); res != 2 {
		t.Error("size on root object not working")
	}
	if res, _ := Walk("foo.size()", map[string]any{"foo": []int{1, 2, 3}}, NewFunctions()); res != 3 {
		t.Error("size on array not working")
	}
	if _, err := Walk("foo.size()", map[string]any{"foo": 22}, NewFunctions()); err == nil {
		t.Error("error not reported for unsupported size")
	}
	if _, err := Walk("foo.size()", map[string]any{"foo": nil}, NewFunctions()); err == nil {
		t.Error("error not reported for nil")
	}
}

func TestSplitFunction(t *testing.T) {
	if res, _ := Walk("foo.split(|)", map[string]string{"foo": "bar|bananas"}, NewFunctions()); res.([]string)[0] != "bar" {
		t.Error("split with pipe string not working")
	}
}
