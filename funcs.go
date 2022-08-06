package gowalker

import (
	"strconv"
	"strings"
)

type Functions map[string]func(data ...interface{}) (interface{}, error)

func (f *Functions) Add(key string, function func(data ...interface{}) (interface{}, error)) *Functions {
	(*f)[key] = function
	return f
}

func ExtractFunctionName(expr string) string {
	names := functionExtractorRegex.FindStringSubmatch(expr)
	if names != nil && len(names) > 1 {
		return names[1]
	}
	return ""
}

func ExtractParameterString(expr string) string {
	names := functionExtractorRegex.FindStringSubmatch(expr)
	if names != nil && len(names) > 2 {
		return names[3]
	}
	return ""
}
func ExtractParameters(signature string) []string {
	paramString := ExtractParameterString(signature)
	params := paramExtractRegex.FindAllString(paramString, -1)
	return params
}

func RunFunction(expr string, data interface{}, functions Functions) (bool, interface{}, error) {
	if fx := ExtractFunctionName(expr); fx != "" {
		params := append([]interface{}{data}, RetypeParams(ExtractParameters(expr))...)
		if function, ok := functions[fx]; ok {
			res, err := function(params...)
			return true, res, err
		} else {
			return true, expr, nil
		}
	}
	return false, "", nil
}

func RetypeParams(params []string) []interface{} {
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
