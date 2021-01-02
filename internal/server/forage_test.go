package server

import (
	"fmt"
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
	"github.com/SOMAS2020/SOMAS2020/internal/common/disasters"
	"github.com/SOMAS2020/SOMAS2020/internal/common/foraging"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type mockClientForage struct {
	baseclient.Client
	forageDecision     shared.ForageDecision
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

var envConf = config.DisasterConfig{
	XMax:            10,
	YMax:            10,
	GlobalProb:      0.1,
	SpatialPDFType:  shared.Uniform,
	MagnitudeLambda: 1.0,
}

var deerConf = config.DeerHuntConfig{
	MaxDeerPerHunt:     4,
	BernoulliProb:      0.95,
	ResourceMultiplier: 1,
	ExponentialRate:    1.0,
}

var cpConf = config.CommonPoolConfig{} // define if necessary but generally not required here

func TestForagingCallsForageUpdate(t *testing.T) {
	cases := []shared.ForageType{
		shared.DeerForageType,
		shared.FishForageType,
	}

	contribs := []shared.Resources{0.0, 1.0, 8.0} // test zero resource contribution first off

	for _, contrib := range contribs {
		for _, tc := range cases {
			t.Run(fmt.Sprintf("%v", tc), func(t *testing.T) {
				forageDecision := shared.ForageDecision{
					Type:         tc,
					Contribution: contrib,
				}

				client := mockClientForage{
					forageDecision:     forageDecision,
					forageUpdateCalled: false,
				}

				clientMap := map[shared.ClientID]baseclient.Client{
					shared.Team1: &client,
				}

				clientIDs := make([]shared.ClientID, 0, len(clientMap))
				for k := range clientMap {
					clientIDs = append(clientIDs, k)
				}

				s := SOMASServer{
					gameState: gamestate.GameState{
						ClientInfos: map[shared.ClientID]gamestate.ClientInfo{
							shared.Team1: {
								LifeStatus: shared.Alive,
								Resources:  100,
							},
						},
						ForagingHistory: map[shared.ForageType][]foraging.ForagingReport{
							shared.DeerForageType: make([]foraging.ForagingReport, 0),
							shared.FishForageType: make([]foraging.ForagingReport, 0),
						},
						DeerPopulation: foraging.CreateDeerPopulationModel(deerConf),
						Environment:    disasters.InitEnvironment(clientIDs, envConf, cpConf),
					},
					clientMap: clientMap,
				}

				err := s.runForage()

				if err != nil {
					t.Errorf("runForage error: %v", err)
				}
				if contrib > 0 {
					if !client.forageUpdateCalled {
						t.Errorf("ForageUpdate was not called")
					}
					if client.gotForageDecision != forageDecision {
						t.Errorf("ForageUpdate got the wrong forageDecision")
					}
				}
			})
		}

	}
}
