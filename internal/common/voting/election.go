package voting

type Election struct {

}

// ProposeMotion sets the role to be voted on
func (e *Election) ProposeMotion() {

}

// OpenBallot sets the islands eligible to vote.
func (e *Election) OpenBallot(clientIDs []shared.ClientID){
	s.GetVoteForRule(ruleToVote)
}

// Vote gets votes from eligible islands.
func (e *Election) Vote() {

}

// CloseBallot counts the votes received and returns the result.
func (e *Election) CloseBallot() {

}