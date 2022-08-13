package gowalker

import (
	"context"
	"errors"
	"reflect"
	"strings"
)

type mapOfFunctions map[string]func(ctx context.Context, scope any, params ...string) (any, error)

// Functions is a map of actual Golang functions the expression can call.
// When invoked, a function receives a variadic argument in which the first position is always the current selected
// element in the expression, while the following ones are provided as params.
// Functions will return a value and an error.
// mapOfFunctions is the actual map of key=function
// functionScope is an extra scope a function can access if part of Functions
type Functions struct {
	mapOfFunctions
	functionScope map[string]any
}

// NewFunctions is the constructor of Functions and adds some very basic implementations
func NewFunctions() *Functions {
	fx := Functions{mapOfFunctions{}, map[string]any{}}
	fx.Add("size", fx.size)
	fx.Add("split", fx.split)
	fx.Add("collect", fx.collect)
	fx.Add("render", fx.render)
	fx.Add("renderEach", fx.renderEach)
	return &fx
}

// Add adds a function ot the Functions' data structure
func (f *Functions) Add(key string, function func(ctx context.Context, data any, params ...string) (any, error)) *Functions {
	f.mapOfFunctions[key] = function
	return f
}

// size is one of the base functions for the user to invoke.
// It returns the size of maps, slices and strings
func (f *Functions) size(_ context.Context, scope any, _ ...string) (any, error) {
	// if scope is nil, then we return an error
	if scope == nil {
		return nil, errors.New("nil reference to size function")
	}
	val := reflect.ValueOf(scope)
	kind := reflect.TypeOf(scope).Kind()
	switch kind {
	// maps, slices and strings, all support Len
	case reflect.Map, reflect.Slice, reflect.String:
		return val.Len(), nil
	default:
		// in any other cases, we return an error
		return nil, errors.New("size not supported for: " + kind.String())
	}
}

// split is one of the base functions for the user to invoke.
// It splits a string into an array, given a separator
func (f *Functions) split(_ context.Context, scope any, params ...string) (any, error) {
	// if the scope is a string, then we can proceed with the split
	if val, ok := scope.(string); ok {
		return strings.Split(val, params[0]), nil
	} else {
		// returning an error if attempting a split on a data type that is not a string
		return nil, errors.New("split only supported for strings")
	}
}

// render will render a sub-template against the selected scope. It requires one param that is the name of the
// sub-template
func (f *Functions) render(ctx context.Context, scope any, params ...string) (any, error) {
	// returning an error if the sub-template name was not provided
	if len(params) < 1 || len(params[0]) == 0 {
		return nil, errors.New("template not provided")
	}
	// if the sub-template name is found, we can run Render against it
	if templ, ok := f.functionScope["_"+params[0]]; ok {
		return Render(ctx, templ.(string), scope, f)
	} else {
		// returning an error if the template was not found
		return nil, errors.New("template not found")
	}
}

// renderEach will render a sub-template against each element in the provided scope, assuming it's an array.
// It requires one param that is the name of the sub-template. Additionally, it accepts a second param that is a
// separator string to append at each iteration.
func (f *Functions) renderEach(ctx context.Context, scope any, params ...string) (any, error) {
	// if there are no params, it's an error
	if len(params) < 1 {
		return nil, errors.New("template not provided")
	}
	// if it has two params, we have a separator string
	sep := ""
	if len(params) == 2 {
		sep = params[1]
	}
	// if the sub-template exists
	if templ, ok := f.functionScope["_"+params[0]]; ok {
		// and the scope is a slice
		if reflect.TypeOf(scope).Kind() == reflect.Slice {
			sliceVal := reflect.ValueOf(scope)
			res := ""
			// against each item in the slice
			for i := 0; i < sliceVal.Len(); i++ {
				// we render the sub-template
				if tmp, err := Render(ctx, templ.(string), sliceVal.Index(i).Interface(), f); err == nil {
					res = res + tmp
					// if this is not the last item in the list, we print the separator character
					if i < sliceVal.Len()-1 {
						res += sep
					}
				} else {
					// returning an error if Render failed
					return nil, err
				}
			}
			// returning the collected strings
			return res, nil
		} else {
			// returning an error if the data type of the scope was not a slice
			return nil, errors.New("cannot iterate on a data type that is not an array")
		}
	} else {
		// returning an error if the template was not found
		return nil, errors.New("template not found")
	}
}

// collect expects an array of objects to be the scope. Params is a list of fields we're interested in.
// The function will produce a derivative array of objects containing only the fields expressed in params
func (f *Functions) collect(_ context.Context, scope any, params ...string) (any, error) {
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
		res := make([]map[string]any, size)
		// iterating the original array
		for i := 0; i < size; i++ {
			item := data.Index(i).Interface()
			// all elements have to be maps
			if reflect.TypeOf(item).Kind() == reflect.Map {
				// creating derivative element
				block := map[string]any{}
				// for each parameter expressed in the arguments
				for _, p := range params {
					// if we find an attribute with that name
					found := reflect.ValueOf(item).MapIndex(reflect.ValueOf(p))
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
func runFunction(ctx context.Context, expr string, data any, functions *Functions) (bool, any, error) {
	// Extracting the function name. If empty, then this is not a function call
	if fx := extractFunctionName(expr); fx != "" {
		// If it's a function call, though, we extract the parameters
		params := extractParameters(expr)
		// If the provided functions do contain the one being invoked...
		if function, ok := functions.mapOfFunctions[fx]; ok {
			// ... we can run it and return the result
			res, err := function(ctx, data, params...)
			return true, res, err
		} else {
			// otherwise, we still report that the function was detected, but as it was not found, the function call
			// is returned as value, so it can be printed.
			return true, expr, errors.New("function not found")
		}
	}
	// If no functions were found, we report back
	return false, "", nil
}
