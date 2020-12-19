package roles

type BaseSpeaker struct {
	id	int
	budget	int
	judgeSalary int
	ruleToVote	int
}

func (s *BaseSpeaker) WithdrawJudgeSalary() {

}

func (s *BaseSpeaker) PayJudge() {

}

// Receive a rule to call a vote on
func (s *BaseSpeaker) SetRuleToVote(r int) {
	s.ruleToVote = r
}

func (s *BaseSpeaker) RunVote() {

	if s.ruleToVote == -1 {
		// No rules were proposed by the islands
		return
	} else{
		//Run the vote
		//Creating the object = opening the ballot?
		v := voting.VoteRule{s.ruleToVote}
		//TODO: updateTurnHistory of rule given to vote on vs ruleToVote

		//Providing optional array of islandID-s = closing ballot early/
		//controlling which islands are allowed to vote if some are not permitted?
		//Providing optional string of result counting method

		//Receive result, number of islands alive, number of islands that voted
		ballots := v.CallVote(s.id)
		//Speaker Id passed in for robust logging

		result := v.CountVotes(votes, "majority")
		s.DeclareResult(result, s.ruleToVote)

	}
}

func (s *BaseSpeaker) DeclareResult(result bool, rule int){

	s.UpdateRules(result, rule)
}

func (s *BaseSpeaker) UpdateRules() {

}

func (s *BaseSpeaker) voteNewJudge() {
	
}
