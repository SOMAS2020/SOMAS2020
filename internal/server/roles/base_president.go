package roles

import (
	"math/rand"
)

//base President Object
type basePresident struct {
	id               int
	clientPresident  President
	budget           int
	speakerSalary    int
	rulesProposals   []string
	resourceRequests map[int]int
	//resourceAllocation map[int]int
	//ruleToVote         string
	//taxAmountMap       map[int]int
}

// Set allowed resource allocation based on each islands requests
func (p *basePresident) evaluateAllocationRequests(resourceRequest map[int]int, availCommonPool int) (map[int]int, error) {
	var requestSum int
	resourceAllocation := make(map[int]int)

	for _, request := range resourceRequest {
		requestSum += request
	}

	if requestSum < 0.75*availCommonPool {
		resourceAllocation = resourceRequest
	} else {
		for id, request := range resourceRequest {
			resourceAllocation[id] = int((request / requestSum)* 0.75 * availCommonPool)
	}
	return resourceAllocation, nil
}

// Chose a rule proposal from all the proposals
// need to pass in since this is now functional for the sake of client side
func (p *basePresident) pickRuleToVote(rulesProposals []string) (string, error) {
	if len(rulesProposals) == 0 {
		// No rules were proposed by the islands
		return "", nil
	}
	return rulesProposals[rand.Intn(len(rulesProposals))], nil
}

// Get rule proposals to be voted on from remaining islands
// Called by orchestration
func (p *basePresident) setRuleProposals(rulesProposals []string) {
	p.rulesProposals = rulesProposals
}

// Set approved resources request for all the remaining islands
// Called by orchestration
func (p *basePresident) setAllocationRequest(resourceRequests map[int]int) {
	p.resourceRequests = resourceRequests
}

// Set taxation amount for all of the living islands
// island_resources: map of all the living islands and their remaing resources
func (p *basePresident) setTaxationAmount(islandsResources map[int]int) (map[int]int, error) {
	taxAmountMap := make(map[int]int)
	for id, resourceLeft := range islandsResources {
		taxAmountMap[id] = rand.Intn(resourceLeft)
	}
	return taxAmountMap, nil
}

// Get rules to be voted on to Speaker
// Called by orchestration at the end of the turn
func (p *basePresident) getRuleForSpeaker() string {
	if p.clientPresident != nil {
		result, error := p.clientPresident.pickRuleToVote(p.rulesProposals)
		if error == nil {
			return result
		}
	}
	result, _ := p.pickRuleToVote(p.rulesProposals)
	return result
}

// Send Tax map all the remaining islands
// Called by orchestration at the end of the turn
func (p *basePresident) getTaxMap(islandsResources map[int]int) map[int]int {
	if p.clientPresident != nil {
		result, error := p.clientPresident.setTaxationAmount(islandsResources)
		if error == nil {
			return result
		}
	}
	result, _ := p.setTaxationAmount(islandsResources)
	return result
}

// Send Tax map all the remaining islands
// Called by orchestration at the end of the turn
func (p *basePresident) getAllocationRequests() map[int]int {
	if p.clientPresident != nil {
		result, error := p.clientPresident.evaluateAllocationRequests(p.resourceRequests)
		if error == nil {
			return result
		}
	}
	result, _ := p.evaluateAllocationRequests(p.resourceRequests)
	return result
}

func (p *basePresident) appointNextSpeaker() int {
	return rand.Intn(5)
}
