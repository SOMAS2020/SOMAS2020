package voting

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type Election struct {
	roleToElect   baseclient.Role
	votingMethod  ElectionVotingMethod
	islandsToVote []shared.ClientID
	votes         [][]shared.ClientID
}

type ElectionVotingMethod = int

const (
	BordaCount = iota
	Plurality
	Majority
)

// ProposeMotion sets the role to be voted on
func (e *Election) ProposeElection(role baseclient.Role, method ElectionVotingMethod) {
	e.roleToElect = role
	e.votingMethod = method
}

// OpenBallot sets the islands eligible to vote.
func (e *Election) OpenBallot(clientIDs []shared.ClientID) {
	e.islandsToVote = clientIDs
}

// Vote gets votes from eligible islands.
func (e *Election) Vote(clientMap map[shared.ClientID]baseclient.Client) {
	for _, island := range e.islandsToVote {
		e.votes = append(e.votes, clientMap[island].GetVoteForElection(e.roleToElect))
	}
}

// CloseBallot counts the votes received and returns the result.
func (e *Election) CloseBallot() shared.ClientID {

	var result shared.ClientID

	switch e.votingMethod {
	case BordaCount:
		result = e.bordaCountResult()
	case Plurality:
		result = e.pluralityResult()
	case Majority:
		result = e.majorityResult()
	}
	return result
}

func (e *Election) bordaCountResult() shared.ClientID {
	// Implement Borda count winner selection method
	//(may need to modify if methods design is changed)
	var candidatesNumber int = 0
	var islandsNumber int
	var votesLayoutElect map[int][]int
	votesSliceSquare := e.votes
	candidatesNumber = len(e.islandsToVote)
	islandsNumber = len(votesSliceSquare)
	scoreInit := candidatesNumber + 1

	//Transfer e.votes to a preference map with type "int"
	for i := 0; i < islandsNumber; i++ {
		for j := 0; j < candidatesNumber; j++ {
			for k := 0; k < candidatesNumber; k++ {
				if votesSliceSquare[i][j] == e.islandsToVote[k] {
					votesLayoutElect[i][k] = scoreInit
					scoreInit = scoreInit - 1
					break
				}
			}
		}
	}

	//Sort the preference map in order.
	var order []int
	var index []int
	var score []int
	preferenceMap := make(map[int][]int)
	for k, v := range votesLayoutElect {
		j := 0

		for t := 0; t < candidatesNumber; t++ {
			order[t] = v[t]
		}

		for i := 0; i < candidatesNumber; i++ {

			maxlim := 100

			for j = 0; j < candidatesNumber; j++ {
				if maxlim > order[j] {
					maxlim = order[j]
					index[i] = j
				}
			}

			j = index[i]
			order[j] = 101
		}

		itrans := 0
		s := 1
		for i := 0; i < candidatesNumber; i++ {
			itrans = index[i]
			score[itrans] = s
			s++
		}

		preferenceMap[k] = score

	}

	//Calculate the final score for all candidates and ditermine the winner.
	var FinalScore []int
	for _, v := range preferenceMap {
		for i := 0; i < candidatesNumber; i++ {
			FinalScore[i] = FinalScore[i] + v[i]
		}
	}

	maxscore := 0
	var winnerIndex int
	winnerIndex = 0
	for i := 0; i < candidatesNumber; i++ {
		if maxscore < FinalScore[i] {
			maxscore = FinalScore[i]
			winnerIndex = i
		}
	}
	var winner shared.ClientID
	type client shared.ClientID
	winner = e.islandsToVote[winnerIndex]

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
