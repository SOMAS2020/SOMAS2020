package rules

import "errors"

type VariableValuePair struct {
	VariableName string
	Multivalued  bool
	SingleValue  float64
	MultiValue   []float64
}

var VariableMap = map[string]VariableValuePair{}

func RegisterNewVariable(pair VariableValuePair) error {
	if _, ok := VariableMap[pair.VariableName]; ok {
		return errors.New("attempted to re-register a variable that had already been registered")
	} else {
		VariableMap[pair.VariableName] = pair
		return nil
	}
}

func UpdateVariable(variableName string, newValue VariableValuePair) error {
	if _, ok := VariableMap[variableName]; ok {
		VariableMap[variableName] = newValue
		return nil
	} else {
		return errors.New("attempted to modify a variable has not been defined")
	}
}
