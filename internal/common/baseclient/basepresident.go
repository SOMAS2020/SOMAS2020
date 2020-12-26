package baseclient

import (
	"math/rand"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type BasePresident struct{}

// EvaluateAllocationRequests sets allowed resource allocation based on each islands requests
func (p *BasePresident) EvaluateAllocationRequests(resourceRequest map[shared.ClientID]shared.Resources, availCommonPool shared.Resources) (map[shared.ClientID]shared.Resources, bool) {
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
	return resourceAllocation, true
}

// PickRuleToVote chooses a rule proposal from all the proposals
// need to pass in since this is now functional for the sake of client side
func (p *BasePresident) PickRuleToVote(rulesProposals []string) (string, bool) {
	if len(rulesProposals) == 0 {
		// No rules were proposed by the islands
		return "", false
	}
	return rulesProposals[rand.Intn(len(rulesProposals))], true
}

// SetTaxationAmount sets taxation amount for all of the living islands
// islandsResources: map of all the living islands and their remaining resources
func (p *BasePresident) SetTaxationAmount(islandsResources map[shared.ClientID]shared.Resources) (map[shared.ClientID]shared.Resources, bool) {
	taxAmountMap := make(map[shared.ClientID]shared.Resources)
	for id, resourceLeft := range islandsResources {
		taxAmountMap[id] = shared.Resources(float64(resourceLeft) * rand.Float64())
	}
	return taxAmountMap, true
}

// PaySpeaker pays the speaker a salary.
func (p *BasePresident) PaySpeaker(salary shared.Resources) (shared.Resources, bool) {
	// TODO : Implement opinion based salary payment.
	return salary, true
}
