package roles

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
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

func (s *baseSpeaker) updateRules(ruleName string, ruleVotedIn bool) {
	//TODO: check with Neelesh: maybe PullRuleInto and OutOf play shouldn't return error when a rule is already in/out?
	if ruleVotedIn {
		rules.PullRuleIntoPlay(ruleName)
	} else {
		rules.PullRuleOutOfPlay(ruleName)
	}
}

func (s *baseSpeaker) voteNewJudge() {

}
