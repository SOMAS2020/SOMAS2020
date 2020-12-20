package roles

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/pkg/errors"
)

// base Judge object
type BaseJudge struct {
	id                int
	budget            int
	presidentSalary   int
	BallotID          int
	ResAllocID        int
	speakerID         int
	presidentID       int
	evaluationResults map[int]EvaluationReturn
	clientJudge       Judge
}

func (j *BaseJudge) init() {
	j.BallotID = 0
	j.ResAllocID = 0
}

// Withdraw president salary from the common pool
// Call common withdraw function with president as parameter
func (j *BaseJudge) withdrawPresidentSalary(gameState *common.GameState) error {
	var presidentSalary = int(rules.VariableMap["presidentSalary"].Values[0])
	var withdrawError = WithdrawFromCommonPool(presidentSalary, gameState)
	if withdrawError != nil {
		Base_judge.presidentSalary = presidentSalary
	}
	return withdrawError
}

// Pay the president
func (j *BaseJudge) payPresident(gameState *common.GameState) {
	Base_President.budget = Base_judge.presidentSalary
	Base_judge.presidentSalary = 0
}

func (j *BaseJudge) setSpeakerAndPresidentIDs(speakerId int, presidentId int) {
	j.speakerID = speakerId
	j.presidentID = presidentId
}

type EvaluationReturn struct {
	rules       []rules.RuleMatrix
	evaluations []bool
}

func (j *BaseJudge) inspectHistoryInternal() {
	outputMap := map[int]EvaluationReturn{}
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
			tempTemp := EvaluationReturn{
				rules:       []rules.RuleMatrix{},
				evaluations: []bool{},
			}
			outputMap[clientID] = tempTemp
		}
		tempReturn := outputMap[clientID]
		for _, v3 := range rulesAffected {
			evaluation, _ := rules.BasicBooleanRuleEvaluator(v3)
			tempReturn.rules = append(tempReturn.rules, rules.RulesInPlay[v3])
			tempReturn.evaluations = append(tempReturn.evaluations, evaluation)
		}
	}
	j.evaluationResults = outputMap
}

func (j *BaseJudge) inspectHistory() (map[int]EvaluationReturn, error) {
	if j.clientJudge != nil {
		outputMap, err := j.clientJudge.inspectHistory()
		if err != nil {
			j.inspectHistoryInternal()
		} else {
			j.evaluationResults = outputMap
		}
	} else {
		j.inspectHistoryInternal()
	}
	return j.evaluationResults, nil
}

func (j *BaseJudge) inspectBallot() (bool, error) {
	// 1. Evaluate difference between newRules and oldRules to check
	//    rule changes are in line with ruleToVote in previous ballot
	// 2. Compare each ballot action adheres to rules in ruleSet matrix
	rulesAffectedBySpeaker := j.evaluationResults[j.speakerID]
	indexOfBallotRule, err := searchForRule("inspect_ballot_rule", rulesAffectedBySpeaker.rules)
	if err == nil {
		return rulesAffectedBySpeaker.evaluations[indexOfBallotRule], nil
	} else {
		return true, errors.Errorf("Speaker did not conduct any ballots")
	}
}

func (j *BaseJudge) inspectAllocation() (bool, error) {
	// 1. Evaluate difference between commonPoolNew and commonPoolOld
	//    to check resource allocation changes are in line with resourceRequests
	//    in previous resourceAllocation
	// 2. Compare each resource allocation action adheres to rules in ruleSet
	//    matrix
	rulesAffectedByPresident := j.evaluationResults[j.presidentID]
	indexOfAllocRule, err := searchForRule("inspect_allocation_rule", rulesAffectedByPresident.rules)
	if err == nil {
		return rulesAffectedByPresident.evaluations[indexOfAllocRule], nil
	} else {
		return true, errors.Errorf("President didn't conduct any allocations")
	}
}

func searchForRule(ruleName string, listOfRuleMatrices []rules.RuleMatrix) (int, error) {
	for i, v := range listOfRuleMatrices {
		if v.RuleName == ruleName {
			return i, nil
		}
	}
	return 0, errors.Errorf("The rule name '%v' was not found", ruleName)
}

func (j *BaseJudge) declareSpeakerPerformanceWrapped() {

	var BID int
	var result bool
	var SID int
	var checkRole bool
	var err error

	if j.clientJudge != nil {
		BID, result, SID, checkRole, err = j.clientJudge.declareSpeakerPerformance()
		if err == nil {
			BID, result, SID, checkRole, err = j.declareSpeakerPerformanceInternal()
		}
	} else {
		BID, result, SID, checkRole, err = j.declareSpeakerPerformanceInternal()
	}
	if err == nil {
		message := generateSpeakerPerformanceMessage(BID, result, SID, checkRole)
		broadcastToAllIslands(j.id, message)
	}
}

func (j *BaseJudge) declareSpeakerPerformanceInternal() (int, bool, int, bool, error) {
	j.BallotID++
	result, err := j.inspectBallot()

	conductedRole := err == nil

	return j.BallotID, result, j.speakerID, conductedRole, nil
}

func (j *BaseJudge) declarePresidentPerformanceWrapped() {

	var RID int
	var result bool
	var PID int
	var checkRole bool
	var err error

	if j.clientJudge != nil {
		RID, result, PID, checkRole, err = j.clientJudge.declarePresidentPerformance()
		if err == nil {
			RID, result, PID, checkRole, err = j.declarePresidentPerformanceInternal()
		}
	} else {
		RID, result, PID, checkRole, err = j.declarePresidentPerformanceInternal()
	}
	if err == nil {
		message := generatePresidentPerformanceMessage(RID, result, PID, checkRole)
		broadcastToAllIslands(j.id, message)
	}
}

func (j *BaseJudge) declarePresidentPerformanceInternal() (int, bool, int, bool, error) {

	j.ResAllocID++
	result, err := j.inspectAllocation()

	conductedRole := err == nil

	return j.ResAllocID, result, j.presidentID, conductedRole, nil
}

func generateSpeakerPerformanceMessage(BID int, result bool, SID int, conductedRole bool) map[int]DataPacket {
	returnMap := map[int]DataPacket{}

	returnMap[BallotID] = DataPacket{
		integerData: BID,
	}
	returnMap[SpeakerBallotCheck] = DataPacket{
		booleanData: result,
	}
	returnMap[SpeakerID] = DataPacket{
		integerData: SID,
	}
	returnMap[RoleConducted] = DataPacket{
		booleanData: conductedRole,
	}
	return returnMap
}

func generatePresidentPerformanceMessage(RID int, result bool, PID int, conductedRole bool) map[int]DataPacket {
	returnMap := map[int]DataPacket{}

	returnMap[ResAllocID] = DataPacket{
		integerData: RID,
	}
	returnMap[PresidentAllocationCheck] = DataPacket{
		booleanData: result,
	}
	returnMap[PresidentID] = DataPacket{
		integerData: PID,
	}
	returnMap[RoleConducted] = DataPacket{
		booleanData: conductedRole,
	}
	return returnMap
}
