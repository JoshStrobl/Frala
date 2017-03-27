// This file contains functionality responsible for parsing a file or syntax

package frala

import (
	"errors"
	"github.com/StroblIndustries/coreutils"
	"io/ioutil"
	"strings"
)

// MultiParse
// Parses multiple files provided and return a map of ParseResponses
func MultiParse(files []string) []ParseResponse {
	var parseResponses []ParseResponse

	for _, file := range files { // For each file
		if strings.HasSuffix(file, ".html") { // If this is an HTML file
			parseResponses = append(parseResponses, Parse(file))
		}
	}

	return parseResponses
}

// MultilingualParse
// Parses all provided files using all available languages
// Returns a map of languages cooresponding to an array of ParseResponse
func MultilingualParse(files []string) map[string][]ParseResponse {
	parserResponses := make(map[string][]ParseResponse)

	for _, lang := range Config.Languages { // For each of our languages
		Config.CurrentLanguage = lang
		parserResponses[lang] = MultiParse(files)
	}

	return parserResponses
}

// Parse
// Parses a file
func Parse(file string) ParseResponse {
	var contentBytes []byte
	var parseResponse ParseResponse

	parseResponse.Name = file
	contentBytes, parseResponse.Error = ioutil.ReadFile(file) // Read the file, set any read error to parseResponse.Error

	if len(contentBytes) == 0 || parseResponse.Error != nil {
		parseResponse.Error = errors.New("Failed to read: " + file)
		return parseResponse
	}

	content := string(contentBytes[:])
	lines := strings.Split(content, "\n") // Split the string up into lines
	newLines := []string{} // Create a newLines slice that'll contain our content

	for _, line := range lines { // For each line
		if strings.Contains(line, "{{") && strings.Contains(line, "}}") { // Frala syntax
			var newLineContent string // Define newLineContent as the new line content we will return

			multiSyntaxSplit := strings.Split(line, "{{") // Split the line into segments based on {{

			for _, lineSegment := range multiSyntaxSplit { // For each segment of the line
				var parsedContext string
				var context Context

				syntaxEndSplit := strings.Split(lineSegment, "}}")                              // Split the segment based on the ending of the Frala syntax
				fralaSyntaxContent := strings.TrimSpace(syntaxEndSplit[0])
				fralaSyntaxProperties := strings.Split(fralaSyntaxContent, " ") // Split on whitespace

				for _, property := range fralaSyntaxProperties {
					propertyValueSplit := strings.Split(strings.Replace(property, "\"", "", -1), "=")
					property := propertyValueSplit[0]
					value := propertyValueSplit[1]

					if property == "lang" {
						context.Lang = value
					} else if property == "src" {
						context.Source = value
					} else if property == "type" {
						context.Type = value
					}
				}

				parsedContext += context.Parse()
				contentAfterFralaSyntax := syntaxEndSplit[1]

				newLineContent += parsedContext + contentAfterFralaSyntax
			}

			newLines = append(newLines, newLineContent)
		} else { // Not Frala syntax
			newLines = append(newLines, line)
		}
	}

	parseResponse.Content = strings.Join(newLines, "\n")
	return parseResponse
}

// Context Parse
// Parses a Frala context and returns a string
func (c *Context) Parse() string {
	if c.Lang == "" && c.Type == "term" { // If Lang isn't set for term
		c.Lang = Config.CurrentLanguage // Set to the current parsing language.
	} else if c.Source == "" { // No Source set
		return "Source for this Frala syntax not specified."
	} else if c.Type == "" { // No Type set
		return "Type for this Frala syntax not specified."
	}

	var parsedContext string

	if (c.Type == "fragment") { // If this is a Fragment
		if c.Source != CurrentParsingFile { // If we're not doing some crazy import fragment within itself sorcery
			restoreFileName := CurrentParsingFile                                                              // Set restoreFileName to CurrentParsingFile before doing any potential crazy business
			c.Source = coreutils.AbsPath(CurrentParsingFile) + c.Source

			fragmentParserResponse := Parse(c.Source) // Pass the Fragment file path to our Parse
			CurrentParsingFile = restoreFileName                 // Restore file name back to original state

			if fragmentParserResponse.Error == nil { // If there was no error reading the fragment file
				parsedContext = fragmentParserResponse.Content // Set parsedContext to fragment ParserResponse Content
			} else { // If the fragment file does not exist
				parsedContext = fragmentParserResponse.Error.Error() // Get the parser response error
			}
		} else { // If we're attempting Fragment inception
			parsedContext = "Cannot import " + c.Source + " within itself."
		}
	} else if (c.Type == "term")  { // If this is a term
		switch (c.Source) {
			case "frala.CurrentLanguage":
				parsedContext = Config.CurrentLanguage
				break;
			case "frala.DefaultLanguage":
				parsedContext = Config.DefaultLanguage
				break;
			case "frala.Direction":
				parsedContext = Config.Direction
				break;
			case "frala.Languages":
				if len(Config.Languages) != 0 { // If there was languages defined in the Config
					parsedContext = strings.Join(Config.Languages, ",")
				} else { // If there are no languages defined in the Config.Languages
					parsedContext = Config.DefaultLanguage // Return the DefaultLanguage instead
				}
				break;
			default:
				parsedContext = GetValue(c.Source, c.Lang) // Get the Language value of this Source in Terms
				break;
		}
	} else {
		parsedContext = c.Type + " is not a valid type."
	}

	return parsedContext
}