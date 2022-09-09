package gowalker

import (
	"context"
	"errors"
	"strings"
)

// Render renders a template, using the provided map as scope. Will return the rendered template or an error
func Render(ctx context.Context, template string, data any, functions *Functions) (string, error) {
	if deadlineMet(ctx) {
		return "", errors.New("deadline exceeded")
	}
	if hasCancelled(ctx) {
		return "", errors.New("cancelled")
	}
	// let's first find all the template markers
	items := templateFinderRegex.FindAllStringSubmatch(template, -1)
	// for each marker...
	for _, item := range items {
		// matcher is the exact expression, including dollar sign and brackets
		matcher := item[0]
		// expr is what's within the brackets
		expr := item[1]
		// let's walk the path for the expression against the provided data
		if val, err := Walk(ctx, expr, data, functions); err == nil {
			// if the value is not nil...
			if val != nil {
				// then we replace the matcher with what we've found.
				// If the value is nil, this won't happen and the matcher will be left untouched
				template = strings.Replace(template, matcher, convertDataToString(val), 1)
			}

		} else {
			// if there was an error, we return it
			return template, err
		}
	}
	// returning the results of our effort
	return template, nil
}

// RenderAll will render the provided templates, making subTemplates available for complex rendering
func RenderAll(ctx context.Context, template string, subTemplates SubTemplates, data any, functions *Functions) (string, error) {
	if subTemplates == nil {
		subTemplates = NewSubTemplates()
	}
	for k, v := range subTemplates {
		functions.functionScope["_"+k] = v
	}
	return Render(ctx, template, data, functions)
}
