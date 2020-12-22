package iigointernal

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"math/rand"

	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
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
	evaluationResults map[int]roles.EvaluationReturn
	clientJudge       roles.Judge
}

func (j *BaseJudge) init() {
	j.BallotID = 0
	j.ResAllocID = 0
}

// Withdraw president salary from the common pool
// Call common withdraw function with president as parameter
func (j *BaseJudge) withdrawPresidentSalary(gameState *gamestate.GameState) error {
	var presidentSalary = int(rules.VariableMap["presidentSalary"].Values[0])
	var withdrawError = WithdrawFromCommonPool(presidentSalary, gameState)
	if withdrawError != nil {
		featureJudge.presidentSalary = presidentSalary
	}
	return withdrawError
}

func (j *BaseJudge) sendPresidentSalary() {
	if j.clientJudge != nil {
		amount, err := j.clientJudge.PayPresident()
		if err == nil {
			featurePresident.budget = amount
			return
		}
	}
	amount, _ := j.payPresident()
	featurePresident.budget = amount
}

// Pay the president
func (j *BaseJudge) payPresident() (int, error) {
	hold := j.presidentSalary
	j.presidentSalary = 0
	return hold, nil
}

func (j *BaseJudge) setSpeakerAndPresidentIDs(speakerId int, presidentId int) {
	j.speakerID = speakerId
	j.presidentID = presidentId
}

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
	j.evaluationResults = outputMap
}

func (j *BaseJudge) inspectHistory() (map[int]roles.EvaluationReturn, error) {
	j.budget -= 10
	if j.clientJudge != nil {
		outputMap, err := j.clientJudge.InspectHistory()
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
	j.budget -= 10
	rulesAffectedBySpeaker := j.evaluationResults[j.speakerID]
	indexOfBallotRule, err := searchForRule("inspect_ballot_rule", rulesAffectedBySpeaker.Rules)
	if err == nil {
		return rulesAffectedBySpeaker.Evaluations[indexOfBallotRule], nil
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
	j.budget -= 10
	rulesAffectedByPresident := j.evaluationResults[j.presidentID]
	indexOfAllocRule, err := searchForRule("inspect_allocation_rule", rulesAffectedByPresident.Rules)
	if err == nil {
		return rulesAffectedByPresident.Evaluations[indexOfAllocRule], nil
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

func (j *BaseJudge) declareSpeakerPerformanceInternal() (int, bool, int, bool, error) {
	j.BallotID++
	result, err := j.inspectBallot()

	conductedRole := err == nil

	return j.BallotID, result, j.speakerID, conductedRole, nil
}

func (j *BaseJudge) declareSpeakerPerformance() (int, bool, int, bool, error) {

	j.budget -= 10
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

func (j *BaseJudge) declareSpeakerPerformanceWrapped() {

	BID, result, SID, checkRole, err := j.declareSpeakerPerformance()

	if err == nil {
		message := generateSpeakerPerformanceMessage(BID, result, SID, checkRole)
		broadcastToAllIslands(j.id, message)
	}
}

func (j *BaseJudge) declarePresidentPerformance() (int, bool, int, bool, error) {

	j.budget -= 10
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

func (j *BaseJudge) declarePresidentPerformanceWrapped() {

	RID, result, PID, checkRole, err := j.declarePresidentPerformance()

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

func (j *BaseJudge) appointNextPresident() int {
	j.budget -= 10
	return rand.Intn(5)
}

func generateSpeakerPerformanceMessage(BID int, result bool, SID int, conductedRole bool) map[int]baseclient.Communication {
	returnMap := map[int]baseclient.Communication{}

	returnMap[BallotID] = baseclient.Communication{
		IntegerData: BID,
	}
	returnMap[SpeakerBallotCheck] = baseclient.Communication{
		BooleanData: result,
	}
	returnMap[SpeakerID] = baseclient.Communication{
		IntegerData: SID,
	}
	returnMap[RoleConducted] = baseclient.Communication{
		BooleanData: conductedRole,
	}
	return returnMap
}

func generatePresidentPerformanceMessage(RID int, result bool, PID int, conductedRole bool) map[int]baseclient.Communication {
	returnMap := map[int]baseclient.Communication{}

	returnMap[ResAllocID] = baseclient.Communication{
		IntegerData: RID,
	}
	returnMap[PresidentAllocationCheck] = baseclient.Communication{
		BooleanData: result,
	}
	returnMap[PresidentID] = baseclient.Communication{
		IntegerData: PID,
	}
	returnMap[RoleConducted] = baseclient.Communication{
		BooleanData: conductedRole,
	}
	return returnMap
}
