package gowalker

import (
	"context"
	"testing"
	"time"
)

func TestWalk(t *testing.T) {
	ctx := context.Background()
	if rx, _ := Walk(ctx, "foo.double_foo", map[string]any{"foo": map[string]any{"double_foo": "bar"}}, nil); rx != "bar" {
		t.Error("basic map navigation failing")
	}
	if rx, _ := Walk(ctx, "foo", map[string]any{"foo": map[string]any{"double_foo": "bar"}}, nil); rx.(map[string]any)["double_foo"] != "bar" {
		t.Error("addressing map as return value failed")
	}
	if rx, _ := Walk(ctx, "foo[0]", map[string]any{"foo": []any{"bar1", "bar2"}}, nil); rx != "bar1" {
		t.Error("array navigation not working")
	}
	if _, err := Walk(ctx, "foo[3]", map[string]any{"foo": []any{"bar1", "bar2"}}, nil); err == nil {
		t.Error("expression should index out of bounds")
	}
	if rx, _ := Walk(ctx, "foo", map[string]any{"foo": []any{"bar1", "bar2"}}, nil); rx.([]any)[0] != "bar1" {
		t.Error("returning entire array not working")
	}
	if rx, _ := Walk(ctx, "foo[0].gino", map[string]any{"foo": []any{map[string]any{"gino": "pino"}, "bar2"}}, nil); rx != "pino" {
		t.Error("navigating in object past array not working")
	}
	if _, err := Walk(ctx, "foo.bananas", map[string]any{"foo": []any{map[string]any{"gino": "pino"}, "bar2"}}, nil); err == nil {
		t.Error("referencing an attribute in an array is not returning an error")
	}
	if rx, _ := Walk(ctx, "foo[0][1]", map[string]any{"foo": []any{[]any{"foo", "bar"}}}, nil); rx != "bar" {
		t.Error("nested array selection fails")
	}
	if rx, _ := Walk(ctx, "foo[0][1].foo", map[string]any{"foo": []any{[]any{"foo",
		map[string]any{"foo": "bar"}}}}, nil); rx != "bar" {
		t.Error("nested array selection with more digging into sub-object failed")
	}
	if res, _ := Walk(ctx, "foo.dawg.bar", map[string]any{"foo": map[string]any{"bar": "dawg"}}, nil); res != nil {
		t.Error("nil in the path broke something")
	}
	if res, _ := Walk(ctx, "foo.bar", map[string]any{"foo": nil}, nil); res != nil {
		t.Error("nil as legit value not working")
	}
	if res, _ := Walk(ctx, "", map[string]string{"foo": "bar"}, nil); res.(map[string]string)["foo"] != "bar" {
		t.Error("empty selector not working")
	}
}

func TestWalkWithNestedTypes(t *testing.T) {
	ctx := context.Background()
	if rx, _ := Walk(ctx, "foo.double_foo", map[string]map[string]string{"foo": {"double_foo": "bar"}}, nil); rx != "bar" {
		t.Error("basic map navigation failing")
	}
}

func TestExtractIndex(t *testing.T) {
	if partial, index := extractIndexes("foo[0]"); partial != "foo" || index[0] != 0 {
		t.Error("could not extract 1 digit index or partial")
	}
	if partial, index := extractIndexes("foo[29]"); partial != "foo" || index[0] != 29 {
		t.Error("could not extract 2 digits index or partial")
	}
	if partial, index := extractIndexes("foo[]"); partial != "foo[]" || index != nil {
		t.Error("error parsing empty square brackets")
	}
	if partial, index := extractIndexes("foo"); partial != "foo" || index != nil {
		t.Error("could not extract no index partial")
	}
	if partial, _ := extractIndexes("foo[bar]"); partial != "foo[bar]" {
		t.Error("an index with alpha characters should be parsed as a segment")
	}
}

func TestWalkWithFunctions(t *testing.T) {
	ctx := context.Background()
	f0 := func(ctx context.Context, data any, params ...string) (any, error) {
		return "hello world", nil
	}
	f1 := func(ctx context.Context, data any, params ...string) (any, error) {
		return "I'm crazy: " + data.(map[string]any)["double_foo"].(string), nil
	}
	functions := NewFunctions()
	functions.Add("hello", f0).Add("debug", f1)
	if res, _ := Walk(ctx, "hello()", map[string]any{"foo": map[string]any{"double_foo": "bar"}}, functions); res != "hello world" {
		t.Error("basic function calling not working")
	}
	if res, _ := Walk(ctx, "foo.debug(\"bananas\")", map[string]any{"foo": map[string]any{"double_foo": "bar"}}, functions); res != "I'm crazy: bar" {
		t.Error("reflexive function calling not working")
	}
	if res, _ := Walk(ctx, "dawg()", map[string]any{"foo": map[string]any{"double_foo": "bar"}}, functions); res != "dawg()" {
		t.Error("calling a function that does not exist not working")
	}

}

func TestWalkWithChainedFunctionAndIndex(t *testing.T) {
	ctx := context.Background()
	functions := NewFunctions()
	if res, _ := Walk(ctx, "foo.split(|)[0]", map[string]string{"foo": "hello|world"}, functions); res != "hello" {
		t.Error("basic function chaining not working")
	}
}

func TestWalkWithChainedFunctions(t *testing.T) {
	ctx := context.Background()
	functions := NewFunctions()
	if res, _ := Walk(ctx, "foo.split(|).size()", map[string]string{"foo": "hello|world"}, functions); res != 2 {
		t.Error("basic function chaining not working")
	}

}

func TestWalkWithChainedIndexAndFunction(t *testing.T) {
	ctx := context.Background()
	if res, _ := Walk(ctx, "arr[0].size()", map[string]any{"arr": []string{"foo", "bar"}}, NewFunctions()); res != 3 {
		t.Error("basic function chaining not working")
	}
}

func TestWalkerWithFunctionAndDeadline(t *testing.T) {
	ctx, cancel := context.WithDeadline(context.TODO(), time.Now().Add(5*time.Millisecond))
	defer cancel()
	functions := NewFunctions()
	functions.Add("wait", func(ctx context.Context, data any, params ...string) (any, error) {
		time.Sleep(10 * time.Millisecond)
		return data, nil
	})
	if _, err := Walk(ctx, "wait().foo", map[string]string{"foo": "bar"}, functions); err.Error() != "deadline exceeded" {
		t.Error("deadline not working")
	}
}

func TestWalkerWithFunctionAndCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	functions := NewFunctions()
	functions.Add("wait", func(ctx context.Context, data any, params ...string) (any, error) {
		time.Sleep(10 * time.Millisecond)
		return data, nil
	})
	go func() {
		time.Sleep(5 * time.Millisecond)
		cancel()
	}()
	if _, err := Walk(ctx, "wait().foo", map[string]string{"foo": "bar"}, functions); err.Error() != "cancelled" {
		t.Error("cancellation not working")
	}
}
