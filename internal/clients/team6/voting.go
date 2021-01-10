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

func (c *client) rolesInfro() map[shared.Role]shared.ClientID {
	rolesInfro := make(map[shared.Role]shared.ClientID)
	rolesInfro[shared.President] = c.ServerReadHandle.GetGameState().PresidentID
	rolesInfro[shared.Speaker] = c.ServerReadHandle.GetGameState().SpeakerID
	rolesInfro[shared.Judge] = c.ServerReadHandle.GetGameState().JudgeID
	return rolesInfro
}

func (c *client) doWeHaveRoles(roleToElect shared.Role) bool {
	doWeHaveRoles := false
	numOfRoles := 0
	for role, roleID := range c.rolesInfro() {
		if roleID == id {
			numOfRoles++
			if role == roleToElect {
				numOfRoles--
			}
		}
	}
	if numOfRoles > 0 {
		doWeHaveRoles = true
	}
	return doWeHaveRoles
}

// VoteForElection returns the client's Borda vote for the role to be elected.
func (c *client) VoteForElection(roleToElect shared.Role, candidateList []shared.ClientID) []shared.ClientID {
	//Create a slice for id
	idToSlice := []shared.ClientID{id}

	doWeHaveRoles := c.doWeHaveRoles(roleToElect)
	
	//Rank candidates except yourself according to friendship level
	for i, candidateID := range candidateList {
		if candidateID == id {
			candidateList = append(candidateList[:i], candidateList[i+1:]...)
			break
		}
	} 

	for i := 0; i < len(candidateList); i++ {
		for j := i; j < len(candidateList); j++ {
			if c.friendship[candidateList[j]] > c.friendship[candidateList[i]] {
				candidateList[i], candidateList[j] = candidateList[j], candidateList[i]
			}
		}
	}

	//Votes for ourselves according to the preference for different roles when more than 3 islands alive and we don't have roles
	//put ourselves in the first, second or third place in the preference list according to rolePreference
	//Otherwise put ourselves at the last place of the list(3 or less than 3 islands alive)
	rolePreference := []shared.Role{shared.President, shared.Speaker, shared.Judge}
	if len(candidateList) > 3 && doWeHaveRoles == false {
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
