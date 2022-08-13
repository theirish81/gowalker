package gowalker

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
)

func TestFromJSON(t *testing.T) {
	ctx := context.Background()
	var scope map[string]any
	_ = json.Unmarshal([]byte(`{
								"id":"banana",
								"meta": {
									"counter":11,
									"price": 2.99,
									"available": true
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
	"available": ${meta.available},
	"price": ${meta.price},
	"first_item": "${items[0]}",
	"all_items": ${items},
	"item_count": ${items.size()},
	"something": ${items[0].split(\,)},
	"more_something": ${more_items.collect(pino,cane)}
}`, scope, nil); res != `{
	"name":"banana",
	"availability": 11,
	"available": true,
	"price": 2.99,
	"first_item": "foo,bar",
	"all_items": ["foo,bar","bar"],
	"item_count": 2,
	"something": ["foo","bar"],
	"more_something": [{"cane":5,"pino":10},{"cane":5,"pino":10}]
}` {
		fmt.Println(res)
		t.Error("template with data from JSON did not work")
	}
}
