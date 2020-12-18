package roles

import (
	"math/rand"
)

//base President Objectlist
type BasePresident struct {
	id                 int
	budget             int
	speakerSalary      int
	rulesProposals     []int
	resourceRequests   map[int]int
	resourceAllocation map[int]int
	ruleToVote         int
	taxAmountMap       map[int]int
}

// Set taxation amount for all of the living islands
// island_resources: map of all the living islands and their remaing resources
func (p *BasePresident) setTaxationAmount(islands_resources map[int]int) {
	taxAmountMap := make(map[int]int)
	for id, resource_left := range islands_resources {
		taxAmountMap[id] = rand.Intn(resource_left)
	}
	p.taxAmountMap = taxAmountMap
}

// Set allowed resource allocation based on each islands requests
func (p *BasePresident) setAllocationRequest() {
	resourceAllocation := make(map[int]int)
	for id, request := range p.resourceRequests {
		resourceAllocation[id] = rand.Intn(request)
	}
	p.resourceAllocation = resourceAllocation
}

// Chose a rule proposal from all the proposals
func (p *BasePresident) choseRuleFromProposals() {
	if len(p.rulesProposals) == 0 {
		// No rules were proposed by the islands
		p.ruleToVote = -1
	} else {
		p.ruleToVote = rand.Intn(len(p.rulesProposals))
	}
}

// Send rule to be voted on to Speaker
// Called by orchestration at the end of the turn
func (p *BasePresident) GetRuleForSpeaker() int {
	return p.ruleToVote
}

// Send Tax map all the remaining islands
// Called by orchestration at the end of the turn
func (p *BasePresident) GetTaxMap() map[int]int {
	return p.taxAmountMap
}

// Send approved resources request for all the remaining islands
// Called by orchestration at the end of the turn
func (p *BasePresident) GetAllocationRequest() map[int]int {
	return p.resourceAllocation
}
