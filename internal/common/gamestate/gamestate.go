// Package gamestate contains information about the current game state.
package gamestate

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/disasters"
	"github.com/SOMAS2020/SOMAS2020/internal/common/foraging"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// GameState represents the game's state.
type GameState struct {
	// Season represents the current (1-index) season of the game.
	Season uint

	// Turn represents the current (1-index) Turn of the game.
	Turn uint

	// CommonPool represents the amount of resources in the common pool.
	CommonPool shared.Resources

	// ClientInfos map from the shared.ClientID to ClientInfo.
	ClientInfos    map[shared.ClientID]ClientInfo
	Environment    disasters.Environment
	DeerPopulation foraging.DeerPopulationModel

	// Foraging History
	ForagingHistory map[shared.ForageType][]foraging.ForagingReport

	// IIGO History: indexed by turn
	IIGOHistory map[uint][]shared.Accountability

	// IIGO roles budget (initialised in orchestration.go)
	IIGORolesBudget map[shared.Role]shared.Resources
	// IITO Transactions
	IITOTransactions map[shared.ClientID]shared.GiftResponseDict

	// Orchestration
	SpeakerID   shared.ClientID
	JudgeID     shared.ClientID
	PresidentID shared.ClientID

	// [INFRA] add more details regarding state of game here
	// REMEMBER TO EDIT `Copy` IF YOU ADD ANY REFERENCE TYPES (maps, slices, channels, functions etc.)
}

// Copy returns a deep copy of the GameState.
func (g GameState) Copy() GameState {
	ret := g
	ret.ClientInfos = copyClientInfos(g.ClientInfos)
	ret.Environment = g.Environment.Copy()
	ret.DeerPopulation = g.DeerPopulation.Copy()
	ret.ForagingHistory = copyForagingHistory(g.ForagingHistory)
	ret.IIGOHistory = copyIIGOHistory(g.IIGOHistory)
	ret.IIGORolesBudget = copyRolesBudget(g.IIGORolesBudget)
	return ret
}

// GetClientGameStateCopy returns the ClientGameState for the client having the id.
func (g *GameState) GetClientGameStateCopy(id shared.ClientID) ClientGameState {
	clientLifeStatuses := map[shared.ClientID]shared.ClientLifeStatus{}
	for id, clientInfo := range g.ClientInfos {
		clientLifeStatuses[id] = clientInfo.LifeStatus
	}

	return ClientGameState{
		Season:             g.Season,
		Turn:               g.Turn,
		ClientInfo:         g.ClientInfos[id].Copy(),
		ClientLifeStatuses: clientLifeStatuses,
		CommonPool:         g.CommonPool,
		IslandLocations:    g.Environment.Geography.Islands,
		SpeakerID:          g.SpeakerID,
		JudgeID:            g.JudgeID,
		PresidentID:        g.PresidentID,
		IIGORolesBudget:    copyRolesBudget(g.IIGORolesBudget),
	}
}

func copyClientInfos(m map[shared.ClientID]ClientInfo) map[shared.ClientID]ClientInfo {
	ret := make(map[shared.ClientID]ClientInfo, len(m))
	for k, v := range m {
		ret[k] = v.Copy()
	}
	return ret
}

func copyRolesBudget(m map[shared.Role]shared.Resources) map[shared.Role]shared.Resources {
	ret := make(map[shared.Role]shared.Resources, len(m))
	for k, v := range m {
		ret[k] = v
	}
	return ret
}

func copyIIGOHistory(iigoHistory map[uint][]shared.Accountability) map[uint][]shared.Accountability {
	targetMap := make(map[uint][]shared.Accountability)
	for key, value := range iigoHistory {
		targetMap[key] = copySingleIIGOEntry(value)
	}
	return targetMap
}

func copySingleIIGOEntry(input []shared.Accountability) []shared.Accountability {
	ret := make([]shared.Accountability, len(input))
	copy(ret, input)
	return ret
}

func copyForagingHistory(fHist map[shared.ForageType][]foraging.ForagingReport) map[shared.ForageType][]foraging.ForagingReport {
	ret := make(map[shared.ForageType][]foraging.ForagingReport, len(fHist))
	for k, v := range fHist { // iterate over different foraging types
		ret[k] = make([]foraging.ForagingReport, len(v))
		for i, el := range v { // iterate over reports across time for given foraging type
			ret[k][i] = el.Copy()
		}
	}
	return ret
}

// ClientInfo contains the client struct as well as the client's attributes
type ClientInfo struct {
	// Resources contains the amount of resources owned by the client.
	Resources shared.Resources // nice

	LifeStatus shared.ClientLifeStatus
	// CriticalConsecutiveTurnsCounter is the number of consecutive turns the client is in critical state.
	// Client will die if this Counter reaches config.MaxCriticalConsecutiveTurns
	CriticalConsecutiveTurnsCounter uint

	// [INFRA] add more client information here
	// REMEMBER TO EDIT `Copy` IF YOU ADD ANY REFERENCE TYPES (maps, slices, channels, functions etc.)
}

// Copy returns a deep copy of the ClientInfo.
func (c ClientInfo) Copy() ClientInfo {
	ret := c
	return ret
}
