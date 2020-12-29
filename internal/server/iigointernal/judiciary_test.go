package iigointernal

import (
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func TestJudgeIncureServiceCharge(t *testing.T) {
	cases := []struct {
		name                string
		bJudge              judiciary // base
		input               shared.Resources
		expectedReturn      bool
		expectedCommonPool  shared.Resources
		expectedJudgeBudget shared.Resources
	}{
		{
			name: "Excess pay",
			bJudge: judiciary{
				JudgeID: shared.Team1,
				gameState: &gamestate.GameState{
					CommonPool: 400,
					IIGORolesBudget: map[string]shared.Resources{
						"president": 10,
						"speaker":   10,
						"judge":     100,
					},
				},
			},
			input:               50,
			expectedReturn:      true,
			expectedCommonPool:  350,
			expectedJudgeBudget: 50,
		},
		{
			name: "Negative Budget",
			bJudge: judiciary{
				JudgeID: shared.Team1,
				gameState: &gamestate.GameState{
					CommonPool: 400,
					IIGORolesBudget: map[string]shared.Resources{
						"president": 10,
						"speaker":   10,
						"judge":     10,
					},
				},
			},
			input:               50,
			expectedReturn:      true,
			expectedCommonPool:  350,
			expectedJudgeBudget: -40,
		},
		{
			name: "Limited common pool",
			bJudge: judiciary{
				JudgeID: shared.Team1,
				gameState: &gamestate.GameState{
					CommonPool: 40,
					IIGORolesBudget: map[string]shared.Resources{
						"president": 10,
						"speaker":   10,
						"judge":     10,
					},
				},
			},
			input:               50,
			expectedReturn:      false,
			expectedCommonPool:  40,
			expectedJudgeBudget: 10,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			returned := tc.bJudge.incurServiceCharge(tc.input)
			commonPool := tc.bJudge.gameState.CommonPool
			judgeBudget := tc.bJudge.gameState.IIGORolesBudget["judge"]
			if returned != tc.expectedReturn ||
				commonPool != tc.expectedCommonPool ||
				judgeBudget != tc.expectedJudgeBudget {
				t.Errorf("%v - Failed. Got '%v, %v, %v', but expected '%v, %v, %v'",
					tc.name, returned, commonPool, judgeBudget,
					tc.expectedReturn, tc.expectedCommonPool, tc.expectedJudgeBudget)
			}
		})
	}
}
