package vote

import (
	"math/rand"
	"time"
)

func GetVotesForRule(ruleID string, islandID int) map[int][]int {
	votes_layout := make(map[int][]int)
	rand.Seed(time.Now().UTC().UnixNano())
	i := 1
	for {
		votes_layout[i] = []int{rand.Intn(100), rand.Intn(100)}
		i++
		if i > 6 {
			break
		}
	}
	return votes_layout
}

func VoteRule(ruleID string, islandID int) (int, int, bool) {

	//Get votes from islands for the rule changes.
	votes_layout := make(map[int][]int)
	votes_layout = Get_votes_for_rule(ruleID, islandID)

	//Calculate results of each island.
	results_of_all_islands := make(map[int][]int)
	j := 1
	for {
		if votes_layout[j][0] > votes_layout[j][1] {
			results_of_all_islands[j] = []int{1, 0}
		}
		if votes_layout[j][0] <= votes_layout[j][1] {
			results_of_all_islands[j] = []int{0, 1}
		}
		j++
		if j > 6 {
			break
		}
	}

	//Announce voting result
	num_acc := 0
	num_rej := 0
	var ans bool
	for _, v := range results_of_all_islands {
		if v[0] == 1 {
			num_acc++
		}
		if v[1] == 1 {
			num_rej++
		}
	}
	if num_acc > num_rej {
		ans = true
	}
	if num_acc <= num_rej {
		ans = false
	}
	return num_acc, num_rej, ans
}

func Get_votes_for_elect() map[int][]int {
	rand.Seed(time.Now().UTC().UnixNano())
	votes_layout_e := make(map[int][]int)
	i := 1
	for {
		votes_layout_e[i] = []int{rand.Intn(100), rand.Intn(100),
			rand.Intn(100), rand.Intn(100), rand.Intn(100), rand.Intn(100)}
		i++
		if i > 6 {
			break
		}
	}

	return votes_layout_e
}

func Vote_Elect() (int, [6]int, map[int][]int) {

	//Get votes from each island for the election.
	votes_layout_e := make(map[int][]int)
	votes_layout_e = Get_votes_for_elect()

	//Calculate the preference map.
	var order [6]int
	var index [6]int
	var score [6]int
	PreferenceMap := make(map[int][]int)
	for k, v := range votes_layout_e {
		t := 0
		i := 0
		j := 0

		for {
			order[t] = v[t]
			t++
			if t > 5 {
				break
			}
		}

		for {

			p := 100
			j = 0

			for {
				if p > order[j] {
					p = order[j]
					index[i] = j
				}

				j++
				if j > 5 {
					break
				}
			}

			j = index[i]
			order[j] = 101

			i++
			if i > 5 {
				break
			}
		}

		o := 0
		r := 0
		s := 1
		for {
			r = index[o]
			score[r] = s
			s++
			o++
			if o > 5 {
				break
			}
		}

		PreferenceMap[k] = []int{score[0], score[1], score[2],
			score[3], score[4], score[5]}

	}

	//Calculate the final score for six island and ditermine the winner.
	var Final_score = [6]int{0, 0, 0, 0, 0, 0}
	for _, v := range PreferenceMap {
		Final_score[0] = Final_score[0] + v[0]
		Final_score[1] = Final_score[1] + v[1]
		Final_score[2] = Final_score[2] + v[2]
		Final_score[3] = Final_score[3] + v[3]
		Final_score[4] = Final_score[4] + v[4]
		Final_score[5] = Final_score[5] + v[5]
	}
	h := 0
	l := 0
	win := 0
	for {
		if l < Final_score[h] {
			l = Final_score[h]
			win = h
		}
		h++
		if h > 5 {
			break
		}
	}

	win = win + 1

	return win, Final_score, PreferenceMap
}

func VoteElectJudge() (int, [6]int, map[int][]int) {
	win, Final_score, PreferenceMap := Vote_Elect()
	return win, Final_score, PreferenceMap
}

func VoteElectSpeaker() (int, [6]int, map[int][]int) {
	win, Final_score, PreferenceMap := Vote_Elect()
	return win, Final_score, PreferenceMap
}

func VoteElectPresident() (int, [6]int, map[int][]int) {
	win, Final_score, PreferenceMap := Vote_Elect()
	return win, Final_score, PreferenceMap
}
