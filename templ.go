package gowalker

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Render renders a template, using the provided map as scope. Will return the rendered template or an error
func Render(template string, data map[string]interface{}, functions Functions) (string, error) {
	// let's first find all the template markers
	items := templateFinderRegex.FindAllStringSubmatch(template, -1)
	// for each marker...
	for _, item := range items {
		// matcher is the exact expression, including dollar sign and brackets
		matcher := item[0]
		// expr is what's within the brackets
		expr := item[1]
		// let's walk the path for the expression against the provided data
		if val, err := Walk(expr, data, functions); err == nil {
			// if the value is not nil...
			if val != nil {
				// then we replace the matcher with what we've found.
				// If the value is nil, this won't happen and the matcher will be left untouched
				template = strings.Replace(template, matcher, convertData(val), 1)
			}

		} else {
			// if there was an error, we return it
			return template, err
		}
	}
	// returning the results of our effort
	return template, nil
}

// convertData converts the provided data into a string for the template
func convertData(data interface{}) string {
	switch t := data.(type) {
	case float64:
		return fmt.Sprintf("%f", t)
	case int:
		return fmt.Sprintf("%d", t)
	case bool:
		return strconv.FormatBool(t)
	case []interface{}:
		d, _ := json.Marshal(data)
		return string(d)
	case map[string]interface{}:
		d, _ := json.Marshal(data)
		return string(d)
	case string:
		return t
	default:
		return data.(reflect.Value).String()
	}
}
