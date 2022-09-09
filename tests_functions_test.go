package gowalker

import (
	"context"
	"reflect"
	"testing"
)

func TestRunFunction(t *testing.T) {
	ctx := context.Background()
	if found, data, err := runFunction(ctx, "split(\\,)", "foo,bar", NewFunctions()); !found || err != nil || data.([]string)[0] != "foo" {
		t.Error("something went wrong while running function")
	}
	if found, _, _ := runFunction(ctx, "split(\\,", "foo,bar", NewFunctions()); found {
		t.Error("a function was found where there was none")
	}
	if found, _, _ := runFunction(ctx, "()", "foo,bar", NewFunctions()); found {
		t.Error("a function was found where there was none")
	}
}

func TestPassDotToFunction(t *testing.T) {
	ctx := context.Background()
	if res, _ := Walk(ctx, "foo.split(.)", map[string]string{"foo": "bar.bar"}, NewFunctions()); len(res.([]string)) != 2 {
		t.Error("passing a dot to a function breaks stuff")
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
	ctx := context.Background()
	if res, _ := Walk(ctx, "size()", map[string]string{"foo": "bar", "foo2": "bar2"}, NewFunctions()); res != 2 {
		t.Error("size on root object not working")
	}
	if res, _ := Walk(ctx, "foo.size()", map[string]any{"foo": []int{1, 2, 3}}, NewFunctions()); res != 3 {
		t.Error("size on array not working")
	}
	if _, err := Walk(ctx, "foo.size()", map[string]any{"foo": 22}, NewFunctions()); err == nil {
		t.Error("error not reported for unsupported size")
	}
	if _, err := Walk(ctx, "foo.size()", map[string]any{"foo": nil}, NewFunctions()); err == nil {
		t.Error("error not reported for nil")
	}
}

func TestSplitFunction(t *testing.T) {
	ctx := context.Background()
	if res, _ := Walk(ctx, "foo.split(|)", map[string]string{"foo": "bar|bananas"}, NewFunctions()); res.([]string)[0] != "bar" {
		t.Error("split with pipe string not working")
	}

	if _, err := Walk(ctx, "foo.split(|)", map[string]int{"foo": 22}, NewFunctions()); err == nil {
		t.Error("split on non-string does not return an error")
	}
}

func TestCollectFunction(t *testing.T) {
	ctx := context.Background()
	res, _ := Walk(ctx, "foo.collect(foo,bar)", map[string][]map[string]int{"foo": {{"foo": 1, "bar": 2, "gino": 3}, {"foo": 4, "bar": 5, "gino": 6}}}, NewFunctions())
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
	res, _ = Walk(ctx, "foo.collect(foo,bar)", map[string][]map[string]int{"foo": {{"foo": 1, "bar": 2, "gino": 3}, {"bar": 5, "gino": 6}}}, NewFunctions())
	rx = res.([]map[string]any)
	if len(rx[0]) != 2 || len(rx[1]) != 1 {
		t.Error("not the exact number of attributes on missing key")
	}

	if _, err := Walk(ctx, "foo.collect(foo,bar)", map[string]string{"foo": "bar"}, NewFunctions()); err == nil {
		t.Error("did not return error in collect when selected element is not an array")
	}
	if _, err := Walk(ctx, "foo.collect(foo,bar)", map[string][]string{"foo": {"bar"}}, NewFunctions()); err == nil {
		t.Error("did not return error in collect when child elements are not maps")
	}
}

func TestToVar(t *testing.T) {
	ctx := context.Background()
	fx := NewFunctions()
	fx.GetScope()["foo"] = map[string]string{"dawg": "bar"}
	if res, _ := Walk(ctx, "toVar(foo.dawg)", map[string]any{}, fx); res != "bar" {
		t.Error("renderVar not working")
	}
	fx.GetScope()["foo"] = map[string][]string{"dawg": {"bar", "yay"}}
	if res, _ := Walk(ctx, "toVar(foo.dawg[1])", map[string]any{}, fx); res != "yay" {
		t.Error("renderVar not working")
	}
}

func TestToString(t *testing.T) {
	data, _ := Walk(context.Background(), "foo.bar.toString()", map[string]map[string]any{"foo": {"bar": 2}}, nil)
	if reflect.ValueOf(data).Kind().String() != "string" {
		t.Error("toString is not working as expected")
	}
}

func TestEq(t *testing.T) {
	data, _ := Walk(context.Background(), "foo.bar.eq(22)", map[string]map[string]any{"foo": {"bar": 22}}, nil)
	if data == false {
		t.Error("positive equality did not work")
	}
	data, _ = Walk(context.Background(), "foo.bar.eq(22)", map[string]map[string]any{"foo": {"bar": true}}, nil)
	if data == true {
		t.Error("negative equality did not work")
	}
	data, _ = Walk(context.Background(), "foo.eq(22)", map[string]map[string]any{"foo": {"bar": true}}, nil)
	if data == true {
		t.Error("negative equality did not work")
	}
	if _, err := Walk(context.Background(), "foo.eq()", map[string]string{"foo": "bar"}, nil); err == nil {
		t.Error("empty parameter should return an error")
	}
}

func TestJsonEscape(t *testing.T) {
	f := NewFunctions()
	if res, _ := f.jsonEscape(context.Background(), "foo\""); res != "foo\\\"" {
		t.Error("could not json escape a string")
	}
	if _, err := f.jsonEscape(context.Background(), 22); err == nil {
		t.Error("json escape of a non string should return an error")
	}
}
