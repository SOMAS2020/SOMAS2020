package foraging

import (
	"fmt"
	"math"
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func TestDeerUtilityTier(t *testing.T) {

	decay := 0.8 // use an arbitrary decay param
	maxDeer := 4

	var tests = []struct {
		inputR shared.ForageContribution // cumulative resource input from hunt participants
		want   int                       // output tier
	}{
		{0.0, 0},
		{0.99, 0},
		{1.0, 1},
		{1.8, 2},
		{2.45, 3},
		{2.96, 4},
		{1000.0, 4},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("%.3f", tt.inputR)
		t.Run(testname, func(t *testing.T) {
			ans := deerUtilityTier(tt.inputR, uint(maxDeer), decay)
			if ans != uint(tt.want) {
				t.Errorf("got %d, want %d", ans, tt.want)
			}
		})
	}
}

func TestTotalInput(t *testing.T) {
	huntParticipants := map[shared.ClientID]shared.ForageContribution{shared.Team1: 1.0, shared.Team2: 0.9} // arbitrarily chosen for test
	hunt, _ := CreateDeerHunt(huntParticipants)
	ans := hunt.TotalInput()
	if ans != 1.9 {
		t.Errorf("TotalInput() = %.2f; want 1.9", ans)
	}
}

func TestDeerReturn(t *testing.T) {
	params := deerHuntParams{p: 0.95, lam: 1.0}
	avReturn := 0.0
	for i := 1; i <= 1000; i++ { // calculate empirical mean return over 1000 trials
		d := deerReturn(params)
		avReturn = (avReturn*(float64(i)-1) + d) / float64(i)
	}
	expectedReturn := params.p * (1 + 1/params.lam) // theoretical mean based on def of expectation
	if math.Abs(1-expectedReturn/avReturn) > 0.05 {
		t.Errorf("Empirical mean return deviated from theoretical by > 5 percent: got %.3f, want %.3f", avReturn, expectedReturn)
	}
}
