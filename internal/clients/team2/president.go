package team2

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// const Attitude int

// const (
// 	FreeRider Attitude = iota
// 	FairSharer
// 	Altruist
// )

type President struct {
	*baseclient.BasePresident
	c *client
}

func (p *President) EvaluateAllocationRequests(resourceRequest map[shared.ClientID]shared.Resources, availCommonPool shared.Resources) shared.PresidentReturnContent {
	var modeMult, requestSum shared.Resources
	resourceAllocation := make(map[shared.ClientID]shared.Resources)

	for _, request := range resourceRequest {
		requestSum += request
	}
	switch p.MethodOfPlay() {
	case 0:
		modeMult = 0.6 //pool is struggling, give less
	case 1:
		modeMult = 0.7 //default
	case 2:
		modeMult = 0.8 //pool has surplus, give a bit more than usual
	}

	if requestSum < modeMult*availCommonPool || requestSum == 0 { //enough resources to go round
		resourceAllocation = resourceRequest
	} else { // In this case we only give out resources to critical islands
		//order islands in terms of most critical to least critical
		resourcesToGive := modeMult * availCommonPool
		rankedIslands := rankIslands(p, resourceRequest)
		for _, islandID := range rankedIslands {
			if resourcesToGive-resourceRequest[islandID] > 0 {
				resourceAllocation[islandID] = shared.Resources(resourceRequest[islandID])
				resourcesToGive -= resourceRequest[islandID]
				continue
			}
			if resourcesToGive > 0 {
				resourceAllocation[islandID] = shared.Resources(resourcesToGive)
				resourcesToGive = 0
				break
			}
			break
		}
	}

	return shared.PresidentReturnContent{
		ContentType: shared.PresidentAllocation,
		ResourceMap: resourceAllocation,
		ActionTaken: true,
	}
}

func rankIslands(p *President, resourceRequest map[shared.ClientID]shared.Resources) []shared.ClientID {
	// If islands have requested something then include them in ranking
	islandsToRank := make([]shared.ClientID, 0)
	for islandID, _ := range resourceRequest {
		islandsToRank = append(islandsToRank, islandID)
	}

	// Rank these islands according to how many turns they have been critical for first
	// TODO: get rounds crit for before being dead
	criticalTurnsToIslands := make(map[uint][]shared.ClientID)
	for turnsCriticalFor := p.c.gameConfig().MaxCriticalConsecutiveTurns; turnsCriticalFor >= 0; turnsCriticalFor-- {
		for _, islandID := range islandsToRank {
			p.c.gameState().ClientInfo.CriticalConsecutiveTurnsCounter
		}
	}

	// Secondary ranking on how many resources they have
}

func (p *President) PickRuleToVote(rulesProposals []rules.RuleMatrix) shared.PresidentReturnContent {
}

func (p *President) SetTaxationAmount(islandsResources map[shared.ClientID]shared.ResourcesReport) shared.PresidentReturnContent {
	// The islands are taxed based on how "Critical" they are --> more critical == less tax
	// How close a disaster is from occurring  - sometimes known, else use Hamish's function to get prediction from combinedPrediction (in disaster branch)
	// How well off the CP is
	// How much the roles' budget is
}

//don't think we need to overwrite
func (p *President) PaySpeaker() shared.PresidentReturnContent {
}

func (p *President) CallSpeakerElection(monitoring shared.MonitorResult, turnsInPower int, allIslands []shared.ClientID) shared.ElectionSettings {
	// example implementation calls an election if monitoring was performed and the result was negative
	var electionsettings = shared.ElectionSettings{
		VotingMethod:  shared.InstantRunoff,
		IslandsToVote: allIslands,
		HoldElection:  false,
	}
	if monitoring.Performed && !monitoring.Result {
		electionsettings.HoldElection = true
	}
	if turnsInPower >= 2 { //keep default because speaker is not too relevant
		electionsettings.HoldElection = true
	}
	return electionsettings
}
