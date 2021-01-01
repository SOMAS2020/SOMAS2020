package voting

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type Election struct {
	roleToElect   shared.Role
	votingMethod  shared.ElectionVotingMethod
	candidateList []shared.ClientID
	voterList     []shared.ClientID
	votes         [][]shared.ClientID
}

// ProposeMotion sets the role to be voted on
func (e *Election) ProposeElection(role shared.Role, method shared.ElectionVotingMethod) {
	e.roleToElect = role
	e.votingMethod = method
}

// OpenBallot sets the islands eligible to vote.
func (e *Election) OpenBallot(clientIDs []shared.ClientID) {
	e.voterList = clientIDs
	//TODO candidate list in each election need to be determined, as it won't always equals to voter list.
	e.candidateList = e.voterList
}

// Vote gets votes from eligible islands.
func (e *Election) Vote(clientMap map[shared.ClientID]baseclient.Client) {
	for i := 0; i < len(e.voterList); i++ {
		e.votes = append(e.votes, clientMap[e.voterList[i]].VoteForElection(e.roleToElect, e.candidateList))
	}
}

// CloseBallot counts the votes received and returns the result.
func (e *Election) CloseBallot() shared.ClientID {

	var result shared.ClientID

	switch e.votingMethod {
	case shared.BordaCount:
		result = e.bordaCountResult()
	case shared.Plurality:
		result = e.pluralityResult()
	case shared.Majority:
		result = e.majorityResult()
	}
	return result
}

func (e *Election) scoreCalculator(totalVotes [][]shared.ClientID, candidateList []shared.ClientID) ([]int, []float32) {
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
		for j := 0; j < candidatesNumber; j++ {
			for k := 0; k < candidatesNumber; k++ {
				if votesSliceSquare[i][j] == e.candidateList[k] {
					votesLayoutElect[i+1][k] = scoreInit
					scoreInit--
				}
			}
		}
	}

	//Sort the preference map in order.
	order := make([]int, candidatesNumber)
	index := make([]int, candidatesNumber)
	score := make([]int, candidatesNumber)
	preferenceMap := make(map[int][]int)
	for k, v := range votesLayoutElect {
		j := 0

		for t := 0; t < candidatesNumber; t++ {
			order[t] = v[t]
		}

		for i := 0; i < candidatesNumber; i++ {

			sum := 0
			for t := 0; t < candidatesNumber; t++ {
				sum = sum + v[t]
			}

			searcher := sum

			for k := 0; k < candidatesNumber; k++ {
				if searcher > order[k] {
					searcher = order[k]
					index[i] = k
				}
			}

			j = index[i]
			order[j] = sum
		}

		itrans := 0
		for i := 0; i < candidatesNumber; i++ {
			itrans = index[i]
			score[itrans] = i + 1
		}

		for i := 0; i < candidatesNumber; i++ {
			preferenceMap[k] = append(preferenceMap[k], score[i])
		}

	}

	//Calculate the final score for all candidates.
	finalScore := make([]int, candidatesNumber)
	for _, v := range preferenceMap {
		for i := 0; i < candidatesNumber; i++ {
			finalScore[i] += v[i]
		}
	}
	//variance is needed when two or more candidates have equal votes. 
	variance := make([]float32, candidatesNumber)
	for _, v := range preferenceMap {
		for i := 0; i < candidatesNumber; i++ {
			variance[i] += (float32(v[i]) - float32(finalScore[i])/float32(candidatesNumber)) *
				(float32(v[i]) - float32(finalScore[i])/float32(candidatesNumber))
		}
	}

	return finalScore, variance
}

func (e *Election) bordaCountResult() shared.ClientID {
	// Implement Borda count winner selection method
	candidatesNumber := len(e.candidateList)
	finalScore, variance := e.scoreCalculator(e.votes, e.candidateList)

	maxScore := 0
	var winnerIndex int
	winnerIndex = 0
	for i := 0; i < candidatesNumber; i++ {
		if maxScore < finalScore[i] {
			maxScore = finalScore[i]
			winnerIndex = i
		}
		if maxScore == finalScore[i] {
			if variance[winnerIndex] < variance[i] {
				winnerIndex = i
			}
		}
	}
	var winner shared.ClientID
	winner = e.candidateList[winnerIndex]

	return winner
}

func (e *Election) runOffResult(clientMap map[shared.ClientID]baseclient.Client) shared.ClientID {
	var winner shared.ClientID
	//Round one
	finalScore, variance := e.scoreCalculator(e.votes, e.candidateList)
	candidatesNumber := len(e.candidateList)
	rOneCandidateList := e.candidateList
	voterNumber := len(e.voterList)

	var totalScore float32 = 0
	for i := 0; i < candidatesNumber; i++ {
		totalScore += float32(finalScore[i])
	}
	halfTotalScore := 0.5 * totalScore
	maxScoreIndex := 0
	maxScore := 0
	var maxVariance float32 = 0

	for i := 0; i < candidatesNumber; i++ {
		if finalScore[i] > maxScore {
			maxScore = finalScore[i]
			maxScoreIndex = i
		}
		if finalScore[i] == maxScore {
			if variance[i] > maxVariance {
				maxScoreIndex = i
			}
		}
	}

	if float32(maxScore) > halfTotalScore {
		winner = rOneCandidateList[maxScoreIndex]
	} else {
		//Round two
		competitorScore := 0
		competitorIndex := 0
		remainNumber := 0
		changeNumber := 0
		finalScore[maxScoreIndex] = 0
		for i := 0; i < candidatesNumber; i++ {
			if finalScore[i] > competitorScore {
				competitorScore = finalScore[i]
				competitorIndex = i
			}
			if finalScore[i] == maxScore {
				if variance[i] > maxVariance {
					competitorIndex = i
				}
			}
		}
		rTwoCandidateList := []shared.ClientID{rOneCandidateList[maxScoreIndex], rOneCandidateList[competitorIndex]}

		var rTwoVotes [][]shared.ClientID
		for i := 0; i < voterNumber; i++ {
			rTwoVotes = append(rTwoVotes, clientMap[e.voterList[i]].VoteForElection(e.roleToElect, rTwoCandidateList))
		}
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
	return winner
}

func (e *Election) instantRunoffResult(clientMap map[shared.ClientID]baseclient.Client) shared.ClientID {
	var winner shared.ClientID
	candidateNumber := len(e.candidateList)
	candidateList := e.candidateList
	totalVotes := e.votes
	maxScore := 0
	maxScoreIndex := 0
	var totalScore float32 = 0
	var halfTotalScore float32 = 0

	for {
		scoreList, variance := e.scoreCalculator(totalVotes, candidateList)

		for i := 0; i < candidateNumber; i++ {
			totalScore += float32(scoreList[i])
		}
		halfTotalScore = 0.5 * totalScore

		for i := 0; i < candidateNumber; i++ {
			if scoreList[i] > maxScore {
				maxScore = scoreList[i]
				maxScoreIndex = i
			}
			if scoreList[i] == maxScore {
				if variance[i] > variance[maxScoreIndex] {
					maxScoreIndex = i
				}
			}
		}
		//keep eliminating the least popular one untill the most popular one has more than half of the total score. 
		if float32(maxScore) > halfTotalScore {
			winner = candidateList[maxScoreIndex]
			break
		}

		minScore := int(totalScore)
		minScoreIndex := 0

		for i := 0; i < len(candidateList); i++ {
			if scoreList[i] < minScore {
				minScore = scoreList[i]
				minScoreIndex = i
			}
			if scoreList[i] == minScore {
				if variance[i] < variance[minScoreIndex] {
					minScoreIndex = i
				}
			}
		}

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
		for i := 1; i < len(e.voterList); i++ {
			totalVotes = append(totalVotes, clientMap[e.voterList[i]].VoteForElection(e.roleToElect, candidateList))
		}
		totalVotes = totalVotes[(len(totalVotes) - len(e.voterList)):]

		//Re initialize parameters
		maxScore = 0
		maxScoreIndex = 0
		totalScore = 0
		halfTotalScore = 0
	}
	return winner
}

func (e *Election) approvalResult() shared.ClientID {
	var winner shared.ClientID
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
	winner = candidateList[maxScoreIndex]
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
