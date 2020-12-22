package voting

type RoleVote struct {

}

// ProposeMotion sets the role to be voted on
func (v *RoleVote) ProposeMotion() {

}

// OpenBallot sets the islands eligible to vote.
func (v *RoleVote) OpenBallot(clientIDs []shared.ClientID){
	s.GetVoteForRule(ruleToVote)
}

// Vote gets votes from eligible islands.
func (v *RoleVote) Vote() {

}

// CloseBallot counts the votes received and returns the result.
func (v *RoleVote) CloseBallot() {

}