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
		if islandID != shared.Team5 {
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
	// c.Logf("[DEBUG] - Ballot: %v", ballot)

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
	sortedTeamByOpinion = append(sortedTeamByOpinion, shared.Team5)
	for _, clientID := range rank {
		sortedTeamByOpinion = append(sortedTeamByOpinion, shared.ClientID(clientID))
	}
	return sortedTeamByOpinion
}
