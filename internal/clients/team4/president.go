package team4

import (
	"math/rand"
	"sort"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/config"

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
	actionOverBudget := p.parent.getRoleBudget(shared.President)-
		p.parent.getIIGOConfig().ReplyAllocationRequestsActionCost < 0

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

	//Limit resources to be distributed from the CP in a linear manner
	disasterPrediction := p.parent.obs.iifoObs.finalDisasterPrediction
	resourceThreshold := shared.Resources(0)
	if p.parent.getSeason() != 1 && disasterPrediction.Confidence >= 0.5 {
		if disasterPrediction.TimeLeft != 0 {
			resourceThreshold = shared.Resources(
				float64(p.parent.getTurn()) * float64(p.parent.getMinimumThreshold()) /
					(float64(disasterPrediction.TimeLeft) + float64(p.parent.getTurn())))
		} else if disasterPrediction.TimeLeft == 0 {
			resourceThreshold = p.parent.getMinimumThreshold()
		}
	}
	remainingResources := p.parent.getCommonPool() - resourceThreshold

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
	ret := shared.Resources(0)
	if request.Value*5 >= 5*p.parent.getCostOfLiving() {
		ret = 5 * p.parent.getCostOfLiving()
	}
	//weigh by trust
	ret = shared.Resources(float64(ret) * 2 * p.parent.getTrust(request.Key))
	return ret
}

// PickRuleToVote chooses a rule proposal from all the proposals
func (p *president) PickRuleToVote(rulesProposals []rules.RuleMatrix) shared.PresidentReturnContent {
	_, budgetRuleInPay := p.GameState.RulesInfo.CurrentRulesInPlay["president_over_budget"]
	actionOverBudget := p.parent.getRoleBudget(shared.President)-
		p.parent.getIIGOConfig().GetRuleForSpeakerActionCost < 0

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
	//TODO: Pick rules close to ideals
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
	actionOverBudget := p.parent.getRoleBudget(shared.President)-
		p.parent.getIIGOConfig().BroadcastTaxationActionCost < 0

	if budgetRuleInPay && actionOverBudget {
		return shared.PresidentReturnContent{
			ContentType: shared.PresidentTaxation,
			ActionTaken: false,
		}
	}

	taxAmountMap := make(map[shared.ClientID]shared.Resources)

	for clientID, clientReport := range islandsResources {
		weAreRichThreshold := p.generateWeAreRichThreshold()
		if p.parent.getCommonPool() > weAreRichThreshold {
			//Excuse everyone if its not season one and there are a lot of resources
			taxAmountMap[clientID] = shared.Resources(0)
		} else if p.parent.ServerReadHandle.GetGameState().ClientLifeStatuses[clientID] == shared.Critical {
			//Excuse if the client is critical
			taxAmountMap[clientID] = shared.Resources(0)
		} else if clientReport.Reported {
			//Excuse if the reports have to be true and the island is poor
			//TODO: remove magik numbers
			_, reportTruthfulnessRuleInPay := p.parent.ServerReadHandle.GetGameState().RulesInfo.CurrentRulesInPlay["island_must_report_actual_private_resource"]
			if reportTruthfulnessRuleInPay && clientReport.ReportedAmount*3 <= p.parent.getCostOfLiving() {
				taxAmountMap[clientID] = shared.Resources(0)
			} else if clientReport.ReportedAmount*1.5 <= p.parent.getCostOfLiving() {
				taxAmountMap[clientID] = shared.Resources(0)
			}
			taxAmountMap[clientID] = p.generateTaxAmount(clientID, clientReport.ReportedAmount)
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

func (p *president) generateWeAreRichThreshold() shared.Resources {
	disasterPrediction := p.parent.obs.iifoObs.finalDisasterPrediction
	if p.parent.getSeason() == 1 {
		return p.parent.getMinimumThreshold()
	} else {
		return shared.Resources(disasterPrediction.Magnitude * disasterPrediction.Confidence)
	}
}

func (p *president) generateTaxAmount(clientID shared.ClientID, clientReport shared.Resources) shared.Resources {
	if p.parent.getSeason() == 1 {
		//Static tax for when there is little information
		return 0.2 * clientReport
	}
	disasterPrediction := p.parent.obs.iifoObs.finalDisasterPrediction
	var averageTax shared.Resources
	if disasterPrediction.TimeLeft == 0 {
		averageTax = shared.Resources((disasterPrediction.Magnitude - float64(p.parent.getCommonPool())) / p.numIslandsAlive())
	} else {
		averageTax = shared.Resources((disasterPrediction.Magnitude - float64(p.parent.getCommonPool())) /
			(p.numIslandsAlive()) * float64(disasterPrediction.TimeLeft))
	}
	return averageTax

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
	actionOverBudget := p.parent.getRoleBudget(shared.President)-p.parent.getIIGOConfig().AppointNextSpeakerActionCost < 0

	if budgetRuleInPay && actionOverBudget {
		return shared.ElectionSettings{
			VotingMethod:  shared.InstantRunoff,
			IslandsToVote: allIslands,
			HoldElection:  false,
		}
	}

	//Compliance with election rule
	_, electionRuleInPay := p.GameState.RulesInfo.CurrentRulesInPlay["roles_must_hold_election"]
	speakerOverTerm := p.parent.getTurnsInPower(shared.Speaker) >
		p.parent.getTermLength(shared.Speaker)
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
		//If we trust the judge
		if p.parent.getTrust(p.parent.ServerReadHandle.GetGameState().JudgeID) > 0.4 {
			return shared.ElectionSettings{
				VotingMethod:  shared.InstantRunoff,
				IslandsToVote: allIslands,
				HoldElection:  true,
			}
		}
	}

	//If we dont trust the Speaker
	if p.parent.getTrust(p.parent.ServerReadHandle.GetGameState().SpeakerID) < 0.4 {
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
	_, resultRuleInPay := p.parent.ServerReadHandle.GetGameState().RulesInfo.CurrentRulesInPlay["must_appoint_elected_island"]
	if !resultRuleInPay {
		if p.parent.getTrust(winner) < 0.3 {
			//Assign at random if we dont like the result
			return p.selectRandomAliveIsland(winner)
		}
		return winner
	}
	return winner
}

func (p *president) numIslandsAlive() float64 {
	ret := 0
	for _, alive := range p.parent.ServerReadHandle.GetGameState().ClientLifeStatuses {
		if alive == shared.Alive {
			ret++
		}
	}
	return float64(ret)
}

func (p *president) selectRandomAliveIsland(winner shared.ClientID) shared.ClientID {
	var islands []shared.ClientID
	for islandID, status := range p.parent.ServerReadHandle.GetGameState().ClientLifeStatuses {
		if status != shared.Dead {
			islands = append(islands, islandID)
		}
	}
	if len(islands) > 0 {
		return islands[rand.Intn(int(p.numIslandsAlive()))]
	}
	return winner
}
