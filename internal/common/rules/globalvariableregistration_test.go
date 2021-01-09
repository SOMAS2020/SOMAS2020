package rules

import "testing"

// TestGlobalVariableRegistration checks whether global cache contains all required variable
func TestGlobalVariableRegistration(t *testing.T) {
	variablesToFind := []VariableFieldName{
		NumberOfIslandsContributingToCommonPool,
		NumberOfFailedForages,
		NumberOfBrokenAgreements,
		MaxSeverityOfSanctions,
	}

	for _, v := range variablesToFind {
		if _, ok := InitialVarRegistration()[v]; !ok {
			t.Errorf("Required variable '%v' not found", v)
		}
	}
}
