package rules

import "testing"

// TestGlobalRuleRegistration checks whether required rules are registered
func TestGlobalRuleRegistration(t *testing.T) {
	rulesToFind := []string{
		"Kinda Complicated Rule",
	}

	avail, _ := InitialRuleRegistration(false)

	for _, v := range rulesToFind {
		if _, ok := avail[v]; !ok {
			t.Errorf("Required rule '%v' not found", v)
		}
	}
}
