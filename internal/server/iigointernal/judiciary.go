package iigointernal

import (
	"fmt"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"github.com/SOMAS2020/SOMAS2020/internal/common/voting"
)

// to be moved to paramters
const sanctionCacheDepth = 3
const historyCacheDepth = 3

// to be changed
const sanctionLength = 2

type judiciary struct {
	JudgeID               shared.ClientID
	budget                shared.Resources
	presidentSalary       shared.Resources
	evaluationResults     map[shared.ClientID]roles.EvaluationReturn
	clientJudge           roles.Judge
	presidentTurnsInPower int
	sanctionRecord        map[shared.ClientID]roles.IIGOSanctionScore
	sanctionThresholds    map[roles.IIGOSanctionTier]roles.IIGOSanctionScore
	ruleViolationSeverity map[string]roles.IIGOSanctionScore
	localSanctionCache    map[int][]roles.Sanction
	localHistoryCache     map[int][]shared.Accountability
}

// Loads ruleViolationSeverity and sanction thresholds
func (j *judiciary) loadSanctionConfig() {
	j.sanctionThresholds = softMergeSanctionThresholds(j.clientJudge.GetSanctionThresholds())
	j.ruleViolationSeverity = j.clientJudge.GetRuleViolationSeverity()
}

// loadClientJudge checks client pointer is good and if not panics
func (j *judiciary) loadClientJudge(clientJudgePointer roles.Judge) {
	if clientJudgePointer == nil {
		panic(fmt.Sprintf("Client '%v' has loaded a nil judge pointer", j.JudgeID))
	}
	j.clientJudge = clientJudgePointer
}

// returnPresidentSalary returns the salary to the common pool.
func (j *judiciary) returnPresidentSalary() shared.Resources {
	x := j.presidentSalary
	j.presidentSalary = 0
	return x
}

// withdrawPresidentSalary withdraws the president's salary from the common pool.
func (j *judiciary) withdrawPresidentSalary(gameState *gamestate.GameState) bool {
	var presidentSalary = shared.Resources(rules.VariableMap[rules.PresidentSalary].Values[0])
	var withdrawAmount, withdrawSuccesful = WithdrawFromCommonPool(presidentSalary, gameState)
	j.presidentSalary = withdrawAmount
	return withdrawSuccesful
}

// sendPresidentSalary sends the president's salary to the president.
func (j *judiciary) sendPresidentSalary(executiveBranch *executive) {
	if j.clientJudge != nil {
		amount, payPresident := j.clientJudge.PayPresident(j.presidentSalary)
		if payPresident {
			executiveBranch.budget = amount
		}
		return
	}
	amount := j.PayPresident()
	executiveBranch.budget = amount
}

// PayPresident pays the president salary.
func (j *judiciary) PayPresident() shared.Resources {
	hold := j.presidentSalary
	j.presidentSalary = 0
	return hold
}

// inspectHistory checks all actions that happened in the last turn and audits them.
// This can be overridden by clients.
func (j *judiciary) inspectHistory(iigoHistory []shared.Accountability) (map[shared.ClientID]roles.EvaluationReturn, bool) {
	j.budget -= serviceCharge
	success := false
	tempResults := map[shared.ClientID]roles.EvaluationReturn{}
	finalResults := getBaseEvalResults()
	if j.clientJudge.HistoricalRetributionEnabled() {
		for _, v := range j.localHistoryCache {
			res, rsuccess := j.clientJudge.InspectHistory(v)
			if !rsuccess {
				success = false
			} else {
				for key, accounts := range res {
					curr := finalResults[key]
					curr.Evaluations = append(curr.Evaluations, accounts.Evaluations...)
					curr.Rules = append(curr.Rules, accounts.Rules...)
					finalResults[key] = curr
				}
			}
		}
		j.localHistoryCache = defaultInitLocalHistoryCache(historyCacheDepth)
	}
	tempResults, success = j.clientJudge.InspectHistory(iigoHistory)
	finalResults = mergeEvalResults(tempResults, finalResults)
	entryForHistoryCache := cullCheckedRules(iigoHistory, finalResults, rules.RulesInPlay, rules.VariableMap)
	j.cycleHistoryCache(entryForHistoryCache)
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
func (j *judiciary) appointNextPresident(currentPresident shared.ClientID, allIslands []shared.ClientID) shared.ClientID {
	var election voting.Election
	var nextPresident shared.ClientID
	electionsettings := j.clientJudge.CallPresidentElection(j.presidentTurnsInPower, allIslands)
	if electionsettings.HoldElection {
		// TODO: deduct the cost of holding an election
		election.ProposeElection(baseclient.President, electionsettings.VotingMethod)
		election.OpenBallot(electionsettings.IslandsToVote)
		election.Vote(iigoClients)
		j.presidentTurnsInPower = 0
		nextPresident = election.CloseBallot()
		nextPresident = j.clientJudge.DecideNextPresident(nextPresident)
	} else {
		j.presidentTurnsInPower++
		nextPresident = currentPresident
	}
	return nextPresident
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
	for islandId, rulesBroken := range transgressions {
		totalIslandTurnScore := roles.IIGOSanctionScore(0)
		for _, ruleBroken := range rulesBroken {
			if score, ok := j.ruleViolationSeverity[ruleBroken]; ok {
				totalIslandTurnScore += score
			} else {
				totalIslandTurnScore += roles.IIGOSanctionScore(1)
			}
		}
		j.sanctionRecord[islandId] += totalIslandTurnScore
	}
}

// applySanctions uses RulesInPlay and it's versions of the sanction rules to work out how much to sanction an island
func (j *judiciary) applySanctions() {
	j.cycleSanctionCache()
	var currentSanctions []roles.Sanction
	for islandID, sanctionScore := range j.sanctionRecord {
		islandSanctionTier := getIslandSanctionTier(sanctionScore, j.sanctionThresholds)
		sanctionEntry := roles.Sanction{
			ClientID:     islandID,
			SanctionTier: islandSanctionTier,
			TurnsLeft:    sanctionLength,
		}
		currentSanctions = append(currentSanctions, sanctionEntry)
	}
	j.localSanctionCache[0] = currentSanctions
}

// sanctionEvaluate allows the clients to effectively pardon islands, levy and communicate sanctions
func (j *judiciary) sanctionEvaluate(reportedIslandResources map[shared.ClientID]shared.Resources) {
	pardons := j.clientJudge.GetPardonedIslands(j.localSanctionCache)
	pardonsValid, communications, newSanctionMap := checkPardons(j.localSanctionCache, pardons)
	if pardonsValid {
		for _, communicationList := range communications {
			for _, v := range communicationList {
				broadcastToAllIslands(j.JudgeID, v)
			}
		}
	}
	j.localSanctionCache = newSanctionMap
	totalSanctionPerAgent := map[shared.ClientID]shared.Resources{}
	for _, sanctionList := range j.localSanctionCache {
		for _, sanction := range sanctionList {
			stringName := getTierSanctionMap()[sanction.SanctionTier]
			ruleMat := rules.RulesInPlay[stringName]
			sanctionVal := evaluateSanction(ruleMat, map[rules.VariableFieldName]rules.VariableValuePair{
				rules.IslandReportedResources: {
					VariableName: rules.IslandReportedResources,
					Values:       []float64{float64(reportedIslandResources[sanction.ClientID])},
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
		}
	}
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

// cycleSanctionCache rolls the sanction cahce one turn forward (effectively dropping any sanctions longer than the depth)
func (j *judiciary) cycleSanctionCache() {
	newMap := j.localSanctionCache
	delete(newMap, sanctionCacheDepth-1)
	newMapCache := newMap
	for i := 0; i < sanctionCacheDepth-1; i++ {
		newMap[i+1] = newMapCache[i]
	}
	newMap[0] = []roles.Sanction{}
	j.localSanctionCache = newMap
}

// cycleHistoryCache rolls the history cache (for retributive justice) forward
func (j *judiciary) cycleHistoryCache(iigoHistory []shared.Accountability) {
	newMap := j.localHistoryCache
	delete(newMap, historyCacheDepth-1)
	newMapCache := newMap
	for i := 0; i < historyCacheDepth-1; i++ {
		newMap[i+1] = newMapCache[i]
	}
	newMap[0] = iigoHistory
	j.localHistoryCache = newMap
}

// clearHistoryCache wipes the history cache (when retributive justice has happened)
func (j *judiciary) clearHistoryCache() {
	j.localHistoryCache = defaultInitLocalHistoryCache(historyCacheDepth)
}

// Helper functions //

func decrementSanctionTime(sanctions map[int][]roles.Sanction) (updatedSanctions map[int][]roles.Sanction) {
	for _, v := range sanctions {
		for _, sanction := range v {
			sanction.TurnsLeft -= 1
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
	defaultMap := getDefaultSanctionThresholds()
	for k := range defaultMap {
		if clientVal, ok := clientSanctionMap[k]; ok {
			defaultMap[k] = clientVal
		}
	}
	return defaultMap
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
	if islandScore <= scoreMap[roles.SanctionTier1] {
		return roles.NoSanction
	} else if islandScore <= scoreMap[roles.SanctionTier2] {
		return roles.SanctionTier1
	} else if islandScore <= scoreMap[roles.SanctionTier3] {
		return roles.SanctionTier2
	} else if islandScore <= scoreMap[roles.SanctionTier4] {
		return roles.SanctionTier3
	} else if islandScore <= scoreMap[roles.SanctionTier5] {
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

func implementPardons(sanctionCache map[int][]roles.Sanction, pardons map[int][]bool, allTeamIds []shared.ClientID) (bool, map[int][]roles.Sanction, map[shared.ClientID][]map[shared.CommunicationFieldName]shared.CommunicationContent) {
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

func generateEmptyCommunicationsMap(allTeamIds []shared.ClientID) map[shared.ClientID][]map[shared.CommunicationFieldName]shared.CommunicationContent {
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

func processSingleTimeStep(sanctions []roles.Sanction, pardons []bool, allTeamIds []shared.ClientID) (sanctionsAfterPardons []roles.Sanction, commsForPardons map[shared.ClientID][]map[shared.CommunicationFieldName]shared.CommunicationContent) {
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

// removeSanctions is a helper function to remove a sanction element from a slice
func removeSanctions(slice []roles.Sanction, s int) []roles.Sanction {
	var output []roles.Sanction
	for i, v := range slice {
		if i != s {
			output = append(output, v)
		}
	}
	return output
}

// getDifferenceInLength helper function to get difference in length between two lists
func getDifferenceInLength(slice1 []roles.Sanction, slice2 []roles.Sanction) int {
	return len(slice1) - len(slice2)
}

// pickUpRulesByVariable returns a list of rule_id's which are affected by certain variables.
func pickUpRulesByVariable(variableName rules.VariableFieldName, ruleStore map[string]rules.RuleMatrix, variableMap map[rules.VariableFieldName]rules.VariableValuePair) ([]string, bool) {
	var Rules []string
	if _, ok := variableMap[variableName]; ok {
		for k, v := range ruleStore {
			_, found := searchForVariableInArray(variableName, v.RequiredVariables)
			if found {
				Rules = append(Rules, k)
			}
		}
		return Rules, true
	} else {
		// fmt.Sprintf("Variable name '%v' was not found in the variable cache", variableName)
		return []string{}, false
	}
}

func searchForVariableInArray(val rules.VariableFieldName, array []rules.VariableFieldName) (int, bool) {
	for i, v := range array {
		if v == val {
			return i, true
		}
	}
	return -1, false
}

func getBaseEvalResults() map[shared.ClientID]roles.EvaluationReturn {
	return map[shared.ClientID]roles.EvaluationReturn{
		shared.Team1: {
			Rules:       []rules.RuleMatrix{},
			Evaluations: []bool{},
		},
		shared.Team2: {
			Rules:       []rules.RuleMatrix{},
			Evaluations: []bool{},
		},
		shared.Team3: {
			Rules:       []rules.RuleMatrix{},
			Evaluations: []bool{},
		},
		shared.Team4: {
			Rules:       []rules.RuleMatrix{},
			Evaluations: []bool{},
		},
		shared.Team5: {
			Rules:       []rules.RuleMatrix{},
			Evaluations: []bool{},
		},
		shared.Team6: {
			Rules:       []rules.RuleMatrix{},
			Evaluations: []bool{},
		},
	}
}

// cullCheckedRules removes any entries in the history that have been evaluated in evalResults (for historical retribution)
func cullCheckedRules(iigoHistory []shared.Accountability, evalResults map[shared.ClientID]roles.EvaluationReturn, rulesCache map[string]rules.RuleMatrix, variableCache map[rules.VariableFieldName]rules.VariableValuePair) []shared.Accountability {
	reducedAccountability := []shared.Accountability{}
	for _, v := range iigoHistory {
		pairsAffected := v.Pairs
		var allRulesAffected []string
		for _, pair := range pairsAffected {
			additionalRules, success := pickUpRulesByVariable(pair.VariableName, rulesCache, variableCache)
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

// mergeEvalResults takes two evaluation results and merges them to provide a unified set
func mergeEvalResults(set1 map[shared.ClientID]roles.EvaluationReturn, set2 map[shared.ClientID]roles.EvaluationReturn) map[shared.ClientID]roles.EvaluationReturn {
	for key, val := range set1 {
		if set2Val, ok := set2[key]; ok {
			defRules := set2Val.Rules
			defEvals := set2Val.Evaluations
			finalRules := append(defRules, val.Rules...)
			finalEvals := append(defEvals, val.Evaluations...)
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
