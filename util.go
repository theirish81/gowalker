package gowalker

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
)

// convertDataToString converts the provided data into a string for the template
func convertDataToString(data any) string {
	if data == nil {
		return "null"
	}
	switch reflect.TypeOf(data).Kind() {
	case reflect.Int, reflect.Int64:
		return fmt.Sprintf("%d", data)
	case reflect.Float32:
		return strconv.FormatFloat(float64(data.(float32)), 'f', -1, 64)
	case reflect.Float64:
		return strconv.FormatFloat(data.(float64), 'f', -1, 64)
	case reflect.Bool:
		return strconv.FormatBool(data.(bool))
	case reflect.Slice, reflect.Map, reflect.Struct:
		// Slices and maps are rendered as JSON
		d, _ := json.Marshal(data)
		return string(d)
	case reflect.Pointer:
		t := reflect.ValueOf(data).Elem().Interface()
		return convertDataToString(t)
	}
	return data.(string)
}

// convertStringToSameType tries to convert val to the same type of sample
func convertStringToSameType(sample any, val string) (any, error) {
	if sample == nil {
		return val, errors.New("sample is nil")
	}
	switch reflect.TypeOf(sample).Kind() {
	case reflect.Int, reflect.Int64:
		return strconv.Atoi(val)
	case reflect.Float32:
		return strconv.ParseFloat(val, 32)
	case reflect.Float64:
		return strconv.ParseFloat(val, 64)
	case reflect.Bool:
		return strconv.ParseBool(val)
	default:
		return val, nil
	}
}

func LoadTemplatesFromDisk(filePath string) (string, SubTemplates, error) {
	data, err := os.ReadFile(filePath)
	template := string(data)
	if err != nil {
		return "", nil, err
	}

	templDir := path.Dir(filePath)
	subTemplates := NewSubTemplates()
	files, err := os.ReadDir(templDir)
	if err != nil {
		return "", nil, err
	}
	rootTemplateName := filepath.Base(filePath)
	for _, file := range files {
		if !file.IsDir() && !strings.HasPrefix(file.Name(), ".") && file.Name() != rootTemplateName {
			data, _ = os.ReadFile(path.Join(templDir, file.Name()))
			subTemplateName := file.Name()[0:strings.LastIndex(file.Name(), ".")]
			subTemplates[subTemplateName] = string(data)
		}
	}
	return template, subTemplates, nil
}
