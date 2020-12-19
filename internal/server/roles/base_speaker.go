package roles

type baseSpeaker struct {
	id          	int
	budget      	int
	judgeSalary 	int
	ruleToVote  	string
	votingResult	bool
	clientSpeaker 	Speaker
}

func (s *baseSpeaker) WithdrawJudgeSalary() {

}

func (s *baseSpeaker) PayJudge() {

}

// Receive a rule to call a vote on
func (s *baseSpeaker) SetRuleToVote(r int) {
	s.ruleToVote = r
}

func (s *baseSpeaker) RunVote() {
	if s.clientSpeaker != nil {
		//TODO:
		result, err := s.clientSpeaker.RunVote(s.ruleToVote)
		if err != nil {
			s.votingResult = s.runVoteInternal()
		} else {
			s.votingResult = result
		}
	} else{
		s.votingResult = s.runVoteInternal()
	}
}

func (s *baseSpeaker) runVoteInternal() bool{
	if s.ruleToVote == -1 {
		// No rules were proposed by the islands
		return false
	} else{
		//Run the vote
		//TODO: updateTurnHistory of rule given to vote on vs , so need to pass in
		v := voting.VoteRule{s.ruleToVote}

		//Receive ballots
		//Speaker Id passed in for logging
		//TODO:
		ballots := v.CallVote(s.id)

		//TODO:
		return v.CountVotes(ballots, "majority")
	}
}

func (s *baseSpeaker) DeclareResult(){
	if s.clientSpeaker != nil {
		//Power to change what is declared completely
		//TODO:
		rule, result, err := s.clientSpeaker.DeclareResult(ruleToVote)
		if err != nil {
			broadcastToAllIslands(s.id, generateVotingResultMessage(s.ruleToVote, s.votingResult))
			//TODO:
			s.UpdateRules(s.ruleToVote, s.votingResult)
		} else {
			broadcastToAllIslands(s.id, generateVotingResultMessage(rule, result))
			//TODO:
			s.UpdateRules(rule, result)
		}
	} else{
		broadcastToAllIslands(s.id, generateVotingResultMessage(s.ruleToVote, s.votingResult))
		//TODO:
		s.UpdateRules(s.ruleToVote, s.votingResult)
	}

	//Reset
	s.ruleToVote = -1
	s.votingResult = false

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

func (s *baseSpeaker) UpdateRules() {

}

func (s *baseSpeaker) voteNewJudge() {

}
