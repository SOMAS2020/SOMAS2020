package iigointernal

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"github.com/SOMAS2020/SOMAS2020/internal/common/voting"
	"github.com/pkg/errors"
)

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
	ID:               0,
	speakerSalary:    0,
	ResourceRequests: nil,
}

// SpeakerIDGlobal is the single source of truth for speaker ID (MVP)
var SpeakerIDGlobal shared.ClientID = 0

// JudgeIDGlobal is the single source of truth for judge ID (MVP)
var JudgeIDGlobal shared.ClientID = 1

// PresidentIDGlobal is the single source of truth for president ID (MVP)
var PresidentIDGlobal shared.ClientID = 2

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
func RunIIGO(g *gamestate.GameState, clientMap *map[shared.ClientID]baseclient.Client) error {

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
	judicialBranch.JudgeID = JudgeIDGlobal
	legislativeBranch.SpeakerID = SpeakerIDGlobal
	executiveBranch.ID = PresidentIDGlobal // TODO:- change ID to president id

	// Set judgePointer
	judgePointer = iigoClients[JudgeIDGlobal].GetClientJudgePointer()
	// Set speakerPointer
	speakerPointer = iigoClients[SpeakerIDGlobal].GetClientSpeakerPointer()
	// Set presidentPointer
	presidentPointer = iigoClients[PresidentIDGlobal].GetClientPresidentPointer()

	// Initialise iigointernal with their clientVersions
	judicialBranch.clientJudge = judgePointer
	executiveBranch.clientPresident = presidentPointer
	legislativeBranch.clientSpeaker = speakerPointer

	// Pay salaries into budgets

	errorJudicial := judicialBranch.sendPresidentSalary()
	errorLegislative := legislativeBranch.sendJudgeSalary()
	errorExecutive := executiveBranch.sendSpeakerSalary()
	// Throw error
	if errorJudicial != nil || errorLegislative != nil || errorExecutive != nil {
		return errors.Errorf("Cannot pay IIGO salary")
	}

	// 1 Judge actions - inspect history
	_, historyInspected := judicialBranch.inspectHistory()

	// 2 President actions
	resourceReports := map[shared.ClientID]shared.Resources{}
	var aliveClientIds []shared.ClientID
	for _, v := range rules.VariableMap[rules.IslandsAlive].Values {
		aliveClientIds = append(aliveClientIds, shared.ClientID(int(v)))
		resourceReports[shared.ClientID(int(v))] = iigoClients[shared.ClientID(int(v))].ResourceReport()
	}

	// Throw error if any of the actions returns error
	insufficientBudget := executiveBranch.broadcastTaxation(resourceReports)
	if insufficientBudget == nil {
		insufficientBudget = executiveBranch.requestAllocationRequest()
	}
	if insufficientBudget == nil {
		insufficientBudget = executiveBranch.replyAllocationRequest(g.CommonPool)
	}
	if insufficientBudget == nil {
		insufficientBudget = executiveBranch.requestRuleProposal()
	}
	if insufficientBudget != nil {
		return errors.Errorf("Common pool resources insufficient for executiveBranch actions")
	}
	ruleToVote, ruleSelected := executiveBranch.getRuleForSpeaker()

	// 3 Speaker actions

	//TODO:- shouldn't updateRules be called here?
	if ruleSelected {
		insufficientBudget := legislativeBranch.setRuleToVote(ruleToVote)
		if insufficientBudget == nil {
			_, insufficientBudget = legislativeBranch.setVotingResult(aliveClientIds)
		}
		if insufficientBudget == nil {
			insufficientBudget = legislativeBranch.announceVotingResult()
		}
		if insufficientBudget != nil {
			return errors.Errorf("Common pool resources insufficient for legislativeBranch actions")
		}
	}

	// 4 Declare performance (Judge) (in future all the iigointernal)
	if historyInspected {
		judicialBranch.declarePresidentPerformanceWrapped()

		judicialBranch.declareSpeakerPerformanceWrapped()
	}

	// TODO:- at the moment, these are action (and cost resources) but should they?
	// Get new Judge ID
	JudgeIDGlobal, _ = legislativeBranch.appointNextJudge(aliveClientIds)
	// Get new Speaker ID
	SpeakerIDGlobal, _ = executiveBranch.appointNextSpeaker(aliveClientIds)
	// Get new President ID
	PresidentIDGlobal, _ = judicialBranch.appointNextPresident(aliveClientIds)

	return nil
}
