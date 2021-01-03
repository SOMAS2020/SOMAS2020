package team5

import (
	"math/rand"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func (c *client) GetVoteForElection(roleToElect shared.Role) []shared.ClientID {
	// No dynamic strategy right now. Team 5 always get the first vote. 
	var vote []shared.ClientID
	var initialVote []shared.ClientID 
	vote = append(vote, shared.Team5) // Team 5 alaways first
	initialVote = append(initialVote, shared.Team1)
	initialVote = append(initialVote, shared.Team2)
	initialVote = append(initialVote, shared.Team4)
	initialVote = append(initialVote, shared.Team6)
	// shuffle those teams
	rand.Shuffle(len(initialVote), func(i, j int) { initialVote[i], initialVote[j] = initialVote[j], initialVote[i] })
	initialVote = append(initialVote, shared.Team3)
	vote = append(vote, initialVote...)
	c.Logf("Team 5 vote: %v", vote)
	return vote
}