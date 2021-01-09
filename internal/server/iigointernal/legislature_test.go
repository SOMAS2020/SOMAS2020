package iigointernal

import (
	"reflect"
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"github.com/SOMAS2020/SOMAS2020/pkg/testutils"
	"gonum.org/v1/gonum/mat"
)

var requiredVariables = []rules.VariableFieldName{
	rules.IslandReportedResources,
	rules.ConstSanctionAmount,
	rules.TurnsLeftOnSanction,
}

func genRuleMatrixExample1(ruleName string) rules.RuleMatrix {
	var v1 = []float64{1, -1, 0}
	var coreMatrix1 = mat.NewDense(1, 3, v1)
	var aux1 = []float64{0}
	var auxiliaryVector1 = mat.NewVecDense(1, aux1)
	return rules.RuleMatrix{RuleName: ruleName, ApplicableMatrix: *coreMatrix1, AuxiliaryVector: *auxiliaryVector1, Mutable: true, RequiredVariables: requiredVariables}
}
func genRuleMatrixExample2(ruleName string) rules.RuleMatrix {
	var v2 = []float64{5, 4, 3}
	var coreMatrix2 = mat.NewDense(1, 3, v2)
	var aux2 = []float64{3}
	var auxiliaryVector2 = mat.NewVecDense(1, aux2)
	return rules.RuleMatrix{RuleName: ruleName, ApplicableMatrix: *coreMatrix2, AuxiliaryVector: *auxiliaryVector2, Mutable: true, RequiredVariables: requiredVariables}
}
func TestRuleVotedIn(t *testing.T) {
	avail, inPlay := generateRulesTestStores()
	fakeGameState := gamestate.GameState{
		CommonPool: 400,
		IIGORolesBudget: map[shared.Role]shared.Resources{
			shared.President: 10,
			shared.Speaker:   10,
			shared.Judge:     10,
		},
		RulesInfo: gamestate.RulesContext{
			VariableMap:        nil,
			AvailableRules:     avail,
			CurrentRulesInPlay: inPlay,
		},
	}
	s := legislature{
		gameState: &fakeGameState,
		gameConf:  &config.IIGOConfig{},
	}
	cases := []struct {
		name          string
		rule          rules.RuleMatrix
		expectedError bool
		want          rules.RuleErrorType
	}{
		{
			name:          "normal working",
			rule:          genRuleMatrixExample1("Kinda Test Rule"),
			expectedError: false,
		},
		{
			name:          "unidentified rule name",
			rule:          genRuleMatrixExample1("Unknown Rule"),
			expectedError: true,
			want:          rules.RuleNotInAvailableRulesCache,
		},
		{
			name:          "Rule already in play",
			rule:          genRuleMatrixExample1("Kinda Test Rule 2"),
			expectedError: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := s.updateRules(tc.rule, true)
			if tc.expectedError {
				if ruleErr, ok := got.(*rules.RuleError); ok {
					if ruleErr.Type() != tc.want {
						t.Errorf("Expected error type '%v' got error type '%v'", tc.want, ruleErr.Type())
					}
				} else {
					t.Errorf("Unrecognised Error format received, with message: '%v'", ruleErr.Error())
				}
			} else {
				testutils.CompareTestErrors(nil, got, t)
			}
		})
	}
	expectedRulesInPlay := map[string]rules.RuleMatrix{
		"Kinda Test Rule":   genRuleMatrixExample1("Kinda Test Rule"),
		"Kinda Test Rule 2": genRuleMatrixExample1("Kinda Test Rule 2"),
	}
	eq := reflect.DeepEqual(inPlay, expectedRulesInPlay)
	if !eq {
		t.Errorf("The rules in play are not the same as expected, expected '%v', got '%v'", expectedRulesInPlay, inPlay)
	}
}

func TestRuleVotedOut(t *testing.T) {
	avail, inPlay := generateRulesTestStores()
	fakeGameState := gamestate.GameState{
		CommonPool: 400,
		IIGORolesBudget: map[shared.Role]shared.Resources{
			shared.President: 10,
			shared.Speaker:   10,
			shared.Judge:     10,
		},
		RulesInfo: gamestate.RulesContext{
			VariableMap:        nil,
			AvailableRules:     avail,
			CurrentRulesInPlay: inPlay,
		},
	}
	s := legislature{
		gameState: &fakeGameState,
		gameConf:  &config.IIGOConfig{},
	}
	cases := []struct {
		name          string
		rule          rules.RuleMatrix
		expectedError bool
		want          rules.RuleErrorType
	}{
		{
			name:          "normal working",
			rule:          genRuleMatrixExample1("Kinda Test Rule"),
			expectedError: false,
		},
		{
			name:          "unidentified rule name",
			rule:          genRuleMatrixExample1("Unknown Rule"),
			expectedError: true,
			want:          rules.RuleNotInAvailableRulesCache,
		},
		{
			name:          "Rule already in play",
			rule:          genRuleMatrixExample1("Kinda Test Rule 2"),
			expectedError: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := s.updateRules(tc.rule, false)
			if tc.expectedError {
				if ruleErr, ok := got.(*rules.RuleError); ok {
					if ruleErr.Type() != tc.want {
						t.Errorf("Expected error type '%v' got error type '%v'", tc.want, ruleErr.Type())
					}
				} else {
					t.Errorf("Unrecognised Error format received, with message: '%v'", ruleErr.Error())
				}
			} else {
				testutils.CompareTestErrors(nil, got, t)
			}
		})
	}
	expectedRulesInPlay := map[string]rules.RuleMatrix{}
	eq := reflect.DeepEqual(inPlay, expectedRulesInPlay)
	if !eq {
		t.Errorf("The rules in play are not the same as expected, expected '%v', got '%v'", expectedRulesInPlay, inPlay)
	}
}

func TestModifiedRuleVotedIn(t *testing.T) {
	avail, inPlay := generateRulesTestStores()
	fakeGameState := gamestate.GameState{
		CommonPool: 400,
		IIGORolesBudget: map[shared.Role]shared.Resources{
			shared.President: 10,
			shared.Speaker:   10,
			shared.Judge:     10,
		},
		RulesInfo: gamestate.RulesContext{
			VariableMap:        nil,
			AvailableRules:     avail,
			CurrentRulesInPlay: inPlay,
		},
	}
	s := legislature{
		gameState: &fakeGameState,
		gameConf:  &config.IIGOConfig{},
	}
	cases := []struct {
		name          string
		rule          rules.RuleMatrix
		expectedError bool
		want          rules.RuleErrorType
	}{
		{
			name:          "normal working",
			rule:          genRuleMatrixExample2("Kinda Test Rule"),
			expectedError: false,
		},
		{
			name:          "unidentified rule name",
			rule:          genRuleMatrixExample2("Unknown Rule"),
			expectedError: true,
			want:          rules.RuleNotInAvailableRulesCache,
		},
		{
			name:          "Rule already in play",
			rule:          genRuleMatrixExample2("Kinda Test Rule 2"),
			expectedError: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := s.updateRules(tc.rule, true)
			if tc.expectedError {
				if ruleErr, ok := got.(*rules.RuleError); ok {
					if ruleErr.Type() != tc.want {
						t.Errorf("Expected error type '%v' got error type '%v'", tc.want, ruleErr.Type())
					}
				} else {
					t.Errorf("Unrecognised Error format received, with message: '%v'", ruleErr.Error())
				}
			} else {
				testutils.CompareTestErrors(nil, got, t)
			}
		})
	}
	expectedRulesInPlay := map[string]rules.RuleMatrix{
		"Kinda Test Rule 2": genRuleMatrixExample2("Kinda Test Rule 2"),
	}
	eq := reflect.DeepEqual(inPlay, expectedRulesInPlay)
	if !eq {
		t.Errorf("The rules in play are not the same as expected, expected '%v', got '%v'", expectedRulesInPlay, inPlay)
	}
	expectedAvailbleRules := map[string]rules.RuleMatrix{
		"Kinda Test Rule":   genRuleMatrixExample2("Kinda Test Rule"),
		"Kinda Test Rule 2": genRuleMatrixExample2("Kinda Test Rule 2"),
		"Kinda Test Rule 3": genRuleMatrixExample1("Kinda Test Rule 3"),
		"TestingRule1":      genRuleMatrixExample1("TestingRule1"),
		"TestingRule2":      genRuleMatrixExample1("TestingRule2"),
	}
	eq = reflect.DeepEqual(avail, expectedAvailbleRules)
	if !eq {
		t.Errorf("The rules in play are not the same as expected, expected '%v', got '%v'", expectedAvailbleRules, avail)
	}
}

func TestModifiedRuleVotedOut(t *testing.T) {
	avail, inPlay := generateRulesTestStores()
	fakeGameState := gamestate.GameState{
		CommonPool: 400,
		IIGORolesBudget: map[shared.Role]shared.Resources{
			shared.President: 10,
			shared.Speaker:   10,
			shared.Judge:     10,
		},
		RulesInfo: gamestate.RulesContext{
			VariableMap:        nil,
			AvailableRules:     avail,
			CurrentRulesInPlay: inPlay,
		},
	}
	s := legislature{
		gameState: &fakeGameState,
		gameConf:  &config.IIGOConfig{},
	}
	cases := []struct {
		name          string
		rule          rules.RuleMatrix
		expectedError bool
		want          rules.RuleErrorType
	}{
		{
			name:          "normal working",
			rule:          genRuleMatrixExample2("Kinda Test Rule"),
			expectedError: false,
		},
		{
			name:          "unidentified rule name",
			rule:          genRuleMatrixExample2("Unknown Rule"),
			expectedError: true,
			want:          rules.RuleNotInAvailableRulesCache,
		},
		{
			name:          "Rule already in play",
			rule:          genRuleMatrixExample2("Kinda Test Rule 2"),
			expectedError: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := s.updateRules(tc.rule, false)
			if tc.expectedError {
				if ruleErr, ok := got.(*rules.RuleError); ok {
					if ruleErr.Type() != tc.want {
						t.Errorf("Expected error type '%v' got error type '%v'", tc.want, ruleErr.Type())
					}
				} else {
					t.Errorf("Unrecognised Error format received, with message: '%v'", ruleErr.Error())
				}
			} else {
				testutils.CompareTestErrors(nil, got, t)
			}
		})
	}
	expectedRulesInPlay := map[string]rules.RuleMatrix{
		"Kinda Test Rule 2": genRuleMatrixExample1("Kinda Test Rule 2"),
	}
	eq := reflect.DeepEqual(inPlay, expectedRulesInPlay)
	if !eq {
		t.Errorf("The rules in play are not the same as expected, expected '%v', got '%v'", expectedRulesInPlay, inPlay)
	}
	expectedAvailbleRules := map[string]rules.RuleMatrix{
		"Kinda Test Rule":   genRuleMatrixExample1("Kinda Test Rule"),
		"Kinda Test Rule 2": genRuleMatrixExample1("Kinda Test Rule 2"),
		"Kinda Test Rule 3": genRuleMatrixExample1("Kinda Test Rule 3"),
		"TestingRule1":      genRuleMatrixExample1("TestingRule1"),
		"TestingRule2":      genRuleMatrixExample1("TestingRule2"),
	}
	eq = reflect.DeepEqual(avail, expectedAvailbleRules)
	if !eq {
		t.Errorf("The rules in play are not the same as expected, expected '%v', got '%v'", expectedAvailbleRules, avail)
	}
}

func generateRulesTestStores() (map[string]rules.RuleMatrix, map[string]rules.RuleMatrix) {
	return map[string]rules.RuleMatrix{
			"Kinda Test Rule":   genRuleMatrixExample1("Kinda Test Rule"),
			"Kinda Test Rule 2": genRuleMatrixExample1("Kinda Test Rule 2"),
			"Kinda Test Rule 3": genRuleMatrixExample1("Kinda Test Rule 3"),
			"TestingRule1":      genRuleMatrixExample1("TestingRule1"),
			"TestingRule2":      genRuleMatrixExample1("TestingRule2"),
		},
		map[string]rules.RuleMatrix{
			"Kinda Test Rule 2": genRuleMatrixExample1("Kinda Test Rule 2"),
		}

}

func TestSpeakerIncureServiceCharge(t *testing.T) {
	cases := []struct {
		name                  string
		bSpeaker              legislature // base
		input                 shared.Resources
		expectedReturn        bool
		expectedCommonPool    shared.Resources
		expectedSpeakerBudget shared.Resources
	}{
		{
			name: "Excess pay",
			bSpeaker: legislature{
				SpeakerID: shared.Team1,
				gameState: &gamestate.GameState{
					CommonPool: 400,
					IIGORolesBudget: map[shared.Role]shared.Resources{
						shared.President: 10,
						shared.Speaker:   100,
						shared.Judge:     10,
					},
				},
			},
			input:                 50,
			expectedReturn:        true,
			expectedCommonPool:    350,
			expectedSpeakerBudget: 50,
		},
		{
			name: "Negative Budget",
			bSpeaker: legislature{
				SpeakerID: shared.Team1,
				gameState: &gamestate.GameState{
					CommonPool: 400,
					IIGORolesBudget: map[shared.Role]shared.Resources{
						shared.President: 10,
						shared.Speaker:   10,
						shared.Judge:     10,
					},
				},
			},
			input:                 50,
			expectedReturn:        true,
			expectedCommonPool:    350,
			expectedSpeakerBudget: -40,
		},
		{
			name: "Limited common pool",
			bSpeaker: legislature{
				SpeakerID: shared.Team1,
				gameState: &gamestate.GameState{
					CommonPool: 40,
					IIGORolesBudget: map[shared.Role]shared.Resources{
						shared.President: 10,
						shared.Speaker:   10,
						shared.Judge:     10,
					},
				},
			},
			input:                 50,
			expectedReturn:        false,
			expectedCommonPool:    40,
			expectedSpeakerBudget: 10,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			returned := tc.bSpeaker.incurServiceCharge(tc.input)
			commonPool := tc.bSpeaker.gameState.CommonPool
			presidentBudget := tc.bSpeaker.gameState.IIGORolesBudget[shared.Speaker]
			if returned != tc.expectedReturn ||
				commonPool != tc.expectedCommonPool ||
				presidentBudget != tc.expectedSpeakerBudget {
				t.Errorf("%v - Failed. Got '%v, %v, %v', but expected '%v, %v, %v'",
					tc.name, returned, commonPool, presidentBudget,
					tc.expectedReturn, tc.expectedCommonPool, tc.expectedSpeakerBudget)
			}
		})
	}
}
