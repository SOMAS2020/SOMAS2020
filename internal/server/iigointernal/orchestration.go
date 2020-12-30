package iigointernal

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"github.com/SOMAS2020/SOMAS2020/internal/common/voting"
)

// Single source of truth for all action
var actionCost config.IIGOConfig

// featureJudge is an instantiation of the Judge interface
// with both the Base Judge features and a reference to client judges
var judicialBranch = judiciary{
	gameState:         nil,
	JudgeID:           0,
	presidentSalary:   0,
	BallotID:          0,
	ResAllocID:        0,
	speakerID:         0,
	presidentID:       0,
	EvaluationResults: nil,
}

// featureSpeaker is an instantiation of the Speaker interface
// with both the baseSpeaker features and a reference to client speakers
var legislativeBranch = legislature{
	gameState:    nil,
	SpeakerID:    0,
	judgeSalary:  0,
	ruleToVote:   "",
	ballotBox:    voting.BallotBox{},
	votingResult: false,
}

// featurePresident is an instantiation of the President interface
// with both the basePresident features and a reference to client presidents
var executiveBranch = executive{
	gameState:        nil,
	PresidentID:      0,
	speakerSalary:    0,
	ResourceRequests: nil,
}

// TaxAmountMapExport is a local tax amount cache for checking of rules
var TaxAmountMapExport map[shared.ClientID]shared.Resources

// AllocationAmountMapExport is a local allocation map for checking of rules
var AllocationAmountMapExport map[shared.ClientID]shared.Resources

// Pointers allow clients to customise implementations of mutable functions
var judgePointer roles.Judge = nil
var speakerPointer roles.Speaker = nil
var presidentPointer roles.President = nil

// iigoClients holds pointers to all the clients
var iigoClients map[shared.ClientID]baseclient.Client

// RunIIGO runs all iigo function in sequence
func RunIIGO(g *gamestate.GameState, clientMap *map[shared.ClientID]baseclient.Client) (IIGOSuccessful bool, StatusDescription string) {

	iigoClients = *clientMap

	// Increments the budget by a constant 100
	// TODO:- the constant should be retrieved from the rules
	g.IIGORolesBudget["president"] += 100
	g.IIGORolesBudget["judge"] += 100
	g.IIGORolesBudget["speaker"] += 100

	// Pass in gamestate -
	// So that we don't have to pass gamestate as arguement in every function in roles
	judicialBranch.gameState = g
	legislativeBranch.gameState = g
	executiveBranch.gameState = g

	// Initialise IDs
	judicialBranch.JudgeID = g.JudgeID
	legislativeBranch.SpeakerID = g.SpeakerID
	executiveBranch.PresidentID = g.PresidentID

	// Set judgePointer
	judgePointer = iigoClients[g.JudgeID].GetClientJudgePointer()
	// Set speakerPointer
	speakerPointer = iigoClients[g.SpeakerID].GetClientSpeakerPointer()
	// Set presidentPointer
	presidentPointer = iigoClients[g.PresidentID].GetClientPresidentPointer()

	// Initialise iigointernal with their clientVersions
	judicialBranch.loadClientJudge(judgePointer)
	executiveBranch.loadClientPresident(presidentPointer)
	legislativeBranch.loadClientSpeaker(speakerPointer)

	// Pay salaries into budgets

	errorJudicial := judicialBranch.sendPresidentSalary()
	errorLegislative := legislativeBranch.sendJudgeSalary()
	errorExecutive := executiveBranch.sendSpeakerSalary()
	// Throw error
	if errorJudicial != nil || errorLegislative != nil || errorExecutive != nil {
		return false, "Cannot pay IIGO salary"
	}

	// 1 Judge actions - inspect history
	_, historyInspected := judicialBranch.inspectHistory(g.IIGOHistory)

	// 2 President actions
	resourceReports := map[shared.ClientID]shared.ResourcesReport{}
	aliveClientIds := []shared.ClientID{}
	for clientID, clientGameState := range g.ClientInfos {
		if clientGameState.LifeStatus != shared.Dead {
			aliveClientIds = append(aliveClientIds, clientID)
			resourceReports[clientID] = iigoClients[clientID].ResourceReport()
		}
	}

	// Throw error if any of the actions returns error
	insufficientBudget := executiveBranch.broadcastTaxation(resourceReports, aliveClientIds)
	if insufficientBudget == nil {
		insufficientBudget = executiveBranch.requestAllocationRequest(aliveClientIds)
	}
	if insufficientBudget == nil {
		insufficientBudget = executiveBranch.replyAllocationRequest(g.CommonPool, aliveClientIds)
	}
	if insufficientBudget == nil {
		insufficientBudget = executiveBranch.requestRuleProposal(aliveClientIds)
	}
	if insufficientBudget != nil {
		return false, "Common pool resources insufficient for executiveBranch actions"
	}
	ruleToVoteReturn := executiveBranch.getRuleForSpeaker()

	// 3 Speaker actions

	//TODO:- shouldn't updateRules be called here?
	if ruleToVoteReturn.ActionTaken && ruleToVoteReturn.ContentType == shared.PresidentRuleProposal {
		insufficientBudget := legislativeBranch.setRuleToVote(ruleToVoteReturn.ProposedRule)
		if insufficientBudget == nil {
			insufficientBudget = legislativeBranch.setVotingResult(aliveClientIds)
		}
		if insufficientBudget == nil {
			insufficientBudget = legislativeBranch.announceVotingResult()
		}
		if insufficientBudget != nil {
			return false, "Common pool resources insufficient for legislativeBranch actions"
		}
	}

	// 4 Declare performance (Judge) (in future all the iigointernal)
	if historyInspected {
		judicialBranch.declarePresidentPerformanceWrapped()

		judicialBranch.declareSpeakerPerformanceWrapped()
	}

	// TODO:- at the moment, these are action (and cost resources) but should they?
	var appointJudgeError, appointSpeakerError, appointPresidentError error
	// Get new Judge ID
	g.JudgeID, appointJudgeError = legislativeBranch.appointNextJudge(aliveClientIds)
	if appointJudgeError != nil {
		return false, "Judge was not apointed by the Speaker. Insufficient budget"
	}
	// Get new Speaker ID
	g.SpeakerID, appointSpeakerError = executiveBranch.appointNextSpeaker(aliveClientIds)
	if appointSpeakerError != nil {
		return false, "Speaker was not apointed by the President. Insufficient budget"
	}
	// Get new President ID
	g.PresidentID, appointPresidentError = judicialBranch.appointNextPresident(aliveClientIds)
	if appointPresidentError != nil {
		return false, "President was not apointed by the Judge. Insufficient budget"
	}
	return true, "IIGO Run Successful"
}
