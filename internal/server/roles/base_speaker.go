package roles

type baseSpeaker struct {
	id          int
	budget      int
	judgeSalary int
	ruleToVote  int
}

func (s *baseSpeaker) WithdrawJudgeSalary() {

}

func (s *baseSpeaker) PayJudge() {

}

func (s *baseSpeaker) RunVote() {

}

func (s *baseSpeaker) UpdateRules() {

}

func (s *baseSpeaker) voteNewJudge() {

}
