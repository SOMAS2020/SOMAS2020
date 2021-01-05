package team5

import (
	"sort"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// func (c *client) GetVoteForElection(roleToElect shared.Role) []shared.ClientID {
// 	// No dynamic strategy right now. Team 5 always get the first vote.
// 	var vote []shared.ClientID
// 	var initialVote []shared.ClientID
// 	vote = append(vote, shared.Team5) // Team 5 alaways first
// 	initialVote = append(initialVote, shared.Team1)
// 	initialVote = append(initialVote, shared.Team2)
// 	initialVote = append(initialVote, shared.Team4)
// 	initialVote = append(initialVote, shared.Team6)
// 	// shuffle those teams
// 	rand.Shuffle(len(initialVote), func(i, j int) { initialVote[i], initialVote[j] = initialVote[j], initialVote[i] })
// 	initialVote = append(initialVote, shared.Team3)
// 	vote = append(vote, initialVote...)
// 	c.Logf("Team 5 vote: %v", vote)
// 	return vote
// }

// Generate borda vote based on opinion score
func (c *client) GetVoteForElection(roleToElect shared.Role) []shared.ClientID {
	c.evaluateRoles() //evaluate roles
	//transfer current opinion score into array
	opinionSortByTeam := opinionSortByTeam()
	opinionSortByScore := opnionSortByScore()
	sortedTeamByOpinion := sortedMapOfOpinion(findIndexOfScore(opinionSortByScore, opinionSortByTeam)) 
	return sortedTeamByOpinion
}

// opinion score sort by team number 1,2,3,4,5,6
func (c *client) opinionSortByTeam() (opnionSortByTeam []opinion) {
	opnionSortByTeam = append(opnionSortByTeam, *c.opinions[shared.Team1], *c.opinions[shared.Team2])
	opnionSortByTeam = append(opnionSortByTeam, *c.opinions[shared.Team3], *c.opinions[shared.Team4])
	opnionSortByTeam = append(opnionSortByTeam, *c.opinions[shared.Team5], *c.opinions[shared.Team6])
	return
}

//opinion sorts by scores from min to max
func (c *client) opinionSortByScore() (opinionSortByScore []opinion) {
	opinionSortByScore = c.opinionSortByTeam()
	sort.Slice(opinionSortByScore, func(i, j int) bool {
		return opinionSortByScore[i].score < opinionSortByScore[j].score
	})
	return
}

func findIndexOfScore(opinionSortByScore []opinion, opinionSortByTeam []opinion) (rank []int){
	for i:=len(opinionSortByScore)-1; i>=0; i-- {
		for j := len(opinionSortByScore)-1; j>=0; j-- {
			if opinionSortByTeam[j] == opinionSortByScore[i] {
				rank = append(rank, j)
				opinionSortByTeam[j].score = -100 //number never get used so don't mix 
			}
		}
	}
	return
}
 
func sortedMapOfOpinion(rank []int) (sortedMap []shared.ClientID) {
	// translate int to shared.clientID but put our ID first
	sortedMap = append(sortedMap, shared.Team5)
	for i:= 0; i<len(rank); i++ {
		if rank[i] == 0 {
			sortedMap = append(sortedMap, shared.Team1)
		} else if rank[i] == 1 {
			sortedMap = append(sortedMap, shared.Team2)
		} else if rank [i] == 2 {
			sortedMap = append(sortedMap, shared.Team3)
		} else if rank[i] == 3 {
			sortedMap = append(sortedMap, shared.Team4)
		} else if rank [i] == 5 {
			sortedMap = append(sortedMap, shared.Team6)
	}
	return sortedMap
}