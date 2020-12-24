package server

import (
	"fmt"
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type mockClientForage struct {
	baseclient.Client

	forageDecision        shared.ForageDecision
	got_forageUpdateValue shared.Resources
}

func (c mockClientForage) DecideForage() (shared.ForageDecision, error) {
	return c.forageDecision, nil
}

func (c *mockClientForage) ForageUpdate(resources shared.Resources) {
	c.got_forageUpdateValue = resources
}

func TestForagingCallsForageUpdate(t *testing.T) {
	cases := []shared.ForageType{
		shared.DeerForageType,
		// TODO: Uncomment when fish foraging implemented
		// shared.FishForageType,
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%v", tc), func(t *testing.T) {
			client := mockClientForage{
				forageDecision: shared.ForageDecision{
					Type:         tc,
					Contribution: 10,
				},

				got_forageUpdateValue: -1,
			}

			s := SOMASServer{
				gameState: gamestate.GameState{
					ClientInfos: map[shared.ClientID]gamestate.ClientInfo{
						shared.Team1: {
							LifeStatus: shared.Alive,
							Resources: 100,
						},
					},
				},
				clientMap: map[shared.ClientID]baseclient.Client{
					shared.Team1: &client,
				},
			}

			s.runForage()

			if client.got_forageUpdateValue < 0 {
				t.Errorf("ForageUpdate was not called (%v)", client.got_forageUpdateValue)
			}

		})
	}
}
