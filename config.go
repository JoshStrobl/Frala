// This file contains functionality for manipulating the frala Config or frala.json file.

package frala

import (
	"encoding/json"
	"errors"
	"github.com/StroblIndustries/coreutils"
	"io/ioutil"
)

// ReadConfig reads any frala.json file and update the Config
func ReadConfig() error {
	var configContent []byte
	var readError error

	configContent, readError = ioutil.ReadFile("frala.json") // Read the frala.json configuration

	if readError == nil { // If there was no read error
		readError = json.Unmarshal(configContent, &Config) // Decode configContent into Config

		if readError != nil {
			readError = errors.New("Unable to decode frala.json: " + readError.Error())
		}
	} else { // If there was a read error
		readError = errors.New("frala.json file does not exist")
	}

	return readError
}

// SaveConfig saves the Config to frala.json
func SaveConfig() error {
	var configContent []byte
	var saveError error

	configContent, saveError = json.MarshalIndent(Config, "", "\t") // Encode the Config into configContent, ensure it maintains pretty formatting. If encoding fails, set saveError

	if saveError == nil { // If there was no error encoding
		saveError = coreutils.WriteOrUpdateFile("frala.json", configContent, coreutils.NonGlobalFileMode) // Attempt to write the configContent to frala.json
	} else { // If we failed to encode the Config to JSON
		saveError = errors.New("Failed to encode the Config to JSON: " + saveError.Error())
	}

	return saveError
}
