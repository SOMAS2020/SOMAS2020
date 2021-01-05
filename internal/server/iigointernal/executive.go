package iigointernal

import (
	"fmt"

	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"github.com/SOMAS2020/SOMAS2020/internal/common/voting"
	"github.com/pkg/errors"
)

type executive struct {
	gameState           *gamestate.GameState
	gameConf            *config.IIGOConfig
	PresidentID         shared.ClientID
	clientPresident     roles.President
	speakerSalary       shared.Resources
	RulesProposals      []string
	ResourceRequests    map[shared.ClientID]shared.Resources
	speakerTurnsInPower int
}

// loadClientPresident checks client pointer is good and if not panics
func (e *executive) loadClientPresident(clientPresidentPointer roles.President) {
	if clientPresidentPointer == nil {
		panic(fmt.Sprintf("Client '%v' has loaded a nil president pointer", e.PresidentID))
	}
	e.clientPresident = clientPresidentPointer
}

func (e *executive) syncWithGame(gameState *gamestate.GameState, gameConf *config.IIGOConfig) {
	e.gameState = gameState
	e.gameConf = gameConf
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

// setGameState is used for setting the game state of the executive branch
// Called for testing.
func (e *executive) setGameState(g *gamestate.GameState) {
	e.gameState = g
}

// Get rules to be voted on to Speaker
// Called by orchestration at the end of the turn
func (e *executive) getRuleForSpeaker() (shared.PresidentReturnContent, error) {
	if !CheckEnoughInCommonPool(e.gameConf.GetRuleForSpeakerActionCost, e.gameState) {
		return shared.PresidentReturnContent{ContentType: shared.PresidentRuleProposal, ProposedRule: "", ActionTaken: false},
			errors.Errorf("Insufficient Budget in common Pool: broadcastTaxation")
	}

	returnRule := e.clientPresident.PickRuleToVote(e.RulesProposals)

	if returnRule.ActionTaken && (returnRule.ContentType == shared.PresidentRuleProposal) {
		if !e.incurServiceCharge(e.gameConf.GetRuleForSpeakerActionCost) {
			return returnRule, errors.Errorf("Insufficient Budget in common Pool: getRuleForSpeaker")
		}
	}

	return returnRule, nil
}

// broadcastTaxation broadcasts the tax amount decided by the president to all island still in the game.
func (e *executive) broadcastTaxation(islandsResources map[shared.ClientID]shared.ResourcesReport, aliveIslands []shared.ClientID) error {
	if !CheckEnoughInCommonPool(e.gameConf.BroadcastTaxationActionCost, e.gameState) {
		return errors.Errorf("Insufficient Budget in common Pool: broadcastTaxation")
	}
	taxMapReturn := e.getTaxMap(islandsResources)
	if taxMapReturn.ActionTaken && taxMapReturn.ContentType == shared.PresidentTaxation {
		if !e.incurServiceCharge(e.gameConf.BroadcastTaxationActionCost) {
			return errors.Errorf("Insufficient Budget in common Pool: broadcastTaxation")
		}
		for islandID, resourceAmount := range taxMapReturn.ResourceMap {
			if Contains(aliveIslands, islandID) {
				d := shared.CommunicationContent{T: shared.CommunicationInt, IntegerData: int(resourceAmount)}
				data := make(map[shared.CommunicationFieldName]shared.CommunicationContent)
				data[shared.TaxAmount] = d
				communicateWithIslands(islandID, shared.TeamIDs[e.PresidentID], data)
			}
		}
	}
	return nil
}

// Send Tax map all the remaining islands
// Called by orchestration at the end of the turn
func (e *executive) getAllocationRequests(commonPool shared.Resources) shared.PresidentReturnContent {
	return e.clientPresident.EvaluateAllocationRequests(e.ResourceRequests, commonPool)
}

func (e *executive) requestAllocationRequest(aliveIslands []shared.ClientID) error {
	// Listening to islands requests incurs cost to the president. Set the cost to 0?
	if !e.incurServiceCharge(e.gameConf.RequestAllocationRequestActionCost) {
		return errors.Errorf("Insufficient Budget in common Pool: requestAllocationRequest")
	}
	allocRequests := make(map[shared.ClientID]shared.Resources)
	for _, islandID := range aliveIslands {
		allocRequests[islandID] = iigoClients[islandID].CommonPoolResourceRequest()
	}
	AllocationAmountMapExport = allocRequests
	e.setAllocationRequest(allocRequests)
	return nil
}

// replyAllocationRequest broadcasts the allocation of resources decided by the president
// to all islands alive
func (e *executive) replyAllocationRequest(commonPool shared.Resources) (bool, error) {
	if !CheckEnoughInCommonPool(e.gameConf.ReplyAllocationRequestsActionCost, e.gameState) {
		return false, errors.Errorf("Insufficient Budget in common Pool: replyAllocationRequest")
	}
	returnContent := e.getAllocationRequests(commonPool)
	allocationsMade := false
	if returnContent.ActionTaken {
		if !e.incurServiceCharge(e.gameConf.ReplyAllocationRequestsActionCost) {
			return false, errors.Errorf("Insufficient Budget in common Pool: replyAllocationRequest")
		}
		allocationsMade = true
		for _, island := range getIslandAlive() {
			d := shared.CommunicationContent{T: shared.CommunicationInt, IntegerData: int(returnContent.ResourceMap[shared.ClientID(int(island))])}
			data := make(map[shared.CommunicationFieldName]shared.CommunicationContent)
			data[shared.AllocationAmount] = d
			communicateWithIslands(shared.TeamIDs[int(island)], shared.TeamIDs[int(e.PresidentID)], data)
		}
	}
	return allocationsMade, nil
}

// appointNextSpeaker returns the island ID of the island appointed to be Speaker in the next turn
func (e *executive) appointNextSpeaker(monitoring shared.MonitorResult, currentSpeaker shared.ClientID, allIslands []shared.ClientID) (shared.ClientID, error) {
	var election voting.Election
	var nextSpeaker shared.ClientID
	electionsettings := e.clientPresident.CallSpeakerElection(monitoring, e.speakerTurnsInPower, allIslands)
	if electionsettings.HoldElection {
		if !e.incurServiceCharge(e.gameConf.AppointNextSpeakerActionCost) {
			return e.gameState.SpeakerID, errors.Errorf("Insufficient Budget in common Pool: appointNextSpeaker")
		}
		election.ProposeElection(shared.Speaker, electionsettings.VotingMethod)
		election.OpenBallot(electionsettings.IslandsToVote, iigoClients)
		election.Vote(iigoClients)
		e.speakerTurnsInPower = 0
		nextSpeaker = election.CloseBallot(iigoClients)
		nextSpeaker = e.clientPresident.DecideNextSpeaker(nextSpeaker)
	} else {
		e.speakerTurnsInPower++
		nextSpeaker = currentSpeaker
	}
	return nextSpeaker, nil
}

// sendSpeakerSalary conduct the transaction based on amount from client implementation
func (e *executive) sendSpeakerSalary() error {
	if e.clientPresident != nil {
		amountReturn := e.clientPresident.PaySpeaker(e.speakerSalary)
		if amountReturn.ActionTaken && amountReturn.ContentType == shared.PresidentSpeakerSalary {
			// Subtract from common resources pool
			amountWithdraw, withdrawSuccess := WithdrawFromCommonPool(amountReturn.SpeakerSalary, e.gameState)

			if withdrawSuccess {
				// Pay into the client private resources pool
				depositIntoClientPrivatePool(amountWithdraw, e.gameState.SpeakerID, e.gameState)
				return nil
			}
		}
	}
	return errors.Errorf("Cannot perform sendSpeakerSalary")
}

func (e *executive) reset(val string) {
	e.PresidentID = 0
	e.clientPresident = nil
	e.ResourceRequests = map[shared.ClientID]shared.Resources{}
	e.RulesProposals = []string{}
	e.speakerSalary = 0
}

// Helper functions:
func (e *executive) getTaxMap(islandsResources map[shared.ClientID]shared.ResourcesReport) shared.PresidentReturnContent {
	return e.clientPresident.SetTaxationAmount(islandsResources)
}

//requestRuleProposal asks each island alive for its rule proposal.
func (e *executive) requestRuleProposal() error {
	if !e.incurServiceCharge(e.gameConf.RequestRuleProposalActionCost) {
		return errors.Errorf("Insufficient Budget in common Pool: broadcastTaxation")
	}

	var ruleProposals []string
	for _, island := range getIslandAlive() {
		proposedRule := iigoClients[shared.ClientID(int(island))].RuleProposal()
		if checkRuleIsValid(proposedRule, rules.AvailableRules) {
			ruleProposals = append(ruleProposals, proposedRule)
		}
	}

	e.setRuleProposals(ruleProposals)
	return nil
}

func checkRuleIsValid(ruleName string, rulesCache map[string]rules.RuleMatrix) bool {
	_, valid := rulesCache[ruleName]
	return valid
}

func getIslandAlive() []float64 {
	return rules.VariableMap[rules.IslandsAlive].Values
}

// incur charges in both budget and commonpool for performing an actions
// actionID is typically the function name of the action
// only return error if we can't withdraw from common pool
func (e *executive) incurServiceCharge(cost shared.Resources) bool {
	_, ok := WithdrawFromCommonPool(cost, e.gameState)
	if ok {
		e.gameState.IIGORolesBudget[shared.President] -= cost
	}
	return ok
}
