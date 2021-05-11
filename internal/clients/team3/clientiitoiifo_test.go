package team3

import (
	"reflect"
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func TestGetGiftRequests(t *testing.T) {
	cases := []struct {
		name             string
		serverReadHandle baseclient.ServerReadHandle
		ourClient        client
		expectedVal      shared.GiftRequestDict
	}{
		{
			name: "Basic test: all islands alive and trusted equally",
			serverReadHandle: mockServerReadHandle{
				gameState: gamestate.ClientGameState{
					ClientLifeStatuses: map[shared.ClientID]shared.ClientLifeStatus{
						shared.Teams["Team1"]: shared.Alive,
						shared.Teams["Team2"]: shared.Alive,
						shared.Teams["Team3"]: shared.Alive,
						shared.Teams["Team4"]: shared.Alive,
						shared.Teams["Team5"]: shared.Alive,
						shared.Teams["Team6"]: shared.Alive,
					},
					ClientInfo: gamestate.ClientInfo{Resources: 600.0}}},
			ourClient: client{
				BaseClient: baseclient.NewClient(shared.Teams["Team3"]),
				params: islandParams{
					giftInflationPercentage: 0.1,
					riskFactor:              0.2,
				},
				trustScore: map[shared.ClientID]float64{
					shared.Teams["Team1"]: 50,
					shared.Teams["Team2"]: 50,
					shared.Teams["Team4"]: 50,
					shared.Teams["Team5"]: 50,
					shared.Teams["Team6"]: 50,
				},
				requestedGiftAmounts: map[shared.ClientID]shared.GiftRequest{
					shared.Teams["Team1"]: 3.644540246477594,
					shared.Teams["Team2"]: 3.644540246477594,
					shared.Teams["Team4"]: 3.644540246477594,
					shared.Teams["Team5"]: 3.644540246477594,
					shared.Teams["Team6"]: 3.644540246477594,
				},
				initialResourcesAtStartOfGame: 100,
			},
			expectedVal: shared.GiftRequestDict{
				shared.Teams["Team1"]: 3.644540246477594,
				shared.Teams["Team2"]: 3.644540246477594,
				shared.Teams["Team4"]: 3.644540246477594,
				shared.Teams["Team5"]: 3.644540246477594,
				shared.Teams["Team6"]: 3.644540246477594,
			},
		},
		{
			name: "Basic test: Half islands are critical, other half are dead and trusted equally",
			serverReadHandle: mockServerReadHandle{gameState: gamestate.ClientGameState{
				ClientLifeStatuses: map[shared.ClientID]shared.ClientLifeStatus{
					shared.Teams["Team1"]: shared.Critical,
					shared.Teams["Team2"]: shared.Dead,
					shared.Teams["Team3"]: shared.Critical,
					shared.Teams["Team4"]: shared.Dead,
					shared.Teams["Team5"]: shared.Critical,
					shared.Teams["Team6"]: shared.Dead,
				},
				ClientInfo: gamestate.ClientInfo{Resources: 600.0}}},
			ourClient: client{
				BaseClient: baseclient.NewClient(shared.Teams["Team3"]),
				params: islandParams{
					giftInflationPercentage: 0.1,
					riskFactor:              0.2,
				},
				trustScore: map[shared.ClientID]float64{
					shared.Teams["Team1"]: 50,
					shared.Teams["Team2"]: 50,
					shared.Teams["Team4"]: 50,
					shared.Teams["Team5"]: 50,
					shared.Teams["Team6"]: 50,
				},
				requestedGiftAmounts: map[shared.ClientID]shared.GiftRequest{
					shared.Teams["Team1"]: 0,
					shared.Teams["Team2"]: 0,
					shared.Teams["Team4"]: 0,
					shared.Teams["Team5"]: 0,
					shared.Teams["Team6"]: 0,
				},
				initialResourcesAtStartOfGame: 100,
			},
			expectedVal: shared.GiftRequestDict{
				shared.Teams["Team1"]: 0,
				shared.Teams["Team2"]: 0,
				shared.Teams["Team4"]: 0,
				shared.Teams["Team5"]: 0,
				shared.Teams["Team6"]: 0,
			},
		},
		{
			name: "Complex test: All islands are alive, but are trusted differently",
			serverReadHandle: mockServerReadHandle{gameState: gamestate.ClientGameState{
				ClientLifeStatuses: map[shared.ClientID]shared.ClientLifeStatus{
					shared.Teams["Team1"]: shared.Alive,
					shared.Teams["Team2"]: shared.Alive,
					shared.Teams["Team3"]: shared.Alive,
					shared.Teams["Team4"]: shared.Alive,
					shared.Teams["Team5"]: shared.Alive,
					shared.Teams["Team6"]: shared.Alive,
				},
				ClientInfo: gamestate.ClientInfo{Resources: 10.0}}},
			ourClient: client{
				BaseClient: baseclient.NewClient(shared.Teams["Team3"]),
				params: islandParams{
					giftInflationPercentage: 0.1,
					riskFactor:              0.2,
				},
				trustScore: map[shared.ClientID]float64{
					shared.Teams["Team1"]: 50,
					shared.Teams["Team2"]: 60,
					shared.Teams["Team4"]: 70,
					shared.Teams["Team5"]: 80,
					shared.Teams["Team6"]: 90,
				},
				requestedGiftAmounts: map[shared.ClientID]shared.GiftRequest{
					shared.Teams["Team1"]: 43.87594114979616,
					shared.Teams["Team2"]: 45.92210178127672,
					shared.Teams["Team4"]: 47.726375540564796,
					shared.Teams["Team5"]: 49.34650978030029,
					shared.Teams["Team6"]: 50.821159755976886,
				},
				initialResourcesAtStartOfGame: 100,
			},
			expectedVal: shared.GiftRequestDict{
				shared.Teams["Team1"]: 36.08094844012818,
				shared.Teams["Team2"]: 37.4208970618899,
				shared.Teams["Team4"]: 38.59255681756235,
				shared.Teams["Team5"]: 39.637106321387236,
				shared.Teams["Team6"]: 40.58190651651451,
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.ourClient.BaseClient.ServerReadHandle = tc.serverReadHandle
			res := tc.ourClient.GetGiftRequests()
			if !reflect.DeepEqual(res, tc.expectedVal) {
				t.Errorf("Expected final transgressions to be %v got %v", tc.expectedVal, res)
				if !reflect.DeepEqual(tc.ourClient.requestedGiftAmounts, tc.expectedVal) {
					t.Errorf("Expected internal resources map to be %v got %v", tc.expectedVal, tc.ourClient.requestedGiftAmounts)
				}
			}
		})
	}
}

func TestGetGiftOffers(t *testing.T) {
	cases := []struct {
		name             string
		ourClient        client
		receivedRequests shared.GiftRequestDict
		giftBudget       shared.GiftOffer
		expectedVal      shared.GiftOfferDict
	}{
		{
			name: "Basic test: our island is critical so offers nothing to all requested islands",
			ourClient: client{
				BaseClient: &baseclient.BaseClient{
					ServerReadHandle: mockServerReadHandle{gameState: gamestate.ClientGameState{
						ClientLifeStatuses: map[shared.ClientID]shared.ClientLifeStatus{
							shared.Teams["Team1"]: shared.Alive,
							shared.Teams["Team2"]: shared.Alive,
							shared.Teams["Team3"]: shared.Critical,
							shared.Teams["Team4"]: shared.Alive,
							shared.Teams["Team5"]: shared.Alive,
							shared.Teams["Team6"]: shared.Alive,
						},
						ClientInfo: gamestate.ClientInfo{Resources: 600.0}}}},
				params: islandParams{
					giftInflationPercentage: 0.1,
					selfishness:             0.3,
				},
				trustScore: map[shared.ClientID]float64{
					shared.Teams["Team1"]: 50,
					shared.Teams["Team2"]: 50,
					shared.Teams["Team4"]: 50,
					shared.Teams["Team5"]: 50,
					shared.Teams["Team6"]: 50,
				},
				requestedGiftAmounts: map[shared.ClientID]shared.GiftRequest{
					shared.Teams["Team1"]: 15,
					shared.Teams["Team2"]: 15,
					shared.Teams["Team4"]: 15,
					shared.Teams["Team5"]: 15,
					shared.Teams["Team6"]: 15,
				},
				initialResourcesAtStartOfGame: 150,
			},
			receivedRequests: shared.GiftRequestDict{
				shared.Teams["Team1"]: 50,
				shared.Teams["Team2"]: 500,
				shared.Teams["Team4"]: 50,
			},
			giftBudget: 376.25,
			expectedVal: shared.GiftOfferDict{
				shared.Teams["Team1"]: 0,
				shared.Teams["Team2"]: 0,
				shared.Teams["Team4"]: 0,
				shared.Teams["Team5"]: 0,
				shared.Teams["Team6"]: 0,
			},
		},

		{
			name: "Basic test: all islands are alive, trusted equally and we send offers",
			ourClient: client{
				BaseClient: &baseclient.BaseClient{
					ServerReadHandle: mockServerReadHandle{gameState: gamestate.ClientGameState{
						ClientLifeStatuses: map[shared.ClientID]shared.ClientLifeStatus{
							shared.Teams["Team1"]: shared.Alive,
							shared.Teams["Team2"]: shared.Alive,
							shared.Teams["Team3"]: shared.Alive,
							shared.Teams["Team4"]: shared.Alive,
							shared.Teams["Team5"]: shared.Alive,
							shared.Teams["Team6"]: shared.Alive,
						},
						ClientInfo: gamestate.ClientInfo{Resources: 1000.0}}}},
				params: islandParams{
					giftInflationPercentage: 0.1,
					selfishness:             0.3,
				},
				trustScore: map[shared.ClientID]float64{
					shared.Teams["Team1"]: 50,
					shared.Teams["Team2"]: 60,
					shared.Teams["Team4"]: 75,
					shared.Teams["Team5"]: 85,
					shared.Teams["Team6"]: 100,
				},
				requestedGiftAmounts: map[shared.ClientID]shared.GiftRequest{
					shared.Teams["Team1"]: 15,
					shared.Teams["Team2"]: 15,
					shared.Teams["Team4"]: 15,
					shared.Teams["Team5"]: 15,
					shared.Teams["Team6"]: 15,
				},
				initialResourcesAtStartOfGame: 150,
			},
			receivedRequests: shared.GiftRequestDict{
				shared.Teams["Team1"]: 100,
				shared.Teams["Team2"]: 100,
				shared.Teams["Team4"]: 100,
				shared.Teams["Team5"]: 100,
				shared.Teams["Team6"]: 100,
			},
			giftBudget: 376.25,
			expectedVal: shared.GiftOfferDict{
				shared.Teams["Team1"]: 0,
				shared.Teams["Team2"]: 0,
				shared.Teams["Team4"]: 0,
				shared.Teams["Team5"]: 0,
				shared.Teams["Team6"]: 0,
			},
		},

		{
			name: "Basic test: all islands are alive, trusted differently and we send offers",
			ourClient: client{
				BaseClient: &baseclient.BaseClient{
					ServerReadHandle: mockServerReadHandle{gameState: gamestate.ClientGameState{
						ClientLifeStatuses: map[shared.ClientID]shared.ClientLifeStatus{
							shared.Teams["Team1"]: shared.Alive,
							shared.Teams["Team2"]: shared.Alive,
							shared.Teams["Team3"]: shared.Alive,
							shared.Teams["Team4"]: shared.Alive,
							shared.Teams["Team5"]: shared.Alive,
							shared.Teams["Team6"]: shared.Alive,
						},
						ClientInfo: gamestate.ClientInfo{Resources: 1000.0}}}},
				params: islandParams{
					giftInflationPercentage: 0.1,
					selfishness:             0.3,
				},
				trustScore: map[shared.ClientID]float64{
					shared.Teams["Team1"]: 50,
					shared.Teams["Team2"]: 50,
					shared.Teams["Team4"]: 80,
					shared.Teams["Team5"]: 50,
					shared.Teams["Team6"]: 20,
				},
				requestedGiftAmounts: map[shared.ClientID]shared.GiftRequest{
					shared.Teams["Team1"]: 15,
					shared.Teams["Team2"]: 15,
					shared.Teams["Team4"]: 15,
					shared.Teams["Team5"]: 15,
					shared.Teams["Team6"]: 15,
				},
				initialResourcesAtStartOfGame: 150,
			},
			receivedRequests: shared.GiftRequestDict{
				shared.Teams["Team1"]: 100,
				shared.Teams["Team2"]: 100,
				shared.Teams["Team4"]: 100,
				shared.Teams["Team5"]: 100,
				shared.Teams["Team6"]: 100,
			},
			giftBudget: 376.25,
			expectedVal: shared.GiftOfferDict{
				shared.Teams["Team1"]: 0,
				shared.Teams["Team2"]: 0,
				shared.Teams["Team4"]: 0,
				shared.Teams["Team5"]: 0,
				shared.Teams["Team6"]: 0,
			},
		},

		{
			name: "Complex test: all islands are alive, trusted differently but we send offers to all including those that did not",
			ourClient: client{
				BaseClient: &baseclient.BaseClient{
					ServerReadHandle: mockServerReadHandle{gameState: gamestate.ClientGameState{
						ClientLifeStatuses: map[shared.ClientID]shared.ClientLifeStatus{
							shared.Teams["Team1"]: shared.Alive,
							shared.Teams["Team2"]: shared.Alive,
							shared.Teams["Team3"]: shared.Alive,
							shared.Teams["Team4"]: shared.Alive,
							shared.Teams["Team5"]: shared.Alive,
							shared.Teams["Team6"]: shared.Alive,
						},
						ClientInfo: gamestate.ClientInfo{Resources: 1000.0}}}},
				params: islandParams{
					giftInflationPercentage: 0.1,
					selfishness:             0.3,
				},
				trustScore: map[shared.ClientID]float64{
					shared.Teams["Team1"]: 50,
					shared.Teams["Team2"]: 50,
					shared.Teams["Team4"]: 80,
					shared.Teams["Team5"]: 50,
					shared.Teams["Team6"]: 20,
				},
				requestedGiftAmounts: map[shared.ClientID]shared.GiftRequest{
					shared.Teams["Team1"]: 15,
					shared.Teams["Team2"]: 15,
					shared.Teams["Team4"]: 15,
					shared.Teams["Team5"]: 15,
					shared.Teams["Team6"]: 15,
				},
				initialResourcesAtStartOfGame: 150,
			},
			receivedRequests: shared.GiftRequestDict{
				shared.Teams["Team1"]: 100,
				shared.Teams["Team4"]: 100,
				shared.Teams["Team6"]: 100,
			},
			giftBudget: 376.25,
			expectedVal: shared.GiftOfferDict{
				shared.Teams["Team1"]: 0,
				shared.Teams["Team2"]: 0,
				shared.Teams["Team4"]: 0,
				shared.Teams["Team5"]: 0,
				shared.Teams["Team6"]: 0,
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var sum shared.GiftOffer
			res := tc.ourClient.GetGiftOffers(tc.receivedRequests)
			for _, ans := range res {
				sum += ans
			}
			if sum > tc.giftBudget {
				t.Errorf("Total sum of offered gifts (%f) is greater gift budget (%f)", sum, tc.giftBudget)
			}
		})
	}
}

func TestGetGiftResponses(t *testing.T) {
	cases := []struct {
		name           string
		ourClient      client
		receivedOffers shared.GiftOfferDict
		expectedVal    shared.GiftResponseDict
	}{
		{
			name: "Basic test: all islands are alive so all offers are accepted",
			ourClient: client{
				BaseClient: &baseclient.BaseClient{
					ServerReadHandle: mockServerReadHandle{gameState: gamestate.ClientGameState{
						ClientLifeStatuses: map[shared.ClientID]shared.ClientLifeStatus{
							shared.Teams["Team1"]: shared.Alive,
							shared.Teams["Team2"]: shared.Alive,
							shared.Teams["Team3"]: shared.Alive,
							shared.Teams["Team4"]: shared.Alive,
							shared.Teams["Team5"]: shared.Alive,
							shared.Teams["Team6"]: shared.Alive,
						},
						ClientInfo: gamestate.ClientInfo{Resources: 600.0}}}},
			},
			receivedOffers: shared.GiftOfferDict{
				shared.Teams["Team1"]: 50,
				shared.Teams["Team2"]: 500,
				shared.Teams["Team4"]: 50,
				shared.Teams["Team5"]: 100.78,
				shared.Teams["Team6"]: 78.987,
			},
			expectedVal: shared.GiftResponseDict{
				shared.Teams["Team1"]: {AcceptedAmount: 50, Reason: 0},
				shared.Teams["Team2"]: {AcceptedAmount: 500, Reason: 0},
				shared.Teams["Team4"]: {AcceptedAmount: 50, Reason: 0},
				shared.Teams["Team5"]: {AcceptedAmount: 100.78, Reason: 0},
				shared.Teams["Team6"]: {AcceptedAmount: 78.987, Reason: 0},
			},
		},

		{
			name: "Basic test: few island are critical so their offers are not accepted",
			ourClient: client{
				BaseClient: &baseclient.BaseClient{
					ServerReadHandle: mockServerReadHandle{gameState: gamestate.ClientGameState{
						ClientLifeStatuses: map[shared.ClientID]shared.ClientLifeStatus{
							shared.Teams["Team1"]: shared.Alive,
							shared.Teams["Team2"]: shared.Critical,
							shared.Teams["Team3"]: shared.Alive,
							shared.Teams["Team4"]: shared.Critical,
							shared.Teams["Team5"]: shared.Alive,
							shared.Teams["Team6"]: shared.Critical,
						},
						ClientInfo: gamestate.ClientInfo{Resources: 600.0}}}},
			},
			receivedOffers: shared.GiftOfferDict{
				shared.Teams["Team1"]: 50,
				shared.Teams["Team2"]: 500,
				shared.Teams["Team4"]: 50,
				shared.Teams["Team5"]: 100.78,
				shared.Teams["Team6"]: 78.987,
			},
			expectedVal: shared.GiftResponseDict{
				shared.Teams["Team1"]: {AcceptedAmount: 50, Reason: 0},
				shared.Teams["Team2"]: {AcceptedAmount: 0, Reason: 1},
				shared.Teams["Team4"]: {AcceptedAmount: 0, Reason: 1},
				shared.Teams["Team5"]: {AcceptedAmount: 100.78, Reason: 0},
				shared.Teams["Team6"]: {AcceptedAmount: 0, Reason: 1},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res := tc.ourClient.GetGiftResponses(tc.receivedOffers)
			if !reflect.DeepEqual(res, tc.expectedVal) {
				t.Errorf("Expected final transgressions to be %v got %v", tc.expectedVal, res)
			}
		})
	}
}

func TestUpdateGiftInfo(t *testing.T) {
	cases := []struct {
		name              string
		ourClient         client
		receivedResponses shared.GiftResponseDict
		expectedVal       map[shared.ClientID][]float64
	}{
		{
			name: "Basic test: all islands are alive and all requested amounts are complied to",
			ourClient: client{
				BaseClient: &baseclient.BaseClient{
					ServerReadHandle: mockServerReadHandle{gameState: gamestate.ClientGameState{
						ClientLifeStatuses: map[shared.ClientID]shared.ClientLifeStatus{
							shared.Teams["Team1"]: shared.Alive,
							shared.Teams["Team2"]: shared.Alive,
							shared.Teams["Team3"]: shared.Alive,
							shared.Teams["Team4"]: shared.Alive,
							shared.Teams["Team5"]: shared.Alive,
							shared.Teams["Team6"]: shared.Alive,
						},
						ClientInfo: gamestate.ClientInfo{Resources: 600.0}}}},
				trustMapAgg: map[shared.ClientID][]float64{
					0: {},
					1: {},
					3: {},
					4: {},
					5: {},
				},
				requestedGiftAmounts: map[shared.ClientID]shared.GiftRequest{
					shared.Teams["Team1"]: 40,
					shared.Teams["Team2"]: 50,
					shared.Teams["Team4"]: 60,
					shared.Teams["Team5"]: 70,
					shared.Teams["Team6"]: 80,
				},
			},
			receivedResponses: shared.GiftResponseDict{
				shared.Teams["Team1"]: {AcceptedAmount: 40, Reason: 0},
				shared.Teams["Team2"]: {AcceptedAmount: 50, Reason: 0},
				shared.Teams["Team4"]: {AcceptedAmount: 60, Reason: 0},
				shared.Teams["Team5"]: {AcceptedAmount: 70, Reason: 0},
				shared.Teams["Team6"]: {AcceptedAmount: 80, Reason: 0},
			},
			expectedVal: map[shared.ClientID][]float64{
				0: {},
				1: {},
				3: {},
				4: {},
				5: {},
			},
		},

		{
			name: "Basic test: all islands are alive but different decline reasons given",
			ourClient: client{
				BaseClient: &baseclient.BaseClient{
					ServerReadHandle: mockServerReadHandle{gameState: gamestate.ClientGameState{
						ClientLifeStatuses: map[shared.ClientID]shared.ClientLifeStatus{
							shared.Teams["Team1"]: shared.Alive,
							shared.Teams["Team2"]: shared.Alive,
							shared.Teams["Team3"]: shared.Alive,
							shared.Teams["Team4"]: shared.Alive,
							shared.Teams["Team5"]: shared.Alive,
							shared.Teams["Team6"]: shared.Alive,
						},
						ClientInfo: gamestate.ClientInfo{Resources: 600.0}}}},
				trustMapAgg: map[shared.ClientID][]float64{
					0: {},
					1: {},
					3: {},
					4: {},
					5: {},
				},
				requestedGiftAmounts: map[shared.ClientID]shared.GiftRequest{
					shared.Teams["Team1"]: 40,
					shared.Teams["Team2"]: 40,
					shared.Teams["Team4"]: 40,
					shared.Teams["Team5"]: 40,
					shared.Teams["Team6"]: 40,
				},
			},
			receivedResponses: shared.GiftResponseDict{
				shared.Teams["Team1"]: {AcceptedAmount: 50, Reason: 0},
				shared.Teams["Team2"]: {AcceptedAmount: 500, Reason: 0},
				shared.Teams["Team4"]: {AcceptedAmount: 30, Reason: 2},
				shared.Teams["Team5"]: {AcceptedAmount: 39.999, Reason: 2},
				shared.Teams["Team6"]: {AcceptedAmount: 40.001, Reason: 0},
			},
			expectedVal: map[shared.ClientID][]float64{
				0: {},
				1: {},
				3: {-5},
				4: {-5},
				5: {},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.ourClient.UpdateGiftInfo(tc.receivedResponses)
			if !reflect.DeepEqual(tc.ourClient.trustMapAgg, tc.expectedVal) {
				t.Errorf("Expected final transgressions to be %v got %v", tc.expectedVal, tc.ourClient.trustMapAgg)
			}
		})
	}
}

func TestReceivedGift(t *testing.T) {
	cases := []struct {
		name        string
		ourClient   client
		received    shared.Resources
		from        shared.ClientID
		expectedVal map[shared.ClientID][]float64
	}{
		{
			name: "Basic test: all islands are alive and received more than requested",
			ourClient: client{
				BaseClient: &baseclient.BaseClient{
					ServerReadHandle: mockServerReadHandle{gameState: gamestate.ClientGameState{
						ClientLifeStatuses: map[shared.ClientID]shared.ClientLifeStatus{
							shared.Teams["Team1"]: shared.Alive,
							shared.Teams["Team2"]: shared.Alive,
							shared.Teams["Team3"]: shared.Alive,
							shared.Teams["Team4"]: shared.Alive,
							shared.Teams["Team5"]: shared.Alive,
							shared.Teams["Team6"]: shared.Alive,
						},
						ClientInfo: gamestate.ClientInfo{Resources: 600.0}}}},
				trustMapAgg: map[shared.ClientID][]float64{
					0: {},
					1: {},
					3: {},
					4: {},
					5: {},
				},
				requestedGiftAmounts: map[shared.ClientID]shared.GiftRequest{
					shared.Teams["Team1"]: 40,
					shared.Teams["Team2"]: 50,
					shared.Teams["Team4"]: 60,
					shared.Teams["Team5"]: 70,
					shared.Teams["Team6"]: 80,
				},
			},
			received: 100,
			from:     shared.Teams["Team1"],
			expectedVal: map[shared.ClientID][]float64{
				0: {22},
				1: {},
				3: {},
				4: {},
				5: {},
			},
		},

		{
			name: "Basic test: some islands are critical and hence their trust score not adjusted",
			ourClient: client{
				BaseClient: &baseclient.BaseClient{
					ServerReadHandle: mockServerReadHandle{gameState: gamestate.ClientGameState{
						ClientLifeStatuses: map[shared.ClientID]shared.ClientLifeStatus{
							shared.Teams["Team1"]: shared.Critical,
							shared.Teams["Team2"]: shared.Alive,
							shared.Teams["Team3"]: shared.Alive,
							shared.Teams["Team4"]: shared.Alive,
							shared.Teams["Team5"]: shared.Alive,
							shared.Teams["Team6"]: shared.Critical,
						},
						ClientInfo: gamestate.ClientInfo{Resources: 600.0}}}},
				trustMapAgg: map[shared.ClientID][]float64{
					0: {},
					1: {},
					3: {},
					4: {},
					5: {},
				},
				requestedGiftAmounts: map[shared.ClientID]shared.GiftRequest{
					shared.Teams["Team1"]: 40,
					shared.Teams["Team2"]: 50,
					shared.Teams["Team4"]: 60,
					shared.Teams["Team5"]: 70,
					shared.Teams["Team6"]: 80,
				},
			},
			received: 40,
			from:     shared.Teams["Team4"],
			expectedVal: map[shared.ClientID][]float64{
				0: {},
				1: {},
				3: {-6},
				4: {},
				5: {},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tc.ourClient.ReceivedGift(tc.received, tc.from)
			if !reflect.DeepEqual(tc.ourClient.trustMapAgg, tc.expectedVal) {
				t.Errorf("Expected final transgressions to be %v got %v", tc.expectedVal, tc.ourClient.trustMapAgg)
			}
		})
	}
}
