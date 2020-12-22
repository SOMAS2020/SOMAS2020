package vote

import "github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"

//output of this function is number of people for the rule, against the rule,
//and number of people for and againet the rule on each island
func VoteRule(ruleID int, numOfIslands int) (int, int, bool, map[int][]int) {

	//Get votes from islands for the rule changes.
	var clientVoteForRule baseclient.Client
	clientVoteForRule = new(baseclient.BaseClient)
	votesLayoutRule := clientVoteForRule.GetVotesForRule(ruleID, numOfIslands)

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
		if j > numOfIslands {
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
	return numacc, numrej, ans, votesLayoutRule
}

//ouput of this function is winnerID, final score of each candidate,
//and preference list provided by each island
func VoteElect(numOfIslands int) (int, []int, map[int][]int) {

	//Get votes from each island for the election.
	//var votesLayoutElect map[int][]int
	var clientVoteForElect baseclient.Client
	clientVoteForElect = new(baseclient.BaseClient)
	votesLayoutElect := clientVoteForElect.GetVotesForElect(numOfIslands)

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
			if t > numOfIslands-1 {
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
				if j > numOfIslands-1 {
					break
				}
			}

			j = index[i]
			order[j] = 101

			i++
			if i > numOfIslands-1 {
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
			if i > numOfIslands-1 {
				break
			}
		}

		preferenceMap[k] = score

	}

	//Calculate the final score for six island and ditermine the winner.
	var FinalScore []int
	for _, v := range preferenceMap {
		i := 0
		for {
			FinalScore[i] = FinalScore[i] + v[i]
			i++
			if i > numOfIslands-1 {
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
		if i > numOfIslands-1 {
			break
		}
	}

	win = win + 1

	return win, FinalScore, preferenceMap
}

//vote for judge
func VoteElectJudge(numOfIslands int) (int, []int, map[int][]int) {
	win, FinalScore, preferenceMap := VoteElect(numOfIslands)
	return win, FinalScore, preferenceMap
}

//vote for speaker
func VoteElectSpeaker(numOfIslands int) (int, []int, map[int][]int) {
	win, FinalScore, preferenceMap := VoteElect(numOfIslands)
	return win, FinalScore, preferenceMap
}

//vote for president
func VoteElectPresident(numOfIslands int) (int, []int, map[int][]int) {
	win, FinalScore, preferenceMap := VoteElect(numOfIslands)
	return win, FinalScore, preferenceMap
}
