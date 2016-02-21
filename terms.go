// This file contains functionality for manipulating Terms

package frala

// GetValue
// This function will get the value of a language from a Term, if it exists
func GetValue(termName, language string) string {
	if language == "" { // If no language is defined
		language = Config.DefaultLanguage // Set to Default Language
	}

	SetTerm(termName)                                 // Automatically set the term if it doesn't exist already
	value, exists := Config.Terms[termName][language] // Get the language value of this term in Terms

	if !exists { // If the language key/val does not exist
		value = "Term " + termName + " is not translated into " + language
	}

	return value
}

// SetTerm
// This function will enable you to set a Term to Terms
func SetTerm(termName string) {
	if termName != "" { // If the termName passed isn't empty
		if _, exists := Config.Terms[termName]; !exists { // If the Term doesn't exist already
			Config.Terms[termName] = Term{} // Set this termName to be equivelant to a new Term
		}
	}
}

// SetValue
// This function will enable you to set the value of a Term language
func SetValue(termName, language, value string) {
	SetTerm(termName)                 // Automatically set the term if it doesn't exist already
	term, _ := Config.Terms[termName] // Get the term if it exists

	term[language] = value        // Set the value of a particular language to this term
	Config.Terms[termName] = term // Update the Terms
}

// DeleteTerm
// This function will delete a Term from Terms
func DeleteTerm(termName string) {
	delete(Config.Terms, termName) // Simply call builtin delete
}

// DeleteValue
// This function will delete a language / value from a Term
func DeleteValue(termName, language string) {
	term, exists := Config.Terms[termName] // Get the term if it exists

	if exists { // If the term exists
		delete(term, language)        // Delete from the term the language key/val
		Config.Terms[termName] = term // Update the Terms
	}
}