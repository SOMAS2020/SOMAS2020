package iigointernal

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"github.com/SOMAS2020/SOMAS2020/internal/common/voting"
	"github.com/pkg/errors"
)

type executive struct {
	gameState        *gamestate.GameState
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

// ---------Actions----------------

// Get rules to be voted on to Speaker
// Called by orchestration at the end of the turn
func (e *executive) getRuleForSpeaker() (string, bool) {
	if !e.incurServiceCharge("getRuleForSpeaker") {
		return "", false
	}

	return e.clientPresident.PickRuleToVote(e.RulesProposals)
}

// Send Tax map all the remaining islands
// Called by orchestration at the end of the turn
// broadcastTaxation broadcasts the tax amount decided by the president to all island still in the game.
func (e *executive) broadcastTaxation(islandsResources map[shared.ClientID]shared.Resources) error {
	if e.incurServiceCharge("broadcastTaxation") {
		return errors.Errorf("Insufficient Budget in common Pool: broadcastTaxation")
	}
	taxAmountMap, taxesCollected := e.getTaxMap(islandsResources)
	if taxesCollected {
		for _, island := range getIslandAlive() {
			d := shared.CommunicationContent{T: shared.CommunicationInt, IntegerData: int(taxAmountMap[shared.ClientID(int(island))])}
			data := make(map[shared.CommunicationFieldName]shared.CommunicationContent)
			data[shared.TaxAmount] = d
			communicateWithIslands(shared.TeamIDs[int(island)], shared.TeamIDs[e.ID], data)
		}
	}
	return nil
}

// Send Tax map all the remaining islands
// Called by orchestration at the end of the turn
func (e *executive) getAllocationRequests(commonPool shared.Resources) (map[shared.ClientID]shared.Resources, bool) {
	return e.clientPresident.EvaluateAllocationRequests(e.ResourceRequests, commonPool)
}

//requestRuleProposal asks each island alive for its rule proposal.
func (e *executive) requestRuleProposal() error {
	if !e.incurServiceCharge("requestRuleProposal") {
		return errors.Errorf("Insufficient Budget in common Pool: broadcastTaxation")
	}

	var rules []string
	for _, island := range getIslandAlive() {
		rules = append(rules, iigoClients[shared.ClientID(int(island))].RuleProposal())
	}

	e.setRuleProposals(rules)
	return nil
}

// ------ Actions end -----------
func (e *executive) requestAllocationRequest() {
	allocRequests := make(map[shared.ClientID]shared.Resources)
	for _, island := range getIslandAlive() {
		allocRequests[shared.ClientID(int(island))] = iigoClients[shared.ClientID(int(island))].CommonPoolResourceRequest()
	}
	AllocationAmountMapExport = allocRequests
	e.setAllocationRequest(allocRequests)

}

// replyAllocationRequest broadcasts the allocation of resources decided by the president
// to all islands alive
func (e *executive) replyAllocationRequest(commonPool shared.Resources) error {
	// If request costs, why does the reply cost? (Need to update return types)
	if !e.incurServiceCharge("replyAllocationRequest") {
		return errors.Errorf("Insufficient Budget in common Pool: replyAllocationRequest")
	}

	allocationMap, requestsEvaluated := e.getAllocationRequests(commonPool)
	if requestsEvaluated {
		for _, island := range getIslandAlive() {
			d := shared.CommunicationContent{T: shared.CommunicationInt, IntegerData: int(allocationMap[shared.ClientID(int(island))])}
			data := make(map[shared.CommunicationFieldName]shared.CommunicationContent)
			data[shared.AllocationAmount] = d
			communicateWithIslands(shared.TeamIDs[int(island)], shared.TeamIDs[e.ID], data)
		}
	}
	return nil
}

// appointNextSpeaker returns the island id of the island appointed to be speaker in the next turn.
// appointing new role should be free now
func (e *executive) appointNextSpeaker(clientIDs []shared.ClientID) (shared.ClientID, error) {
	if !e.incurServiceCharge("appointNextSpeaker") {
		return e.ID, errors.Errorf("Insufficient Budget in common Pool: appointNextSpeaker")
	}
	var election voting.Election
	election.ProposeElection(baseclient.Speaker, voting.Plurality)
	election.OpenBallot(clientIDs)
	election.Vote(iigoClients)
	return election.CloseBallot(), nil
}

// withdrawSpeakerSalary withdraws the salary for speaker from the common pool.
func (e *executive) withdrawSpeakerSalary() bool {
	var speakerSalary = shared.Resources(rules.VariableMap[rules.SpeakerSalary].Values[0])
	var withdrawnAmount, withdrawSuccesful = WithdrawFromCommonPool(speakerSalary, e.gameState)
	e.speakerSalary = withdrawnAmount
	return withdrawSuccesful
}

// sendSpeakerSalary send speaker's salary to the speaker.
func (e *executive) sendSpeakerSalary(legislativeBranch *legislature) {
	if e.clientPresident != nil {
		amount, salaryPaid := e.clientPresident.PaySpeaker(e.speakerSalary)
		if salaryPaid {
			legislativeBranch.budget = amount
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

// Helper functions:
func (e *executive) getTaxMap(islandsResources map[shared.ClientID]shared.Resources) (map[shared.ClientID]shared.Resources, bool) {
	return e.clientPresident.SetTaxationAmount(islandsResources)
}

func getIslandAlive() []float64 {
	return rules.VariableMap[rules.IslandsAlive].Values
}

// incur charges in both budget and commonpool for performing an actions
// actionID is typically the function name of the action
// only return error if we can't withdraw from common pool
func (e *executive) incurServiceCharge(actionID string) bool {
	cost := config.GameConfig().IIGOConfig.ExecutiveActionCost[actionID]
	_, ok := WithdrawFromCommonPool(cost, e.gameState)
	if ok {
		e.budget -= cost
	}
	return ok
}
