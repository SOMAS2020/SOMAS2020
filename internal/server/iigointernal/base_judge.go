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
	Id                shared.ClientID
	budget            shared.Resources
	presidentSalary   shared.Resources
	BallotID          int
	ResAllocID        int
	speakerID         shared.ClientID
	presidentID       shared.ClientID
	EvaluationResults map[shared.ClientID]roles.EvaluationReturn
	clientJudge       roles.Judge
}

func (j *BaseJudge) init() {
	j.BallotID = 0
	j.ResAllocID = 0
}

// Withdraw president salary from the common pool
// Call common withdraw function with president as parameter
func (j *BaseJudge) withdrawPresidentSalary(gameState *gamestate.GameState) error {
	var presidentSalary = shared.Resources(rules.VariableMap["presidentSalary"].Values[0])
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
	amount, _ := j.PayPresident()
	featurePresident.budget = amount
}

// Pay the president
func (j *BaseJudge) PayPresident() (shared.Resources, error) {
	hold := j.presidentSalary
	j.presidentSalary = 0
	return hold, nil
}

func (j *BaseJudge) setSpeakerAndPresidentIDs(speakerId shared.ClientID, presidentId shared.ClientID) {
	j.speakerID = speakerId
	j.presidentID = presidentId
}

func (j *BaseJudge) inspectHistoryInternal() {
	outputMap := map[shared.ClientID]roles.EvaluationReturn{}
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

func (j *BaseJudge) InspectHistory() (map[shared.ClientID]roles.EvaluationReturn, error) {
	j.budget -= 10
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

func (j *BaseJudge) inspectBallot() (bool, error) {
	// 1. Evaluate difference between newRules and oldRules to check
	//    rule changes are in line with RuleToVote in previous ballot
	// 2. Compare each ballot action adheres to rules in ruleSet matrix
	j.budget -= 10
	rulesAffectedBySpeaker := j.EvaluationResults[j.speakerID]
	indexOfBallotRule, err := searchForRule("inspect_ballot_rule", rulesAffectedBySpeaker.Rules)
	if err == nil {
		return rulesAffectedBySpeaker.Evaluations[indexOfBallotRule], nil
	} else {
		return true, errors.Errorf("Speaker did not conduct any ballots")
	}
}

func (j *BaseJudge) inspectAllocation() (bool, error) {
	// 1. Evaluate difference between commonPoolNew and commonPoolOld
	//    to check resource allocation changes are in line with ResourceRequests
	//    in previous resourceAllocation
	// 2. Compare each resource allocation action adheres to rules in ruleSet
	//    matrix
	j.budget -= 10
	rulesAffectedByPresident := j.EvaluationResults[j.presidentID]
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

func (j *BaseJudge) declareSpeakerPerformanceInternal() (BID int, result bool, SID shared.ClientID, checkRole bool, err error) {
	j.BallotID++
	result, err = j.inspectBallot()

	conductedRole := err == nil

	return j.BallotID, result, j.speakerID, conductedRole, nil
}

func (j *BaseJudge) DeclareSpeakerPerformance() (BID int, result bool, SID shared.ClientID, checkRole bool, err error) {
	j.budget -= 10

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

	BID, result, SID, checkRole, err := j.DeclareSpeakerPerformance()

	if err == nil {
		message := generateSpeakerPerformanceMessage(BID, result, SID, checkRole)
		broadcastToAllIslands(shared.TeamIDs[j.Id], message)
	}
}

func (j *BaseJudge) DeclarePresidentPerformance() (RID int, result bool, PID shared.ClientID, checkRole bool, err error) {
	j.budget -= 10

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

	RID, result, PID, checkRole, err := j.DeclarePresidentPerformance()

	if err == nil {
		message := generatePresidentPerformanceMessage(RID, result, PID, checkRole)
		broadcastToAllIslands(shared.TeamIDs[j.Id], message)
	}
}

func (j *BaseJudge) declarePresidentPerformanceInternal() (int, bool, shared.ClientID, bool, error) {

	j.ResAllocID++
	result, err := j.inspectAllocation()

	conductedRole := err == nil

	return j.ResAllocID, result, j.presidentID, conductedRole, nil
}

func (j *BaseJudge) appointNextPresident(clientIDs []shared.ClientID) shared.ClientID {
	j.budget -= 10
	var election voting.Election
	election.ProposeElection(baseclient.President, voting.Plurality)
	election.OpenBallot(clientIDs)
	election.Vote(iigoClients)
	return election.CloseBallot()
}

func generateSpeakerPerformanceMessage(BID int, result bool, SID shared.ClientID, conductedRole bool) map[int]baseclient.Communication {
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
		IntegerData: int(SID),
	}
	returnMap[RoleConducted] = baseclient.Communication{
		T:           baseclient.CommunicationBool,
		BooleanData: conductedRole,
	}
	return returnMap
}

func generatePresidentPerformanceMessage(RID int, result bool, PID shared.ClientID, conductedRole bool) map[int]baseclient.Communication {
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
		IntegerData: int(PID),
	}
	returnMap[RoleConducted] = baseclient.Communication{
		T:           baseclient.CommunicationBool,
		BooleanData: conductedRole,
	}
	return returnMap
}
