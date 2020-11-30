package rules

import (
	"testing"
)

func TestRegisterNewVariable(t *testing.T) {
	registerTestVariable()
	if val, ok := VariableMap["Test variable"]; !ok {
		t.Errorf("Global variable map unable to register new variables")
	} else {
		if val.Values[0] != 5 {
			t.Errorf("Global variable map didn't register correct value for variable, wanted '%v' got '%v'", 5, val.Values[0])
		}
	}
}

func registerTestVariable() {
	pair := VariableValuePair{
		VariableName: "Test variable",
		Values:       []float64{5},
	}
	_ = RegisterNewVariable(pair)
}
