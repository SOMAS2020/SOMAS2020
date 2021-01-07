package team5

import (
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func TestBestHistoryForaging(t *testing.T) {
	// Check foraging method
	c := MakeTestClient(gamestate.ClientGameState{
		Turn: 10,
		ClientInfo: gamestate.ClientInfo{
			LifeStatus: shared.Critical,
		},
	})
	cases := []struct {
		name              string
		expectedVal       shared.ForageType
		forageHistoryTest forageHistory
	}{
		{
			name:        "Basic test",
			expectedVal: shared.ForageType(-1), // Nul
			forageHistoryTest: forageHistory{
				shared.DeerForageType: {
					{team: 4, turn: 10, input: shared.Resources(20), output: shared.Resources(1)},
				}, // End of Shared.DeerForageType data
				shared.FishForageType: {
					{team: 5, turn: 11, input: shared.Resources(20), output: shared.Resources(1)},
				}, // End of Shared.FishForageType data
			}, // End of Foraging History
		}, // End of Basic Test

		{
			name:        "Deer Test",
			expectedVal: shared.ForageType(0), // Deer
			forageHistoryTest: forageHistory{
				shared.DeerForageType: {
					{team: 4, turn: 10, input: shared.Resources(1), output: shared.Resources(20)},
				}, // End of Shared.DeerForageType data
				shared.FishForageType: {
					{team: 5, turn: 11, input: shared.Resources(20), output: shared.Resources(1)},
				}, // End of Shared.FishForageType data
			}, // End of Foraging History
		}, // End of Basic

		{
			name:        "Fish Test",
			expectedVal: shared.ForageType(1), // Fish
			forageHistoryTest: forageHistory{
				shared.DeerForageType: {
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
			ans := c.bestHistoryForaging(tc.forageHistoryTest)
			if ans != tc.expectedVal {
				t.Errorf("Expected final transgressions to be %v got %v", tc.expectedVal, ans)
			}
		})
	}
}

func TestInitialForage(t *testing.T) {
	// Check foraging method
	c := MakeTestClient(gamestate.ClientGameState{
		Turn: 1,
		ClientInfo: gamestate.ClientInfo{
			Resources: 100,
		},
	})

	cases := []struct {
		name                    string
		expectedVal             shared.ForageType
		MinimumForagePercentage float64
		NormalForagePercentage  float64
		JBForagePercentage      float64
	}{
		{
			name:                    "Test",
			MinimumForagePercentage: 0.05,
			NormalForagePercentage:  0.10,
			JBForagePercentage:      0.15,
		}, // End of Basic Test
		{
			name:                    "Imperial Student Test",
			MinimumForagePercentage: 0.10,
			NormalForagePercentage:  0.20,
			JBForagePercentage:      0.40,
		}, // End of Basic Test

	} // End of test data

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c.config.MinimumForagePercentage = tc.MinimumForagePercentage
			c.config.NormalForagePercentage = tc.NormalForagePercentage
			c.config.JBForagePercentage = tc.JBForagePercentage

			ans := c.InitialForage()
			if ans.Contribution < shared.Resources(tc.MinimumForagePercentage)*ans.Contribution &&
				ans.Contribution > shared.Resources(tc.NormalForagePercentage)*ans.Contribution {
				t.Errorf("Expected final transgressions to be between %v and %v got %v",
					tc.MinimumForagePercentage, tc.NormalForagePercentage, ans)
			}
			if !(ans.Type == shared.DeerForageType || ans.Type == shared.FishForageType) {
				t.Errorf("Expected final transgressions to be %v or %v got %v",
					shared.DeerForageType, shared.FishForageType, ans)
			}
		})
	}

}

func TestForageHistorySize(t *testing.T) {
	// Check foraging method
	c := MakeTestClient(gamestate.ClientGameState{
		Turn: 10,
		ClientInfo: gamestate.ClientInfo{
			LifeStatus: shared.Critical,
		},
	})

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
					{team: 4, turn: 10, input: shared.Resources(20), output: shared.Resources(1)},
					{team: 4, turn: 10, input: shared.Resources(20), output: shared.Resources(1)},
					{team: 4, turn: 10, input: shared.Resources(20), output: shared.Resources(1)},
					{team: 4, turn: 10, input: shared.Resources(20), output: shared.Resources(1)},
					{team: 4, turn: 10, input: shared.Resources(20), output: shared.Resources(1)},
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
					{team: 5, turn: 11, input: shared.Resources(20), output: shared.Resources(1)},
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

func TestForageUpdate(t *testing.T) {
	// Check foraging method
	c := MakeTestClient(gamestate.ClientGameState{
		Turn: 1,
		ClientInfo: gamestate.ClientInfo{
			LifeStatus: shared.Critical,
		},
	})

	cases := []struct {
		name           string
		forageDecision shared.ForageDecision
		outputFunction shared.Resources
		expectedVal    forageHistory

		foragetype shared.ForageType
		team       shared.ClientID
		turn       uint
		input      shared.Resources
		output     shared.Resources
	}{
		{
			name:           "Basic test",
			forageDecision: shared.ForageDecision{Type: shared.DeerForageType, Contribution: 1},
			outputFunction: 20,

			foragetype: shared.DeerForageType,
			turn:       1,
			team:       shared.Team5,
			input:      shared.Resources(1),
			output:     shared.Resources(20),
		}, // End of Basic Test

		{
			name:           "Basic test",
			forageDecision: shared.ForageDecision{Type: shared.FishForageType, Contribution: 15},
			outputFunction: 200,

			foragetype: shared.FishForageType,
			turn:       1,
			team:       shared.Team5,
			input:      shared.Resources(15),
			output:     shared.Resources(200),
		}, // End of Basic Test

	}

	forageHistory := forageHistory{}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			c.ForageUpdate(tc.forageDecision, tc.output)
			for forageType, FOutcome := range forageHistory { // For the whole foraging history
				for _, returns := range FOutcome {
					if !(forageType == tc.foragetype && //Deer Hunters
						returns.team == tc.team && // Last turn
						returns.input != tc.input &&
						returns.output != tc.output) { // Not including us
						t.Errorf("Someting Wong")
					}
				}
			}
		})
	}
}
