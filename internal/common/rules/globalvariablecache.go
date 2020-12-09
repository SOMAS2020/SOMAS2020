package rules

import "github.com/pkg/errors"

type VariableValuePair struct {
	VariableName string
	Values       []float64
}

var VariableMap = map[string]VariableValuePair{}

// RegisterNewVariable Registers the provided variable in the global variable cache
func RegisterNewVariable(pair VariableValuePair) error {
	if _, ok := VariableMap[pair.VariableName]; ok {
		return errors.Errorf("attempted to re-register a variable that had already been registered")
	} else {
		VariableMap[pair.VariableName] = pair
		return nil
	}
}

// UpdateVariable Updates variable in global cache with new value
func UpdateVariable(variableName string, newValue VariableValuePair) error {
	if _, ok := VariableMap[variableName]; ok {
		VariableMap[variableName] = newValue
		return nil
	} else {
		return errors.Errorf("attempted to modify a variable has not been defined")
	}
}
