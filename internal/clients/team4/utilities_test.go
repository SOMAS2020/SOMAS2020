package team4

import (
	"reflect"
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func makeSingleVar(variable rules.VariableFieldName, value float64) rules.VariableValuePair {
	return rules.MakeVariableValuePair(variable, []float64{value})
}

func TestBuildHistoryInfo(t *testing.T) {
	cases := []struct {
		name       string
		valuePairs map[rules.VariableFieldName]float64
		ok         bool
		expected   judgeHistoryInfo
	}{
		{
			name: "Only required ok",
			valuePairs: map[rules.VariableFieldName]float64{
				rules.IslandReportedPrivateResources: 12.2,
				rules.IslandActualPrivateResources:   12.2,
				rules.IslandTaxContribution:          1.2,
				rules.ExpectedTaxContribution:        12.2,
				rules.IslandAllocation:               5463.1,
				rules.ExpectedAllocation:             865.124,
			},
			ok: true,
			expected: judgeHistoryInfo{
				Resources: valuePair{
					actual:   12.2,
					expected: 12.2,
				},
				Taxation: valuePair{
					actual:   1.2,
					expected: 12.2,
				},
				Allocation: valuePair{
					actual:   5463.1,
					expected: 865.124,
				},
				LawfulRatio: 0,
			},
		},
		{
			name: "More than required required ok",
			valuePairs: map[rules.VariableFieldName]float64{
				rules.IslandReportedPrivateResources: 12.2,
				rules.IslandActualPrivateResources:   12.2,
				rules.IslandTaxContribution:          1.2,
				rules.ExpectedTaxContribution:        12.2,
				rules.IslandAllocation:               5463.1,
				rules.ExpectedAllocation:             865.124,
				rules.AllocationMade:                 865.124,
				rules.AnnouncementResultMatchesVote:  865.124,
			},
			ok: true,
			expected: judgeHistoryInfo{
				Resources: valuePair{
					actual:   12.2,
					expected: 12.2,
				},
				Taxation: valuePair{
					actual:   1.2,
					expected: 12.2,
				},
				Allocation: valuePair{
					actual:   5463.1,
					expected: 865.124,
				},
				LawfulRatio: 0,
			},
		},
		{
			name: "not ok",
			valuePairs: map[rules.VariableFieldName]float64{
				rules.IslandReportedPrivateResources: 12.2,
				rules.IslandActualPrivateResources:   12.2,
				rules.IslandTaxContribution:          1.2,
				rules.ExpectedTaxContribution:        12.2,
				rules.IslandAllocation:               5463.1,
			},
			ok:       false,
			expected: judgeHistoryInfo{},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			pairs := []rules.VariableValuePair{}
			for name, val := range tc.valuePairs {
				pairs = append(pairs, makeSingleVar(name, val))
			}
			got, ok := buildHistoryInfo(pairs)

			if ok != tc.ok {
				t.Errorf("Function was expected to return %v, returned %v", tc.ok, ok)
			} else if ok {
				if !reflect.DeepEqual(got, tc.expected) {
					t.Errorf("Function was expected to return %v, returned %v", tc.expected, got)
				}
			}

		})
	}
}

func TestCheckIfClientIsInList(t *testing.T) {
	cases := []struct {
		name     string
		lst      []shared.ClientID
		client   shared.ClientID
		expected bool
	}{
		{
			name:     "Empty list",
			lst:      []shared.ClientID{},
			client:   shared.ClientID(0),
			expected: false,
		},
		{
			name:     "Single Entry true",
			lst:      []shared.ClientID{shared.ClientID(1)},
			client:   shared.ClientID(1),
			expected: true,
		},
		{
			name:     "Single Entry false",
			lst:      []shared.ClientID{shared.ClientID(1)},
			client:   shared.ClientID(2),
			expected: false,
		},
		{
			name:     "Longer list false",
			lst:      []shared.ClientID{shared.ClientID(1), shared.ClientID(2), shared.ClientID(4), shared.ClientID(5)},
			client:   shared.ClientID(3),
			expected: false,
		},
		{
			name:     "Longer list true",
			lst:      []shared.ClientID{shared.ClientID(1), shared.ClientID(2), shared.ClientID(4), shared.ClientID(5)},
			client:   shared.ClientID(5),
			expected: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ret := checkIfClientIsInList(tc.lst, tc.client)
			if ret != tc.expected {
				t.Errorf("%v - Failed. Input (%v, %v). Expected %v, got %v",
					tc.name, tc.lst, tc.client, tc.expected, ret)
			}
		})
	}
}
