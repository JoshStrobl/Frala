// These are structs used by Frala

package frala

// ConfigOptions is the configuration for Frala
type ConfigOptions struct {
	DefaultLanguage string          // Default Language string, if not declared, default to en
	Languages       []string        // Languages is a list of languages (string)
	Terms           map[string]Term // Terms is a map of strings (term names) to individual Terms
}

// Term is a map[string]string, as each Term has a map of language -> value (where language is a string and value is a string)
type Term map[string]string

// Context is a struct that has properties relating to the type and type's associated information.
// Used by the line parser.
type Context struct {
	Lang   string `json:"lang"` // Language of the term (if not a fragment)
	Source string `json:"src"`  // Source such as the word or link to fragment
	Type   string `json:"type"` // Type of the Context (fragment or term)
}

// ParseResponse is a struct that contains both the content of a file and associated parsing error
type ParseResponse struct {
	Content string // Content of the parsed file
	Error   error  // Error that occurred during parsing
}
