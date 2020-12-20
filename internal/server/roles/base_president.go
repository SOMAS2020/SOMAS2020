package roles

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"math/rand"

	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
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
	p.budget -= 10
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
	p.budget -= 10
	if len(rulesProposals) == 0 {
		// No rules were proposed by the islands
		return "", nil
	}
	return rulesProposals[rand.Intn(len(rulesProposals))], nil
}

func (p *basePresident) requestRuleProposal() {
	p.budget -= 10
	var rules []string
	// use this mock
	// mockfunction -> getIslandRuleProposed(islandID int)string
	// TODO (neelesh): get this working pls
	//for _, v := range getIslandAlive() {
	//	rules = append(rules, getIslandRuleProposed(int(v)))
	//}
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
	p.budget -= 10
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
	p.budget -= 10
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
	p.budget -= 10
	taxAmountMap := p.getTaxMap(islandsResources)
	for _, v := range getIslandAlive() {
		d := DataPacket{integerData: taxAmountMap[int(v)]}
		data := make(map[int]DataPacket)
		data[TaxAmount] = d
		communicateWithIslands(int(v), p.id, data)
	}
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
	//for _, v := range getIslandAlive() {
	//	allocRequests[int(v)] = getIslandRequest(int(v)
	//}
	p.setAllocationRequest(allocRequests)

}

func (p *basePresident) replyAllocationRequest(commonPool int) {
	p.budget -= 10
	allocationMap := p.getAllocationRequests(commonPool)
	for _, v := range getIslandAlive() {
		d := DataPacket{integerData: allocationMap[int(v)]}
		data := make(map[int]DataPacket)
		data[AllocationAmount] = d
		communicateWithIslands(int(v), p.id, data)
	}
}

func (p *basePresident) appointNextSpeaker() int {
	p.budget -= 10
	return rand.Intn(5)
}

func (p *basePresident) withdrawSpeakerSalary(gameState *gamestate.GameState) error {
	var speakerSalary = int(rules.VariableMap["speakerSalary"].Values[0])
	var withdrawError = WithdrawFromCommonPool(speakerSalary, gameState)
	if withdrawError != nil {
		Base_President.speakerSalary = speakerSalary
	}
	return withdrawError
}

// Pay the speaker
func (p *basePresident) paySpeaker() {
	Base_speaker.budget = Base_President.speakerSalary
	Base_President.speakerSalary = 0
}

func getIslandAlive() []float64 {
	return rules.VariableMap["islands_alive"].Values
}

func (p *basePresident) reset(val string) error {
	p.id = 0
	p.clientPresident = nil
	p.budget = 0
	p.resourceRequests = map[int]int{}
	p.rulesProposals = []string{}
	p.speakerSalary = 0
	return nil
}
