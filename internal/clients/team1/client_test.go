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

func TestRequestsResourcesWhenDesperate(t *testing.T) {
	c := MakeTestClient(gamestate.ClientGameState{
		ClientInfo: gamestate.ClientInfo{
			LifeStatus: shared.Critical,
		},
	})

	gotRequest := c.CommonPoolResourceRequest()
	if gotRequest <= 0 {
		t.Errorf("Client did not make a resource request while desperate (%v)", gotRequest)
	}
}
