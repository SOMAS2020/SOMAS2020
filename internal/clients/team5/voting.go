package team5

import (
	"sort"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// Generate borda vote based on opinion score for roles
// Team 5 will always be top
func (c *client) GetVoteForElection(roleToElect shared.Role) []shared.ClientID {
	c.evaluateRoles() //temporarily. Not ideal because election don't happen everyday
	opinionSortByTeam := c.opinionSortByTeam()
	opinionSortByScore := c.opinionSortByScore()
	sortedTeamByOpinion := sortedMapOfOpinion(findIndexOfScore(opinionSortByScore, opinionSortByTeam))
	return sortedTeamByOpinion
}

// opinion score in order of team number 1,2,3,4,5,6
func (c *client) opinionSortByTeam() (opnionSortByTeam []opinion) {
	opnionSortByTeam = append(opnionSortByTeam, *c.opinions[shared.Team1], *c.opinions[shared.Team2])
	opnionSortByTeam = append(opnionSortByTeam, *c.opinions[shared.Team3], *c.opinions[shared.Team4])
	opnionSortByTeam = append(opnionSortByTeam, *c.opinions[shared.Team5], *c.opinions[shared.Team6])
	return
}

//opinion sorted by scores from min to max
func (c *client) opinionSortByScore() (opinionSortByScore []opinion) {
	opinionSortByScore = c.opinionSortByTeam()
	sort.Slice(opinionSortByScore, func(i, j int) bool {
		return opinionSortByScore[i].score < opinionSortByScore[j].score
	})
	return
}

// arrange teams corresponding to their opinion score from max -> min
func findIndexOfScore(opinionSortByScore []opinion, opinionSortByTeam []opinion) (rank []int) {
	for i := len(opinionSortByScore) - 1; i >= 0; i-- {
		for j := len(opinionSortByScore) - 1; j >= 0; j-- {
			if opinionSortByTeam[j] == opinionSortByScore[i] {
				rank = append(rank, j)
				opinionSortByTeam[j].score = -100 // assign random number we know we won't have so it won't mixed up the index
			}
		}
	}
	return
}

// translate int to shared.clientID but put our ID first and someone last
func sortedMapOfOpinion(rank []int) (sortedMap []shared.ClientID) {
	sortedMap = append(sortedMap, shared.Team5)
	for i := 0; i < len(rank); i++ {
		if rank[i] == 0 {
			sortedMap = append(sortedMap, shared.Team1)
		} else if rank[i] == 1 {
			sortedMap = append(sortedMap, shared.Team2)
		} else if rank[i] == 3 {
			sortedMap = append(sortedMap, shared.Team4)
		} else if rank[i] == 5 {
			sortedMap = append(sortedMap, shared.Team6)
		}
	}
	sortedMap = append(sortedMap, shared.Team3)
	return sortedMap
}
