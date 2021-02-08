package team1

import (
	"math"
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type testServerHandle struct {
	clientGameState  gamestate.ClientGameState
	clientGameConfig config.ClientConfig
}

func (h testServerHandle) GetGameState() gamestate.ClientGameState {
	return h.clientGameState
}

func (h testServerHandle) GetGameConfig() config.ClientConfig {
	return h.clientGameConfig
}

func MakeTestClient(gamestate gamestate.ClientGameState) client {
	c := DefaultClient(shared.Team1)
	c.Initialise(testServerHandle{
		clientGameState: gamestate,
		clientGameConfig: config.ClientConfig{
			CostOfLiving: 10,
		},
	})
	return *c.(*client)
}

func TestRegression(t *testing.T) {
	expectCoeffs := []float64{6.0, 2.0, 5.0}

	outcomes := []ForageOutcome{}
	for i := 0.0; i < 50; i++ {
		outcomes = append(outcomes, ForageOutcome{
			contribution: shared.Resources(i),
			revenue: shared.Resources(
				expectCoeffs[2]*math.Pow(i, 2.0) +
					expectCoeffs[1]*i +
					expectCoeffs[0],
			),
		})
	}

	gotRegression, err := outcomeRegression(outcomes)
	if err != nil {
		t.Errorf("Regression error: %v", err)
	}

	for i, expectedCoeff := range expectCoeffs {
		gotCoeff := gotRegression.Coeff(i)
		diff := math.Abs(expectedCoeff - gotCoeff)
		if diff >= 0.0000001 {
			t.Errorf("Incorrect coeffecient %v:\n\tExpected: %v\n\tGot: %v",
				i,
				expectedCoeff,
				gotCoeff,
			)
		}
	}
}
