package iigointernal

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"gonum.org/v1/gonum/mat"
)

// TestSearchEvalReturnForRuleName checks whether the searchEvalReturnForRuleName is able to pick up which rules
// have been evaluated in a set of EvaluationReturns
func TestSearchEvalReturnForRuleName(t *testing.T) {
	availableRules := generateDummyRuleMatrices()
	cases := []struct {
		name     string
		ruleName string
		evalRet  roles.EvaluationReturn
		expected bool
	}{
		{
			name:     "Simple search",
			ruleName: "iigo_economic_sanction_5",
			evalRet: roles.EvaluationReturn{
				Rules: []rules.RuleMatrix{
					availableRules[0],
					availableRules[8],
				},
				Evaluations: []bool{true, true},
			},
			expected: true,
		},
		{
			name:     "Expect search to return false",
			ruleName: "iigo_economic_sanction_5",
			evalRet: roles.EvaluationReturn{
				Rules: []rules.RuleMatrix{
					availableRules[0],
					availableRules[1],
				},
				Evaluations: []bool{true, true},
			},
			expected: false,
		},
		{
			name:     "More complex search",
			ruleName: "iigo_economic_sanction_5",
			evalRet: roles.EvaluationReturn{
				Rules: []rules.RuleMatrix{
					availableRules[0],
					availableRules[1],
					availableRules[3],
					availableRules[4],
					availableRules[8],
				},
				Evaluations: []bool{true, true, false, true, false},
			},
			expected: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res := searchEvalReturnForRuleName(tc.ruleName, tc.evalRet)
			if res != tc.expected {
				t.Errorf("Expected %v got %v", tc.expected, res)
			}
		})
	}
}

// TestCullCheckedRules checks whether cullCheckedRules is able to remove elements of the history of a turn that
// have already been checked by a client Judge
func TestCullCheckedRules(t *testing.T) {
	availableRules := generateDummyRuleMatrices()
	cases := []struct {
		name     string
		history  []shared.Accountability
		evalRes  map[shared.ClientID]roles.EvaluationReturn
		expected []shared.Accountability
	}{
		{
			name: "Basic cull test",
			history: []shared.Accountability{
				{
					ClientID: shared.Team1,
					Pairs: []rules.VariableValuePair{
						{
							VariableName: rules.SanctionPaid,
							Values:       []float64{5},
						},
						{
							VariableName: rules.SanctionExpected,
							Values:       []float64{5},
						},
					},
				},
			},
			evalRes: map[shared.ClientID]roles.EvaluationReturn{
				shared.Team1: {
					Rules: []rules.RuleMatrix{
						availableRules[9],
					},
					Evaluations: []bool{true},
				},
			},
			expected: []shared.Accountability{},
		},
		{
			name: "More Advanced Cull",
			history: []shared.Accountability{
				{
					ClientID: shared.Team1,
					Pairs: []rules.VariableValuePair{
						{
							VariableName: rules.SanctionPaid,
							Values:       []float64{5},
						},
						{
							VariableName: rules.SanctionExpected,
							Values:       []float64{5},
						},
					},
				},
				{
					ClientID: shared.Team2,
					Pairs: []rules.VariableValuePair{
						{
							VariableName: rules.SanctionPaid,
							Values:       []float64{5},
						},
						{
							VariableName: rules.SanctionExpected,
							Values:       []float64{5},
						},
					},
				},
			},
			evalRes: map[shared.ClientID]roles.EvaluationReturn{
				shared.Team1: {
					Rules: []rules.RuleMatrix{
						availableRules[9],
					},
					Evaluations: []bool{true},
				},
			},
			expected: []shared.Accountability{
				{
					ClientID: shared.Team2,
					Pairs: []rules.VariableValuePair{
						{
							VariableName: rules.SanctionPaid,
							Values:       []float64{5},
						},
						{
							VariableName: rules.SanctionExpected,
							Values:       []float64{5},
						},
					},
				},
			},
		},
		{
			name: "Even More Advanced Cull",
			history: []shared.Accountability{
				{
					ClientID: shared.Team1,
					Pairs: []rules.VariableValuePair{
						{
							VariableName: rules.SanctionPaid,
							Values:       []float64{5},
						},
						{
							VariableName: rules.SanctionExpected,
							Values:       []float64{5},
						},
					},
				},
				{
					ClientID: shared.Team2,
					Pairs: []rules.VariableValuePair{
						{
							VariableName: rules.SanctionPaid,
							Values:       []float64{5},
						},
						{
							VariableName: rules.SanctionExpected,
							Values:       []float64{5},
						},
					},
				},
			},
			evalRes: map[shared.ClientID]roles.EvaluationReturn{
				shared.Team1: {
					Rules: []rules.RuleMatrix{
						availableRules[9],
					},
					Evaluations: []bool{true},
				},
				shared.Team2: {
					Rules: []rules.RuleMatrix{
						availableRules[8],
					},
					Evaluations: []bool{true},
				},
			},
			expected: []shared.Accountability{
				{
					ClientID: shared.Team2,
					Pairs: []rules.VariableValuePair{
						{
							VariableName: rules.SanctionPaid,
							Values:       []float64{5},
						},
						{
							VariableName: rules.SanctionExpected,
							Values:       []float64{5},
						},
					},
				},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res := cullCheckedRules(tc.history, tc.evalRes, generateRuleStore(), generateDummyVariableCache())
			if !reflect.DeepEqual(res, tc.expected) {
				t.Errorf("Expected %v got %v", tc.expected, res)
			}
		})
	}
}

// TestPickUpRulesByVariable checks whether the pickUpRulesByVariable is able to correctly identify what rules
// a particular variable has affected
func TestPickUpRulesByVariable(t *testing.T) {
	ruleMap := generateRuleStore()
	cases := []struct {
		name        string
		searchVar   rules.VariableFieldName
		ruleStore   map[string]rules.RuleMatrix
		expectedRes bool
		expectedArr []string
	}{
		{
			name:        "Basic Pickup Test",
			searchVar:   rules.SanctionExpected,
			ruleStore:   ruleMap,
			expectedRes: true,
			expectedArr: []string{"check_sanction_rule"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			arr, res := pickUpRulesByVariable(tc.searchVar, tc.ruleStore, generateDummyVariableCache())
			if !reflect.DeepEqual(arr, tc.expectedArr) {
				t.Errorf("Expected %v got %v", tc.expectedArr, arr)
			} else if res != tc.expectedRes {
				t.Errorf("Expected %v got %v", tc.expectedRes, res)
			}
		})
	}
}

// TestCheckPardons checks whether checkPardons is able to correctly identify pardons issued on sanctions that
// don't exist, and then subtracts these pardoned sanctions from the sanctionCache
func TestImplementPardons(t *testing.T) {
	TeamIDs := []shared.ClientID{shared.Team1, shared.Team2, shared.Team3, shared.Team4, shared.Team5, shared.Team6}
	cases := []struct {
		name          string
		sanctionCache map[int][]roles.Sanction
		pardons       map[int][]bool
		expValidity   bool
		expComms      map[shared.ClientID][]map[shared.CommunicationFieldName]shared.CommunicationContent
		expFinCache   map[int][]roles.Sanction
	}{
		{
			name: "No Pardons Check",
			sanctionCache: map[int][]roles.Sanction{
				1: {
					{
						ClientID:     shared.Team1,
						SanctionTier: roles.SanctionTier3,
						TurnsLeft:    3,
					},
				},
			},
			pardons: map[int][]bool{
				1: {
					false,
				},
			},
			expValidity: true,
			expComms:    generateEmptyCommunicationsMap(TeamIDs),
			expFinCache: map[int][]roles.Sanction{
				1: {
					{
						ClientID:     shared.Team1,
						SanctionTier: roles.SanctionTier3,
						TurnsLeft:    3,
					},
				},
			},
		},
		{
			name: "Simple Pardons Case",
			sanctionCache: map[int][]roles.Sanction{
				1: {
					{
						ClientID:     shared.Team1,
						SanctionTier: roles.SanctionTier3,
						TurnsLeft:    3,
					},
					{
						ClientID:     shared.Team2,
						SanctionTier: roles.SanctionTier3,
						TurnsLeft:    3,
					},
				},
			},
			pardons: map[int][]bool{
				1: {
					true,
					false,
				},
			},
			expValidity: true,
			expComms: map[shared.ClientID][]map[shared.CommunicationFieldName]shared.CommunicationContent{
				shared.Team1: {
					{
						shared.PardonClientID: {
							T:           shared.CommunicationInt,
							IntegerData: int(shared.Team1),
						},
						shared.PardonTier: {
							T:           shared.CommunicationInt,
							IntegerData: int(roles.SanctionTier3),
						},
					},
				},
				shared.Team2: {},
				shared.Team3: {},
				shared.Team4: {},
				shared.Team5: {},
				shared.Team6: {},
			},
			expFinCache: map[int][]roles.Sanction{
				1: {
					{
						ClientID:     shared.Team2,
						SanctionTier: roles.SanctionTier3,
						TurnsLeft:    3,
					},
				},
			},
		},
		{
			name: "Complex Pardons Case",
			sanctionCache: map[int][]roles.Sanction{
				1: {
					{
						ClientID:     shared.Team1,
						SanctionTier: roles.SanctionTier3,
						TurnsLeft:    3,
					},
					{
						ClientID:     shared.Team2,
						SanctionTier: roles.SanctionTier3,
						TurnsLeft:    3,
					},
					{
						ClientID:     shared.Team3,
						SanctionTier: roles.SanctionTier2,
						TurnsLeft:    2,
					},
				},
			},
			pardons: map[int][]bool{
				1: {
					true,
					false,
					true,
				},
			},
			expValidity: true,
			expComms: map[shared.ClientID][]map[shared.CommunicationFieldName]shared.CommunicationContent{
				shared.Team1: {
					{
						shared.PardonClientID: {
							T:           shared.CommunicationInt,
							IntegerData: int(shared.Team1),
						},
						shared.PardonTier: {
							T:           shared.CommunicationInt,
							IntegerData: int(roles.SanctionTier3),
						},
					},
				},
				shared.Team3: {
					{
						shared.PardonClientID: {
							T:           shared.CommunicationInt,
							IntegerData: int(shared.Team3),
						},
						shared.PardonTier: {
							T:           shared.CommunicationInt,
							IntegerData: int(roles.SanctionTier2),
						},
					},
				},
				shared.Team2: {},
				shared.Team4: {},
				shared.Team5: {},
				shared.Team6: {},
			},
			expFinCache: map[int][]roles.Sanction{
				1: {
					{
						ClientID:     shared.Team2,
						SanctionTier: roles.SanctionTier3,
						TurnsLeft:    3,
					},
				},
			},
		},
		{
			name: "Complex Pardons Failure Case",
			sanctionCache: map[int][]roles.Sanction{
				1: {
					{
						ClientID:     shared.Team1,
						SanctionTier: roles.SanctionTier3,
						TurnsLeft:    3,
					},
					{
						ClientID:     shared.Team2,
						SanctionTier: roles.SanctionTier3,
						TurnsLeft:    3,
					},
					{
						ClientID:     shared.Team3,
						SanctionTier: roles.SanctionTier2,
						TurnsLeft:    2,
					},
				},
			},
			pardons: map[int][]bool{
				1: {
					false,
				},
			},
			expValidity: false,
			expComms:    nil,
			expFinCache: map[int][]roles.Sanction{
				1: {
					{
						ClientID:     shared.Team1,
						SanctionTier: roles.SanctionTier3,
						TurnsLeft:    3,
					},
					{
						ClientID:     shared.Team2,
						SanctionTier: roles.SanctionTier3,
						TurnsLeft:    3,
					},
					{
						ClientID:     shared.Team3,
						SanctionTier: roles.SanctionTier2,
						TurnsLeft:    2,
					},
				},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			val, finalCache, comms := implementPardons(tc.sanctionCache, tc.pardons, TeamIDs)
			if val != tc.expValidity {
				t.Errorf("Expected validity %v got %v", tc.expValidity, val)
			} else if !reflect.DeepEqual(comms, tc.expComms) {
				t.Errorf("Expected Communications to be %v got %v", tc.expComms, comms)
			} else if !reflect.DeepEqual(finalCache, tc.expFinCache) {
				t.Errorf("Expected Final Cache to be %v got %v", tc.expFinCache, finalCache)
			}
		})
	}
}

// TestSoftMergeSanctionThreshold checks whether the softMergeSanctionThresholds is able to correctly merge
// the clients threshols with the default ones
func TestSoftMergeSanctionThreshold(t *testing.T) {
	cases := []struct {
		name              string
		clientSanctionMap map[roles.IIGOSanctionTier]roles.IIGOSanctionScore
		expectedVal       map[roles.IIGOSanctionTier]roles.IIGOSanctionScore
	}{
		{
			name: "Basic soft merge",
			clientSanctionMap: map[roles.IIGOSanctionTier]roles.IIGOSanctionScore{
				roles.SanctionTier1: roles.IIGOSanctionScore(3),
			},
			expectedVal: map[roles.IIGOSanctionTier]roles.IIGOSanctionScore{
				roles.SanctionTier1: 3,
				roles.SanctionTier2: 5,
				roles.SanctionTier3: 10,
				roles.SanctionTier4: 20,
				roles.SanctionTier5: 30,
			},
		},
		{
			name:              "No merge",
			clientSanctionMap: map[roles.IIGOSanctionTier]roles.IIGOSanctionScore{},
			expectedVal: map[roles.IIGOSanctionTier]roles.IIGOSanctionScore{
				roles.SanctionTier1: 1,
				roles.SanctionTier2: 5,
				roles.SanctionTier3: 10,
				roles.SanctionTier4: 20,
				roles.SanctionTier5: 30,
			},
		},
		{
			name: "More complicated merge",
			clientSanctionMap: map[roles.IIGOSanctionTier]roles.IIGOSanctionScore{
				roles.SanctionTier1: 7,
				roles.SanctionTier2: 9,
				roles.SanctionTier5: 400,
			},
			expectedVal: map[roles.IIGOSanctionTier]roles.IIGOSanctionScore{
				roles.SanctionTier1: 7,
				roles.SanctionTier2: 9,
				roles.SanctionTier3: 10,
				roles.SanctionTier4: 20,
				roles.SanctionTier5: 400,
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res := softMergeSanctionThresholds(tc.clientSanctionMap)
			if !reflect.DeepEqual(res, tc.expectedVal) {
				t.Errorf("Expected final transgressions to be %v got %v", tc.expectedVal, res)
			}
		})
	}
}

// Utility functions for testing //

func generateRuleStore() map[string]rules.RuleMatrix {
	returnMap := map[string]rules.RuleMatrix{}
	availableRules := generateDummyRuleMatrices()
	for _, v := range availableRules {
		returnMap[v.RuleName] = v

	}
	return returnMap
}

func generateDummyVariableCache() map[rules.VariableFieldName]rules.VariableValuePair {
	returnMap := map[rules.VariableFieldName]rules.VariableValuePair{}
	availableVariables := getStaticVariables()
	for _, v := range availableVariables {
		returnMap[v.VariableName] = v
	}
	return returnMap
}

func getStaticVariables() []rules.VariableValuePair {
	return []rules.VariableValuePair{
		{
			VariableName: rules.NumberOfIslandsContributingToCommonPool,
			Values:       []float64{5},
		},
		{
			VariableName: rules.NumberOfFailedForages,
			Values:       []float64{0.5},
		},
		{
			VariableName: rules.NumberOfBrokenAgreements,
			Values:       []float64{1},
		},
		{
			VariableName: rules.MaxSeverityOfSanctions,
			Values:       []float64{2},
		},
		{
			VariableName: rules.NumberOfIslandsAlive,
			Values:       []float64{6},
		},
		{
			VariableName: rules.NumberOfBallotsCast,
			Values:       []float64{6},
		},
		{
			VariableName: rules.NumberOfAllocationsSent,
			Values:       []float64{6},
		},
		{
			VariableName: rules.IslandsAlive,
			Values:       []float64{0, 1, 2, 3, 4, 5},
		},
		{
			VariableName: rules.SpeakerSalary,
			Values:       []float64{50},
		},
		{
			VariableName: rules.JudgeSalary,
			Values:       []float64{50},
		},
		{
			VariableName: rules.PresidentSalary,
			Values:       []float64{50},
		},
		{
			VariableName: rules.ExpectedTaxContribution,
			Values:       []float64{0},
		},
		{
			VariableName: rules.ExpectedAllocation,
			Values:       []float64{0},
		},
		{
			VariableName: rules.IslandTaxContribution,
			Values:       []float64{0},
		},
		{
			VariableName: rules.IslandAllocation,
			Values:       []float64{0},
		},
		{
			VariableName: rules.IslandReportedResources,
			Values:       []float64{0},
		},
		{
			VariableName: rules.ConstSanctionAmount,
			Values:       []float64{0},
		},
		{
			VariableName: rules.TurnsLeftOnSanction,
			Values:       []float64{0},
		},
		{
			VariableName: rules.SanctionPaid,
			Values:       []float64{0},
		},
		{
			VariableName: rules.SanctionExpected,
			Values:       []float64{0},
		},
	}
}

func generateDummyRuleMatrices() []rules.RuleMatrix {
	ruleSpecs := []struct {
		name    string
		reqVar  []rules.VariableFieldName
		v       []float64
		aux     []float64
		mutable bool
	}{
		{
			name: "inspect_ballot_rule",
			reqVar: []rules.VariableFieldName{
				rules.NumberOfIslandsAlive,
				rules.NumberOfBallotsCast,
			},
			v:       []float64{1, -1, 0},
			aux:     []float64{0},
			mutable: false,
		},
		{
			name: "inspect_allocation_rule",
			reqVar: []rules.VariableFieldName{
				rules.NumberOfIslandsAlive,
				rules.NumberOfAllocationsSent,
			},
			v:       []float64{1, -1, 0},
			aux:     []float64{0},
			mutable: false,
		},
		{
			name: "check_taxation_rule",
			reqVar: []rules.VariableFieldName{
				rules.IslandTaxContribution,
				rules.ExpectedTaxContribution,
			},
			v:       []float64{1, -1, 0},
			aux:     []float64{2},
			mutable: false,
		},
		{
			name: "check_allocation_rule",
			reqVar: []rules.VariableFieldName{
				rules.IslandAllocation,
				rules.ExpectedAllocation,
			},
			v:       []float64{1, -1, 0},
			aux:     []float64{0},
			mutable: false,
		},
		{
			name: "iigo_economic_sanction_1",
			reqVar: []rules.VariableFieldName{
				rules.IslandReportedResources,
				rules.ConstSanctionAmount,
				rules.TurnsLeftOnSanction,
			},
			v:       []float64{0, 0, 1, 0, 0, 0, 0, 0},
			aux:     []float64{1, 4},
			mutable: true,
		},
		{
			name: "iigo_economic_sanction_2",
			reqVar: []rules.VariableFieldName{
				rules.IslandReportedResources,
				rules.ConstSanctionAmount,
				rules.TurnsLeftOnSanction,
			},
			v:       []float64{0, 0, 1, 0, 0.1, 1, 0, 0},
			aux:     []float64{1, 4},
			mutable: true,
		},
		{
			name: "iigo_economic_sanction_3",
			reqVar: []rules.VariableFieldName{
				rules.IslandReportedResources,
				rules.ConstSanctionAmount,
				rules.TurnsLeftOnSanction,
			},
			v:       []float64{0, 0, 1, 0, 0.3, 1, 0, 0},
			aux:     []float64{1, 4},
			mutable: true,
		},
		{
			name: "iigo_economic_sanction_4",
			reqVar: []rules.VariableFieldName{
				rules.IslandReportedResources,
				rules.ConstSanctionAmount,
				rules.TurnsLeftOnSanction,
			},
			v:       []float64{0, 0, 1, 0, 0.5, 1, 0, 0},
			aux:     []float64{1, 4},
			mutable: true,
		},
		{
			name: "iigo_economic_sanction_5",
			reqVar: []rules.VariableFieldName{
				rules.IslandReportedResources,
				rules.ConstSanctionAmount,
				rules.TurnsLeftOnSanction,
			},
			v:       []float64{0, 0, 1, 0, 0.8, 1, 0, 0},
			aux:     []float64{1, 4},
			mutable: true,
		},
		{
			name: "check_sanction_rule",
			reqVar: []rules.VariableFieldName{
				rules.SanctionPaid,
				rules.SanctionExpected,
			},
			v:       []float64{1, -1, 0},
			aux:     []float64{0},
			mutable: true,
		},
	}
	var outputArray []rules.RuleMatrix
	for _, rs := range ruleSpecs {
		rowLength := len(rs.reqVar) + 1
		if len(rs.v)%rowLength != 0 {
			panic(fmt.Sprintf("Rule '%v' was registered without correct matrix dimensions", rs.name))
		}
		nrows := len(rs.v) / rowLength
		CoreMatrix := mat.NewDense(nrows, rowLength, rs.v)
		AuxiliaryVector := mat.NewVecDense(nrows, rs.aux)
		ruleMat := rules.RuleMatrix{
			RuleName:          rs.name,
			RequiredVariables: rs.reqVar,
			ApplicableMatrix:  *CoreMatrix,
			AuxiliaryVector:   *AuxiliaryVector,
			Mutable:           rs.mutable,
		}
		outputArray = append(outputArray, ruleMat)
	}
	return outputArray
}
