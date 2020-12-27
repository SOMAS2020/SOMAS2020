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

	forageDecision shared.ForageDecision

	forageUpdateCalled bool
	gotForageDecision  shared.ForageDecision
}

func (c mockClientForage) DecideForage() (shared.ForageDecision, error) {
	return c.forageDecision, nil
}

func (c *mockClientForage) ForageUpdate(forageDecision shared.ForageDecision, resources shared.Resources) {
	c.forageUpdateCalled = true
	c.gotForageDecision = forageDecision
}

func TestForagingCallsForageUpdate(t *testing.T) {
	cases := []shared.ForageType{
		shared.DeerForageType,
		// TODO: Uncomment when fish foraging implemented
		// shared.FishForageType,
	}

	for _, tc := range cases {
		t.Run(fmt.Sprintf("%v", tc), func(t *testing.T) {
			forageDecision := shared.ForageDecision{
				Type:         tc,
				Contribution: 10,
			}

			client := mockClientForage{
				forageDecision:     forageDecision,
				forageUpdateCalled: false,
			}

			s := SOMASServer{
				gameState: gamestate.GameState{
					ClientInfos: map[shared.ClientID]gamestate.ClientInfo{
						shared.Team1: {
							LifeStatus: shared.Alive,
							Resources:  100,
						},
					},
				},
				clientMap: map[shared.ClientID]baseclient.Client{
					shared.Team1: &client,
				},
			}

			err := s.runForage()

			if err != nil {
				t.Errorf("runForage error: %v", err)
			}
			if !client.forageUpdateCalled {
				t.Errorf("ForageUpdate was not called")
			}
			if client.gotForageDecision != forageDecision {
				t.Errorf("ForageUpdate got the wrong forageDecision")
			}

		})
	}
}
