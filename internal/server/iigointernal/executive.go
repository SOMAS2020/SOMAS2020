package iigointernal

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"github.com/SOMAS2020/SOMAS2020/internal/common/voting"
)

type executive struct {
	ID               shared.ClientID
	clientPresident  roles.President
	budget           shared.Resources
	speakerSalary    shared.Resources
	RulesProposals   []string
	ResourceRequests map[shared.ClientID]shared.Resources
}

// returnSpeakerSalary returns the salary to the common pool.
func (e *executive) returnSpeakerSalary() shared.Resources {
	x := e.speakerSalary
	e.speakerSalary = 0
	return x
}

// Get rule proposals to be voted on from remaining islands
// Called by orchestration
func (e *executive) setRuleProposals(rulesProposals []string) {
	e.RulesProposals = rulesProposals
}

// Set approved resources request for all the remaining islands
// Called by orchestration
func (e *executive) setAllocationRequest(resourceRequests map[shared.ClientID]shared.Resources) {
	e.ResourceRequests = resourceRequests
}

// Get rules to be voted on to Speaker
// Called by orchestration at the end of the turn
func (e *executive) getRuleForSpeaker() (string, bool) {
	e.budget -= serviceCharge
	return e.clientPresident.PickRuleToVote(e.RulesProposals)
}

// Send Tax map all the remaining islands
// Called by orchestration at the end of the turn
func (e *executive) getTaxMap(islandsResources map[shared.ClientID]shared.Resources) (map[shared.ClientID]shared.Resources, bool) {
	e.budget -= serviceCharge
	return e.clientPresident.SetTaxationAmount(islandsResources)
}

// broadcastTaxation broadcasts the tax amount decided by the president to all island still in the game.
func (e *executive) broadcastTaxation(islandsResources map[shared.ClientID]shared.Resources) {
	e.budget -= serviceCharge
	taxAmountMap, taxesCollected := e.getTaxMap(islandsResources)
	if taxesCollected {
		for _, v := range getIslandAlive() {
			d := shared.Communication{T: shared.CommunicationInt, IntegerData: int(taxAmountMap[shared.ClientID(int(v))])}
			data := make(map[shared.CommunicationFieldName]shared.Communication)
			data[shared.TaxAmount] = d
			communicateWithIslands(shared.TeamIDs[int(v)], shared.TeamIDs[e.ID], data)
		}
	}
}

// Send Tax map all the remaining islands
// Called by orchestration at the end of the turn
func (e *executive) getAllocationRequests(commonPool shared.Resources) (map[shared.ClientID]shared.Resources, bool) {
	e.budget -= serviceCharge
	return e.clientPresident.EvaluateAllocationRequests(e.ResourceRequests, commonPool)
}

func (e *executive) requestAllocationRequest() {
	allocRequests := make(map[shared.ClientID]shared.Resources)
	for _, v := range getIslandAlive() {
		allocRequests[shared.ClientID(int(v))] = iigoClients[shared.ClientID(int(v))].CommonPoolResourceRequest()
	}
	AllocationAmountMapExport = allocRequests
	e.setAllocationRequest(allocRequests)

}

// replyAllocationRequest broadcasts the allocation of resources decided by the president
// to all islands alive
func (e *executive) replyAllocationRequest(commonPool shared.Resources) {
	e.budget -= serviceCharge
	allocationMap, requestsEvaluated := e.getAllocationRequests(commonPool)
	if requestsEvaluated {
		for _, v := range getIslandAlive() {
			d := shared.Communication{T: shared.CommunicationInt, IntegerData: int(allocationMap[shared.ClientID(int(v))])}
			data := make(map[shared.CommunicationFieldName]shared.Communication)
			data[shared.AllocationAmount] = d
			communicateWithIslands(shared.TeamIDs[int(v)], shared.TeamIDs[e.ID], data)
		}
	}
}

// appointNextSpeaker returns the island id of the island appointed to be speaker in the next turn.
func (e *executive) appointNextSpeaker(clientIDs []shared.ClientID) shared.ClientID {
	e.budget -= serviceCharge
	var election voting.Election
	election.ProposeElection(baseclient.Speaker, voting.Plurality)
	election.OpenBallot(clientIDs)
	election.Vote(iigoClients)
	return election.CloseBallot()
}

// withdrawSpeakerSalary withdraws the salary for speaker from the common pool.
func (e *executive) withdrawSpeakerSalary(gameState *gamestate.GameState) error {
	var speakerSalary = shared.Resources(rules.VariableMap[rules.SpeakerSalary].Values[0])
	var withdrawError = WithdrawFromCommonPool(speakerSalary, gameState)
	if withdrawError != nil {
		e.speakerSalary = speakerSalary
	}
	return withdrawError
}

// sendSpeakerSalary send speaker's salary to the speaker.
func (e *executive) sendSpeakerSalary(legislativeBranch *legislature) {
	amount, salaryPaid := e.clientPresident.PaySpeaker(e.speakerSalary)
	if salaryPaid {
		legislativeBranch.budget = amount
	}
}

func (e *executive) reset(val string) error {
	e.ID = 0
	e.clientPresident = nil
	e.budget = 0
	e.ResourceRequests = map[shared.ClientID]shared.Resources{}
	e.RulesProposals = []string{}
	e.speakerSalary = 0
	return nil
}

//requestRuleProposal asks each island alive for its rule proposal.
func (e *executive) requestRuleProposal() {
	e.budget -= serviceCharge
	var rules []string
	for _, v := range getIslandAlive() {
		rules = append(rules, iigoClients[shared.ClientID(int(v))].RuleProposal())
	}

	e.setRuleProposals(rules)
}

func getIslandAlive() []float64 {
	return rules.VariableMap[rules.IslandsAlive].Values
}
