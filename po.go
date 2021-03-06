// This file contains functionality for converting to and from gettext PO files.

package frala

import (
	"github.com/robfig/gettext-go/gettext/po" // Support for reading / writing GNU PO files
	"sort"
	"strings"
)

// ConvertFromPo reads a .po file and convert its content to Frala Terms, automatically adding them to the config
func ConvertFromPo(fileName string) error {
	var conversionError error
	var poFile *po.File
	poFile, conversionError = po.Load(fileName)

	if conversionError == nil { // If there was no issue loading the po file
		for _, message := range poFile.Messages { // For each po.Message struc in poFile.Messages
			poLanguage := Sanitize(poFile.MimeHeader.Language)  // Ensure the Po file's MimeHeader Language is sanitized
			SetValue(message.MsgId, poLanguage, message.MsgStr) // Set the msg ID / val as term / value for the language of the file

			if !strings.Contains(strings.Join(Config.Languages, ",")+",", poLanguage+",") { // If the Languages array doesn't contain this Po file lang
				Config.Languages = append(Config.Languages, poLanguage) // Append the language of the Po file to the Languages
				sort.Strings(Config.Languages)                          // Sort the Languages
			}
		}
	}

	return conversionError
}

// ConvertToPo converts Frala Terms into msgid / msgstr context for usage in a .po file
func ConvertToPo(language string) string {
	var poFile po.File
	poFile.MimeHeader.Language = language // Language in MimeHeader as language provided

	for termName := range Config.Terms { // For each termName and term in Terms
		poFile.Messages = append(poFile.Messages, po.Message{MsgId: termName, MsgStr: GetValue(termName, language)}) // Append a new po.Message
	}

	return poFile.String() // Return po format file string
}
