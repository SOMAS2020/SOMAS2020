package rules

import (
	"testing"

	"gonum.org/v1/gonum/mat"
)

// TestRegisterNewRule Tests whether the global rule cache is able to register new rules
func TestRegisterNewRule(t *testing.T) {
	AvailableRulesTesting, _ := generateRulesTestStores()
	registerTestRule(AvailableRulesTesting)
	if _, ok := AvailableRulesTesting["Kinda Test Rule"]; !ok {
		t.Errorf("Global rule register unable to register new rules")
	}
}

func TestPullRuleIntoPlay(t *testing.T) {
	AvailableRulesTesting, RulesInPlayTesting := generateRulesTestStores()
	registerTestRule(AvailableRulesTesting)
	_, _ = pullRuleIntoPlayInternal("Kinda Test Rule 2", AvailableRulesTesting, RulesInPlayTesting)
	cases := []struct {
		name    string
		rule    string
		success bool
		message RuleError
	}{
		{
			name:    "normal working",
			rule:    "Kinda Test Rule",
			success: true,
			message: None,
		},
		{
			name:    "unidentified rule name",
			rule:    "Unknown Rule",
			success: false,
			message: RuleNotInAvailableRulesCache,
		},
		{
			name:    "Rule already in play",
			rule:    "Kinda Test Rule 2",
			success: false,
			message: RuleIsAlreadyInPlay,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			success, errorMessage := pullRuleIntoPlayInternal(tc.rule, AvailableRulesTesting, RulesInPlayTesting)
			if success != tc.success {
				t.Errorf("Pulling rule into play test error, expected '%v' got '%v'", success, errorMessage)
			}
			if errorMessage != tc.message {
				t.Errorf("Pulling rule into play test error, expected message '%v' got '%v'", tc.message, errorMessage)
			}
		})
	}
}

func TestPullRuleOutOfPlay(t *testing.T) {
	AvailableRulesTesting, RulesInPlayTesting := generateRulesTestStores()
	registerTestRule(AvailableRulesTesting)
	_, _ = pullRuleIntoPlayInternal("Kinda Test Rule", AvailableRulesTesting, RulesInPlayTesting)
	cases := []struct {
		name    string
		rule    string
		success bool
		message RuleError
	}{
		{
			name:    "normal working",
			rule:    "Kinda Test Rule",
			success: true,
			message: None,
		},
		{
			name:    "Rule already in play",
			rule:    "Kinda Test Rule 2",
			success: false,
			message: RuleIsNotInPlay,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			success, errorMessage := pullRuleOutOfPlayInternal(tc.rule, AvailableRulesTesting, RulesInPlayTesting)
			if success != tc.success {
				t.Errorf("Pulling Rule out of play status wanted '%v' got '%v'", tc.success, success)
			}
			if errorMessage != tc.message {
				t.Errorf("Pulling Rule out of play error wanted: '%v' got '%v'", tc.message, errorMessage)
			}
		})
	}
}

func TestModifyRule(t *testing.T) {
	AvailableRulesTesting, RulesInPlayTesting := generateRulesTestStores()
	registerTestRule(AvailableRulesTesting)
	RulesInPlayTesting["Kinda Test Rule 2"] = AvailableRulesTesting["Kinda Test Rule 2"]
	BasicVect := []float64{1, 0, 0, 0, -4, 0, -1, -1, 0, 2, 0, 0, 0, 1, -2, 0, 0, 1, 0, -1}
	BasicAux := []float64{2, 3, 3, 2}
	cases := []struct {
		name           string
		rule           string
		modifiedMatrix mat.Dense
		modifiedAux    mat.VecDense
		expResult      bool
		expMessage     RuleError
	}{
		{
			name:           "Normal Rule Modification",
			rule:           "Kinda Test Rule 2",
			modifiedMatrix: *mat.NewDense(4, 5, BasicVect),
			modifiedAux:    *mat.NewVecDense(4, BasicAux),
			expResult:      true,
			expMessage:     None,
		},
		{
			name:           "Advanced Rule Modification",
			rule:           "Kinda Test Rule 2",
			modifiedMatrix: *mat.NewDense(5, 5, append(BasicVect, []float64{1, 1, 1, 1, 0}...)),
			modifiedAux:    *mat.NewVecDense(5, append(BasicAux, 1)),
			expResult:      true,
			expMessage:     None,
		},
		{
			name:           "Testing Rule not found error",
			rule:           "Fake Rule name",
			modifiedMatrix: *mat.NewDense(4, 5, BasicVect),
			modifiedAux:    *mat.NewVecDense(4, BasicAux),
			expResult:      false,
			expMessage:     RuleNotInAvailableRulesCache,
		},
		{
			name:           "Testing Matrix dimension mismatch",
			rule:           "Kinda Test Rule 2",
			modifiedMatrix: *mat.NewDense(5, 4, BasicVect),
			modifiedAux:    *mat.NewVecDense(4, BasicAux),
			expResult:      false,
			expMessage:     ModifiedRuleMatrixDimensionMismatch,
		},
		{
			name:           "Testing Auxiliary vector dimension mismatch",
			rule:           "Kinda Test Rule 2",
			modifiedMatrix: *mat.NewDense(4, 5, BasicVect),
			modifiedAux:    *mat.NewVecDense(5, append(BasicAux, 1)),
			expResult:      false,
			expMessage:     AuxVectorDimensionDontMatchRuleMatrix,
		},
		{
			name:           "Testing Immutable Rules stay immutable",
			rule:           "Kinda Test Rule",
			modifiedMatrix: *mat.NewDense(4, 5, BasicVect),
			modifiedAux:    *mat.NewVecDense(4, BasicAux),
			expResult:      false,
			expMessage:     RuleRequestedForModificationWasImmutable,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, message := modifyRuleInternal(tc.rule, tc.modifiedMatrix, tc.modifiedAux, AvailableRulesTesting, RulesInPlayTesting)
			if result != tc.expResult {
				t.Errorf("Rule modification result expected '%v' got '%v'", tc.expResult, result)
			}
			if message != tc.expMessage {
				t.Errorf("Rule modification message expected '%v' got '%v'", tc.expMessage, message)
			}
		})
	}
}

func generateRulesTestStores() (map[string]RuleMatrix, map[string]RuleMatrix) {
	return map[string]RuleMatrix{}, map[string]RuleMatrix{}
}

func registerTestRule(rulesStore map[string]RuleMatrix) {

	//A very contrived rule//
	name := "Kinda Test Rule"
	reqVar := []VariableFieldName{
		NumberOfIslandsContributingToCommonPool,
		NumberOfFailedForages,
		NumberOfBrokenAgreements,
		MaxSeverityOfSanctions,
	}

	v := []float64{1, 0, 0, 0, -4, 0, -1, -1, 0, 2, 0, 0, 0, 1, -2, 0, 0, 1, 0, -1}
	CoreMatrix := mat.NewDense(4, 5, v)
	aux := []float64{2, 3, 3, 2}
	AuxiliaryVector := mat.NewVecDense(4, aux)

	_, success, _ := registerNewRuleInternal(name, reqVar, *CoreMatrix, *AuxiliaryVector, rulesStore, false)
	// Check internal/clients/team3/client.go for an implementation of a basic evaluator for this rule
	//A very contrived rule//
	if !success {
		panic("Couldn't register Rule during test")
	}
	name = "Kinda Test Rule 2"

	_, success, _ = registerNewRuleInternal(name, reqVar, *CoreMatrix, *AuxiliaryVector, rulesStore, true)
	if !success {
		panic("Couldn't register Rule during test")
	}
}
