package team6

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// VoteForRule returns the client's vote in favour of or against a rule.
func (c *client) VoteForRule(ruleMatrix rules.RuleMatrix) shared.RuleVoteType {
	// TODO implement decision on voting that considers the rule
	return shared.Abstain
}

// VoteForElection returns the client's Borda vote for the role to be elected.
func (c *client) VoteForElection(roleToElect shared.Role, candidateList []shared.ClientID) []shared.ClientID {
	//Sort candidates according to friendship level as a preference list
	idToSlice := []shared.ClientID{id}
	
	for i, candidateID := range candidateList {
		if candidateID == id {
			candidateList = append(candidateList[:i], candidateList[i+1:]...)
			break
		}
	} 

	//Rank candidates except yourself according to friendship level
	for i := 0; i < len(candidateList); i++ {
		for j := i; j < len(candidateList); j++ {
			if c.friendship[candidateList[j]] > c.friendship[candidateList[i]] {
				candidateList[i], candidateList[j] = candidateList[j], candidateList[i]
			}
		}
	}

	//Votes for itself according to the preference for different roles when more than 3 islands alive,
	//put yourself in the first, second or third place in the preference list according to rolePreference
	//Otherwise put yourself at the last place of the list(3 or less than 3 islands alive)
	rolePreference := []shared.Role{shared.President, shared.Speaker, shared.Judge}
	if len(candidateList) > 3 {
		for i, role := range rolePreference {
			if role == roleToElect {
				insertID := append(candidateList[:i], idToSlice...)
				candidateList = append(insertID, candidateList[i:]...)
			}
		}
	} else {
		candidateList = append(candidateList, idToSlice...)
	}

	preferenceList := candidateList

	return preferenceList
}
