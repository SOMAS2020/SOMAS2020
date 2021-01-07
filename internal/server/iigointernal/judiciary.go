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

type judiciary struct {
	gameState             *gamestate.GameState
	gameConf              *config.IIGOConfig
	JudgeID               shared.ClientID
	evaluationResults     map[shared.ClientID]roles.EvaluationReturn
	clientJudge           roles.Judge
	sanctionRecord        map[shared.ClientID]roles.IIGOSanctionScore
	sanctionThresholds    map[roles.IIGOSanctionTier]roles.IIGOSanctionScore
	ruleViolationSeverity map[string]roles.IIGOSanctionScore
	localSanctionCache    map[int][]roles.Sanction
	localHistoryCache     map[int][]shared.Accountability
	monitoring            *monitor
	logger                shared.Logger
}

func (j *judiciary) Logf(format string, a ...interface{}) {
	j.logger("[JUDICIARY]: %v", fmt.Sprintf(format, a...))
}

// Loads ruleViolationSeverity and sanction thresholds
func (j *judiciary) loadSanctionConfig() {
	j.sanctionThresholds = softMergeSanctionThresholds(j.clientJudge.GetSanctionThresholds())
	j.ruleViolationSeverity = j.clientJudge.GetRuleViolationSeverity()
	j.broadcastSanctionConfig()
}

func (j *judiciary) syncWithGame(gameState *gamestate.GameState, gameConf *config.IIGOConfig) {
	j.gameState = gameState
	j.gameConf = gameConf
	j.resetCaches()
}

func (j *judiciary) resetCaches() {
	if len(j.localSanctionCache) != int(j.gameConf.SanctionCacheDepth) {
		j.localSanctionCache = defaultInitLocalSanctionCache(int(j.gameConf.SanctionCacheDepth))
	}
	if len(j.localHistoryCache) != int(j.gameConf.HistoryCacheDepth) {
		j.localHistoryCache = defaultInitLocalHistoryCache(int(j.gameConf.HistoryCacheDepth))
	}
}

func (j *judiciary) broadcastSanctionConfig() {
	broadcastGeneric(j.JudgeID, createBroadcastsForSanctionThresholds(j.sanctionThresholds))
	broadcastGeneric(j.JudgeID, createBroadcastsForRuleViolationPenalties(j.ruleViolationSeverity))
}

// loadClientJudge checks client pointer is good and if not panics
func (j *judiciary) loadClientJudge(clientJudgePointer roles.Judge) {
	if clientJudgePointer == nil {
		panic(fmt.Sprintf("Client '%v' has loaded a nil judge pointer", j.JudgeID))
	}
	j.clientJudge = clientJudgePointer
}

// sendPresidentSalary conduct the transaction based on amount from client implementation
func (j *judiciary) sendPresidentSalary() error {
	if j.clientJudge != nil {
		amountReturn, presidentPaid := j.clientJudge.PayPresident()
		if presidentPaid {
			// Subtract from common resources po
			amountWithdraw, withdrawSuccess := WithdrawFromCommonPool(amountReturn, j.gameState)

			if withdrawSuccess {
				// Pay into the client private resources pool
				depositIntoClientPrivatePool(amountWithdraw, j.gameState.PresidentID, j.gameState)

				variablesToCache := []rules.VariableFieldName{rules.PresidentPayment}
				valuesToCache := [][]float64{{float64(amountWithdraw)}}
				j.monitoring.addToCache(j.JudgeID, variablesToCache, valuesToCache)
				return nil
			}
		}
		variablesToCache := []rules.VariableFieldName{rules.PresidentPaid}
		valuesToCache := [][]float64{{boolToFloat(presidentPaid)}}
		j.monitoring.addToCache(j.JudgeID, variablesToCache, valuesToCache)
	}
	return errors.Errorf("Cannot perform sendJudgeSalary")
}

// inspectHistory checks all actions that happened in the last turn and audits them.
// This can be overridden by clients.
func (j *judiciary) inspectHistory(iigoHistory []shared.Accountability) (map[shared.ClientID]roles.EvaluationReturn, bool) {
	if !CheckEnoughInCommonPool(j.gameConf.InspectHistoryActionCost, j.gameState) {
		return nil, false
	}
	finalResults := getBaseEvalResults(shared.TeamIDs)
	tempResults, success := j.clientJudge.InspectHistory(iigoHistory, 0)

	//Log rule: "Judge has the obligation to inspect history"
	variablesToCache := []rules.VariableFieldName{rules.JudgeInspectionPerformed}
	valuesToCache := [][]float64{{boolToFloat(success)}}
	j.monitoring.addToCache(j.JudgeID, variablesToCache, valuesToCache)

	if success {
		//Quit if taking resources goes wrong
		if !j.incurServiceCharge(j.gameConf.InspectHistoryActionCost) {
			return nil, false
		}
	}

	//Quit early if CP does not have enough resources for historical Retribution
	if success && !CheckEnoughInCommonPool(j.gameConf.HistoricalRetributionActionCost, j.gameState) {
		finalResults = mergeEvaluationReturn(tempResults, finalResults)
		entryForHistoryCache := cullCheckedRules(iigoHistory, finalResults, rules.RulesInPlay, rules.VariableMap)
		j.cycleHistoryCache(entryForHistoryCache, int(j.gameConf.HistoryCacheDepth))
		j.evaluationResults = finalResults
		return j.evaluationResults, success
	}

	//Perform historical checking
	decisionOfHistoricalRetribution := j.clientJudge.HistoricalRetributionEnabled()
	if decisionOfHistoricalRetribution {
		if !j.incurServiceCharge(j.gameConf.InspectHistoryActionCost) {
			//Quit if taking resources goes wrong
			if success {
				finalResults = mergeEvaluationReturn(tempResults, finalResults)
				entryForHistoryCache := cullCheckedRules(iigoHistory, finalResults, rules.RulesInPlay, rules.VariableMap)
				j.cycleHistoryCache(entryForHistoryCache, int(j.gameConf.HistoryCacheDepth))
				j.evaluationResults = finalResults
				return j.evaluationResults, success
			}
			return nil, false
		}
		for turnsAgo, v := range j.localHistoryCache {
			res, rsuccess := j.clientJudge.InspectHistory(v, turnsAgo+1)
			if rsuccess {
				for key, accounts := range res {
					curr := finalResults[key]
					curr.Evaluations = append(curr.Evaluations, accounts.Evaluations...)
					curr.Rules = append(curr.Rules, accounts.Rules...)
					finalResults[key] = curr
				}
			}
		}
		j.localHistoryCache = defaultInitLocalHistoryCache(int(j.gameConf.HistoryCacheDepth))
	}

	//Log rule: "Judge is not allowed to perform historical retribution"
	variablesToCache = []rules.VariableFieldName{rules.JudgeHistoricalRetributionPerformed}
	valuesToCache = [][]float64{{boolToFloat(decisionOfHistoricalRetribution)}}
	j.monitoring.addToCache(j.JudgeID, variablesToCache, valuesToCache)

	finalResults = mergeEvaluationReturn(tempResults, finalResults)
	entryForHistoryCache := cullCheckedRules(iigoHistory, finalResults, rules.RulesInPlay, rules.VariableMap)
	j.cycleHistoryCache(entryForHistoryCache, int(j.gameConf.HistoryCacheDepth))
	j.evaluationResults = finalResults
	return j.evaluationResults, success
}

// searchForRule searches for a given rule in the RuleMatrix
func searchForRule(ruleName string, listOfRuleMatrices []rules.RuleMatrix) (int, bool) {
	for i, v := range listOfRuleMatrices {
		if v.RuleName == ruleName {
			return i, true
		}
	}
	return -1, false
}

// appointNextPresident returns the island ID of the island appointed to be President in the next turn
func (j *judiciary) appointNextPresident(monitoring shared.MonitorResult, currentPresident shared.ClientID, allIslands []shared.ClientID) (shared.ClientID, error) {
	var election = voting.Election{
		Logger: j.logger,
	}
	var appointedPresident shared.ClientID
	electionSettings := j.clientJudge.CallPresidentElection(monitoring, int(j.gameState.IIGOTurnsInPower[shared.President]), allIslands)

	//Log election rule
	termCondition := j.gameState.IIGOTurnsInPower[shared.President] > j.gameConf.IIGOTermLengths[shared.President]
	variablesToCache := []rules.VariableFieldName{rules.TermEnded, rules.ElectionHeld}
	valuesToCache := [][]float64{{boolToFloat(termCondition)}, {boolToFloat(electionSettings.HoldElection)}}
	j.monitoring.addToCache(j.JudgeID, variablesToCache, valuesToCache)

	if electionSettings.HoldElection {
		if !j.incurServiceCharge(j.gameConf.InspectHistoryActionCost) {
			return j.gameState.PresidentID, errors.Errorf("Insufficient Budget in common Pool: appointNextPresident")
		}
		election.ProposeElection(shared.President, electionSettings.VotingMethod)
		election.OpenBallot(electionSettings.IslandsToVote, allIslands)
		election.Vote(iigoClients)
		j.gameState.IIGOTurnsInPower[shared.President] = 0
		electedPresident := election.CloseBallot(iigoClients)
		appointedPresident = j.clientJudge.DecideNextPresident(electedPresident)

		//Log rule: Must appoint elected role
		appointmentMatchesVote := appointedPresident == electedPresident
		variablesToCache := []rules.VariableFieldName{rules.AppointmentMatchesVote}
		valuesToCache := [][]float64{{boolToFloat(appointmentMatchesVote)}}
		j.monitoring.addToCache(j.JudgeID, variablesToCache, valuesToCache)
		j.Logf("Result of election for new President: %v", appointedPresident)
	} else {
		appointedPresident = currentPresident
	}
	return appointedPresident, nil
}

func (j *judiciary) incurServiceCharge(cost shared.Resources) bool {
	_, ok := WithdrawFromCommonPool(cost, j.gameState)
	if ok {
		j.gameState.IIGORolesBudget[shared.Judge] -= cost
		if j.monitoring != nil {
			variablesToCache := []rules.VariableFieldName{rules.JudgeLeftoverBudget}
			valuesToCache := [][]float64{{float64(j.gameState.IIGORolesBudget[shared.Judge])}}
			j.monitoring.addToCache(j.JudgeID, variablesToCache, valuesToCache)
		}
	}
	return ok
}

// updateSanctionScore uses results of InspectHistory to assign sanction scores to clients
// This function relies upon (in part) the client provided RuleViolationSeverity
func (j *judiciary) updateSanctionScore() {
	//Clearing sanction map
	j.sanctionRecord = map[shared.ClientID]roles.IIGOSanctionScore{}
	islandTransgressions := map[shared.ClientID][]string{}
	for k, v := range j.evaluationResults {
		islandTransgressions[k] = unpackSingleIslandTransgressions(v)
	}
	j.scoreIslandTransgressions(islandTransgressions)
}

// scoreIslandTransgressions uses the client provided ruleViolationSeverity map to score island transgressions
func (j *judiciary) scoreIslandTransgressions(transgressions map[shared.ClientID][]string) {
	for islandID, rulesBroken := range transgressions {
		totalIslandTurnScore := roles.IIGOSanctionScore(0)
		for _, ruleBroken := range rulesBroken {
			if score, ok := j.ruleViolationSeverity[ruleBroken]; ok {
				totalIslandTurnScore += score
			} else {
				totalIslandTurnScore += roles.IIGOSanctionScore(1)
			}
			j.Logf("Rule: %v, broken by: %v", ruleBroken, islandID)
		}
		j.sanctionRecord[islandID] += totalIslandTurnScore
	}
}

// applySanctions uses RulesInPlay and it's versions of the sanction rules to work out how much to sanction an island
func (j *judiciary) applySanctions() {
	j.cycleSanctionCache(int(j.gameConf.SanctionCacheDepth))
	var currentSanctions []roles.Sanction
	for islandID, sanctionScore := range j.sanctionRecord {
		islandSanctionTier := getIslandSanctionTier(sanctionScore, j.sanctionThresholds)
		sanctionEntry := roles.Sanction{
			ClientID:     islandID,
			SanctionTier: islandSanctionTier,
			TurnsLeft:    int(j.gameConf.SanctionLength),
		}
		currentSanctions = append(currentSanctions, sanctionEntry)
		broadcastToAllIslands(j.JudgeID, createBroadcastForSanction(islandID, islandSanctionTier))
	}
	j.localSanctionCache[0] = currentSanctions
}

// sanctionEvaluate allows the clients to effectively pardon islands, levy and communicate sanctions
func (j *judiciary) sanctionEvaluate(reportedIslandResources map[shared.ClientID]shared.ResourcesReport) {
	pardons := j.clientJudge.GetPardonedIslands(j.localSanctionCache)
	pardonsValid, newSanctionMap, communications := implementPardons(j.localSanctionCache, pardons, shared.TeamIDs)
	if pardonsValid {
		broadcastPardonCommunications(j.JudgeID, communications)
	}
	j.localSanctionCache = newSanctionMap
	totalSanctionPerAgent := runEvaluationRulesOnSanctions(j.localSanctionCache, reportedIslandResources, rules.RulesInPlay, j.gameConf.AssumedResourcesNoReport)
	SanctionAmountMapExport = totalSanctionPerAgent
	for clientID, sanctionedResources := range totalSanctionPerAgent {
		communicateWithIslands(j.JudgeID, clientID, map[shared.CommunicationFieldName]shared.CommunicationContent{
			shared.SanctionAmount: {
				T:           shared.CommunicationInt,
				IntegerData: int(sanctionedResources),
			},
		})
	}
	j.localSanctionCache = decrementSanctionTime(j.localSanctionCache)
}

// cycleSanctionCache rolls the sanction cache one turn forward (effectively dropping any sanctions longer than the depth)
func (j *judiciary) cycleSanctionCache(sanctionCacheDepth int) {
	oldMap := j.localSanctionCache
	delete(oldMap, sanctionCacheDepth-1)
	newMapReturn := defaultInitLocalSanctionCache(sanctionCacheDepth)
	for i := 0; i < sanctionCacheDepth-1; i++ {
		newMapReturn[i+1] = oldMap[i]
	}
	newMapReturn[0] = []roles.Sanction{}
	j.localSanctionCache = newMapReturn
}

// cycleHistoryCache rolls the history cache (for retributive justice) forward
func (j *judiciary) cycleHistoryCache(iigoHistory []shared.Accountability, historyCacheDepth int) {
	oldMap := j.localHistoryCache
	delete(oldMap, historyCacheDepth-1)
	newMapReturn := defaultInitLocalHistoryCache(historyCacheDepth)
	for i := 0; i < historyCacheDepth-1; i++ {
		newMapReturn[i+1] = oldMap[i]
	}
	newMapReturn[0] = iigoHistory
	j.localHistoryCache = newMapReturn
}

// Helper functions //

func broadcastGeneric(judgeID shared.ClientID, itemsForbroadcast []map[shared.CommunicationFieldName]shared.CommunicationContent) {
	for _, item := range itemsForbroadcast {
		broadcastToAllIslands(judgeID, item)
	}
}

func createBroadcastForSanction(clientID shared.ClientID, sanctionTier roles.IIGOSanctionTier) map[shared.CommunicationFieldName]shared.CommunicationContent {
	return map[shared.CommunicationFieldName]shared.CommunicationContent{
		shared.SanctionClientID: {
			T:           shared.CommunicationInt,
			IntegerData: int(clientID),
		},
		shared.IIGOSanctionTier: {
			T:           shared.CommunicationInt,
			IntegerData: int(sanctionTier),
		},
	}
}

func createBroadcastsForSanctionThresholds(thresholds map[roles.IIGOSanctionTier]roles.IIGOSanctionScore) []map[shared.CommunicationFieldName]shared.CommunicationContent {
	var outputBroadcast []map[shared.CommunicationFieldName]shared.CommunicationContent
	for tier, score := range thresholds {
		outputBroadcast = append(outputBroadcast, map[shared.CommunicationFieldName]shared.CommunicationContent{
			shared.IIGOSanctionTier: {
				T:           shared.CommunicationInt,
				IntegerData: int(tier),
			},
			shared.IIGOSanctionScore: {
				T:           shared.CommunicationInt,
				IntegerData: int(score),
			},
		})
	}
	return outputBroadcast
}

func createBroadcastsForRuleViolationPenalties(penalties map[string]roles.IIGOSanctionScore) []map[shared.CommunicationFieldName]shared.CommunicationContent {
	var outputBroadcast []map[shared.CommunicationFieldName]shared.CommunicationContent
	for ruleName, score := range penalties {
		outputBroadcast = append(outputBroadcast, map[shared.CommunicationFieldName]shared.CommunicationContent{
			shared.RuleName: {
				T:        shared.CommunicationString,
				TextData: ruleName,
			},
			shared.IIGOSanctionScore: {
				T:           shared.CommunicationInt,
				IntegerData: int(score),
			},
		})
	}
	return outputBroadcast
}

// runEvaluationRulesOnSanctions uses the custom sanction evaluator calculate how much each island should be paying in sanctions
func runEvaluationRulesOnSanctions(localSanctionCache map[int][]roles.Sanction, reportedIslandResources map[shared.ClientID]shared.ResourcesReport, rulesCache map[string]rules.RuleMatrix, maxNoReport shared.Resources) map[shared.ClientID]shared.Resources {
	totalSanctionPerAgent := map[shared.ClientID]shared.Resources{}
	for _, sanctionList := range localSanctionCache {
		for _, sanction := range sanctionList {
			ruleName := getTierSanctionMap()[sanction.SanctionTier]
			if ruleMat, ok := rulesCache[ruleName]; ok {
				resources := maxNoReport
				if reportedIslandResources[sanction.ClientID].Reported {
					resources = reportedIslandResources[sanction.ClientID].ReportedAmount
				}
				sanctionVal := evaluateSanction(ruleMat, map[rules.VariableFieldName]rules.VariableValuePair{
					rules.IslandReportedResources: {
						VariableName: rules.IslandReportedResources,
						Values:       []float64{float64(resources)},
					},
					rules.ConstSanctionAmount: {
						VariableName: rules.ConstSanctionAmount,
						Values:       []float64{0},
					},
					rules.TurnsLeftOnSanction: {
						VariableName: rules.TurnsLeftOnSanction,
						Values:       []float64{float64(sanction.TurnsLeft)},
					},
				})
				totalSanctionPerAgent[sanction.ClientID] += sanctionVal
			} // TODO: When logger PR is available, pass through here and log the missing sanction
		}
	}
	return totalSanctionPerAgent
}

func broadcastPardonCommunications(judgeID shared.ClientID, communications map[shared.ClientID][]map[shared.CommunicationFieldName]shared.CommunicationContent) {
	for _, communicationList := range communications {
		for _, v := range communicationList {
			broadcastToAllIslands(judgeID, v)
		}
	}
}

func decrementSanctionTime(sanctions map[int][]roles.Sanction) (updatedSanctions map[int][]roles.Sanction) {
	for k, v := range sanctions {
		for index, sanction := range v {
			sanctions[k][index].TurnsLeft = sanction.TurnsLeft - 1
		}
	}
	return sanctions
}

// getDefaultSanctionThresholds provides default thresholds for sanctions
func getDefaultSanctionThresholds() map[roles.IIGOSanctionTier]roles.IIGOSanctionScore {
	return map[roles.IIGOSanctionTier]roles.IIGOSanctionScore{
		roles.SanctionTier1: 1,
		roles.SanctionTier2: 5,
		roles.SanctionTier3: 10,
		roles.SanctionTier4: 20,
		roles.SanctionTier5: 30,
	}
}

// softMergeSanctionThresholds merges the default sanction thresholds with a (preferred) client version
func softMergeSanctionThresholds(clientSanctionMap map[roles.IIGOSanctionTier]roles.IIGOSanctionScore) map[roles.IIGOSanctionTier]roles.IIGOSanctionScore {
	outputMap := getDefaultSanctionThresholds()
	for k := range outputMap {
		if clientVal, ok := clientSanctionMap[k]; ok {
			outputMap[k] = clientVal
		}
	}
	if checkMonotonicityOfSanctionThresholds(outputMap) {
		return outputMap
	} else {
		return getDefaultSanctionThresholds()
	}
}

func checkMonotonicityOfSanctionThresholds(clientSanctionMap map[roles.IIGOSanctionTier]roles.IIGOSanctionScore) bool {
	sanctionsInOrder := []roles.IIGOSanctionTier{
		roles.SanctionTier1,
		roles.SanctionTier2,
		roles.SanctionTier3,
		roles.SanctionTier4,
		roles.SanctionTier5,
	}
	hold := roles.IIGOSanctionScore(0)
	for _, v := range sanctionsInOrder {
		if clientSanctionMap[v] < hold {
			return false
		}
		hold = clientSanctionMap[v]
	}
	return true
}

// unpackSingleIslandTransgressions fetches the rule names of any broken rules as per EvaluationReturns
func unpackSingleIslandTransgressions(evaluationReturn roles.EvaluationReturn) []string {
	transgressions := []string{}
	for i, v := range evaluationReturn.Rules {
		if !evaluationReturn.Evaluations[i] {
			transgressions = append(transgressions, v.RuleName)
		}
	}
	return transgressions
}

// getIslandSanctionTier if statement based evaluator for which sanction tier a particular score is in
func getIslandSanctionTier(islandScore roles.IIGOSanctionScore, scoreMap map[roles.IIGOSanctionTier]roles.IIGOSanctionScore) roles.IIGOSanctionTier {
	if islandScore < scoreMap[roles.SanctionTier1] {
		return roles.NoSanction
	} else if islandScore < scoreMap[roles.SanctionTier2] {
		return roles.SanctionTier1
	} else if islandScore < scoreMap[roles.SanctionTier3] {
		return roles.SanctionTier2
	} else if islandScore < scoreMap[roles.SanctionTier4] {
		return roles.SanctionTier3
	} else if islandScore < scoreMap[roles.SanctionTier5] {
		return roles.SanctionTier4
	} else {
		return roles.SanctionTier5
	}
}

// getTierSanctionMap basic mapping between sanction tier and rule that governs it
func getTierSanctionMap() map[roles.IIGOSanctionTier]string {
	return map[roles.IIGOSanctionTier]string{
		roles.SanctionTier1: "iigo_economic_sanction_1",
		roles.SanctionTier2: "iigo_economic_sanction_2",
		roles.SanctionTier3: "iigo_economic_sanction_3",
		roles.SanctionTier4: "iigo_economic_sanction_4",
		roles.SanctionTier5: "iigo_economic_sanction_5",
	}
}

// defaultInitLocalSanctionCache generates a blank sanction cache
func defaultInitLocalSanctionCache(depth int) map[int][]roles.Sanction {
	returnMap := map[int][]roles.Sanction{}
	for i := 0; i < depth; i++ {
		returnMap[i] = []roles.Sanction{}
	}
	return returnMap
}

func implementPardons(sanctionCache map[int][]roles.Sanction, pardons map[int][]bool, allTeamIds [len(shared.TeamIDs)]shared.ClientID) (bool, map[int][]roles.Sanction, map[shared.ClientID][]map[shared.CommunicationFieldName]shared.CommunicationContent) {
	if validatePardons(sanctionCache, pardons) {
		finalSanctionCache := sanctionCache
		communicationsAboutPardons := generateEmptyCommunicationsMap(allTeamIds)
		for timeStep := range sanctionCache {
			sanctionsAfterPardons, communicationsForTimeStep := processSingleTimeStep(sanctionCache[timeStep], pardons[timeStep], allTeamIds)
			finalSanctionCache = knitSanctions(finalSanctionCache, timeStep, sanctionsAfterPardons)
			communicationsAboutPardons = knitPardonCommunications(communicationsAboutPardons, communicationsForTimeStep)
		}
		return true, finalSanctionCache, communicationsAboutPardons
	}
	return false, sanctionCache, nil
}

// defaultInitLocalHistoryCache generates a blank history cache
func defaultInitLocalHistoryCache(depth int) map[int][]shared.Accountability {
	returnMap := map[int][]shared.Accountability{}
	for i := 0; i < depth; i++ {
		returnMap[i] = []shared.Accountability{}
	}
	return returnMap
}

func generateEmptyCommunicationsMap(allTeamIds [len(shared.TeamIDs)]shared.ClientID) map[shared.ClientID][]map[shared.CommunicationFieldName]shared.CommunicationContent {
	commsMap := map[shared.ClientID][]map[shared.CommunicationFieldName]shared.CommunicationContent{}
	for _, clientID := range allTeamIds {
		commsMap[clientID] = []map[shared.CommunicationFieldName]shared.CommunicationContent{}
	}
	return commsMap
}

func knitSanctions(sanctionCache map[int][]roles.Sanction, loc int, newSanctions []roles.Sanction) map[int][]roles.Sanction {
	sanctionCache[loc] = newSanctions
	return sanctionCache
}

func knitPardonCommunications(originalCommunications map[shared.ClientID][]map[shared.CommunicationFieldName]shared.CommunicationContent, newCommunications map[shared.ClientID][]map[shared.CommunicationFieldName]shared.CommunicationContent) map[shared.ClientID][]map[shared.CommunicationFieldName]shared.CommunicationContent {
	for clientID, oldComms := range originalCommunications {
		if newComms, ok := newCommunications[clientID]; ok {
			originalCommunications[clientID] = append(oldComms, newComms...)
		}
	}
	return originalCommunications
}

func processSingleTimeStep(sanctions []roles.Sanction, pardons []bool, allTeamIds [6]shared.ClientID) (sanctionsAfterPardons []roles.Sanction, commsForPardons map[shared.ClientID][]map[shared.CommunicationFieldName]shared.CommunicationContent) {
	finalSanctions := []roles.Sanction{}
	finalComms := generateEmptyCommunicationsMap(allTeamIds)
	for entry, pardoned := range pardons {
		if pardoned {
			finalComms[sanctions[entry].ClientID] = append(finalComms[sanctions[entry].ClientID], generateCommunication(sanctions[entry].ClientID, sanctions[entry].SanctionTier))
		} else {
			finalSanctions = append(finalSanctions, sanctions[entry])
		}
	}
	return finalSanctions, finalComms
}

func generateCommunication(clientID shared.ClientID, pardonTier roles.IIGOSanctionTier) map[shared.CommunicationFieldName]shared.CommunicationContent {
	return map[shared.CommunicationFieldName]shared.CommunicationContent{
		shared.PardonClientID: {
			T:           shared.CommunicationInt,
			IntegerData: int(clientID),
		},
		shared.PardonTier: {
			T:           shared.CommunicationInt,
			IntegerData: int(pardonTier),
		},
	}
}

func validatePardons(sanctionCache map[int][]roles.Sanction, pardons map[int][]bool) bool {
	return checkKeys(sanctionCache, pardons) && checkSizes(sanctionCache, pardons)
}

func checkKeys(sanctionCache map[int][]roles.Sanction, pardons map[int][]bool) bool {
	for k := range sanctionCache {
		if _, ok := pardons[k]; !ok {
			return false
		}
	}
	return true
}

func checkSizes(sanctionCache map[int][]roles.Sanction, pardons map[int][]bool) bool {
	for k := range sanctionCache {
		if len(sanctionCache[k]) != len(pardons[k]) {
			return false
		}
	}
	return true
}

func getBaseEvalResults(teamIDs [6]shared.ClientID) map[shared.ClientID]roles.EvaluationReturn {
	baseResults := map[shared.ClientID]roles.EvaluationReturn{}
	for _, teamID := range teamIDs {
		baseResults[teamID] = roles.EvaluationReturn{
			Rules:       []rules.RuleMatrix{},
			Evaluations: []bool{},
		}
	}
	return baseResults
}

// cullCheckedRules removes any entries in the history that have been evaluated in evalResults (for historical retribution)
func cullCheckedRules(iigoHistory []shared.Accountability, evalResults map[shared.ClientID]roles.EvaluationReturn, rulesCache map[string]rules.RuleMatrix, variableCache map[rules.VariableFieldName]rules.VariableValuePair) []shared.Accountability {
	reducedAccountability := []shared.Accountability{}
	for _, v := range iigoHistory {
		pairsAffected := v.Pairs
		var allRulesAffected []string
		for _, pair := range pairsAffected {
			additionalRules, success := rules.PickUpRulesByVariable(pair.VariableName, rulesCache, variableCache)
			if success {
				allRulesAffected = append(allRulesAffected, additionalRules...)
			}
		}
		allRulesAffected = streamlineRulesAffected(allRulesAffected)
		for _, ruleAff := range allRulesAffected {
			found := searchEvalReturnForRuleName(ruleAff, evalResults[v.ClientID])
			if !found {
				reducedAccountability = append(reducedAccountability, v)
			}
		}
	}
	return reducedAccountability
}

// streamlineRulesAffected removes duplicate rules
func streamlineRulesAffected(input []string) []string {
	streamlineMap := map[string]bool{}
	for _, v := range input {
		streamlineMap[v] = true
	}
	var returnArray []string
	for key := range streamlineMap {
		returnArray = append(returnArray, key)
	}
	return returnArray
}

func searchEvalReturnForRuleName(name string, evaluationReturn roles.EvaluationReturn) bool {
	rulesAffected := evaluationReturn.Rules
	for _, v := range rulesAffected {
		if v.RuleName == name {
			return true
		}
	}
	return false
}

// mergeEvaluationReturn takes two evaluation results and merges them to provide a unified set
func mergeEvaluationReturn(set1 map[shared.ClientID]roles.EvaluationReturn, set2 map[shared.ClientID]roles.EvaluationReturn) map[shared.ClientID]roles.EvaluationReturn {
	for key, val := range set1 {
		if set2Val, ok := set2[key]; ok {
			defRules := set2Val.Rules
			defEvals := set2Val.Evaluations
			var finalRules = make([]rules.RuleMatrix, len(defRules))
			copy(finalRules, defRules)
			var finalEvals = make([]bool, len(defEvals))
			copy(finalEvals, defEvals)
			finalRules = append(finalRules, val.Rules...)
			finalEvals = append(finalEvals, val.Evaluations...)
			resolved := roles.EvaluationReturn{
				Rules:       finalRules,
				Evaluations: finalEvals,
			}
			set2[key] = resolved
		} else {
			set2[key] = set1[key]
		}
	}
	return set2
}
