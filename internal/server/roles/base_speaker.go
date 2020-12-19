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

<<<<<<< HEAD
func (s* BaseSpeaker) DeclareResult() {

}

func (s *BaseSpeaker) UpdateRules() {
=======
func (s *baseSpeaker) UpdateRules() {
>>>>>>> origin/orchestration

}

func (s *baseSpeaker) voteNewJudge() {

}
