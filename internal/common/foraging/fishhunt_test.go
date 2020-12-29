package foraging

import (
	"fmt"
	"math"
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

//Checks if the fish utility is correct
func TestFishUtilityTier(t *testing.T) {

	decayFish := 0.8 // use an arbitrary decay param
	maxFish := 6

	var tests = []struct {
		inputF shared.Resources // cumulative resource input from hunt participants
		wantF  int              // output tier
	}{
		// Tiers and coressponding thresholds and cumlative cost
		// Tiers		 			0		1		2		3			4
		// Cumulative cost 			0.0		1.0		1.8		2.44 		2.952		...
		// Incremental cost			0.0		1.0		0.8		0.64		0.512
		{0.0, 0},
		{0.99, 0},
		{1.52, 1},
		{2.1, 2},
		{2.45, 3},
		{2.99, 4},
		{1000.0, 6},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("%.3f", tt.inputF)
		t.Run(testname, func(t *testing.T) {
			ans := utilityTier(tt.inputF, uint(maxFish), decayFish)
			if ans != uint(tt.wantF) {
				t.Errorf("got %d, want %d", ans, tt.wantF)
			}
		})
	}
}

// Checks if total fish input is correct
func TestTotalFishInput(t *testing.T) {
	huntParticipants := map[shared.ClientID]shared.Resources{shared.Team1: 1.0, shared.Team2: 0.9} // arbitrarily chosen for test
	huntF, _ := CreateFishingExpedition(huntParticipants)
	ans := huntF.TotalInput()
	if ans != 1.9 {
		t.Errorf("TotalInput() = %.2f; want 1.9", ans)
	}
}

func TestFishReturn(t *testing.T) {
	params := fishingParams{Mu: 0.9, Sigma: 0.2}
	avReturn := 0.0
	for i := 1; i <= 1000; i++ { // calculate empirical mean return over 1000 trials
		d := fishingReturn(params)
		avReturn = (avReturn*(float64(i)-1) + float64(d)) / float64(i)
	}
	expectedReturn := params.Mu                                 // theoretical mean based on defined expectation
	if math.Abs(float64(1.0-expectedReturn/avReturn)) > 0.012 { // The mean deviation from the theoretical mean with 99.7% confidence
		// (needs to be 3standard deviations out to return error)
		t.Errorf("Empirical mean return deviated from theoretical for > 1.2 percent: got %.5f, want %.5f", avReturn, expectedReturn)
	}
}
