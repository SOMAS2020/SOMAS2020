package rules

import (
	"testing"

	"gonum.org/v1/gonum/mat"
)

// TestBasicRuleEvaluatorPositive Checks whether rule we expect to evaluate as true actually evaluates as such
func TestBasicRuleEvaluatorPositive(t *testing.T) {
	ret := EvaluateRuleFromCaches("Kinda Complicated Rule", AvailableRules, generateMockVarCache())
	if !ret.RulePasses {
		t.Errorf("Rule evaluation came as false, when it was expected to be true, potential error with value '%v'", ret.EvalError)
	}
}

// TestBasicRuleEvaluatorNegative Checks whether rule we expect to evaluate as false actually evaluates as such
func TestBasicRuleEvaluatorNegative(t *testing.T) {
	registerTestRule(AvailableRules)
	ret := EvaluateRuleFromCaches("Kinda Test Rule", AvailableRules, generateMockVarCache())
	if ret.RulePasses || ret.EvalError != nil {
		t.Errorf("Rule evaluation came as true, when it was expected to be false, potential error with value '%v'", ret.EvalError)
	}
}

func TestBasicRealValuedRuleEvaluator(t *testing.T) {
	registerNewRealValuedRule(t)
	ret := EvaluateRuleFromCaches("Real Test rule", AvailableRules, generateMockVarCache())
	if !ret.RulePasses && ret.RealOutputVal != 2.0 {
		t.Errorf("Real values rule evaluation error, expected true got '%v', value expected '2' got '%v'", ret.RulePasses, ret.RealOutputVal)
	}
}

func TestBasicLinkedRuleEvaluator(t *testing.T) {
	registerNewLinkedRule(t)
	ret := EvaluateRuleFromCaches("Linked test rule", AvailableRules, generateMockVarCache())
	if ret.EvalError != nil {
		t.Errorf("Linked rule evaluation error: %v", ret.EvalError)
	}
	if !ret.RulePasses {
		t.Errorf("Linked rule evaluated to %v expected %v", ret.RulePasses, true)
	}
}

func generateMockVarCache() map[VariableFieldName]VariableValuePair {
	return map[VariableFieldName]VariableValuePair{
		NumberOfIslandsContributingToCommonPool: {
			NumberOfIslandsContributingToCommonPool,
			[]float64{5},
		},
		NumberOfFailedForages: {
			NumberOfFailedForages,
			[]float64{0.5},
		},
		NumberOfBrokenAgreements: {
			NumberOfBrokenAgreements,
			[]float64{1},
		},
		MaxSeverityOfSanctions: {
			MaxSeverityOfSanctions,
			[]float64{2},
		},
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
