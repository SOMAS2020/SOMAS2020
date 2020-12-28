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
	sanctionThresholds    map[roles.IIGOSanctionTier]int
	ruleViolationSeverity map[string]roles.IIGOSanctionScore
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
	return j.clientJudge.InspectHistory(iigoHistory)
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

func getDefaultSanctionThresholds() map[roles.IIGOSanctionTier]int {
	return map[roles.IIGOSanctionTier]int{
		roles.SanctionTier1: 0,
		roles.SanctionTier2: 1,
		roles.SanctionTier3: 10,
		roles.SanctionTier4: 20,
	}
}

func softMergeSanctionThresholds(clientSanctionMap map[roles.IIGOSanctionTier]int) map[roles.IIGOSanctionTier]int {
	defaultMap := getDefaultSanctionThresholds()
	for k, _ := range defaultMap {
		if clientVal, ok := clientSanctionMap[k]; ok {
			defaultMap[k] = clientVal
		}
	}
	return defaultMap
}

func unpackSingleIslandTransgressions(evaluationReturn roles.EvaluationReturn) []string {
	var transgressions []string
	for i, v := range evaluationReturn.Rules {
		if !evaluationReturn.Evaluations[i] {
			transgressions = append(transgressions, v.RuleName)
		}
	}
	return transgressions
}
