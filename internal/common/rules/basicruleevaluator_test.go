package rules

import (
	"testing"
)

// TestBasicRuleEvaluatorPositive Checks whether rule we expect to evaluate as true actually evaluates as such
func TestBasicRuleEvaluatorPositive(t *testing.T) {
	result, err := BasicBooleanRuleEvaluator("Kinda Complicated Rule")
	if !result {
		t.Errorf("Rule evaluation came as false, when it was expected to be true, potential error with value '%v'", err)
	}
}

// TestBasicRuleEvaluatorNegative Checks whether rule we expect to evaluate as false actually evaluates as such
func TestBasicRuleEvaluatorNegative(t *testing.T) {
	registerTestRule()
	result, err := BasicBooleanRuleEvaluator("Kinda Test Rule")
	if result || err != nil {
		t.Errorf("Rule evaluation came as true, when it was expected to be false, potential error with value '%v'", err)
	}
}
