package gowalker

import (
	"context"
	"encoding/json"
	"testing"
)

func TestFromJSON(t *testing.T) {
	ctx := context.Background()
	var scope map[string]any
	_ = json.Unmarshal([]byte(`{
								"id":"banana",
								"meta": {
									"counter":11	
								},
								"items": [ "foo,bar","bar" ],
								"more_items": [
									{ "gino":22, "pino":10, "cane":5},
									{ "gino":22, "pino":10, "cane":5}
								]
							}`), &scope)
	if res, _ := Render(ctx, `{
	"name":"${id}",
	"availability": ${meta.counter},
	"first_item": "${items[0]}",
	"all_items": ${items},
	"item_count": ${items.size()},
	"something": ${items[0].split(\,)},
	"more_something": ${more_items.collect(pino,cane)}
}`, scope, nil); res != `{
	"name":"banana",
	"availability": 11,
	"first_item": "foo,bar",
	"all_items": ["foo,bar","bar"],
	"item_count": 2,
	"something": ["foo","bar"],
	"more_something": [{"cane":5,"pino":10},{"cane":5,"pino":10}]
}` {
		t.Error("template with data from JSON did not work")
	}
}
