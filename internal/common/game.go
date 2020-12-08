package common

import (
	"fmt"
	"log"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// ClientInfo contains the client struct as well as the client's attributes
type ClientInfo struct {
	// Resources contains the amount of resources owned by the client.
	// Made an integer so an island can "owe" resources.
	Resources int

	// Normal Condition: Alive = true, Critical = false, CriticalConsecutiveTurnsLeft = config.MaxCriticalConsecutiveTurnsLeft
	// Critical Condition: Alive = true, Critical = false, CriticalConsecutiveTurnsLeft <= config.MaxCriticalConsecutiveTurnsLeft
	// Dead: Alive = false
	Alive                        bool
	Critical                     bool
	CriticalConsecutiveTurnsLeft uint

	// [INFRA] add more client information here
	// REMEMBER TO EDIT `Copy` IF YOU ADD ANY REFERENCE TYPES (maps, slices, channels, functions etc.)
}

// Copy returns a deep copy of the ClientInfo.
func (c ClientInfo) Copy() ClientInfo {
	ret := c
	return ret
}

// GameState represents the game's state.
type GameState struct {
	// Season represents the current (1-index) season of the game.
	Season uint

	// Turn represents the current (1-index) Turn of the game.
	Turn uint

	// ClientInfos map from the shared.ClientID to ClientInfo.
	// EXTRA note: Golang maps are made to be random!
	ClientInfos map[shared.ClientID]ClientInfo

	// [INFRA] add more details regarding state of game here
	// REMEMBER TO EDIT `Copy` IF YOU ADD ANY REFERENCE TYPES (maps, slices, channels, functions etc.)
}

// Copy returns a deep copy of the GameState.
func (g GameState) Copy() GameState {
	ret := g
	ret.ClientInfos = copyClientInfos(g.ClientInfos)
	return ret
}

func copyClientInfos(m map[shared.ClientID]ClientInfo) map[shared.ClientID]ClientInfo {
	ret := make(map[shared.ClientID]ClientInfo, len(m))
	for k, v := range m {
		ret[k] = v.Copy()
	}
	return ret
}

func (g GameState) logf(format string, a ...interface{}) {
	log.Printf("[GAMESTATE]: %v", fmt.Sprintf(format, a...))
}
