package voting

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type RuleVote struct {
	ruleToVote    string
	islandsToVote []shared.ClientID
	votes         []bool
}

type RuleVoteResult struct {
	votesInFavour uint
	votesAgainst  uint
	result        bool
}

// ProposeMotion is called by Speaker to set the rule to be voted on.
func (v *RuleVote) ProposeMotion(rule string) {
	v.ruleToVote = rule
}

// OpenBallot is called by Speaker to set the islands eligible to vote.
func (v *RuleVote) OpenBallot(clientIDs []shared.ClientID) {
	v.islandsToVote = clientIDs
}

// Vote is called by Speaker to get votes from clients.
func (v *RuleVote) Vote(clientMap map[shared.ClientID]baseclient.Client) {
	for _, island := range v.islandsToVote {
		// TODO how the hell do we get clientMap in here
		v.votes = append(v.votes, clientMap[island].GetVoteForRule(v.ruleToVote))
	}
}

// CloseBallot is called by Speaker to count votes received.
func (v *RuleVote) CloseBallot() RuleVoteResult {

	var outcome RuleVoteResult
	for _, vote := range v.votes {
		if vote == true {
			outcome.votesInFavour += 1
		} else if vote == false {
			outcome.votesAgainst += 1
		}
	}
	outcome.result = outcome.votesInFavour >= outcome.votesAgainst
	return outcome
}
