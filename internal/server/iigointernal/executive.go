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
	"gonum.org/v1/gonum/mat"
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

type conversionType int

const (
	tax conversionType = iota
	allocation
)

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
		for islandID, amount := range taxMapReturn.ResourceMap {
			if Contains(aliveIslands, islandID) {
				e.sendTax(islandID, amount)
			}
		}
	} else {
		// default case when president doesn't take an action. send tax = 0
		for _, islandID := range aliveIslands {
			e.sendNoTax(islandID)
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
	if returnContent.ActionTaken && returnContent.ContentType == shared.PresidentAllocation {
		if !e.incurServiceCharge(e.gameConf.ReplyAllocationRequestsActionCost) {
			return false, errors.Errorf("Insufficient Budget in common Pool: replyAllocationRequest")
		}
		allocationsMade = true
		for islandID, amount := range returnContent.ResourceMap {
			e.sendAllocation(islandID, amount)
		}
	} else {
		for islandID := range e.ResourceRequests {
			e.sendNoAllocation(islandID)
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
		election.OpenBallot(electionsettings.IslandsToVote)
		election.Vote(iigoClients)
		e.speakerTurnsInPower = 0
		nextSpeaker = election.CloseBallot()
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

// convertAmount takes the amount of tax/allocation and converts it into appropriate variable and rule ready to be sent to the client
func convertAmount(amount shared.Resources, amountType conversionType) (rules.VariableValuePair, rules.RuleMatrix) {
	var reqVar rules.VariableFieldName
	var sentVar rules.VariableFieldName
	var ruleVariables []rules.VariableFieldName
	if amountType == tax {
		reqVar, sentVar = rules.IslandTaxContribution, rules.ExpectedTaxContribution
		// Rule in form IslandTaxContribution - ExpectedTaxContribution >= 0
		ruleVariables = []rules.VariableFieldName{reqVar, sentVar}
	} else if amountType == allocation {
		reqVar, sentVar = rules.IslandAllocation, rules.ExpectedAllocation
		// Rule in form ExpectedAllocation - IslandAllocation >= 0
		ruleVariables = []rules.VariableFieldName{sentVar, reqVar}
	}

	v := []float64{1, -1, 0}
	aux := []float64{2}

	rowLength := len(ruleVariables) + 1
	nrows := len(v) / rowLength

	CoreMatrix := mat.NewDense(nrows, rowLength, v)
	AuxiliaryVector := mat.NewVecDense(nrows, aux)

	retRule := rules.RuleMatrix{
		RuleName:          fmt.Sprintf("%s = %.2f", sentVar, amount),
		RequiredVariables: ruleVariables,
		ApplicableMatrix:  *CoreMatrix,
		AuxiliaryVector:   *AuxiliaryVector,
		Mutable:           false,
	}

	retVar := rules.VariableValuePair{
		VariableName: sentVar,
		Values:       []float64{float64(amount)},
	}

	return retVar, retRule
}

func (e *executive) sendTax(islandID shared.ClientID, taxAmount shared.Resources) {
	data := make(map[shared.CommunicationFieldName]shared.CommunicationContent)
	expectedVariable, taxRule := convertAmount(taxAmount, tax)
	taxToSend := shared.TaxDecision{
		TaxAmount:   taxAmount,
		TaxRule:     taxRule,
		ExpectedTax: expectedVariable,
		TaxDecided:  true,
	}

	data[shared.Tax] = shared.CommunicationContent{T: shared.CommunicationTax, TaxDecision: taxToSend}
	communicateWithIslands(islandID, shared.TeamIDs[e.PresidentID], data)
}

func (e *executive) sendNoTax(islandID shared.ClientID) {
	data := make(map[shared.CommunicationFieldName]shared.CommunicationContent)
	taxToSend := shared.TaxDecision{TaxDecided: false}
	data[shared.Tax] = shared.CommunicationContent{T: shared.CommunicationTax, TaxDecision: taxToSend}
	communicateWithIslands(islandID, shared.TeamIDs[e.PresidentID], data)
}

func (e *executive) sendAllocation(islandID shared.ClientID, allocationAmount shared.Resources) {
	data := make(map[shared.CommunicationFieldName]shared.CommunicationContent)
	expectedVariable, allocationRule := convertAmount(allocationAmount, allocation)
	allocationToSend := shared.AllocationDecision{
		AllocationAmount:   allocationAmount,
		AllocationRule:     allocationRule,
		ExpectedAllocation: expectedVariable,
		AllocationDecided:  true,
	}

	data[shared.Allocation] = shared.CommunicationContent{T: shared.CommunicationAllocation, AllocationDecision: allocationToSend}
	communicateWithIslands(islandID, shared.TeamIDs[e.PresidentID], data)
}

func (e *executive) sendNoAllocation(islandID shared.ClientID) {
	data := make(map[shared.CommunicationFieldName]shared.CommunicationContent)
	allocationToSend := shared.AllocationDecision{AllocationDecided: false}
	data[shared.Allocation] = shared.CommunicationContent{T: shared.CommunicationAllocation, AllocationDecision: allocationToSend}
	communicateWithIslands(islandID, shared.TeamIDs[e.PresidentID], data)
}
