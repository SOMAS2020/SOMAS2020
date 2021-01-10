package team3

// IIGO President functions testing

import (
	"reflect"
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
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
				criticalThreshold: 30,
				iigoInfo:          iigoCommunicationInfo{commonPoolAllocation: shared.Resources(10)},
				params:            islandParams{resourcesSkew: 1.3, selfishness: 0.3, equity: 0.1, riskFactor: 0.2, saveCriticalIsland: false},
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
				BaseClient: &baseclient.BaseClient{ServerReadHandle: mockServerReadHandle{
					gameState: gamestate.ClientGameState{
						ClientInfo: gamestate.ClientInfo{LifeStatus: shared.Alive, Resources: 100},
						CommonPool: shared.Resources(40),
					},
					gameConfig: config.ClientConfig{
						IIGOClientConfig: config.IIGOConfig{
							GetRuleForSpeakerActionCost:        shared.Resources(2),
							BroadcastTaxationActionCost:        shared.Resources(2),
							ReplyAllocationRequestsActionCost:  shared.Resources(2),
							RequestAllocationRequestActionCost: shared.Resources(2),
							RequestRuleProposalActionCost:      shared.Resources(2),
							AppointNextSpeakerActionCost:       shared.Resources(2),
							InspectHistoryActionCost:           shared.Resources(2),
							HistoricalRetributionActionCost:    shared.Resources(2),
							InspectBallotActionCost:            shared.Resources(2),
							InspectAllocationActionCost:        shared.Resources(2),
							AppointNextPresidentActionCost:     shared.Resources(2),
							AssumedResourcesNoReport:           shared.Resources(2),
							SetVotingResultActionCost:          shared.Resources(2),
							SetRuleToVoteActionCost:            shared.Resources(2),
							AnnounceVotingResultActionCost:     shared.Resources(2),
							UpdateRulesActionCost:              shared.Resources(2),
							AppointNextJudgeActionCost:         shared.Resources(2),
						},
					},
				}},
				criticalThreshold: 10,
				params:            islandParams{selfishness: 0.3, riskFactor: 0.5, resourcesSkew: 1.3},
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
				0: {ReportedAmount: 100, Reported: true},
				1: {ReportedAmount: 100, Reported: true},
				2: {ReportedAmount: 100, Reported: true},
				3: {ReportedAmount: 100, Reported: true},
				4: {ReportedAmount: 100, Reported: true},
				5: {ReportedAmount: 100, Reported: true},
			},
			expected: map[shared.ClientID]shared.Resources{
				0: 11,
				1: 15,
				2: 15,
				3: 15,
				4: 15,
				5: 15,
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
