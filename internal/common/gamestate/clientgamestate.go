package gamestate

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/disasters"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

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

	// Island Locations
	Geography disasters.ArchipelagoGeography

	// Islands holding IIGO roles
	SpeakerID   shared.ClientID
	JudgeID     shared.ClientID
	PresidentID shared.ClientID

	// IIGO roles budget (initialised in orchestration.go)
	IIGORolesBudget map[shared.Role]shared.Resources

	// IIGO turns in power (incremented and set by monitoring)
	IIGOTurnsInPower map[shared.Role]uint

	// RuleInfo contains the global rules information for clients to access
	RulesInfo RulesContext
}
