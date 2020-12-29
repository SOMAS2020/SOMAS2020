package voting

import (
	"fmt"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"github.com/SOMAS2020/SOMAS2020/pkg/miscutils"
)

type Election struct {
	roleToElect   baseclient.Role
	votingMethod  ElectionVotingMethod
	islandsToVote []shared.ClientID
	votes         [][]shared.ClientID
}

// ElectionVotingMethod provides enumerated type for selection of voting system to be used
type ElectionVotingMethod int

const (
	BordaCount = iota
	Plurality
	Majority
)

func (e ElectionVotingMethod) String() string {
	strs := [...]string{
		"BordaCount",
		"Plurality",
		"Majority",
	}
	if e >= 0 && int(e) < len(strs) {
		return strs[e]
	}
	return fmt.Sprintf("UNKNOWN ElectionVotingMethod '%v'", int(e))
}

// GoString implements GoStringer
func (e ElectionVotingMethod) GoString() string {
	return e.String()
}

// MarshalText implements TextMarshaler
func (e ElectionVotingMethod) MarshalText() ([]byte, error) {
	return miscutils.MarshalTextForString(e.String())
}

// MarshalJSON implements RawMessage
func (e ElectionVotingMethod) MarshalJSON() ([]byte, error) {
	return miscutils.MarshalJSONForString(e.String())
}

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
	var votesLayoutElect map[int][]int
	votesSliceSquare := e.votes
	candidatesNumber := len(e.islandsToVote)
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
				if votesSliceSquare[i][j] == e.islandsToVote[k] {
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

	//Calculate the final score for all candidates and ditermine the winner.
	finalScore := make([]int, candidatesNumber)
	for _, v := range preferenceMap {
		for i := 0; i < candidatesNumber; i++ {
			finalScore[i] = finalScore[i] + v[i]
		}
	}

	maxscore := 0
	var winnerIndex int
	winnerIndex = 0
	for i := 0; i < candidatesNumber; i++ {
		if maxscore < finalScore[i] {
			maxscore = finalScore[i]
			winnerIndex = i
		}
	}
	var winner shared.ClientID
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
