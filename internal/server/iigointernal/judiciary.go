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
	JudgeID           shared.ClientID
	budget            shared.Resources
	presidentSalary   shared.Resources
	BallotID          int
	ResAllocID        int
	speakerID         shared.ClientID
	presidentID       shared.ClientID
	EvaluationResults map[shared.ClientID]roles.EvaluationReturn
	clientJudge       roles.Judge
}

func (j *judiciary) init() {
	j.BallotID = 0
	j.ResAllocID = 0
}

// returnPresidentSalary returns the salary to the common pool.
func (j *judiciary) returnPresidentSalary() shared.Resources {
	x := j.presidentSalary
	j.presidentSalary = 0
	return x
}

// withdrawPresidentSalary withdraws the president's salary from the common pool.
func (j *judiciary) withdrawPresidentSalary(gameState *gamestate.GameState) error {
	var presidentSalary = shared.Resources(rules.VariableMap["presidentSalary"].Values[0])
	var withdrawError = WithdrawFromCommonPool(presidentSalary, gameState)
	if withdrawError == nil {
		j.presidentSalary = presidentSalary
	}
	return withdrawError
}

// sendPresidentSalary sends the president's salary to the president.
func (j *judiciary) sendPresidentSalary() {
	if j.clientJudge != nil {
		amount := j.clientJudge.PayPresident()
		executiveBranch.budget = amount
		return
	}
	amount, _ := j.PayPresident()
	executiveBranch.budget = amount
}

// PayPresident pays the president salary.
func (j *judiciary) PayPresident() (shared.Resources, error) {
	hold := j.presidentSalary
	j.presidentSalary = 0
	return hold, nil
}

// setSpeakerAndPresidentIDs set the speaker and president IDs.
func (j *judiciary) setSpeakerAndPresidentIDs(speakerId shared.ClientID, presidentId shared.ClientID) {
	j.speakerID = speakerId
	j.presidentID = presidentId
}

// InspectHistory checks all actions that happened in the last turn and audits them.
// This can be overridden by clients.
func (j *judiciary) inspectHistory() (map[shared.ClientID]roles.EvaluationReturn, error) {
	j.budget -= 10
	return j.clientJudge.InspectHistory()
}

// inspectBallot checks each ballot action adheres to the rules
func (j *judiciary) inspectBallot() (bool, error) {
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
func (j *judiciary) inspectAllocation() (bool, error) {
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
	j.budget -= 10
	var election voting.Election
	election.ProposeElection(baseclient.President, voting.Plurality)
	election.OpenBallot(clientIDs)
	election.Vote(iigoClients)
	return election.CloseBallot()
}

// generateSpeakerPerformanceMessage generates the appropriate communication required regarding
// speaker performance to be sent to clients
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

// generatePresidentPerformanceMessage generated the appropriate communication required regarding
// president performance to be sent to clients
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
