package voting

type Election struct {
	islandsToVote 	[]shared.ClientID
	votes 			[]bool
}

type ElectionResult struct {

}

// ProposeMotion sets the role to be voted on
func (e *Election) ProposeMotion() {

}

// OpenBallot sets the islands eligible to vote.
func (e *Election) OpenBallot(clientIDs []shared.ClientID){
	islandsToVote := clientIDs
}

// Vote gets votes from eligible islands.
func (e *Election) Vote() {
	s.GetVoteForElection(ruleToVote)
}

// CloseBallot counts the votes received and returns the result.
func (e *Election) CloseBallot() {

}