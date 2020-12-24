package iigointernal

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/voting"
	"math/rand"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

//base President Object
type basePresident struct {
	ID               shared.ClientID
	clientPresident  roles.President
	budget           shared.Resources
	speakerSalary    shared.Resources
	RulesProposals   []string
	ResourceRequests map[shared.ClientID]shared.Resources
	//resourceAllocation map[shared.ClientID]shared.Resources
	//RuleToVote         string
	//taxAmountMap       map[shared.ClientID]shared.Resources
}

// returnSpeakerSalary returns the salary to the common pool.
func (p *basePresident) returnSpeakerSalary() shared.Resources {
	x := p.speakerSalary
	p.speakerSalary = 0
	return x
}

// Set allowed resource allocation based on each islands requests
func (p *basePresident) EvaluateAllocationRequests(resourceRequest map[shared.ClientID]shared.Resources, availCommonPool shared.Resources) (map[shared.ClientID]shared.Resources, error) {
	p.budget -= 10
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
	return resourceAllocation, nil
}

// Chose a rule proposal from all the proposals
// need to pass in since this is now functional for the sake of client side
func (p *basePresident) PickRuleToVote(rulesProposals []string) (string, error) {
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
	for _, v := range getIslandAlive() {
		rules = append(rules, iigoClients[shared.ClientID(int(v))].RuleProposal())
	}
	p.setRuleProposals(rules)
}

// Get rule proposals to be voted on from remaining islands
// Called by orchestration
func (p *basePresident) setRuleProposals(rulesProposals []string) {
	p.RulesProposals = rulesProposals
}

// Set approved resources request for all the remaining islands
// Called by orchestration
func (p *basePresident) setAllocationRequest(resourceRequests map[shared.ClientID]shared.Resources) {
	p.ResourceRequests = resourceRequests
}

// Set taxation amount for all of the living islands
// island_resources: map of all the living islands and their remaining resources
func (p *basePresident) SetTaxationAmount(islandsResources map[shared.ClientID]shared.Resources) (map[shared.ClientID]shared.Resources, error) {
	p.budget -= 10
	taxAmountMap := make(map[shared.ClientID]shared.Resources)
	for id, resourceLeft := range islandsResources {
		taxAmountMap[id] = shared.Resources(float64(resourceLeft) * rand.Float64())
	}
	TaxAmountMapExport = taxAmountMap
	return taxAmountMap, nil
}

// Get rules to be voted on to Speaker
// Called by orchestration at the end of the turn
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

// Send Tax map all the remaining islands
// Called by orchestration at the end of the turn
func (p *basePresident) getTaxMap(islandsResources map[shared.ClientID]shared.Resources) map[shared.ClientID]shared.Resources {
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

func (p *basePresident) broadcastTaxation(islandsResources map[shared.ClientID]shared.Resources) {
	p.budget -= 10
	taxAmountMap := p.getTaxMap(islandsResources)
	for _, v := range getIslandAlive() {
		d := baseclient.Communication{T: baseclient.CommunicationInt, IntegerData: int(taxAmountMap[shared.ClientID(int(v))])}
		data := make(map[int]baseclient.Communication)
		data[TaxAmount] = d
		communicateWithIslands(shared.TeamIDs[int(v)], shared.TeamIDs[p.ID], data)
	}
}

// Send Tax map all the remaining islands
// Called by orchestration at the end of the turn
func (p *basePresident) getAllocationRequests(commonPool shared.Resources) map[shared.ClientID]shared.Resources {
	if p.clientPresident != nil {
		result, error := p.clientPresident.EvaluateAllocationRequests(p.ResourceRequests, commonPool)
		if error == nil {
			return result
		}
	}
	result, _ := p.EvaluateAllocationRequests(p.ResourceRequests, commonPool)
	return result
}

func (p *basePresident) requestAllocationRequest() {
	allocRequests := make(map[shared.ClientID]shared.Resources)
	for _, v := range getIslandAlive() {
		allocRequests[shared.ClientID(int(v))] = iigoClients[shared.ClientID(int(v))].CommonPoolResourceRequest()
	}
	AllocationAmountMapExport = allocRequests
	p.setAllocationRequest(allocRequests)

}

func (p *basePresident) replyAllocationRequest(commonPool shared.Resources) {
	p.budget -= 10
	allocationMap := p.getAllocationRequests(commonPool)
	for _, v := range getIslandAlive() {
		d := baseclient.Communication{T: baseclient.CommunicationInt, IntegerData: int(allocationMap[shared.ClientID(int(v))])}
		data := make(map[int]baseclient.Communication)
		data[AllocationAmount] = d
		communicateWithIslands(shared.TeamIDs[int(v)], shared.TeamIDs[p.ID], data)
	}
}

func (p *basePresident) appointNextSpeaker(clientIDs []shared.ClientID) shared.ClientID {
	p.budget -= 10
	var election voting.Election
	election.ProposeElection(baseclient.Speaker, voting.Plurality)
	election.OpenBallot(clientIDs)
	election.Vote(iigoClients)
	return election.CloseBallot()
}

func (p *basePresident) withdrawSpeakerSalary(gameState *gamestate.GameState) error {
	var speakerSalary = shared.Resources(rules.VariableMap["speakerSalary"].Values[0])
	var withdrawError = WithdrawFromCommonPool(speakerSalary, gameState)
	if withdrawError != nil {
		featurePresident.speakerSalary = speakerSalary
	}
	return withdrawError
}

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

// Pay the speaker
func (p *basePresident) PaySpeaker() (shared.Resources, error) {
	hold := p.speakerSalary
	p.speakerSalary = 0
	return hold, nil
}

func getIslandAlive() []float64 {
	return rules.VariableMap["islands_alive"].Values
}

func (p *basePresident) Reset(val string) error {
	p.ID = 0
	p.clientPresident = nil
	p.budget = 0
	p.ResourceRequests = map[shared.ClientID]shared.Resources{}
	p.RulesProposals = []string{}
	p.speakerSalary = 0
	return nil
}
