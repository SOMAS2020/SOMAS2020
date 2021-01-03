package voting

import (
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func TestElection(t *testing.T) {
	var ele Election
	ele.roleToElect = 0
	ele.votingMethod = 0
	ele.candidateList = []shared.ClientID{shared.Team1, shared.Team2, shared.Team3, shared.Team4, shared.Team5, shared.Team6}
	ele.voterList = ele.candidateList
	ele.votes = [][]shared.ClientID{{shared.Team3, shared.Team2, shared.Team1, shared.Team6, shared.Team5, shared.Team4},
		{shared.Team4, shared.Team6, shared.Team5, shared.Team2, shared.Team3, shared.Team1},
		{shared.Team3, shared.Team6, shared.Team1, shared.Team4, shared.Team5, shared.Team2},
		{shared.Team2, shared.Team5, shared.Team6, shared.Team4, shared.Team3, shared.Team1},
		{shared.Team6, shared.Team4, shared.Team1, shared.Team5, shared.Team2, shared.Team3},
		{shared.Team5, shared.Team2, shared.Team3, shared.Team6, shared.Team1, shared.Team4}}
	clientMap := make(map[shared.ClientID]baseclient.Client)

	var c baseclient.BaseClient

	for i := 0; i < len(ele.voterList); i++ {
		clientMap[ele.voterList[i]] = &c
	}

	ele.bordaCountTest()
	ele.runOffTest(clientMap)
	ele.instantRunoffTest(clientMap)
	ele.approvalTest()

}

func (e *Election) bordaCountTest() ([]int, shared.ClientID) {
	// Implement Borda count winner selection method
	candidatesNumber := len(e.candidateList)
	scoreList, variance, _ := e.scoreCalculator(e.votes, e.candidateList)

	maxScore := 0
	var winnerIndex int
	winnerIndex = 0
	for i := 0; i < candidatesNumber; i++ {
		if maxScore < scoreList[i] {
			maxScore = scoreList[i]
			winnerIndex = i
		}
		if maxScore == scoreList[i] {
			if variance[winnerIndex] < variance[i] {
				winnerIndex = i
			}
		}
	}
	winner := e.candidateList[winnerIndex]

	return scoreList, winner
}

func (e *Election) runOffTest(clientMap map[shared.ClientID]baseclient.Client) (map[int][]shared.ClientID, shared.ClientID) {
	var winner shared.ClientID
	//Round one
	scoreList, variance, totalScore := e.scoreCalculator(e.votes, e.candidateList)
	rOneCandidateList := e.candidateList
	voterNumber := len(e.voterList)
	roundCandidateMap := make(map[int][]shared.ClientID)
	roundCandidateMap[1] = rOneCandidateList

	halfTotalScore := 0.5 * totalScore

	maxScore, maxScoreIndex := findMaxScore(scoreList, variance)
	fmt.Println(maxScore, maxScoreIndex, totalScore)

	if float64(maxScore) > halfTotalScore {
		winner = rOneCandidateList[maxScoreIndex]
	} else {
		//Round two
		remainNumber := 0
		changeNumber := 0
		scoreList[maxScoreIndex] = 0

		_, competitorIndex := findMaxScore(scoreList, variance)

		rTwoCandidateList := []shared.ClientID{rOneCandidateList[maxScoreIndex], rOneCandidateList[competitorIndex]}
		roundCandidateMap[2] = rTwoCandidateList

		fmt.Println(rTwoCandidateList)
		var rTwoVotes [][]shared.ClientID
		for i := 0; i < voterNumber; i++ {
			rTwoVotes = append(rTwoVotes, clientMap[e.voterList[i]].VoteForElection(e.roleToElect, rTwoCandidateList))
		}
		fmt.Println(rTwoVotes)
		for i := 0; i < voterNumber; i++ {
			if rTwoVotes[i][0] == rOneCandidateList[maxScoreIndex] {
				remainNumber++
			}
			if rTwoVotes[i][0] == rOneCandidateList[competitorIndex] {
				changeNumber++
			}
		}
		if changeNumber > remainNumber {
			winner = rOneCandidateList[competitorIndex]
		} else {
			winner = rOneCandidateList[maxScoreIndex]
		}
	}
	return roundCandidateMap, winner
}

func (e *Election) instantRunoffTest(clientMap map[shared.ClientID]baseclient.Client) (map[int][]shared.ClientID, map[int][]int, shared.ClientID) {
	var winner shared.ClientID
	candidateNumber := len(e.candidateList)
	candidateList := e.candidateList
	totalVotes := e.votes
	var halfTotalScore float64 = 0
	roundCandidateMap := make(map[int][]shared.ClientID)
	roundScoreList := make(map[int][]int)
	roundCount := 1
	roundCandidateMap[roundCount] = candidateList

	for {
		scoreList, variance, totalScore := e.scoreCalculator(totalVotes, candidateList)

		roundScoreList[roundCount] = scoreList

		halfTotalScore = 0.5 * totalScore

		maxScore, maxScoreIndex := findMaxScore(scoreList, variance)

		//keep eliminating the least popular one untill the most popular one has more than half of the total score.
		if float64(maxScore) > halfTotalScore {
			winner = candidateList[maxScoreIndex]
			break
		}

		roundCount++

		_, minScoreIndex := findMinScore(scoreList, variance)

		//Eliminate the least popular candidate
		if minScoreIndex == 0 {
			candidateList = candidateList[minScoreIndex+1:]
		} else if minScoreIndex == candidateNumber-1 {
			candidateList = candidateList[:minScoreIndex]
		} else {
			candidateList = append(candidateList[:minScoreIndex], candidateList[minScoreIndex+1:]...)
		}
		candidateNumber--

		//New round voting status update
		for i := 0; i < len(e.voterList); i++ {
			totalVotes = append(totalVotes, clientMap[e.voterList[i]].VoteForElection(e.roleToElect, candidateList))
		}
		totalVotes = totalVotes[(len(totalVotes) - len(e.voterList)):]

	}
	return roundCandidateMap, roundScoreList, winner
}

func (e *Election) approvalTest() ([][]shared.ClientID, shared.ClientID) {
	candidateList := e.candidateList
	scoreList := make([]int, len(candidateList))
	//If there are more than two candidates has the highest score, then the winner will be randomly chosen.
	for i := 0; i < len(e.votes); i++ {
		for j := 0; j < len(e.votes[i]); j++ {
			for p := 0; p < len(candidateList); p++ {
				if candidateList[p] == e.votes[i][j] {
					scoreList[p] += 1
				}
			}
		}
	}
	maxScore := 0
	maxScoreIndex := 0
	for i := 0; i < len(candidateList); i++ {
		if scoreList[i] > maxScore {
			maxScore = scoreList[i]
			maxScoreIndex = i
		}
	}
	winner := candidateList[maxScoreIndex]
	return e.votes, winner
}
