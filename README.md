# GoWalker
**Status:**
[![CircleCI](https://dl.circleci.com/status-badge/img/gh/theirish81/gowalker/tree/main.svg?style=svg)](https://dl.circleci.com/status-badge/redirect/gh/theirish81/gowalker/tree/main)

GoWalker is two things:

## A data path expression interpreter
Given a data structure like this:
```json
{
  "name": "Joe",
  "age": 22,
  "friends": [
    {
      "name": "billy",
      "age": 27
    } ,
    {
      "name": "john",
      "age": 23
    }
  ],
  "items": [
    "keys",
    "wallet"
  ]
}
```
You can use the `Walk` function to easily navigate the data structure, by providing a string that represents the path
to the data, as in:
```go
ctx := context.TODO()
Walk(ctx, "name",data,nil)               // returns `Joe`
Walk(ctx, "items[1]",data,nil)           // returns `wallet`
Walk(ctx, "friends[0].name",data,nil)    // returns `billy`
Walk(ctx, "items",data,nil)              // returns ["keys","wallet"]
```
The library uses no code evaluations therefore it's super safe.

### Expressions
Expressions are actually pretty easy. A few notes:
* the path separator through maps is the `.` (dot). No square-bracket notation is supported or required
* The index expression in arrays uses the square-bracket (`[n]`) notation
* the `.` (dot) alone in an expression refers to the whole scope

### Functions
Expressions also support the use of functions.
From the expression parser standpoint, assertions work as follows:
* at any point of the expression you can invoke a function
* functions can be reflexive and can only operate on the piece of data they've been called upon
* functions can receive comma separated parameters. Quotation is not required as data typing will be handled by the
  function implementation
* running a function without a preceding expression will make the function operate on the full scope
* you can chain functions, object and index selectors

Examples:
```text
foo.bar.size()
```
Will evaluate the size of `bar`.
```text
foo.myString.split(|)
```
Will split `myString` using pipe as separator.
```text
foo.myArray.collect(banana,mango)
```
Where myArray is an array of objects, it will collect all the fields named `banana` and `mango`.

### Implementing functions
The engine comes with just a few of default functions for demonstration purposes, such as:
* `size()`: returns the size of the object in scope
* `split(sep)`: splits the string in scope, using a separator
* `collect(...)`: given an array containing maps, it will return an array of maps in which the maps only show the
  provided keys
* `toVar(varName)`: will return a variable from the *Functions extra variables* and ignore the provided data 

You can implement more by passing the `functions` parameter when invoking `Walk`.
Example:

Assuming you have a data structure as follows:
```json
{
  "items": [
    "foo",
    "bar"
  ]
}
```

```go
functions := NewFunctions()
functions.Add("sayHello",func (context context.Context, scope any, params ...string) (any, error) {
	if len(params) < 1 {
		return nil,errors.New("not enough parameters")
    }
	if data,ok := scope.(string); ok {
        return "hello "data+" from "+params[0]	
    } else {
        return nil, errors.New("cannot run sayHello against a data type that is not string")
    }
})
//...
ctx := context.TODO()
Walk(ctx, "items[0].sayHello(Barney)", data,functions)
```
will return:
```text
hello foo from Barney
```

### Functions extra variables
Functions can also access another map of variables, unrelated to the data they're evaluating. This may be useful if
your custom functions need to interact with other pieces of information beyond the data itself, such as request params.
This map of variables can be accessed by invoking `getScope()` in a `Functions` instance.

If, for example, you wanted to add a variable to the scope, you could simply:
```go
functions := NewFunctions()
functions.GetScope()["foo"] = "bar"
```

## A simple template engine
Powered by the same path expression interpreter, this tiny template engine allows you to substitute strings with
data coming from a map. As in:
```text
{
  "name": "${name}",
  "first_item": "${items[0]}",
  "all_items": ${items}
}
```
When a complex object is referenced in an expression, the rendering engine will automatically convert it to its
JSON counterpart.

Just call:
```go
data := map[string]any{"name": "pino", "items": []any{"keys", "wallet"}}
templ := `{
    "name": "${name}",
    "first_item": "${items[0]}",
    "all_items": ${items}
}`
ctx := context.TODO()
res, _ := Render(ctx, templ, data, nil)
```
and you're set. You can, of course, pass a `Functions` instance as third parameter.

### Sub-templates
Sometimes you need to split your templates into multiple files. There are typically two scenarios when this is
recommended in GoWalker:
* When you want to share a sub-template across multiple master templates
* When you need to run a template against each item in an array

Here's an example of simple template splitting. It uses the `render` function against `items`
```go
t1 := "this is a test ${items.render(t2)}"
t2 := "T2 ${.}"
templates := NewTemplates()
templates.Add("t2",t2)
ctx := context.TODO()
res, _ := RenderAll(ctx, t1, templates, map[string]any{"items": []string{"foo", "bar"}}, NewFunctions())
// prints:
// `this is a test T2 ["foo","bar"]`
}
```

* `render(templateName)`: renders a sub-template against the variable it was run against

And here's an example where we iterate over an array. It uses the `renderEach` function against `items`:
```go
t1 := "this is a test ${items.renderEach(t2,\\,)}"
t2 := "\nT2 ${.}"
templates := NewTemplates()
templates.Add("t2",t2)
ctx := context.TODO()
res, _ := RenderAll(ctx, t1, templates, map[string]any{"items": []string{"foo", "bar"}}, NewFunctions())
// prints:
// this is a test
// T2 foo
// T2 bar
```

* `renderEach(templateName,sep?)`: renders a sub-template against each item of the array it was run against.
  Additionally, you can provide an optional separator string that will be printed between an iteration and the next


## Cancellation and deadlines
As rendering large templates (or selecting complex paths) can be memory and CPU intensive, all functions now receive
a context as first parameter, supporting both deadlines and cancellations.