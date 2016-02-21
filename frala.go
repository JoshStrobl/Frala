// This is the main functionality of Frala

package frala

var Config FralaConfig        // Define Config as a Frala Config struct
var CurrentParsingFile string // Define CurrentParsingFile as the file we're currently parsing
var InitError error           // Define InitError as any potential error from initializing Frala

// init
func init() {
	InitError = ReadConfig() // Read the config, setting it's content to Config and any error to readError

	if Config.DefaultLanguage == "" { // If no DefaultLanguage was provided
		Config.DefaultLanguage = "en" // Default language to English
	}
}
