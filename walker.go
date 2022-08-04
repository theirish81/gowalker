package gowalker

import (
	"errors"
	"strconv"
	"strings"
)

func Walk(expr string, data interface{}) (interface{}, error) {
	return walkImpl(expr, data, -1)
}

func walkImpl(expr string, data interface{}, index int) (interface{}, error) {
	switch t := data.(type) {
	case map[string]interface{}:
		if len(expr) == 0 {
			return data, nil
		}
		items := strings.SplitN(expr, ".", 2)
		next := ""
		if len(items) > 0 {
			if len(items) > 1 {
				next = items[1]
			}
			partial, index := ExtractIndex(items[0])
			return walkImpl(next, t[partial], index)
		} else {
			return handleStraightValue(items[0]), nil
		}
	case []interface{}:
		if index > -1 {
			if index < len(t) {
				return walkImpl(expr, t[index], -1)
			} else {
				return t, errors.New("index out of bounds")
			}
		}
		if index > -1 && len(expr) > 0 {
			return walkImpl(expr, t, -1)
		}
		if len(expr) > 0 {
			return nil, errors.New("cannot access attributes from an array")
		}
		return handleStraightValue(t), nil
	default:
		return handleStraightValue(data), nil
	}
}

func handleStraightValue(data interface{}) interface{} {
	return data
}

func ExtractIndex(expr string) (string, int) {
	bits := indexExtractorRegex.FindStringSubmatch(expr)
	if bits == nil {
		return expr, -1
	}
	partial := indexExtractorRegex.ReplaceAllString(expr, "")
	if len(bits) > 1 {
		index, _ := strconv.Atoi(bits[1])
		return partial, index
	}
	return partial, -1
}
