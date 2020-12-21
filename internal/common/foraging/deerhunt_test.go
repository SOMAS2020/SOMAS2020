package foraging

import (
	"fmt"
	"math"
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

var huntParticipants = map[shared.ClientID]float64{shared.Team1: 1.0, shared.Team2: 0.9} // arbitrarily chosen for test
var params = DeerHuntParams{P: 0.95, Lam: 1.0}
var hunt = DeerHunt{Participants: huntParticipants, Params: params}

func TestDeerUtilityTier(t *testing.T) {
	dI := []float64{1.0, 0.75, 0.5, 0.25} // define the incremental resource input requirements across tiers

	var tests = []struct {
		inputR float64 // cumulative resource input from hunt participants
		want   int     // output tier
	}{
		{0.0, 0},
		{dI[0] - 0.001, 0},
		{dI[0] + 0.001, 1},
		{dI[0] + dI[1], 2},
		{dI[0] + dI[1] + 0.8*dI[2], 2},
		{dI[0] + dI[1] + dI[2], 3},
		{dI[0] + dI[1] + dI[2] + dI[3], 4},
		{dI[0] + dI[1] + dI[2] + dI[3]*100, len(dI)},
	}
	for i, tt := range tests {
		testname := fmt.Sprintf("%.3f", tt.inputR)
		t.Run(testname, func(t *testing.T) {
			ans := deerUtilityTier(tt.inputR, dI)
			fmt.Println(i)
			if ans != tt.want {
				t.Errorf("got %d, want %d", ans, tt.want)
			}
		})
	}
}

func TestTotalInput(t *testing.T) {
	ans := hunt.TotalInput()
	if ans != 1.9 {
		t.Errorf("TotalInput() = %.2f; want 1.9", ans)
	}
}

func TestDeerReturn(t *testing.T) {
	avReturn := 0.0
	for i := 1; i <= 1000; i++ { // calculate empirical mean return over 1000 trials
		d := deerReturn(params)
		avReturn = (avReturn*(float64(i)-1) + d) / float64(i)
	}
	expectedReturn := params.P * (1 + 1/params.Lam) // theoretical mean based on def of expectation
	if math.Abs(1-expectedReturn/avReturn) > 0.05 {
		t.Errorf("Empirical mean return deviated from theoretical by > 5 percent: got %.3f, want %.3f", avReturn, expectedReturn)
	}
}
