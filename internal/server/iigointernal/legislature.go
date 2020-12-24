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
	RuleToVote    string
	VotingResult  bool
	clientSpeaker roles.Speaker
}

// returnJudgeSalary returns the salary to the common pool.
func (l *legislature) returnJudgeSalary() shared.Resources {
	x := l.judgeSalary
	l.judgeSalary = 0
	return x
}

func (l *legislature) withdrawJudgeSalary(gameState *gamestate.GameState) error {
	var judgeSalary = shared.Resources(rules.VariableMap["judgeSalary"].Values[0])
	var withdrawError = WithdrawFromCommonPool(judgeSalary, gameState)
	if withdrawError != nil {
		l.judgeSalary = judgeSalary
	}
	return withdrawError
}

func (l *legislature) sendJudgeSalary() {
	amount, _ := l.clientSpeaker.PayJudge()
	judicialBranch.budget = amount
}

// Receive a rule to call a vote on
func (l *legislature) setRuleToVote(r string) {
	l.RuleToVote = r
}

//Asks islands to vote on a rule
//Called by orchestration
func (l *legislature) setVotingResult(iigoClients map[shared.ClientID]baseclient.Client) {

	//TODO: Separate and clearly define speaker held information and vote held information (see ruleToVote)
	//TODO: Remove tests
	ruleVote := voting.RuleVote{}
	ruleVote.ProposeMotion(l.RuleToVote)

	//TODO: for loop should not be done here
	var clientIDs []shared.ClientID
	for id := range getIslandAlive() {
		clientIDs = append(clientIDs, shared.ClientID(id))
	}
	ruleVote.OpenBallot(clientIDs)
	ruleVote.Vote(iigoClients)
	l.VotingResult = ruleVote.CloseBallot().Result

}

//Speaker declares a result of a vote (see spec to see conditions on what this means for a rule-abiding speaker)
//Called by orchestration
func (l *legislature) announceVotingResult() error {

	rule, result, err := l.clientSpeaker.DecideAnnouncement(l.RuleToVote, l.VotingResult)

	if err == nil {
		//Deduct action cost
		l.budget -= 10

		//Reset
		l.RuleToVote = ""
		l.VotingResult = false

		//Perform announcement
		broadcastToAllIslands(shared.TeamIDs[l.SpeakerID], generateVotingResultMessage(rule, result))
		return l.updateRules(rule, result)
	}
	return nil
}

//Example of the client implementation of DecideAnnouncement
//A well behaved speaker announces what had been voted on and the corresponding result
//Return "", _ for no announcement to occur
func (l *legislature) DecideAnnouncement(ruleId string, result bool) (string, bool, error) {
	return ruleId, result, nil
}

func generateVotingResultMessage(ruleID string, result bool) map[int]baseclient.Communication {
	returnMap := map[int]baseclient.Communication{}

	returnMap[RuleName] = baseclient.Communication{
		T:        baseclient.CommunicationString,
		TextData: ruleID,
	}
	returnMap[RuleVoteResult] = baseclient.Communication{
		T:           baseclient.CommunicationBool,
		BooleanData: result,
	}

	return returnMap
}

func (l *legislature) updateRules(ruleName string, ruleVotedIn bool) error {
	l.budget -= 10
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
	l.budget -= 10
	var election voting.Election
	election.ProposeElection(baseclient.Judge, voting.Plurality)
	election.OpenBallot(clientIDs)
	election.Vote(iigoClients)
	return election.CloseBallot()
}
