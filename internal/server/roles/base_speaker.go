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

func (s *BaseSpeaker) RunVote() error {

	if s.ruleToVote == -1 {
		// No rules were proposed by the islands
		return nil
	} else{
		//Run the vote
		//Creating the object = opening the ballot?
		v := voting.Vote{s.ruleToVote}
		//TODO: updateTurnHistory of rule given to vote on vs ruleToVote

		//Providing optional array of islandID-s = closing ballot early/
		//controlling which islands are allowed to vote if some are not permitted?
		//Providing optional string of result counting method

		//Receive result, number of islands alive, number of islands that voted
		votes, nia, niv := voting.CallVote(s.id)
		//Speaker Id passed in for robust logging

		//Log result TODO: Can this log be moved to inside voting?
		noIslandAlive := rules.VariableValuePair{
			VariableName: "no_islands_alive",
			Values:       []float64{nia},
		}
		noIslandsVoting := rules.VariableValuePair{
			VariableName: "no_islands_voted",
			Values:       []float64{niv},
		}

		err := updateTurnHistory(s.id, []rules.VariableValuePair{noIslandAlive, noIslandsVoting})

		if err != nil {
			return err
		} else {
			//Count the votes
			result := voting.CountVotes(votes, "majority")

			//Declare the result = UpdateRules?
			s.UpdateRules(s.ruleToVote, result)
		}

	}
}

func (s *BaseSpeaker) UpdateRules() {

}

func (s *BaseSpeaker) voteNewJudge() {
	
}
