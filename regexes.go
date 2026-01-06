package xvalidator

import (
	"regexp"
	"sync"
)

// Regex patterns
const (
	// e164RegexString matches E.164 phone numbers (international format).
	e164RegexString = "^\\+[1-9]?[0-9]{7,14}$"
)

// lazyRegexCompile returns a function that compiles a regex pattern only once using sync.Once.
// This pattern provides thread-safe lazy initialization for regex compilation.
func lazyRegexCompile(pattern string) func() *regexp.Regexp {
	var regex *regexp.Regexp
	var once sync.Once
	return func() *regexp.Regexp {
		once.Do(func() {
			regex = regexp.MustCompile(pattern)
		})
		return regex
	}
}

// Pre-compiled regex functions
var (
	// E164Regex returns a compiled regex for validating E.164 phone numbers.
	E164Regex = lazyRegexCompile(e164RegexString)
)
