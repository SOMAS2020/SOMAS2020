package team5

import (
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func TestForageHistorySize(t *testing.T) {

	cases := []struct {
		name              string
		expectedVal       uint
		forageHistoryTest forageHistory
	}{
		{
			name:        "Basic test",
			expectedVal: 6,
			forageHistoryTest: forageHistory{
				shared.DeerForageType: {
					{team: 4, turn: 1, input: shared.Resources(20), output: shared.Resources(1)},
					{team: 4, turn: 16, input: shared.Resources(20), output: shared.Resources(1)},
					{team: 4, turn: 12, input: shared.Resources(20), output: shared.Resources(1)},
					{team: 4, turn: 11, input: shared.Resources(20), output: shared.Resources(1)},
					{team: 4, turn: 18, input: shared.Resources(20), output: shared.Resources(1)},
				}, // End of Shared.DeerForageType data
				shared.FishForageType: {
					{team: 5, turn: 11, input: shared.Resources(20), output: shared.Resources(1)},
				}, // End of Shared.FishForageType data
			}, // End of Foraging History
		}, // End of Basic Test

		{
			name:        "Semi Filled Test",
			expectedVal: 2,
			forageHistoryTest: forageHistory{
				shared.DeerForageType: {
					{team: 4, turn: 10, input: shared.Resources(1), output: shared.Resources(0)},
				}, // End of Shared.DeerForageType data
				shared.FishForageType: {
					{team: 5, turn: 11, input: shared.Resources(0), output: shared.Resources(1)},
				}, // End of Shared.FishForageType data
			}, // End of Foraging History
		}, // End of Basic

		{
			name:        "Incorrect Type Test",
			expectedVal: 3,
			forageHistoryTest: forageHistory{
				shared.ForageType(-1): {
					{team: 5, turn: 10, input: shared.Resources(20), output: shared.Resources(1)},
					{team: 5, turn: 10, input: shared.Resources(20), output: shared.Resources(1)},
				}, // End of Shared.DeerForageType data
				shared.FishForageType: {
					{team: 1, turn: 11, input: shared.Resources(1), output: shared.Resources(20)},
				}, // End of Shared.FishForageType data
			}, // End of Foraging History
		}, // End of Basic

	} // End of test data
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c.forageHistory = tc.forageHistoryTest
			ans := c.forageHistorySize()
			if ans != tc.expectedVal {
				t.Errorf("Expected final transgressions to be %v got %v", tc.expectedVal, ans)
			}
		})
	}
}
