package main

import (
	"fmt"
	"github.com/JoshStrobl/frala"
	"io/ioutil"
	"os"
)

func main() {
	fmt.Println("Testing Frala Parsing")
	parsedResponse := frala.Parse("page.html") // Parse page.html

	if parsedResponse.Error == nil { // If there was no parse error, because we're awesome
		fmt.Println(parsedResponse.Content) // Output parsedContent
	} else { // If we failed to parse
		fmt.Println("You fool, you doomed us all! Okay, not really, but here is the error message: ", parsedResponse.Error)
	}

	fmt.Println("\nTesting Po functionality")

	fmt.Println("Running ConvertToPo for lang: en")
	fmt.Println(frala.ConvertToPo("en"))

	fmt.Println("\nRunning ConvertToPo for lang: fi")
	finnishPoContent := frala.ConvertToPo("fi")
	fmt.Println(finnishPoContent)

	if len(finnishPoContent) != 0 {
		ioutil.WriteFile("finnish.po", []byte(finnishPoContent), 0777) // Write the contents to finnish.po
		fmt.Println("\nRunning ConvertFromPo")
		convertError := frala.ConvertFromPo("finnish.po")

		if convertError == nil { // If there was no conversion error
			fmt.Println("No conversion error. Terms below:\n")
			fmt.Println(frala.Config.Terms)
		} else { // Error converting from PO
			fmt.Println("There was an issue converting from Po: ", convertError)
		}

		os.Remove("finnish.po")
	} else { // No convent in finnishPoContent
		fmt.Println("No content from ConvertToPo for Finnish.")
	}
}
