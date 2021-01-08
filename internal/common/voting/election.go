package voting

import (
	"fmt"
	"math"
	"sort"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type Election struct {
	roleToElect   shared.Role
	votingMethod  shared.ElectionVotingMethod
	candidateList []shared.ClientID
	voterList     []shared.ClientID
	votes         [][]shared.ClientID
	Logger        shared.Logger
}

// Logf is the Election logger
func (e *Election) Logf(format string, a ...interface{}) {
	e.Logger("[ELECTION]: %v", fmt.Sprintf(format, a...))
}

// ProposeMotion sets the role to be voted on
func (e *Election) ProposeElection(role shared.Role, method shared.ElectionVotingMethod) {
	e.roleToElect = role
	e.votingMethod = method
}

// OpenBallot sets the islands eligible to vote.
func (e *Election) OpenBallot(clientIDs []shared.ClientID, allIslands []shared.ClientID) {
	e.voterList = clientIDs
	//Get candidate list in sorted order.
	sort.SliceStable(allIslands, func(i, j int) bool {
		return int(allIslands[i]) < int(allIslands[j])
	})
	e.candidateList = allIslands
}

// Vote gets votes from eligible islands.
func (e *Election) Vote(clientMap map[shared.ClientID]baseclient.Client) {
	for i := 0; i < len(e.voterList); i++ {
		e.votes = append(e.votes, clientMap[e.voterList[i]].VoteForElection(e.roleToElect, e.candidateList))
	}
	e.Logf("Votes: %v", e.votes)
}

// CloseBallot counts the votes received and returns the result.
func (e *Election) CloseBallot(clientMap map[shared.ClientID]baseclient.Client) shared.ClientID {

	var result shared.ClientID

	switch e.votingMethod {
	case shared.BordaCount:
		result = e.bordaCountResult()
	case shared.Runoff:
		result = e.runOffResult(clientMap)
	case shared.InstantRunoff:
		result = e.instantRunoffResult(clientMap)
	case shared.Approval:
		result = e.approvalResult()
	}
	return result
}

//func (e *Election) completePreferenceMap()

func scoreCalculator(totalVotes [][]shared.ClientID, candidateList []shared.ClientID) ([]float64, []float64, float64) {
	votesLayoutElect := make(map[int][]int)
	votesSliceSquare := totalVotes
	candidatesNumber := len(candidateList)
	islandsNumber := len(votesSliceSquare)

	//Initialize votesLayoutMap
	for i := 1; i < islandsNumber+1; i++ {
		for j := 0; j < candidatesNumber; j++ {
			votesLayoutElect[i] = append(votesLayoutElect[i], 0)
		}
	}
	//Transfer e.votes to votesLayoutMap with type "int"
	for i := 0; i < islandsNumber; i++ {
		scoreInit := candidatesNumber + 1
		for j := 0; j < len(votesSliceSquare[i]); j++ {
			for k := 0; k < candidatesNumber; k++ {
				if votesSliceSquare[i][j] == candidateList[k] {
					votesLayoutElect[i+1][k] = scoreInit
					scoreInit--
				}
			}
		}
	}

	//Sort the preference map in order.
	order := make([]int, candidatesNumber)
	index := make([]int, candidatesNumber)
	score := make([]float64, candidatesNumber)
	scoreMap := make(map[int][]float64)
	for k, v := range votesLayoutElect {
		j := 0
		var absentNum float64 = 0

		for i := 0; i < candidatesNumber; i++ {
			order[i] = v[i]
		}

		for i := 0; i < candidatesNumber; i++ {

			sum := 0
			for t := 0; t < candidatesNumber; t++ {
				sum = sum + v[t]
			}

			searcher := sum

			for p := 0; p < candidatesNumber; p++ {
				if searcher > order[p] {
					searcher = order[p]
					index[i] = p
				}
			}

			j = index[i]
			order[j] = sum
		}

		itrans := 0
		for i := 0; i < candidatesNumber; i++ {
			itrans = index[i]
			score[itrans] = float64(i + 1)
		}

		for i := 0; i < candidatesNumber; i++ {
			scoreMap[k] = append(scoreMap[k], score[i])
		}

		for i := 0; i < len(v); i++ {
			if v[i] == 0 {
				absentNum++
			}
		}

		for i := 0; i < len(v); i++ {
			if v[i] == 0 {
				scoreMap[k][i] = (1 + absentNum) / 2
			}
		}

	}

	//Calculate the final score for all candidates.
	finalScore := make([]float64, candidatesNumber)
	for _, v := range scoreMap {
		for i := 0; i < candidatesNumber; i++ {
			finalScore[i] += v[i]
		}
	}
	//variance is needed when two or more candidates have equal votes.
	variance := make([]float64, candidatesNumber)
	for _, v := range scoreMap {
		for i := 0; i < candidatesNumber; i++ {
			cN := float64(candidatesNumber)
			variance[i] += math.Pow((v[i] - finalScore[i]/cN), 2)
		}
	}

	var totalScore float64 = 0
	for _, v := range finalScore {
		totalScore += v
	}

	return finalScore, variance, totalScore
}

func findMaxScore(scoreList []float64, variance []float64) (float64, int) {
	var maxScore float64 = 0
	maxScoreIndex := 0
	for i := 0; i < len(scoreList); i++ {
		if scoreList[i] > maxScore {
			maxScore = scoreList[i]
			maxScoreIndex = i
		} else if scoreList[i] == maxScore {
			if variance[i] > variance[maxScoreIndex] {
				maxScoreIndex = i
			}
		}
	}
	return maxScore, maxScoreIndex
}

func findMinScore(scoreList []float64, variance []float64) (float64, int) {
	var minScore float64 = 0
	for i := 0; i < len(scoreList); i++ {
		minScore += scoreList[i]
	}
	minScoreIndex := 0
	for i := 0; i < len(scoreList); i++ {
		if scoreList[i] < minScore {
			minScore = scoreList[i]
			minScoreIndex = i
		} else if scoreList[i] == minScore {
			if variance[i] < variance[minScoreIndex] {
				minScoreIndex = i
			}
		}
	}
	return minScore, minScoreIndex
}

func (e *Election) bordaCountResult() shared.ClientID {
	// Implement Borda count winner selection method
	candidatesNumber := len(e.candidateList)
	finalScore, variance, _ := scoreCalculator(e.votes, e.candidateList)

	var maxScore float64 = 0
	var winnerIndex int
	winnerIndex = 0
	for i := 0; i < candidatesNumber; i++ {
		if maxScore < finalScore[i] {
			maxScore = finalScore[i]
			winnerIndex = i
		} else if maxScore == finalScore[i] {
			if variance[winnerIndex] < variance[i] {
				winnerIndex = i
			}
		}
	}
	winner := e.candidateList[winnerIndex]

	return winner
}

func (e *Election) runOffResult(clientMap map[shared.ClientID]baseclient.Client) shared.ClientID {
	var winner shared.ClientID
	//Round one
	scoreList, variance, totalScore := scoreCalculator(e.votes, e.candidateList)
	rOneCandidateList := e.candidateList
	voterNumber := len(e.voterList)

	halfTotalScore := 0.5 * totalScore

	maxScore, maxScoreIndex := findMaxScore(scoreList, variance)

	if maxScore > halfTotalScore {
		winner = rOneCandidateList[maxScoreIndex]
	} else {
		//Round two
		remainNumber := 0
		changeNumber := 0
		scoreList[maxScoreIndex] = 0

		_, competitorIndex := findMaxScore(scoreList, variance)

		rTwoCandidateList := []shared.ClientID{rOneCandidateList[maxScoreIndex], rOneCandidateList[competitorIndex]}

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
	return winner
}

func (e *Election) instantRunoffResult(clientMap map[shared.ClientID]baseclient.Client) shared.ClientID {
	var winner shared.ClientID
	candidateNumber := len(e.candidateList)
	candidateList := e.candidateList
	totalVotes := e.votes
	var halfTotalScore float64 = 0

	for {
		scoreList, variance, totalScore := scoreCalculator(totalVotes, candidateList)

		halfTotalScore = 0.5 * totalScore

		maxScore, maxScoreIndex := findMaxScore(scoreList, variance)

		//Eliminate the least popular one until the most popular one has more than half of the total score.
		if maxScore > halfTotalScore {
			winner = candidateList[maxScoreIndex]
			break
		}

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
	return winner
}

//Election method only considering the number of times the candidate appears on the preference list.
func (e *Election) approvalResult() shared.ClientID {
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
	return winner
}

func (e *Election) pluralityResult() shared.ClientID {

	// How many first place votes did each island get
	votesPerIsland := map[shared.ClientID]int{}
	for _, ranking := range e.votes {
		votesPerIsland[ranking[0]] += 1
	}

	// Who got the most first place votes
	winVote := 0
	winner := shared.ClientID(1)
	for island, votes := range votesPerIsland {
		if votes >= winVote {
			winVote = votes
			winner = island
		}
	}
	return winner
}

func (e *Election) majorityResult() shared.ClientID {
	// TODO implement majority winner selection method.
	return e.pluralityResult()
}
