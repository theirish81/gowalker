package gowalker

import (
	"errors"
	"strconv"
	"strings"
)

type Functions map[string]func(data ...any) (any, error)

func NewFunctions() Functions {
	fx := Functions{}
	fx.Add("size", fx.size)
	return fx
}

func (f *Functions) Add(key string, function func(data ...any) (any, error)) *Functions {
	(*f)[key] = function
	return f
}

func (f *Functions) size(data ...any) (any, error) {
	scope := data[0]
	switch t := scope.(type) {
	case int, float64, bool, nil:
		return nil, errors.New("size no available for data type")
	case map[string]any:
		return len(t), nil
	case []any:
		return len(t), nil
	default:
		return 0, nil
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
		if strings.HasPrefix(p, "\"") && strings.HasSuffix(p, "\"") {
			res = append(res, p[1:len(p)-1])
		} else {
			if rx, err := strconv.Atoi(p); err == nil {
				res = append(res, rx)
			}
			if rx, err := strconv.ParseBool(p); err == nil {
				res = append(res, rx)
			}
		}
	}
	return res
}
