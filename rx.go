package gowalker

import (
	"github.com/dlclark/regexp2"
	"regexp"
)

// exprSplitterRegex is a regex to properly split the expressions on dots
var exprSplitterRegex = regexp2.MustCompile("\\.(?![^(]*\\))", regexp2.RE2)

// templateFinderRegex will find the template markers in a string
var templateFinderRegex, _ = regexp.Compile("\\$\\{(.*?)\\}")

// indexExtractorRegex will find the index in an array index accessor
var indexExtractorRegex, _ = regexp.Compile("\\[([0-9]+)\\]")

// functionExtractorRegex will try to extract the function name from a string
var functionExtractorRegex, _ = regexp.Compile("(^.+)(\\((.*)\\))$")

// paramExtractRegex will try to collect and split parameters from a comma separated list of values
var paramExtractRegex, _ = regexp.Compile("([a-zA-Z0-9\\$\\?_\\-\\!|\\/;:\\.\"]|(\\\\,?)*)*")
