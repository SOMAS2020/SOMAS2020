package baseclient

import (
	"math/rand"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type BasePresident struct{}

// EvaluateAllocationRequests sets allowed resource allocation based on each islands requests
// COMPULSORY: decide how much each island can take from the common pool.
func (p *BasePresident) EvaluateAllocationRequests(resourceRequest map[shared.ClientID]shared.Resources, availCommonPool shared.Resources) (map[shared.ClientID]shared.Resources, bool) {
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
	return resourceAllocation, true
}

// PickRuleToVote chooses a rule proposal from all the proposals
// COMPULSORY: decide a method for picking a rule to vote on
func (p *BasePresident) PickRuleToVote(rulesProposals []string) (string, bool) {
	if len(rulesProposals) == 0 {
		// No rules were proposed by the islands
		return "", false
	}
	return rulesProposals[rand.Intn(len(rulesProposals))], true
}

// SetTaxationAmount sets taxation amount for all of the living islands
// islandsResources: map of all the living islands and their remaining resources
// COMPULSORY: decide how much to tax islands, using rules as a guide if you wish.
func (p *BasePresident) SetTaxationAmount(islandsResources map[shared.ClientID]shared.Resources) (map[shared.ClientID]shared.Resources, bool) {
	taxAmountMap := make(map[shared.ClientID]shared.Resources)
	for id, resourceLeft := range islandsResources {
		taxAmountMap[id] = shared.Resources(float64(resourceLeft) * rand.Float64())
	}
	return taxAmountMap, true
}

// PaySpeaker pays the speaker a salary.
// OPTIONAL: override to pay the Speaker less than the full amount.
func (p *BasePresident) PaySpeaker(salary shared.Resources) (shared.Resources, bool) {
	// TODO : Implement opinion based salary payment.
	return salary, true
}

// CallSpeakerElection is called by the executive to decide on power-transfer
// COMPULSORY: decide when to call an election following relevant rulesInPlay if you wish
func (p *BasePresident) CallSpeakerElection(monitoring shared.MonitorResult, turnsInPower int, allIslands []shared.ClientID) shared.ElectionSettings {
	// example implementation calls an election if monitoring was performed and the result was negative
	// or if the number of turnsInPower exceeds 3
	var electionsettings = shared.ElectionSettings{
		VotingMethod:  shared.Plurality,
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
