// This is the main initialization of the Frala Tool

package main

import (
	"fmt"
	"github.com/JoshStrobl/frala"
	"github.com/JoshStrobl/nflag"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// OriginalLanguage is the original value of frala
var OriginalLanguage string

// TargetDirectory is the directory we should save parsed files or Po content to
var TargetDirectory string

// Initialization
func init() {
	nflag.Configure(nflag.ConfigOptions{ShowHelpIfNoArgs: true, ProgramDescription: "Frala CLI tool for file parsing and Po conversion."})

	nflag.Set("convert-terms", nflag.Flag{ // Create the convert-terms flag
		Type:         "bool",
		Descriptor:   "Use convert-terms to declare the action of Terms to Po conversion. Don't pass for Po to Terms conversion.",
		AllowNothing: true,
	})

	nflag.Set("lang", nflag.Flag{ // Create the lang Flag
		Descriptor:   "Language to parse files, Terms to Po conversion, or Po to Terms conversion.",
		Type:         "string",
		DefaultValue: frala.Config.DefaultLanguage, // Set DefaultValue to whatever is already set in frala
		AllowNothing: true,                         // Allow nothing, since we'll just default to DefaultValue
	})

	nflag.Set("parse", nflag.Flag{ // Create the parse flag
		Descriptor: "Files to parse. Accepts comma-separated values.",
		Type:       "string",
	})

	nflag.Set("po", nflag.Flag{ // Create the po flag
		Descriptor: "Po file you wish to convert to Terms or save Terms to.",
		Type:       "string",
	})

	currentWorkingDirectory, getWdErr := os.Getwd()

	if getWdErr != nil { // If there was an error getting the working directory
		absolutePathToDir, absolutePathErr := filepath.Abs(".") // Attempt to get the absolute path of the directory

		if absolutePathErr == nil { // If there was no issue getting the absolute path of the directory
			currentWorkingDirectory = absolutePathToDir // Change currentWorkingDirectory to absolutePathToDir
		} else { // If there was an issue getting the the current working directory
			fmt.Println("Unable to determine the current working directory. Exiting.")
			os.Exit(1)
		}
	}

	nflag.Set("target-dir", nflag.Flag{ // Create the target-dir flag
		Descriptor:   "Target directory to save parsed files or Po content to.",
		DefaultValue: currentWorkingDirectory, // Current working directory
		Type:         "string",
	})
}

func main() {
	nflag.Parse() // Parse nflag

	OriginalLanguage = frala.Config.DefaultLanguage                     // Set OriginalLanguage before doing any parsing or conversion
	frala.Config.DefaultLanguage, _ = nflag.GetAsString("lang")         // Change DefaultLanguage to whatever may be set as lang. This may not change if no value is passed to lang
	frala.Config.Direction = frala.GetDirection(frala.Config.Direction) // Ensure we have an accurate Direction

	parseFiles, parseFilesErr := nflag.GetAsString("parse")
	poFile, poFileErr := nflag.GetAsString("po")

	TargetDirectory, _ = nflag.GetAsString("target-dir")
	TargetDirectory += "/" // Ensure / is appending to end of path

	if (parseFilesErr == nil) && (parseFiles != "") { // If we are parsing files
		ParseFiles(parseFiles) // Call ParseFiles
	} else if (poFileErr == nil) && (poFile != "") { // If we are doing Po conversion
		convertTermsVal, convertTermsGetErr := nflag.GetAsBool("convert-terms") // Get the boolean value of convert-terms

		if convertTermsGetErr != nil { // If there was an error getting convertTermsVal
			convertTermsVal = false // Don't convert from Terms to Po
		}

		PoConversion(convertTermsVal, poFile) // Do Po conversion
	} else { // If we are neither parsing nor doing Po conversion
		nflag.PrintFlags() // Print the nflag flags
	}
}

// ParseFiles will parse the files provided and export parsed to target directory
func ParseFiles(parseFiles string) {
	filesToParse := strings.Split(parseFiles, ",")   // Split the comma-separated parseFiles to filesToParse
	parseResponses := frala.MultiParse(filesToParse) // Parse each file

	for fileName, parseResponse := range parseResponses { // For each parseResponse
		if parseResponse.Error == nil { // If there was no issue parsing this file
			if parseResponse.Content != "" { // If the content is not empty
				fileNameNoPath := filepath.Base(fileName) // Get the base filename
				fmt.Println("Writing content to: " + TargetDirectory + fileNameNoPath)
				os.MkdirAll(TargetDirectory, 0755) // Ensure the directory exists
				ioutil.WriteFile(TargetDirectory+fileNameNoPath, []byte(parseResponse.Content), 0755)
			} else { // If there was no content in the parseResponse.Content
				fmt.Println("No content provided via parsing: " + fileName)
			}
		} else { // If there was an issue parsing this file
			fmt.Println(parseResponse.Error) // Print the error
		}
	}
}

// PoConversion will handle conversion to/from Po
func PoConversion(convertTerms bool, poFile string) {
	if strings.HasSuffix(poFile, ".po") { // If the file is a .po file
		if convertTerms { // If we are converting Terms to Po
			os.MkdirAll(TargetDirectory, 0755)                                    // Ensure the target directory exists
			poFileContent := frala.ConvertToPo(frala.Config.DefaultLanguage)      // Convert the Terms of the language defined (or default language)
			ioutil.WriteFile(TargetDirectory+poFile, []byte(poFileContent), 0755) // Save the Po contents to the poFile in the target directory
		} else { // If we are converting Po to Terms
			conversionError := frala.ConvertFromPo(poFile) // Convert the poFile to Frala Terms

			if conversionError == nil { // If there was no conversion error
				frala.Config.DefaultLanguage = OriginalLanguage // Ensure the default language is maintained after importing from Po
				frala.SaveConfig()                              // Save the config
			} else { // If there was a conversion error
				fmt.Println(conversionError)
			}
		}
	} else {
		fmt.Println(poFile + " does not appear to be a Po file. Please ensure the extension is .po")
	}
}
