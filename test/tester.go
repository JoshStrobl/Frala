package main

import (
	"fmt"
	"github.com/JoshStrobl/frala"
)

func main() {
	parsedResponse := frala.Parse("page.html") // Parse page.html

	if parsedResponse.Error == nil { // If there was no parse error, because we're awesome
		fmt.Println(parsedResponse.Content) // Output parsedContent
	} else { // If we failed to parse
		fmt.Println("You fool, you doomed us all! Okay, not really, but here is the error message: ", parsedResponse.Error)
	}
}
