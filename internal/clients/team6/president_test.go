package team6

import (
	"reflect"
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// ########################
// ######  Testing  #######
// ########################
/*
type stubServerReadHandle struct {
	gameState  gamestate.ClientGameState
	gameConfig config.ClientConfig
}

func (s stubServerReadHandle) GetGameState() gamestate.ClientGameState {
	return s.gameState
}
func (s stubServerReadHandle) GetGameConfig() config.ClientConfig {
	return s.gameConfig
}
*/
func TestEvaluateAllocationRequests(t *testing.T) {
	tests := []struct {
		testname    string
		testClient  president
		testRequest map[shared.ClientID]shared.Resources
		testCPR     shared.Resources
		want        shared.PresidentReturnContent
	}{
		{
			testname:   "normalStrategy",
			testClient: president{},
			testRequest: map[shared.ClientID]shared.Resources{
				shared.Team1: 10,
				shared.Team2: 10,
				shared.Team3: 10,
				shared.Team4: 10,
				shared.Team5: 10,
				shared.Team6: 10,
			},
			testCPR: 100,
			want: shared.PresidentReturnContent{
				ResourceMap: map[shared.ClientID]shared.Resources{
					shared.Team1: 10,
					shared.Team2: 10,
					shared.Team3: 10,
					shared.Team4: 10,
					shared.Team5: 10,
					shared.Team6: 10,
				},
				ActionTaken: true,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.testname, func(t *testing.T) {
			got := tc.testClient.EvaluateAllocationRequests(tc.testRequest, tc.testCPR)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}
