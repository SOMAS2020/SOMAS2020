package voting

import (
	"reflect"
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func TestElection(t *testing.T) {
	var ele Election
	ele.roleToElect = 0
	ele.votingMethod = 0
	ele.candidateList = []shared.ClientID{shared.Teams["Team1"], shared.Teams["Team2"], shared.Teams["Team3"], shared.Teams["Team4"], shared.Teams["Team5"], shared.Teams["Team6"]}
	ele.voterList = ele.candidateList
	ele.votes = [][]shared.ClientID{{shared.Teams["Team3"], shared.Teams["Team2"], shared.Teams["Team1"], shared.Teams["Team6"], shared.Teams["Team5"], shared.Teams["Team4"]},
		{shared.Teams["Team4"], shared.Teams["Team6"], shared.Teams["Team5"], shared.Teams["Team2"], shared.Teams["Team3"], shared.Teams["Team1"]},
		{shared.Teams["Team3"], shared.Teams["Team6"], shared.Teams["Team1"], shared.Teams["Team4"], shared.Teams["Team5"], shared.Teams["Team2"]},
		{shared.Teams["Team2"], shared.Teams["Team5"], shared.Teams["Team6"], shared.Teams["Team4"], shared.Teams["Team3"], shared.Teams["Team1"]},
		{shared.Teams["Team6"], shared.Teams["Team4"], shared.Teams["Team1"], shared.Teams["Team5"], shared.Teams["Team2"], shared.Teams["Team3"]},
		{shared.Teams["Team5"], shared.Teams["Team2"], shared.Teams["Team3"], shared.Teams["Team6"], shared.Teams["Team1"], shared.Teams["Team4"]}}
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

func (e *Election) bordaCountTest() ([]float64, shared.ClientID) {
	// Implement Borda count winner selection method
	candidatesNumber := len(e.candidateList)
	scoreList, variance, _ := scoreCalculator(e.votes, e.candidateList)

	var maxScore float64 = 0
	var winnerIndex int
	winnerIndex = 0
	for i := 0; i < candidatesNumber; i++ {
		if maxScore < scoreList[i] {
			maxScore = scoreList[i]
			winnerIndex = i
		} else if maxScore == scoreList[i] {
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
	scoreList, variance, totalScore := scoreCalculator(e.votes, e.candidateList)
	rOneCandidateList := e.candidateList
	voterNumber := len(e.voterList)
	roundCandidateMap := make(map[int][]shared.ClientID)
	roundCandidateMap[1] = rOneCandidateList

	halfTotalScore := 0.5 * totalScore

	maxScore, maxScoreIndex := findMaxScore(scoreList, variance)

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

		var rTwoVotes [][]shared.ClientID
		for i := 0; i < voterNumber; i++ {
			rTwoVotes = append(rTwoVotes, clientMap[e.voterList[i]].VoteForElection(e.roleToElect, rTwoCandidateList))
		}

		for i := 0; i < voterNumber; i++ {
			if rTwoVotes[i][0] == rOneCandidateList[maxScoreIndex] {
				remainNumber++
			} else if rTwoVotes[i][0] == rOneCandidateList[competitorIndex] {
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

func (e *Election) instantRunoffTest(clientMap map[shared.ClientID]baseclient.Client) (map[int][]shared.ClientID, map[int][]float64, shared.ClientID) {
	var winner shared.ClientID
	candidateNumber := len(e.candidateList)
	candidateList := e.candidateList
	totalVotes := e.votes
	var halfTotalScore float64 = 0
	roundCandidateMap := make(map[int][]shared.ClientID)
	roundScoreList := make(map[int][]float64)
	roundCount := 1
	roundCandidateMap[roundCount] = candidateList

	for {
		scoreList, variance, totalScore := scoreCalculator(totalVotes, candidateList)

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

func TestOpenBallot(t *testing.T) {
	cases := []struct {
		name               string
		voters             []shared.ClientID
		allIslands         []shared.ClientID
		expectedVoters     []shared.ClientID
		expectedCandidates []shared.ClientID
	}{
		{
			name:               "basic_test",
			voters:             []shared.ClientID{shared.ClientID(1), shared.ClientID(3)},
			allIslands:         []shared.ClientID{shared.ClientID(1), shared.ClientID(3)},
			expectedVoters:     []shared.ClientID{shared.ClientID(1), shared.ClientID(3)},
			expectedCandidates: []shared.ClientID{shared.ClientID(1), shared.ClientID(3)},
		},
		{
			name:               "mismatched_voters_and_candidates_test",
			voters:             []shared.ClientID{shared.ClientID(1), shared.ClientID(3)},
			allIslands:         []shared.ClientID{shared.ClientID(1), shared.ClientID(2), shared.ClientID(3)},
			expectedVoters:     []shared.ClientID{shared.ClientID(1), shared.ClientID(3)},
			expectedCandidates: []shared.ClientID{shared.ClientID(1), shared.ClientID(2), shared.ClientID(3)},
		},
		{
			name:               "mismatched_and_disordered_voters_and_candidates_test",
			voters:             []shared.ClientID{shared.ClientID(1), shared.ClientID(3)},
			allIslands:         []shared.ClientID{shared.ClientID(2), shared.ClientID(3), shared.ClientID(1)},
			expectedVoters:     []shared.ClientID{shared.ClientID(1), shared.ClientID(3)},
			expectedCandidates: []shared.ClientID{shared.ClientID(1), shared.ClientID(2), shared.ClientID(3)},
		},
		{
			name:               "empty_lists_test",
			voters:             []shared.ClientID{},
			allIslands:         []shared.ClientID{},
			expectedVoters:     []shared.ClientID{},
			expectedCandidates: []shared.ClientID{},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			testelection := &Election{
				candidateList: []shared.ClientID{},
				voterList:     []shared.ClientID{},
			}
			testelection.OpenBallot(tc.voters, tc.allIslands)
			resVoters := testelection.voterList
			resCandidates := testelection.candidateList
			if !reflect.DeepEqual(resVoters, tc.expectedVoters) {
				t.Errorf("Expected voters to be %v got %v", tc.expectedVoters, resVoters)
			}
			if !reflect.DeepEqual(resCandidates, tc.expectedCandidates) {
				t.Errorf("Expected candidates to be %v got %v", tc.expectedCandidates, resCandidates)
			}
		})
	}
}
