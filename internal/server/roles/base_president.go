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

	if requestSum < 3*availCommonPool/4 {
		resourceAllocation = resourceRequest
	} else {
		for id, request := range resourceRequest {
			resourceAllocation[id] = int(request * availCommonPool * 3 / (4 * requestSum))
		}
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

func (p *basePresident) requestRuleProposal() {
	var rules []string
	//TODO: request Island for the rules (this function havnt been implemented
	//just create some mock function we need so they can write it up for us accordingly)
	// use this mock
	// mockfunction -> getIslandRuleProposed(islandID)string

	p.setRuleProposals(rules)
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

func (p *basePresident) broadcastTaxation(islandsResources map[int]int) {
	taxAmountMap := p.getTaxMap(islandsResources)
	//TODO: broadcastTaxation to every island
}

// Send Tax map all the remaining islands
// Called by orchestration at the end of the turn
func (p *basePresident) getAllocationRequests(commonPool int) map[int]int {
	if p.clientPresident != nil {
		result, error := p.clientPresident.evaluateAllocationRequests(p.resourceRequests, commonPool)
		if error == nil {
			return result
		}
	}
	result, _ := p.evaluateAllocationRequests(p.resourceRequests, commonPool)
	return result
}

func (p *basePresident) requestAllocationRequest() {
	allocRequests := make(map[int]int)
	//TODO: use mock function that to get the requests so neelesh can deal
	// mockfunction -> getIslandRequest(islandID)int
	//
	p.setAllocationRequest(allocRequests)

}

func (p *basePresident) replyAllocationRequest(commonPool int) {
	allocation := p.getAllocationRequests(commonPool)
	//TODO: broadcast the result
}

func (p *basePresident) appointNextSpeaker() int {
	return rand.Intn(5)
}

func (p *basePresident) withdrawSpeakerSalary(int) {
	//TODO: need to discuss with neelesh on how to be integrated
}

//TODO (optional): you can write a helper function (either put it here or orchestration.go) that you can use to broadcast or reply island
