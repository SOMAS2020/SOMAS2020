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
	i := 0
	for {
		j := 0
		for {
			k := 0
			for {
				if votesSliceSquare[i][j] == e.islandsToVote[k] {
					votesLayoutElect[i][k] = scoreInit
					scoreInit = scoreInit - 1
					break
				}
				k++
				if k > candidatesNumber-1 {
					break
				}
			}
			j++
			if j > candidatesNumber-1 {
				break
			}
		}
		i++
		if i > islandsNumber-1 {
			break
		}
	}

	//Calculate the preference map.
	var order []int
	var index []int
	var score []int
	preferenceMap := make(map[int][]int)
	for k, v := range votesLayoutElect {
		t := 0
		i := 0
		j := 0

		for {
			order[t] = v[t]
			t++
			if t > candidatesNumber-1 {
				break
			}
		}

		for {

			maxlim := 100
			j = 0

			for {
				if maxlim > order[j] {
					maxlim = order[j]
					index[i] = j
				}

				j++
				if j > candidatesNumber-1 {
					break
				}
			}

			j = index[i]
			order[j] = 101

			i++
			if i > candidatesNumber-1 {
				break
			}
		}

		i = 0
		itrans := 0
		s := 1
		for {
			itrans = index[i]
			score[itrans] = s
			s++
			i++
			if i > candidatesNumber-1 {
				break
			}
		}

		preferenceMap[k] = score

	}

	//Calculate the final score for all candidates and ditermine the winner.
	var FinalScore []int
	for _, v := range preferenceMap {
		i := 0
		for {
			FinalScore[i] = FinalScore[i] + v[i]
			i++
			if i > candidatesNumber-1 {
				break
			}
		}
	}
	i = 0
	maxscore := 0
	var winnerIndex int
	winnerIndex = 0
	for {
		if maxscore < FinalScore[i] {
			maxscore = FinalScore[i]
			winnerIndex = i
		}
		i++
		if i > candidatesNumber-1 {
			break
		}
	}
	var winner shared.ClientID
	type client shared.ClientID
	winner = shared.TeamIDs[winnerIndex]

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
