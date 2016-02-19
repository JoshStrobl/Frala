// This file contains functionality for manipulating the frala Config or frala.json file.

package frala

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

// ReadConfig
// This function will read any frala.json file and update the Config
func ReadConfig() error {
	var configContent []byte
	var readError error

	configContent, readError = ioutil.ReadFile("frala.json") // Read the frala.json configuration

	if readError == nil { // If there was no read error
		readError = json.Unmarshal(configContent, &Config) // Decode configContent into Config

		if readError != nil {
			readError = errors.New("Unable to decode frala.json into the appropriate Frala configuration structure. Please verify the correctness of your config.")
		}
	} else { // If there was a read error
		readError = errors.New("frala.json file does not exist.")
	}

	return readError
}

// SaveConfig
// This function will save the Config to frala.json
func SaveConfig() error {
	var configContent []byte
	var saveError error

	configContent, saveError = json.Marshal(Config) // Encode the Config into configContent. If encoding fails, set saveError

	if saveError == nil { // If there was no error encoding
		saveError = ioutil.WriteFile("frala.json", configContent, 0755) // Attempt to write the configContent to frala.json

		if saveError != nil {
			saveError = errors.New("Failed to save the Config to frala.json.")
		}
	} else { // If we failed to encode the Config to JSON
		saveError = errors.New("Failed to encode the Config to JSON.")
	}

	return saveError
}
