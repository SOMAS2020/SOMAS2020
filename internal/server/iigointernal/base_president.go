package iigointernal

import (
	"math/rand"

	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

//base President Object
type basePresident struct {
	id               int
	clientPresident  roles.President
	budget           int
	speakerSalary    int
	rulesProposals   []string
	resourceRequests map[int]int
}

// evaluateAllocationRequests takes in resource requests from all islands and
// the available common pool. Returns map of allocated resources in which either
// meets demand if <75% of common pool is used, else scales allocations such that
// the 75% threshold usage is met.
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

// pickRuleToVote choses a rule proposal from all the proposals made by the islands.
// It returns an empty string if no rules were proposed, and the rule name (string)
// of the selected rule otherwise.
func (p *basePresident) pickRuleToVote(rulesProposals []string) (string, error) {
	p.budget -= 10
	if len(rulesProposals) == 0 {
		// No rules were proposed by the islands
		return "", nil
	}
	return rulesProposals[rand.Intn(len(rulesProposals))], nil
}

//requestRuleProposal asks each island alive for its rule proposal
func (p *basePresident) requestRuleProposal() {
	p.budget -= 10
	var rules []string
	for _, v := range getIslandAlive() {
		rules = append(rules, iigoClients[shared.ClientID(int(v))].RuleProposal())
	}
	p.setRuleProposals(rules)
}

// setRuleProposals takes in a array of rules, uses that array to set
// the class p.rulesProposal class attribute to it.
func (p *basePresident) setRuleProposals(rulesProposals []string) {
	p.rulesProposals = rulesProposals
}

// setAllocationRequests a map of the resourceRequests, uses that array to set
// the class p.resourceRequests class attribute to it.
func (p *basePresident) setAllocationRequest(resourceRequests map[int]int) {
	p.resourceRequests = resourceRequests
}

// setTaxationAmount sets the tax rate for all of the living islands. It requires
// the remaining resources as an input to ensure the tax required is not higher than what
// islands can afford. It returns a map indexed with the island IDs, containing the amount of
// resources required as tax in the form of an integer.
func (p *basePresident) setTaxationAmount(islandsResources map[int]int) (map[int]int, error) {
	p.budget -= 10
	taxAmountMap := make(map[int]int)
	for id, resourceLeft := range islandsResources {
		taxAmountMap[id] = rand.Intn(resourceLeft)
	}
	TaxAmountMapExport = taxAmountMap
	return taxAmountMap, nil
}

// getRuleForSpeaker returns the rule to be voted on to Speaker.
// It returns the name of the rule to be voted on as a string.
func (p *basePresident) getRuleForSpeaker() string {
	if p.clientPresident != nil {
		result, error := p.clientPresident.PickRuleToVote(p.rulesProposals)
		if error == nil {
			return result
		}
	}
	result, _ := p.pickRuleToVote(p.rulesProposals)
	return result
}

// getTaxMap returns the taxation map calculated by setTaxationAmount as a map of integers.
func (p *basePresident) getTaxMap(islandsResources map[int]int) map[int]int {
	p.budget -= 10
	if p.clientPresident != nil {
		result, error := p.clientPresident.SetTaxationAmount(islandsResources)
		if error == nil {
			return result
		}
	}
	result, _ := p.setTaxationAmount(islandsResources)
	return result
}

// broadcastTaxation broadcasts the tax amount decided by the president to all island still in the game
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

// getAllocationRequests returns the allowed allocation requests from the common
// pool calculated by evaluateAllocationRequests. It returns the allocation requests
// as a map of ints indexed by islands ids.
func (p *basePresident) getAllocationRequests(commonPool int) map[int]int {
	if p.clientPresident != nil {
		result, error := p.clientPresident.EvaluateAllocationRequests(p.resourceRequests, commonPool)
		if error == nil {
			return result
		}
	}
	result, _ := p.evaluateAllocationRequests(p.resourceRequests, commonPool)
	return result
}

// requestAllocationRequest asks all alive islands for its resource allocation request
func (p *basePresident) requestAllocationRequest() {
	allocRequests := make(map[int]int)
	for _, v := range getIslandAlive() {
		allocRequests[int(v)] = iigoClients[shared.ClientID(int(v))].CommonPoolResourceRequest()
	}
	AllocationAmountMapExport = allocRequests
	p.setAllocationRequest(allocRequests)

}

// replyAllocationRequest broadcasts the allocation of resources decided by the president
// to all islands alive
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

// appointNextSpeaker returns the island id of the island appointed to be speaker in the next turn
func (p *basePresident) appointNextSpeaker() int {
	p.budget -= 10
	return rand.Intn(5)
}

// withdrawSpeakerSalary withdraws the salary for speaker from the common pool
func (p *basePresident) withdrawSpeakerSalary(gameState *gamestate.GameState) error {
	var speakerSalary = int(rules.VariableMap["speakerSalary"].Values[0])
	var withdrawError = WithdrawFromCommonPool(speakerSalary, gameState)
	if withdrawError != nil {
		featurePresident.speakerSalary = speakerSalary
	}
	return withdrawError
}

// sendSpeakerSalary send speaker's salary to the speaker
func (p *basePresident) sendSpeakerSalary() {
	if p.clientPresident != nil {
		amount, err := p.clientPresident.PaySpeaker()
		if err == nil {
			featureSpeaker.budget = amount
			return
		}
	}
	amount, _ := p.paySpeaker()
	featureSpeaker.budget = amount
}

// paySpeaker pays speaker salary
func (p *basePresident) paySpeaker() (int, error) {
	hold := p.speakerSalary
	p.speakerSalary = 0
	return hold, nil
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
