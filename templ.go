package gowalker

import (
	"encoding/json"
	"strconv"
	"strings"
)

func Render(template string, data map[string]interface{}) (string, error) {
	items := templateFinderRegex.FindAllStringSubmatch(template, -1)
	for _, item := range items {
		matcher := item[0]
		expr := item[1]
		if val, err := Walk(expr, data); err == nil {
			template = strings.Replace(template, matcher, convertData(val), 1)
		} else {
			return template, err
		}
	}
	return template, nil
}

func convertData(data interface{}) string {
	switch t := data.(type) {
	case int:
		return strconv.Itoa(t)
	case bool:
		return strconv.FormatBool(t)
	case []interface{}, map[string]interface{}:
		d, _ := json.Marshal(data)
		return string(d)
	default:
		return data.(string)
	}
}
