package iigointernal

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"github.com/SOMAS2020/SOMAS2020/internal/common/voting"
)

const serviceCharge = shared.Resources(10)

// featureJudge is an instantiation of the Judge interface
// with both the Base Judge features and a reference to client judges
var judicialBranch = judiciary{
	JudgeID:           0,
	budget:            0,
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
	SpeakerID:    0,
	budget:       0,
	judgeSalary:  0,
	ruleToVote:   "",
	ballotBox:    voting.BallotBox{},
	votingResult: false,
}

// featurePresident is an instantiation of the President interface
// with both the basePresident features and a reference to client presidents
var executiveBranch = executive{
	ID:               0,
	budget:           0,
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
func RunIIGO(g *gamestate.GameState, clientMap *map[shared.ClientID]baseclient.Client) (IIGOSuccessful bool, StatusDescription string) {

	iigoClients = *clientMap

	// Initialise IDs
	judicialBranch.JudgeID = JudgeIDGlobal
	legislativeBranch.SpeakerID = SpeakerIDGlobal
	executiveBranch.ID = PresidentIDGlobal

	// Set judgePointer
	judgePointer = iigoClients[JudgeIDGlobal].GetClientJudgePointer()
	// Set speakerPointer
	speakerPointer = iigoClients[SpeakerIDGlobal].GetClientSpeakerPointer()
	// Set presidentPointer
	presidentPointer = iigoClients[PresidentIDGlobal].GetClientPresidentPointer()

	// Initialise iigointernal with their clientVersions
	judicialBranch.loadClientJudge(judgePointer)
	executiveBranch.loadClientPresident(presidentPointer)
	legislativeBranch.loadClientSpeaker(speakerPointer)

	// Withdraw the salaries
	presidentWithdrawSuccess := judicialBranch.withdrawPresidentSalary(g)
	judgeWithdrawSuccess := legislativeBranch.withdrawJudgeSalary(g)
	speakerWithdrawSuccess := executiveBranch.withdrawSpeakerSalary(g)

	// Handle the lack of resources
	if !presidentWithdrawSuccess {
		returnWithdrawnSalariesToCommonPool(g, &executiveBranch, &legislativeBranch, &judicialBranch)
		return false, "Could not run IIGO since President has no resources to spend"
	}
	if !judgeWithdrawSuccess {
		returnWithdrawnSalariesToCommonPool(g, &executiveBranch, &legislativeBranch, &judicialBranch)
		return false, "Could not run IIGO since Judge has no resources to spend"
	}
	if !speakerWithdrawSuccess {
		returnWithdrawnSalariesToCommonPool(g, &executiveBranch, &legislativeBranch, &judicialBranch)
		return false, "Could not run IIGO since Speaker has no resources to spend"
	}

	// Pay salaries into budgets
	judicialBranch.sendPresidentSalary(&executiveBranch)
	legislativeBranch.sendJudgeSalary(&judicialBranch)
	executiveBranch.sendSpeakerSalary(&legislativeBranch)

	// 1 Judge actions - inspect history
	_, historyInspected := judicialBranch.inspectHistory(g.IIGOHistory)

	// 2 President actions
	resourceReports := map[shared.ClientID]shared.Resources{}

	// TODO get alive clients differently
	var aliveClientIds []shared.ClientID
	for _, v := range rules.VariableMap[rules.IslandsAlive].Values {
		aliveClientIds = append(aliveClientIds, shared.ClientID(int(v)))
		resourceReports[shared.ClientID(int(v))] = iigoClients[shared.ClientID(int(v))].ResourceReport()
	}
	executiveBranch.broadcastTaxation(resourceReports, aliveClientIds)
	executiveBranch.requestAllocationRequest(aliveClientIds)
	executiveBranch.replyAllocationRequest(g.CommonPool, aliveClientIds)
	executiveBranch.requestRuleProposal(aliveClientIds)
	ruleToVoteReturn := executiveBranch.getRuleForSpeaker()

	// 3 Speaker actions
	if ruleToVoteReturn.ActionTaken && ruleToVoteReturn.ContentType == shared.PresidentRuleProposal {
		legislativeBranch.setRuleToVote(ruleToVoteReturn.ProposedRule)
		legislativeBranch.setVotingResult(aliveClientIds)
		legislativeBranch.announceVotingResult()
	}

	// 4 Declare performance (Judge) (in future all the iigointernal)
	if historyInspected {
		judicialBranch.declarePresidentPerformanceWrapped()

		judicialBranch.declareSpeakerPerformanceWrapped()
	}

	// Get new Judge ID
	JudgeIDGlobal = legislativeBranch.appointNextJudge(aliveClientIds)
	// Get new Speaker ID
	SpeakerIDGlobal = executiveBranch.appointNextSpeaker(aliveClientIds)
	// Get new President ID
	PresidentIDGlobal = judicialBranch.appointNextPresident(aliveClientIds)

	return true, "IIGO Run Successful"
}

func returnWithdrawnSalariesToCommonPool(state *gamestate.GameState, executiveBranch *executive, legislativeBranch *legislature, judicialBranch *judiciary) {
	returnVal := executiveBranch.returnSpeakerSalary() + legislativeBranch.returnJudgeSalary() + judicialBranch.returnPresidentSalary()
	depositIntoCommonPool(returnVal, state)
}
