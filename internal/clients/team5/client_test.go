package team5

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type testServerHandle struct {
	clientGameState  gamestate.ClientGameState
	clientGameConfig config.ClientConfig
}

func (h testServerHandle) GetGameState() gamestate.ClientGameState {
	return h.clientGameState
}

func (h testServerHandle) GetGameConfig() config.ClientConfig {
	return h.clientGameConfig
}

func MakeTestClient(gamestate gamestate.ClientGameState) client {
	c := NewTestClient(shared.Team5)
	c.Initialise(testServerHandle{
		clientGameState: gamestate,
	})
	return *c.(*client)
}

// NewTestClient is a client for testing purposes
func NewTestClient(clientID shared.ClientID) baseclient.Client {
	return &client{
		BaseClient:              baseclient.NewClient(clientID),
		cpRequestHistory:        cpRequestHistory{},
		cpAllocationHistory:     cpAllocationHistory{},
		forageHistory:           forageHistory{},
		resourceHistory:         resourceHistory{},
		team5President:          president{},
		giftHistory:             map[shared.ClientID]giftExchange{},
		forecastHistory:         forecastHistory{},
		receivedForecastHistory: receivedForecastHistory{},
		disasterHistory:         disasterHistory{},

		config: getClientConfig(),
	}
}
