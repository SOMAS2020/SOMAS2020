package baseclient

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"math/rand"

	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type BasePresident struct {
	GameState gamestate.ClientGameState
}

// EvaluateAllocationRequests sets allowed resource allocation based on each islands requests
func (p *BasePresident) EvaluateAllocationRequests(resourceRequest map[shared.ClientID]shared.Resources, availCommonPool shared.Resources) shared.PresidentReturnContent {
	var requestSum shared.Resources
	resourceAllocation := make(map[shared.ClientID]shared.Resources)

	for _, request := range resourceRequest {
		requestSum += request
	}

	if requestSum < 0.75*availCommonPool || requestSum == 0 {
		resourceAllocation = resourceRequest
	} else {
		for id, request := range resourceRequest {
			resourceAllocation[id] = shared.Resources(request * availCommonPool * 3 / (4 * requestSum))
		}
	}

	return shared.PresidentReturnContent{
		ContentType: shared.PresidentAllocation,
		ResourceMap: resourceAllocation,
		ActionTaken: true,
	}
}

// PickRuleToVote chooses a rule proposal from all the proposals
func (p *BasePresident) PickRuleToVote(rulesProposals []rules.RuleMatrix) shared.PresidentReturnContent {
	// DefaulContentType: No rules were proposed by the islands
	proposedRuleMatrix := rules.RuleMatrix{}
	actionTaken := false

	// if some rules were proposed
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
func (p *BasePresident) SetTaxationAmount(islandsResources map[shared.ClientID]shared.ResourcesReport) shared.PresidentReturnContent {
	taxAmountMap := make(map[shared.ClientID]shared.Resources)

	for clientID, clientReport := range islandsResources {
		if clientReport.Reported {
			taxRate := 0.1
			taxAmountMap[clientID] = shared.Resources(float64(clientReport.ReportedAmount) * taxRate)
		} else {
			taxAmountMap[clientID] = 15 //flat tax rate
		}
	}
	return shared.PresidentReturnContent{
		ContentType: shared.PresidentTaxation,
		ResourceMap: taxAmountMap,
		ActionTaken: true,
	}
}

// PaySpeaker pays the speaker a salary.
func (p *BasePresident) PaySpeaker() shared.PresidentReturnContent {
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
// COMPULSORY: decide when to call an election following relevant rulesInPlay if you wish
func (p *BasePresident) CallSpeakerElection(monitoring shared.MonitorResult, turnsInPower int, allIslands []shared.ClientID) shared.ElectionSettings {
	// example implementation calls an election if monitoring was performed and the result was negative
	// or if the number of turnsInPower exceeds 3
	var electionsettings = shared.ElectionSettings{
		VotingMethod:  shared.InstantRunoff,
		IslandsToVote: allIslands,
		HoldElection:  false,
	}
	if monitoring.Performed && !monitoring.Result {
		electionsettings.HoldElection = true
	}
	if turnsInPower >= 2 {
		electionsettings.HoldElection = true
	}
	return electionsettings
}

// DecideNextSpeaker returns the ID of chosen next Speaker
// OPTIONAL: override to manipulate the result of the election
func (p *BasePresident) DecideNextSpeaker(winner shared.ClientID) shared.ClientID {
	return winner
}
