package team2

import (
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// Testing logic, sue me for bad practice time is of essence haha
func TestOtherHunters(t *testing.T) {
	foragingReturnsHist := map[shared.ClientID][]int{
		shared.Team1: {0, 0, 1, 0, 1},
		shared.Team2: {1, 1, 1, 1, 1},
		shared.Team3: {1, 1, 1, 1, 1},
	}
	HuntNum := 0.00                       //this the average number of likely hunters
	for id := range foragingReturnsHist { //loop through every agent
		for _, forageInfo := range foragingReturnsHist[id] { //loop through the agents array and add their average to HuntNum
			HuntNum += float64(forageInfo) / float64(len(foragingReturnsHist[id])) //add the agents decision to HuntNum and then average
		}
	}
	print(HuntNum)
	if HuntNum != 2.4 {
		t.Error("Otherhunters does calculate the average number of hunters correctly")
	}
}

func TestMinInt(t *testing.T) {
	bufferLen := MinInt(0, MinInt(5, 3))
	print("Bufferlen %v", bufferLen)
	if bufferLen != 0 {
		t.Error()
	}
}
