package team1

import (
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"

	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type mockGetRecommendation struct {
	success bool
}

func (c *mockGetRecommendation) GetRecommendation(variable rules.VariableFieldName) (compliantValue rules.VariableValuePair, success bool) {
	c.success = true
	return rules.MakeVariableValuePair(variable, []float64{0}), c.success
}

func TestStealsResourcesWhenDesperate(t *testing.T) {
	c := MakeTestClient(gamestate.ClientGameState{
		ClientInfo: gamestate.ClientInfo{
			LifeStatus: shared.Critical,
		},
		RulesInfo: gamestate.RulesContext{
			VariableMap: map[rules.VariableFieldName]rules.VariableValuePair{
				rules.AllocationMade: {
					rules.AllocationMade,
					[]float64{0},
				},
			},
			AvailableRules:     nil,
			CurrentRulesInPlay: nil,
		},
		CommonPool: 300,
	})
	c.config.desperateStealAmount = 50
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
	c.config.anxietyThreshold = 100

	gotRequest := c.CommonPoolResourceRequest()
	if gotRequest <= 0 {
		t.Errorf("Client did not make a resource request while anxious (%v)", gotRequest)
	}
}
