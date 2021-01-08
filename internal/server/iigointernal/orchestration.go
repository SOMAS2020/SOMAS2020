package iigointernal

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
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

//monitoring holds the monitoring object that is used across turns
var monitoring monitor

// RunIIGO runs all iigo function in sequence
func RunIIGO(logger shared.Logger, g *gamestate.GameState, clientMap *map[shared.ClientID]baseclient.Client, gameConf *config.Config) (IIGOSuccessful bool, StatusDescription string) {

	removeDeadBodiesFromOffice(g)

	logger("President %v, Speaker %v, Judge %v", g.PresidentID, g.SpeakerID, g.JudgeID)

	var judicialBranch = judiciary{
		gameState:          nil,
		gameConf:           nil,
		JudgeID:            0,
		evaluationResults:  nil,
		localSanctionCache: defaultInitLocalSanctionCache(3),
		localHistoryCache:  defaultInitLocalHistoryCache(3),
		logger:             logger,
	}

	var legislativeBranch = legislature{
		gameState:    nil,
		gameConf:     nil,
		SpeakerID:    0,
		ruleToVote:   rules.RuleMatrix{},
		ballotBox:    voting.BallotBox{},
		votingResult: false,
		logger:       logger,
	}

	var executiveBranch = executive{
		gameState:        nil,
		gameConf:         nil,
		PresidentID:      0,
		ResourceRequests: nil,
		logger:           logger,
	}

	monitoring.logger = logger

	executiveBranch.monitoring = &monitoring
	legislativeBranch.monitoring = &monitoring
	judicialBranch.monitoring = &monitoring

	iigoClients = *clientMap

	// Increments the budget according to increment_budget_role rules
	PresidentIncRule, ok := rules.RulesInPlay["increment_budget_president"]
	if ok {
		PresidentBudgetInc := PresidentIncRule.ApplicableMatrix.At(0, 1)
		g.IIGORolesBudget[shared.President] += shared.Resources(PresidentBudgetInc)
	}
	JudgeIncRule, ok := rules.RulesInPlay["increment_budget_judge"]
	if ok {
		JudgeBudgetInc := JudgeIncRule.ApplicableMatrix.At(0, 1)
		g.IIGORolesBudget[shared.Judge] += shared.Resources(JudgeBudgetInc)
	}
	SpeakerIncRule, ok := rules.RulesInPlay["increment_budget_speaker"]
	if ok {
		SpeakerBudgetInc := SpeakerIncRule.ApplicableMatrix.At(0, 1)
		g.IIGORolesBudget[shared.Speaker] += shared.Resources(SpeakerBudgetInc)
	}

	//Increment the turns in Power for each role
	g.IIGOTurnsInPower[shared.President]++
	g.IIGOTurnsInPower[shared.Speaker]++
	g.IIGOTurnsInPower[shared.Judge]++

	// Pass in gamestate and IIGO configs
	// So that we don't have to pass gamestate as arguments in every function in roles
	judicialBranch.syncWithGame(g, &gameConf.IIGOConfig)
	legislativeBranch.syncWithGame(g, &gameConf.IIGOConfig)
	executiveBranch.syncWithGame(g, &gameConf.IIGOConfig)

	// Initialise IDs
	judicialBranch.JudgeID = g.JudgeID
	legislativeBranch.SpeakerID = g.SpeakerID
	executiveBranch.PresidentID = g.PresidentID

	monitoring.judgeID = g.JudgeID
	monitoring.presidentID = g.PresidentID
	monitoring.speakerID = g.SpeakerID

	// Set judgePointer
	judgePointer := iigoClients[g.JudgeID].GetClientJudgePointer()
	// Set speakerPointer
	speakerPointer := iigoClients[g.SpeakerID].GetClientSpeakerPointer()
	// Set presidentPointer
	presidentPointer := iigoClients[g.PresidentID].GetClientPresidentPointer()

	// Initialise iigointernal with their clientVersions
	judicialBranch.loadClientJudge(judgePointer)
	executiveBranch.loadClientPresident(presidentPointer)
	legislativeBranch.loadClientSpeaker(speakerPointer)

	// 1 Judge action - inspect history
	judicialBranch.loadSanctionConfig()
	if g.Turn > 0 {
		//TODO: handle return types, quit IIGO if no moneyz
		judicialBranch.inspectHistory(g.IIGOHistory[g.Turn-1])
		judicialBranch.updateSanctionScore()
		judicialBranch.applySanctions()
	}

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
	if !ruleToVoteReturn.ProposedRuleMatrix.RuleMatrixIsEmpty() {
		ruleSelected = true
	}

	variablesToCache := []rules.VariableFieldName{rules.AllocationMade}
	valuesToCache := [][]float64{{boolToFloat(allocationsMade)}}
	monitoring.addToCache(g.PresidentID, variablesToCache, valuesToCache)

	// 3 Speaker actions

	//TODO:- shouldn't updateRules be called somewhere?
	insufficientBudget = legislativeBranch.setRuleToVote(ruleToVoteReturn.ProposedRuleMatrix)

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

	// Pay salaries into budgets
	errorJudicial := judicialBranch.sendPresidentSalary()
	errorLegislative := legislativeBranch.sendJudgeSalary()
	errorExecutive := executiveBranch.sendSpeakerSalary()
	// Return false only after attempting to pay all roles their salary
	if errorJudicial != nil || errorLegislative != nil || errorExecutive != nil {
		return false, "Cannot pay IIGO salary"
	}

	presidentMonitored := monitoring.monitorRole(g, iigoClients[g.SpeakerID])
	judgeMonitored := monitoring.monitorRole(g, iigoClients[g.PresidentID])
	speakerMonitored := monitoring.monitorRole(g, iigoClients[g.JudgeID])
	//Save to gameState and clear
	g.IIGOCache = monitoring.internalIIGOCache
	monitoring.clearCache()

	// TODO:- at the moment, these are action (and cost resources) but should they?
	// Get new Role ID

	var appointJudgeError, appointSpeakerError, appointPresidentError error
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

	oldCacheSave := monitoring.internalIIGOCache
	monitoring = monitor{
		speakerID:         g.SpeakerID,
		presidentID:       g.PresidentID,
		judgeID:           g.JudgeID,
		internalIIGOCache: oldCacheSave,
		TermLengths:       gameConf.IIGOConfig.IIGOTermLengths,
	}

	return true, "IIGO Run Successful"
}
