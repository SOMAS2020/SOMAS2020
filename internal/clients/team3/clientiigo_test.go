package team3

// IIGO client functions testing

import (
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func TestRequestAllocation(t *testing.T) {
	cases := []struct {
		name      string
		ourClient client
		expected  shared.Resources
	}{
		{
			name: "Get critical difference",
			ourClient: client{
				BaseClient: &baseclient.BaseClient{ServerReadHandle: mockServerReadHandle{gameState: gamestate.ClientGameState{
					ClientInfo: gamestate.ClientInfo{LifeStatus: shared.Critical}}}},
				criticalStatePrediction: criticalStatePrediction{upperBound: 70, lowerBound: 30},
				iigoInfo:                iigoCommunicationInfo{commonPoolAllocation: shared.Resources(10)},
				params:                  islandParams{escapeCritcaIsland: true, selfishness: 0.3},
			},
			expected: shared.Resources(40),
		},
		{
			name: "Non-escape critical, non-cheat",
			ourClient: client{
				BaseClient: &baseclient.BaseClient{ServerReadHandle: mockServerReadHandle{gameState: gamestate.ClientGameState{
					ClientInfo: gamestate.ClientInfo{LifeStatus: shared.Critical}}}},
				compliance:              1.0,
				criticalStatePrediction: criticalStatePrediction{upperBound: 70, lowerBound: 30},
				iigoInfo:                iigoCommunicationInfo{commonPoolAllocation: shared.Resources(10)},
				params:                  islandParams{escapeCritcaIsland: false, selfishness: 0.3},
			},
			expected: shared.Resources(10),
		},
		{
			name: "Cheating",
			ourClient: client{
				BaseClient: &baseclient.BaseClient{ServerReadHandle: mockServerReadHandle{gameState: gamestate.ClientGameState{
					ClientInfo: gamestate.ClientInfo{LifeStatus: shared.Alive}}}},
				compliance:              0.0,
				criticalStatePrediction: criticalStatePrediction{upperBound: 70, lowerBound: 30},
				iigoInfo:                iigoCommunicationInfo{commonPoolAllocation: shared.Resources(10)},
				params:                  islandParams{escapeCritcaIsland: false, selfishness: 0.3},
			},
			expected: shared.Resources(10 + 10*0.3),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ans := tc.ourClient.RequestAllocation()
			if ans != tc.expected {
				t.Errorf("got %f, want %f", ans, tc.expected)
			}
		})
	}
}

func TestCommonPoolResourceRequest(t *testing.T) {
	cases := []struct {
		name      string
		ourClient client
		expected  shared.Resources
	}{
		{
			name: "Request critical difference",
			ourClient: client{
				BaseClient: &baseclient.BaseClient{ServerReadHandle: mockServerReadHandle{gameState: gamestate.ClientGameState{
					ClientInfo: gamestate.ClientInfo{LifeStatus: shared.Critical, Resources: 20}}}},
				compliance:              1.0,
				criticalStatePrediction: criticalStatePrediction{upperBound: 70, lowerBound: 30},
				iigoInfo:                iigoCommunicationInfo{},
				params:                  islandParams{minimumRequest: 10, escapeCritcaIsland: true, selfishness: 0.3},
			},
			expected: shared.Resources(30),
		},
		{
			name: "Non-escape critical, non-cheat",
			ourClient: client{
				BaseClient: &baseclient.BaseClient{ServerReadHandle: mockServerReadHandle{gameState: gamestate.ClientGameState{
					ClientInfo: gamestate.ClientInfo{LifeStatus: shared.Alive, Resources: 20}}}},
				compliance:              1.0,
				criticalStatePrediction: criticalStatePrediction{upperBound: 70, lowerBound: 30},
				iigoInfo:                iigoCommunicationInfo{},
				params:                  islandParams{minimumRequest: 50, escapeCritcaIsland: false, selfishness: 0.3},
			},
			expected: shared.Resources(50),
		},
		{
			name: "Cheating",
			ourClient: client{
				BaseClient: &baseclient.BaseClient{ServerReadHandle: mockServerReadHandle{gameState: gamestate.ClientGameState{
					ClientInfo: gamestate.ClientInfo{LifeStatus: shared.Alive}}}},
				compliance:              0.0,
				criticalStatePrediction: criticalStatePrediction{upperBound: 70, lowerBound: 30},
				iigoInfo:                iigoCommunicationInfo{},
				params:                  islandParams{minimumRequest: 50, escapeCritcaIsland: false, selfishness: 0.3},
			},
			expected: shared.Resources(50 + 50*0.3),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ans := tc.ourClient.CommonPoolResourceRequest()
			if ans != tc.expected {
				t.Errorf("got %f, want %f", ans, tc.expected)
			}
		})
	}
}
