package yogsot

import (
	"errors"

	"gopkg.in/yaml.v2"
)

// Parses a Template file
// Functions that are provided
// Equals
// Variable
// Parameters
func parseTemplate(template []byte) (createStackInput, error) {
	csi := createStackInput{}
	err := yaml.Unmarshal(template, &csi)
	if err != nil {
		return csi, errors.New("error happened while unmarshaling tempalte: " + err.Error())
	}

	return csi, nil
}