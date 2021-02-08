package team6

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
)

func findSameVariables(originalOne []rules.VariableFieldName, newOne []rules.VariableFieldName) []rules.VariableFieldName {
	sharedVariables := []rules.VariableFieldName{}
	for _, oVar := range originalOne {
		for _, nVar := range newOne {
			if oVar == nVar {
				sharedVariables = append(sharedVariables, oVar)
			}
		}
	}
	return sharedVariables
}

//VoteForRule returns the client's vote in favour of or against a rule.
//If the corresponding varlables are what we care about, then we vote approve, otherwise abstain
func (c *client) VoteForRule(ruleMatrix rules.RuleMatrix) shared.RuleVoteType {
	ourAttitude := shared.Abstain
	isItImportant := false
	//doesItFitUs := false
	numOfIslandsAlive := 6
	variablesIncreasePrefered := []rules.VariableFieldName{}
	variablesDecreasePrefered := []rules.VariableFieldName{}

	if numOfIslandsAlive > 3 {
		variablesIncreasePrefered = append(variablesIncreasePrefered, rules.ExpectedAllocation)
		variablesDecreasePrefered = append(variablesDecreasePrefered, rules.ExpectedTaxContribution)
	} else {
		variablesIncreasePrefered = append(variablesIncreasePrefered, rules.ExpectedTaxContribution)
		variablesDecreasePrefered = append(variablesDecreasePrefered, rules.ExpectedAllocation)
	}

	variablesIncreasePrefered = append(variablesIncreasePrefered, rules.SanctionExpected)

	for role, roleID := range c.rolesInfro() {
		if role == shared.President && roleID == id {
			variablesIncreasePrefered = append(variablesIncreasePrefered, rules.PresidentSalary)
			variablesDecreasePrefered = append(variablesDecreasePrefered, rules.PresidentPayment)
		} else if role == shared.Judge && roleID == id {
			variablesIncreasePrefered = append(variablesIncreasePrefered, rules.JudgeSalary)
			variablesDecreasePrefered = append(variablesDecreasePrefered, rules.JudgePayment)
		} else if role == shared.Speaker && roleID == id {
			variablesIncreasePrefered = append(variablesIncreasePrefered, rules.SpeakerSalary)
			variablesDecreasePrefered = append(variablesDecreasePrefered, rules.SpeakerPayment)
		}
	}

	variablesWeCareAbout := append(variablesIncreasePrefered, variablesDecreasePrefered...)
	if len(findSameVariables(ruleMatrix.RequiredVariables, variablesWeCareAbout)) > 0 {
		isItImportant = true
	}


	if isItImportant {
		ourAttitude = shared.Approve
	} 

	return ourAttitude
}

func (c *client) rolesInfro() map[shared.Role]shared.ClientID {
	rolesInfro := make(map[shared.Role]shared.ClientID)
	rolesInfro[shared.President] = c.ServerReadHandle.GetGameState().PresidentID
	rolesInfro[shared.Speaker] = c.ServerReadHandle.GetGameState().SpeakerID
	rolesInfro[shared.Judge] = c.ServerReadHandle.GetGameState().JudgeID
	return rolesInfro
}

//Figure out whether we already have role(s) except for the one to be changed. 
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
	roleIndex := 0
	if len(candidateList) > 3 && doWeHaveRoles == false {
		for i, role := range rolePreference {
			if role == roleToElect {
				roleIndex = i
			}
		}
		insertID := append(idToSlice, candidateList[roleIndex:]...)
		candidateList = append(candidateList[:roleIndex], insertID...)
	} else {
		candidateList = append(candidateList, idToSlice...)
	}

	preferenceList := candidateList

	return preferenceList
}
