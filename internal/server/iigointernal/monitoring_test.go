package iigointernal

import (
	"reflect"
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func TestAddToCache(t *testing.T) {
	cases := []struct {
		name        string
		roleID      shared.ClientID
		variables   []rules.VariableFieldName
		values      [][]float64
		expectedVal []shared.Accountability
	}{
		{
			name:      "Basic adding variables with corresponding values",
			roleID:    shared.ClientID(1),
			variables: []rules.VariableFieldName{rules.RuleSelected, rules.VoteCalled},
			values:    [][]float64{{1}, {1}},
			expectedVal: []shared.Accountability{
				{
					ClientID: shared.ClientID(1),
					Pairs: []rules.VariableValuePair{
						rules.MakeVariableValuePair(rules.RuleSelected, []float64{1}),
						rules.MakeVariableValuePair(rules.VoteCalled, []float64{1}),
					},
				},
			},
		},
		{
			name:        "Adding a variable and too many values",
			roleID:      shared.ClientID(1),
			variables:   []rules.VariableFieldName{rules.RuleSelected},
			values:      [][]float64{{1}, {1}},
			expectedVal: []shared.Accountability{},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			monitor := &monitor{
				internalIIGOCache: []shared.Accountability{},
			}
			monitor.addToCache(tc.roleID, tc.variables, tc.values)
			res := monitor.internalIIGOCache
			if !reflect.DeepEqual(res, tc.expectedVal) {
				t.Errorf("Expected internalIIGOCache to be %v got %v", tc.expectedVal, res)
			}
		})
	}
}

func TestEvaluateCache(t *testing.T) {
	cases := []struct {
		name        string
		roleID      shared.ClientID
		iigoCache   []shared.Accountability
		expectedVal bool
	}{
		{
			name:   "Basic evaluation of compliant President",
			roleID: shared.ClientID(1),
			iigoCache: []shared.Accountability{
				{
					ClientID: shared.ClientID(1),
					Pairs: []rules.VariableValuePair{
						rules.MakeVariableValuePair(rules.RuleSelected, []float64{1}),
						rules.MakeVariableValuePair(rules.VoteCalled, []float64{1}),
					},
				},
			},
			expectedVal: true,
		},
		{
			name:   "Basic evaluation of non compliant Speaker",
			roleID: shared.ClientID(1),
			iigoCache: []shared.Accountability{
				{
					ClientID: shared.ClientID(1),
					Pairs: []rules.VariableValuePair{
						rules.MakeVariableValuePair(rules.RuleSelected, []float64{0}),
						rules.MakeVariableValuePair(rules.VoteCalled, []float64{1}),
					},
				},
			},
			expectedVal: false,
		},
		{
			name:        "Evaluating with empty cache",
			roleID:      shared.ClientID(1),
			iigoCache:   []shared.Accountability{},
			expectedVal: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			monitor := &monitor{
				internalIIGOCache: tc.iigoCache,
			}
			res := monitor.evaluateCache(tc.roleID)
			if !reflect.DeepEqual(res, tc.expectedVal) {
				t.Errorf("Expected evaluation of internalIIGOCache to be %v got %v", tc.expectedVal, res)
			}
		})
	}
}
