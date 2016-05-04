// This file contains functionality responsible for parsing a file or syntax

package frala

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// MultiParse parses multiple files provided and return a map of ParseResponses
func MultiParse(files []string) map[string]ParseResponse {
	parserResponses := make(map[string]ParseResponse)

	for _, file := range files { // For each file
		if strings.HasSuffix(file, ".html") { // If this is an HTML file
			parserResponses[file] = Parse(file) // Define this file in parserResponses to ParserResponse provided by Parse
		}
	}

	return parserResponses
}

// Parse parses a file provided and return a ParseResponse
func Parse(file string) ParseResponse {
	parserResponse := ParseResponse{}
	fileContentBytes, fileContentError := ioutil.ReadFile(file) // Read the file content and push error to fileContentError

	if fileContentError == nil { /// If there was no error reading the file
		var parsedString string

		CurrentParsingFile = file // Set CurrentParsingFile to this file
		fileContent := string(fileContentBytes[:])
		fileSplitLines := strings.Split(fileContent, "\n") // Split by new line

		for _, line := range fileSplitLines { // For each line
			parsedString += ParseLine(line) + "\n"
		}

		parserResponse.Content = parsedString
	} else {
		parserResponse.Error = errors.New(file + " does not exist.")
	}

	return parserResponse
}

// ParseLine parses an individual line
func ParseLine(lineContent string) string {
	parsedLineContent := lineContent                                                // Default parsedLineContent as the lineContent
	if strings.Contains(lineContent, "{{") && strings.Contains(lineContent, "}}") { // If this has Frala syntax in it (has both { and })
		var newLineContent string // Define newLineContent as the new line content we will return

		multiSyntaxSplit := strings.Split(lineContent, "{{") // Split the lineContent into segments based on {

		for _, lineSegment := range multiSyntaxSplit { // For each segment of the line
			if strings.Contains(lineSegment, "}}") { // If this segment of the line contains }
				syntaxEndSplit := strings.Split(lineSegment, "}}")                              // Split the segment based on the ending of the Frala syntax
				parsedSyntax := ParseSyntax("{{" + strings.TrimSpace(syntaxEndSplit[0]) + "}}") // Parse the syntax and return it
				contentAfterFralaSyntax := syntaxEndSplit[1]                                    // Set Frala

				newLineContent += parsedSyntax + contentAfterFralaSyntax
			} else { // If this segment does not contain an end-syntax, meaning it is likely a segment prior to the syntax
				newLineContent += lineSegment // Add the lineSegment to the newLineContent
			}
		}

		parsedLineContent = newLineContent
	}

	return parsedLineContent
}

// ParseSyntax parses a Frala syntax string and return the appropriate (if any) associated HTML content or term
func ParseSyntax(fralaSyntax string) string {
	parsedString := fralaSyntax // Default parsedString to fralaSyntax

	var fralaContext Context // Define fralaContext as a Context struct

	// #region Convert Frala syntax to JSON

	fralaSyntaxJSON := fralaSyntax                                                      // Initially set fralaSyntaxJSON to fralaSyntax
	searchStrings := []string{"{{", "}}", "\" ", "=", "lang", "src", "type"}            // Array of things we need to search for
	replaceStrings := []string{"{", "}", "\",", ":", "\"lang\"", "\"src\"", "\"type\""} // Array of things we'll replace

	for pos, searchString := range searchStrings { // For each searchString in searchStrings
		fralaSyntaxJSON = strings.Replace(fralaSyntaxJSON, searchString, replaceStrings[pos], -1) // Replace the searchString with the cooresponding replaceString
	}

	// #endregion

	decodeErr := json.Unmarshal([]byte(fralaSyntaxJSON), &fralaContext) // Decode fralaSyntaxJSON into fralaContext

	if decodeErr == nil { // If there was no decode error
		if (fralaContext.Type == "") { // If no type was defined
			return "No type provided for this Frala context."
		} else if (fralaContext.Source == "") { // If no source was provided
			return "No source provided for this Frala context."
		}

		if (fralaContext.Type == "fragment") { // If this is a Fragment
			if fralaContext.Source != CurrentParsingFile { // If we're not doing some crazy import fragment within itself sorcery
				restoreFileName := CurrentParsingFile                                                              // Set restoreFileName to CurrentParsingFile before doing any potential crazy business
				fralaContext.Source = filepath.Clean(filepath.Dir(CurrentParsingFile) + "/" + fralaContext.Source) // Ensure we have prepend the directory of the current parsing file

				fragmentParserResponse := Parse(fralaContext.Source) // Attempt to read the fragment
				CurrentParsingFile = restoreFileName                 // Restore file name back to original state

				if fragmentParserResponse.Error == nil { // If there was no error reading the fragment file
					parsedString = fragmentParserResponse.Content // Set parsedString to fragment ParserResponse Content
				} else { // If the fragment file does not exist
					parsedString = fragmentParserResponse.Error.Error() // Get the parser response error
				}
			} else { // If we're attempting Fragment inception
				parsedString = "I can't do that Dave. (Importing Fragment within itself)"
			}
		} else if (fralaContext.Type == "term")  { // If this is a term
			if strings.HasPrefix(fralaContext.Source, "frala.") { // If we are actually fetching an option from Frala
				switch (fralaContext.Source) {
					case "frala.DefaultLanguage":
						parsedString = Config.DefaultLanguage
						break;
					case "frala.Direction":
						parsedString = Config.Direction
						break;
					case "frala.Languages":
						if len(Config.Languages) != 0 { // If there was languages defined in the Config
							parsedString = strings.Join(Config.Languages, ",")
						} else { // If there are no languages defined in the Config.Languages
							parsedString = Config.DefaultLanguage // Return the DefaultLanguage instead
						}
						break;
					default:
						parsedString = fralaContext.Source + " is not a valid Frala Built-in term."
						break;
				}
			} else { // If this is a "normal" term
				parsedString = GetValue(fralaContext.Source, fralaContext.Lang) // Get the Language value of this Source in Terms
			}
		}
	}

	return parsedString
}
