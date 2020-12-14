package rules

import (
	"gonum.org/v1/gonum/mat"
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

func TestBasicRealValuedRuleEvaluator(t *testing.T) {
	registerNewRealValuedRule()
	result, value, _ := BasicRealValuedRuleEvaluator("Real Test rule")
	if !result && value != 2.0 {
		t.Errorf("Real values rule evaluation error, expected true got '%v', value expected '2' got '%v'", result, value)
	}
}

func registerNewRealValuedRule() {
	//A very contrived rule//
	name := "Real Test rule"
	reqVar := []string{
		"number_of_islands_contributing_to_common_pool",
		"number_of_failed_forages",
		"number_of_broken_agreements",
		"max_severity_of_sanctions",
	}

	v := []float64{1, 0, 0, 0, -4, 0, -1, -1, 0, 2, 0, 0, 0, 0, 2, 0, 0, 1, 0, -1}
	CoreMatrix := mat.NewDense(4, 5, v)
	aux := []float64{1, 1, 3, 0}
	AuxiliaryVector := mat.NewVecDense(4, aux)

	RegisterNewRule(name, reqVar, *CoreMatrix, *AuxiliaryVector)
	// Check internal/clients/team3/client.go for an implementation of a basic evaluator for this rule
}
