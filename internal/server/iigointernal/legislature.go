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

type legislature struct {
	SpeakerID     shared.ClientID
	budget        shared.Resources
	judgeSalary   shared.Resources
	ruleToVote    string
	ballotBox     voting.BallotBox
	votingResult  bool
	clientSpeaker roles.Speaker
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
	amount, judgePaid := l.clientSpeaker.PayJudge()
	if judgePaid {
		judicialBranch.budget = amount
	}
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
func (l *legislature) setVotingResult(clientIDs []shared.ClientID) {

	ruleID, participatingIslands, voteDecided := l.clientSpeaker.DecideVote(l.ruleToVote, clientIDs)
	if !voteDecided {
		return
	}
	l.ballotBox = l.RunVote(ruleID, participatingIslands)
	l.votingResult = l.ballotBox.CountVotesMajority()

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

func generateVotingResultMessage(ruleID string, result bool) map[shared.CommunicationFieldName]shared.Communication {
	returnMap := map[shared.CommunicationFieldName]shared.Communication{}

	returnMap[shared.RuleName] = shared.Communication{
		T:        shared.CommunicationString,
		TextData: ruleID,
	}
	returnMap[shared.RuleVoteResult] = shared.Communication{
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
	//TODO: might want to log the errors as normal messages rather than completely ignoring them? But then Speaker needs access to client's logger
	notInRulesCache := errors.Errorf("Rule '%v' is not available in rules cache", ruleName)
	if ruleVotedIn {
		// _ = rules.PullRuleIntoPlay(ruleName)
		err := rules.PullRuleIntoPlay(ruleName)
		if err != nil {
			if err.Error() == notInRulesCache.Error() {
				return err
			}
		}
	} else {
		// _ = rules.PullRuleOutOfPlay(ruleName)
		err := rules.PullRuleOutOfPlay(ruleName)
		if err != nil {
			if err.Error() == notInRulesCache.Error() {
				return err
			}
		}

	}
	return nil

}

func (l *legislature) appointNextJudge(clientIDs []shared.ClientID) shared.ClientID {
	l.budget -= serviceCharge
	var election voting.Election
	election.ProposeElection(baseclient.Judge, voting.Plurality)
	election.OpenBallot(clientIDs)
	election.Vote(iigoClients)
	return election.CloseBallot()
}
