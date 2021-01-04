package rules

import (
	"fmt"
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

func TestBasicLinkedRuleEvaluator(t *testing.T) {
	registerNewLinkedRule(t)
	result, err := BasicLinkedRuleEvaluator("Linked test rule")
	if err != nil {
		t.Errorf("Linked rule evaluation error: %v", err)
	}
	if !result {
		t.Errorf("Linked rule evaluated to %v expected %v", result, true)
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

func TestBasicLocalRuleEvaluatorPositive(t *testing.T) {
	cases := []struct {
		name            string
		tax             float64
		taxContribution float64
		want            bool
	}{
		{
			name:            "Equal test",
			tax:             21.37,
			taxContribution: 21.37,
			want:            true,
		},
		{
			name:            "Greater than test",
			tax:             21.37,
			taxContribution: 22.22,
			want:            true,
		},
		{
			name:            "Smaller than test",
			tax:             5,
			taxContribution: 2.2,
			want:            false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rule, expectedVariable := createTaxRule(tc.tax)

			islandTaxVariable := VariableValuePair{
				VariableName: IslandTaxContribution,
				Values:       []float64{tc.taxContribution},
			}

			variableMap := map[VariableFieldName]VariableValuePair{}
			variableMap[expectedVariable.VariableName] = expectedVariable
			variableMap[islandTaxVariable.VariableName] = islandTaxVariable

			ok, err := BasicLocalBooleanRuleEvaluator(rule, variableMap)

			if err != nil {
				t.Error(err)
			}

			if ok != tc.want {
				t.Errorf("Mismatch in expected and evaluated")
			}
		})
	}

}
