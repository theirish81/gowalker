package gowalker

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

// Walk "walks" the provided data using the provided expression
func Walk(expr string, data any, functions Functions) (any, error) {
	if functions == nil {
		functions = NewFunctions()
	}
	return walkImpl(expr, data, nil, functions)
}

// walkImpl is the actual recursive implementation of the walker
func walkImpl(expr string, data any, indexes []int, functions Functions) (any, error) {
	switch reflect.TypeOf(data).Kind() { /*
		case map[string]int: {
			return walkImpl(expr,convertMap[int](t),indexes,functions)
		}
		case map[string]float64: {
			return walkImpl(expr,convertMap[float64](t),indexes,functions)
		}
		case map[string]string: {
			return walkImpl(expr,convertMap[string](t),indexes,functions)
		}
		case map[string]bool: {
			return walkImpl(expr,convertMap[bool](t),indexes,functions)
		}*/
	// if it's a map...
	case reflect.Map:
		t := reflect.ValueOf(data)
		// if there's no expression left, we can return the data we got as input.
		// This is the case in which the user wants a whole map returned
		if len(expr) == 0 {
			return data, nil
		}
		// splitting the current segment from the rest
		items := strings.SplitN(expr, ".", 2)
		next := ""
		// if we got at least one item, it means we're still selecting
		if len(items) > 0 {
			if found, res, err := runFunction(items[0], data, functions); err != nil {
				return res, err
			} else {
				if found {
					return res, nil
				}
			}

			// if we got more than 1 item, it means that not only we're still selecting, but there will be more
			// segments to select after. So we take the "next" part of the expression for the following recursion.
			if len(items) > 1 {
				next = items[1]
			}
			// If the segment contains one or more indexing blocks for arrays, then we separate the selector and
			//the indexes. If it doesn't contain indexes, then partial is still the correct selector, and indexes is null
			partial, indexes, err := extractIndexes(items[0])
			// if there was an error in the extraction of the indexes, then we return
			if err != nil {
				return data, err
			}
			val := t.MapIndex(reflect.ValueOf(partial))
			if val.IsValid() && !val.IsZero() {
				// recursion passing the selected value
				return walkImpl(next, t.MapIndex(reflect.ValueOf(partial)).Interface(), indexes, functions)
			} else {
				return nil, nil
			}
		} else {
			// we're not selecting anymore, we can return the value
			return items[0], nil
		}
	// if it's an array
	case reflect.Slice:
		t := reflect.ValueOf(data)
		// if there's one or more index selectors
		if indexes != nil || len(indexes) > 0 {
			// we pick the first index in the array
			nextIndex := indexes[0]
			// making sure that its value does not exceed the array size
			if indexes[0] < t.Len() {
				// popping the current index
				if len(indexes) == 1 {
					indexes = nil
				} else {
					indexes = indexes[1:]
				}

				// we select the indexed item and move forward
				return walkImpl(expr, t.Index(nextIndex).Interface(), indexes, functions)
			} else {
				// if the index exceeds the array size, we return an out-of-bounds error
				return t, errors.New("index out of bounds")
			}
		}
		// if someone is trying to access a property in an array...
		if len(expr) > 0 {
			if found, res, err := runFunction(expr, data, functions); err != nil {
				return res, err
			} else {
				if found {
					return res, nil
				}
			}
			//... then they're doing something wrong
			return nil, errors.New("cannot access attributes from an array")
		}
		// if this has no index, it means the user wants to return the entire array
		return t.Interface(), nil
	// all other data types
	default:
		return data, nil
	}
}

// extractIndexes tries to extract the index from an index notation. Will return the partial expression and an array
// of indexes as separate return values. If no index was found, then the indexes will be nil. Indexes is an array
// in case a user is selecting nested arrays, such as array[0][1]
func extractIndexes(expr string) (string, []int, error) {
	// we find the indexing notation blocks
	bits := indexExtractorRegex.FindAllStringSubmatch(expr, 100)
	// no indexing notation block?
	if bits == nil || len(bits) == 0 {
		// then the expression has no indexing notation. We return the expression and -1
		return expr, nil, nil
	}
	// otherwise, we take care of removing the entire indexing notation from the string. We should be left with
	// the expression alone
	partial := indexExtractorRegex.ReplaceAllString(expr, "")

	// converting each found index to an integer and composing the final indexes array
	indexes := make([]int, 0)
	for _, bx := range bits {
		index, err := strconv.Atoi(bx[1])
		if err != nil {
			return partial, indexes, err
		}
		indexes = append(indexes, index)
	}
	// and return
	return partial, indexes, nil
}
