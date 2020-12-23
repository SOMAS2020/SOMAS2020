package iigointernal

import (
	"math/rand"

	"github.com/SOMAS2020/SOMAS2020/internal/common/voting"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

//base President Object
type basePresident struct {
	Id               int
	clientPresident  roles.President
	budget           int
	speakerSalary    int
	RulesProposals   []string
	ResourceRequests map[int]int
	//resourceAllocation map[int]int
	//RuleToVote         string
	//taxAmountMap       map[int]int
}

// Set allowed resource allocation based on each islands requests
func (p *basePresident) EvaluateAllocationRequests(resourceRequest map[int]int, availCommonPool int) (map[int]int, error) {
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
func (p *basePresident) PickRuleToVote(rulesProposals []string) (string, error) {
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
	p.RulesProposals = rulesProposals
}

// setAllocationRequests a map of the resourceRequests, uses that array to set
// the class p.resourceRequests class attribute to it.
func (p *basePresident) setAllocationRequest(resourceRequests map[int]int) {
	p.ResourceRequests = resourceRequests
}

// setTaxationAmount sets the tax rate for all of the living islands. It requires
// the remaining resources as an input to ensure the tax required is not higher than what
// islands can afford. It returns a map indexed with the island IDs, containing the amount of
// resources required as tax in the form of an integer.
func (p *basePresident) SetTaxationAmount(islandsResources map[int]int) (map[int]int, error) {
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
		result, error := p.clientPresident.PickRuleToVote(p.RulesProposals)
		if error == nil {
			return result
		}
	}
	result, _ := p.PickRuleToVote(p.RulesProposals)
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
	result, _ := p.SetTaxationAmount(islandsResources)
	return result
}

// broadcastTaxation broadcasts the tax amount decided by the president to all island still in the game
func (p *basePresident) broadcastTaxation(islandsResources map[int]int) {
	p.budget -= 10
	taxAmountMap := p.getTaxMap(islandsResources)
	for _, v := range getIslandAlive() {
		d := baseclient.Communication{T: baseclient.CommunicationInt, IntegerData: taxAmountMap[int(v)]}
		data := make(map[int]baseclient.Communication)
		data[TaxAmount] = d
		communicateWithIslands(shared.TeamIDs[int(v)], shared.TeamIDs[p.Id], data)
	}
}

// getAllocationRequests returns the allowed allocation requests from the common
// pool calculated by evaluateAllocationRequests. It returns the allocation requests
// as a map of ints indexed by islands ids.
func (p *basePresident) getAllocationRequests(commonPool int) map[int]int {
	if p.clientPresident != nil {
		result, error := p.clientPresident.EvaluateAllocationRequests(p.ResourceRequests, commonPool)
		if error == nil {
			return result
		}
	}
	result, _ := p.EvaluateAllocationRequests(p.ResourceRequests, commonPool)
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
		d := baseclient.Communication{T: baseclient.CommunicationInt, IntegerData: allocationMap[int(v)]}
		data := make(map[int]baseclient.Communication)
		data[AllocationAmount] = d
		communicateWithIslands(shared.TeamIDs[int(v)], shared.TeamIDs[p.Id], data)
	}
}

// appointNextSpeaker returns the island id of the island appointed to be speaker in the next turn
func (p *basePresident) appointNextSpeaker(clientIDs []shared.ClientID) int {
	p.budget -= 10
	var election voting.Election
	election.ProposeElection(baseclient.Speaker, voting.Plurality)
	election.OpenBallot(clientIDs)
	election.Vote(iigoClients)
	return int(election.CloseBallot())
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
	amount, _ := p.PaySpeaker()
	featureSpeaker.budget = amount
}

// paySpeaker pays speaker salary
func (p *basePresident) PaySpeaker() (int, error) {
	hold := p.speakerSalary
	p.speakerSalary = 0
	return hold, nil
}

func getIslandAlive() []float64 {
	return rules.VariableMap["islands_alive"].Values
}

func (p *basePresident) Reset(val string) error {
	p.Id = 0
	p.clientPresident = nil
	p.budget = 0
	p.ResourceRequests = map[int]int{}
	p.RulesProposals = []string{}
	p.speakerSalary = 0
	return nil
}
