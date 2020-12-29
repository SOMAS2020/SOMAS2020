package team1

import (
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type testServerHandle struct {
	clientGameState gamestate.ClientGameState
}

func (h testServerHandle) GetGameState() gamestate.ClientGameState {
	return h.clientGameState
}

func MakeTestClient(gamestate gamestate.ClientGameState) client {
	c := NewClient(shared.Team1)
	c.Initialise(testServerHandle{
		clientGameState: gamestate,
	})
	return *c.(*client)
}

func TestRequestsResourcesWhenAnxious(t *testing.T) {
	c := MakeTestClient(gamestate.ClientGameState{
		ClientInfo: gamestate.ClientInfo{
			Resources: 99,
		},
	})
	c.config.anxietyThreshold = 100

	gotRequest := c.CommonPoolResourceRequest()
	if gotRequest <= 0 {
		t.Errorf("Client did not make a resource request while anxious (%v)", gotRequest)
	}
}

func TestStealsResourcesWhenDesperate(t *testing.T) {
	c := MakeTestClient(gamestate.ClientGameState{
		ClientInfo: gamestate.ClientInfo{
			LifeStatus: shared.Critical,
		},
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
