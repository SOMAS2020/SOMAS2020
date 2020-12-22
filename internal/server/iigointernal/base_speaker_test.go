package iigointernal

import (
	"reflect"
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/pkg/testutils"
	"github.com/pkg/errors"
)

// TestRegisterNewRule Tests whether the global rule cache is able to register new rules
// func TestRegisterNewRule(t *testing.T) {
// 	AvailableRulesTesting, _ := generateRulesTestStores()
// 	registerTestRule(AvailableRulesTesting)
// 	if _, ok := AvailableRulesTesting["Kinda Test Rule"]; !ok {
// 		t.Errorf("Global rule register unable to register new rules")
// 	}
// }
var ruleMatrixExample rules.RuleMatrix

func TestRuleVotedIn(t *testing.T) {
	rules.AvailableRules, rules.RulesInPlay = generateRulesTestStores()
	s := baseSpeaker{}
	cases := []struct {
		name string
		rule string
		want error
	}{
		{
			name: "normal working",
			rule: "Kinda Test Rule",
			want: nil,
		},
		{
			name: "unidentified rule name",
			rule: "Unknown Rule",
			want: errors.Errorf("Rule '%v' is not available in rules cache", "Unknown Rule"),
		},
		{
			name: "Rule already in play",
			rule: "Kinda Test Rule 2",
			want: nil,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := s.updateRules(tc.rule, true)
			testutils.CompareTestErrors(tc.want, got, t)
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
	s := baseSpeaker{}
	cases := []struct {
		name string
		rule string
		want error
	}{
		{
			name: "normal working",
			rule: "Kinda Test Rule",
			want: nil,
		},
		{
			name: "unidentified rule name",
			rule: "Unknown Rule",
			want: errors.Errorf("Rule '%v' is not available in rules cache", "Unknown Rule"),
		},
		{
			name: "Rule already in play",
			rule: "Kinda Test Rule 2",
			want: nil,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := s.updateRules(tc.rule, false)
			testutils.CompareTestErrors(tc.want, got, t)
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

type speakerState struct {
	ruleToVote   string
	VotingResult bool
}

func TestVoting(t *testing.T) {
	rules.AvailableRules, rules.RulesInPlay = generateRulesTestStores()
	s := baseSpeaker{clientSpeaker: nil}
	cases := []struct {
		name           string
		ruleID         string
		expectedStates []speakerState
		want           error
	}{
		{
			name:           "Rule given",
			ruleID:         "TestingRule1",
			expectedStates: []speakerState{{"TestingRule1", false}, {"TestingRule1", true}},
			want:           nil,
		},
		{
			name:           "Another rule given",
			ruleID:         "TestingRule2",
			expectedStates: []speakerState{{"TestingRule2", false}, {"TestingRule2", true}},
			want:           nil,
		},
		{
			name:           "No rule given",
			ruleID:         "",
			expectedStates: []speakerState{{"", false}, {"", false}},
			want:           nil,
		},
	}
	var stateTransfer [][]speakerState
	var expectedStateTransfer [][]speakerState
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			s.setRuleToVote(tc.ruleID)
			state1 := speakerState{s.RuleToVote, s.VotingResult}
			s.setVotingResult()
			state2 := speakerState{s.RuleToVote, s.VotingResult}
			got := s.announceVotingResultTest()

			stateTransfer = append(stateTransfer, []speakerState{state1, state2})
			tc.expectedStates[1].VotingResult = state2.VotingResult //Result is random
			expectedStateTransfer = append(expectedStateTransfer, tc.expectedStates)

			testutils.CompareTestErrors(tc.want, got, t)
		})
	}

	eq := reflect.DeepEqual(stateTransfer, expectedStateTransfer)
	if !eq {
		t.Errorf("The rules in play are not the same as expected, expected '%v', got '%v'", expectedStateTransfer, stateTransfer)
	}
}

func (s *baseSpeaker) announceVotingResultTest() error {

	var rule string
	var result bool
	var err error

	if s.clientSpeaker != nil {
		//Power to change what is declared completely, return "", _ for no announcement to occur
		rule, result, err = s.clientSpeaker.DecideAnnouncement(s.RuleToVote, s.VotingResult)
		//TODO: log of given vs. returned rule and result
		if err != nil {
			rule, result, _ = s.DecideAnnouncement(s.RuleToVote, s.VotingResult)
		}
	} else {
		rule, result, _ = s.DecideAnnouncement(s.RuleToVote, s.VotingResult)
	}

	//Reset
	s.RuleToVote = ""
	s.VotingResult = false

	if rule != "" {
		//Deduct action cost
		s.budget -= 10

		//Perform announcement
		//broadcastToAllIslands(s.Id, generateVotingResultMessage(rule, result))
		return s.updateRules(rule, result)
	}

	return nil
}
