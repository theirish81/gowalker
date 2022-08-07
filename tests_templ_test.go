package gowalker

import (
	"testing"
)

func TestRender(t *testing.T) {

	data := map[string]any{"name": "pino", "age": 22}
	templ := "my name is: ${name}, my age is ${age}"
	if res, _ := Render(templ, data, nil); res != "my name is: pino, my age is 22" {
		t.Error("basic template not working")
	}

	data = map[string]any{"name": "pino", "items": []any{"keys", "wallet"}}
	templ = `{
	"name": "${name}",
    "first_item": "${items[0]}",
	"all_items": ${items}
}`
	if res, _ := Render(templ, data, nil); res != `{
	"name": "pino",
    "first_item": "keys",
	"all_items": ["keys","wallet"]
}` {
		t.Error("array navigation not working")
	}
	data = map[string]any{"user": map[string]any{"name": "pino", "age": 22, "items": []any{"keys", "wallet"}}}
	templ = `{
	"data": ${user}
}`
	if res, _ := Render(templ, data, nil); res != `{
	"data": {"age":22,"items":["keys","wallet"],"name":"pino"}
}` {
		t.Error("printing maps does not work")
	}
	if res, _ := Render("foo bar", map[string]any{}, nil); res != "foo bar" {
		t.Error("something went wrong when no template tags are present")
	}
	if res, _ := Render("foo bar", nil, nil); res != "foo bar" {
		t.Error("something went wrong when scope is nil")
	}
	if res, _ := Render("${foo}", map[string]any{"bar": "bar"}, nil); res != "${foo}" {
		t.Error("something went wrong while rendering a template referencing a missing variable")
	}
}

func TestRenderWithFunctions(t *testing.T) {
	functions := NewFunctions()
	functions.Add("hello", func(data any, params ...string) (any, error) {
		return "hello world", nil
	})
	functions.Add("first", func(data any, params ...string) (any, error) {
		return data.([]any)[0], nil
	})
	if res, _ := Render("What do we all say? ${hello()}", map[string]any{}, functions); res != "What do we all say? hello world" {
		t.Error("simple function in template not working")
	}

	if res, _ := Render("First element in the array is: ${myArray.first()}", map[string]any{"myArray": []any{0, 1, 2, 3}}, functions); res != "First element in the array is: 0" {
		t.Error("reflexive function not working")
	}

	if res, _ := Render("Splitting and printing ${foo.split(\\,)}", map[string]interface{}{"foo": "bar,dawg"}, functions); res != "Splitting and printing [\"bar\",\"dawg\"]" {
		t.Error("error in running split function in template")
	}

}
