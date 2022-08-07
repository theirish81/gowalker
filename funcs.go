package gowalker

import (
	"errors"
	"reflect"
	"strings"
)

// Functions is a map of actual Golang functions the expression can call.
// When invoked, a function receives a variadic argument in which the first position is always the current selected
// element in the expression, while the following ones are provided as params.
// Functions will return a value and an error
type Functions map[string]func(scope any, params ...string) (any, error)

// NewFunctions is the constructor of Functions and adds some very basic implementations
func NewFunctions() Functions {
	fx := Functions{}
	fx.Add("size", fx.size)
	fx.Add("split", fx.split)
	return fx
}

// Add adds a function ot the Functions' data structure
func (f *Functions) Add(key string, function func(data any, params ...string) (any, error)) *Functions {
	(*f)[key] = function
	return f
}

// size is one of the base functions for the user to invoke.
// It returns the size of maps, slices and strings
func (f *Functions) size(scope any, _ ...string) (any, error) {
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
func (f *Functions) split(scope any, params ...string) (any, error) {
	if len(params) < 1 {
		return nil, errors.New("separator not provided")
	}
	if val, ok := scope.(string); ok {
		return strings.Split(val, params[0]), nil
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
	for i, p := range params {
		params[i] = strings.ReplaceAll(p, "\\,", ",")
	}
	return params
}

func runFunction(expr string, data interface{}, functions Functions) (bool, interface{}, error) {
	if fx := extractFunctionName(expr); fx != "" {
		params := extractParameters(expr)
		if function, ok := functions[fx]; ok {
			res, err := function(data, params...)
			return true, res, err
		} else {
			return true, expr, nil
		}
	}
	return false, "", nil
}
