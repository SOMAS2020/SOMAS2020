package iigointernal

import (
	"fmt"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
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
	evaluationResults     map[shared.ClientID]shared.EvaluationReturn
	clientJudge           roles.Judge
	sanctionRecord        map[shared.ClientID]shared.IIGOSanctionsScore
	sanctionThresholds    map[shared.IIGOSanctionsTier]shared.IIGOSanctionsScore
	ruleViolationSeverity map[string]shared.IIGOSanctionsScore
	iigoClients           map[shared.ClientID]baseclient.Client
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
}

func (j *judiciary) broadcastSanctionConfig() {
	broadcastGeneric(j.iigoClients, j.JudgeID, createBroadcastsForSanctionThresholds(j.sanctionThresholds), *j.gameState)
	broadcastGeneric(j.iigoClients, j.JudgeID, createBroadcastsForRuleViolationPenalties(j.ruleViolationSeverity), *j.gameState)
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
func (j *judiciary) inspectHistory(iigoHistory []shared.Accountability) (map[shared.ClientID]shared.EvaluationReturn, bool) {
	if !CheckEnoughInCommonPool(j.gameConf.InspectHistoryActionCost, j.gameState) {
		return nil, false
	}
	finalResults := getBaseEvalResults(shared.TeamIDs)
	tempResults, actionTakenByClient := j.clientJudge.InspectHistory(iigoHistory, 0)

	if actionTakenByClient {
		for island, results := range tempResults {
			for index, eval := range results.Evaluations {
				if !eval {
					j.gameState.RulesBrokenByIslands[island] = append(j.gameState.RulesBrokenByIslands[island], results.Rules[index].RuleName)
				}
			}
		}
	}

	//Log rule: "Judge has the obligation to inspect history"
	variablesToCache := []rules.VariableFieldName{rules.JudgeInspectionPerformed}
	valuesToCache := [][]float64{{boolToFloat(actionTakenByClient)}}
	j.monitoring.addToCache(j.JudgeID, variablesToCache, valuesToCache)

	if actionTakenByClient {
		//Quit if taking resources goes wrong
		if !j.incurServiceCharge(j.gameConf.InspectHistoryActionCost) {
			return nil, false
		}
	}

	rulesInPlay := j.gameState.RulesInfo.CurrentRulesInPlay
	if !CheckEnoughInCommonPool(j.gameConf.HistoricalRetributionActionCost, j.gameState) {
		if actionTakenByClient {
			finalResults = mergeEvaluationReturn(tempResults, finalResults)
			entryForHistoryCache := cullCheckedRules(iigoHistory, finalResults, rulesInPlay, j.gameState.RulesInfo.VariableMap)
			j.cycleHistoryCache(entryForHistoryCache, int(j.gameConf.HistoryCacheDepth))
			j.evaluationResults = finalResults
			return j.evaluationResults, actionTakenByClient
		}
		return nil, false
	}

	//Quit early if CP does not have enough resources for historical Retribution

	//Perform historical checking
	decisionOfHistoricalRetribution := j.clientJudge.HistoricalRetributionEnabled()
	if decisionOfHistoricalRetribution {
		j.incurServiceCharge(j.gameConf.InspectHistoryActionCost)
		for turnsAgo, v := range j.gameState.IIGOHistoryCache {
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
		j.gameState.IIGOHistoryCache = DefaultInitLocalHistoryCache(int(j.gameConf.HistoryCacheDepth))
	}

	//Log rule: "Judge is not allowed to perform historical retribution"
	variablesToCache = []rules.VariableFieldName{rules.JudgeHistoricalRetributionPerformed}
	valuesToCache = [][]float64{{boolToFloat(decisionOfHistoricalRetribution)}}
	j.monitoring.addToCache(j.JudgeID, variablesToCache, valuesToCache)

	finalResults = mergeEvaluationReturn(tempResults, finalResults)
	entryForHistoryCache := cullCheckedRules(iigoHistory, finalResults, rulesInPlay, j.gameState.RulesInfo.VariableMap)
	j.cycleHistoryCache(entryForHistoryCache, int(j.gameConf.HistoryCacheDepth))
	j.evaluationResults = finalResults
	return j.evaluationResults, actionTakenByClient
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
	allIslandsCopy1 := copyClientList(allIslands)
	electionSettings := j.clientJudge.CallPresidentElection(monitoring, int(j.gameState.IIGOTurnsInPower[shared.President]), allIslandsCopy1)

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
		allIslandsCopy2 := copyClientList(allIslands)
		election.OpenBallot(electionSettings.IslandsToVote, allIslandsCopy2)
		election.Vote(j.iigoClients)
		j.gameState.IIGOTurnsInPower[shared.President] = 0
		electedPresident := election.CloseBallot(j.iigoClients)
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
	j.gameState.IIGOElection = append(j.gameState.IIGOElection, election.GetVotingInfo())
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
	j.sanctionRecord = map[shared.ClientID]shared.IIGOSanctionsScore{}
	islandTransgressions := map[shared.ClientID][]string{}
	for k, v := range j.evaluationResults {
		islandTransgressions[k] = unpackSingleIslandTransgressions(v)
	}
	j.scoreIslandTransgressions(islandTransgressions)
}

// scoreIslandTransgressions uses the client provided ruleViolationSeverity map to score island transgressions
func (j *judiciary) scoreIslandTransgressions(transgressions map[shared.ClientID][]string) {
	for islandID, rulesBroken := range transgressions {
		totalIslandTurnScore := shared.IIGOSanctionsScore(0)
		for _, ruleBroken := range rulesBroken {
			if score, ok := j.ruleViolationSeverity[ruleBroken]; ok {
				totalIslandTurnScore += score
			} else {
				totalIslandTurnScore += j.gameConf.DefaultSanctionScore
			}
			j.Logf("Rule: %v, broken by: %v", ruleBroken, islandID)
		}
		j.sanctionRecord[islandID] += totalIslandTurnScore
	}
}

// applySanctions uses RulesInPlay and it's versions of the sanction rules to work out how much to sanction an island
func (j *judiciary) applySanctions() {
	j.cycleSanctionCache(int(j.gameConf.SanctionCacheDepth))
	var currentSanctions []shared.Sanction
	for islandID, sanctionScore := range j.sanctionRecord {
		islandSanctionTier := getIslandSanctionTier(sanctionScore, j.sanctionThresholds)
		sanctionEntry := shared.Sanction{
			ClientID:     islandID,
			SanctionTier: islandSanctionTier,
			TurnsLeft:    int(j.gameConf.SanctionLength),
		}
		currentSanctions = append(currentSanctions, sanctionEntry)
		broadcastToAllIslands(j.iigoClients, j.JudgeID, createBroadcastForSanction(islandID, islandSanctionTier), *j.gameState)
	}
	j.gameState.IIGOSanctionCache[0] = currentSanctions
}

// sanctionEvaluate allows the clients to effectively pardon islands, levy and communicate sanctions
func (j *judiciary) sanctionEvaluate(reportedIslandResources map[shared.ClientID]shared.ResourcesReport) {
	pardons := j.clientJudge.GetPardonedIslands(j.gameState.IIGOSanctionCache)
	pardonsValid, newSanctionMap, communications := implementPardons(j.gameState.IIGOSanctionCache, pardons, shared.TeamIDs)
	if pardonsValid {
		broadcastPardonCommunications(j.iigoClients, j.JudgeID, communications, *j.gameState)
	}
	j.gameState.IIGOSanctionCache = newSanctionMap
	totalSanctionPerAgent := runEvaluationRulesOnSanctions(j.gameState.IIGOSanctionCache, reportedIslandResources, j.gameState.RulesInfo.CurrentRulesInPlay, j.gameConf.AssumedResourcesNoReport)
	j.gameState.IIGOSanctionMap = totalSanctionPerAgent
	for clientID, sanctionedResources := range totalSanctionPerAgent {
		communicateWithIslands(j.iigoClients, clientID, j.JudgeID, map[shared.CommunicationFieldName]shared.CommunicationContent{
			shared.SanctionAmount: {
				T:           shared.CommunicationInt,
				IntegerData: int(sanctionedResources),
			},
		})
	}
	j.gameState.IIGOSanctionCache = decrementSanctionTime(j.gameState.IIGOSanctionCache)
}

// cycleSanctionCache rolls the sanction cache one turn forward (effectively dropping any sanctions longer than the depth)
func (j *judiciary) cycleSanctionCache(sanctionCacheDepth int) {
	oldMap := j.gameState.IIGOSanctionCache
	delete(oldMap, sanctionCacheDepth-1)
	newMapReturn := DefaultInitLocalSanctionCache(sanctionCacheDepth)
	for i := 0; i < sanctionCacheDepth-1; i++ {
		newMapReturn[i+1] = oldMap[i]
	}
	newMapReturn[0] = []shared.Sanction{}
	j.gameState.IIGOSanctionCache = newMapReturn
}

// cycleHistoryCache rolls the history cache (for retributive justice) forward
func (j *judiciary) cycleHistoryCache(iigoHistory []shared.Accountability, historyCacheDepth int) {
	oldMap := j.gameState.IIGOHistoryCache
	delete(oldMap, historyCacheDepth-1)
	newMapReturn := DefaultInitLocalHistoryCache(historyCacheDepth)
	for i := 0; i < historyCacheDepth-1; i++ {
		newMapReturn[i+1] = oldMap[i]
	}
	newMapReturn[0] = iigoHistory
	j.gameState.IIGOHistoryCache = newMapReturn
}

// Helper functions //

func broadcastGeneric(clients map[shared.ClientID]baseclient.Client, judgeID shared.ClientID, itemsForbroadcast []map[shared.CommunicationFieldName]shared.CommunicationContent, state gamestate.GameState) {
	for _, item := range itemsForbroadcast {
		broadcastToAllIslands(clients, judgeID, item, state)
	}
}

func createBroadcastForSanction(clientID shared.ClientID, sanctionTier shared.IIGOSanctionsTier) map[shared.CommunicationFieldName]shared.CommunicationContent {
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

func createBroadcastsForSanctionThresholds(thresholds map[shared.IIGOSanctionsTier]shared.IIGOSanctionsScore) []map[shared.CommunicationFieldName]shared.CommunicationContent {
	var outputBroadcast []map[shared.CommunicationFieldName]shared.CommunicationContent
	for tier, score := range thresholds {
		outputBroadcast = append(outputBroadcast, map[shared.CommunicationFieldName]shared.CommunicationContent{
			shared.IIGOSanctionTier: {
				T:           shared.CommunicationInt,
				IntegerData: int(tier),
			},
			shared.RuleSanctionPenalty: {
				T:           shared.CommunicationInt,
				IntegerData: int(score),
			},
		})
	}
	return outputBroadcast
}

func createBroadcastsForRuleViolationPenalties(penalties map[string]shared.IIGOSanctionsScore) []map[shared.CommunicationFieldName]shared.CommunicationContent {
	var outputBroadcast []map[shared.CommunicationFieldName]shared.CommunicationContent
	for ruleName, score := range penalties {
		outputBroadcast = append(outputBroadcast, map[shared.CommunicationFieldName]shared.CommunicationContent{
			shared.RuleName: {
				T:        shared.CommunicationString,
				TextData: ruleName,
			},
			shared.RuleSanctionPenalty: {
				T:           shared.CommunicationInt,
				IntegerData: int(score),
			},
		})
	}
	return outputBroadcast
}

// runEvaluationRulesOnSanctions uses the custom sanction evaluator calculate how much each island should be paying in sanctions
func runEvaluationRulesOnSanctions(localSanctionCache map[int][]shared.Sanction, reportedIslandResources map[shared.ClientID]shared.ResourcesReport, rulesCache map[string]rules.RuleMatrix, maxNoReport shared.Resources) map[shared.ClientID]shared.Resources {
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

func broadcastPardonCommunications(clients map[shared.ClientID]baseclient.Client, judgeID shared.ClientID, communications map[shared.ClientID][]map[shared.CommunicationFieldName]shared.CommunicationContent, state gamestate.GameState) {
	for _, communicationList := range communications {
		for _, v := range communicationList {
			broadcastToAllIslands(clients, judgeID, v, state)
		}
	}
}

func decrementSanctionTime(sanctions map[int][]shared.Sanction) (updatedSanctions map[int][]shared.Sanction) {
	for k, v := range sanctions {
		for index, sanction := range v {
			sanctions[k][index].TurnsLeft = sanction.TurnsLeft - 1
		}
	}
	return sanctions
}

// getDefaultSanctionThresholds provides default thresholds for sanctions
func getDefaultSanctionThresholds() map[shared.IIGOSanctionsTier]shared.IIGOSanctionsScore {
	return map[shared.IIGOSanctionsTier]shared.IIGOSanctionsScore{
		shared.SanctionTier1: 1,
		shared.SanctionTier2: 5,
		shared.SanctionTier3: 10,
		shared.SanctionTier4: 20,
		shared.SanctionTier5: 30,
	}
}

// softMergeSanctionThresholds merges the default sanction thresholds with a (preferred) client version
func softMergeSanctionThresholds(clientSanctionMap map[shared.IIGOSanctionsTier]shared.IIGOSanctionsScore) map[shared.IIGOSanctionsTier]shared.IIGOSanctionsScore {
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

func checkMonotonicityOfSanctionThresholds(clientSanctionMap map[shared.IIGOSanctionsTier]shared.IIGOSanctionsScore) bool {
	sanctionsInOrder := []shared.IIGOSanctionsTier{
		shared.SanctionTier1,
		shared.SanctionTier2,
		shared.SanctionTier3,
		shared.SanctionTier4,
		shared.SanctionTier5,
	}
	hold := shared.IIGOSanctionsScore(0)
	for _, v := range sanctionsInOrder {
		if clientSanctionMap[v] < hold {
			return false
		}
		hold = clientSanctionMap[v]
	}
	return true
}

// unpackSingleIslandTransgressions fetches the rule names of any broken rules as per EvaluationReturns
func unpackSingleIslandTransgressions(evaluationReturn shared.EvaluationReturn) []string {
	transgressions := []string{}
	for i, v := range evaluationReturn.Rules {
		if !evaluationReturn.Evaluations[i] {
			transgressions = append(transgressions, v.RuleName)
		}
	}
	return transgressions
}

// getIslandSanctionTier if statement based evaluator for which sanction tier a particular score is in
func getIslandSanctionTier(islandScore shared.IIGOSanctionsScore, scoreMap map[shared.IIGOSanctionsTier]shared.IIGOSanctionsScore) shared.IIGOSanctionsTier {
	if islandScore < scoreMap[shared.SanctionTier1] {
		return shared.NoSanction
	} else if islandScore < scoreMap[shared.SanctionTier2] {
		return shared.SanctionTier1
	} else if islandScore < scoreMap[shared.SanctionTier3] {
		return shared.SanctionTier2
	} else if islandScore < scoreMap[shared.SanctionTier4] {
		return shared.SanctionTier3
	} else if islandScore < scoreMap[shared.SanctionTier5] {
		return shared.SanctionTier4
	} else {
		return shared.SanctionTier5
	}
}

// getTierSanctionMap basic mapping between sanction tier and rule that governs it
func getTierSanctionMap() map[shared.IIGOSanctionsTier]string {
	return map[shared.IIGOSanctionsTier]string{
		shared.SanctionTier1: "iigo_economic_sanction_1",
		shared.SanctionTier2: "iigo_economic_sanction_2",
		shared.SanctionTier3: "iigo_economic_sanction_3",
		shared.SanctionTier4: "iigo_economic_sanction_4",
		shared.SanctionTier5: "iigo_economic_sanction_5",
	}
}

func implementPardons(sanctionCache map[int][]shared.Sanction, pardons map[int][]bool, allTeamIds [len(shared.TeamIDs)]shared.ClientID) (bool, map[int][]shared.Sanction, map[shared.ClientID][]map[shared.CommunicationFieldName]shared.CommunicationContent) {
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

func generateEmptyCommunicationsMap(allTeamIds [len(shared.TeamIDs)]shared.ClientID) map[shared.ClientID][]map[shared.CommunicationFieldName]shared.CommunicationContent {
	commsMap := map[shared.ClientID][]map[shared.CommunicationFieldName]shared.CommunicationContent{}
	for _, clientID := range allTeamIds {
		commsMap[clientID] = []map[shared.CommunicationFieldName]shared.CommunicationContent{}
	}
	return commsMap
}

func knitSanctions(sanctionCache map[int][]shared.Sanction, loc int, newSanctions []shared.Sanction) map[int][]shared.Sanction {
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

func processSingleTimeStep(sanctions []shared.Sanction, pardons []bool, allTeamIds [6]shared.ClientID) (sanctionsAfterPardons []shared.Sanction, commsForPardons map[shared.ClientID][]map[shared.CommunicationFieldName]shared.CommunicationContent) {
	finalSanctions := []shared.Sanction{}
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

func generateCommunication(clientID shared.ClientID, pardonTier shared.IIGOSanctionsTier) map[shared.CommunicationFieldName]shared.CommunicationContent {
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

func validatePardons(sanctionCache map[int][]shared.Sanction, pardons map[int][]bool) bool {
	return checkKeys(sanctionCache, pardons) && checkSizes(sanctionCache, pardons)
}

func checkKeys(sanctionCache map[int][]shared.Sanction, pardons map[int][]bool) bool {
	for k := range sanctionCache {
		if _, ok := pardons[k]; !ok {
			return false
		}
	}
	return true
}

func checkSizes(sanctionCache map[int][]shared.Sanction, pardons map[int][]bool) bool {
	for k := range sanctionCache {
		if len(sanctionCache[k]) != len(pardons[k]) {
			return false
		}
	}
	return true
}

func getBaseEvalResults(teamIDs [6]shared.ClientID) map[shared.ClientID]shared.EvaluationReturn {
	baseResults := map[shared.ClientID]shared.EvaluationReturn{}
	for _, teamID := range teamIDs {
		baseResults[teamID] = shared.EvaluationReturn{
			Rules:       []rules.RuleMatrix{},
			Evaluations: []bool{},
		}
	}
	return baseResults
}

// cullCheckedRules removes any entries in the history that have been evaluated in evalResults (for historical retribution)
func cullCheckedRules(iigoHistory []shared.Accountability, evalResults map[shared.ClientID]shared.EvaluationReturn, rulesCache map[string]rules.RuleMatrix, variableCache map[rules.VariableFieldName]rules.VariableValuePair) []shared.Accountability {
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

func searchEvalReturnForRuleName(name string, evaluationReturn shared.EvaluationReturn) bool {
	rulesAffected := evaluationReturn.Rules
	for _, v := range rulesAffected {
		if v.RuleName == name {
			return true
		}
	}
	return false
}

// mergeEvaluationReturn takes two evaluation results and merges them to provide a unified set
func mergeEvaluationReturn(set1 map[shared.ClientID]shared.EvaluationReturn, set2 map[shared.ClientID]shared.EvaluationReturn) map[shared.ClientID]shared.EvaluationReturn {
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
			resolved := shared.EvaluationReturn{
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
