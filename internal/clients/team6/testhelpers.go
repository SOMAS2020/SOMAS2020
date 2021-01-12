package team6

// some tools for testing

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type mockInit struct {
	serverReadHandle      stubServerReadHandle
	friendship            Friendship
	trustRank             TrustRank
	giftsSentHistory      GiftsSentHistory
	giftsReceivedHistory  GiftsReceivedHistory
	giftsRequestedHistory GiftsRequestedHistory
	disastersHistory      DisastersHistory
	disasterPredictions   DisasterPredictions
	forageHistory         ForageHistory
	taxDemanded           shared.Resources
}

func newMockClient(clientID shared.ClientID, init mockInit) client {
	mockClient := client{
		BaseClient:   baseclient.NewClient(clientID),
		clientConfig: getClientConfig(),
	}

	mockClient.ServerReadHandle = init.serverReadHandle
	mockClient.friendship = init.friendship
	mockClient.trustRank = init.trustRank
	mockClient.giftsSentHistory = init.giftsSentHistory
	mockClient.giftsReceivedHistory = init.giftsReceivedHistory
	mockClient.giftsRequestedHistory = init.giftsRequestedHistory
	mockClient.disastersHistory = init.disastersHistory
	mockClient.disasterPredictions = init.disasterPredictions
	mockClient.forageHistory = init.forageHistory
	mockClient.taxDemanded = init.taxDemanded

	return mockClient
}

type stubServerReadHandle struct {
	gameState  gamestate.ClientGameState
	gameConfig config.ClientConfig
}

func (s stubServerReadHandle) GetGameState() gamestate.ClientGameState {
	return s.gameState
}
func (s stubServerReadHandle) GetGameConfig() config.ClientConfig {
	return s.gameConfig
}
