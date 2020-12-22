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
		inputF float64 // cumulative resource input from hunt participants
		wantF  int     // output tier
	}{
		// Tiers and coressponding thresholds and cumlative cost
		// Tiers		 					0			1			2			3				4				5					6
		// Thresholds 					0.0		0.8		0.64	0.512		0.4096	0.32768		0.262144
		// Cumlative cost 			0.0		0.8		1.44	1.952		2.3616	2.68928		2.951424
		{0.0, 0},
		{0.1, 0},
		{0.99, 1},
		{1.52, 2},
		{2.1, 3},
		{2.42, 4},
		{2.9, 5},
		{1000.0, 6},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("%.3f", tt.inputF)
		t.Run(testname, func(t *testing.T) {
			ans := fishUtilityTier(tt.inputF, uint(maxFish), decayFish)
			if ans != uint(tt.wantF) {
				t.Errorf("got %d, want %d", ans, tt.wantF)
			}
		})
	}
}

// Checks if total fish input is correct
func TestTotalFishInput(t *testing.T) {
	huntParticipants := map[shared.ClientID]float64{shared.Team1: 1.0, shared.Team2: 0.9} // arbitrarily chosen for test
	huntF, _ := CreateFishHunt(huntParticipants)
	ans := huntF.TotalInput()
	if ans != 1.9 {
		t.Errorf("TotalInput() = %.2f; want 1.9", ans)
	}
}

func TestFishReturn(t *testing.T) {
	params := FishHuntParams{Mu: 0.9, Sigma: 0.2}
	avReturn := 0.0
	for i := 1; i <= 1000; i++ { // calculate empirical mean return over 1000 trials
		d := fishReturn(params)
		avReturn = (avReturn*(float64(i)-1) + d) / float64(i)
	}
	expectedReturn := params.Mu // theoretical mean based on def of expectation
	if math.Abs(1-expectedReturn/avReturn) > 0.006 {
		t.Errorf("Empirical mean return deviated from theoretical by > 5 percent: got %.3f, want %.3f", avReturn, expectedReturn)
	}
}
