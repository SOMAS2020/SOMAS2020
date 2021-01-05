package team3

// IIGO President functions testing

import (
	"reflect"
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
			ansMap := tc.ourPresident.EvaluateAllocationRequests(tc.requests, tc.availCommonPool).ResourceMap
			for _, ans := range ansMap {
				sum += ans
			}
			if sum > tc.availCommonPool {
				t.Errorf("total allocation sum (%f) is greater than common pool(%f)", sum, tc.availCommonPool)
			}
		})
	}
}

func TestSetTaxationAmount(t *testing.T) {
	cases := []struct {
		name              string
		president         president
		declaredResources map[shared.ClientID]shared.ResourcesReport
		expected          map[shared.ClientID]shared.Resources
	}{
		{
			name: "Normal",
			president: president{c: &client{
				BaseClient: &baseclient.BaseClient{ServerReadHandle: mockServerReadHandle{gameState: gamestate.ClientGameState{
					ClientInfo: gamestate.ClientInfo{LifeStatus: shared.Alive,
						Resources: 100,
					},
					CommonPool: shared.Resources(40),
				}}},
				criticalStatePrediction: criticalStatePrediction{upperBound: 10, lowerBound: 0},
				params:                  islandParams{escapeCritcaIsland: true, selfishness: 0.3, riskFactor: 0.5, resourcesSkew: 1.3},
				trustScore: map[shared.ClientID]float64{
					0: 50,
					1: 50,
					2: 50,
					3: 50,
					4: 50,
					5: 50,
				},
				compliance: 1,
			}},
			declaredResources: map[shared.ClientID]shared.ResourcesReport{
				0: {100, true},
				1: {100, true},
				2: {100, true},
				3: {100, true},
				4: {100, true},
				5: {100, true},
			},
			expected: map[shared.ClientID]shared.Resources{
				0: 7,
				1: 10,
				2: 10,
				3: 10,
				4: 10,
				5: 10,
			}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ans := tc.president.SetTaxationAmount(tc.declaredResources).ResourceMap
			if !reflect.DeepEqual(ans, tc.expected) {
				t.Errorf("got %v, want %v", ans, tc.expected)
			}
		})
	}
}
