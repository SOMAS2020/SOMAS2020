package team6

import (
	"reflect"
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func TestEvaluateAllocationRequests(t *testing.T) {

	tests := []struct {
		testname            string
		testClient          client
		testResourceRequest map[shared.ClientID]shared.Resources
		testAvailCommonPool shared.Resources
		want                shared.PresidentReturnContent
	}{
		{
			testname: "selfish president having larger request than CP",
			testClient: newMockClient(shared.Team6, mockInit{
				serverReadHandle: stubServerReadHandle{
					gameState: gamestate.ClientGameState{
						ClientInfo: gamestate.ClientInfo{
							Resources: shared.Resources(49.0),
						},
					},
				},
			}),
			testResourceRequest: map[shared.ClientID]shared.Resources{
				shared.Team1: 20,
				shared.Team2: 10,
				shared.Team3: 10,
				shared.Team4: 20,
				shared.Team5: 40,
				shared.Team6: 310,
			},
			testAvailCommonPool: 400,
			want: shared.PresidentReturnContent{
				ContentType: shared.PresidentAllocation,
				ResourceMap: map[shared.ClientID]shared.Resources{
					shared.Team1: 15,
					shared.Team2: 7.5,
					shared.Team3: 7.5,
					shared.Team4: 15,
					shared.Team5: 30,
					shared.Team6: 300,
				},
				ActionTaken: true,
			},
		},
		{
			testname: "normal president with request sum larger than CP",
			testClient: newMockClient(shared.Team6, mockInit{
				serverReadHandle: stubServerReadHandle{
					gameState: gamestate.ClientGameState{
						ClientInfo: gamestate.ClientInfo{
							Resources: shared.Resources(51.0),
						},
					},
				},
			}),
			testResourceRequest: map[shared.ClientID]shared.Resources{
				shared.Team1: 40,
				shared.Team2: 20,
				shared.Team3: 10,
				shared.Team4: 10,
				shared.Team5: 10,
				shared.Team6: 10,
			},
			testAvailCommonPool: 100,
			want: shared.PresidentReturnContent{
				ContentType: shared.PresidentAllocation,
				ResourceMap: map[shared.ClientID]shared.Resources{
					shared.Team1: 30,
					shared.Team2: 15,
					shared.Team3: 7.5,
					shared.Team4: 7.5,
					shared.Team5: 7.5,
					shared.Team6: 7.5,
				},
				ActionTaken: true,
			},
		},
		{
			testname: "selfish president",
			testClient: newMockClient(shared.Team6, mockInit{
				serverReadHandle: stubServerReadHandle{
					gameState: gamestate.ClientGameState{
						ClientInfo: gamestate.ClientInfo{
							Resources: shared.Resources(49.0),
						},
					},
				},
			}),
			testResourceRequest: map[shared.ClientID]shared.Resources{
				shared.Team1: 75,
				shared.Team2: 5,
				shared.Team3: 10,
				shared.Team4: 10,
				shared.Team5: 50,
				shared.Team6: 200,
			},
			testAvailCommonPool: 400,
			want: shared.PresidentReturnContent{
				ContentType: shared.PresidentAllocation,
				ResourceMap: map[shared.ClientID]shared.Resources{
					shared.Team1: 75,
					shared.Team2: 5,
					shared.Team3: 10,
					shared.Team4: 10,
					shared.Team5: 50,
					shared.Team6: 200,
				},
				ActionTaken: true,
			},
		},
		{
			testname: "normal president with request sum less than CP",
			testClient: newMockClient(shared.Team6, mockInit{
				serverReadHandle: stubServerReadHandle{
					gameState: gamestate.ClientGameState{
						ClientInfo: gamestate.ClientInfo{
							Resources: shared.Resources(51.0),
						},
					},
				},
			}),
			testResourceRequest: map[shared.ClientID]shared.Resources{
				shared.Team1: 10,
				shared.Team2: 20,
				shared.Team3: 30,
				shared.Team4: 5,
				shared.Team5: 5,
				shared.Team6: 0,
			},
			testAvailCommonPool: 100,
			want: shared.PresidentReturnContent{
				ContentType: shared.PresidentAllocation,
				ResourceMap: map[shared.ClientID]shared.Resources{
					shared.Team1: 10,
					shared.Team2: 20,
					shared.Team3: 30,
					shared.Team4: 5,
					shared.Team5: 5,
					shared.Team6: 0,
				},
				ActionTaken: true,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.testname, func(t *testing.T) {
			testPresident := president{client: &tc.testClient}
			got := testPresident.EvaluateAllocationRequests(tc.testResourceRequest, tc.testAvailCommonPool)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}
