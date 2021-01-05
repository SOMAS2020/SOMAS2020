package iigointernal

import (
	"fmt"
	"reflect"

	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"github.com/SOMAS2020/SOMAS2020/internal/common/voting"
	"github.com/pkg/errors"
)

type legislature struct {
	gameState         *gamestate.GameState
	gameConf          *config.IIGOConfig
	SpeakerID         shared.ClientID
	judgeSalary       shared.Resources
	ruleToVote        rules.RuleMatrix
	ballotBox         voting.BallotBox
	votingResult      bool
	clientSpeaker     roles.Speaker
	judgeTurnsInPower int
	monitoring        *monitor
}

// loadClientSpeaker checks client pointer is good and if not panics
func (l *legislature) loadClientSpeaker(clientSpeakerPointer roles.Speaker) {
	if clientSpeakerPointer == nil {
		panic(fmt.Sprintf("Client '%v' has loaded a nil speaker pointer", l.SpeakerID))
	}
	l.clientSpeaker = clientSpeakerPointer
}

func (l *legislature) syncWithGame(gameState *gamestate.GameState, gameConf *config.IIGOConfig) {
	l.gameState = gameState
	l.gameConf = gameConf
}

// sendJudgeSalary conduct the transaction based on amount from client implementation
func (l *legislature) sendJudgeSalary() error {
	if l.clientSpeaker != nil {
		amountReturn := l.clientSpeaker.PayJudge(l.judgeSalary)
		if amountReturn.ActionTaken && amountReturn.ContentType == shared.SpeakerJudgeSalary {
			// Subtract from common resources pool
			amountWithdraw, withdrawSuccess := WithdrawFromCommonPool(amountReturn.JudgeSalary, l.gameState)

			if withdrawSuccess {
				// Pay into the client private resources pool
				depositIntoClientPrivatePool(amountWithdraw, l.gameState.JudgeID, l.gameState)
				return nil
			}
		}
	}
	return errors.Errorf("Cannot perform sendJudgeSalary")
}

// Receive a rule to call a vote on
func (l *legislature) setRuleToVote(ruleMatrix rules.RuleMatrix) error {
	if !CheckEnoughInCommonPool(l.gameConf.SetRuleToVoteActionCost, l.gameState) {
		return errors.Errorf("Insufficient Budget in common Pool: setRuleToVote")
	}

	agendaReturn := l.clientSpeaker.DecideAgenda(ruleMatrix)
	if agendaReturn.ActionTaken && agendaReturn.ContentType == shared.SpeakerAgenda {
		if !l.incurServiceCharge(l.gameConf.SetRuleToVoteActionCost) {
			return errors.Errorf("Insufficient Budget in common Pool: setRuleToVote")
		}
		l.ruleToVote = agendaReturn.RuleMatrix
	}
	return nil
}

//Asks islands to vote on a rule
//Called by orchestration
func (l *legislature) setVotingResult(clientIDs []shared.ClientID) (bool, error) {
	voteCalled := false
	if !CheckEnoughInCommonPool(l.gameConf.SetVotingResultActionCost, l.gameState) {
		return voteCalled, errors.Errorf("Insufficient Budget in common Pool: announceVotingResult")
	}

	returnVote := l.clientSpeaker.DecideVote(l.ruleToVote, clientIDs)
	if returnVote.ActionTaken && returnVote.ContentType == shared.SpeakerVote {
		if !l.incurServiceCharge(l.gameConf.SetVotingResultActionCost) {
			return voteCalled, errors.Errorf("Insufficient Budget in common Pool: setVotingResult")
		}
		l.ballotBox = l.RunVote(returnVote.RuleMatrix, returnVote.ParticipatingIslands)

		l.votingResult = l.ballotBox.CountVotesMajority()
		voteCalled = true
	}
	return voteCalled, nil
}

//RunVote creates the voting object, returns votes by category (for, against) in BallotBox.
//Passing in empty ruleID or empty clientIDs results in no vote occurring
func (l *legislature) RunVote(ruleMatrix rules.RuleMatrix, clientIDs []shared.ClientID) voting.BallotBox {

	if ruleMatrix.RuleMatrixIsEmpty() || len(clientIDs) == 0 {
		return voting.BallotBox{}
	}

	ruleVote := voting.RuleVote{}

	//TODO: check if rule is valid, otherwise return empty ballot, raise error?
	ruleVote.SetRule(ruleMatrix)

	//TODO: intersection of islands alive and islands chosen to vote in case of client error
	//TODO: check if remaining slice is >0, otherwise return empty ballot, raise error?
	ruleVote.SetVotingIslands(clientIDs)

	ruleVote.GatherBallots(iigoClients)
	//TODO: log of vote occurring with ruleMatrix, clientIDs
	//TODO: log of clientIDs vs islandsAllowedToVote
	//TODO: log of ruleMatrix vs s.RuleToVote

	variablesToCache := []rules.VariableFieldName{rules.IslandsAllowedToVote}
	valuesToCache := [][]float64{{float64(len(clientIDs))}}
	l.monitoring.addToCache(l.SpeakerID, variablesToCache, valuesToCache)

	rulesEqual := false
	if ruleID == l.ruleToVote {
		rulesEqual = true
	}

	variablesToCache = []rules.VariableFieldName{rules.SpeakerProposedPresidentRule}
	valuesToCache = [][]float64{{boolToFloat(rulesEqual)}}
	l.monitoring.addToCache(l.SpeakerID, variablesToCache, valuesToCache)

	return ruleVote.GetBallotBox()
}

//Speaker declares a result of a vote (see spec to see conditions on what this means for a rule-abiding speaker)
//Called by orchestration
func (l *legislature) announceVotingResult() (bool, error) {
	resultAnnounced := false
	if !CheckEnoughInCommonPool(l.gameConf.AnnounceVotingResultActionCost, l.gameState) {
		return resultAnnounced, errors.Errorf("Insufficient Budget in common Pool: announceVotingResult")
	}

	returnAnouncement := l.clientSpeaker.DecideAnnouncement(l.ruleToVote, l.votingResult)

	if returnAnouncement.ActionTaken && returnAnouncement.ContentType == shared.SpeakerAnnouncement {
		//Deduct action cost
		if !l.incurServiceCharge(l.gameConf.AnnounceVotingResultActionCost) {
			return resultAnnounced, errors.Errorf("Insufficient Budget in common Pool: announceVotingResult")
		}

		//Perform announcement
		broadcastToAllIslands(shared.TeamIDs[l.SpeakerID], generateVotingResultMessage(returnAnouncement.RuleMatrix, returnAnouncement.VotingResult))
		resultAnnounced = true
	}
	return resultAnnounced, nil
}

func generateVotingResultMessage(ruleMatrix rules.RuleMatrix, result bool) map[shared.CommunicationFieldName]shared.CommunicationContent {
	returnMap := map[shared.CommunicationFieldName]shared.CommunicationContent{}

	returnMap[shared.RuleName] = shared.CommunicationContent{
		T:              shared.CommunicationString,
		RuleMatrixData: ruleMatrix,
	}
	returnMap[shared.RuleVoteResult] = shared.CommunicationContent{
		T:           shared.CommunicationBool,
		BooleanData: result,
	}

	return returnMap
}


// updateRules updates the rules in play according to the result of a vote.
func (l *legislature) updateRules(ruleMatrix rules.RuleMatrix, ruleIsVotedIn bool) error {
	if !l.incurServiceCharge(l.gameConf.UpdateRulesActionCost) {
		return errors.Errorf("Insufficient Budget in common Pool: updateRules")
	}
	//TODO: might want to log the errors as logging messages too?
	//notInRulesCache := errors.Errorf("Rule '%v' is not available in rules cache", ruleMatrix)
	if _, ok := rules.AvailableRules[ruleMatrix.RuleName]; !ok || reflect.DeepEqual(ruleMatrix, rules.AvailableRules[ruleMatrix.RuleName]) { //if the proposed ruleMatrix has the same content as the rule with the same name in AvailableRules, the proposal is for putting a rule in/out of play.
		if ruleIsVotedIn {
			err := rules.PullRuleIntoPlay(ruleMatrix.RuleName)
			if ruleErr, ok := err.(*rules.RuleError); ok {
				if ruleErr.Type() == rules.RuleNotInAvailableRulesCache {
					return ruleErr
				}
			}
		} else {
			err := rules.PullRuleOutOfPlay(ruleMatrix.RuleName)
			if ruleErr, ok := err.(*rules.RuleError); ok {
				if ruleErr.Type() == rules.RuleNotInAvailableRulesCache {
					return ruleErr
				}
			}

		}
	} else { //if the proposed ruleMatrix has different content to the rule with the same name in AvailableRules, the proposal is for modifying the rule in the rule caches. It doesn't put a rule in/out of play.
		if ruleIsVotedIn {
			err := rules.ModifyRule(ruleMatrix.RuleName, ruleMatrix.ApplicableMatrix, ruleMatrix.AuxiliaryVector)
			return err
		}
	}

	return nil

}

// appointNextJudge returns the island ID of the island appointed to be Judge in the next turn
func (l *legislature) appointNextJudge(monitoring shared.MonitorResult, currentJudge shared.ClientID, allIslands []shared.ClientID) (shared.ClientID, error) {
	var election voting.Election
	var nextJudge shared.ClientID
	electionsettings := l.clientSpeaker.CallJudgeElection(monitoring, l.judgeTurnsInPower, allIslands)
	if electionsettings.HoldElection {
		if !l.incurServiceCharge(l.gameConf.AppointNextJudgeActionCost) {
			return l.gameState.JudgeID, errors.Errorf("Insufficient Budget in common Pool: appointNextJudge")
		}
		election.ProposeElection(shared.Judge, electionsettings.VotingMethod)
		election.OpenBallot(electionsettings.IslandsToVote, iigoClients)
		election.Vote(iigoClients)
		l.judgeTurnsInPower = 0
		nextJudge = election.CloseBallot(iigoClients)
		nextJudge = l.clientSpeaker.DecideNextJudge(nextJudge)
	} else {
		l.judgeTurnsInPower++
		nextJudge = currentJudge
	}
	return nextJudge, nil
}

func (l *legislature) incurServiceCharge(cost shared.Resources) bool {
	_, ok := WithdrawFromCommonPool(cost, l.gameState)
	if ok {
		l.gameState.IIGORolesBudget[shared.Speaker] -= cost
	}
	return ok
}
