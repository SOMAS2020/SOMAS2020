package roles

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
)

type baseSpeaker struct {
	id          int
	budget      int
	judgeSalary int
	ruleToVote  int
}

func (s *baseSpeaker) withdrawJudgeSalary(gameState *common.GameState) error {
	var judgeSalary = int(rules.VariableMap["judgeSalary"].Values[0])
	var withdrawError = WithdrawFromCommonPool(judgeSalary, gameState)
	if withdrawError != nil {
		Base_speaker.judgeSalary = judgeSalary
	}
	return withdrawError
}

// Pay the judge
func (s *baseSpeaker) payJudge() {
	Base_judge.budget = Base_speaker.judgeSalary
	Base_speaker.judgeSalary = 0
}

func (s *baseSpeaker) RunVote() {

}

func (s *baseSpeaker) UpdateRules() {

}

func (s *baseSpeaker) voteNewJudge() {

}
