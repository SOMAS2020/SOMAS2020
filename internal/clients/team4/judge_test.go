package team4

import (
	"math/rand"
	"reflect"
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func makeHistory(clientVars map[shared.ClientID][]rules.VariableFieldName) (history []shared.Accountability, expected map[shared.ClientID]judgeHistoryInfo) {
	history = []shared.Accountability{}
	expected = map[shared.ClientID]judgeHistoryInfo{
		shared.Team1: {},
		shared.Team2: {},
		shared.Team3: {},
		shared.Team4: {},
		shared.Team5: {},
		shared.Team6: {},
	}
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
			expected[client] = got
		}
	}

	return history, expected
}

func TestSaveHistoryInfo(t *testing.T) {
	cases := []struct {
		name            string
		clientVariables map[shared.ClientID][]rules.VariableFieldName
		lieCounts       map[shared.ClientID]int
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
			lieCounts: map[shared.ClientID]int{
				shared.Team1: 6,
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
			lieCounts: map[shared.ClientID]int{
				shared.Team1: 5,
				shared.Team2: 12,
				shared.Team3: 0,
			},
			turn: 1,
			reps: 1,
		},
		{
			name: "Single client empty",
			clientVariables: map[shared.ClientID][]rules.VariableFieldName{
				shared.Team1: {
					rules.IslandReportedPrivateResources,
					rules.IslandActualPrivateResources,
					rules.IslandTaxContribution,
					rules.ExpectedTaxContribution,
					rules.IslandAllocation,
				},
			},
			lieCounts: map[shared.ClientID]int{
				shared.Team1: 6,
			},
			turn: 12,
			reps: 1,
		},
		{
			name: "Multiple clients empty",
			clientVariables: map[shared.ClientID][]rules.VariableFieldName{
				shared.Team1: {
					rules.IslandReportedPrivateResources,
					rules.IslandActualPrivateResources,
					rules.IslandTaxContribution,
					rules.ExpectedTaxContribution,
					rules.ExpectedAllocation,
				},
				shared.Team2: {
					rules.IslandReportedPrivateResources,
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
				},
			},
			lieCounts: map[shared.ClientID]int{
				shared.Team1: 1,
				shared.Team2: 7,
				shared.Team3: 22,
			},
			turn: 112,
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
			lieCounts: map[shared.ClientID]int{
				shared.Team1: 1,
				shared.Team2: 7,
				shared.Team3: 22,
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
			lieCounts: map[shared.ClientID]int{
				shared.Team1: 1,
				shared.Team2: 7,
				shared.Team3: 22,
			},
			turn: 13,
			reps: 3,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			testClient := client{
				BaseClient:    baseclient.NewClient(id),
				clientJudge:   judge{BaseJudge: &baseclient.BaseJudge{}, t: t},
				clientSpeaker: speaker{BaseSpeaker: &baseclient.BaseSpeaker{}},
				yes:           "",
				obs:           &observation{},
				internalParam: &internalParameters{},
				savedHistory:  map[uint]map[shared.ClientID]judgeHistoryInfo{},
			}
			testClient.clientJudge.parent = &testClient

			j := testClient.clientJudge

			wholeHistory := map[uint]map[shared.ClientID]judgeHistoryInfo{}

			for i := uint(0); i < tc.reps; i++ {
				fakeHistory, expected := makeHistory(tc.clientVariables)
				turn := tc.turn + i
				wholeHistory[turn] = expected

				j.saveHistoryInfo(&fakeHistory, &tc.lieCounts, turn)

				if reflect.DeepEqual(expected, testClient.savedHistory[turn]) {
					t.Errorf("Single history failed: %v", testClient.savedHistory[turn])
				}
			}

			if reflect.DeepEqual(wholeHistory, testClient.savedHistory) {
				t.Errorf("Whole history comparison failed. Saved history: %v", testClient.savedHistory)
			}

		})
	}
}
