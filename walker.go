package gowalker

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

// Walk "walks" the provided data using the provided expression
func Walk(expr string, data any, functions *Functions) (any, error) {
	if functions == nil {
		functions = NewFunctions()
	}
	return walkImpl(expr, data, nil, functions)
}

// walkImpl is the actual recursive implementation of the walker
func walkImpl(expr string, data any, indexes []int, functions *Functions) (any, error) {
	// if data is nil, then check if there's a function to run against it. This generally does not happen, but you
	// never know someone wants to do something with that nil
	if data == nil {
		found, res, err := runFunction(expr, data, functions)
		if found || err != nil {
			return res, err
		}
		return data, nil
	}
	if expr == "." {
		return data, nil
	}
	// Let's check the kind of data
	switch reflect.TypeOf(data).Kind() {
	// if it's a map...
	case reflect.Map:
		t := reflect.ValueOf(data)
		// if there's no expression left, we can return the data we got as input.
		// This is the case in which the user wants a whole map returned
		if len(expr) == 0 {
			return data, nil
		}
		// splitting the current segment from the rest
		current, next := getSegments(expr)
		// if we got at least one item, it means we're still selecting
		if current != "" {

			// If the segment contains one or more indexing blocks for arrays, then we separate the selector and
			//the indexes. If it doesn't contain indexes, then partial is still the correct selector, and indexes is null
			partial, indexes, err := extractIndexes(current)

			// if there was an error in the extraction of the indexes, then we return
			if err != nil {
				return data, err
			}

			if found, res, err := runFunction(partial, data, functions); err != nil {
				return res, err
			} else {
				if found {
					return walkImpl(next, res, indexes, functions)
				}
			}

			val := t.MapIndex(reflect.ValueOf(partial))
			if val.IsValid() && !val.IsZero() {
				// recursion passing the selected value
				return walkImpl(next, t.MapIndex(reflect.ValueOf(partial)).Interface(), indexes, functions)
			} else {
				return walkImpl(next, nil, indexes, functions)
			}
		} else {
			// we're not selecting anymore, we can return the value
			return current, nil
		}
	// if it's a slice...
	case reflect.Slice:
		t := reflect.ValueOf(data)
		// if there's one or more index selectors
		if indexes != nil || len(indexes) > 0 {
			// we pick the first index in the array
			nextIndex, indexes := sliceOneOff(indexes)
			// making sure that its value does not exceed the array size
			if nextIndex < t.Len() {
				// we select the indexed item and move forward
				return walkImpl(expr, t.Index(nextIndex).Interface(), indexes, functions)
			} else {
				// if the index exceeds the array size, we return an out-of-bounds error
				return t, errors.New("index out of bounds")
			}
		}
		// if someone is trying to access a property in an array...
		if len(expr) > 0 {
			// we try to understand if it's one fo the available functions, as it's totally legit
			found, res, err := runFunction(expr, data, functions)
			if found || err != nil {
				return res, err
			}
			//... if it's not a function, they're probably doing something wrong
			return nil, errors.New("cannot access attributes from an array")
		}
		// if this has no index, it means the user wants to return the entire array
		return t.Interface(), nil
	// all other data types
	default:
		current, next := getSegments(expr)

		// If the segment contains one or more indexing blocks for arrays, then we separate the selector and
		//the indexes. If it doesn't contain indexes, then partial is still the correct selector, and indexes is null
		partial, indexes, err := extractIndexes(current)

		// if there was an error in the extraction of the indexes, then we return
		if err != nil {
			return data, err
		}

		// let's check if we need to run a function against it
		if found, res, err := runFunction(partial, data, functions); err != nil {
			return res, err
		} else {
			if found {
				return walkImpl(next, res, indexes, functions)
			}
		}
		// otherwise, we just return the value
		return data, nil
	}
}

// getSegments will receive an expression, do a one-split on the dot and return the results.
// The first returned value is the "current" segment being evaluated, while the second is the "remaining part" of the
// expression. In absence of a current element or a remaining part, empty strings will be returned
func getSegments(expr string) (string, string) {
	items := strings.SplitN(expr, ".", 2)
	current := ""
	next := ""
	if len(items) > 0 {
		current = items[0]
		if len(items) > 1 {
			next = items[1]
		}
	}
	return current, next
}

// sliceOneOff will take an array of indexes, take the head off, and return the resulting array
func sliceOneOff(indexes []int) (int, []int) {
	if len(indexes) == 1 {
		return indexes[0], nil
	} else {
		return indexes[0], indexes[1:]
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
		if index, err := strconv.Atoi(bx[1]); err == nil {
			indexes = append(indexes, index)
		} else {
			return partial, indexes, err
		}
	}
	// and return
	return partial, indexes, nil
}
