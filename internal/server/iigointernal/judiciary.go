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
	gameState             *gamestate.GameState
	JudgeID               shared.ClientID
	presidentSalary       shared.Resources
	EvaluationResults     map[shared.ClientID]roles.EvaluationReturn
	clientJudge           roles.Judge
	presidentTurnsInPower int
	sanctionRecord        map[shared.ClientID]roles.IIGOSanctionScore
	sanctionThresholds    map[roles.IIGOSanctionTier]roles.IIGOSanctionScore
	ruleViolationSeverity map[string]roles.IIGOSanctionScore
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
func (j *judiciary) appointNextPresident(currentPresident shared.ClientID, allIslands []shared.ClientID) (shared.ClientID, error) {
	var election voting.Election
	var nextPresident shared.ClientID
	electionsettings := j.clientJudge.CallPresidentElection(j.presidentTurnsInPower, allIslands)
	if electionsettings.HoldElection {
		if !j.incurServiceCharge(actionCost.InspectHistoryActionCost) {
			return j.gameState.PresidentID, errors.Errorf("Insufficient Budget in common Pool: appointNextPresident")
		}
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
	return nextPresident, nil
}

func (j *judiciary) incurServiceCharge(cost shared.Resources) bool {
	_, ok := WithdrawFromCommonPool(cost, j.gameState)
	if ok {
		j.gameState.IIGORolesBudget["judge"] -= cost
	}
	return ok
}

// Helper functions //

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
