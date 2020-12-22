package voting

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type Election struct {
	roleToElect	baseclient.Role
	islandsToVote []shared.ClientID
	votes         [][]shared.ClientID
}

// ProposeMotion sets the role to be voted on
func (e *Election) ProposeElection(role baseclient.Role) {
	roleToElect = role
}

// OpenBallot sets the islands eligible to vote.
func (e *Election) OpenBallot(clientIDs []shared.ClientID) {
	e.islandsToVote = clientIDs
}

// Vote gets votes from eligible islands.
func (e *Election) Vote(clientMap map[shared.ClientID]baseclient.Client) {
	for _, island := range v.islandsToVote {
		v.votes = append(v.votes, clientMap[island].GetVoteForElection(v.roleToElect))
	}
}

// CloseBallot counts the votes received and returns the result.
func (e *Election) CloseBallot() shared.ClientID {
	
}
