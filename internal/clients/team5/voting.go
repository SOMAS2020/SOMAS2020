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
	sortedTeamByOpinion := c.sortedMapOfOpinion(findIndexOfScore(opinionSortByScore, opinionSortByTeam))
	c.Logf("opinionSortByTeam(): %v", opinionSortByTeam)
	c.Logf("opinionSortByScore(): %v", opinionSortByScore)
	c.Logf("sortedTeamByOpinion(): %v", sortedTeamByOpinion)
	return sortedTeamByOpinion
}

// opinion score in order of team number 1,2,3,4,5,6
func (c *client) opinionSortByTeam() (opnionSortByTeam []opinionScore) {
	opnionSortByTeam = append(opnionSortByTeam, c.opinions[shared.Team1].getScore(), c.opinions[shared.Team2].getScore())
	opnionSortByTeam = append(opnionSortByTeam, c.opinions[shared.Team3].getScore(), c.opinions[shared.Team4].getScore())
	opnionSortByTeam = append(opnionSortByTeam, c.opinions[shared.Team5].getScore(), c.opinions[shared.Team6].getScore())
	return
}

//opinion sorted by scores from min to max
func (c *client) opinionSortByScore() (opinionSortByScore []opinionScore) {
	opinionSortByScore = c.opinionSortByTeam()
	sort.Slice(opinionSortByScore, func(i, j int) bool {
		return opinionSortByScore[i] < opinionSortByScore[j]
	})
	return
}

// arrange teams corresponding to their opinion score from max -> min
func findIndexOfScore(opinionSortByScore []opinionScore, opinionSortByTeam []opinionScore) (rank []int) {
	for i := len(opinionSortByScore) - 1; i >= 0; i-- {
		for j := len(opinionSortByScore) - 1; j >= 0; j-- {
			if opinionSortByTeam[j] == opinionSortByScore[i] {
				rank = append(rank, j)
				opinionSortByTeam[j] = 100 // assign random number we know we won't have so it won't mixed up the index
			}
		}
	}
	return
}

// translate int to shared.clientID but put our ID first and someone last
func (c *client) sortedMapOfOpinion(rank []int) (sortedMap []shared.ClientID) {
	c.Logf("rank: %v", rank)
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

//Evaluate if the roles are corrupted or not based on their budget spending versus total tax paid to common pool
//Either everyone is corrupted or not
func (c *client) evaluateRoles() {
	speakerID := c.ServerReadHandle.GetGameState().SpeakerID
	judgeID := c.ServerReadHandle.GetGameState().JudgeID
	presidentID := c.ServerReadHandle.GetGameState().PresidentID
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
