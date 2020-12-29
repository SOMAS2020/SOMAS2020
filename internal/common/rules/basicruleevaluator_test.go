package rules

import (
	"testing"

	"gonum.org/v1/gonum/mat"
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
	registerTestRule(AvailableRules)
	result, err := BasicBooleanRuleEvaluator("Kinda Test Rule")
	if result || err != nil {
		t.Errorf("Rule evaluation came as true, when it was expected to be false, potential error with value '%v'", err)
	}
}

func TestBasicRealValuedRuleEvaluator(t *testing.T) {
	registerNewRealValuedRule(t)
	result, value, _ := BasicRealValuedRuleEvaluator("Real Test rule")
	if !result && value != 2.0 {
		t.Errorf("Real values rule evaluation error, expected true got '%v', value expected '2' got '%v'", result, value)
	}
}

func registerNewRealValuedRule(t *testing.T) {
	//A very contrived rule//
	name := "Real Test rule"
	reqVar := []VariableFieldName{
		NumberOfIslandsContributingToCommonPool,
		NumberOfFailedForages,
		NumberOfBrokenAgreements,
		MaxSeverityOfSanctions,
	}

	v := []float64{1, 0, 0, 0, -4, 0, -1, -1, 0, 2, 0, 0, 0, 0, 2, 0, 0, 1, 0, -1}
	CoreMatrix := mat.NewDense(4, 5, v)
	aux := []float64{1, 1, 3, 0}
	AuxiliaryVector := mat.NewVecDense(4, aux)

	_, ruleError := RegisterNewRule(name, reqVar, *CoreMatrix, *AuxiliaryVector, false)
	if ruleError != nil {
		t.Errorf("Problem with registering new real valued rule in test, error message : '%v'", ruleError.Error())
	}
	// Check internal/clients/team3/client.go for an implementation of a basic evaluator for this rule
}
