package rules

import (
	"testing"
)

// TestRegisterNewVariable tests whther the global variable cache can register new values
func TestRegisterNewVariable(t *testing.T) {
	VariableMapTesting := generateTestVariableStore()
	registerTestVariable(VariableMapTesting)
	if val, ok := VariableMapTesting["Test variable"]; !ok {
		t.Errorf("Global variable map unable to register new variables")
	} else {
		if val.Values[0] != 5 {
			t.Errorf("Global variable map didn't register correct value for variable, wanted '%v' got '%v'", 5, val.Values[0])
		}
	}
}

func generateTestVariableStore() map[string]VariableValuePair {
	return map[string]VariableValuePair{}
}

func registerTestVariable(variableStore map[string]VariableValuePair) {
	pair := VariableValuePair{
		VariableName: "Test variable",
		Values:       []float64{5},
	}
	_ = registerNewVariableInternal(pair, variableStore)
}
