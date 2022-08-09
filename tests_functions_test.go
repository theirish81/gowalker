package gowalker

import (
	"testing"
)

func TestRunFunction(t *testing.T) {
	if found, data, err := runFunction("split(\\,)", "foo,bar", NewFunctions()); !found || err != nil || data.([]string)[0] != "foo" {
		t.Error("something went wrong while running function")
	}
	if found, _, _ := runFunction("split(\\,", "foo,bar", NewFunctions()); found {
		t.Error("a function was found where there was none")
	}
	if found, _, _ := runFunction("()", "foo,bar", NewFunctions()); found {
		t.Error("a function was found where there was none")
	}
}

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

func TestCollectFunction(t *testing.T) {
	res, _ := Walk("foo.collect(foo,bar)", map[string][]map[string]int{"foo": {{"foo": 1, "bar": 2, "gino": 3}, {"foo": 4, "bar": 5, "gino": 6}}}, NewFunctions())
	rx := res.([]map[string]any)
	if len(rx) != 2 {
		t.Error("did not return all items in array")
	}
	if len(rx[0]) != 2 || len(rx[1]) != 2 {
		t.Error("did not select all items in the array")
	}
	if rx[0]["foo"] != 1 || rx[1]["foo"] != 4 {
		t.Error("did not collect the right values")
	}
	res, _ = Walk("foo.collect(foo,bar)", map[string][]map[string]int{"foo": {{"foo": 1, "bar": 2, "gino": 3}, {"bar": 5, "gino": 6}}}, NewFunctions())
	rx = res.([]map[string]any)
	if len(rx[0]) != 2 || len(rx[1]) != 1 {
		t.Error("not the exact number of attributes on missing key")
	}

}
