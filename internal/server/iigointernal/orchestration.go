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
	JudgeID:            0,
	budget:             0,
	presidentSalary:    0,
	evaluationResults:  nil,
	localSanctionCache: defaultInitLocalSanctionCache(sanctionCacheDepth),
	localHistoryCache:  defaultInitLocalHistoryCache(historyCacheDepth),
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

// TaxAmountMapExport is a local tax amount cache for checking of rules
var TaxAmountMapExport map[shared.ClientID]shared.Resources

// AllocationAmountMapExport is a local allocation map for checking of rules
var AllocationAmountMapExport map[shared.ClientID]shared.Resources

// SanctionAmountMapExport
var SanctionAmountMapExport map[shared.ClientID]shared.Resources

// Pointers allow clients to customise implementations of mutable functions
var judgePointer roles.Judge = nil
var speakerPointer roles.Speaker = nil
var presidentPointer roles.President = nil

// iigoClients holds pointers to all the clients
var iigoClients map[shared.ClientID]baseclient.Client

// RunIIGO runs all iigo function in sequence
func RunIIGO(g *gamestate.GameState, clientMap *map[shared.ClientID]baseclient.Client) (IIGOSuccessful bool, StatusDescription string) {

	var monitoring = monitor{
		speakerID:         g.SpeakerID,
		presidentID:       g.PresidentID,
		judgeID:           g.JudgeID,
		internalIIGOCache: []shared.Accountability{},
	}
	iigoClients = *clientMap

	// Initialise IDs
	judicialBranch.JudgeID = g.JudgeID
	legislativeBranch.SpeakerID = g.SpeakerID
	executiveBranch.ID = g.PresidentID

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

	// 1 Judge action - inspect history
	_, historyInspected := judicialBranch.inspectHistory(g.IIGOHistory)

	variablesToCache := []rules.VariableFieldName{rules.JudgeInspectionPerformed}
	valuesToCache := [][]float64{{boolToFloat(historyInspected)}}
	monitoring.addToCache(g.PresidentID, variablesToCache, valuesToCache)

	judgeMonitored := monitoring.monitorRole(iigoClients[g.PresidentID])

	// 2 President actions
	resourceReports := map[shared.ClientID]shared.Resources{}
	var aliveClientIds []shared.ClientID
	for _, v := range rules.VariableMap[rules.IslandsAlive].Values {
		aliveClientIds = append(aliveClientIds, shared.ClientID(int(v)))
		resourceReports[shared.ClientID(int(v))] = iigoClients[shared.ClientID(int(v))].ResourceReport()
	}
	executiveBranch.broadcastTaxation(resourceReports)
	executiveBranch.requestAllocationRequest()
	allocationsMade := executiveBranch.replyAllocationRequest(g.CommonPool)
	executiveBranch.requestRuleProposal()
	ruleToVote, ruleSelected := executiveBranch.getRuleForSpeaker()

	variablesToCache = []rules.VariableFieldName{rules.AllocationMade}
	valuesToCache = [][]float64{{boolToFloat(allocationsMade)}}
	monitoring.addToCache(g.PresidentID, variablesToCache, valuesToCache)

	presidentMonitored := monitoring.monitorRole(iigoClients[g.SpeakerID])

	// 3 Speaker actions
	legislativeBranch.setRuleToVote(ruleToVote)
	voteCalled := legislativeBranch.setVotingResult(aliveClientIds)
	legislativeBranch.announceVotingResult()

	variablesToCache = []rules.VariableFieldName{rules.RuleSelected, rules.VoteCalled}
	valuesToCache = [][]float64{{boolToFloat(ruleSelected)}, {boolToFloat(voteCalled)}}
	monitoring.addToCache(g.SpeakerID, variablesToCache, valuesToCache)

	speakerMonitored := monitoring.monitorRole(iigoClients[g.JudgeID])

	// Get new Judge ID
	g.JudgeID = legislativeBranch.appointNextJudge(judgeMonitored, g.JudgeID, aliveClientIds)
	// Get new Speaker ID
	g.SpeakerID = executiveBranch.appointNextSpeaker(speakerMonitored, g.SpeakerID, aliveClientIds)
	// Get new President ID
	g.PresidentID = judicialBranch.appointNextPresident(presidentMonitored, g.PresidentID, aliveClientIds)

	return true, "IIGO Run Successful"
}

func returnWithdrawnSalariesToCommonPool(state *gamestate.GameState, executiveBranch *executive, legislativeBranch *legislature, judicialBranch *judiciary) {
	returnVal := executiveBranch.returnSpeakerSalary() + legislativeBranch.returnJudgeSalary() + judicialBranch.returnPresidentSalary()
	depositIntoCommonPool(returnVal, state)
}
