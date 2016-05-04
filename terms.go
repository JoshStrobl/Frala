// This file contains functionality for manipulating Terms

package frala

// GetValue gets the value of a language from a Term, if it exists
func GetValue(termName, language string) string {
	if language != "" { // If a language is defined
		language = Sanitize(language) // Ensure it is sanitized
	} else { // If a language is not defined
		language = Config.DefaultLanguage // Set to Default Language
	}

	SetTerm(termName)                                 // Automatically set the term if it doesn't exist already
	value, exists := Config.Terms[termName][language] // Get the language value of this term in Terms

	if !exists { // If the language key/val does not exist
		value = "Term " + termName + " is not translated into " + language
	}

	return value
}

// SetTerm enables you to set a Term to Terms
func SetTerm(termName string) {
	if termName != "" { // If the termName passed isn't empty
		if _, exists := Config.Terms[termName]; !exists { // If the Term doesn't exist already
			Config.Terms[termName] = Term{} // Set this termName to be equivelant to a new Term
		}
	}
}

// SetValue enables you to set the value of a Term language
func SetValue(termName, language, value string) {
	SetTerm(termName)                 // Automatically set the term if it doesn't exist already
	term, _ := Config.Terms[termName] // Get the term if it exists
	language = Sanitize(language)     // Ensure the language is sanitized

	term[language] = value        // Set the value of a particular language to this term
	Config.Terms[termName] = term // Update the Terms
}

// DeleteTerm deletes a Term from Terms
func DeleteTerm(termName string) {
	delete(Config.Terms, termName) // Simply call builtin delete
}

// DeleteValue deletes a language / value from a Term
func DeleteValue(termName, language string) {
	term, exists := Config.Terms[termName] // Get the term if it exists

	if exists { // If the term exists
		delete(term, language)        // Delete from the term the language key/val
		Config.Terms[termName] = term // Update the Terms
	}
}
