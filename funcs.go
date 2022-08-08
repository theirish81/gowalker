package gowalker

import (
	"errors"
	"reflect"
	"strings"
)

type mapOfFunctions map[string]func(scope any, params ...string) (any, error)

// Functions is a map of actual Golang functions the expression can call.
// When invoked, a function receives a variadic argument in which the first position is always the current selected
// element in the expression, while the following ones are provided as params.
// Functions will return a value and an error
type Functions struct {
	mapOfFunctions
	functionScope map[string]interface{}
}

// NewFunctions is the constructor of Functions and adds some very basic implementations
func NewFunctions() *Functions {
	fx := Functions{mapOfFunctions{}, map[string]interface{}{}}
	fx.Add("size", fx.size)
	fx.Add("split", fx.split)
	fx.Add("collect", fx.collect)
	fx.Add("render", fx.render)
	fx.Add("renderEach", fx.renderEach)
	return &fx
}

// Add adds a function ot the Functions' data structure
func (f *Functions) Add(key string, function func(data any, params ...string) (any, error)) *Functions {
	f.mapOfFunctions[key] = function
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

func (f *Functions) render(scope any, params ...string) (any, error) {
	if len(params) < 1 {
		return nil, errors.New("template not provided")
	}
	if templ, ok := f.functionScope["_"+params[0]]; ok {
		return Render(templ.(string), scope, f)
	} else {
		return nil, errors.New("template not found")
	}
}

func (f *Functions) renderEach(scope any, params ...string) (any, error) {
	if len(params) < 1 {
		return nil, errors.New("template not provided")
	}
	sep := ""
	if len(params) == 2 {
		sep = params[1]
	}
	if templ, ok := f.functionScope["_"+params[0]]; ok {
		if reflect.TypeOf(scope).Kind() == reflect.Slice {
			sliceVal := reflect.ValueOf(scope)
			res := ""
			for i := 0; i < sliceVal.Len(); i++ {
				if tmp, err := Render(templ.(string), sliceVal.Index(i).Interface(), f); err == nil {
					res = res + tmp
					if i < sliceVal.Len()-1 {
						res += sep
					}
				} else {
					return nil, err
				}
			}
			return res, nil
		} else {
			return nil, errors.New("cannot iterate on a data type that is not an array")
		}
	} else {
		return nil, errors.New("template not found")
	}
}

// collect expects an array of objects to be the scope. Params is a list of fields we're interested in.
// The function will produce a derivative array of objects containing only the fields expressed in params
func (f *Functions) collect(scope any, params ...string) (any, error) {
	// if no params are passed, then we return an error
	if len(params) < 1 {
		return nil, errors.New("list of fields not provided")
	}
	kind := reflect.TypeOf(scope).Kind()
	// if it's a slice, then the scope probably is the correct data type
	if kind == reflect.Slice {
		data := reflect.ValueOf(scope)
		size := data.Len()
		// let's create an array of maps that's going to hold the results
		res := make([]map[string]interface{}, size)
		// iterating the original array
		for i := 0; i < size; i++ {
			item := data.Index(i)
			// all elements have to be maps
			if item.Kind() == reflect.Map {
				// creating derivative element
				block := map[string]interface{}{}
				// for each parameter expressed in the arguments
				for _, p := range params {
					// if we find an attribute with that name
					found := item.MapIndex(reflect.ValueOf(p))
					if found.IsValid() && !found.IsZero() && found.CanInterface() {
						// we copy the value over
						block[p] = found.Interface()
					}
				}
				// assigning the newly created map to the new array
				res[i] = block
			} else {
				// if the `kind` of a child object is not a map, then it's an error
				return nil, errors.New("at least one item in the array is not a map")
			}
		}
		// returning the derivative array
		return res, nil
	} else {
		// if the given scope is not even an array, then we return an error
		return nil, errors.New("operation can only be applied to arrays of maps")
	}
}

// extractFunctionName will extract the function name from an expression. If the expression doesn't look like a function
// call, then it returns an empty string
func extractFunctionName(expr string) string {
	names := functionExtractorRegex.FindStringSubmatch(expr)
	if names != nil && len(names) > 1 {
		return names[1]
	}
	return ""
}

// extractParameterString will extract the portion of the parameters from an expression that looks like a function call
func extractParameterString(expr string) string {
	names := functionExtractorRegex.FindStringSubmatch(expr)
	if names != nil && len(names) > 2 {
		return names[3]
	}
	return ""
}

// extractParameters will try to extract the parameters from a function call string
func extractParameters(signature string) []string {
	paramString := extractParameterString(signature)
	params := paramExtractRegex.FindAllString(paramString, -1)
	for i, p := range params {
		params[i] = strings.ReplaceAll(p, "\\,", ",")
	}
	return params
}

// runFunction will try to run the function expressed by expr, against data, with the provided functions.
// The first return value will be `true` if a function was indeed found in expr and the execution of the function
// was attempted. If the functions ran, the second return value will be the result of the function execution.
// The third parameter is an error, in case the function failed
func runFunction(expr string, data interface{}, functions *Functions) (bool, interface{}, error) {
	if fx := extractFunctionName(expr); fx != "" {
		params := extractParameters(expr)
		if function, ok := functions.mapOfFunctions[fx]; ok {
			res, err := function(data, params...)
			return true, res, err
		} else {
			return true, expr, nil
		}
	}
	return false, "", nil
}
