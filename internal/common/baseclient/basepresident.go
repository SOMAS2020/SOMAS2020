package baseclient

import (
	"math/rand"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type BasePresident struct {
	FlatTaxRate float64
}

// EvaluateAllocationRequests sets allowed resource allocation based on each islands requests
func (p *BasePresident) EvaluateAllocationRequests(resourceRequest map[shared.ClientID]shared.Resources, availCommonPool shared.Resources) shared.PresidentReturnContent {
	var requestSum shared.Resources
	resourceAllocation := make(map[shared.ClientID]shared.Resources)

	for _, request := range resourceRequest {
		requestSum += request
	}

	if requestSum < 0.75*availCommonPool {
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
func (p *BasePresident) PickRuleToVote(rulesProposals []string) shared.PresidentReturnContent {
	// DefaulContentType: No rules were proposed by the islands
	proposedRule := ""
	actionTaken := false

	// if some rules were proposed
	if len(rulesProposals) != 0 {
		proposedRule = rulesProposals[rand.Intn(len(rulesProposals))]
		actionTaken = true
	}

	return shared.PresidentReturnContent{
		ContentType:  shared.PresidentRuleProposal,
		ProposedRule: proposedRule,
		ActionTaken:  actionTaken,
	}
}

// SetTaxationAmount sets taxation amount for all of the living islands
// islandsResources: map of all the living islands and their reported resources
func (p *BasePresident) SetTaxationAmount(islandsResources map[shared.ClientID]shared.ResourcesReport) shared.PresidentReturnContent {
	taxAmountMap := make(map[shared.ClientID]shared.Resources)

	for clientID, clientReport := range islandsResources {
		if clientReport.Reported {
			taxAmountMap[clientID] = shared.Resources(float64(clientReport.ReportedAmount) * rand.Float64())
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
func (p *BasePresident) PaySpeaker(salary shared.Resources) shared.PresidentReturnContent {
	// TODO : Implement opinion based salary payment.
	return shared.PresidentReturnContent{
		ContentType:   shared.PresidentSpeakerSalary,
		SpeakerSalary: 0,
		ActionTaken:   true,
	}
}
