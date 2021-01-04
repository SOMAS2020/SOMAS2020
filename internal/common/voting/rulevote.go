package voting

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type RuleVote struct {
	//Checked by RuleVote
	ruleToVote    rules.RuleMatrix
	islandsToVote []shared.ClientID
	//Held by RuleVote
	ballots []bool
}

type BallotBox struct {
	VotesInFavour uint
	VotesAgainst  uint
}

// SetRule is called by baseSpeaker to set the rule to be voted on.
func (v *RuleVote) SetRule(ruleMatrix rules.RuleMatrix) {
	v.ruleToVote = ruleMatrix
}

// SetVotingIslands is called by baseSpeaker to set the islands eligible to vote.
func (v *RuleVote) SetVotingIslands(clientIDs []shared.ClientID) {
	//TODO: intersection of islands alive and islands chosen to vote
	v.islandsToVote = clientIDs
}

// GatherBallots is called by baseSpeaker to get votes from clients.
func (v *RuleVote) GatherBallots(clientMap map[shared.ClientID]baseclient.Client) {
	//Gather N ballots from islands
	if !v.ruleToVote.RuleMatrixIsEmpty() && len(v.islandsToVote) > 0 {
		for _, island := range v.islandsToVote {
			v.ballots = append(v.ballots, clientMap[island].GetVoteForRule(v.ruleToVote))
		}
	}
}

//GetBallotBox is called by baseSpeaker and
//returns the BallotBox with n votesInFavour and N-n votesAgainst
func (v *RuleVote) GetBallotBox() BallotBox {
	//The following is in accordance with anonymous voting
	var outcome BallotBox
	for _, vote := range v.ballots {
		if vote {
			outcome.VotesInFavour += 1
		} else if !vote {
			outcome.VotesAgainst += 1
		}
	}
	return outcome
}

//CountVotesMajority is called by baseSpeaker and
//returns the majority result of the BallotBox
func (b *BallotBox) CountVotesMajority() bool {
	return b.VotesInFavour >= b.VotesAgainst
}
