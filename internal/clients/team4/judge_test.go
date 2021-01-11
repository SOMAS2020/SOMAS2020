package team4

import (
	"math/rand"
	"reflect"
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func makeHistory(clientVars map[shared.ClientID][]rules.VariableFieldName, lawfulness map[shared.ClientID]float64) (history []shared.Accountability, expected map[shared.ClientID]judgeHistoryInfo) {
	history = []shared.Accountability{}
	expected = map[shared.ClientID]judgeHistoryInfo{}
	for client, vars := range clientVars {
		pairs := []rules.VariableValuePair{}
		for _, v := range vars {
			newVal := rand.ExpFloat64()
			newPair := makeSingleVar(v, newVal)
			newAcc := shared.Accountability{
				ClientID: client,
				Pairs:    []rules.VariableValuePair{newPair},
			}
			pairs = append(pairs, newPair)
			history = append(history, newAcc)
		}
		got, ok := buildHistoryInfo(pairs)
		if ok {
			got.LawfulRatio = lawfulness[client]
			expected[client] = got
		}
	}

	return history, expected
}

func TestSaveHistoryInfo(t *testing.T) {
	cases := []struct {
		name            string
		clientVariables map[shared.ClientID][]rules.VariableFieldName
		truthfulness    map[shared.ClientID]float64
		turn            uint
		reps            uint
	}{
		{
			name: "Single client simple",
			clientVariables: map[shared.ClientID][]rules.VariableFieldName{
				shared.Team1: {
					rules.IslandReportedPrivateResources,
					rules.IslandActualPrivateResources,
					rules.IslandTaxContribution,
					rules.ExpectedTaxContribution,
					rules.IslandAllocation,
					rules.ExpectedAllocation,
				},
			},
			truthfulness: map[shared.ClientID]float64{
				shared.Team1: 0.6,
			},
			turn: 12,
			reps: 1,
		},
		{
			name: "Multiple clients simple",
			clientVariables: map[shared.ClientID][]rules.VariableFieldName{
				shared.Team1: {
					rules.IslandReportedPrivateResources,
					rules.IslandActualPrivateResources,
					rules.IslandTaxContribution,
					rules.ExpectedTaxContribution,
					rules.IslandAllocation,
					rules.ExpectedAllocation,
				},
				shared.Team2: {
					rules.IslandReportedPrivateResources,
					rules.IslandActualPrivateResources,
					rules.IslandTaxContribution,
					rules.ExpectedTaxContribution,
					rules.IslandAllocation,
					rules.ExpectedAllocation,
				},
				shared.Team3: {
					rules.IslandReportedPrivateResources,
					rules.IslandActualPrivateResources,
					rules.IslandTaxContribution,
					rules.ExpectedTaxContribution,
					rules.IslandAllocation,
					rules.ExpectedAllocation,
				},
			},
			truthfulness: map[shared.ClientID]float64{
				shared.Team1: 0.99,
				shared.Team2: 1,
				shared.Team3: 0.2,
			},
			turn: 1,
			reps: 1,
		},
		{
			name: "Multiple clients mixed",
			clientVariables: map[shared.ClientID][]rules.VariableFieldName{
				shared.Team1: {
					rules.IslandReportedPrivateResources,
					rules.IslandActualPrivateResources,
					rules.IslandTaxContribution,
					rules.ExpectedTaxContribution,
					rules.IslandAllocation,
					rules.ExpectedAllocation,
				},
				shared.Team2: {
					rules.IslandReportedPrivateResources,
					rules.IslandActualPrivateResources,
					rules.IslandTaxContribution,
					rules.ExpectedTaxContribution,
					rules.IslandAllocation,
					rules.ExpectedAllocation,
					rules.AllocationMade,
					rules.ConstSanctionAmount,
					rules.HasIslandReportPrivateResources,
				},
				shared.Team3: {
					rules.IslandReportedPrivateResources,
					rules.IslandActualPrivateResources,
					rules.IslandTaxContribution,
					rules.ExpectedTaxContribution,
					rules.IslandAllocation,
				},
			},
			truthfulness: map[shared.ClientID]float64{
				shared.Team1: 0.99,
				shared.Team2: 1,
				shared.Team3: 0.2,
			},
			turn: 13,
			reps: 1,
		},
		{
			name: "Multiple clients mixed repeat",
			clientVariables: map[shared.ClientID][]rules.VariableFieldName{
				shared.Team1: {
					rules.IslandReportedPrivateResources,
					rules.IslandActualPrivateResources,
					rules.IslandTaxContribution,
					rules.ExpectedTaxContribution,
					rules.IslandAllocation,
					rules.ExpectedAllocation,
				},
				shared.Team2: {
					rules.IslandReportedPrivateResources,
					rules.IslandActualPrivateResources,
					rules.IslandTaxContribution,
					rules.ExpectedTaxContribution,
					rules.IslandAllocation,
					rules.ExpectedAllocation,
					rules.AllocationMade,
					rules.ConstSanctionAmount,
					rules.HasIslandReportPrivateResources,
				},
				shared.Team3: {
					rules.IslandReportedPrivateResources,
					rules.IslandActualPrivateResources,
					rules.IslandTaxContribution,
					rules.ExpectedTaxContribution,
					rules.IslandAllocation,
				},
			},
			truthfulness: map[shared.ClientID]float64{
				shared.Team1: 0.99,
				shared.Team2: 1,
				shared.Team3: 0.5,
			},
			turn: 13,
			reps: 3,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			testClient := newClientInternal(shared.Team4, forTesting)
			j := testClient.clientJudge

			wholeHistory := map[uint]map[shared.ClientID]judgeHistoryInfo{}

			for i := uint(0); i < tc.reps; i++ {
				fakeHistory, expected := makeHistory(tc.clientVariables, tc.truthfulness)
				turn := tc.turn + i
				wholeHistory[turn] = expected

				j.saveHistoryInfo(&fakeHistory, &tc.truthfulness, turn)

				clientHistory := *testClient.savedHistory

				if !reflect.DeepEqual(expected, clientHistory.history[turn]) {
					t.Errorf("Single history failed. expected %v,\n got %v", expected, clientHistory.history[turn])
				}

				if !clientHistory.updated {
					t.Errorf("Single history failed. History was not updated")
				}

				if clientHistory.updatedTurn != turn {
					t.Errorf("Single history failed. Expected updated turn: %v, got turn: %v", turn, clientHistory.updatedTurn)
				}
			}

			if !reflect.DeepEqual(wholeHistory, testClient.savedHistory.history) {
				t.Errorf("Whole history comparison failed. Saved history: %v", *testClient.savedHistory)
			}

		})
	}
}

func TestCallPresidentElection(t *testing.T) {
	cases := []struct {
		name               string
		monitoring         shared.MonitorResult
		turnsInPower       int
		termLength         uint
		electionRuleInPlay bool
		expectedElection   bool
	}{
		{
			name:               "no conditions",
			monitoring:         shared.MonitorResult{Performed: false, Result: false},
			turnsInPower:       1,
			termLength:         4,
			electionRuleInPlay: false,
			expectedElection:   false,
		},
		{
			name:               "term length exceeded. no rule",
			monitoring:         shared.MonitorResult{Performed: false, Result: false},
			turnsInPower:       2,
			termLength:         1,
			electionRuleInPlay: false,
			expectedElection:   false,
		},
		{
			name:               "term length exceeded. rule in play",
			monitoring:         shared.MonitorResult{Performed: false, Result: false},
			turnsInPower:       7,
			termLength:         5,
			electionRuleInPlay: true,
			expectedElection:   true,
		},
		{
			name:               "termLength=turnsInPower. rule in play",
			monitoring:         shared.MonitorResult{Performed: false, Result: false},
			turnsInPower:       5,
			termLength:         5,
			electionRuleInPlay: true,
			expectedElection:   false,
		},
		{
			name:               "monitoring done",
			monitoring:         shared.MonitorResult{Performed: true, Result: true},
			turnsInPower:       2,
			termLength:         3,
			electionRuleInPlay: true,
			expectedElection:   false,
		},
		{
			name:               "monitoring done and cheated",
			monitoring:         shared.MonitorResult{Performed: true, Result: false},
			turnsInPower:       5,
			termLength:         7,
			electionRuleInPlay: true,
			expectedElection:   true,
		},
		{
			name:               "monitoring done and cheated. term ended",
			monitoring:         shared.MonitorResult{Performed: true, Result: false},
			turnsInPower:       5,
			termLength:         4,
			electionRuleInPlay: true,
			expectedElection:   true,
		},
	}

	allTeams := []shared.ClientID{}
	for _, client := range shared.TeamIDs {
		allTeams = append(allTeams, client)
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			testServer := fakeServerHandle{
				TermLengths: map[shared.Role]uint{
					shared.President: tc.termLength,
				},
				ElectionRuleInPlay: tc.electionRuleInPlay,
			}

			testClient := newClientInternal(shared.Team4, honest)

			testClient.Initialise(testServer)

			j := testClient.GetClientJudgePointer()

			got := j.CallPresidentElection(tc.monitoring, tc.turnsInPower, allTeams)

			if got.HoldElection != tc.expectedElection {
				t.Errorf("Expected holdElection: %v. Got holdElection: %v", tc.expectedElection, got.HoldElection)
			}
		})
	}
}

func TestGetPardonedIslands(t *testing.T) {
	cases := []struct {
		name      string
		sanctions map[int][]shared.Sanction
		pardons   map[int][]bool
		trust     map[shared.ClientID]float64
	}{
		{
			name:      "empty sanctions",
			sanctions: make(map[int][]shared.Sanction),
			pardons:   make(map[int][]bool),
			trust:     make(map[shared.ClientID]float64),
		},
		{
			name: "no sanction",
			sanctions: map[int][]shared.Sanction{
				0: {
					{
						ClientID:     shared.Team1,
						SanctionTier: shared.NoSanction,
						TurnsLeft:    5,
					},
					{
						ClientID:     shared.Team2,
						SanctionTier: shared.NoSanction,
						TurnsLeft:    5,
					},
				},
			},
			pardons: map[int][]bool{
				0: {false, false},
			},
			trust: map[shared.ClientID]float64{
				shared.Team1: 0.5,
				shared.Team2: 0.5,
			},
		},
		{
			name: "get pardon",
			sanctions: map[int][]shared.Sanction{
				0: {
					{
						ClientID:     shared.Team1,
						SanctionTier: shared.SanctionTier1,
						TurnsLeft:    3,
					},
				},
			},
			pardons: map[int][]bool{
				0: {true},
			},
			trust: map[shared.ClientID]float64{
				shared.Team1: 0.8,
				shared.Team2: 0.2,
			},
		},
		{
			name: "get one pardon",
			sanctions: map[int][]shared.Sanction{
				0: {
					{
						ClientID:     shared.Team1,
						SanctionTier: shared.SanctionTier1,
						TurnsLeft:    3,
					},
					{
						ClientID:     shared.Team2,
						SanctionTier: shared.SanctionTier4, // this will prevent getting a pardon
						TurnsLeft:    2,
					},
				},
			},
			pardons: map[int][]bool{
				0: {true, false},
			},
			trust: map[shared.ClientID]float64{
				shared.Team1: 0.7,
				shared.Team2: 0.3,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			clients := []shared.ClientID{}
			for clientID := range tc.trust {
				clients = append(clients, clientID)
			}

			testServer := fakeServerHandle{clients: clients}
			testClient := newClientInternal(shared.Team4, forTesting)
			testClient.Initialise(testServer)

			for client, clientTrust := range tc.trust {
				testClient.trustMatrix.SetClientTrust(client, clientTrust)
			}

			j := testClient.GetClientJudgePointer()

			pardons := j.GetPardonedIslands(tc.sanctions)

			if reflect.DeepEqual(tc.pardons, map[int][]bool{}) {
				if !reflect.DeepEqual(pardons, map[int][]bool{}) {
					t.Errorf("GetPardonedIslands failed. expected: empty map, got %v", pardons)
				}
			}

			if !reflect.DeepEqual(pardons, tc.pardons) {
				t.Errorf("GetPardonedIslands failed. expected: %v, got %v", tc.pardons, pardons)
			}

		})
	}
}
