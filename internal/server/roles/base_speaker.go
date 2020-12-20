package roles

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/pkg/errors"
)

type baseSpeaker struct {
	id            int
	budget        int
	judgeSalary   int
	ruleToVote    int
	clientSpeaker Speaker
}

func (s *baseSpeaker) WithdrawJudgeSalary() {

}

func (s *baseSpeaker) PayJudge() {

}

func (s *baseSpeaker) RunVote() {

}

func (s *baseSpeaker) DeclareResult() {

}

func (s *baseSpeaker) updateRules(ruleName string, ruleVotedIn bool) error {
	//TODO: might want to log the errors as normal messages rather than completely ignoring them? But then Speaker needs access to client's logger
	notInRulesCache := errors.Errorf("Rule '%v' not available in rules cache", ruleName)
	if ruleVotedIn {
		err := rules.PullRuleIntoPlay(ruleName)
		if err.Error() == notInRulesCache.Error() {
			return err
		}
	} else {
		err := rules.PullRuleOutOfPlay(ruleName)
		if err.Error() == notInRulesCache.Error() {
			return err
		}
	}
	return nil
}

func (s *baseSpeaker) voteNewJudge() {

}
