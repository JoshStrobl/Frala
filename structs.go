// These are structs used by Frala

package frala

// FralaConfig struct
// This is the configuration for Frala
type FralaConfig struct {
	DefaultLanguage string          // Default Language string, if not declared, default to en
	Terms           map[string]Term // Terms is a map of strings (term names) to individual Terms
}

// Term
// Term is a map[string]string, as each Term has a map of language -> value (where language is a string and value is a string)
type Term map[string]string

// Context
// Context is a struct that has properties relating to the type and type's associated information.
// Used by the line parser.
type Context struct {
	Lang   string `json:"lang"` // Language of the term (if not a fragment)
	Source string `json:"src"`  // Source such as the word or link to fragment
	Type   string `json:"type"` // Type of the Context (fragment or term)
}
