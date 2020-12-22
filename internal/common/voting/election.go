package voting

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type Election struct {
	islandsToVote []shared.ClientID
	votes         map[shared.ClientID][]bool
}

type ElectionResult struct {
}

// ProposeMotion sets the role to be voted on
func (e *Election) ProposeMotion() {

}

// OpenBallot sets the islands eligible to vote.
func (e *Election) OpenBallot(clientIDs []shared.ClientID) {
	e.islandsToVote = clientIDs
}

// Vote gets votes from eligible islands.
func (e *Election) Vote(clientMap map[shared.ClientID]baseclient.Client) {
	for i, v := range clientMap {
		e.votes[i] = v.GetVoteForElection(len(e.islandsToVote))
	}
}

// CloseBallot counts the votes received and returns the result.
func (e *Election) CloseBallot() {

}
