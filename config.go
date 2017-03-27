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
	if configContent, readErr := ioutil.ReadFile("frala.json"); readErr == nil { // Read the frala.json configuration, set any err to readErr
		decodeErr := json.Unmarshal(configContent, &Config)

		if decodeErr != nil { // Decode configContent into Config
			decodeErr = errors.New("Unable to decode frala.json: " + decodeErr.Error())
		}

		return decodeErr
	} else { // If there was a read error
		return errors.New("Failed to read frala.json: " + readErr.Error())
	}
}

// SaveConfig saves the Config to frala.json
func SaveConfig() error {
	if configContent, encodeErr := json.MarshalIndent(Config, "", "\t"); encodeErr == nil { // Encode the Config into configContent, ensure it maintains pretty formatting.
		return coreutils.WriteOrUpdateFile("frala.json", configContent, coreutils.NonGlobalFileMode) // Attempt to write the configContent to frala.json
	} else { // If we failed to encode the Config to JSON
		return errors.New("Failed to encode the Config to JSON: " + encodeErr.Error())
	}
}
