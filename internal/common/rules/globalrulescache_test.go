package rules

import (
	"testing"

	"github.com/SOMAS2020/SOMAS2020/pkg/testutils"
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
	_ = PullRuleIntoPlayInternal("Kinda Test Rule 2", AvailableRulesTesting, RulesInPlayTesting)
	cases := []struct {
		name          string
		rule          string
		errorExpected bool
		errorType     RuleErrorType
	}{
		{
			name:          "normal working",
			rule:          "Kinda Test Rule",
			errorExpected: false,
		},
		{
			name:          "unidentified rule name",
			rule:          "Unknown Rule",
			errorExpected: true,
			errorType:     RuleNotInAvailableRulesCache,
		},
		{
			name:          "Rule already in play",
			rule:          "Kinda Test Rule 2",
			errorExpected: true,
			errorType:     RuleIsAlreadyInPlay,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := PullRuleIntoPlayInternal(tc.rule, AvailableRulesTesting, RulesInPlayTesting)

			if tc.errorExpected {
				if ruleErr, ok := err.(*RuleError); ok {
					if ruleErr.ErrorType != tc.errorType {
						t.Errorf("Expected error type '%v' got error type '%v'", tc.errorType.String(), ruleErr.ErrorType.String())
					}
				} else {
					t.Errorf("Unrecognised Error format recieved, with message: '%v'", ruleErr.Error())
				}
			} else {
				testutils.CompareTestErrors(nil, err, t)
			}
		})
	}
}

func TestPullRuleOutOfPlay(t *testing.T) {
	AvailableRulesTesting, RulesInPlayTesting := generateRulesTestStores()
	registerTestRule(AvailableRulesTesting)
	_ = PullRuleIntoPlayInternal("Kinda Test Rule", AvailableRulesTesting, RulesInPlayTesting)
	cases := []struct {
		name          string
		rule          string
		errorExpected bool
		errorType     RuleErrorType
	}{
		{
			name:          "normal working",
			rule:          "Kinda Test Rule",
			errorExpected: false,
		},
		{
			name:          "Rule already in play",
			rule:          "Kinda Test Rule 2",
			errorExpected: true,
			errorType:     RuleIsNotInPlay,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := PullRuleOutOfPlayInternal(tc.rule, AvailableRulesTesting, RulesInPlayTesting)

			if tc.errorExpected {
				if ruleErr, ok := err.(*RuleError); ok {
					if ruleErr.ErrorType != tc.errorType {
						t.Errorf("Expected error type '%v' got error type '%v'", tc.errorType.String(), ruleErr.ErrorType.String())
					}
				} else {
					t.Errorf("Unrecognised Error format received, with message: '%v'", ruleErr.Error())
				}
			} else {
				testutils.CompareTestErrors(nil, err, t)
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
		errorExpected  bool
		errType        RuleErrorType
	}{
		{
			name:           "Normal Rule Modification",
			rule:           "Kinda Test Rule 2",
			modifiedMatrix: *mat.NewDense(4, 5, BasicVect),
			modifiedAux:    *mat.NewVecDense(4, BasicAux),
			errorExpected:  false,
		},
		{
			name:           "Advanced Rule Modification",
			rule:           "Kinda Test Rule 2",
			modifiedMatrix: *mat.NewDense(5, 5, append(BasicVect, []float64{1, 1, 1, 1, 0}...)),
			modifiedAux:    *mat.NewVecDense(5, append(BasicAux, 1)),
			errorExpected:  false,
		},
		{
			name:           "Testing Rule not found error",
			rule:           "Fake Rule name",
			modifiedMatrix: *mat.NewDense(4, 5, BasicVect),
			modifiedAux:    *mat.NewVecDense(4, BasicAux),
			errorExpected:  true,
			errType:        RuleNotInAvailableRulesCache,
		},
		{
			name:           "Testing Matrix dimension mismatch",
			rule:           "Kinda Test Rule 2",
			modifiedMatrix: *mat.NewDense(5, 4, BasicVect),
			modifiedAux:    *mat.NewVecDense(4, BasicAux),
			errorExpected:  true,
			errType:        ModifiedRuleMatrixDimensionMismatch,
		},
		{
			name:           "Testing Auxiliary vector dimension mismatch",
			rule:           "Kinda Test Rule 2",
			modifiedMatrix: *mat.NewDense(4, 5, BasicVect),
			modifiedAux:    *mat.NewVecDense(5, append(BasicAux, 1)),
			errorExpected:  true,
			errType:        AuxVectorDimensionDontMatchRuleMatrix,
		},
		{
			name:           "Testing Immutable Rules stay immutable",
			rule:           "Kinda Test Rule",
			modifiedMatrix: *mat.NewDense(4, 5, BasicVect),
			modifiedAux:    *mat.NewVecDense(4, BasicAux),
			errorExpected:  true,
			errType:        RuleRequestedForModificationWasImmutable,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ModifyRuleInternal(tc.rule, tc.modifiedMatrix, tc.modifiedAux, AvailableRulesTesting, RulesInPlayTesting)

			if tc.errorExpected {
				if ruleErr, ok := err.(*RuleError); ok {
					if ruleErr.ErrorType != tc.errType {
						t.Errorf("Expected error type '%v' got error type '%v'", tc.errType.String(), ruleErr.ErrorType.String())
					}
				} else {
					t.Errorf("Unrecognised Error format received, with message: '%v'", ruleErr.Error())
				}
			} else {
				testutils.CompareTestErrors(nil, err, t)
			}
		})
	}
}

func TestCheckLinking(t *testing.T) {
	cases := []struct {
		name               string
		ruleName           string
		rulesCache         map[string]RuleMatrix
		expectedLinkedRule string
		expectedLinked     bool
	}{
		{
			name:     "Basic non linked rule",
			ruleName: "Some rule",
			rulesCache: map[string]RuleMatrix{
				"Some rule": {
					RuleName: "Some rule",
					Link: RuleLink{
						Linked: false,
					},
				},
			},
			expectedLinkedRule: "",
			expectedLinked:     false,
		},
		{
			name:     "Basic linked rule",
			ruleName: "Some rule",
			rulesCache: map[string]RuleMatrix{
				"Some rule": {
					RuleName: "Some rule",
					Link: RuleLink{
						Linked:     true,
						LinkedRule: "Linker Rule",
					},
				},
				"Linker Rule": {
					RuleName: "Linker Rule",
					Link: RuleLink{
						Linked: false,
					},
				},
			},
			expectedLinkedRule: "Linker Rule",
			expectedLinked:     true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			lnkRule, linked := checkLinking(tc.ruleName, tc.rulesCache)

			if lnkRule != tc.expectedLinkedRule {
				t.Errorf("Expected: %v got %v", tc.expectedLinkedRule, lnkRule)
			}
			if linked != tc.expectedLinked {
				t.Errorf("Expected: %v got %v", tc.expectedLinked, linked)
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

	_, err := RegisterNewRuleInternal(name, reqVar, *CoreMatrix, *AuxiliaryVector, rulesStore, false, RuleLink{
		Linked: false,
	})
	// Check internal/clients/team3/client.go for an implementation of a basic evaluator for this rule
	//A very contrived rule//
	if err != nil {
		panic("Couldn't register Rule during test")
	}
	name = "Kinda Test Rule 2"

	_, err = RegisterNewRuleInternal(name, reqVar, *CoreMatrix, *AuxiliaryVector, rulesStore, true, RuleLink{
		Linked: false,
	})
	if err != nil {
		panic("Couldn't register Rule during test")
	}
}
