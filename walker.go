package gowalker

import (
	"errors"
	"strconv"
	"strings"
)

//Walk "walks" the provided data using the provided expression
func Walk(expr string, data interface{}) (interface{}, error) {
	return walkImpl(expr, data, nil)
}

// walkImpl is the actual recursive implementation of the walker
func walkImpl(expr string, data interface{}, indexes []int) (interface{}, error) {
	switch t := data.(type) {
	// if it's a map...
	case map[string]interface{}:
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
			// if we got more than 1 item, it means that not only we're still selecting, but there will be more
			// segments to select after. So we take the "next" part of the expression for the following recursion.
			if len(items) > 1 {
				next = items[1]
			}
			// If the segment contains one or more indexing blocks for arrays, then we separate the selector and
			//the indexes. If it doesn't contain indexes, then partial is still the correct selector, and indexes is null
			partial, indexes, err := ExtractIndexes(items[0])
			// if there was an error in the extraction of the indexes, then we return
			if err != nil {
				return data, err
			}
			// recursion passing the selected value
			return walkImpl(next, t[partial], indexes)
		} else {
			// we're not selecting anymore, we can return the value
			return items[0], nil
		}
	// if it's an array
	case []interface{}:
		// if there's one or more index selectors
		if indexes != nil || len(indexes) > 0 {
			// we pick the first index in the array
			nextIndex := indexes[0]
			// making sure that its value does not exceed the array size
			if indexes[0] < len(t) {
				// popping the current index
				if len(indexes) == 1 {
					indexes = nil
				} else {
					indexes = indexes[1:]
				}

				// we select the indexed item and move forward
				return walkImpl(expr, t[nextIndex], indexes)
			} else {
				// if the index exceeds the array size, we return an out-of-bounds error
				return t, errors.New("index out of bounds")
			}
		}
		// if someone is trying to access a property in an array...
		if len(expr) > 0 {
			//... then they're doing something wrong
			return nil, errors.New("cannot access attributes from an array")
		}
		// if this has no index, it means the user wants to return the entire array
		return t, nil
	// all other data types
	default:
		return data, nil
	}
}

// ExtractIndexes tries to extract the index from an index notation. Will return the partial expression and an array
// of indexes as separate return values. If no index was found, then the indexes will be nil. Indexes is an array
// in case a user is selecting nested arrays, such as array[0][1]
func ExtractIndexes(expr string) (string, []int, error) {
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
