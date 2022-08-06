package gowalker

import (
	"testing"
)

func TestWalk(t *testing.T) {
	if rx, _ := Walk("foo.double_foo", map[string]interface{}{"foo": map[string]interface{}{"double_foo": "bar"}}, nil); rx != "bar" {
		t.Error("basic map navigation failing")
	}
	if rx, _ := Walk("foo", map[string]interface{}{"foo": map[string]interface{}{"double_foo": "bar"}}, nil); rx.(map[string]interface{})["double_foo"] != "bar" {
		t.Error("addressing map as return value failed")
	}
	if rx, _ := Walk("foo[0]", map[string]interface{}{"foo": []interface{}{"bar1", "bar2"}}, nil); rx != "bar1" {
		t.Error("array navigation not working")
	}
	if _, err := Walk("foo[3]", map[string]interface{}{"foo": []interface{}{"bar1", "bar2"}}, nil); err == nil {
		t.Error("expression should index out of bounds")
	}
	if rx, _ := Walk("foo", map[string]interface{}{"foo": []interface{}{"bar1", "bar2"}}, nil); rx.([]interface{})[0] != "bar1" {
		t.Error("returning entire array not working")
	}
	if rx, _ := Walk("foo[0].gino", map[string]interface{}{"foo": []interface{}{map[string]interface{}{"gino": "pino"}, "bar2"}}, nil); rx != "pino" {
		t.Error("navigating in object past array not working")
	}
	if _, err := Walk("foo.bananas", map[string]interface{}{"foo": []interface{}{map[string]interface{}{"gino": "pino"}, "bar2"}}, nil); err == nil {
		t.Error("referencing an attribute in an array is not returning an error")
	}
	if rx, _ := Walk("foo[0][1]", map[string]interface{}{"foo": []interface{}{[]interface{}{"foo", "bar"}}}, nil); rx != "bar" {
		t.Error("nested array selection fails")
	}
	if rx, _ := Walk("foo[0][1].foo", map[string]interface{}{"foo": []interface{}{[]interface{}{"foo",
		map[string]interface{}{"foo": "bar"}}}}, nil); rx != "bar" {
		t.Error("nested array selection with more digging into sub-object failed")
	}
	if res, _ := Walk("foo.dawg.bar", map[string]interface{}{"foo": map[string]interface{}{"bar": "dawg"}}, nil); res != nil {
		t.Error("nil in the path broke something")
	}
}

func TestExtractIndex(t *testing.T) {
	if partial, index, _ := ExtractIndexes("foo[0]"); partial != "foo" || index[0] != 0 {
		t.Error("could not extract 1 digit index or partial")
	}
	if partial, index, _ := ExtractIndexes("foo[29]"); partial != "foo" || index[0] != 29 {
		t.Error("could not extract 2 digits index or partial")
	}
	if partial, index, _ := ExtractIndexes("foo[]"); partial != "foo[]" || index != nil {
		t.Error("error parsing empty square brackets")
	}
	if partial, index, _ := ExtractIndexes("foo"); partial != "foo" || index != nil {
		t.Error("could not extract no index partial")
	}
	if partial, _, _ := ExtractIndexes("foo[bar]"); partial != "foo[bar]" {
		t.Error("an index with alpha characters should be parsed as a segment")
	}
}

func TestWalkWithFunctions(t *testing.T) {

	f0 := func(data ...interface{}) (interface{}, error) {
		return "hello world", nil
	}
	f1 := func(data ...interface{}) (interface{}, error) {
		return "I'm crazy: " + data[0].(map[string]interface{})["double_foo"].(string), nil
	}
	functions := Functions{}
	functions.Add("hello", f0).Add("debug", f1)
	if res, _ := Walk("hello()", map[string]interface{}{"foo": map[string]interface{}{"double_foo": "bar"}}, functions); res != "hello world" {
		t.Error("basic function calling not working")
	}
	if res, _ := Walk("foo.debug(\"bananas\")", map[string]interface{}{"foo": map[string]interface{}{"double_foo": "bar"}}, functions); res != "I'm crazy: bar" {
		t.Error("reflexive function calling not working")
	}
	if res, _ := Walk("dawg()", map[string]interface{}{"foo": map[string]interface{}{"double_foo": "bar"}}, functions); res != "dawg()" {
		t.Error("calling a function that does not exist not working")
	}
}
