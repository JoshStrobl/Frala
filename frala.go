// This is the main functionality of Frala

package frala

import (
	"fmt"
)

var Config FralaConfig        // Define Config as a Frala Config struct
var CurrentParsingFile string // Define CurrentParsingFile as the file we're currently parsing

// init
func init() {
	readError := ReadConfig() // Read the config, setting it's content to Config and any error to readError

	if readError != nil { // If there was a read error
		fmt.Println(readError) // Output the error
	}

	if Config.DefaultLanguage == "" { // If no DefaultLanguage was provided
		Config.DefaultLanguage = "en" // Default language to English
	}
}
