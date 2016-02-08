// This is the main functionality of Frala

package frala

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var Config FralaConfig        // Define Config as a Frala Config struct
var CurrentParsingFile string // Define CurrentParsingFile as the file we're currently parsing

// init
func init() {
	configContent, configReadError := ioutil.ReadFile("frala.json") // Read the frala.json configuration

	if configReadError == nil { // If there was no read error
		unmarshalError := json.Unmarshal(configContent, &Config) // Decode configContent into Config

		if unmarshalError != nil {
			fmt.Println("Unable to decode frala.json into the appropriate Frala configuration structure. Please verify the correctness of your config.")
			os.Exit(1) // Die
		}
	}

	if Config.DefaultLanguage == "" { // If no DefaultLanguage was provided
		Config.DefaultLanguage = "en" // Default language to English
	}
}

// Parse
// This function will parse a file provided and return either parsed contents or an error
func Parse(file string) (string, error) {
	var parsedString string
	var parseError error
	fileContentBytes, fileContentError := ioutil.ReadFile(file) // Read the file content and push error to fileContentError

	if fileContentError == nil { /// If there was no error reading the file
		CurrentParsingFile = file // Set CurrentParsingFile to this file
		fileContent := string(fileContentBytes[:])
		fileSplitLines := strings.Split(fileContent, "\n") // Split by new line

		for _, line := range fileSplitLines { // For each line
			parsedLineContent := line // Default parsedLineContent to existing line content

			if strings.Contains(line, "{") { // If the string contains Frala syntax
				parsedLineContent = ParseLine(line) // Parse the line content
			}

			parsedString += parsedLineContent + "\n"
		}
	} else {
		parseError = errors.New("File does not exist.")
	}

	return parsedString, parseError
}

// ParseLine
// This function will parse an individual line
func ParseLine(lineContent string) string {
	parsedLineContent := lineContent        // Default parsedLineContent as the lineContent
	if strings.Contains(lineContent, "{") { // If the string contains Frala syntax
		var newLineContent string // Define newLineContent as the new line content we will return

		multiSyntaxSplit := strings.Split(lineContent, "{") // Split the lineContent into segments based on {

		for _, lineSegment := range multiSyntaxSplit { // For each segment of the line
			if strings.Contains(lineSegment, "}") { // If this segment of the line contains }
				syntaxEndSplit := strings.Split(lineSegment, "}")                             // Split the segment based on the ending of the Frala syntax
				parsedSyntax := ParseSyntax("{" + strings.TrimSpace(syntaxEndSplit[0]) + "}") // Parse the syntax and return it
				contentAfterFralaSyntax := syntaxEndSplit[1]                                  // Set Frala

				newLineContent += parsedSyntax + contentAfterFralaSyntax
			} else { // If this segment does not contain an end-syntax, meaning it is likely a segment prior to the syntax
				newLineContent += lineSegment // Add the lineSegment to the newLineContent
			}
		}

		parsedLineContent = newLineContent
	}

	return parsedLineContent
}

// ParseSyntax
// This function will parse a Frala syntax string and return the appropriate (if any) associated HTML content or term
func ParseSyntax(fralaSyntax string) string {
	parsedString := fralaSyntax // Default parsedString to fralaSyntax

	var fralaContext Context // Define fralaContext as a Context struct

	// #region Convert Frala syntax to JSON

	fralaSyntaxJSON := fralaSyntax                                          // Initially set fralaSyntaxJSON to fralaSyntax
	searchStrings := []string{" ", "=", "lang", "src", "type"}              // Array of things we need to search for
	replaceStrings := []string{",", ":", "\"lang\"", "\"src\"", "\"type\""} // Array of things we'll replace

	for pos, searchString := range searchStrings { // For each searchString in searchStrings
		fralaSyntaxJSON = strings.Replace(fralaSyntaxJSON, searchString, replaceStrings[pos], -1) // Replace the searchString with the cooresponding replaceString
	}

	// #endregion

	decodeErr := json.Unmarshal([]byte(fralaSyntaxJSON), &fralaContext) // Decode fralaSyntaxJSON into fralaContext

	if decodeErr == nil { // If there was no decode error
		if (fralaContext.Type == "fragment") && (fralaContext.Source != "") { // If this is a Fragment
			if fralaContext.Source != CurrentParsingFile { // If we're not doing some crazy import fragment within itself sorcery
				restoreFileName := CurrentParsingFile // Set restoreFileName to CurrentParsingFile before doing any potential crazy business

				fragmentContentBytes, fragmentContentIOErr := Parse(fralaContext.Source) // Attempt to read the fragment
				CurrentParsingFile = restoreFileName                                     // Restore file name back to original state

				if fragmentContentIOErr == nil { // If there was no error reading the fragment file
					parsedString = string(fragmentContentBytes[:]) // Set parsedString to fragment file content
				} else { // If the fragment file does not exist
					parsedString = fralaContext.Source + " does not exist."
				}
			} else { // If we're attempting Fragment inception
				parsedString = "I can't do that Dave. (Importing Fragment within itself)"
			}
		} else if (fralaContext.Type == "term") && (fralaContext.Source != "") { // If this is a term
			var termExists bool          // Define termExists as a bool during our Config.Terms key/val check
			if fralaContext.Lang == "" { // If no language is defined
				fralaContext.Lang = Config.DefaultLanguage // Set to Default Language
			}

			parsedString, termExists = Config.Terms[fralaContext.Source][fralaContext.Lang] // Get the Language value of this Source in Terms

			if !termExists { // If the term does not exist
				parsedString = "Term " + fralaContext.Source + " is not translated into " + fralaContext.Lang // Change parsedString to err string
			}
		} else { // If the necessary syntax elements were not provided
			parsedString = "Necessary Frala Syntax elements were not provided for this context."
		}
	}

	return parsedString
}
