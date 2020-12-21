package vote

import (
	"math/rand"
	"time"
)

func Getvotesforrule(ruleID string, Numofislands int) map[int][]int {
	voteslayout := make(map[int][]int)
	rand.Seed(time.Now().UTC().UnixNano())
	i := 1
	for {
		voteslayout[i] = []int{rand.Intn(100), rand.Intn(100)}
		i++
		if i > Numofislands {
			break
		}
	}
	return voteslayout
}

func VoteRule(ruleID string, Numofislands int) (int, int, bool) {

	//Get votes from islands for the rule changes.
	voteslayout := make(map[int][]int)
	voteslayout = Getvotesforrule(ruleID, Numofislands)

	//Calculate results of each island.
	resultsofallislands := make(map[int][]int)
	j := 1
	for {
		if voteslayout[j][0] > voteslayout[j][1] {
			resultsofallislands[j] = []int{1, 0}
		}
		if voteslayout[j][0] <= voteslayout[j][1] {
			resultsofallislands[j] = []int{0, 1}
		}
		j++
		if j > Numofislands {
			break
		}
	}

	//Announce voting result
	numacc := 0
	numrej := 0
	var ans bool
	for _, v := range resultsofallislands {
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

func Getvotesforelect(Numofislands int) map[int][]int {
	rand.Seed(time.Now().UTC().UnixNano())
	voteslayoute := make(map[int][]int)
	i := 1
	for {
		voteslayoute[i] = []int{rand.Intn(100), rand.Intn(100),
			rand.Intn(100), rand.Intn(100), rand.Intn(100), rand.Intn(100)}
		i++
		if i > Numofislands {
			break
		}
	}

	return voteslayoute
}

func VoteElect(Numofislands int) (int, []int, map[int][]int) {

	//Get votes from each island for the election.
	voteslayoute := make(map[int][]int)
	voteslayoute = Getvotesforelect(Numofislands)

	//Calculate the preference map.
	var order []int
	var index []int
	var score []int
	PreferenceMap := make(map[int][]int)
	for k, v := range voteslayoute {
		t := 0
		i := 0
		j := 0

		for {
			order[t] = v[t]
			t++
			if t > Numofislands-1 {
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
				if j > Numofislands-1 {
					break
				}
			}

			j = index[i]
			order[j] = 101

			i++
			if i > Numofislands-1 {
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
			if i > Numofislands-1 {
				break
			}
		}

		PreferenceMap[k] = score

	}

	//Calculate the final score for six island and ditermine the winner.
	var Finalscore []int
	for _, v := range PreferenceMap {
		i := 0
		for {
			Finalscore[i] = Finalscore[i] + v[i]
			i++
			if i > Numofislands-1 {
				break
			}
		}
	}
	i := 0
	maxscore := 0
	win := 0
	for {
		if maxscore < Finalscore[i] {
			maxscore = Finalscore[i]
			win = i
		}
		i++
		if i > Numofislands-1 {
			break
		}
	}

	win = win + 1

	return win, Finalscore, PreferenceMap
}

func VoteElectJudge(Numofislands int) (int, []int, map[int][]int) {
	win, Finalscore, PreferenceMap := VoteElect(Numofislands)
	return win, Finalscore, PreferenceMap
}

func VoteElectSpeaker(Numofislands int) (int, []int, map[int][]int) {
	win, Finalscore, PreferenceMap := VoteElect(Numofislands)
	return win, Finalscore, PreferenceMap
}

func VoteElectPresident(Numofislands int) (int, []int, map[int][]int) {
	win, Finalscore, PreferenceMap := VoteElect(Numofislands)
	return win, Finalscore, PreferenceMap
}
