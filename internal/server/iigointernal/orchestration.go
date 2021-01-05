package iigointernal

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"github.com/SOMAS2020/SOMAS2020/internal/common/voting"
)

// TaxAmountMapExport is a local tax amount cache for checking of rules
var TaxAmountMapExport map[shared.ClientID]shared.Resources

// AllocationAmountMapExport is a local allocation map for checking of rules
var AllocationAmountMapExport map[shared.ClientID]shared.Resources

// SanctionAmountMapExport is a local sanction map for sanctions
var SanctionAmountMapExport map[shared.ClientID]shared.Resources

// iigoClients holds pointers to all the clients
var iigoClients map[shared.ClientID]baseclient.Client

// RunIIGO runs all iigo function in sequence
func RunIIGO(g *gamestate.GameState, clientMap *map[shared.ClientID]baseclient.Client, gameConf *config.Config) (IIGOSuccessful bool, StatusDescription string) {

	// Pointers allow clients to customise implementations of mutable functions
	var judgePointer roles.Judge = nil
	var speakerPointer roles.Speaker = nil
	var presidentPointer roles.President = nil

	// featureJudge is an instantiation of the Judge interface
	// with both the Base Judge features and a reference to client judges
	var judicialBranch = judiciary{
		gameState:          nil,
		gameConf:           nil,
		JudgeID:            0,
		evaluationResults:  nil,
		localSanctionCache: defaultInitLocalSanctionCache(3),
		localHistoryCache:  defaultInitLocalHistoryCache(3),
	}

	// featureSpeaker is an instantiation of the Speaker interface
	// with both the baseSpeaker features and a reference to client speakers
	var legislativeBranch = legislature{
		gameState:    nil,
		gameConf:     nil,
		SpeakerID:    0,
		ruleToVote:   "",
		ballotBox:    voting.BallotBox{},
		votingResult: false,
	}

	// featurePresident is an instantiation of the President interface
	// with both the basePresident features and a reference to client presidents
	var executiveBranch = executive{
		gameState:        nil,
		gameConf:         nil,
		PresidentID:      0,
		ResourceRequests: nil,
	}

	var monitoring = monitor{
		speakerID:         g.SpeakerID,
		presidentID:       g.PresidentID,
		judgeID:           g.JudgeID,
		internalIIGOCache: []shared.Accountability{},
	}
	executiveBranch.monitoring = &monitoring
	legislativeBranch.monitoring = &monitoring
	judicialBranch.monitoring = &monitoring

	iigoClients = *clientMap

	// Increments the budget by a constant 100
	// TODO:- the constant should be retrieved from the rules
	g.IIGORolesBudget[shared.President] += 100
	g.IIGORolesBudget[shared.Judge] += 100
	g.IIGORolesBudget[shared.Speaker] += 100

	// Pass in gamestate and IIGO configs
	// So that we don't have to pass gamestate as arguments in every function in roles
	judicialBranch.syncWithGame(g, &gameConf.IIGOConfig)
	legislativeBranch.syncWithGame(g, &gameConf.IIGOConfig)
	executiveBranch.syncWithGame(g, &gameConf.IIGOConfig)

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

	// 1 Judge action - inspect history
	judicialBranch.loadSanctionConfig()
	historyInspected := true
	if g.Turn > 0 {
		_, historyInspected = judicialBranch.inspectHistory(g.IIGOHistory[g.Turn-1])
		judicialBranch.updateSanctionScore()
		judicialBranch.applySanctions()
	}

	variablesToCache := []rules.VariableFieldName{rules.JudgeInspectionPerformed}
	valuesToCache := [][]float64{{boolToFloat(historyInspected)}}
	monitoring.addToCache(g.PresidentID, variablesToCache, valuesToCache)

	judgeMonitored := monitoring.monitorRole(iigoClients[g.PresidentID])

	// 2 President actions
	resourceReports := map[shared.ClientID]shared.ResourcesReport{}
	aliveClientIds := []shared.ClientID{}
	for clientID, clientGameState := range g.ClientInfos {
		if clientGameState.LifeStatus != shared.Dead {
			aliveClientIds = append(aliveClientIds, clientID)
			resourceReports[clientID] = iigoClients[clientID].ResourceReport()

			// Update Variables in Rules (updateIIGOTurnHistory)
			g.IIGOHistory[g.Turn] = append(g.IIGOHistory[g.Turn],
				shared.Accountability{
					ClientID: clientID,
					Pairs: []rules.VariableValuePair{
						{
							VariableName: rules.HasIslandReportPrivateResources,
							Values:       []float64{boolToFloat(resourceReports[clientID].Reported)},
						},
						{
							VariableName: rules.IslandReportedPrivateResources,
							Values:       []float64{float64(resourceReports[clientID].ReportedAmount)},
						},
						{
							VariableName: rules.IslandActualPrivateResources,
							Values:       []float64{float64(g.ClientInfos[clientID].Resources)},
						},
					},
				})
		}
	}

	// Judge uses resourceReports
	if g.Turn > 0 {
		judicialBranch.sanctionEvaluate(resourceReports)
	}

	// Throw error if any of the actions returns error
	insufficientBudget := executiveBranch.broadcastTaxation(resourceReports, aliveClientIds)
	if insufficientBudget != nil {
		return false, "Common pool resources insufficient for executiveBranch broadcastTaxation"
	}
	//var ruleToVoteReturn shared.PresidentReturnContent
	insufficientBudget = executiveBranch.requestAllocationRequest(aliveClientIds)
	if insufficientBudget != nil {
		return false, "Common pool resources insufficient for executiveBranch requestAllocationRequest"
	}

	allocationsMade, insufficientBudget := executiveBranch.replyAllocationRequest(g.CommonPool)
	if insufficientBudget != nil {
		return false, "Common pool resources insufficient for executiveBranch replyAllocationRequest"
	}

	insufficientBudget = executiveBranch.requestRuleProposal()
	if insufficientBudget != nil {
		return false, "Common pool resources insufficient for executiveBranch requestRuleProposal"
	}

	ruleToVoteReturn, insufficientBudget := executiveBranch.getRuleForSpeaker()
	if insufficientBudget != nil {
		return false, "Common pool resources insufficient for executiveBranch getRuleForSpeaker"
	}

	ruleSelected := false
	if ruleToVoteReturn.ActionTaken && ruleToVoteReturn.ProposedRule != "" {
		ruleSelected = true
	}

	variablesToCache = []rules.VariableFieldName{rules.AllocationMade}
	valuesToCache = [][]float64{{boolToFloat(allocationsMade)}}
	monitoring.addToCache(g.PresidentID, variablesToCache, valuesToCache)

	presidentMonitored := monitoring.monitorRole(iigoClients[g.SpeakerID])

	// 3 Speaker actions

	insufficientBudget = legislativeBranch.setRuleToVote(ruleToVoteReturn.ProposedRule)
	if insufficientBudget != nil {
		return false, "Common pool resources insufficient for legislativeBranch setRuleToVote"
	}
	voteCalled, insufficientBudget := legislativeBranch.setVotingResult(aliveClientIds)
	if insufficientBudget != nil {
		return false, "Common pool resources insufficient for legislativeBranch setVotingResult"
	}
	resultAnnounced, insufficientBudget := legislativeBranch.announceVotingResult()
	if insufficientBudget != nil {
		return false, "Common pool resources insufficient for legislativeBranch announceVotingResult"
	}

	variablesToCache = []rules.VariableFieldName{rules.RuleSelected, rules.VoteCalled}
	valuesToCache = [][]float64{{boolToFloat(ruleSelected)}, {boolToFloat(voteCalled)}}
	monitoring.addToCache(g.SpeakerID, variablesToCache, valuesToCache)

	variablesToCache = []rules.VariableFieldName{rules.VoteCalled, rules.VoteResultAnnounced}
	valuesToCache = [][]float64{{boolToFloat(voteCalled)}, {boolToFloat(resultAnnounced)}}
	monitoring.addToCache(g.SpeakerID, variablesToCache, valuesToCache)

	speakerMonitored := monitoring.monitorRole(iigoClients[g.JudgeID])

	// TODO:- at the moment, these are action (and cost resources) but should they?
	var appointJudgeError, appointSpeakerError, appointPresidentError error
	// Get new Judge ID
	actionCost := gameConf.IIGOConfig
	costOfElection := actionCost.AppointNextSpeakerActionCost + actionCost.AppointNextJudgeActionCost + actionCost.AppointNextPresidentActionCost
	if !CheckEnoughInCommonPool(costOfElection, g) {
		return false, "Insufficient budget to run IIGO elections"
	}
	g.JudgeID, appointJudgeError = legislativeBranch.appointNextJudge(judgeMonitored, g.JudgeID, aliveClientIds)
	if appointJudgeError != nil {
		return false, "Judge was not apointed by the Speaker. Insufficient budget"
	}
	// Get new Speaker ID
	g.SpeakerID, appointSpeakerError = executiveBranch.appointNextSpeaker(speakerMonitored, g.SpeakerID, aliveClientIds)
	if appointSpeakerError != nil {
		return false, "Speaker was not apointed by the President. Insufficient budget"
	}
	// Get new President ID
	g.PresidentID, appointPresidentError = judicialBranch.appointNextPresident(presidentMonitored, g.PresidentID, aliveClientIds)
	if appointPresidentError != nil {
		return false, "President was not apointed by the Judge. Insufficient budget"
	}

	// Pay salaries into budgets
	errorJudicial := judicialBranch.sendPresidentSalary()
	errorLegislative := legislativeBranch.sendJudgeSalary()
	errorExecutive := executiveBranch.sendSpeakerSalary()
	// Return false only after attempting to pay all roles their salary
	if errorJudicial != nil || errorLegislative != nil || errorExecutive != nil {
		return false, "Cannot pay IIGO salary"
	}

	return true, "IIGO Run Successful"
}
