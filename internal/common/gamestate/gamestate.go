// Package gamestate contains information about the current game state.
package gamestate

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/disasters"
	"github.com/SOMAS2020/SOMAS2020/internal/common/foraging"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
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

	// All global rules information
	RulesInfo RulesContext

	// IIGO History: indexed by turn
	IIGOHistory map[uint][]shared.Accountability

	// IIGO roles budget (initialised in orchestration.go)
	IIGORolesBudget map[shared.Role]shared.Resources

	// IIGO turns in power (incremented and set by monitoring)
	IIGOTurnsInPower map[shared.Role]uint

	// IIGO Tax Amount Map
	IIGOTaxAmount map[shared.ClientID]shared.Resources

	// IIGO Allocation Amount Map
	IIGOAllocationMap map[shared.ClientID]shared.Resources

	// IIGO Allocation Made
	IIGOAllocationMade bool

	// IIGO Sanction Amount Map
	IIGOSanctionMap map[shared.ClientID]shared.Resources

	// IIGO Sanction Cache is a record of all sanctions
	IIGOSanctionCache map[int][]shared.Sanction

	// IIGO History Cache is a record of all transgressions for historical retribution
	IIGOHistoryCache map[int][]shared.Accountability

	// IIGO Role Monitoring Cache is used in the IIGO accountability cycle
	IIGORoleMonitoringCache []shared.Accountability

	// Rules broken by islands in this turn
	RulesBrokenByIslands map[shared.ClientID][]string

	// Internal IIGO rules broken by IIGO roles in this turn
	IIGORulesBrokenByRoles map[shared.Role][]string

	// IIGO Role Voting
	IIGOElection []VotingInfo

	// IIGO Run Status
	IIGORunStatus string

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
	ret.RulesInfo = copyRulesContext(g.RulesInfo)
	ret.IIGOHistory = copyIIGOHistory(g.IIGOHistory)
	ret.IIGORolesBudget = copyRolesBudget(g.IIGORolesBudget)
	ret.IIGOTurnsInPower = copyTurnsInPower(g.IIGOTurnsInPower)
	ret.IIGOTaxAmount = copyIIGOClientIDResourceMap(g.IIGOTaxAmount)
	ret.IIGOAllocationMap = copyIIGOClientIDResourceMap(g.IIGOAllocationMap)
	ret.IIGOSanctionMap = copyIIGOClientIDResourceMap(g.IIGOSanctionMap)
	ret.IIGOHistoryCache = copyIIGOHistoryCache(g.IIGOHistoryCache)
	ret.IIGOSanctionCache = copyIIGOSanctionCache(g.IIGOSanctionCache)
	ret.IIGORoleMonitoringCache = copySingleIIGOEntry(g.IIGORoleMonitoringCache)
	ret.IITOTransactions = copyIITOTransactions(g.IITOTransactions)
	ret.IIGOElection = copyIIGOElection(g.IIGOElection)
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
		Geography:          g.Environment.Geography,
		SpeakerID:          g.SpeakerID,
		JudgeID:            g.JudgeID,
		PresidentID:        g.PresidentID,
		IIGORolesBudget:    copyRolesBudget(g.IIGORolesBudget),
		IIGOTurnsInPower:   copyTurnsInPower(g.IIGOTurnsInPower),
		RulesInfo:          copyRulesContext(g.RulesInfo),
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

func copyTurnsInPower(m map[shared.Role]uint) map[shared.Role]uint {
	ret := make(map[shared.Role]uint, len(m))
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

func copyIIGOSanctionCache(m map[int][]shared.Sanction) map[int][]shared.Sanction {
	targetMap := make(map[int][]shared.Sanction)
	for k, v := range m {
		entry := make([]shared.Sanction, len(v))
		copy(entry, v)
		targetMap[k] = entry
	}
	return targetMap
}

func copyIIGOHistoryCache(iigoHistory map[int][]shared.Accountability) map[int][]shared.Accountability {
	targetMap := make(map[int][]shared.Accountability)
	for key, value := range iigoHistory {
		targetMap[key] = copySingleIIGOEntry(value)
	}
	return targetMap
}

func copyRulesContext(oldContext RulesContext) RulesContext {
	return RulesContext{
		VariableMap:        copyVariableMap(oldContext.VariableMap),
		CurrentRulesInPlay: rules.CopyRulesMap(oldContext.CurrentRulesInPlay),
		AvailableRules:     rules.CopyRulesMap(oldContext.AvailableRules),
	}
}

func copyVariableMap(varMap map[rules.VariableFieldName]rules.VariableValuePair) map[rules.VariableFieldName]rules.VariableValuePair {
	targetMap := make(map[rules.VariableFieldName]rules.VariableValuePair)
	for key, value := range varMap {
		targetMap[key] = rules.VariableValuePair{
			VariableName: value.VariableName,
			Values:       value.Values,
		}
	}
	return targetMap
}

func copySingleIIGOEntry(input []shared.Accountability) []shared.Accountability {
	ret := make([]shared.Accountability, len(input))
	copy(ret, input)
	return ret
}

func copyIIGOClientIDResourceMap(m map[shared.ClientID]shared.Resources) map[shared.ClientID]shared.Resources {
	ret := make(map[shared.ClientID]shared.Resources, len(m))
	for k, v := range m {
		ret[k] = v
	}
	return ret
}

func copyIITOTransactions(iitoHistory map[shared.ClientID]shared.GiftResponseDict) map[shared.ClientID]shared.GiftResponseDict {
	targetMap := make(map[shared.ClientID]shared.GiftResponseDict)
	for key, giftResponseDict := range iitoHistory {
		giftResponses := make(shared.GiftResponseDict)
		for clientID, giftResponse := range giftResponseDict {
			giftResponses[clientID] = giftResponse
		}
		targetMap[key] = giftResponses
	}
	return targetMap
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

func copyIIGOElection(input []VotingInfo) []VotingInfo {
	ret := make([]VotingInfo, len(input))
	copy(ret, input)
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

// VotingInfo contains all the information necessary to visualise voting
type VotingInfo struct {
	RoleToElect  shared.Role
	VotingMethod shared.ElectionVotingMethod
	VoterList    []shared.ClientID
	Votes        [][]shared.ClientID
}

// Copy returns a deep copy of the ClientInfo.
func (c ClientInfo) Copy() ClientInfo {
	ret := c
	return ret
}

// RulesContext contains the global rules information at any given time
type RulesContext struct {
	// Map of global variable values
	VariableMap map[rules.VariableFieldName]rules.VariableValuePair

	//Map of All available rules
	AvailableRules map[string]rules.RuleMatrix

	// Rules Currently In Play
	CurrentRulesInPlay map[string]rules.RuleMatrix
}
