package rules

import (
	"testing"
)

// TestRegisterNewVariable tests whether the global variable cache can register new values
func TestRegisterNewVariable(t *testing.T) {
	VariableMapTesting := generateTestVariableStore()
	registerTestVariable(VariableMapTesting)
	if val, ok := VariableMapTesting[TestVariable]; !ok {
		t.Errorf("Global variable map unable to register new variables")
	} else {
		if val.Values[0] != 5 {
			t.Errorf("Global variable map didn't register correct value for variable, wanted '%v' got '%v'", 5, val.Values[0])
		}
	}
}

func generateTestVariableStore() map[VariableFieldName]VariableValuePair {
	return map[VariableFieldName]VariableValuePair{}
}

func registerTestVariable(variableStore map[VariableFieldName]VariableValuePair) {
	pair := VariableValuePair{
		VariableName: TestVariable,
		Values:       []float64{5},
	}
	_ = RegisterNewVariableInternal(pair, variableStore)
}
