package iigointernal

import (
	"reflect"
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"github.com/SOMAS2020/SOMAS2020/pkg/testutils"
	"github.com/pkg/errors"
	"gonum.org/v1/gonum/mat"
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
			gamestate := gamestate.GameState{
				IIGORoleMonitoringCache: []shared.Accountability{},
			}
			monitor := &monitor{
				gameState: &gamestate,
			}
			monitor.addToCache(tc.roleID, tc.variables, tc.values)
			res := gamestate.IIGORoleMonitoringCache
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
	var logging shared.Logger = func(format string, a ...interface{}) {}
	ruleStore := registerMonitoringTestRule()
	tempCache, _ := rules.InitialRuleRegistration(false)
	avail := ruleStore
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			monitoring := &monitor{
				logger: logging,
				gameState: &gamestate.GameState{
					RulesInfo: gamestate.RulesContext{
						VariableMap:    generateDummyVariableCache(),
						AvailableRules: avail,
					},
					IIGORoleMonitoringCache: tc.iigoCache,
					IIGORulesBrokenByRoles:  map[shared.Role][]string{},
				},
			}
			res := monitoring.evaluateCache(tc.roleID, shared.President, ruleStore)
			if !reflect.DeepEqual(res, tc.expectedVal) {
				t.Errorf("Expected evaluation of internalIIGOCache to be %v got %v", tc.expectedVal, res)
			}
		})
	}
	avail = tempCache
}

func TestFindRoleToMonitor(t *testing.T) {
	cases := []struct {
		name            string
		roleAccountable shared.ClientID
		expectedRoleID  shared.ClientID
		expectedRole    shared.Role
		expectedError   error
	}{
		{
			name:            "Test Speaker to perform monitoring",
			roleAccountable: shared.ClientID(1),
			expectedRoleID:  shared.ClientID(2),
			expectedRole:    shared.President,
			expectedError:   nil,
		},
		{
			name:            "Test President to perform monitoring",
			roleAccountable: shared.ClientID(2),
			expectedRoleID:  shared.ClientID(3),
			expectedRole:    shared.Judge,
			expectedError:   nil,
		},
		{
			name:            "Test Judge to perform monitoring",
			roleAccountable: shared.ClientID(3),
			expectedRoleID:  shared.ClientID(1),
			expectedRole:    shared.Speaker,
			expectedError:   nil,
		},
		{
			name:            "Test non IIGO role trying to perform monitoring",
			roleAccountable: shared.ClientID(4),
			expectedRoleID:  shared.ClientID(-1),
			expectedRole:    shared.Speaker,
			expectedError:   errors.Errorf("Monitoring by island that is not an IIGO Role"),
		},
	}
	monitoring := &monitor{
		gameState: &gamestate.GameState{
			SpeakerID:   1,
			PresidentID: 2,
			JudgeID:     3,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			id, role, err := monitoring.findRoleToMonitor(tc.roleAccountable)
			if !(reflect.DeepEqual(tc.expectedRoleID, id) && reflect.DeepEqual(tc.expectedRole, role)) {
				t.Errorf("Expected role to monitor to be %v got %v", tc.expectedRoleID, id)
				t.Errorf("Expected role to monitor to be %v got %v", tc.expectedRole, role)
			}
			testutils.CompareTestErrors(tc.expectedError, err, t)
		})
	}
}

func registerMonitoringTestRule() map[string]rules.RuleMatrix {

	rulesStore := map[string]rules.RuleMatrix{}

	name := "vote_called_rule"
	reqVar := []rules.VariableFieldName{
		rules.RuleSelected,
		rules.VoteCalled,
	}

	v := []float64{1, -1, 0}
	CoreMatrix := mat.NewDense(1, 3, v)
	aux := []float64{0}
	AuxiliaryVector := mat.NewVecDense(1, aux)

	rm := rules.RuleMatrix{RuleName: name, RequiredVariables: reqVar, ApplicableMatrix: *CoreMatrix, AuxiliaryVector: *AuxiliaryVector, Mutable: false}
	rulesStore[name] = rm
	return rulesStore
}
