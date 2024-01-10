package echosec

import "strings"

// Config is the middleware configuration.
// PathMapping contains a list of validation functions, grouped by path and method.
// DefaultValidation is the default validation action taken if no mapping is matched
type Config struct {
	PathMapping       PathItems
	DefaultValidation ValidationFunc
}

// PathItem is a validation item.
// Patterns is a list of URL patterns to which this validation PathItem responds to
// Methods is a list of mappings based on methods. This can be NIL.
// PathValidation is the default validation for this path, if all Methods validations did not find a match
type PathItem struct {
	Methods        ValidationMap
	Patterns       Patterns
	PathValidation ValidationFunc
}

// PathItems is a collection of PathItem
type PathItems []PathItem

// Patterns ia s list of patterns
type Patterns []string

// ValidationMap maps string keys to validation functions
type ValidationMap map[string]ValidationFunc

// MatchPattern will return true if a path pattern matches the provided path
func (i PathItem) MatchPattern(path string) bool {
	for _, p := range i.Patterns {
		if p == path {
			return true
		}
	}
	return false
}

// FindMethodValidator looks for a method validator that matches the provided method.
// It will return NIL if Methods is NIL or if no method matchers are found
func (i PathItem) FindMethodValidator(method string) ValidationFunc {
	if i.Methods != nil {
		for k, v := range i.Methods {
			for _, m := range strings.Split(k, ",") {
				if strings.ToLower(strings.TrimSpace(m)) == strings.ToLower(method) {
					return v
				}
			}
		}
	}
	return nil
}
