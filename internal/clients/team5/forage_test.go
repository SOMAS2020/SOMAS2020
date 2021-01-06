package team5

import (
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// func TestFishUtilityTier(t *testing.T) {

// 	decayFish := 0.8 // use an arbitrary decay param
// 	maxFish := 6

// 	var tests = []struct {
// 		inputF shared.Resources // cumulative resource input from hunt participants
// 		wantF  int              // output tier
// 	}{
// 		// Tiers and coressponding thresholds and cumlative cost
// 		// Tiers		 			0		1		2		3			4
// 		// Cumulative cost 			0.0		1.0		1.8		2.44 		2.952		...
// 		// Incremental cost			0.0		1.0		0.8		0.64		0.512
// 		{0.0, 0},
// 		{0.99, 0},
// 		{1.52, 1},
// 		{2.1, 2},
// 		{2.45, 3},
// 		{2.99, 4},
// 		{1000.0, 6},
// 	}
// 	for _, tt := range tests {
// 		testname := fmt.Sprintf("%.3f", tt.inputF)
// 		t.Run(testname, func(t *testing.T) {
// 			ans := utilityTier(tt.inputF, uint(maxFish), decayFish, 1)
// 			if ans != uint(tt.wantF) {
// 				t.Errorf("got %d, want %d", ans, tt.wantF)
// 			}
// 		})
// 	}
// }

func TestBestHistoryForaging(t *testing.T) {
	// Check foraging method
	c := createClient()
	cases := []struct {
		name              string
		expectedVal       shared.ForageType
		forageHistoryTest forageHistory
	}{
		{
			name:        "Basic test",
			expectedVal: shared.ForageType(-1), // Nul
			forageHistoryTest: map[shared.ForageType][]forageOutcome{
				shared.DeerForageType: {
					{team: 5, turn: 1, input: shared.Resources(20), output: shared.Resources(1)},
				}, // End of Shared.DeerForageType data
				shared.FishForageType: {
					{team: 5, turn: 2, input: shared.Resources(20), output: shared.Resources(1)},
				}, // End of Shared.FishForageType data
			}, // End of Foraging History
		}, // End of Basic Test

		{
			name:        "Deer Test",
			expectedVal: shared.ForageType(0), // Deer
			forageHistoryTest: map[shared.ForageType][]forageOutcome{
				shared.DeerForageType: {
					{team: 5, turn: 1, input: shared.Resources(1), output: shared.Resources(20)},
				}, // End of Shared.DeerForageType data
				shared.FishForageType: {
					{team: 5, turn: 2, input: shared.Resources(20), output: shared.Resources(1)},
				}, // End of Shared.FishForageType data
			}, // End of Foraging History
		}, // End of Basic

		{
			name:        "Fish Test",
			expectedVal: shared.ForageType(1), // Fish
			forageHistoryTest: map[shared.ForageType][]forageOutcome{
				shared.DeerForageType: {
					{team: 5, turn: 1, input: shared.Resources(20), output: shared.Resources(1)},
				}, // End of Shared.DeerForageType data
				shared.FishForageType: {
					{team: 5, turn: 2, input: shared.Resources(1), output: shared.Resources(1)},
				}, // End of Shared.FishForageType data
			}, // End of Foraging History
		}, // End of Basic

	} // End of test data
	for _, tc := range cases {
		foragingHistory := tc.forageHistoryTest
		t.Run(tc.name, func(t *testing.T) {
			ans := c.bestHistoryForaging(foragingHistory)
			if ans != tc.expectedVal {
				t.Errorf("Expected final transgressions to be %v got %v", tc.expectedVal, ans)
			}
		})
	}
}
