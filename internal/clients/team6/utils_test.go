package team6

import (
	"reflect"
	"testing"

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
