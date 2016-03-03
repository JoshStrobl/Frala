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

	language = strings.ToLower(language)       // Lowercase the language
	language = strings.Split(language, "_")[0] // If there is a special format, like en_GB, ensure we only get first part

	if strings.Contains(strings.Join(rtlLanguages, ",")+",", language+",") { // If this is a language specified in rtlLanguages
		direction = "rtl" // Change direction to rtl
	}

	return direction
}
