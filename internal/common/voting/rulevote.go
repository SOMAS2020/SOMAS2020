package voting

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type RuleVote struct {
	//Checked by RuleVote
	ruleToVote string
	voterList  []shared.ClientID
	//Held by RuleVote
	ballots map[int][]bool
}

type BallotBox struct {
	VotesInFavour uint
	VotesAgainst  uint
}

// SetRule is called by baseSpeaker to set the rule to be voted on.
func (v *RuleVote) SetRule(rule string) {
	v.ruleToVote = rule
}

// SetVotingIslands is called by baseSpeaker to set the islands eligible to vote.
func (v *RuleVote) SetVotingIslands(clientIDs []shared.ClientID) {
	//TODO: intersection of islands alive and islands chosen to vote
	v.voterList = clientIDs
}

// GatherBallots is called by baseSpeaker to get votes from clients.
func (v *RuleVote) GatherBallots(clientMap map[shared.ClientID]baseclient.Client) {
	//Gather N ballots from islands
	oneIslandVote := make([]bool, 2)
	if v.ruleToVote != "" && len(v.voterList) > 0 {
		for i := 0; i < len(v.voterList); i++ {
			oneIslandVote[0], oneIslandVote[1] = clientMap[v.voterList[i]].VoteForRule(v.ruleToVote)
			v.ballots[i] = oneIslandVote
		}
	}
}

//GetBallotBox is called by baseSpeaker and
//returns the BallotBox with n votesInFavour and N-n votesAgainst
func (v *RuleVote) GetBallotBox() BallotBox {
	//The following is in accordance with anonymous voting
	//Abstentions will not be considered(vote[1]==true)
	var outcome BallotBox
	for _, vote := range v.ballots {
		if vote[1] == false && vote[0] == true {
			outcome.VotesInFavour += 1
		} else if vote[1] == false && vote[0] == false {
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
