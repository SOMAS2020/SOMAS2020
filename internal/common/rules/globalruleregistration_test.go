package rules

import "testing"

func TestGlobalRuleRegistration(t *testing.T) {
	rulesToFind := []string{
		"Kinda Complicated Rule",
	}

	for _, v := range rulesToFind {
		if _, ok := AvailableRules[v]; !ok {
			t.Errorf("Required rule '%v' not found", v)
		}
	}
}
