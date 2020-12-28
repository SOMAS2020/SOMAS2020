package iigointernal

import (
	"reflect"
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/pkg/testutils"
)

var ruleMatrixExample rules.RuleMatrix

func TestRuleVotedIn(t *testing.T) {
	rules.AvailableRules, rules.RulesInPlay = generateRulesTestStores()
	s := legislature{}
	cases := []struct {
		name          string
		rule          string
		expectedError bool
		want          rules.RuleErrorType
	}{
		{
			name:          "normal working",
			rule:          "Kinda Test Rule",
			expectedError: false,
		},
		{
			name:          "unidentified rule name",
			rule:          "Unknown Rule",
			expectedError: true,
			want:          rules.RuleNotInAvailableRulesCache,
		},
		{
			name:          "Rule already in play",
			rule:          "Kinda Test Rule 2",
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
		"Kinda Test Rule":   ruleMatrixExample,
		"Kinda Test Rule 2": ruleMatrixExample,
	}
	eq := reflect.DeepEqual(rules.RulesInPlay, expectedRulesInPlay)
	if !eq {
		t.Errorf("The rules in play are not the same as expected, expected '%v', got '%v'", expectedRulesInPlay, rules.RulesInPlay)
	}
}

func TestRuleVotedOut(t *testing.T) {
	rules.AvailableRules, rules.RulesInPlay = generateRulesTestStores()
	s := legislature{}
	cases := []struct {
		name          string
		rule          string
		expectedError bool
		want          rules.RuleErrorType
	}{
		{
			name:          "normal working",
			rule:          "Kinda Test Rule",
			expectedError: false,
		},
		{
			name:          "unidentified rule name",
			rule:          "Unknown Rule",
			expectedError: true,
			want:          rules.RuleNotInAvailableRulesCache,
		},
		{
			name:          "Rule already in play",
			rule:          "Kinda Test Rule 2",
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
	eq := reflect.DeepEqual(rules.RulesInPlay, expectedRulesInPlay)
	if !eq {
		t.Errorf("The rules in play are not the same as expected, expected '%v', got '%v'", expectedRulesInPlay, rules.RulesInPlay)
	}
}

func generateRulesTestStores() (map[string]rules.RuleMatrix, map[string]rules.RuleMatrix) {
	return map[string]rules.RuleMatrix{
			"Kinda Test Rule":   ruleMatrixExample,
			"Kinda Test Rule 2": ruleMatrixExample,
			"Kinda Test Rule 3": ruleMatrixExample,
			"TestingRule1":      ruleMatrixExample,
			"TestingRule2":      ruleMatrixExample,
		},
		map[string]rules.RuleMatrix{
			"Kinda Test Rule 2": ruleMatrixExample,
		}

}
