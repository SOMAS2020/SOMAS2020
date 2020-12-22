package voting

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type Election struct {
	roleToElect		baseclient.Role
	votingMethod	ElectionVotingMethod
	islandsToVote []shared.ClientID
	votes         [][]shared.ClientID
}

type ElectionVotingMethod = int

const (
	BordaCount
	Plurality
	Majority
)

// ProposeMotion sets the role to be voted on
func (e *Election) ProposeElection(role baseclient.Role, method ElectionVotingMethod) {
	e.roleToElect = role
	e.votingMethod = method
}

// OpenBallot sets the islands eligible to vote.
func (e *Election) OpenBallot(clientIDs []shared.ClientID) {
	e.islandsToVote = clientIDs
}

// Vote gets votes from eligible islands.
func (e *Election) Vote(clientMap map[shared.ClientID]baseclient.Client) {
	for _, island := range v.islandsToVote {
		e.votes = append(v.votes, clientMap[island].GetVoteForElection(v.roleToElect))
	}
}

// CloseBallot counts the votes received and returns the result.
func (e *Election) CloseBallot() shared.ClientID {

	switch e.votingMethod {
	case BordaCount :
		result = bordaCountResult()
	case Plurality :
		result = pluralityResult()
	case Majority :
		result = majorityResult()
	}
	return result
}

func (e* Election) bordaCountResult() shared.ClientID {

}

func (e* Election) pluralityResult() shared.ClientID {
	
}

func (e* Election) majorityResult() shared.ClientID {
	
}