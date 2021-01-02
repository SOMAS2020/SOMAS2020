package iigointernal

import (
	"fmt"

	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"github.com/SOMAS2020/SOMAS2020/internal/common/voting"
)

type legislature struct {
	SpeakerID         shared.ClientID
	budget            shared.Resources
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

// returnJudgeSalary returns the salary to the common pool.
func (l *legislature) returnJudgeSalary() shared.Resources {
	x := l.judgeSalary
	l.judgeSalary = 0
	return x
}

// withdrawJudgeSalary withdraws the salary for the Judge from the common pool.
func (l *legislature) withdrawJudgeSalary(gameState *gamestate.GameState) bool {
	var judgeSalary = shared.Resources(rules.VariableMap[rules.JudgeSalary].Values[0])
	withdrawAmount, withdrawSuccesful := WithdrawFromCommonPool(judgeSalary, gameState)
	l.judgeSalary = withdrawAmount

	return withdrawSuccesful
}

// sendJudgeSalary sets the budget of the Judge.
func (l *legislature) sendJudgeSalary(judicialBranch *judiciary) {
	if l.clientSpeaker != nil {
		amount, judgePaid := l.clientSpeaker.PayJudge(l.judgeSalary)
		if judgePaid {
			judicialBranch.budget = amount
		}
		return
	}
	judicialBranch.budget = l.judgeSalary
}

// Receive a rule to call a vote on
func (l *legislature) setRuleToVote(r string) {
	ruleToBeVoted, ruleSet := l.clientSpeaker.DecideAgenda(r)
	if ruleSet {
		l.ruleToVote = ruleToBeVoted
	}
}

//Asks islands to vote on a rule
//Called by orchestration
func (l *legislature) setVotingResult(clientIDs []shared.ClientID) bool {

	ruleID, participatingIslands, voteDecided := l.clientSpeaker.DecideVote(l.ruleToVote, clientIDs)
	if !voteDecided {
		return false
	}
	l.ballotBox = l.RunVote(ruleID, participatingIslands)
	l.votingResult = l.ballotBox.CountVotesMajority()
	return true
}

//RunVote creates the voting object, returns votes by category (for, against) in BallotBox.
//Passing in empty ruleID or empty clientIDs results in no vote occurring
func (l *legislature) RunVote(ruleID string, clientIDs []shared.ClientID) voting.BallotBox {

	if ruleID == "" || len(clientIDs) == 0 {
		return voting.BallotBox{}
	}
	l.budget -= serviceCharge
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
func (l *legislature) announceVotingResult() {

	rule, result, announcementDecided := l.clientSpeaker.DecideAnnouncement(l.ruleToVote, l.votingResult)

	if announcementDecided {
		//Deduct action cost
		l.budget -= serviceCharge

		//Reset
		l.ruleToVote = ""
		l.votingResult = false

		//Perform announcement
		broadcastToAllIslands(shared.TeamIDs[l.SpeakerID], generateVotingResultMessage(rule, result))
	}
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
	l.budget -= serviceCharge
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
func (l *legislature) appointNextJudge(monitoring shared.MonitorResult, currentJudge shared.ClientID, allIslands []shared.ClientID) shared.ClientID {
	var election voting.Election
	var nextJudge shared.ClientID
	electionsettings := l.clientSpeaker.CallJudgeElection(monitoring, l.judgeTurnsInPower, allIslands)
	if electionsettings.HoldElection {
		// TODO: deduct the cost of holding an election
		election.ProposeElection(shared.Judge, electionsettings.VotingMethod)
		election.OpenBallot(electionsettings.IslandsToVote)
		election.Vote(iigoClients)
		l.judgeTurnsInPower = 0
		nextJudge = election.CloseBallot()
		nextJudge = l.clientSpeaker.DecideNextJudge(nextJudge)
	} else {
		l.judgeTurnsInPower++
		nextJudge = currentJudge
	}
	return nextJudge
}
