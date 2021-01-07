package team6

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// GetVoteForRule returns the client's vote in favour of or against a rule.
func (c *client) VoteForRule(ruleMatrix rules.RuleMatrix) shared.RuleVoteType {
	// TODO implement decision on voting that considers the rule
	return shared.Abstain
}

// GetVoteForElection returns the client's Borda vote for the role to be elected.
// COMPULSORY: use opinion formation to decide a rank for islands for the role
func (c *client) VoteForElection(roleToElect shared.Role, candidateList []shared.ClientID) []shared.ClientID {
	candidates := map[int]shared.ClientID{}
	for i := 0; i < len(candidateList); i++ {
		candidates[i] = candidateList[i]
	}
	// Recombine map, in shuffled order
	var returnList []shared.ClientID
	for _, v := range candidates {
		returnList = append(returnList, v)
	}
	return returnList
}
