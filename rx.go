package gowalker

import "regexp"

var templateFinderRegex, _ = regexp.Compile("\\$\\{(.*?)\\}")
var indexExtractorRegex, _ = regexp.Compile("\\[([0-9]+)\\]")
