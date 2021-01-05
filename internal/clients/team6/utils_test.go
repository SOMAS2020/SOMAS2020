package team6

import (
	"reflect"
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func TestRaiseFriendshipLevel(t *testing.T) {
	tests := []struct {
		testname      string
		testClient    client
		testTeam      shared.ClientID
		testIncrement FriendshipLevel
		want          Friendship
	}{
		{
			testname: "common test",
			testClient: client{
				friendship: Friendship{
					shared.Team1: 50.0,
					shared.Team2: 50.0,
					shared.Team3: 50.0,
					shared.Team4: 50.0,
					shared.Team5: 50.0,
				},
				config: config,
			},
			testTeam:      shared.Team3,
			testIncrement: FriendshipLevel(50.0),
			want: Friendship{
				shared.Team1: 50.0,
				shared.Team2: 50.0,
				shared.Team3: 60.0,
				shared.Team4: 50.0,
				shared.Team5: 50.0,
			},
		},
		{
			testname: "overflow test",
			testClient: client{
				friendship: Friendship{
					shared.Team1: 50.0,
					shared.Team2: 50.0,
					shared.Team3: 50.0,
					shared.Team4: 50.0,
					shared.Team5: 90.0,
				},
				config: config,
			},
			testTeam:      shared.Team5,
			testIncrement: FriendshipLevel(100.0),
			want: Friendship{
				shared.Team1: 50.0,
				shared.Team2: 50.0,
				shared.Team3: 50.0,
				shared.Team4: 50.0,
				shared.Team5: 100.0,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.testname, func(t *testing.T) {
			tc.testClient.raiseFriendshipLevel(tc.testTeam, tc.testIncrement)
			if !reflect.DeepEqual(tc.testClient.friendship, tc.want) {
				t.Errorf("got %v, want %v", tc.testClient.friendship, tc.want)
			}
		})
	}
}

func TestLowerFriendshipLevel(t *testing.T) {
	tests := []struct {
		testname      string
		testClient    client
		testTeam      shared.ClientID
		testDeduction FriendshipLevel
		want          Friendship
	}{
		{
			testname: "common test",
			testClient: client{
				friendship: Friendship{
					shared.Team1: 50.0,
					shared.Team2: 50.0,
					shared.Team3: 50.0,
					shared.Team4: 50.0,
					shared.Team5: 50.0,
				},
				config: config,
			},
			testTeam:      shared.Team3,
			testDeduction: FriendshipLevel(50.0),
			want: Friendship{
				shared.Team1: 50.0,
				shared.Team2: 50.0,
				shared.Team3: 40.0,
				shared.Team4: 50.0,
				shared.Team5: 50.0,
			},
		},
		{
			testname: "overflow test",
			testClient: client{
				friendship: Friendship{
					shared.Team1: 50.0,
					shared.Team2: 50.0,
					shared.Team3: 50.0,
					shared.Team4: 50.0,
					shared.Team5: 10.0,
				},
				config: config,
			},
			testTeam:      shared.Team5,
			testDeduction: FriendshipLevel(100.0),
			want: Friendship{
				shared.Team1: 50.0,
				shared.Team2: 50.0,
				shared.Team3: 50.0,
				shared.Team4: 50.0,
				shared.Team5: 0.0,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.testname, func(t *testing.T) {
			tc.testClient.lowerFriendshipLevel(tc.testTeam, tc.testDeduction)
			if !reflect.DeepEqual(tc.testClient.friendship, tc.want) {
				t.Errorf("got %v, want %v", tc.testClient.friendship, tc.want)
			}
		})
	}
}

func TestGetFriendshipCoeffs(t *testing.T) {
	tests := []struct {
		testname   string
		testClient client
		want       map[shared.ClientID]float64
	}{
		{
			testname: "test iteration 1",
			testClient: client{
				friendship: Friendship{
					shared.Team1: 0,
					shared.Team2: 3.14,
					shared.Team3: 21.63,
					shared.Team4: 75.28,
					shared.Team5: 99.31,
				},
				config: config,
			},
			want: map[shared.ClientID]float64{
				shared.Team1: 0.0,
				shared.Team2: 0.0314,
				shared.Team3: 0.2163,
				shared.Team4: 0.7528,
				shared.Team5: 0.9931,
			},
		},
		{
			testname: "test iteration 2",
			testClient: client{
				friendship: Friendship{
					shared.Team1: 22.01,
					shared.Team2: 3.17,
					shared.Team3: 98.39,
					shared.Team4: 3.35,
					shared.Team5: 100.0,
				},
				config: config,
			},
			want: map[shared.ClientID]float64{
				shared.Team1: 0.2201,
				shared.Team2: 0.0317,
				shared.Team3: 0.9839,
				shared.Team4: 0.0335,
				shared.Team5: 1,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.testname, func(t *testing.T) {
			res := tc.testClient.getFriendshipCoeffs()

			for team, coeff := range res {
				if coeff-tc.want[team] > 0.0000001 {
					t.Errorf("got %v, want %v", res, tc.want)
				}
			}
		})
	}
}

func TestGetPersonality(t *testing.T) {
	tests := []struct {
		testname      string
		testClient    client
		testClientGamestate gamestate.ClientGameState
		want          Personality
	}{
		{
			testname: "selfish test",
			testClient: client{
				config: config,
			},
			testClientGamestate:  gamestate.ClientGameState{
				ClientInfo: gamestate.ClientInfo{
						Resources: shared.Resources(30.0),
					},
				},
			},
			want: Personality(Selfish),
		},
		{
			testname: "normal test",
			testClient: client{
				config: config,
			},
			testClientGamestate:  gamestate.ClientGameState{
				ClientInfo: gamestate.ClientInfo{
						Resources: shared.Resources(105.0),
					},
				},
			},
			want: Personality(Normal),
		},
		{
			testname: "generous test",
			testClient: client{
				config: config,
			},
			testClientGamestate:  gamestate.ClientGameState{
				ClientInfo: gamestate.ClientInfo{
						Resources: shared.Resources(205.0),
					},
				},
			},
			want: Personality(Generous),
		},
	}

	for _, tc := range tests {
		t.Run(tc.testname, func(t *testing.T) {
			res := tc.testClient.getPersonality()
			if !reflect.DeepEqual(res, tc.want) {
				t.Errorf("got %v, want %v", tc.testClient.friendship, tc.want)
			}
		})
	}
}
