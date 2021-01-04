package team3

// IIGO President functions testing

import (
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func TestEvaluateAllocationRequests(t *testing.T) {
	cases := []struct {
		name            string
		ourPresident    president
		ourClient       client
		availCommonPool shared.Resources
		requests        map[shared.ClientID]shared.Resources
	}{
		{
			name: "Get Avg request",
			ourClient: client{
				BaseClient: &baseclient.BaseClient{ServerReadHandle: mockServerReadHandle{gameState: gamestate.ClientGameState{
					ClientInfo: gamestate.ClientInfo{LifeStatus: shared.Critical}}}},
				criticalStatePrediction: criticalStatePrediction{lowerBound: 30},
				iigoInfo:                iigoCommunicationInfo{commonPoolAllocation: shared.Resources(10)},
				params:                  islandParams{resourcesSkew: 1.3, selfishness: 0.3, equity: 0.1, riskFactor: 0.2, saveCriticalIsland: false},
				trustScore: map[shared.ClientID]float64{
					shared.Team1: 1,
					shared.Team2: 1,
					shared.Team3: 1,
					shared.Team4: 1,
					shared.Team5: 1,
					shared.Team6: 1,
				},
				declaredResources: map[shared.ClientID]shared.Resources{
					shared.Team1: 10,
					shared.Team2: 10,
					shared.Team3: 10,
					shared.Team4: 10,
					shared.Team5: 10,
					shared.Team6: 10,
				},
			},
			requests: map[shared.ClientID]shared.Resources{
				shared.Team1: 15,
				shared.Team2: 15,
				shared.Team3: 15,
				shared.Team4: 15,
				shared.Team5: 15,
				shared.Team6: 15,
			},
			availCommonPool: 100,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var sum shared.Resources
			tc.ourPresident = president{c: &tc.ourClient}
			ansMap, _ := tc.ourPresident.EvaluateAllocationRequests(tc.requests, tc.availCommonPool)
			for _, ans := range ansMap {
				sum += ans
			}
			if sum > tc.availCommonPool {
				t.Errorf("total allocation sum (%f) is greater than common pool(%f)", sum, tc.availCommonPool)
			}
		})
	}
}
