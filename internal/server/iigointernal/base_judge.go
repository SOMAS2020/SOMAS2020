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

// base Judge object
type BaseJudge struct {
	Id                int
	budget            int
	presidentSalary   int
	BallotID          int
	ResAllocID        int
	speakerID         int
	presidentID       int
	EvaluationResults map[int]roles.EvaluationReturn
	clientJudge       roles.Judge
}

func (j *BaseJudge) init() {
	j.BallotID = 0
	j.ResAllocID = 0
}

// returnPresidentSalary returns the salary to the common pool.
func (j *BaseJudge) returnPresidentSalary() int {
	x := j.presidentSalary
	j.presidentSalary = 0
	return x
}

// withdrawPresidentSalary withdraws the president's salary from the common pool.
func (j *BaseJudge) withdrawPresidentSalary(gameState *gamestate.GameState) error {
	var presidentSalary = int(rules.VariableMap["presidentSalary"].Values[0])
	var withdrawError = WithdrawFromCommonPool(presidentSalary, gameState)
	if withdrawError == nil {
		featureJudge.presidentSalary = presidentSalary
	}
	return withdrawError
}

// sendPresidentSalary sends the president's salary to the president.
func (j *BaseJudge) sendPresidentSalary() {
	if j.clientJudge != nil {
		amount, err := j.clientJudge.PayPresident()
		if err == nil {
			featurePresident.budget = amount
			return
		}
	}
	amount, _ := j.PayPresident()
	featurePresident.budget = amount
}

// PayPresident pays the president salary.
func (j *BaseJudge) PayPresident() (int, error) {
	hold := j.presidentSalary
	j.presidentSalary = 0
	return hold, nil
}

// setSpeakerAndPresidentIDs set the speaker and president IDs.
func (j *BaseJudge) setSpeakerAndPresidentIDs(speakerId int, presidentId int) {
	j.speakerID = speakerId
	j.presidentID = presidentId
}

// inspectHistoryInternal is the base implementation of InspectHistory.
func (j *BaseJudge) inspectHistoryInternal() {
	outputMap := map[int]roles.EvaluationReturn{}
	for _, v := range TurnHistory {
		variablePairs := v.pairs
		clientID := v.clientID
		var rulesAffected []string
		for _, v2 := range variablePairs {
			valuesToBeAdded, err := PickUpRulesByVariable(v2.VariableName, rules.RulesInPlay)
			if err == nil {
				rulesAffected = append(rulesAffected, valuesToBeAdded...)
			}
			err = rules.UpdateVariable(v2.VariableName, v2)
			if err != nil {
				return
			}
		}
		if _, ok := outputMap[clientID]; !ok {
			tempTemp := roles.EvaluationReturn{
				Rules:       []rules.RuleMatrix{},
				Evaluations: []bool{},
			}
			outputMap[clientID] = tempTemp
		}
		tempReturn := outputMap[clientID]
		for _, v3 := range rulesAffected {
			evaluation, _ := rules.BasicBooleanRuleEvaluator(v3)
			tempReturn.Rules = append(tempReturn.Rules, rules.RulesInPlay[v3])
			tempReturn.Evaluations = append(tempReturn.Evaluations, evaluation)
		}
	}
	j.EvaluationResults = outputMap
}

// InspectHistory checks all actions that happened in the last turn and audits them.
// This can be overridden by clients.
func (j *BaseJudge) InspectHistory() (map[int]roles.EvaluationReturn, error) {
	j.budget -= 10 // will be removed post-MVP
	if j.clientJudge != nil {
		outputMap, err := j.clientJudge.InspectHistory()
		if err != nil {
			j.inspectHistoryInternal()
		} else {
			j.EvaluationResults = outputMap
		}
	} else {
		j.inspectHistoryInternal()
	}
	return j.EvaluationResults, nil
}

// inspectBallot checks each ballot action adheres to the rules
func (j *BaseJudge) inspectBallot() (bool, error) {
	// 1. Evaluate difference between newRules and oldRules to check
	//    rule changes are in line with RuleToVote in previous ballot
	// 2. Compare each ballot action adheres to rules in ruleSet matrix
	j.budget -= 10 // will be removed post-MVP
	rulesAffectedBySpeaker := j.EvaluationResults[j.speakerID]
	indexOfBallotRule, err := searchForRule("inspect_ballot_rule", rulesAffectedBySpeaker.Rules)
	if err == nil {
		return rulesAffectedBySpeaker.Evaluations[indexOfBallotRule], nil
	} else {
		return true, errors.Errorf("Speaker did not conduct any ballots")
	}
}

// inspectAllocation checks each resource allocation action adheres to the rules
func (j *BaseJudge) inspectAllocation() (bool, error) {
	// 1. Evaluate difference between commonPoolNew and commonPoolOld
	//    to check resource allocation changes are in line with ResourceRequests
	//    in previous resourceAllocation
	// 2. Compare each resource allocation action adheres to rules in ruleSet
	//    matrix
	j.budget -= 10 // will be removed post-MVP
	rulesAffectedByPresident := j.EvaluationResults[j.presidentID]
	indexOfAllocRule, err := searchForRule("inspect_allocation_rule", rulesAffectedByPresident.Rules)
	if err == nil {
		return rulesAffectedByPresident.Evaluations[indexOfAllocRule], nil
	} else {
		return true, errors.Errorf("President didn't conduct any allocations")
	}
}

// searchForRule searches for a given rule in the RuleMatrix
func searchForRule(ruleName string, listOfRuleMatrices []rules.RuleMatrix) (int, error) {
	for i, v := range listOfRuleMatrices {
		if v.RuleName == ruleName {
			return i, nil
		}
	}
	return 0, errors.Errorf("The rule name '%v' was not found", ruleName)
}

// declareSpeakerPerformanceInternal is the base implementation of DeclareSpeakerPerformance
func (j *BaseJudge) declareSpeakerPerformanceInternal() (int, bool, int, bool, error) {
	j.BallotID++
	result, err := j.inspectBallot()

	conductedRole := err == nil

	return j.BallotID, result, j.speakerID, conductedRole, nil
}

// DeclareSpeakerPerformance checks how well the speaker did their job
func (j *BaseJudge) DeclareSpeakerPerformance() (int, bool, int, bool, error) {

	j.budget -= 10 // will be removed post-MVP
	var BID int
	var result bool
	var SID int
	var checkRole bool
	var err error

	if j.clientJudge != nil {
		BID, result, SID, checkRole, err = j.clientJudge.DeclareSpeakerPerformance()
		if err == nil {
			BID, result, SID, checkRole, err = j.declareSpeakerPerformanceInternal()
		}
	} else {
		BID, result, SID, checkRole, err = j.declareSpeakerPerformanceInternal()
	}
	return BID, result, SID, checkRole, err
}

// declareSpeakerPerformanceWrapped wraps the result of DeclareSpeakerPerformance for orchestration
func (j *BaseJudge) declareSpeakerPerformanceWrapped() {

	BID, result, SID, checkRole, err := j.DeclareSpeakerPerformance()

	if err == nil {
		message := generateSpeakerPerformanceMessage(BID, result, SID, checkRole)
		broadcastToAllIslands(shared.TeamIDs[j.Id], message)
	}
}

// DeclarePresidentPerformance checks how well the president did their job
func (j *BaseJudge) DeclarePresidentPerformance() (int, bool, int, bool, error) {
	j.budget -= 10 // will be removed post-MVP
	var RID int
	var result bool
	var PID int
	var checkRole bool
	var err error

	if j.clientJudge != nil {
		RID, result, PID, checkRole, err = j.clientJudge.DeclarePresidentPerformance()
		if err == nil {
			RID, result, PID, checkRole, err = j.declarePresidentPerformanceInternal()
		}
	} else {
		RID, result, PID, checkRole, err = j.declarePresidentPerformanceInternal()
	}

	return RID, result, PID, checkRole, err
}

// declarePresidentPerformanceWrapped wraps the result of DeclarePresidentPerformance for orchestration
func (j *BaseJudge) declarePresidentPerformanceWrapped() {

	RID, result, PID, checkRole, err := j.DeclarePresidentPerformance()

	if err == nil {
		message := generatePresidentPerformanceMessage(RID, result, PID, checkRole)
		broadcastToAllIslands(shared.TeamIDs[j.Id], message)
	}
}

// declarePresidentPerformanceInternal is the base implementation for DeclarePresidentPerformance
func (j *BaseJudge) declarePresidentPerformanceInternal() (int, bool, int, bool, error) {

	j.ResAllocID++
	result, err := j.inspectAllocation()

	conductedRole := err == nil

	return j.ResAllocID, result, j.presidentID, conductedRole, nil
}

// appointNextPresident returns the island ID of the island appointed to be the president in the next turn
func (j *BaseJudge) appointNextPresident(clientIDs []shared.ClientID) int {
	j.budget -= 10 // will be removed post-MVP
	var election voting.Election
	election.ProposeElection(baseclient.President, voting.Plurality)
	election.OpenBallot(clientIDs)
	election.Vote(iigoClients)
	return int(election.CloseBallot())
}

// generateSpeakerPerformanceMessage generates the appropriate communication required regarding
// speaker performance to be sent to clients
func generateSpeakerPerformanceMessage(BID int, result bool, SID int, conductedRole bool) map[int]baseclient.Communication {
	returnMap := map[int]baseclient.Communication{}

	returnMap[BallotID] = baseclient.Communication{
		T:           baseclient.CommunicationInt,
		IntegerData: BID,
	}
	returnMap[SpeakerBallotCheck] = baseclient.Communication{
		T:           baseclient.CommunicationBool,
		BooleanData: result,
	}
	returnMap[SpeakerID] = baseclient.Communication{
		T:           baseclient.CommunicationInt,
		IntegerData: SID,
	}
	returnMap[RoleConducted] = baseclient.Communication{
		T:           baseclient.CommunicationBool,
		BooleanData: conductedRole,
	}
	return returnMap
}

// generatePresidentPerformanceMessage generated the appropriate communication required regarding
// president performance to be sent to clients
func generatePresidentPerformanceMessage(RID int, result bool, PID int, conductedRole bool) map[int]baseclient.Communication {
	returnMap := map[int]baseclient.Communication{}

	returnMap[ResAllocID] = baseclient.Communication{
		T:           baseclient.CommunicationInt,
		IntegerData: RID,
	}
	returnMap[PresidentAllocationCheck] = baseclient.Communication{
		T:           baseclient.CommunicationBool,
		BooleanData: result,
	}
	returnMap[PresidentID] = baseclient.Communication{
		T:           baseclient.CommunicationInt,
		IntegerData: PID,
	}
	returnMap[RoleConducted] = baseclient.Communication{
		T:           baseclient.CommunicationBool,
		BooleanData: conductedRole,
	}
	return returnMap
}
