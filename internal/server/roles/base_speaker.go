package roles

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/pkg/errors"
	"math/rand"
)

type baseSpeaker struct {
	id            int
	budget        int
	judgeSalary   int
	ruleToVote    string
	votingResult  bool
	clientSpeaker Speaker
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
		amount, err := s.clientSpeaker.payJudge()
		if err == nil {
			featureJudge.budget = amount
			return
		}
	}
	amount, _ := s.payJudge()
	featureJudge.budget = amount
}

// Pay the judge
func (s *baseSpeaker) payJudge() (int, error) {
	hold := s.judgeSalary
	s.judgeSalary = 0
	return hold, nil
}

// Receive a rule to call a vote on
func (s *baseSpeaker) SetRuleToVote(r string) {
	s.ruleToVote = r
}

//Asks islands to vote on a rule
//Called by orchestration
func (s *baseSpeaker) setVotingResult() {
	if s.clientSpeaker != nil {
		result, err := s.clientSpeaker.RunVote(s.ruleToVote)
		if err != nil {
			s.votingResult, _ = s.RunVote(s.ruleToVote)
		} else {
			s.votingResult = result
		}
	} else {
		s.votingResult, _ = s.RunVote(s.ruleToVote)
	}
}

//Creates the voting object, collect ballots & count the votes
//Functional so it corresponds to the interface, to the client implementation
func (s *baseSpeaker) RunVote(ruleID string) (bool, error) {
	s.budget -= 10
	if ruleID == "" {
		// No rules were proposed by the islands
		return false, nil
	} else {
		////Run the vote
		////TODO: updateTurnHistory of rule given to vote on vs , so need to pass in
		//v := voting.VoteRule{s.ruleToVote}
		//
		////Receive ballots
		////Speaker Id passed in for logging
		////TODO:
		//ballots := v.CallVote(s.id)
		//
		////TODO:
		//return v.CountVotes(ballots, "majority")

		//For testing while voting is not finished
		return true, nil
	}
}

//Speaker declares a result of a vote (see spec to see conditions on what this means for a rule-abiding speaker)
//Called by orchestration
func (s *baseSpeaker) announceVotingResult() {
	s.budget -= 10
	rule := ""
	result := false
	err := error(nil)

	if s.clientSpeaker != nil {
		//Power to change what is declared completely
		rule, result, err = s.clientSpeaker.DecideAnnouncement(s.ruleToVote, s.votingResult)
		//TODO: log of given vs. returned rule and result
		if err != nil {
			rule, result, _ = s.DecideAnnouncement(s.ruleToVote, s.votingResult)
		}
	} else {
		rule, result, _ = s.DecideAnnouncement(s.ruleToVote, s.votingResult)
	}

	broadcastToAllIslands(s.id, generateVotingResultMessage(rule, result))
	_ = s.updateRules(s.ruleToVote, s.votingResult)

	//Reset
	s.ruleToVote = ""
	s.votingResult = false
}

//Example of the client implementation of DecideAnnouncement
//A well behaved speaker announces what had been voted on and the corresponding result
func (s *baseSpeaker) DecideAnnouncement(ruleId string, result bool) (string, bool, error) {
	return ruleId, result, nil
}

func generateVotingResultMessage(ruleID string, result bool) map[int]DataPacket {
	returnMap := map[int]DataPacket{}

	returnMap[RuleName] = DataPacket{
		textData: ruleID,
	}
	returnMap[RuleVoteResult] = DataPacket{
		booleanData: result,
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

func (s *baseSpeaker) appointNextJudge() int {
	return rand.Intn(5)
}
