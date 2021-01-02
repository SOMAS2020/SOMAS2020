package gamestate

import "github.com/SOMAS2020/SOMAS2020/internal/common/shared"

// ClientGameState contains game state only for a specific client.
type ClientGameState struct {
	// Season represents the current (1-index) season of the game.
	Season uint

	// Turn represents the current (1-index) Turn of the game.
	Turn uint

	// ClientInfo
	ClientInfo ClientInfo

	// AliveStatuses is whether each of the other clients is alive
	ClientLifeStatuses map[shared.ClientID]shared.ClientLifeStatus

	// CommonPool
	CommonPool shared.Resources
}
