package team1

import (
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func TestStealsResourcesWhenDesperate(t *testing.T) {
	availableRules, currentRulesInPlay := rules.InitialRuleRegistration(true)
	c := MakeTestClient(gamestate.ClientGameState{
		ClientInfo: gamestate.ClientInfo{
			LifeStatus: shared.Critical,
		},
		RulesInfo: gamestate.RulesContext{
			VariableMap: map[rules.VariableFieldName]rules.VariableValuePair{
				rules.ExpectedAllocation: {
					VariableName: rules.ExpectedAllocation,
					Values:       []float64{0},
				},
			},
			AvailableRules:     availableRules,
			CurrentRulesInPlay: currentRulesInPlay,
		},
		CommonPool: 300,
	})
	c.config.desperateStealAmount = 50
	c.LocalVariableCache[rules.ExpectedAllocation] = rules.MakeVariableValuePair(
		rules.ExpectedAllocation, []float64{0},
	)
	gotTakeAmt := c.RequestAllocation()
	if gotTakeAmt != c.config.desperateStealAmount {
		t.Errorf(
			`Client took the wrong amount while desperate:
				Expected: %v
				Got: %v`, c.config.desperateStealAmount, gotTakeAmt)
	}

}

func TestRequestsResourcesWhenAnxious(t *testing.T) {
	c := MakeTestClient(gamestate.ClientGameState{
		ClientInfo: gamestate.ClientInfo{
			Resources: 99,
		},
		CommonPool: 300,
	})
	c.config.resourceRequestScale = 1
	c.config.anxietyThreshold = 100

	gotRequest := c.CommonPoolResourceRequest()
	if gotRequest <= 0 {
		t.Errorf("Client did not make a resource request while anxious (%v)", gotRequest)
	}
}
