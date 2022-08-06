package gowalker

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

// Functions is a map of actual Golang functions the expression can call.
// When invoked, a function receives a variadic argument in which the first position is always the current selected
// element in the expression, while the following ones are provided as params.
// Functions will return a value and an error
type Functions map[string]func(data ...any) (any, error)

// NewFunctions is the constructor of Functions and adds some very basic implementations
func NewFunctions() Functions {
	fx := Functions{}
	fx.Add("size", fx.size)
	fx.Add("split", fx.split)
	return fx
}

// Add adds a function ot the Functions' data structure
func (f *Functions) Add(key string, function func(data ...any) (any, error)) *Functions {
	(*f)[key] = function
	return f
}

// size is one of the base functions for the user to invoke.
// It returns the size of maps, slices and strings
func (f *Functions) size(data ...any) (any, error) {
	scope := data[0]
	if scope == nil {
		return nil, errors.New("nil reference to size function")
	}
	val := reflect.ValueOf(scope)
	kind := reflect.TypeOf(scope).Kind()
	switch kind {
	case reflect.Map, reflect.Slice, reflect.String:
		return val.Len(), nil
	default:
		return nil, errors.New("size not supported for: " + kind.String())
	}
}

// split is one of the base functions for the user to invoke.
// It splits a string into an array, given a separator
func (f *Functions) split(data ...any) (any, error) {
	if len(data) < 2 {
		return nil, errors.New("separator not provided")
	}
	if val, ok := data[0].(string); ok {
		if sep, ok := data[1].(string); ok {
			return strings.Split(val, sep), nil
		} else {
			return nil, errors.New("separator not a string")
		}
	} else {
		return nil, errors.New("split only supported for strings")
	}
}

func extractFunctionName(expr string) string {
	names := functionExtractorRegex.FindStringSubmatch(expr)
	if names != nil && len(names) > 1 {
		return names[1]
	}
	return ""
}

func extractParameterString(expr string) string {
	names := functionExtractorRegex.FindStringSubmatch(expr)
	if names != nil && len(names) > 2 {
		return names[3]
	}
	return ""
}
func extractParameters(signature string) []string {
	paramString := extractParameterString(signature)
	params := paramExtractRegex.FindAllString(paramString, -1)
	return params
}

func runFunction(expr string, data interface{}, functions Functions) (bool, interface{}, error) {
	if fx := extractFunctionName(expr); fx != "" {
		params := append([]interface{}{data}, retypeParams(extractParameters(expr))...)
		if function, ok := functions[fx]; ok {
			res, err := function(params...)
			return true, res, err
		} else {
			return true, expr, nil
		}
	}
	return false, "", nil
}

func retypeParams(params []string) []interface{} {
	res := make([]interface{}, 0)
	for _, p := range params {
		if rx, err := strconv.Atoi(p); err == nil {
			res = append(res, rx)
		}
		if rx, err := strconv.ParseBool(p); err == nil {
			res = append(res, rx)
		}
		res = append(res, p)
	}
	return res
}
