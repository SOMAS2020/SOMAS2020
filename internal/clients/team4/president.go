package team4

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
	"math/rand"
	"sort"

	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type president struct {
	*baseclient.BasePresident
	parent *client
	config.ClientConfig
}

type IslandRequestPair struct {
	Key   shared.ClientID
	Value shared.Resources
}

func rankByAllocationSize(allocationRequests map[shared.ClientID]shared.Resources) IslandRequestPairList {
	pl := make(IslandRequestPairList, len(allocationRequests))
	i := 0
	for k, v := range allocationRequests {
		pl[i] = IslandRequestPair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}

// IslandRequestPairList implements sort.Interface for []IslandRequestPair
type IslandRequestPairList []IslandRequestPair

func (p IslandRequestPairList) Len() int           { return len(p) }
func (p IslandRequestPairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p IslandRequestPairList) Less(i, j int) bool { return p[i].Value < p[j].Value }

// EvaluateAllocationRequests sets allowed resource allocation based on each islands requests
func (p *president) EvaluateAllocationRequests(resourceRequest map[shared.ClientID]shared.Resources, availCommonPool shared.Resources) shared.PresidentReturnContent {
	_, budgetRuleInPay := p.GameState.RulesInfo.CurrentRulesInPlay["president_over_budget"]
	actionOverBudget := p.parent.ServerReadHandle.GetGameState().IIGORolesBudget[shared.President]-
		p.parent.ServerReadHandle.GetGameConfig().IIGOClientConfig.ReplyAllocationRequestsActionCost < 0

	if budgetRuleInPay && actionOverBudget {
		return shared.PresidentReturnContent{
			ContentType: shared.PresidentAllocation,
			ActionTaken: false,
		}
	}

	//Separate and sort allocations
	criticalRequests := make(map[shared.ClientID]shared.Resources)
	for islandID, request := range resourceRequest {
		if p.parent.ServerReadHandle.GetGameState().ClientLifeStatuses[islandID] == shared.Critical {
			criticalRequests[islandID] = request
		}
	}
	sortedCriticalIslands := rankByAllocationSize(criticalRequests)

	nonCriticalRequests := make(map[shared.ClientID]shared.Resources)
	for islandID, request := range resourceRequest {
		if p.parent.ServerReadHandle.GetGameState().ClientLifeStatuses[islandID] == shared.Alive {
			nonCriticalRequests[islandID] = request
		}
	}
	sortedNonCriticalRequests := rankByAllocationSize(nonCriticalRequests)

	finalAllocation := make(map[shared.ClientID]shared.Resources)
	//TODO:take into account remaining vs predicted need of resources
	remainingResources := p.parent.ServerReadHandle.GetGameState().CommonPool

	//Allocations for critical islands
	for i := 0; i < len(sortedCriticalIslands); i++ {
		limitedAllocation := p.allocationLimiter(sortedCriticalIslands[i])
		if remainingResources-limitedAllocation >= 0 {
			finalAllocation[sortedCriticalIslands[i].Key] = limitedAllocation
			remainingResources -= limitedAllocation
		} else {
			finalAllocation[sortedCriticalIslands[i].Key] = shared.Resources(0)
		}
	}

	//Allocations for non-critical islands
	for i := 0; i < len(sortedNonCriticalRequests); i++ {
		limitedAllocation := p.allocationLimiter(sortedNonCriticalRequests[i])
		if remainingResources-limitedAllocation >= 0 {
			finalAllocation[sortedNonCriticalRequests[i].Key] = limitedAllocation
			remainingResources -= limitedAllocation
		} else {
			finalAllocation[sortedNonCriticalRequests[i].Key] = shared.Resources(0)
		}
	}

	return shared.PresidentReturnContent{
		ContentType: shared.PresidentAllocation,
		ResourceMap: finalAllocation,
		ActionTaken: true,
	}
}

func (p *president) allocationLimiter(request IslandRequestPair) shared.Resources {
	//TODO: Take into account an island's trust
	if request.Value*5 >= 5*p.parent.ServerReadHandle.GetGameConfig().CostOfLiving {
		return 5 * p.parent.ServerReadHandle.GetGameConfig().CostOfLiving
	}
	return request.Value
}

// PickRuleToVote chooses a rule proposal from all the proposals
func (p *president) PickRuleToVote(rulesProposals []rules.RuleMatrix) shared.PresidentReturnContent {
	_, budgetRuleInPay := p.GameState.RulesInfo.CurrentRulesInPlay["president_over_budget"]
	actionOverBudget := p.parent.ServerReadHandle.GetGameState().IIGORolesBudget[shared.President]-
		p.parent.ServerReadHandle.GetGameConfig().IIGOClientConfig.GetRuleForSpeakerActionCost < 0

	if budgetRuleInPay && actionOverBudget {
		return shared.PresidentReturnContent{
			ContentType: shared.PresidentRuleProposal,
			ActionTaken: false,
		}
	}

	// DefaultContentType: No rules were proposed by the islands
	proposedRuleMatrix := rules.RuleMatrix{}
	actionTaken := false

	// if some rules were proposed
	//TODO: Weigh chance of picking rules by trust in island
	if len(rulesProposals) != 0 {
		proposedRuleMatrix = rulesProposals[rand.Intn(len(rulesProposals))]
		actionTaken = true
	}

	return shared.PresidentReturnContent{
		ContentType:        shared.PresidentRuleProposal,
		ProposedRuleMatrix: proposedRuleMatrix,
		ActionTaken:        actionTaken,
	}
}

// SetTaxationAmount sets taxation amount for all of the living islands
// islandsResources: map of all the living islands and their reported resources
func (p *president) SetTaxationAmount(islandsResources map[shared.ClientID]shared.ResourcesReport) shared.PresidentReturnContent {
	_, budgetRuleInPay := p.GameState.RulesInfo.CurrentRulesInPlay["president_over_budget"]
	actionOverBudget := p.parent.ServerReadHandle.GetGameState().IIGORolesBudget[shared.President]-
		p.parent.ServerReadHandle.GetGameConfig().IIGOClientConfig.BroadcastTaxationActionCost < 0

	if budgetRuleInPay && actionOverBudget {
		return shared.PresidentReturnContent{
			ContentType: shared.PresidentTaxation,
			ActionTaken: false,
		}
	}

	taxAmountMap := make(map[shared.ClientID]shared.Resources)

	for clientID, clientReport := range islandsResources {
		//TODO:excuse everyone if we are 1.5 times over predicted disaster strength

		//Excuse if the client is critical
		if p.parent.ServerReadHandle.GetGameState().ClientLifeStatuses[clientID] == shared.Critical {
			taxAmountMap[clientID] = shared.Resources(0)
		} else if clientReport.Reported {
			//Excuse if the reports have to be true and the island is poor
			//TODO: remove magik numbers
			_, reportTruthfulnessRuleInPay := p.GameState.RulesInfo.CurrentRulesInPlay["island_must_report_actual_private_resource"]
			if reportTruthfulnessRuleInPay && clientReport.ReportedAmount*3 <= p.parent.ServerReadHandle.GetGameConfig().CostOfLiving {
				taxAmountMap[clientID] = shared.Resources(0)
			} else if clientReport.ReportedAmount*1.5 <= p.parent.ServerReadHandle.GetGameConfig().CostOfLiving {
				taxAmountMap[clientID] = shared.Resources(0)
			}
			//derived from
			//affected by how close we are to the predicted threshold * 1
			//how much we trust the island * 0.3
			//how much they have (exponentially)
			taxRate := 0.1
			taxAmountMap[clientID] = shared.Resources(float64(clientReport.ReportedAmount) * taxRate)
		} else {
			taxAmountMap[clientID] = 15 //flat tax rate
		}

		//Excuse ourselves if we are not fair
		if clientID == shared.Team4 {
			taxAmountMap[clientID] = shared.Resources(float64(taxAmountMap[clientID]) * p.parent.internalParam.fairness)
		}
	}

	return shared.PresidentReturnContent{
		ContentType: shared.PresidentTaxation,
		ResourceMap: taxAmountMap,
		ActionTaken: true,
	}
}

// PaySpeaker pays the speaker a salary.
func (p *president) PaySpeaker() shared.PresidentReturnContent {
	SpeakerSalaryRule, ok := p.GameState.RulesInfo.CurrentRulesInPlay["salary_cycle_speaker"]
	var SpeakerSalary shared.Resources = 0
	if ok {
		SpeakerSalary = shared.Resources(SpeakerSalaryRule.ApplicableMatrix.At(0, 1))
	}
	return shared.PresidentReturnContent{
		ContentType:   shared.PresidentSpeakerSalary,
		SpeakerSalary: SpeakerSalary,
		ActionTaken:   true,
	}
}

// CallSpeakerElection is called by the executive to decide on power-transfer
func (p *president) CallSpeakerElection(monitoring shared.MonitorResult, turnsInPower int, allIslands []shared.ClientID) shared.ElectionSettings {
	//Compliance with budget rule
	_, budgetRuleInPay := p.GameState.RulesInfo.CurrentRulesInPlay["president_over_budget"]
	actionOverBudget := p.parent.ServerReadHandle.GetGameState().IIGORolesBudget[shared.President]-
		p.parent.ServerReadHandle.GetGameConfig().IIGOClientConfig.AppointNextSpeakerActionCost < 0

	if budgetRuleInPay && actionOverBudget {
		return shared.ElectionSettings{
			VotingMethod:  shared.InstantRunoff,
			IslandsToVote: allIslands,
			HoldElection:  false,
		}
	}

	//Compliance with election rule
	_, electionRuleInPay := p.GameState.RulesInfo.CurrentRulesInPlay["roles_must_hold_election"]
	speakerOverTerm := p.parent.ServerReadHandle.GetGameState().IIGOTurnsInPower[shared.Speaker] >
		p.parent.ServerReadHandle.GetGameConfig().IIGOClientConfig.IIGOTermLengths[shared.Speaker]
	if electionRuleInPay {
		if speakerOverTerm {
			return shared.ElectionSettings{
				VotingMethod:  shared.InstantRunoff,
				IslandsToVote: allIslands,
				HoldElection:  true,
			}
		}
		return shared.ElectionSettings{
			VotingMethod:  shared.InstantRunoff,
			IslandsToVote: allIslands,
			HoldElection:  false,
		}
	}

	//If we get this far there are no rules that say whether we should hold an election
	//Hold an election if they broke some rules
	if monitoring.Result {
		//TODO: test trust of island announcing monitoring
		return shared.ElectionSettings{
			VotingMethod:  shared.InstantRunoff,
			IslandsToVote: allIslands,
			HoldElection:  true,
		}
	}

	//TODO: use a trust metric instead
	if rand.Float64() <= 0.25 {
		return shared.ElectionSettings{
			VotingMethod:  shared.InstantRunoff,
			IslandsToVote: allIslands,
			HoldElection:  true,
		}
	}
	return shared.ElectionSettings{
		VotingMethod:  shared.InstantRunoff,
		IslandsToVote: allIslands,
		HoldElection:  false,
	}

}

// DecideNextSpeaker returns the ID of chosen next Speaker
func (p *president) DecideNextSpeaker(winner shared.ClientID) shared.ClientID {
	//No manipulate. Our agent trust democracy
	_, resultRuleInPay := p.GameState.RulesInfo.CurrentRulesInPlay["must_appoint_elected_island"]
	if !resultRuleInPay {
		//TODO: change to the most trusted island
		return winner
	}
	return winner
}
