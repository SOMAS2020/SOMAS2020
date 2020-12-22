package voting

type RuleVote struct {
	ruleToVote 		string
	islandsToVote 	[]shared.ClientID
	votes 			[]bool
}

type RuleVoteResult struct {
	votesInFavour	uint
	votesAgainst	uint
	result			bool
}

// ProposeMotion is called by Speaker to set the rule to be voted on.
func (v *RuleVote) ProposeMotion(rule string) {
	ruleToVote := rule
}

// OpenBallot is called by Speaker to set the islands eligible to vote.
func (v *RuleVote) OpenBallot( clientIDs []shared.ClientID){
	islandsToVote := clientIDs
}

// Vote is called by Speaker to get votes from clients.
func (v *RuleVote) Vote() {
	for _, island := range islandsToVote {
		// TODO how the hell do we get clientMap in here
		votes = append(votes, clientMap[island].GetVoteForRule(ruleToVote))
	}
}

// CloseBallot is called by Speaker to count votes received.
func (v *RuleVote) CloseBallot() RuleVoteResult {

	var outcome RuleVoteResult
	for _, vote := range votes {
		if vote == true {
			outcome.votesInFavour += 1
		}
		else if vote == false {
			outcome.votesAgainst += 1
		}
	}
	outcome.result = outcome.votesInFavour >= outcome.votesAgainst
	return outcome
}