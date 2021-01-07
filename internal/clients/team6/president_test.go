package team6

import (
	"reflect"
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// ########################
// ######  Testing  #######
// ########################

func TestEvaluateAllocationRequests(t *testing.T) {
	tests := []struct {
		testname    string
		testP       president
		testC       client
		testRequest map[shared.ClientID]shared.Resources
		testCPR     shared.Resources
		want        shared.PresidentReturnContent
	}{
		{
			testname: "normalStrategy",
			testC: client{
				BaseClient: &baseclient.BaseClient{
					ServerReadHandle: stubServerReadHandle{
						gameState: gamestate.ClientGameState{
							ClientInfo: gamestate.ClientInfo{
								Resources: shared.Resources(149.0),
							},
						},
					},
				},
				clientConfig: clientConfig,
			},
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

		{
			testname: "egoisticStrategy",
			testC: client{
				BaseClient: &baseclient.BaseClient{
					ServerReadHandle: stubServerReadHandle{
						gameState: gamestate.ClientGameState{
							ClientInfo: gamestate.ClientInfo{
								Resources: shared.Resources(1),
							},
						},
					},
				},
				clientConfig: clientConfig,
			},
			testRequest: map[shared.ClientID]shared.Resources{
				shared.Team1: 10,
				shared.Team2: 10,
				shared.Team3: 10,
				shared.Team4: 10,
				shared.Team5: 10,
				shared.Team6: 10,
			},
			testCPR: 20,
			want: shared.PresidentReturnContent{
				ResourceMap: map[shared.ClientID]shared.Resources{
					shared.Team1: 1,
					shared.Team2: 1,
					shared.Team3: 1,
					shared.Team4: 1,
					shared.Team5: 1,
					shared.Team6: 10,
				},
				ActionTaken: true,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.testname, func(t *testing.T) {
			tc.testP = president{client: &tc.testC}
			got := tc.testP.EvaluateAllocationRequests(tc.testRequest, tc.testCPR)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}
