package iigointernal

import (
	"fmt"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"github.com/SOMAS2020/SOMAS2020/internal/common/voting"
	"github.com/pkg/errors"
)

type judiciary struct {
	gameState         *gamestate.GameState
	JudgeID           shared.ClientID
	presidentSalary   shared.Resources
	BallotID          int
	ResAllocID        int
	speakerID         shared.ClientID
	presidentID       shared.ClientID
	EvaluationResults map[shared.ClientID]roles.EvaluationReturn
	clientJudge       roles.Judge
}

// loadClientJudge checks client pointer is good and if not panics
func (j *judiciary) loadClientJudge(clientJudgePointer roles.Judge) {
	if clientJudgePointer == nil {
		panic(fmt.Sprintf("Client '%v' has loaded a nil judge pointer", j.JudgeID))
	}
	j.clientJudge = clientJudgePointer
}

func (j *judiciary) init() {
	j.BallotID = 0
	j.ResAllocID = 0
}

// sendPresidentSalary conduct the transaction based on amount from client implementation
func (j *judiciary) sendPresidentSalary() error {
	if j.clientJudge != nil {
		amount, presidentPaid := j.clientJudge.PayPresident(j.presidentSalary)
		if presidentPaid {
			// Subtract from common resources po
			amountWithdraw, withdrawSuccess := WithdrawFromCommonPool(amount, j.gameState)

			if withdrawSuccess {
				// Pay into the client private resources pool
				depositIntoClientPrivatePool(amountWithdraw, j.gameState.PresidentID, j.gameState)
				return nil
			}
		}
	}
	return errors.Errorf("Cannot perform sendJudgeSalary")
}

// setSpeakerAndPresidentIDs set the speaker and president IDs.
func (j *judiciary) setSpeakerAndPresidentIDs(speakerID shared.ClientID, presidentID shared.ClientID) {
	j.speakerID = speakerID
	j.presidentID = presidentID
}

// InspectHistory checks all actions that happened in the last turn and audits them.
// This can be overridden by clients.
func (j *judiciary) inspectHistory(iigoHistory []shared.Accountability) (map[shared.ClientID]roles.EvaluationReturn, bool) {
	if !CheckEnoughInCommonPool(actionCost.InspectHistoryActionCost, j.gameState) {
		return nil, false
	}
	historyMap, ok := j.clientJudge.InspectHistory(iigoHistory)
	if ok && !j.incurServiceCharge(actionCost.InspectHistoryActionCost) {
		return nil, false
	}
	return historyMap, ok
}

// inspectBallot checks each ballot action adheres to the rules
func (j *judiciary) inspectBallot() (bool, error) {
	// 1. Evaluate difference between newRules and oldRules to check
	//    rule changes are in line with RuleToVote in previous ballot
	// 2. Compare each ballot action adheres to rules in ruleSet matrix

	if !CheckEnoughInCommonPool(actionCost.InspectHistoryActionCost, j.gameState) {
		return false, nil
	}
	rulesAffectedBySpeaker := j.EvaluationResults[j.speakerID]
	indexOfBallotRule, err := searchForRule("inspect_ballot_rule", rulesAffectedBySpeaker.Rules)
	if err {
		return rulesAffectedBySpeaker.Evaluations[indexOfBallotRule], nil
	}
	return true, errors.Errorf("Speaker did not conduct any ballots")
}

// inspectAllocation checks each resource allocation action adheres to the rules
func (j *judiciary) inspectAllocation() (bool, error) {
	// 1. Evaluate difference between commonPoolNew and commonPoolOld
	//    to check resource allocation changes are in line with ResourceRequests
	//    in previous resourceAllocation
	// 2. Compare each resource allocation action adheres to rules in ruleSet
	//    matrix
	if !j.incurServiceCharge(actionCost.InspectAllocationActionCost) {
		return false, errors.Errorf("Insufficient Budget in common Pool: inspectAllocation")
	}

	rulesAffectedByPresident := j.EvaluationResults[j.presidentID]
	indexOfAllocRule, ok := searchForRule("inspect_allocation_rule", rulesAffectedByPresident.Rules)
	if !ok {
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
	return -1, false
}

// declareSpeakerPerformanceWrapped wraps the result of DeclareSpeakerPerformance for orchestration
func (j *judiciary) declareSpeakerPerformanceWrapped() {
	res, err := j.inspectBallot()
	didRole := true
	if err != nil {
		didRole = false
	}
	result, checkRole := j.clientJudge.DeclareSpeakerPerformance(res, didRole)
	message := generateSpeakerPerformanceMessage(j.BallotID, result, j.speakerID, checkRole)
	broadcastToAllIslands(shared.TeamIDs[j.JudgeID], message)

}

// declarePresidentPerformanceWrapped wraps the result of DeclarePresidentPerformance for orchestration
func (j *judiciary) declarePresidentPerformanceWrapped() {
	res, err := j.inspectAllocation()
	didRole := true
	if err != nil {
		didRole = false
	}
	result, checkRole := j.clientJudge.DeclarePresidentPerformance(res, didRole)
	message := generatePresidentPerformanceMessage(j.ResAllocID, result, j.presidentID, checkRole)
	broadcastToAllIslands(shared.TeamIDs[j.JudgeID], message)

}

// appointNextPresident returns the island ID of the island appointed to be the president in the next turn
// appointing new roles should be free
func (j *judiciary) appointNextPresident(clientIDs []shared.ClientID) (shared.ClientID, error) {
	if !j.incurServiceCharge(actionCost.AppointNextPresidentActionCost) {
		return j.JudgeID, errors.Errorf("Insufficient Budget in common Pool: appointNextPresident")
	}
	var election voting.Election
	election.ProposeElection(baseclient.President, voting.Plurality)
	election.OpenBallot(clientIDs)
	election.Vote(iigoClients)
	return election.CloseBallot(), nil
}

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

func (j *judiciary) incurServiceCharge(cost shared.Resources) bool {
	_, ok := WithdrawFromCommonPool(cost, j.gameState)
	if ok {
		j.gameState.IIGORolesBudget["judge"] -= cost
	}
	return ok
}
