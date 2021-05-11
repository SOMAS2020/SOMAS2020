package team6

import (
	"reflect"
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func TestDecideGiftAmount(t *testing.T) {
	tests := []struct {
		testname      string
		testClient    client
		testToTeam    shared.ClientID
		testGiftOffer shared.Resources
		want          shared.Resources
	}{
		{
			testname: "critical little offer test",
			testClient: client{
				BaseClient: &baseclient.BaseClient{
					ServerReadHandle: stubServerReadHandle{
						gameState: gamestate.ClientGameState{
							ClientLifeStatuses: map[shared.ClientID]shared.ClientLifeStatus{
								shared.Teams["Team1"]: shared.Alive,
								shared.Teams["Team2"]: shared.Critical,
								shared.Teams["Team3"]: shared.Alive,
								shared.Teams["Team4"]: shared.Critical,
								shared.Teams["Team5"]: shared.Alive,
								shared.Teams["Team6"]: shared.Critical,
							},
						},
						gameConfig: config.ClientConfig{
							MinimumResourceThreshold: 5.0,
						},
					},
				},
				clientConfig: getClientConfig(),
			},
			testToTeam:    shared.Teams["Team2"],
			testGiftOffer: 3.0,
			want:          5.0,
		},
		{
			testname: "critical enough offer test",
			testClient: client{
				BaseClient: &baseclient.BaseClient{
					ServerReadHandle: stubServerReadHandle{
						gameState: gamestate.ClientGameState{
							ClientLifeStatuses: map[shared.ClientID]shared.ClientLifeStatus{
								shared.Teams["Team1"]: shared.Alive,
								shared.Teams["Team2"]: shared.Critical,
								shared.Teams["Team3"]: shared.Alive,
								shared.Teams["Team4"]: shared.Critical,
								shared.Teams["Team5"]: shared.Alive,
								shared.Teams["Team6"]: shared.Critical,
							},
						},
						gameConfig: config.ClientConfig{
							MinimumResourceThreshold: 5.0,
						},
					},
				},
				clientConfig: getClientConfig(),
			},
			testToTeam:    shared.Teams["Team4"],
			testGiftOffer: 10.0,
			want:          10.0,
		},
		{
			testname: "non critical test",
			testClient: client{
				BaseClient: &baseclient.BaseClient{
					ServerReadHandle: stubServerReadHandle{
						gameState: gamestate.ClientGameState{
							ClientLifeStatuses: map[shared.ClientID]shared.ClientLifeStatus{
								shared.Teams["Team1"]: shared.Alive,
								shared.Teams["Team2"]: shared.Critical,
								shared.Teams["Team3"]: shared.Alive,
								shared.Teams["Team4"]: shared.Critical,
								shared.Teams["Team5"]: shared.Alive,
								shared.Teams["Team6"]: shared.Critical,
							},
						},
						gameConfig: config.ClientConfig{
							MinimumResourceThreshold: 5.0,
						},
					},
				},
				clientConfig: getClientConfig(),
			},
			testToTeam:    shared.Teams["Team5"],
			testGiftOffer: 20.0,
			want:          20.0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.testname, func(t *testing.T) {
			res := tc.testClient.DecideGiftAmount(tc.testToTeam, tc.testGiftOffer)
			if !reflect.DeepEqual(res, tc.want) {
				t.Errorf("got %v, want %v", res, tc.want)
			}
		})
	}
}
