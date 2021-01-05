package team6

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// GetVoteForRule returns the client's vote in favour of or against a rule.
func (c *client) GetVoteForRule(ruleName string) bool {
	for _, val := range c.favourRules {
		if val == ruleName {
			return true
		}
	}
	return false
}

// GetVoteForElection returns the client's Borda vote for the role to be elected.
// COMPULSORY: use opinion formation to decide a rank for islands for the role
func (c *client) GetVoteForElection(roleToElect shared.Role) []shared.ClientID {
	// Done ;)
	// Get all alive islands
	aliveClients := rules.VariableMap[rules.IslandsAlive]
	// Convert to ClientID type and place into unordered map
	aliveClientIDs := map[int]shared.ClientID{}
	for i, v := range aliveClients.Values {
		aliveClientIDs[i] = shared.ClientID(int(v))
	}
	// Recombine map, in shuffled order
	var returnList []shared.ClientID
	for _, v := range aliveClientIDs {
		returnList = append(returnList, v)
	}
	return returnList
}

// ########################
// ######  Voting  ########
// ########################
// GetVoteForRule returns the client's vote in favour of or against a rule.
// func (c *client) GetVoteForRule(ruleName string) bool {
// 	for _, val := range c.favourRules {
// 		if val == ruleName {
// 			return true
// 		}
// 	}
// 	return false
// }
