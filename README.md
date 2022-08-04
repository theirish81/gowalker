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
Walk("name",data)               // returns `Joe`
Walk("items[1]",data)           // returns `wallet`
Walk("friends[0].name",data)    // returns `billy`
Walk("items",data)              // returns ["keys","wallet"]
```
The library uses no code evaluations therefore it's super safe.


## A simple template engine for text
Powered by the same path expression interpreter, this tiny template engine allows you to substitute strings with
data coming from a map. As in:
```
{
  "name": "${name}",
  "first_item": "${items[0]}",
  "all_items": ${items}
}
```