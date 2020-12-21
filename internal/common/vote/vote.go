package vote

import (
	"math/rand"
	"time"
)

func GetVotesForRule(ruleID string, NumOFIslands int) map[int][]int {
	votesLayoutRule := make(map[int][]int)
	rand.Seed(time.Now().UTC().UnixNano())
	i := 1
	for {
		votesLayoutRule[i] = []int{rand.Intn(100), rand.Intn(100)}
		i++
		if i > NumOFIslands {
			break
		}
	}
	return votesLayoutRule
}

func VoteRule(ruleID string, NumOFIslands int) (int, int, bool) {

	//Get votes from islands for the rule changes.
	var votesLayoutRule map[int][]int
	votesLayoutRule = GetVotesForRule(ruleID, NumOFIslands)

	//Calculate results of each island.
	resultsOfAllIslands := make(map[int][]int)
	j := 1
	for {
		if votesLayoutRule[j][0] > votesLayoutRule[j][1] {
			resultsOfAllIslands[j] = []int{1, 0}
		}
		if votesLayoutRule[j][0] <= votesLayoutRule[j][1] {
			resultsOfAllIslands[j] = []int{0, 1}
		}
		j++
		if j > NumOFIslands {
			break
		}
	}

	//Announce voting result
	numacc := 0
	numrej := 0
	var ans bool
	for _, v := range resultsOfAllIslands {
		if v[0] == 1 {
			numacc++
		}
		if v[1] == 1 {
			numrej++
		}
	}
	if numacc > numrej {
		ans = true
	}
	if numacc <= numrej {
		ans = false
	}
	return numacc, numrej, ans
}

func GetVotesForElect(NumOFIslands int) map[int][]int {
	rand.Seed(time.Now().UTC().UnixNano())
	votesLayoutElect := make(map[int][]int)
	i := 1
	for {
		votesLayoutElect[i] = []int{rand.Intn(100), rand.Intn(100),
			rand.Intn(100), rand.Intn(100), rand.Intn(100), rand.Intn(100)}
		i++
		if i > NumOFIslands {
			break
		}
	}

	return votesLayoutElect
}

func VoteElect(NumOFIslands int) (int, []int, map[int][]int) {

	//Get votes from each island for the election.
	var votesLayoutElect map[int][]int
	votesLayoutElect = GetVotesForElect(NumOFIslands)

	//Calculate the preference map.
	var order []int
	var index []int
	var score []int
	PreferenceMap := make(map[int][]int)
	for k, v := range votesLayoutElect {
		t := 0
		i := 0
		j := 0

		for {
			order[t] = v[t]
			t++
			if t > NumOFIslands-1 {
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
				if j > NumOFIslands-1 {
					break
				}
			}

			j = index[i]
			order[j] = 101

			i++
			if i > NumOFIslands-1 {
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
			if i > NumOFIslands-1 {
				break
			}
		}

		PreferenceMap[k] = score

	}

	//Calculate the final score for six island and ditermine the winner.
	var FinalScore []int
	for _, v := range PreferenceMap {
		i := 0
		for {
			FinalScore[i] = FinalScore[i] + v[i]
			i++
			if i > NumOFIslands-1 {
				break
			}
		}
	}
	i := 0
	maxscore := 0
	win := 0
	for {
		if maxscore < FinalScore[i] {
			maxscore = FinalScore[i]
			win = i
		}
		i++
		if i > NumOFIslands-1 {
			break
		}
	}

	win = win + 1

	return win, FinalScore, PreferenceMap
}

func VoteElectJudge(NumOFIslands int) (int, []int, map[int][]int) {
	win, FinalScore, PreferenceMap := VoteElect(NumOFIslands)
	return win, FinalScore, PreferenceMap
}

func VoteElectSpeaker(NumOFIslands int) (int, []int, map[int][]int) {
	win, FinalScore, PreferenceMap := VoteElect(NumOFIslands)
	return win, FinalScore, PreferenceMap
}

func VoteElectPresident(NumOFIslands int) (int, []int, map[int][]int) {
	win, FinalScore, PreferenceMap := VoteElect(NumOFIslands)
	return win, FinalScore, PreferenceMap
}
