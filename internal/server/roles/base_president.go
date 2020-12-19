package roles

import (
	"math/rand"
)

//base President Object
type basePresident struct {
	id                 int
	budget             int
	speakerSalary      int
	rulesProposals     []string
	resourceRequests   map[int]int
	resourceAllocation map[int]int
	ruleToVote         string
	taxAmountMap       map[int]int
}

// Set allowed resource allocation based on each islands requests
func (p *basePresident) evaluateAllocationRequests() {
	resourceAllocation := make(map[int]int)
	for id, request := range p.resourceRequests {
		resourceAllocation[id] = rand.Intn(request)
	}
	p.resourceAllocation = resourceAllocation
}

// Chose a rule proposal from all the proposals
func (p *basePresident) pickRuleToVote() {
	if len(p.rulesProposals) == 0 {
		// No rules were proposed by the islands
		p.ruleToVote = ""
	} else {
		p.ruleToVote = p.rulesProposals[rand.Intn(len(p.rulesProposals))]
	}
}

// Get rule proposals to be voted on from remaining islands
// Called by orchestration
func (p *basePresident) SetRuleProposals(rulesProposals []string) {
	p.rulesProposals = rulesProposals
}

// Set approved resources request for all the remaining islands
// Called by orchestration
func (p *basePresident) SetAllocationRequest(resourceRequests map[int]int) {
	p.resourceRequests = resourceRequests
}

// Set taxation amount for all of the living islands
// island_resources: map of all the living islands and their remaing resources
func (p *basePresident) SetTaxationAmount(islands_resources map[int]int) {
	taxAmountMap := make(map[int]int)
	for id, resource_left := range islands_resources {
		taxAmountMap[id] = rand.Intn(resource_left)
	}
	p.taxAmountMap = taxAmountMap
}

// Get rules to be voted on to Speaker
// Called by orchestration at the end of the turn
func (p *basePresident) GetRuleForSpeaker() string {
	p.pickRuleToVote()
	return p.ruleToVote
}

// Send Tax map all the remaining islands
// Called by orchestration at the end of the turn
func (p *basePresident) GetTaxMap() map[int]int {
	return p.taxAmountMap
}

// Send Tax map all the remaining islands
// Called by orchestration at the end of the turn
func (p *basePresident) GetAllocationRequests() map[int]int {
	p.evaluateAllocationRequests()
	return p.resourceAllocation
}

func (p *basePresident) appointNextSpeaker() int {
	return rand.Intn(5)
}
