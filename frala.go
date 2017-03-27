// This is the main functionality of Frala

package frala

// Config is the configuration of Frala
var Config ConfigOptions

// InitError is any potential error from initializing Frala
var InitError error

// CurrentParsingFile is the file we're currently parsing
var CurrentParsingFile string

func init() {
	InitError = ReadConfig() // Read the config, setting it's content to Config and any error to readError

	if Config.DefaultLanguage == "" { // If no DefaultLanguage was provided
		Config.DefaultLanguage = "en" // Default language to English
	}

	Config.DefaultLanguage = Sanitize(Config.DefaultLanguage) // Sanitize the language
	Config.CurrentLanguage = Config.DefaultLanguage // Set CurrentLanguage to default to DefaultLanguage
	Config.Direction = GetDirection(Config.DefaultLanguage)   // Get the likely direction of the DefaultLanguage
}
