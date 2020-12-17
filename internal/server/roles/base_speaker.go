package roles

type BaseSpeaker struct {
	id					int
	budget				int
	judgeSalary 		int
	ruleToVote			int
}

func (s *BaseSpeaker) WithdrawJudgeSalary() {

}

func (s *BaseSpeaker) PayJudge() {

}

func (s *BaseSpeaker) RunVote() {

}

func (s *BaseSpeaker) UpdateRules() {

}
