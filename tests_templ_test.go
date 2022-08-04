package gowalker

import (
	"testing"
)

func TestRender(t *testing.T) {
	data := map[string]interface{}{"name": "pino", "age": 22}
	templ := "my name is: ${name}, my age is ${age}"
	if res, _ := Render(templ, data); res != "my name is: pino, my age is 22" {
		t.Error("basic template not working")
	}

	data = map[string]interface{}{"name": "pino", "items": []interface{}{"keys", "wallet"}}
	templ = `{
	"name": "${name}",
    "first_item": "${items[0]}",
	"all_items": ${items}
}`
	if res, _ := Render(templ, data); res != `{
	"name": "pino",
    "first_item": "keys",
	"all_items": ["keys","wallet"]
}` {
		t.Error("array navigation not working")
	}
	data = map[string]interface{}{"user": map[string]interface{}{"name": "pino", "age": 22, "items": []interface{}{"keys", "wallet"}}}
	templ = `{
	"data": ${user}
}`
	if res, _ := Render(templ, data); res != `{
	"data": {"age":22,"items":["keys","wallet"],"name":"pino"}
}` {
		t.Error("printing maps does not work")
	}
}
