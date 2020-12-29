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

// to be moved to paramters
const sanctionCacheDepth = 3
const historyCacheDepth = 3

// to be changed
const sanctionLength = 2

type judiciary struct {
	JudgeID               shared.ClientID
	budget                shared.Resources
	presidentSalary       shared.Resources
	BallotID              int
	ResAllocID            int
	speakerID             shared.ClientID
	presidentID           shared.ClientID
	evaluationResults     map[shared.ClientID]roles.EvaluationReturn
	clientJudge           roles.Judge
	sanctionRecord        map[shared.ClientID]roles.IIGOSanctionScore
	sanctionThresholds    map[roles.IIGOSanctionTier]roles.IIGOSanctionScore
	ruleViolationSeverity map[string]roles.IIGOSanctionScore
	localSanctionCache    map[int][]roles.Sanction
	localHistoryCache     map[int][]shared.Accountability
}

func (j *judiciary) init() {
	j.BallotID = 0
	j.ResAllocID = 0
}

// Loads ruleViolationSeverity and sanction thresholds
func (j *judiciary) loadSanctionConfig() {
	j.sanctionThresholds = softMergeSanctionThresholds(j.clientJudge.GetSanctionThresholds())
	j.ruleViolationSeverity = j.clientJudge.GetRuleViolationSeverity()
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
		amount, payPresident := j.clientJudge.PayPresident()
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

// setSpeakerAndPresidentIDs set the speaker and president IDs.
func (j *judiciary) setSpeakerAndPresidentIDs(speakerId shared.ClientID, presidentId shared.ClientID) {
	j.speakerID = speakerId
	j.presidentID = presidentId
}

// InspectHistory checks all actions that happened in the last turn and audits them.
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

// inspectBallot checks each ballot action adheres to the rules
func (j *judiciary) inspectBallot() (bool, error) {
	// 1. Evaluate difference between newRules and oldRules to check
	//    rule changes are in line with RuleToVote in previous ballot
	// 2. Compare each ballot action adheres to rules in ruleSet matrix
	j.budget -= serviceCharge // will be removed post-MVP
	rulesAffectedBySpeaker := j.evaluationResults[j.speakerID]
	indexOfBallotRule, err := searchForRule("inspect_ballot_rule", rulesAffectedBySpeaker.Rules)
	if err {
		return rulesAffectedBySpeaker.Evaluations[indexOfBallotRule], nil
	} else {
		return true, errors.Errorf("Speaker did not conduct any ballots")
	}
}

// inspectAllocation checks each resource allocation action adheres to the rules
func (j *judiciary) inspectAllocation() (bool, error) {
	// 1. Evaluate difference between commonPoolNew and commonPoolOld
	//    to check resource allocation changes are in line with ResourceRequests
	//    in previous resourceAllocation
	// 2. Compare each resource allocation action adheres to rules in ruleSet
	//    matrix
	j.budget -= serviceCharge // will be removed post-MVP
	rulesAffectedByPresident := j.evaluationResults[j.presidentID]
	indexOfAllocRule, err := searchForRule("inspect_allocation_rule", rulesAffectedByPresident.Rules)
	if err {
		return true, errors.Errorf("President didn't conduct any allocations")
	}
	return rulesAffectedByPresident.Evaluations[indexOfAllocRule], nil
}

// searchForRule searches for a given rule in the RuleMatrix
func searchForRule(ruleName string, listOfRuleMatrices []rules.RuleMatrix) (int, bool) {
	for i, v := range listOfRuleMatrices {
		if v.RuleName == ruleName {
			return i, true
		}
	}
	return 0, false
}

// declareSpeakerPerformanceWrapped wraps the result of DeclareSpeakerPerformance for orchestration
func (j *judiciary) declareSpeakerPerformanceWrapped() {
	result, checkRole := j.clientJudge.DeclareSpeakerPerformance()
	message := generateSpeakerPerformanceMessage(j.BallotID, result, j.speakerID, checkRole)
	broadcastToAllIslands(shared.TeamIDs[j.JudgeID], message)

}

// declarePresidentPerformanceWrapped wraps the result of DeclarePresidentPerformance for orchestration
func (j *judiciary) declarePresidentPerformanceWrapped() {
	result, checkRole := j.clientJudge.DeclarePresidentPerformance()
	message := generatePresidentPerformanceMessage(j.ResAllocID, result, j.presidentID, checkRole)
	broadcastToAllIslands(shared.TeamIDs[j.JudgeID], message)

}

// appointNextPresident returns the island ID of the island appointed to be the president in the next turn
func (j *judiciary) appointNextPresident(clientIDs []shared.ClientID) shared.ClientID {
	j.budget -= serviceCharge
	var election voting.Election
	election.ProposeElection(baseclient.President, voting.Plurality)
	election.OpenBallot(clientIDs)
	election.Vote(iigoClients)
	return election.CloseBallot()
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

}

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

func (j *judiciary) clearHistoryCache() {
	j.localHistoryCache = map[int][]shared.Accountability{}
}

// Helper functions for Judiciary branch

// generateSpeakerPerformanceMessage generates the appropriate communication required regarding
// speaker performance to be sent to clients
func generateSpeakerPerformanceMessage(BID int, result bool, SID shared.ClientID, conductedRole bool) map[shared.CommunicationFieldName]shared.CommunicationContent {
	returnMap := map[shared.CommunicationFieldName]shared.CommunicationContent{}

	returnMap[shared.BallotID] = shared.CommunicationContent{
		T:           shared.CommunicationInt,
		IntegerData: BID,
	}
	returnMap[shared.SpeakerBallotCheck] = shared.CommunicationContent{
		T:           shared.CommunicationBool,
		BooleanData: result,
	}
	returnMap[shared.SpeakerID] = shared.CommunicationContent{
		T:           shared.CommunicationInt,
		IntegerData: int(SID),
	}
	returnMap[shared.RoleConducted] = shared.CommunicationContent{
		T:           shared.CommunicationBool,
		BooleanData: conductedRole,
	}
	return returnMap
}

// generatePresidentPerformanceMessage generated the appropriate communication required regarding
// president performance to be sent to clients
func generatePresidentPerformanceMessage(RID int, result bool, PID shared.ClientID, conductedRole bool) map[shared.CommunicationFieldName]shared.CommunicationContent {
	returnMap := map[shared.CommunicationFieldName]shared.CommunicationContent{}

	returnMap[shared.ResAllocID] = shared.CommunicationContent{
		T:           shared.CommunicationInt,
		IntegerData: RID,
	}
	returnMap[shared.PresidentAllocationCheck] = shared.CommunicationContent{
		T:           shared.CommunicationBool,
		BooleanData: result,
	}
	returnMap[shared.PresidentID] = shared.CommunicationContent{
		T:           shared.CommunicationInt,
		IntegerData: int(PID),
	}
	returnMap[shared.RoleConducted] = shared.CommunicationContent{
		T:           shared.CommunicationBool,
		BooleanData: conductedRole,
	}
	return returnMap
}

func getDefaultSanctionThresholds() map[roles.IIGOSanctionTier]roles.IIGOSanctionScore {
	return map[roles.IIGOSanctionTier]roles.IIGOSanctionScore{
		roles.SanctionTier1: 1,
		roles.SanctionTier2: 5,
		roles.SanctionTier3: 10,
		roles.SanctionTier4: 20,
		roles.SanctionTier5: 30,
	}
}

func softMergeSanctionThresholds(clientSanctionMap map[roles.IIGOSanctionTier]roles.IIGOSanctionScore) map[roles.IIGOSanctionTier]roles.IIGOSanctionScore {
	defaultMap := getDefaultSanctionThresholds()
	for k, _ := range defaultMap {
		if clientVal, ok := clientSanctionMap[k]; ok {
			defaultMap[k] = clientVal
		}
	}
	return defaultMap
}

func unpackSingleIslandTransgressions(evaluationReturn roles.EvaluationReturn) []string {
	transgressions := []string{}
	for i, v := range evaluationReturn.Rules {
		if !evaluationReturn.Evaluations[i] {
			transgressions = append(transgressions, v.RuleName)
		}
	}
	return transgressions
}

func getIslandSanctionTier(islandScore roles.IIGOSanctionScore, scoreMap map[roles.IIGOSanctionTier]roles.IIGOSanctionScore) roles.IIGOSanctionTier {
	if islandScore <= scoreMap[roles.SanctionTier1] {
		return roles.None
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

func getTierSanctionMap() map[roles.IIGOSanctionTier]string {
	return map[roles.IIGOSanctionTier]string{
		roles.SanctionTier1: "iigo_economic_sanction_1",
		roles.SanctionTier2: "iigo_economic_sanction_2",
		roles.SanctionTier3: "iigo_economic_sanction_3",
		roles.SanctionTier4: "iigo_economic_sanction_4",
		roles.SanctionTier5: "iigo_economic_sanction_5",
	}
}

func defaultInitLocalSanctionCache(depth int) map[int][]roles.Sanction {
	returnMap := map[int][]roles.Sanction{}
	for i := 0; i < depth; i++ {
		returnMap[i] = []roles.Sanction{}
	}
	return returnMap
}

func defaultInitLocalHistoryCache(depth int) map[int][]shared.Accountability {
	returnMap := map[int][]shared.Accountability{}
	for i := 0; i < depth; i++ {
		returnMap[i] = []shared.Accountability{}
	}
	return returnMap
}

func checkPardons(sanctionCache map[int][]roles.Sanction, pardons map[int]map[int]roles.Sanction) (pardonsValid bool, communications map[shared.ClientID][]map[shared.CommunicationFieldName]shared.CommunicationContent, finalCache map[int][]roles.Sanction) {
	comms := map[shared.ClientID][]map[shared.CommunicationFieldName]shared.CommunicationContent{}
	newSanctionCache := map[int][]roles.Sanction{}
	for k, v := range sanctionCache {
		newSanctionCache[k] = v
	}
	for i, v := range pardons {
		for iSan, vSan := range v {
			if sanctionCache[i][iSan] != vSan {
				return false, map[shared.ClientID][]map[shared.CommunicationFieldName]shared.CommunicationContent{}, sanctionCache
			} else {
				comms[vSan.ClientID] = append(comms[vSan.ClientID], map[shared.CommunicationFieldName]shared.CommunicationContent{
					shared.PardonClientID: {
						T:           shared.CommunicationInt,
						IntegerData: int(vSan.ClientID),
					},
					shared.PardonTier: {
						T:           shared.CommunicationInt,
						IntegerData: int(vSan.SanctionTier),
					},
				})
				copyOfNewSanctionCache := make([]roles.Sanction, len(newSanctionCache[i]))
				copy(copyOfNewSanctionCache, newSanctionCache[i])
				newSanctionCache[i] = removeSanctions(copyOfNewSanctionCache, iSan-getDifferenceInLength(sanctionCache[i], copyOfNewSanctionCache))
			}
		}
	}
	return true, comms, newSanctionCache
}

func removeSanctions(slice []roles.Sanction, s int) []roles.Sanction {
	var output []roles.Sanction
	for i, v := range slice {
		if i != s {
			output = append(output, v)
		}
	}
	return output
}

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

func streamlineRulesAffected(input []string) []string {
	streamlineMap := map[string]bool{}
	for _, v := range input {
		streamlineMap[v] = true
	}
	var returnArray []string
	for key, _ := range streamlineMap {
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
