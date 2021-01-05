package rules

import (
	"fmt"
	"testing"

	"gonum.org/v1/gonum/mat"
)

// TestBasicRuleEvaluatorPositive Checks whether rule we expect to evaluate as true actually evaluates as such
func TestBasicRuleEvaluatorPositive(t *testing.T) {
	ret := EvaluateRule("Kinda Complicated Rule")
	if !ret.RulePasses {
		t.Errorf("Rule evaluation came as false, when it was expected to be true, potential error with value '%v'", ret.EvalError)
	}
}

// TestBasicRuleEvaluatorNegative Checks whether rule we expect to evaluate as false actually evaluates as such
func TestBasicRuleEvaluatorNegative(t *testing.T) {
	registerTestRule(AvailableRules)
	ret := EvaluateRule("Kinda Test Rule")
	if ret.RulePasses || ret.EvalError != nil {
		t.Errorf("Rule evaluation came as true, when it was expected to be false, potential error with value '%v'", ret.EvalError)
	}
}

func TestBasicRealValuedRuleEvaluator(t *testing.T) {
	registerNewRealValuedRule(t)
	ret := EvaluateRule("Real Test rule")
	if !ret.RulePasses && ret.RealOutputVal != 2.0 {
		t.Errorf("Real values rule evaluation error, expected true got '%v', value expected '2' got '%v'", ret.RulePasses, ret.RealOutputVal)
	}
}

func TestBasicLinkedRuleEvaluator(t *testing.T) {
	registerNewLinkedRule(t)
	ret := EvaluateRule("Linked test rule")
	if ret.EvalError != nil {
		t.Errorf("Linked rule evaluation error: %v", ret.EvalError)
	}
	if !ret.RulePasses {
		t.Errorf("Linked rule evaluated to %v expected %v", ret.RulePasses, true)
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
	aux := []float64{1, 1, 4, 0}
	AuxiliaryVector := mat.NewVecDense(4, aux)

	_, ruleError := RegisterNewRule(name, reqVar, *CoreMatrix, *AuxiliaryVector, false, RuleLink{
		Linked: false,
	})
	if ruleError != nil {
		t.Errorf("Problem with registering new real valued rule in test, error message : '%v'", ruleError.Error())
	}
	// Check internal/clients/team3/client.go for an implementation of a basic evaluator for this rule
}

func registerNewLinkedRule(t *testing.T) {
	name := "Linked test rule"
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

	_, ruleError := RegisterNewRule(name, reqVar, *CoreMatrix, *AuxiliaryVector, false, RuleLink{
		Linked:     true,
		LinkType:   ParentFailAutoRulePass,
		LinkedRule: "Kinda Complicated Rule",
	})
	if ruleError != nil {
		t.Errorf("Problem with registering new real valued rule in test, error message : '%v'", ruleError.Error())
	}
}

func createTaxRule(tax float64) (RuleMatrix, VariableValuePair) {
	reqVar := []VariableFieldName{
		IslandTaxContribution,
		ExpectedTaxContribution,
	}
	v := []float64{1, -1, 0}
	aux := []float64{2}

	rowLength := len(reqVar) + 1
	nrows := len(v) / rowLength

	CoreMatrix := mat.NewDense(nrows, rowLength, v)
	AuxiliaryVector := mat.NewVecDense(nrows, aux)

	rule := RuleMatrix{
		RuleName:          fmt.Sprintf("Tax = %.2f", tax),
		RequiredVariables: reqVar,
		ApplicableMatrix:  *CoreMatrix,
		AuxiliaryVector:   *AuxiliaryVector,
		Mutable:           false,
	}

	expectedVariable := VariableValuePair{
		VariableName: ExpectedTaxContribution,
		Values:       []float64{tax},
	}

	return rule, expectedVariable
}
