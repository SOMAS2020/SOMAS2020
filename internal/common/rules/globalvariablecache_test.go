package rules

import (
	"testing"
)

func TestRegisterNewVariable(t *testing.T) {
	registerTestVariable()
	if val, ok := VariableMap["Test variable"]; !ok {
		t.Errorf("Global variable map unable to register new variables")
	} else {
		if val.SingleValue != 5 {
			t.Errorf("Global variable map didn't register correct value for variable, wanted '%v' got '%v'", 5, val.SingleValue)
		}
	}
}

func registerTestVariable() {
	pair := VariableValuePair{
		VariableName: "Test variable",
		Multivalued:  false,
		SingleValue:  5,
		MultiValue:   nil,
	}
	_ = RegisterNewVariable(pair)
}
