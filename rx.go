package gowalker

import "regexp"

// templateFinderRegex will find the template markers in a string
var templateFinderRegex, _ = regexp.Compile("\\$\\{(.*?)\\}")

// indexExtractorRegex will find the index in an array index accessor
var indexExtractorRegex, _ = regexp.Compile("\\[([0-9]+)\\]")
