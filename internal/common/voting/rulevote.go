package voting

import (
	"fmt"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type RuleVote struct {
	//Checked by RuleVote
	ruleToVote rules.RuleMatrix
	voterList  []shared.ClientID
	//Held by RuleVote
	ballots []shared.RuleVoteType
	Logger  shared.Logger
}

type BallotBox struct {
	VotesInFavour uint
	VotesAgainst  uint
}

// Logf is the rule vote logger
func (v *RuleVote) Logf(format string, a ...interface{}) {
	v.Logger("[RULE VOTE]: %v", fmt.Sprintf(format, a...))
}

// SetRule is called by baseSpeaker to set the rule to be voted on.
func (v *RuleVote) SetRule(ruleMatrix rules.RuleMatrix) {
	v.ruleToVote = ruleMatrix
}

// SetVotingIslands is called by baseSpeaker to set the islands eligible to vote.
func (v *RuleVote) SetVotingIslands(clientIDs []shared.ClientID) {
	//TODO: intersection of islands alive and islands chosen to vote
	v.voterList = clientIDs
}

// GatherBallots is called by baseSpeaker to get votes from clients.
func (v *RuleVote) GatherBallots(clientMap map[shared.ClientID]baseclient.Client) {
	//Gather N ballots from islands
	if v.ruleToVote.RuleName != "" && len(v.ruleToVote.RuleName) > 0 {
		for i := 0; i < len(v.voterList); i++ {
			v.ballots = append(v.ballots, clientMap[v.voterList[i]].VoteForRule(v.ruleToVote))
		}
	}
	v.Logf("Votes: %v", v.ballots)
}

//GetBallotBox is called by baseSpeaker and
//returns the BallotBox with n votesInFavour and N-n votesAgainst
func (v *RuleVote) GetBallotBox() BallotBox {
	//The following is in accordance with anonymous voting
	//Abstentions will not be considered(vote[1]==true)
	var outcome BallotBox
	for i := 0; i < len(v.ballots); i++ {
		if v.ballots[i] == shared.Approve {
			outcome.VotesInFavour += 1
		} else if v.ballots[i] == shared.Reject {
			outcome.VotesAgainst += 1
		}
	}
	return outcome
}

//CountVotesMajority is called by baseSpeaker and
//returns the majority result of the BallotBox
func (b *BallotBox) CountVotesMajority() bool {
	return b.VotesInFavour > b.VotesAgainst
}
