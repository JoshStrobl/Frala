// This file contains functionality pertaining to multi-lingual support

package frala

import (
	"strings"
)

// GetDirection gets the likely direction of the language provided
func GetDirection(language string) string {
	direction := "ltr" // Default to Left-to-Right
	rtlLanguages := []string{
		"ar", // Arabic
		"iw", // Hebrew
	}

	language = Sanitize(language)              // Sanitize the language
	language = strings.ToLower(language)       // Lowercase the language
	language = strings.Split(language, "_")[0] // If there is a special format, like en_GB, ensure we only get first part

	if strings.Contains(strings.Join(rtlLanguages, ",")+",", language+",") { // If this is a language specified in rtlLanguages
		direction = "rtl" // Change direction to rtl
	}

	return direction
}

// Sanitize will ensure that language symbols are sanitized correctly
func Sanitize(language string) string {
	filterList := []string{"@"}   // Create a filter list of characters that need to be filtered / replaced
	sanitizeList := []string{"-"} // Create a sanitize list of characters that coorespond with filterList

	for index, filterChar := range filterList {
		language = strings.Replace(language, filterChar, sanitizeList[index], -1) // Set language to replace the filterChar with the cooresponding char in sanitizeList
	}

	return language
}
