package team5

import (
	"sort"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// Generate borda vote based on opinion score for roles
// Team 5 will always be top
func (c *client) VoteForElection(roleToElect shared.Role, candidateList []shared.ClientID) []shared.ClientID {
	var refinedCandidateList []shared.ClientID
	var intRefinedCandidateList []int
	// Take out our ID
	for _, islandID := range candidateList {
		if islandID != shared.Team5 && islandID != shared.Team3 {
			refinedCandidateList = append(refinedCandidateList, islandID)
		}
	}
	// translate to int
	for _, islandID := range refinedCandidateList {
		intRefinedCandidateList = append(intRefinedCandidateList, int(islandID))
	}
	opinionSortByTeam := c.opinionSortByTeam(intRefinedCandidateList)
	opinionSortByScore := c.opinionSortByTeam(intRefinedCandidateList)
	opinionSortByScore = c.opinionSortByScore(opinionSortByScore)
	ballot := c.sortedMapOfOpinion(c.findIndexOfScore(opinionSortByScore, opinionSortByTeam, intRefinedCandidateList), intRefinedCandidateList)
	//last one
	for _, islandID := range candidateList {
		if islandID == shared.Team3 {
			ballot = append(ballot, shared.Team3)
		}
	}
	c.Logf("[DEBUG] - Ballot: %v", ballot)

	return ballot
}

// opinion score in order of team number 1,2,3,4,5,6
func (c *client) opinionSortByTeam(candidateList []int) (opinionSortByTeam []float64) {
	sort.Ints(candidateList)
	for _, island := range candidateList {
		opinionSortByTeam = append(opinionSortByTeam, float64(c.opinions[shared.ClientID(island)].getScore()))
	}
	return opinionSortByTeam
}

//opinion sorted by scores from min to max
func (c *client) opinionSortByScore(opinionSortByTeam []float64) (opinionSortByScore []float64) {
	opinionSortByScore = opinionSortByTeam
	sort.Float64s(opinionSortByScore)
	return
}

// arrange teams corresponding to their opinion score from max -> min
func (c *client) findIndexOfScore(opinionSortByScore []float64, opinionSortByTeam []float64, candidateList []int) (rank []int) {
	for i := len(opinionSortByScore) - 1; i >= 0; i-- {
		//for i := 0; i < len(opinionSortByScore); i++ {
		for j := len(opinionSortByScore) - 1; j >= 0; j-- {
			//for j := 0; j < len(opinionSortByScore); j++ {
			if opinionSortByScore[i] == opinionSortByTeam[j] {
				rank = append(rank, candidateList[j])
				opinionSortByTeam[j] = 100 // assign random number we know we won't have so it won't mixed up the index
			}
		}
	}
	return
}

// translate int to shared.clientID but put our ID first and someone last
// assume that we are always alive when this function is called
func (c *client) sortedMapOfOpinion(rank []int, candidateList []int) (sortedTeamByOpinion []shared.ClientID) {
	var sortedCandidateList []int
	allIsland := []int{0, 1, 2, 3, 4, 5}
	for _, island := range allIsland {
		for _, clientID := range candidateList {
			if clientID == island && clientID != 4 && clientID != 2 {
				sortedCandidateList = append(sortedCandidateList, island)
			}
		}
	}
	sortedTeamByOpinion = append(sortedTeamByOpinion, shared.Team5)
	for _, clientID := range rank {
		sortedTeamByOpinion = append(sortedTeamByOpinion, shared.ClientID(clientID))
	}
	return sortedTeamByOpinion
}

//Evaluate if the roles are corrupted or not based on their budget spending versus total tax paid to common pool
//Either everyone is corrupted or not
func (c *client) evaluateRoles() {
	speakerID := c.gameState().SpeakerID
	judgeID := c.gameState().JudgeID
	presidentID := c.gameState().PresidentID
	c.Logf("[DEBUG] - speakerID: %v, judgeID: %v, presidentID: %v", speakerID, judgeID, presidentID)
	//compute total budget
	budget := c.ServerReadHandle.GetGameState().IIGORolesBudget
	var totalBudget shared.Resources = 0
	for role := range budget {
		totalBudget += budget[role]
	}
	// compute total maximum tax to cp
	var totalTax shared.Resources
	numberAliveTeams := len(c.getAliveTeams(true)) //include us
	for i := 0; i < numberAliveTeams; i++ {
		totalTax += c.taxAmount
	}
	// Not corrupt
	if totalBudget <= totalTax {
		c.opinions[speakerID].updateOpinion(generalBasis, 0.1) //arbitrary number
		c.opinions[judgeID].updateOpinion(generalBasis, 0.1)
		c.opinions[presidentID].updateOpinion(generalBasis, 0.1)

	} else {
		c.opinions[speakerID].updateOpinion(generalBasis, -0.1) //arbitrary number
		c.opinions[judgeID].updateOpinion(generalBasis, -0.1)
		c.opinions[presidentID].updateOpinion(generalBasis, -0.1)
	}
}
