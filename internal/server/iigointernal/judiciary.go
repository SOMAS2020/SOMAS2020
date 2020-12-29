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

type judiciary struct {
	JudgeID               shared.ClientID
	budget                shared.Resources
	presidentSalary       shared.Resources
	EvaluationResults     map[shared.ClientID]roles.EvaluationReturn
	clientJudge           roles.Judge
	presidentTurnsInPower int
}

// loadClientJudge checks client pointer is good and if not panics
func (j *judiciary) loadClientJudge(clientJudgePointer roles.Judge) {
	if clientJudgePointer == nil {
		panic(fmt.Sprintf("Client '%v' has loaded a nil judge pointer", j.JudgeID))
	}
	j.clientJudge = clientJudgePointer
}

// returnPresidentSalary returns the salary to the common pooj.
func (j *judiciary) returnPresidentSalary() shared.Resources {
	x := j.presidentSalary
	j.presidentSalary = 0
	return x
}

// withdrawPresidentSalary withdraws the president's salary from the common pooj.
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

// InspectHistory checks all actions that happened in the last turn and audits them.
// This can be overridden by clients.
func (j *judiciary) inspectHistory(iigoHistory []shared.Accountability) (map[shared.ClientID]roles.EvaluationReturn, bool) {
	j.budget -= serviceCharge
	return j.clientJudge.InspectHistory(iigoHistory)
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
func (j *judiciary) appointNextPresident(currentPresident shared.ClientID) shared.ClientID {
	var election voting.Election
	var nextPresident shared.ClientID
	electionsettings := j.clientJudge.CallPresidentElection(j.presidentTurnsInPower)
	if electionsettings.HoldElection {
		// TODO: deduct the cost of holding an election
		election.ProposeElection(baseclient.President, electionsettings.VotingMethod)
		election.OpenBallot(electionsettings.IslandsToVote)
		election.Vote(iigoClients)
		j.presidentTurnsInPower = 0
		nextPresident = election.CloseBallot()
		nextPresident = j.clientJudge.DecideNextPresident(nextPresident)
	} else {
		j.presidentTurnsInPower += 1
		nextPresident = currentPresident
	}
	return nextPresident
}
