package rules

import "testing"

func TestGlobalVariableRegistration(t *testing.T) {
	variablesToFind := []string{
		"number_of_islands_contributing_to_common_pool",
		"number_of_failed_forages",
		"number_of_broken_agreements",
		"max_severity_of_sanctions",
	}

	for _, v := range variablesToFind {
		if _, ok := VariableMap[v]; !ok {
			t.Errorf("Required variable '%v' not found", v)
		}
	}
}
