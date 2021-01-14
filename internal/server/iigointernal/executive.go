package iigointernal

import (
	"fmt"
	"reflect"

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
	gameConf         *config.IIGOConfig
	PresidentID      shared.ClientID
	clientPresident  roles.President
	RulesProposals   []rules.RuleMatrix
	ResourceRequests map[shared.ClientID]shared.Resources
	iigoClients      map[shared.ClientID]baseclient.Client
	monitoring       *monitor
	logger           shared.Logger
}

func (e *executive) Logf(format string, a ...interface{}) {
	e.logger("[EXECUTIVE]: %v", fmt.Sprintf(format, a...))
}

// loadClientPresident checks client pointer is good and if not panics
func (e *executive) loadClientPresident(clientPresidentPointer roles.President) {
	if clientPresidentPointer == nil {
		panic(fmt.Sprintf("Client '%v' has loaded a nil president pointer", e.PresidentID))
	}
	e.clientPresident = clientPresidentPointer
}

// syncWithGame sets internal game state and configuration. Used to populate the executive struct for testing
func (e *executive) syncWithGame(gameState *gamestate.GameState, gameConf *config.IIGOConfig) {
	e.gameState = gameState
	e.gameConf = gameConf
}

// Get rule proposals to be voted on from remaining islands
// Called by orchestration
func (e *executive) setRuleProposals(rulesProposals []rules.RuleMatrix) {
	e.RulesProposals = rulesProposals
}

// Set approved resources request for all the remaining islands
// Called by orchestration
func (e *executive) setAllocationRequest(resourceRequests map[shared.ClientID]shared.Resources) {
	e.ResourceRequests = resourceRequests
}

// Get rules to be voted on to Speaker
// Called by orchestration at the end of the turn
func (e *executive) getRuleForSpeaker() (shared.PresidentReturnContent, error) {
	if !CheckEnoughInCommonPool(e.gameConf.GetRuleForSpeakerActionCost, e.gameState) {
		return shared.PresidentReturnContent{ContentType: shared.PresidentRuleProposal, ProposedRuleMatrix: rules.RuleMatrix{}, ActionTaken: false},
			errors.Errorf("Insufficient Budget in common Pool: getRuleForSpeaker")
	}

	returnRule := e.clientPresident.PickRuleToVote(e.RulesProposals)

	//Log Rule: obligation to select a rule if there is something in the proposal list
	if e.gameState.IIGORolesBudget[shared.President]-e.gameConf.GetRuleForSpeakerActionCost >= 0 {
		rulesProposed := len(e.RulesProposals) > 0

		variablesToCache := []rules.VariableFieldName{rules.IslandsProposedRules, rules.PresidentRuleProposal}
		valuesToCache := [][]float64{{boolToFloat(rulesProposed)}, {boolToFloat(returnRule.ActionTaken)}}
		e.monitoring.addToCache(e.PresidentID, variablesToCache, valuesToCache)
	}

	if returnRule.ActionTaken && (returnRule.ContentType == shared.PresidentRuleProposal) {
		if !e.incurServiceCharge(e.gameConf.GetRuleForSpeakerActionCost) {
			return returnRule, errors.Errorf("Insufficient Budget in common Pool: getRuleForSpeaker")
		}

	}

	//Log Rule: selected rule must come from rule proposal list
	if returnRule.ActionTaken && (returnRule.ContentType == shared.PresidentRuleProposal) {
		ruleChosenFromProposalList := false
		for _, a := range e.RulesProposals {
			if reflect.DeepEqual(a, returnRule.ProposedRuleMatrix) {
				ruleChosenFromProposalList = true
			}
		}
		variablesToCache := []rules.VariableFieldName{rules.RuleChosenFromProposalList}
		valuesToCache := [][]float64{{boolToFloat(ruleChosenFromProposalList)}}
		e.monitoring.addToCache(e.PresidentID, variablesToCache, valuesToCache)

	}

	return returnRule, nil
}

// broadcastTaxation broadcasts the tax amount decided by the president to all island still in the game.
func (e *executive) broadcastTaxation(islandsResources map[shared.ClientID]shared.ResourcesReport, aliveIslands []shared.ClientID) error {
	if !CheckEnoughInCommonPool(e.gameConf.BroadcastTaxationActionCost, e.gameState) {
		e.gameState.IIGOTaxAmount = make(map[shared.ClientID]shared.Resources)
		return errors.Errorf("Insufficient Budget in common Pool: broadcastTaxation")
	}
	taxMapReturn := e.getTaxMap(islandsResources)
	if taxMapReturn.ActionTaken && taxMapReturn.ContentType == shared.PresidentTaxation {
		if !e.incurServiceCharge(e.gameConf.BroadcastTaxationActionCost) {
			e.gameState.IIGOTaxAmount = make(map[shared.ClientID]shared.Resources)
			return errors.Errorf("Insufficient Budget in common Pool: broadcastTaxation")
		}
		e.gameState.IIGOTaxAmount = taxMapReturn.ResourceMap
		for islandID, amount := range taxMapReturn.ResourceMap {
			if Contains(aliveIslands, islandID) {
				e.sendDecision(islandID, amount, shared.IIGOTaxDecision)
			}
		}
		e.Logf("Tax: %v", taxMapReturn.ResourceMap)
	} else {
		// default case when president doesn't take an action. Update gamestate with tax = 0
		for _, islandID := range aliveIslands {
			e.sendNoDecision(islandID, shared.IIGOTaxDecision)
		}
		e.gameState.IIGOTaxAmount = make(map[shared.ClientID]shared.Resources)
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
		allocRequests[islandID] = e.iigoClients[islandID].CommonPoolResourceRequest()
	}
	e.gameState.IIGOAllocationMap = allocRequests
	e.setAllocationRequest(allocRequests)
	return nil
}

// replyAllocationRequest broadcasts the allocation of resources decided by the president
// to all islands alive
func (e *executive) replyAllocationRequest(commonPool shared.Resources) (bool, error) {
	if !CheckEnoughInCommonPool(e.gameConf.ReplyAllocationRequestsActionCost, e.gameState) {
		e.gameState.IIGOAllocationMade = false
		return false, errors.Errorf("Insufficient Budget in common Pool: replyAllocationRequest")
	}
	returnContent := e.getAllocationRequests(commonPool)
	allocationsMade := false
	if returnContent.ActionTaken && returnContent.ContentType == shared.PresidentAllocation {
		if !e.incurServiceCharge(e.gameConf.ReplyAllocationRequestsActionCost) {
			e.gameState.IIGOAllocationMade = false
			return false, errors.Errorf("Insufficient Budget in common Pool: replyAllocationRequest")
		}
		e.Logf("Resource Allocation: %v", returnContent.ResourceMap)
		allocationsMade = true
		e.gameState.IIGOAllocationMap = returnContent.ResourceMap
		for islandID, amount := range returnContent.ResourceMap {
			e.sendDecision(islandID, amount, shared.IIGOAllocationDecision)
		}
	} else {
		for islandID := range e.ResourceRequests {
			e.sendNoDecision(islandID, shared.IIGOAllocationDecision)
		}
	}
	e.gameState.IIGOAllocationMade = allocationsMade
	return allocationsMade, nil
}

// appointNextSpeaker returns the island ID of the island appointed to be Speaker in the next turn
func (e *executive) appointNextSpeaker(monitoring shared.MonitorResult, currentSpeaker shared.ClientID, allIslands []shared.ClientID) (shared.ClientID, error) {
	var election = voting.Election{
		Logger: e.logger,
	}
	var appointedSpeaker shared.ClientID
	allIslandsCopy1 := copyClientList(allIslands)
	electionSettings := e.clientPresident.CallSpeakerElection(monitoring, int(e.gameState.IIGOTurnsInPower[shared.Speaker]), allIslandsCopy1)

	//Log election rule
	termCondition := e.gameState.IIGOTurnsInPower[shared.Speaker] > e.gameConf.IIGOTermLengths[shared.Speaker]
	variablesToCache := []rules.VariableFieldName{rules.TermEnded, rules.ElectionHeld}
	valuesToCache := [][]float64{{boolToFloat(termCondition)}, {boolToFloat(electionSettings.HoldElection)}}
	e.monitoring.addToCache(e.PresidentID, variablesToCache, valuesToCache)

	if electionSettings.HoldElection {
		if !e.incurServiceCharge(e.gameConf.AppointNextSpeakerActionCost) {
			return e.gameState.SpeakerID, errors.Errorf("Insufficient Budget in common Pool: appointNextSpeaker")
		}
		election.ProposeElection(shared.Speaker, electionSettings.VotingMethod)
		allIslandsCopy2 := copyClientList(allIslands)
		election.OpenBallot(electionSettings.IslandsToVote, allIslandsCopy2)
		election.Vote(e.iigoClients)
		e.gameState.IIGOTurnsInPower[shared.Speaker] = 0
		electedSpeaker := election.CloseBallot(e.iigoClients)
		appointedSpeaker = e.clientPresident.DecideNextSpeaker(electedSpeaker)

		//Log rule: Must appoint elected role
		appointmentMatchesVote := appointedSpeaker == electedSpeaker
		variablesToCache := []rules.VariableFieldName{rules.AppointmentMatchesVote}
		valuesToCache := [][]float64{{boolToFloat(appointmentMatchesVote)}}
		e.monitoring.addToCache(e.PresidentID, variablesToCache, valuesToCache)
		e.Logf("Result of election for new Speaker: %v", appointedSpeaker)
	} else {
		appointedSpeaker = currentSpeaker
	}
	e.gameState.IIGOElection = append(e.gameState.IIGOElection, election.GetVotingInfo())
	return appointedSpeaker, nil
}

// sendSpeakerSalary conduct the transaction based on amount from client implementation
func (e *executive) sendSpeakerSalary() error {
	if e.clientPresident != nil {
		amountReturn := e.clientPresident.PaySpeaker()
		if amountReturn.ActionTaken && amountReturn.ContentType == shared.PresidentSpeakerSalary {
			// Subtract from common resources pool
			amountWithdraw, withdrawSuccess := WithdrawFromCommonPool(amountReturn.SpeakerSalary, e.gameState)

			if withdrawSuccess {
				// Pay into the client private resources pool
				depositIntoClientPrivatePool(amountWithdraw, e.gameState.SpeakerID, e.gameState)

				variablesToCache := []rules.VariableFieldName{rules.SpeakerPayment}
				valuesToCache := [][]float64{{float64(amountWithdraw)}}
				e.monitoring.addToCache(e.PresidentID, variablesToCache, valuesToCache)
				return nil
			}
		}
		variablesToCache := []rules.VariableFieldName{rules.SpeakerPaid}
		valuesToCache := [][]float64{{boolToFloat(amountReturn.ActionTaken)}}
		e.monitoring.addToCache(e.PresidentID, variablesToCache, valuesToCache)
	}
	return errors.Errorf("Cannot perform sendSpeakerSalary")
}

// Helper functions:
func (e *executive) getTaxMap(islandsResources map[shared.ClientID]shared.ResourcesReport) shared.PresidentReturnContent {
	return e.clientPresident.SetTaxationAmount(islandsResources)
}

//requestRuleProposal asks each island alive for its rule proposal.
func (e *executive) requestRuleProposal() error { //TODO: add checks for if immutable rules are changed(not allowed), if rule variables fields are changed(not allowed)
	if !e.incurServiceCharge(e.gameConf.RequestRuleProposalActionCost) {
		return errors.Errorf("Insufficient Budget in common Pool: requestRuleProposal")
	}

	var ruleProposals []rules.RuleMatrix
	for _, island := range e.getIslandAlive() {
		proposedRuleMatrix := e.iigoClients[shared.ClientID(int(island))].RuleProposal()
		if checkRuleIsValid(proposedRuleMatrix.RuleName, e.gameState.RulesInfo.AvailableRules) {
			ruleProposals = append(ruleProposals, proposedRuleMatrix)
		}
	}

	e.setRuleProposals(ruleProposals)
	return nil
}

func checkRuleIsValid(ruleName string, rulesCache map[string]rules.RuleMatrix) bool {
	_, valid := rulesCache[ruleName]
	return valid
}

func (e *executive) getIslandAlive() []float64 {
	return e.gameState.RulesInfo.VariableMap[rules.IslandsAlive].Values
}

// incur charges in both budget and commonpool for performing an actions
// actionID is typically the function name of the action
// only return error if we can't withdraw from common pool
func (e *executive) incurServiceCharge(cost shared.Resources) bool {
	_, ok := WithdrawFromCommonPool(cost, e.gameState)
	if ok {
		e.gameState.IIGORolesBudget[shared.President] -= cost
		if e.monitoring != nil {
			variablesToCache := []rules.VariableFieldName{rules.PresidentLeftoverBudget}
			valuesToCache := [][]float64{{float64(e.gameState.IIGORolesBudget[shared.President])}}
			e.monitoring.addToCache(e.PresidentID, variablesToCache, valuesToCache)
		}
	}
	return ok
}

func (e *executive) sendDecision(islandID shared.ClientID, value shared.Resources, communicationType shared.CommunicationFieldName) {
	data := make(map[shared.CommunicationFieldName]shared.CommunicationContent)
	var expected rules.VariableValuePair
	var decided rules.VariableValuePair

	if communicationType == shared.IIGOTaxDecision {
		expected = rules.MakeVariableValuePair(rules.ExpectedTaxContribution, []float64{float64(value)})
		decided = rules.MakeVariableValuePair(rules.TaxDecisionMade, []float64{boolToFloat(true)})
	} else {
		expected = rules.MakeVariableValuePair(rules.ExpectedAllocation, []float64{float64(value)})
		decided = rules.MakeVariableValuePair(rules.AllocationMade, []float64{boolToFloat(true)})
	}

	allocationToSend := shared.ValueDecision{
		Amount:       value,
		DecisionMade: decided,
		Expected:     expected,
	}

	data[communicationType] = shared.CommunicationContent{T: shared.CommunicationIIGOValue, IIGOValueData: allocationToSend}
	communicateWithIslands(e.iigoClients, islandID, shared.TeamIDs[e.PresidentID], data)
}

func (e *executive) sendNoDecision(islandID shared.ClientID, communicationType shared.CommunicationFieldName) {
	data := make(map[shared.CommunicationFieldName]shared.CommunicationContent)
	var decided rules.VariableValuePair
	if communicationType == shared.IIGOTaxDecision {
		decided = rules.MakeVariableValuePair(rules.TaxDecisionMade, []float64{boolToFloat(false)})
	} else {
		decided = rules.MakeVariableValuePair(rules.AllocationMade, []float64{boolToFloat(false)})
	}
	allocationToSend := shared.ValueDecision{
		DecisionMade: decided,
	}
	data[communicationType] = shared.CommunicationContent{T: shared.CommunicationIIGOValue, IIGOValueData: allocationToSend}
	communicateWithIslands(e.iigoClients, islandID, shared.TeamIDs[e.PresidentID], data)
}
