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

type legislature struct {
	gameState         *gamestate.GameState
	SpeakerID         shared.ClientID
	judgeSalary       shared.Resources
	ruleToVote        string
	ballotBox         voting.BallotBox
	votingResult      bool
	clientSpeaker     roles.Speaker
	judgeTurnsInPower int
}

// loadClientSpeaker checks client pointer is good and if not panics
func (l *legislature) loadClientSpeaker(clientSpeakerPointer roles.Speaker) {
	if clientSpeakerPointer == nil {
		panic(fmt.Sprintf("Client '%v' has loaded a nil speaker pointer", l.SpeakerID))
	}
	l.clientSpeaker = clientSpeakerPointer
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
func (l *legislature) setRuleToVote(r string) error {
	if !CheckEnoughInCommonPool(actionCost.SetRuleToVoteActionCost, l.gameState) {
		return errors.Errorf("Insufficient Budget in common Pool: setRuleToVote")
	}

	agendaReturn := l.clientSpeaker.DecideAgenda(r)
	if agendaReturn.ActionTaken && agendaReturn.ContentType == shared.SpeakerAgenda {
		if !l.incurServiceCharge(actionCost.SetRuleToVoteActionCost) {
			return errors.Errorf("Insufficient Budget in common Pool: setRuleToVote")
		}
		l.ruleToVote = agendaReturn.RuleID
	}
	return nil
}

//Asks islands to vote on a rule
//Called by orchestration
func (l *legislature) setVotingResult(clientIDs []shared.ClientID) error {
	if !CheckEnoughInCommonPool(actionCost.SetVotingResultActionCost, l.gameState) {
		return errors.Errorf("Insufficient Budget in common Pool: announceVotingResult")
	}

	returnVote := l.clientSpeaker.DecideVote(l.ruleToVote, clientIDs)
	if returnVote.ActionTaken && returnVote.ContentType == shared.SpeakerVote {
		if !l.incurServiceCharge(actionCost.SetVotingResultActionCost) {
			return errors.Errorf("Insufficient Budget in common Pool: setVotingResult")
		}
		l.ballotBox = l.RunVote(returnVote.RuleID, returnVote.ParticipatingIslands)

		l.votingResult = l.ballotBox.CountVotesMajority()
	}

	return nil
}

//RunVote creates the voting object, returns votes by category (for, against) in BallotBox.
//Passing in empty ruleID or empty clientIDs results in no vote occurring
func (l *legislature) RunVote(ruleID string, clientIDs []shared.ClientID) voting.BallotBox {

	if ruleID == "" || len(clientIDs) == 0 {
		return voting.BallotBox{}
	}

	ruleVote := voting.RuleVote{}

	//TODO: check if rule is valid, otherwise return empty ballot, raise error?
	ruleVote.SetRule(ruleID)

	//TODO: intersection of islands alive and islands chosen to vote in case of client error
	//TODO: check if remaining slice is >0, otherwise return empty ballot, raise error?
	ruleVote.SetVotingIslands(clientIDs)

	ruleVote.GatherBallots(iigoClients)
	//TODO: log of vote occurring with ruleID, clientIDs
	//TODO: log of clientIDs vs islandsAllowedToVote
	//TODO: log of ruleID vs s.RuleToVote
	return ruleVote.GetBallotBox()
}

//Speaker declares a result of a vote (see spec to see conditions on what this means for a rule-abiding speaker)
//Called by orchestration
func (l *legislature) announceVotingResult() error {
	if !CheckEnoughInCommonPool(actionCost.AnnounceVotingResultActionCost, l.gameState) {
		return errors.Errorf("Insufficient Budget in common Pool: announceVotingResult")
	}

	returnAnouncement := l.clientSpeaker.DecideAnnouncement(l.ruleToVote, l.votingResult)

	if returnAnouncement.ActionTaken && returnAnouncement.ContentType == shared.SpeakerAnnouncement {
		//Deduct action cost
		if !l.incurServiceCharge(actionCost.AnnounceVotingResultActionCost) {
			return errors.Errorf("Insufficient Budget in common Pool: announceVotingResult")
		}

		//Reset
		l.ruleToVote = ""
		l.votingResult = false

		//Perform announcement
		broadcastToAllIslands(shared.TeamIDs[l.SpeakerID], generateVotingResultMessage(returnAnouncement.RuleID, returnAnouncement.VotingResult))
	}
	return nil
}

func generateVotingResultMessage(ruleID string, result bool) map[shared.CommunicationFieldName]shared.CommunicationContent {
	returnMap := map[shared.CommunicationFieldName]shared.CommunicationContent{}

	returnMap[shared.RuleName] = shared.CommunicationContent{
		T:        shared.CommunicationString,
		TextData: ruleID,
	}
	returnMap[shared.RuleVoteResult] = shared.CommunicationContent{
		T:           shared.CommunicationBool,
		BooleanData: result,
	}

	return returnMap
}

//reset resets internal variables for safety
func (l *legislature) reset() {
	l.ruleToVote = ""
	l.ballotBox = voting.BallotBox{}
	l.votingResult = false
}

// updateRules updates the rules in play according to the result of a vote.
func (l *legislature) updateRules(ruleName string, ruleVotedIn bool) error {
	if !l.incurServiceCharge(actionCost.UpdateRulesActionCost) {
		return errors.Errorf("Insufficient Budget in common Pool: updateRules")
	}
	//TODO: might want to log the errors as normal messages rather than completely ignoring them? But then Speaker needs access to client's logger
	//notInRulesCache := errors.Errorf("Rule '%v' is not available in rules cache", ruleName)
	if ruleVotedIn {
		// _ = rules.PullRuleIntoPlay(ruleName)
		err := rules.PullRuleIntoPlay(ruleName)
		if ruleErr, ok := err.(*rules.RuleError); ok {
			if ruleErr.Type() == rules.RuleNotInAvailableRulesCache {
				return ruleErr
			}
		}
	} else {
		// _ = rules.PullRuleOutOfPlay(ruleName)
		err := rules.PullRuleOutOfPlay(ruleName)
		if ruleErr, ok := err.(*rules.RuleError); ok {
			if ruleErr.Type() == rules.RuleNotInAvailableRulesCache {
				return ruleErr
			}
		}

	}
	return nil

}

// appointNextJudge returns the island ID of the island appointed to be Judge in the next turn
func (l *legislature) appointNextJudge(currentJudge shared.ClientID, allIslands []shared.ClientID) (shared.ClientID, error) {
	var election voting.Election
	var nextJudge shared.ClientID
	electionsettings := l.clientSpeaker.CallJudgeElection(l.judgeTurnsInPower, allIslands)
	if electionsettings.HoldElection {
		if !l.incurServiceCharge(actionCost.AppointNextJudgeActionCost) {
			return l.gameState.JudgeID, errors.Errorf("Insufficient Budget in common Pool: appointNextJudge")
		}
		election.ProposeElection(baseclient.President, electionsettings.VotingMethod)
		election.OpenBallot(electionsettings.IslandsToVote)
		election.Vote(iigoClients)
		l.judgeTurnsInPower = 0
		nextJudge = election.CloseBallot()
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
		l.gameState.IIGORolesBudget["speaker"] -= cost
	}
	return ok
}
