package gowalker

import (
	"testing"
)

func TestWalk(t *testing.T) {
	if rx, _ := Walk("foo.double_foo", map[string]interface{}{"foo": map[string]interface{}{"double_foo": "bar"}}); rx != "bar" {
		t.Error("basic map navigation failing")
	}
	if rx, _ := Walk("foo", map[string]interface{}{"foo": map[string]interface{}{"double_foo": "bar"}}); rx.(map[string]interface{})["double_foo"] != "bar" {
		t.Error("addressing map as return value failed")
	}
	if rx, _ := Walk("foo[0]", map[string]interface{}{"foo": []interface{}{"bar1", "bar2"}}); rx != "bar1" {
		t.Error("array navigation not working")
	}
	if _, err := Walk("foo[3]", map[string]interface{}{"foo": []interface{}{"bar1", "bar2"}}); err == nil {
		t.Error("expression should index out of bounds")
	}
	if rx, _ := Walk("foo", map[string]interface{}{"foo": []interface{}{"bar1", "bar2"}}); rx.([]interface{})[0] != "bar1" {
		t.Error("returning entire array not working")
	}
	if rx, _ := Walk("foo[0].gino", map[string]interface{}{"foo": []interface{}{map[string]interface{}{"gino": "pino"}, "bar2"}}); rx != "pino" {
		t.Error("navigating in object past array not working")
	}
	if _, err := Walk("foo.bananas", map[string]interface{}{"foo": []interface{}{map[string]interface{}{"gino": "pino"}, "bar2"}}); err == nil {
		t.Error("referencing an attribute in an array is not returning an error")
	}
}

func TestExtractIndex(t *testing.T) {
	if partial, index := ExtractIndex("foo[0]"); partial != "foo" || index != 0 {
		t.Error("could not extract 1 digit index or partial")
	}
	if partial, index := ExtractIndex("foo[29]"); partial != "foo" || index != 29 {
		t.Error("could not extract 2 digits index or partial")
	}
	if partial, index := ExtractIndex("foo"); partial != "foo" || index != -1 {
		t.Error("could not extract no index partial")
	}
}
