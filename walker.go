package gowalker

import (
	"errors"
	"strconv"
	"strings"
)

//Walk "walks" the provided data using the provided expression
func Walk(expr string, data interface{}) (interface{}, error) {
	return walkImpl(expr, data, -1)
}

// walkImpl is the actual recursive implementation of the walker
func walkImpl(expr string, data interface{}, index int) (interface{}, error) {
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
			// If the segment contains an indexing block for arrays, then we separate the selector and the index.
			// If it doesn't contain an index, then partial is still the correct selector, and index=-1
			partial, index, err := ExtractIndex(items[0])
			// if there was an error in the extraction of the index, then we return
			if err != nil {
				return data, err
			}
			// recursion passing the selected value
			return walkImpl(next, t[partial], index)
		} else {
			// we're not selecting anymore, we can return the value
			return items[0], nil
		}
	// if it's an array
	case []interface{}:
		// if there's an index selector
		if index > -1 {
			// and the index is not overflowing the array
			if index < len(t) {
				// we select the indexed item and move forward
				return walkImpl(expr, t[index], -1)
			} else {
				// otherwise, we return an out-of-bounds error
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

// ExtractIndex tries to extract the index from an index notation. Will return the partial expression and the index
// as separate return values. If no index was found, then the index will be -1
func ExtractIndex(expr string) (string, int, error) {
	// we find the indexing notation block
	bits := indexExtractorRegex.FindStringSubmatch(expr)
	// no indexing notation block?
	if bits == nil {
		// then the expression has no indexing notation. We return the expression and -1
		return expr, -1, nil
	}
	// otherwise, we take care of removing the entire indexing notation from the string. We should be left with
	// the expression alone
	partial := indexExtractorRegex.ReplaceAllString(expr, "")
	// we convert the index to an integer
	index, err := strconv.Atoi(bits[1])
	// and return
	return partial, index, err
}
