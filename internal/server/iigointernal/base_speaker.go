package iigointernal

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"math/rand"
	"time"

	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/voting"
	"github.com/pkg/errors"
)

type baseSpeaker struct {
	Id            int
	budget        int
	judgeSalary   int
	RuleToVote    string
	VotingResult  bool
	clientSpeaker roles.Speaker
}

func (s *baseSpeaker) withdrawJudgeSalary(gameState *gamestate.GameState) error {
	var judgeSalary = int(rules.VariableMap["judgeSalary"].Values[0])
	var withdrawError = WithdrawFromCommonPool(judgeSalary, gameState)
	if withdrawError != nil {
		featureSpeaker.judgeSalary = judgeSalary
	}
	return withdrawError
}

func (s *baseSpeaker) sendJudgeSalary() {
	if s.clientSpeaker != nil {
		amount, err := s.clientSpeaker.PayJudge()
		if err == nil {
			featureJudge.budget = amount
			return
		}
	}
	amount, _ := s.PayJudge()
	featureJudge.budget = amount
}

// Pay the judge
func (s *baseSpeaker) PayJudge() (int, error) {
	hold := s.judgeSalary
	s.judgeSalary = 0
	return hold, nil
}

// Receive a rule to call a vote on
func (s *baseSpeaker) setRuleToVote(r string) {
	s.RuleToVote = r
}

//Asks islands to vote on a rule
//Called by orchestration
func (s *baseSpeaker) setVotingResult(iigoClients map[shared.ClientID]baseclient.Client) {
	//if s.clientSpeaker != nil {
	//	result, err := s.clientSpeaker.RunVote(s.RuleToVote)
	//	if err != nil {
	//		s.VotingResult, _ = s.RunVote(s.RuleToVote)
	//	} else {
	//		s.VotingResult = result
	//	}
	//} else {
	//	s.VotingResult, _ = s.RunVote(s.RuleToVote)
	//}

	//Speaker gets no choice, voting just happens implementation
	//this is not MVP, but it works
	//TODO: Separate and clearly define speaker held information and vote held information (see ruleToVote)
	//TODO: Remove tests
	ruleVote := voting.RuleVote{}
	ruleVote.ProposeMotion(s.RuleToVote)

	//TODO: for loop should not be done here
	var clientIDs []shared.ClientID
	for id := range getIslandAlive() {
		clientIDs = append(clientIDs, shared.ClientID(id))
	}
	ruleVote.OpenBallot(clientIDs)
	ruleVote.Vote(iigoClients)
	s.VotingResult = ruleVote.CloseBallot().Result

}

//Deprecated
//Creates the voting object, collect ballots & count the votes
//Functional so it corresponds to the interface, to the client implementation
//If agent decides not to use voting functions, it is assumed they have not performed them
func (s *baseSpeaker) RunVote(ruleID string) (bool, error) {
	s.budget -= 10
	if ruleID == "" {
		// No rules were proposed by the islands
		return false, nil
	} else {
		////TODO: updateTurnHistory of rule-given-to-vote vs RuleToVote
		//TODO: pass in islandID for log
		//ballotsFor, ballotsAgainst, result = voting.VoteRule(ruleID, getIslandAlive())

		//Return a random result for now
		rand.Seed(time.Now().UnixNano())
		return rand.Int31()&0x01 == 0, nil
	}
}

//Speaker declares a result of a vote (see spec to see conditions on what this means for a rule-abiding speaker)
//Called by orchestration
func (s *baseSpeaker) announceVotingResult() error {

	var rule string
	var result bool
	var err error

	if s.clientSpeaker != nil {
		//Power to change what is declared completely, return "", _ for no announcement to occur
		rule, result, err = s.clientSpeaker.DecideAnnouncement(s.RuleToVote, s.VotingResult)
		//TODO: log of given vs. returned rule and result
		if err != nil {
			rule, result, _ = s.DecideAnnouncement(s.RuleToVote, s.VotingResult)
		}
	} else {
		rule, result, _ = s.DecideAnnouncement(s.RuleToVote, s.VotingResult)
	}

	if rule != "" {
		//Deduct action cost
		s.budget -= 10

		//Reset
		s.RuleToVote = ""
		s.VotingResult = false

		//Perform announcement
		broadcastToAllIslands(shared.TeamIDs[s.Id], generateVotingResultMessage(rule, result))
		return s.updateRules(rule, result)
	}
	return nil
}

//Example of the client implementation of DecideAnnouncement
//A well behaved speaker announces what had been voted on and the corresponding result
//Return "", _ for no announcement to occur
func (s *baseSpeaker) DecideAnnouncement(ruleId string, result bool) (string, bool, error) {
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

func (s *baseSpeaker) updateRules(ruleName string, ruleVotedIn bool) error {
	s.budget -= 10
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

func (s *baseSpeaker) appointNextJudge(clientIDs []shared.ClientID) int {
	s.budget -= 10
	var election voting.Election
	election.ProposeElection(baseclient.Judge, voting.Plurality)
	election.OpenBallot(clientIDs)
	election.Vote(iigoClients)
	return election.CloseBallot()
}
