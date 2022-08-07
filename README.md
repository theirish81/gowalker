# GoWalker
GoWalker is two things:

## A simple path expression interpreter
A very simple data expression interpreter.
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
You can use the `Walk` function to easily navigate the data structure as in:
```go
Walk("name",data,nil)               // returns `Joe`
Walk("items[1]",data,nil)           // returns `wallet`
Walk("friends[0].name",data,nil)    // returns `billy`
Walk("items",data,nil)              // returns ["keys","wallet"]
```
The library uses no code evaluations therefore it's super safe.


## A simple template engine
Powered by the same path expression interpreter, this tiny template engine allows you to substitute strings with
data coming from a map. As in:
```
{
  "name": "${name}",
  "first_item": "${items[0]}",
  "all_items": ${items}
}
```
When a complex object is referenced in an expression, the rendering engine will automatically convert it to its
JSON counterpart.

### Functions
Expressions also support the use of functions.
From the expression parser standpoint, assertions work as follows:
* at any point of the expression you can invoke a function
* the function call must be the last segment in an expression
* function chaining is not supported yet
* functions can be reflective and can only operate on the piece of data they've been called upon
* functions can receive comma separated parameters. Quotation is not required as data typing will be handled by the
  function implementation
* running a function without a preceding expression will make the function operate on the full scope

Examples:
```
foo.bar.size()
```
Will evaluate the size of `bar`.
```
foo.myString.split(|)
```
Will split `myString` using pipe as separator.
```
foo.myArray.collect(banana,mango)
```
Where myArray is an array of objects, it will collect all the fields named `banana` and `mango`.

### Implementing functions
The engine comes with just a few of default functions for demonstration purposes, such as:
* `size()`: returns the size of the object in scope
* `split(sep)`: splits the string in scope, using a separator
* `collect(...)`: given an array containing maps, it will return an array of maps in which the maps only show the
  provided keys

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
functions.Add("sayHello",func (scope any, params ...any) (any, error) {
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
Walk("items[0].sayHello(Barney)",data,functions)
```
will return:
```
hello foo from Barney
```