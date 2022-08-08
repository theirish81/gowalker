package gowalker

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// Render renders a template, using the provided map as scope. Will return the rendered template or an error
func Render(template string, data any, functions *Functions) (string, error) {
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

func RenderAll(template string, subTemplate map[string]string, data map[string]interface{}, functions *Functions) (string, error) {
	for k, v := range subTemplate {
		functions.functionScope["_"+k] = v
	}
	return Render(template, data, functions)
}

// convertData converts the provided data into a string for the template
func convertData(data interface{}) string {
	switch reflect.TypeOf(data).Kind() {
	case reflect.Float64:
		return fmt.Sprintf("%f", data)
	case reflect.Int:
		return fmt.Sprintf("%d", data)
	case reflect.Bool:
		return strconv.FormatBool(data.(bool))
	case reflect.Slice, reflect.Map:
		d, _ := json.Marshal(data)
		return string(d)
	}
	return data.(string)
}
