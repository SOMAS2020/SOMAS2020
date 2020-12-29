package iigointernal

import (
	"fmt"

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

// loadClientJudge checks client pointer is good and if not panics
func (e *executive) loadClientPresident(clientPresidentPointer roles.President) {
	if clientPresidentPointer == nil {
		panic(fmt.Sprintf("Client '%v' has loaded a nil president pointer", e.ID))
	}
	e.clientPresident = clientPresidentPointer
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
func (e *executive) getRuleForSpeaker() shared.PresidentReturnContent {
	e.budget -= serviceCharge
	return e.clientPresident.PickRuleToVote(e.RulesProposals)
}

// Send Tax map all the remaining islands
// Called by orchestration at the end of the turn
func (e *executive) getTaxMap(islandsResources map[shared.ClientID]shared.Resources) shared.PresidentReturnContent {
	e.budget -= serviceCharge
	return e.clientPresident.SetTaxationAmount(islandsResources)
}

// broadcastTaxation broadcasts the tax amount decided by the president to all island still in the game.
func (e *executive) broadcastTaxation(islandsResources map[shared.ClientID]shared.Resources, aliveIslands []shared.ClientID) {
	e.budget -= serviceCharge
	taxMapReturn := e.getTaxMap(islandsResources)
	if taxMapReturn.ActionTaken && taxMapReturn.ContentType == shared.PresidentTaxation {
		for _, islandID := range aliveIslands {
			d := shared.CommunicationContent{T: shared.CommunicationInt, IntegerData: int(taxMapReturn.ResourceMap[islandID])}
			data := make(map[shared.CommunicationFieldName]shared.CommunicationContent)
			data[shared.TaxAmount] = d
			communicateWithIslands(islandID, shared.TeamIDs[e.ID], data)
		}
	}
}

// Send Tax map all the remaining islands
// Called by orchestration at the end of the turn
func (e *executive) getAllocationRequests(commonPool shared.Resources) shared.PresidentReturnContent {
	e.budget -= serviceCharge
	return e.clientPresident.EvaluateAllocationRequests(e.ResourceRequests, commonPool)
}

func (e *executive) requestAllocationRequest(aliveIslands []shared.ClientID) {
	allocRequests := make(map[shared.ClientID]shared.Resources)
	for _, islandID := range aliveIslands {
		allocRequests[islandID] = iigoClients[islandID].CommonPoolResourceRequest()
	}
	AllocationAmountMapExport = allocRequests
	e.setAllocationRequest(allocRequests)

}

// replyAllocationRequest broadcasts the allocation of resources decided by the president
// to all islands alive
func (e *executive) replyAllocationRequest(commonPool shared.Resources, aliveIslands []shared.ClientID) {
	e.budget -= serviceCharge
	allocationMapReturn := e.getAllocationRequests(commonPool)
	if allocationMapReturn.ActionTaken && allocationMapReturn.ContentType == shared.PresidentAllocation {
		for _, islandID := range aliveIslands {
			d := shared.CommunicationContent{T: shared.CommunicationInt, IntegerData: int(allocationMapReturn.ResourceMap[islandID])}
			data := make(map[shared.CommunicationFieldName]shared.CommunicationContent)
			data[shared.AllocationAmount] = d
			communicateWithIslands(islandID, shared.TeamIDs[e.ID], data)
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
func (e *executive) withdrawSpeakerSalary(gameState *gamestate.GameState) bool {
	var speakerSalary = shared.Resources(rules.VariableMap[rules.SpeakerSalary].Values[0])
	var withdrawnAmount, withdrawSuccesful = WithdrawFromCommonPool(speakerSalary, gameState)
	e.speakerSalary = withdrawnAmount
	return withdrawSuccesful
}

// sendSpeakerSalary send speaker's salary to the speaker.
func (e *executive) sendSpeakerSalary(legislativeBranch *legislature) {
	if e.clientPresident != nil {
		amountReturn := e.clientPresident.PaySpeaker(e.speakerSalary)
		if amountReturn.ActionTaken && amountReturn.ContentType == shared.PresidentSpeakerSalary {
			legislativeBranch.budget = amountReturn.SpeakerSalary
		}
		return
	}
	legislativeBranch.budget = e.speakerSalary
}

func (e *executive) reset(val string) {
	e.ID = 0
	e.clientPresident = nil
	e.budget = 0
	e.ResourceRequests = map[shared.ClientID]shared.Resources{}
	e.RulesProposals = []string{}
	e.speakerSalary = 0
}

//requestRuleProposal asks each island alive for its rule proposal.
func (e *executive) requestRuleProposal(aliveIslands []shared.ClientID) {
	e.budget -= serviceCharge
	var rules []string
	for _, islandID := range aliveIslands {
		rules = append(rules, iigoClients[islandID].RuleProposal())
	}

	e.setRuleProposals(rules)
}
