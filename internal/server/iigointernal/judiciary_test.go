package iigointernal

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/roles"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"gonum.org/v1/gonum/mat"
)

/// Judiciary object tests ///

// TestLoadSanctionConfig checks if the judiciary can load sanction information from its client judge
func TestLoadSanctionConfig(t *testing.T) {
	cases := []struct {
		name                          string
		clientJudge                   roles.Judge
		expectedSanctionThresholds    map[shared.IIGOSanctionsTier]shared.IIGOSanctionsScore
		expectedRuleViolationSeverity map[string]shared.IIGOSanctionsScore
	}{
		{
			name: "Basic sanction config load",
			clientJudge: &mockJudge{
				violationSeverity: map[string]shared.IIGOSanctionsScore{
					"Mock Rule": 50,
				},
				sanctionThresholds: map[shared.IIGOSanctionsTier]shared.IIGOSanctionsScore{
					shared.SanctionTier2: 6,
				},
			},
			expectedSanctionThresholds: map[shared.IIGOSanctionsTier]shared.IIGOSanctionsScore{
				shared.SanctionTier1: 1,
				shared.SanctionTier2: 6,
				shared.SanctionTier3: 10,
				shared.SanctionTier4: 20,
				shared.SanctionTier5: 30,
			},
			expectedRuleViolationSeverity: map[string]shared.IIGOSanctionsScore{
				"Mock Rule": 50,
			},
		},
		{
			name: "Checking for monotonicity",
			clientJudge: &mockJudge{
				violationSeverity: map[string]shared.IIGOSanctionsScore{
					"Mock Rule": 50,
				},
				sanctionThresholds: map[shared.IIGOSanctionsTier]shared.IIGOSanctionsScore{
					shared.SanctionTier2: 60,
				},
			},
			expectedSanctionThresholds: map[shared.IIGOSanctionsTier]shared.IIGOSanctionsScore{
				shared.SanctionTier1: 1,
				shared.SanctionTier2: 5,
				shared.SanctionTier3: 10,
				shared.SanctionTier4: 20,
				shared.SanctionTier5: 30,
			},
			expectedRuleViolationSeverity: map[string]shared.IIGOSanctionsScore{
				"Mock Rule": 50,
			},
		},
		{
			name: "Ensuring all rules vals are picked up",
			clientJudge: &mockJudge{
				violationSeverity: map[string]shared.IIGOSanctionsScore{
					"Mock Rule":   50,
					"Mock Rule 2": 100,
				},
				sanctionThresholds: map[shared.IIGOSanctionsTier]shared.IIGOSanctionsScore{
					shared.SanctionTier2: 7,
				},
			},
			expectedSanctionThresholds: map[shared.IIGOSanctionsTier]shared.IIGOSanctionsScore{
				shared.SanctionTier1: 1,
				shared.SanctionTier2: 7,
				shared.SanctionTier3: 10,
				shared.SanctionTier4: 20,
				shared.SanctionTier5: 30,
			},
			expectedRuleViolationSeverity: map[string]shared.IIGOSanctionsScore{
				"Mock Rule":   50,
				"Mock Rule 2": 100,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			judiciaryInst := judiciary{
				clientJudge: tc.clientJudge,
				gameState: &gamestate.GameState{
					RulesInfo: gamestate.RulesContext{
						VariableMap: generateDummyVariableCache(),
					},
				},
			}
			judiciaryInst.loadSanctionConfig()

			if !reflect.DeepEqual(judiciaryInst.sanctionThresholds, tc.expectedSanctionThresholds) {
				t.Errorf("Expected %v got %v", tc.expectedSanctionThresholds, judiciaryInst.sanctionThresholds)
			}
			if !reflect.DeepEqual(judiciaryInst.ruleViolationSeverity, tc.expectedRuleViolationSeverity) {
				t.Errorf("Expected %v got %v", tc.expectedRuleViolationSeverity, judiciaryInst.ruleViolationSeverity)
			}
		})
	}
}

// TODO - Outdated Test? Presidents salary should not be deducted from the budget
// TestSendPresidentSalary checks whether judiciary can correctly send salaries to the executive branch
/*func TestSendPresidentSalary(t *testing.T) {
	cases := []struct {
		name              string
		defaultPresSalary shared.Resources
		clientJudge       roles.Judge
		expectedBudget    shared.Resources
	}{
		{
			name:              "No clientJudge test",
			defaultPresSalary: 50,
			clientJudge:       nil,
			expectedBudget:    50,
		},
		{
			name:              "Fair client judge",
			defaultPresSalary: 50,
			clientJudge: &mockJudge{
				payPresidentVal:    50,
				payPresidentChoice: true,
			},
			expectedBudget: 50,
		},
		{
			name:              "Unfair client Judge",
			defaultPresSalary: 50,
			clientJudge: &mockJudge{
				payPresidentVal:    20,
				payPresidentChoice: true,
			},
			expectedBudget: 20,
		},
		{
			name:              "Client judge not paying",
			defaultPresSalary: 50,
			clientJudge: &mockJudge{
				payPresidentVal:    20,
				payPresidentChoice: false,
			},
			expectedBudget: 0,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dummyExec := executive{}
			judicialInst := judiciary{
				presidentSalary: tc.defaultPresSalary,
				clientJudge:     tc.clientJudge,
			}
			judicialInst.sendPresidentSalary()
			if !reflect.DeepEqual(tc.expectedBudget, dummyExec.budget) {
				t.Errorf("Expected %v got %v", tc.expectedBudget, dummyExec.budget)
			}
		})
	}
}
*/

// TestInspectHistory checks whether inspect history is able to take account of client decisions
func TestInspectHistory(t *testing.T) {
	cases := []struct {
		name            string
		clientJudge     roles.Judge
		iigoHistory     []shared.Accountability
		historicalCache map[int][]shared.Accountability
		expectedResults map[shared.ClientID]shared.EvaluationReturn
		expectedSuccess bool
	}{
		{
			name:            "Empty IIGO history",
			iigoHistory:     []shared.Accountability{},
			historicalCache: map[int][]shared.Accountability{},
			clientJudge: &mockJudge{
				inspectHistoryReturn:  map[shared.ClientID]shared.EvaluationReturn{},
				inspectHistoryChoice:  true,
				historicalRetribution: false,
			},
			expectedResults: getBaseEvalResults(shared.TeamIDs),
			expectedSuccess: true,
		},
		{
			name:            "Non empty IIGO history",
			iigoHistory:     []shared.Accountability{},
			historicalCache: map[int][]shared.Accountability{},
			clientJudge: &mockJudge{
				inspectHistoryReturn:  map[shared.ClientID]shared.EvaluationReturn{},
				inspectHistoryChoice:  false,
				historicalRetribution: false,
			},
			expectedResults: getBaseEvalResults(shared.TeamIDs),
			expectedSuccess: false,
		},
		{
			name:            "Simple Evaluations",
			historicalCache: map[int][]shared.Accountability{},
			iigoHistory: []shared.Accountability{
				{
					ClientID: shared.Team1,
					Pairs:    []rules.VariableValuePair{rules.MakeVariableValuePair(rules.ExpectedAllocation, []float64{5.0})},
				},
			},
			clientJudge: &mockJudge{
				inspectHistoryReturn: map[shared.ClientID]shared.EvaluationReturn{
					shared.Team1: {
						Rules:       []rules.RuleMatrix{generateDummyRuleMatrices()[0]},
						Evaluations: []bool{true},
					},
				},
				inspectHistoryChoice:  true,
				historicalRetribution: false,
			},
			expectedResults: mergeEvaluationReturn(map[shared.ClientID]shared.EvaluationReturn{
				shared.Team1: {
					Rules:       []rules.RuleMatrix{generateDummyRuleMatrices()[0]},
					Evaluations: []bool{true},
				},
			}, getBaseEvalResults(shared.TeamIDs)),
			expectedSuccess: true,
		},
		{
			name:            "Complex Evaluations",
			historicalCache: map[int][]shared.Accountability{},
			iigoHistory: []shared.Accountability{
				{
					ClientID: shared.Team1,
					Pairs:    []rules.VariableValuePair{rules.MakeVariableValuePair(rules.ExpectedAllocation, []float64{5.0})},
				},
				{
					ClientID: shared.Team2,
					Pairs:    []rules.VariableValuePair{rules.MakeVariableValuePair(rules.ExpectedAllocation, []float64{5.0})},
				},
			},
			clientJudge: &mockJudge{
				inspectHistoryReturn: map[shared.ClientID]shared.EvaluationReturn{
					shared.Team1: {
						Rules:       []rules.RuleMatrix{generateDummyRuleMatrices()[0]},
						Evaluations: []bool{true},
					},
					shared.Team2: {
						Rules:       []rules.RuleMatrix{generateDummyRuleMatrices()[1]},
						Evaluations: []bool{false},
					},
				},
				inspectHistoryChoice:  true,
				historicalRetribution: false,
			},
			expectedResults: mergeEvaluationReturn(map[shared.ClientID]shared.EvaluationReturn{
				shared.Team1: {
					Rules:       []rules.RuleMatrix{generateDummyRuleMatrices()[0]},
					Evaluations: []bool{true},
				},
				shared.Team2: {
					Rules:       []rules.RuleMatrix{generateDummyRuleMatrices()[1]},
					Evaluations: []bool{false},
				},
			}, getBaseEvalResults(shared.TeamIDs)),
			expectedSuccess: true,
		},
		{
			name: "Historical Evaluations Tested",
			historicalCache: map[int][]shared.Accountability{
				1: {
					{
						ClientID: shared.Team1,
						Pairs:    []rules.VariableValuePair{rules.MakeVariableValuePair(rules.ExpectedAllocation, []float64{5.0})},
					},
					{
						ClientID: shared.Team2,
						Pairs:    []rules.VariableValuePair{rules.MakeVariableValuePair(rules.ExpectedAllocation, []float64{5.0})},
					},
				},
			},
			iigoHistory: []shared.Accountability{
				{
					ClientID: shared.Team1,
					Pairs:    []rules.VariableValuePair{rules.MakeVariableValuePair(rules.ExpectedAllocation, []float64{5.0})},
				},
				{
					ClientID: shared.Team2,
					Pairs:    []rules.VariableValuePair{rules.MakeVariableValuePair(rules.ExpectedAllocation, []float64{5.0})},
				},
			},
			clientJudge: &mockJudge{
				inspectHistoryReturn: map[shared.ClientID]shared.EvaluationReturn{
					shared.Team1: {
						Rules:       []rules.RuleMatrix{generateDummyRuleMatrices()[0]},
						Evaluations: []bool{true},
					},
				},
				inspectHistoryChoice:  true,
				historicalRetribution: true,
			},
			expectedResults: mergeEvaluationReturn(map[shared.ClientID]shared.EvaluationReturn{
				shared.Team1: {
					Rules:       []rules.RuleMatrix{generateDummyRuleMatrices()[0], generateDummyRuleMatrices()[0]},
					Evaluations: []bool{true, true},
				},
			}, getBaseEvalResults(shared.TeamIDs)),
			expectedSuccess: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			judiciaryInst := defaultInitJudiciary()
			judiciaryInst.clientJudge = tc.clientJudge
			judiciaryInst.gameState.IIGOHistoryCache = tc.historicalCache
			result, success := judiciaryInst.inspectHistory(tc.iigoHistory)
			if success != tc.expectedSuccess {
				t.Errorf("Expected %v got %v", tc.expectedSuccess, success)
			}
			if !reflect.DeepEqual(tc.expectedResults, result) {
				t.Errorf("Expected %v got %v", tc.expectedResults, result)
			}
		})
	}
}

// TestUpdateSanctionScore checks whether the judiciary branch can correctly score an update sanction records for
// agents
func TestUpdateSanctionScore(t *testing.T) {
	cases := []struct {
		name                   string
		evaluationResults      map[shared.ClientID]shared.EvaluationReturn
		ruleViolationPenalties map[string]shared.IIGOSanctionsScore
		expectedIslandScores   map[shared.ClientID]shared.IIGOSanctionsScore
	}{
		{
			name: "Basic update sanction score check",
			evaluationResults: map[shared.ClientID]shared.EvaluationReturn{
				shared.Team1: {
					Rules:       []rules.RuleMatrix{generateDummyRuleMatrices()[0]},
					Evaluations: []bool{false},
				},
			},
			ruleViolationPenalties: map[string]shared.IIGOSanctionsScore{
				"inspect_ballot_rule": 50,
			},
			expectedIslandScores: map[shared.ClientID]shared.IIGOSanctionsScore{
				shared.Team1: 50,
			},
		},
		{
			name: "Normal update sanction score check",
			evaluationResults: map[shared.ClientID]shared.EvaluationReturn{
				shared.Team1: {
					Rules:       []rules.RuleMatrix{generateDummyRuleMatrices()[0]},
					Evaluations: []bool{false},
				},
				shared.Team2: {
					Rules:       []rules.RuleMatrix{generateDummyRuleMatrices()[0], generateDummyRuleMatrices()[1]},
					Evaluations: []bool{false, true},
				},
			},
			ruleViolationPenalties: map[string]shared.IIGOSanctionsScore{
				"inspect_ballot_rule": 50,
			},
			expectedIslandScores: map[shared.ClientID]shared.IIGOSanctionsScore{
				shared.Team1: 50,
				shared.Team2: 50,
			},
		},
		{
			name: "Complex update sanction score scenario",
			evaluationResults: map[shared.ClientID]shared.EvaluationReturn{
				shared.Team1: {
					Rules:       []rules.RuleMatrix{generateDummyRuleMatrices()[0], generateDummyRuleMatrices()[0]},
					Evaluations: []bool{false, false},
				},
				shared.Team2: {
					Rules:       []rules.RuleMatrix{generateDummyRuleMatrices()[0], generateDummyRuleMatrices()[1]},
					Evaluations: []bool{false, false},
				},
				shared.Team3: {
					Rules:       []rules.RuleMatrix{generateDummyRuleMatrices()[4], generateDummyRuleMatrices()[5]},
					Evaluations: []bool{false, true},
				},
			},
			ruleViolationPenalties: map[string]shared.IIGOSanctionsScore{
				"inspect_ballot_rule":      50,
				"iigo_economic_sanction_1": 100,
			},
			expectedIslandScores: map[shared.ClientID]shared.IIGOSanctionsScore{
				shared.Team1: 100,
				shared.Team2: 55,
				shared.Team3: 100,
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			judiciaryInst := defaultInitJudiciary()
			judiciaryInst.evaluationResults = tc.evaluationResults
			judiciaryInst.ruleViolationSeverity = tc.ruleViolationPenalties
			judiciaryInst.updateSanctionScore()
			if !reflect.DeepEqual(tc.expectedIslandScores, judiciaryInst.sanctionRecord) {
				t.Errorf("Expected %v got %v", tc.expectedIslandScores, judiciaryInst.sanctionRecord)
			}
		})
	}
}

func TestApplySanctions(t *testing.T) {
	cases := []struct {
		name               string
		sanctionRecord     map[shared.ClientID]shared.IIGOSanctionsScore
		sanctionThresholds map[shared.IIGOSanctionsTier]shared.IIGOSanctionsScore
		expectedSanctions  []shared.Sanction
		sanctionLength     int
	}{
		{
			name: "Basic sanction scenario",
			sanctionRecord: map[shared.ClientID]shared.IIGOSanctionsScore{
				shared.Team1: 4,
			},
			sanctionThresholds: getDefaultSanctionThresholds(),
			sanctionLength:     2,
			expectedSanctions: []shared.Sanction{
				{
					ClientID:     shared.Team1,
					SanctionTier: shared.SanctionTier1,
					TurnsLeft:    2,
				},
			},
		},
		{
			name: "Normal sanction scenario",
			sanctionRecord: map[shared.ClientID]shared.IIGOSanctionsScore{
				shared.Team1: 4,
				shared.Team2: 5,
				shared.Team3: 10,
			},
			sanctionThresholds: getDefaultSanctionThresholds(),
			sanctionLength:     2,
			expectedSanctions: []shared.Sanction{
				{
					ClientID:     shared.Team1,
					SanctionTier: shared.SanctionTier1,
					TurnsLeft:    2,
				},
				{
					ClientID:     shared.Team2,
					SanctionTier: shared.SanctionTier2,
					TurnsLeft:    2,
				},
				{
					ClientID:     shared.Team3,
					SanctionTier: shared.SanctionTier3,
					TurnsLeft:    2,
				},
			},
		},
		{
			name: "Custom sanction thresholds",
			sanctionRecord: map[shared.ClientID]shared.IIGOSanctionsScore{
				shared.Team1: 4,
				shared.Team2: 5,
				shared.Team3: 10,
			},
			sanctionThresholds: map[shared.IIGOSanctionsTier]shared.IIGOSanctionsScore{
				shared.SanctionTier1: 5,
				shared.SanctionTier2: 8,
				shared.SanctionTier3: 9,
				shared.SanctionTier4: 10,
				shared.SanctionTier5: 30,
			},
			sanctionLength: 5,
			expectedSanctions: []shared.Sanction{
				{
					ClientID:     shared.Team1,
					SanctionTier: shared.NoSanction,
					TurnsLeft:    5,
				},
				{
					ClientID:     shared.Team2,
					SanctionTier: shared.SanctionTier1,
					TurnsLeft:    5,
				},
				{
					ClientID:     shared.Team3,
					SanctionTier: shared.SanctionTier4,
					TurnsLeft:    5,
				},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			judiciaryInst := defaultInitJudiciary()
			judiciaryInst.sanctionRecord = tc.sanctionRecord
			judiciaryInst.sanctionThresholds = tc.sanctionThresholds
			judiciaryInst.gameConf.SanctionLength = uint(tc.sanctionLength)
			judiciaryInst.applySanctions()
			if !checkListOfSanctionEquals(tc.expectedSanctions, judiciaryInst.gameState.IIGOSanctionCache[0]) {
				t.Errorf("Expected %v got %v", tc.expectedSanctions, judiciaryInst.gameState.IIGOSanctionCache[0])
			}
		})
	}
}

/// Unit tests ///

func TestCreateBroadcastForSanction(t *testing.T) {
	cases := []struct {
		name          string
		clientID      shared.ClientID
		tier          shared.IIGOSanctionsTier
		expectedComms map[shared.CommunicationFieldName]shared.CommunicationContent
	}{
		{
			name:     "Basic broadcast",
			clientID: shared.Team1,
			tier:     shared.SanctionTier4,
			expectedComms: map[shared.CommunicationFieldName]shared.CommunicationContent{
				shared.SanctionClientID: {
					T:           shared.CommunicationInt,
					IntegerData: int(shared.Team1),
				},
				shared.IIGOSanctionTier: {
					T:           shared.CommunicationInt,
					IntegerData: int(shared.SanctionTier4),
				},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res := createBroadcastForSanction(tc.clientID, tc.tier)
			if !reflect.DeepEqual(tc.expectedComms, res) {
				t.Errorf("Expected %v got %v", tc.expectedComms, res)
			}
		})
	}
}

func TestCreateBroadcastsForSanctionThresholds(t *testing.T) {
	cases := []struct {
		name                   string
		sanctionThresholds     map[shared.IIGOSanctionsTier]shared.IIGOSanctionsScore
		expectedCommunications []map[shared.CommunicationFieldName]shared.CommunicationContent
	}{
		{
			name: "Sanction thresholds test",
			sanctionThresholds: map[shared.IIGOSanctionsTier]shared.IIGOSanctionsScore{
				shared.SanctionTier1: 1,
			},
			expectedCommunications: []map[shared.CommunicationFieldName]shared.CommunicationContent{
				{
					shared.IIGOSanctionTier: {
						T:           shared.CommunicationInt,
						IntegerData: 0,
					},
					shared.RuleSanctionPenalty: {
						T:           shared.CommunicationInt,
						IntegerData: 1,
					},
				},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res := createBroadcastsForSanctionThresholds(tc.sanctionThresholds)
			if !reflect.DeepEqual(tc.expectedCommunications, res) {
				t.Errorf("Expected %v got %v", tc.expectedCommunications, res)
			}
		})
	}
}

func TestCreateBroadcastsForRuleViolationPenalties(t *testing.T) {
	cases := []struct {
		name                   string
		violationPenalties     map[string]shared.IIGOSanctionsScore
		expectedCommunications []map[shared.CommunicationFieldName]shared.CommunicationContent
	}{
		{
			name: "Sanction thresholds test",
			violationPenalties: map[string]shared.IIGOSanctionsScore{
				"inspect_allocation_rule": 50,
			},
			expectedCommunications: []map[shared.CommunicationFieldName]shared.CommunicationContent{
				{
					shared.RuleName: {
						T:        shared.CommunicationString,
						TextData: "inspect_allocation_rule",
					},
					shared.RuleSanctionPenalty: {
						T:           shared.CommunicationInt,
						IntegerData: int(50),
					},
				},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res := createBroadcastsForRuleViolationPenalties(tc.violationPenalties)
			if !reflect.DeepEqual(tc.expectedCommunications, res) {
				t.Errorf("Expected %v got %v", tc.expectedCommunications, res)
			}
		})
	}
}

// TestRunEvaluationRulesOnSanctions checks whether the runEvaluationsRulesOnSanctions is able to correctly use the
// sanction evaluator to calculate how much each island should be paying in sanctions
func TestRunEvaluationRulesOnSanctions(t *testing.T) {
	cases := []struct {
		name                    string
		localSanctionCache      map[int][]shared.Sanction
		reportedIslandResources map[shared.ClientID]shared.ResourcesReport
		rulesCache              map[string]rules.RuleMatrix
		expectedSanctions       map[shared.ClientID]shared.Resources
	}{
		{
			name:               "Basic evaluations: no sanction",
			localSanctionCache: DefaultInitLocalSanctionCache(3),
			reportedIslandResources: map[shared.ClientID]shared.ResourcesReport{
				shared.Team1: {
					ReportedAmount: 50,
					Reported:       true,
				},
			},
			rulesCache:        generateRuleStore(),
			expectedSanctions: map[shared.ClientID]shared.Resources{},
		},
		{
			name: "Basic evaluations: singleSanction",
			localSanctionCache: augmentBasicSanctionCache(0, []shared.Sanction{
				{
					ClientID:     shared.Team1,
					SanctionTier: shared.SanctionTier1,
					TurnsLeft:    2,
				},
			}),
			reportedIslandResources: map[shared.ClientID]shared.ResourcesReport{
				shared.Team1: {
					ReportedAmount: 50,
					Reported:       true,
				},
			},
			rulesCache: generateRuleStore(),
			expectedSanctions: map[shared.ClientID]shared.Resources{
				shared.Team1: 0,
			},
		},
		{
			name: "Basic evaluations: multiple sanction",
			localSanctionCache: augmentBasicSanctionCache(0, []shared.Sanction{
				{
					ClientID:     shared.Team1,
					SanctionTier: shared.SanctionTier2,
					TurnsLeft:    2,
				},
				{
					ClientID:     shared.Team3,
					SanctionTier: shared.SanctionTier3,
					TurnsLeft:    2,
				},
			}),
			reportedIslandResources: map[shared.ClientID]shared.ResourcesReport{
				shared.Team1: {
					ReportedAmount: 50,
					Reported:       true,
				},
				shared.Team3: {
					ReportedAmount: 100,
					Reported:       true,
				},
			},
			rulesCache: generateRuleStore(),
			expectedSanctions: map[shared.ClientID]shared.Resources{
				shared.Team1: 5,
				shared.Team3: 30,
			},
		},
		{
			name: "Evaluations: multiple sanction with timeout",
			localSanctionCache: augmentBasicSanctionCache(0, []shared.Sanction{
				{
					ClientID:     shared.Team1,
					SanctionTier: shared.SanctionTier2,
					TurnsLeft:    2,
				},
				{
					ClientID:     shared.Team3,
					SanctionTier: shared.SanctionTier3,
					TurnsLeft:    2,
				},
				{
					ClientID:     shared.Team4,
					SanctionTier: shared.SanctionTier3,
					TurnsLeft:    0,
				},
			}),
			reportedIslandResources: map[shared.ClientID]shared.ResourcesReport{
				shared.Team1: {
					Reported:       true,
					ReportedAmount: 50,
				},
				shared.Team3: {
					ReportedAmount: 100,
					Reported:       true,
				},
				shared.Team4: {
					ReportedAmount: 150,
					Reported:       true,
				},
			},
			rulesCache: generateRuleStore(),
			expectedSanctions: map[shared.ClientID]shared.Resources{
				shared.Team1: 5,
				shared.Team3: 30,
				shared.Team4: 0,
			},
		},
		{
			name: "Evaluations: complex case",
			localSanctionCache: augmentBasicSanctionCache(1, []shared.Sanction{
				{
					ClientID:     shared.Team1,
					SanctionTier: shared.SanctionTier2,
					TurnsLeft:    2,
				},
				{
					ClientID:     shared.Team2,
					SanctionTier: shared.SanctionTier2,
					TurnsLeft:    2,
				},
				{
					ClientID:     shared.Team3,
					SanctionTier: shared.SanctionTier3,
					TurnsLeft:    2,
				},
				{
					ClientID:     shared.Team4,
					SanctionTier: shared.SanctionTier4,
					TurnsLeft:    0,
				},
				{
					ClientID:     shared.Team5,
					SanctionTier: shared.SanctionTier5,
					TurnsLeft:    3,
				},
				{
					ClientID:     shared.Team6,
					SanctionTier: shared.SanctionTier1,
					TurnsLeft:    0,
				},
			}),
			reportedIslandResources: map[shared.ClientID]shared.ResourcesReport{
				shared.Team1: {
					ReportedAmount: 50,
					Reported:       true,
				},
				shared.Team2: {
					ReportedAmount: 50,
					Reported:       true,
				},
				shared.Team3: {
					ReportedAmount: 100,
					Reported:       true,
				},
				shared.Team4: {
					ReportedAmount: 150,
					Reported:       true,
				},
				shared.Team5: {
					ReportedAmount: 20,
					Reported:       true,
				},
				shared.Team6: {
					ReportedAmount: 100,
					Reported:       true,
				},
			},
			rulesCache: generateRuleStore(),
			expectedSanctions: map[shared.ClientID]shared.Resources{
				shared.Team1: 5,
				shared.Team2: 5,
				shared.Team3: 30,
				shared.Team4: 0,
				shared.Team5: 16,
				shared.Team6: 0,
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res := runEvaluationRulesOnSanctions(tc.localSanctionCache, tc.reportedIslandResources, tc.rulesCache, 500)
			if !reflect.DeepEqual(res, tc.expectedSanctions) {
				t.Errorf("Expected %v got %v", tc.expectedSanctions, res)
			}
		})
	}
}

// TestDecrementSanctionTime checks whether the decrementSanctionTime function is able to reduce the time
// on the given sanctions by one time step.
func TestDecrementSanctionTime(t *testing.T) {
	cases := []struct {
		name             string
		initialSanctions map[int][]shared.Sanction
		updatedSanctions map[int][]shared.Sanction
	}{
		{
			name: "Basic Decrement Test",
			initialSanctions: map[int][]shared.Sanction{
				1: {
					{
						ClientID:     shared.Team1,
						SanctionTier: shared.SanctionTier2,
						TurnsLeft:    2,
					},
				},
			},
			updatedSanctions: map[int][]shared.Sanction{
				1: {
					{
						ClientID:     shared.Team1,
						SanctionTier: shared.SanctionTier2,
						TurnsLeft:    1,
					},
				},
			},
		},
		{
			name: "Complex Decrement Test",
			initialSanctions: map[int][]shared.Sanction{
				1: {
					{
						ClientID:     shared.Team1,
						SanctionTier: shared.SanctionTier2,
						TurnsLeft:    2,
					},
					{
						ClientID:     shared.Team3,
						SanctionTier: shared.SanctionTier4,
						TurnsLeft:    8,
					},
				},
				2: {
					{
						ClientID:     shared.Team4,
						SanctionTier: shared.SanctionTier4,
						TurnsLeft:    1,
					},
				},
			},
			updatedSanctions: map[int][]shared.Sanction{
				1: {
					{
						ClientID:     shared.Team1,
						SanctionTier: shared.SanctionTier2,
						TurnsLeft:    1,
					},
					{
						ClientID:     shared.Team3,
						SanctionTier: shared.SanctionTier4,
						TurnsLeft:    7,
					},
				},
				2: {
					{
						ClientID:     shared.Team4,
						SanctionTier: shared.SanctionTier4,
						TurnsLeft:    0,
					},
				},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res := decrementSanctionTime(tc.initialSanctions)
			if !reflect.DeepEqual(res, tc.updatedSanctions) {
				t.Errorf("Expected %v got %v", tc.updatedSanctions, res)
			}
		})
	}
}

// TestMergeEvalResults checks whether the mergeEvaluationReturn function can perform a soft merge of two maps of
// the type map[shared.ClientID]shared.EvaluationReturn
func TestMergeEvalResults(t *testing.T) {
	availableRules := generateDummyRuleMatrices()
	cases := []struct {
		name   string
		set1   map[shared.ClientID]shared.EvaluationReturn
		set2   map[shared.ClientID]shared.EvaluationReturn
		expect map[shared.ClientID]shared.EvaluationReturn
	}{
		{
			name: "Simple merge test",
			set1: map[shared.ClientID]shared.EvaluationReturn{
				shared.Team1: {
					Rules: []rules.RuleMatrix{
						availableRules[0],
						availableRules[1],
					},
					Evaluations: []bool{true, false},
				},
			},
			set2: map[shared.ClientID]shared.EvaluationReturn{},
			expect: map[shared.ClientID]shared.EvaluationReturn{
				shared.Team1: {
					Rules: []rules.RuleMatrix{
						availableRules[0],
						availableRules[1],
					},
					Evaluations: []bool{true, false},
				},
			},
		},
		{
			name: "Complex Merge Test",
			set1: map[shared.ClientID]shared.EvaluationReturn{
				shared.Team1: {
					Rules: []rules.RuleMatrix{
						availableRules[0],
						availableRules[1],
					},
					Evaluations: []bool{true, false},
				},
			},
			set2: map[shared.ClientID]shared.EvaluationReturn{
				shared.Team2: {
					Rules: []rules.RuleMatrix{
						availableRules[2],
						availableRules[3],
					},
					Evaluations: []bool{true, false},
				},
			},
			expect: map[shared.ClientID]shared.EvaluationReturn{
				shared.Team1: {
					Rules: []rules.RuleMatrix{
						availableRules[0],
						availableRules[1],
					},
					Evaluations: []bool{true, false},
				},
				shared.Team2: {
					Rules: []rules.RuleMatrix{
						availableRules[2],
						availableRules[3],
					},
					Evaluations: []bool{true, false},
				},
			},
		},
		{
			name: "Patchwork merge test",
			set1: map[shared.ClientID]shared.EvaluationReturn{
				shared.Team1: {
					Rules: []rules.RuleMatrix{
						availableRules[0],
						availableRules[1],
					},
					Evaluations: []bool{true, false},
				},
			},
			set2: map[shared.ClientID]shared.EvaluationReturn{
				shared.Team1: {
					Rules: []rules.RuleMatrix{
						availableRules[0],
						availableRules[1],
					},
					Evaluations: []bool{true, false},
				},
				shared.Team2: {
					Rules: []rules.RuleMatrix{
						availableRules[2],
						availableRules[3],
					},
					Evaluations: []bool{true, false},
				},
			},
			expect: map[shared.ClientID]shared.EvaluationReturn{
				shared.Team1: {
					Rules: []rules.RuleMatrix{
						availableRules[0],
						availableRules[1],
						availableRules[0],
						availableRules[1],
					},
					Evaluations: []bool{true, false, true, false},
				},
				shared.Team2: {
					Rules: []rules.RuleMatrix{
						availableRules[2],
						availableRules[3],
					},
					Evaluations: []bool{true, false},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res := mergeEvaluationReturn(tc.set1, tc.set2)
			if !reflect.DeepEqual(res, tc.expect) {
				t.Errorf("Expected %v got %v", tc.expect, res)
			}
		})
	}
}

// TestSearchEvalReturnForRuleName checks whether the searchEvalReturnForRuleName is able to pick up which rules
// have been evaluated in a set of EvaluationReturns
func TestSearchEvalReturnForRuleName(t *testing.T) {
	availableRules := generateDummyRuleMatrices()
	cases := []struct {
		name     string
		ruleName string
		evalRet  shared.EvaluationReturn
		expected bool
	}{
		{
			name:     "Simple search",
			ruleName: "iigo_economic_sanction_5",
			evalRet: shared.EvaluationReturn{
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
			evalRet: shared.EvaluationReturn{
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
			evalRet: shared.EvaluationReturn{
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
		evalRes  map[shared.ClientID]shared.EvaluationReturn
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
			evalRes: map[shared.ClientID]shared.EvaluationReturn{
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
			evalRes: map[shared.ClientID]shared.EvaluationReturn{
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
			evalRes: map[shared.ClientID]shared.EvaluationReturn{
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
			arr, res := rules.PickUpRulesByVariable(tc.searchVar, tc.ruleStore, generateDummyVariableCache())
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
	cases := []struct {
		name          string
		sanctionCache map[int][]shared.Sanction
		pardons       map[int][]bool
		expValidity   bool
		expComms      map[shared.ClientID][]map[shared.CommunicationFieldName]shared.CommunicationContent
		expFinCache   map[int][]shared.Sanction
	}{
		{
			name: "No Pardons Check",
			sanctionCache: map[int][]shared.Sanction{
				1: {
					{
						ClientID:     shared.Team1,
						SanctionTier: shared.SanctionTier3,
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
			expComms:    generateEmptyCommunicationsMap(shared.TeamIDs),
			expFinCache: map[int][]shared.Sanction{
				1: {
					{
						ClientID:     shared.Team1,
						SanctionTier: shared.SanctionTier3,
						TurnsLeft:    3,
					},
				},
			},
		},
		{
			name: "Simple Pardons Case",
			sanctionCache: map[int][]shared.Sanction{
				1: {
					{
						ClientID:     shared.Team1,
						SanctionTier: shared.SanctionTier3,
						TurnsLeft:    3,
					},
					{
						ClientID:     shared.Team2,
						SanctionTier: shared.SanctionTier3,
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
							IntegerData: int(shared.SanctionTier3),
						},
					},
				},
				shared.Team2: {},
				shared.Team3: {},
				shared.Team4: {},
				shared.Team5: {},
				shared.Team6: {},
			},
			expFinCache: map[int][]shared.Sanction{
				1: {
					{
						ClientID:     shared.Team2,
						SanctionTier: shared.SanctionTier3,
						TurnsLeft:    3,
					},
				},
			},
		},
		{
			name: "Complex Pardons Case",
			sanctionCache: map[int][]shared.Sanction{
				1: {
					{
						ClientID:     shared.Team1,
						SanctionTier: shared.SanctionTier3,
						TurnsLeft:    3,
					},
					{
						ClientID:     shared.Team2,
						SanctionTier: shared.SanctionTier3,
						TurnsLeft:    3,
					},
					{
						ClientID:     shared.Team3,
						SanctionTier: shared.SanctionTier2,
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
							IntegerData: int(shared.SanctionTier3),
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
							IntegerData: int(shared.SanctionTier2),
						},
					},
				},
				shared.Team2: {},
				shared.Team4: {},
				shared.Team5: {},
				shared.Team6: {},
			},
			expFinCache: map[int][]shared.Sanction{
				1: {
					{
						ClientID:     shared.Team2,
						SanctionTier: shared.SanctionTier3,
						TurnsLeft:    3,
					},
				},
			},
		},
		{
			name: "Complex Pardons Failure Case",
			sanctionCache: map[int][]shared.Sanction{
				1: {
					{
						ClientID:     shared.Team1,
						SanctionTier: shared.SanctionTier3,
						TurnsLeft:    3,
					},
					{
						ClientID:     shared.Team2,
						SanctionTier: shared.SanctionTier3,
						TurnsLeft:    3,
					},
					{
						ClientID:     shared.Team3,
						SanctionTier: shared.SanctionTier2,
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
			expFinCache: map[int][]shared.Sanction{
				1: {
					{
						ClientID:     shared.Team1,
						SanctionTier: shared.SanctionTier3,
						TurnsLeft:    3,
					},
					{
						ClientID:     shared.Team2,
						SanctionTier: shared.SanctionTier3,
						TurnsLeft:    3,
					},
					{
						ClientID:     shared.Team3,
						SanctionTier: shared.SanctionTier2,
						TurnsLeft:    2,
					},
				},
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			val, finalCache, comms := implementPardons(tc.sanctionCache, tc.pardons, shared.TeamIDs)
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

// TestUnpackSingleIslandTransgression checks whether unpackSingleIslandTransgressions is able to check an
// EvaluationReturn data-structure and collect all rules that have found to be broken
func TestUnpackSingleIslandTransgression(t *testing.T) {
	availableRules := generateDummyRuleMatrices()
	cases := []struct {
		name       string
		evalReturn shared.EvaluationReturn
		expected   []string
	}{
		{
			name: "Basic unpack test",
			evalReturn: shared.EvaluationReturn{
				Rules: []rules.RuleMatrix{
					availableRules[0],
				},
				Evaluations: []bool{false},
			},
			expected: []string{
				"inspect_ballot_rule",
			},
		},
		{
			name: "None to unpack",
			evalReturn: shared.EvaluationReturn{
				Rules: []rules.RuleMatrix{
					availableRules[0],
				},
				Evaluations: []bool{true},
			},
			expected: []string{},
		},
		{
			name: "Multiple unpack",
			evalReturn: shared.EvaluationReturn{
				Rules: []rules.RuleMatrix{
					availableRules[0],
					availableRules[1],
					availableRules[2],
				},
				Evaluations: []bool{false, false, false},
			},
			expected: []string{availableRules[0].RuleName, availableRules[1].RuleName, availableRules[2].RuleName},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tra := unpackSingleIslandTransgressions(tc.evalReturn)
			if !reflect.DeepEqual(tra, tc.expected) {
				t.Errorf("Expected final transgressions to be %v got %v", tc.expected, tra)
			}
		})
	}
}

// TestSoftMergeSanctionThreshold checks whether the softMergeSanctionThresholds is able to correctly merge
// the clients threshols with the default ones
func TestSoftMergeSanctionThreshold(t *testing.T) {
	cases := []struct {
		name              string
		clientSanctionMap map[shared.IIGOSanctionsTier]shared.IIGOSanctionsScore
		expectedVal       map[shared.IIGOSanctionsTier]shared.IIGOSanctionsScore
	}{
		{
			name: "Basic soft merge",
			clientSanctionMap: map[shared.IIGOSanctionsTier]shared.IIGOSanctionsScore{
				shared.SanctionTier1: shared.IIGOSanctionsScore(3),
			},
			expectedVal: map[shared.IIGOSanctionsTier]shared.IIGOSanctionsScore{
				shared.SanctionTier1: 3,
				shared.SanctionTier2: 5,
				shared.SanctionTier3: 10,
				shared.SanctionTier4: 20,
				shared.SanctionTier5: 30,
			},
		},
		{
			name:              "No merge",
			clientSanctionMap: map[shared.IIGOSanctionsTier]shared.IIGOSanctionsScore{},
			expectedVal: map[shared.IIGOSanctionsTier]shared.IIGOSanctionsScore{
				shared.SanctionTier1: 1,
				shared.SanctionTier2: 5,
				shared.SanctionTier3: 10,
				shared.SanctionTier4: 20,
				shared.SanctionTier5: 30,
			},
		},
		{
			name: "More complicated merge",
			clientSanctionMap: map[shared.IIGOSanctionsTier]shared.IIGOSanctionsScore{
				shared.SanctionTier1: 7,
				shared.SanctionTier2: 9,
				shared.SanctionTier5: 400,
			},
			expectedVal: map[shared.IIGOSanctionsTier]shared.IIGOSanctionsScore{
				shared.SanctionTier1: 7,
				shared.SanctionTier2: 9,
				shared.SanctionTier3: 10,
				shared.SanctionTier4: 20,
				shared.SanctionTier5: 400,
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
		{
			VariableName: rules.RuleSelected,
			Values:       []float64{0},
		},
		{
			VariableName: rules.VoteCalled,
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
		{
			name: "tax_decision",
			reqVar: []rules.VariableFieldName{
				rules.TaxDecisionMade,
			},
			v:       []float64{1, -1},
			aux:     []float64{0},
			mutable: false,
		},
		{
			name: "allocation_decision",
			reqVar: []rules.VariableFieldName{
				rules.AllocationMade,
			},
			v:       []float64{1, -1},
			aux:     []float64{0},
			mutable: false,
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

func TestJudgeIncureServiceCharge(t *testing.T) {
	cases := []struct {
		name                string
		bJudge              judiciary // base
		input               shared.Resources
		expectedReturn      bool
		expectedCommonPool  shared.Resources
		expectedJudgeBudget shared.Resources
	}{
		{
			name: "Excess pay",
			bJudge: judiciary{
				JudgeID: shared.Team1,
				gameState: &gamestate.GameState{
					CommonPool: 400,
					IIGORolesBudget: map[shared.Role]shared.Resources{
						shared.President: 10,
						shared.Speaker:   10,
						shared.Judge:     100,
					},
				},
			},
			input:               50,
			expectedReturn:      true,
			expectedCommonPool:  350,
			expectedJudgeBudget: 50,
		},
		{
			name: "Negative Budget",
			bJudge: judiciary{
				JudgeID: shared.Team1,
				gameState: &gamestate.GameState{
					CommonPool: 400,
					IIGORolesBudget: map[shared.Role]shared.Resources{
						shared.President: 10,
						shared.Speaker:   10,
						shared.Judge:     10,
					},
				},
			},
			input:               50,
			expectedReturn:      true,
			expectedCommonPool:  350,
			expectedJudgeBudget: -40,
		},
		{
			name: "Limited common pool",
			bJudge: judiciary{
				JudgeID: shared.Team1,
				gameState: &gamestate.GameState{
					CommonPool: 40,
					IIGORolesBudget: map[shared.Role]shared.Resources{
						shared.President: 10,
						shared.Speaker:   10,
						shared.Judge:     10,
					},
				},
			},
			input:               50,
			expectedReturn:      false,
			expectedCommonPool:  40,
			expectedJudgeBudget: 10,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			returned := tc.bJudge.incurServiceCharge(tc.input)
			commonPool := tc.bJudge.gameState.CommonPool
			judgeBudget := tc.bJudge.gameState.IIGORolesBudget[shared.Judge]
			if returned != tc.expectedReturn ||
				commonPool != tc.expectedCommonPool ||
				judgeBudget != tc.expectedJudgeBudget {
				t.Errorf("%v - Failed. Got '%v, %v, %v', but expected '%v, %v, %v'",
					tc.name, returned, commonPool, judgeBudget,
					tc.expectedReturn, tc.expectedCommonPool, tc.expectedJudgeBudget)
			}
		})
	}
}

func defaultInitJudiciary() judiciary {
	var logging shared.Logger = func(format string, a ...interface{}) {}
	gamestate := gamestate.GameState{
		CommonPool: 999,
		IIGORolesBudget: map[shared.Role]shared.Resources{
			shared.President: 100,
			shared.Speaker:   10,
			shared.Judge:     10,
		},
		IIGORoleMonitoringCache: []shared.Accountability{},
		RulesBrokenByIslands:    map[shared.ClientID][]string{},
	}
	return judiciary{
		JudgeID:               0,
		evaluationResults:     map[shared.ClientID]shared.EvaluationReturn{},
		clientJudge:           &baseclient.BaseJudge{},
		sanctionRecord:        map[shared.ClientID]shared.IIGOSanctionsScore{},
		sanctionThresholds:    map[shared.IIGOSanctionsTier]shared.IIGOSanctionsScore{},
		ruleViolationSeverity: map[string]shared.IIGOSanctionsScore{},
		gameConf: &config.IIGOConfig{
			DefaultSanctionScore: 5,
		},
		gameState: &gamestate,
		monitoring: &monitor{
			gameState: &gamestate,
		},
		logger: logging,
	}
}

func checkListOfSanctionEquals(list1 []shared.Sanction, list2 []shared.Sanction) bool {
	map1 := map[shared.Sanction]int{}
	map2 := map[shared.Sanction]int{}
	for _, val := range list1 {
		map1[val] = 0
	}
	for _, val := range list2 {
		map2[val] = 0
	}
	return reflect.DeepEqual(map1, map2)
}

func augmentBasicSanctionCache(time int, additionalSanctions []shared.Sanction) map[int][]shared.Sanction {
	basicCache := DefaultInitLocalSanctionCache(3)
	basicCache[time] = append(basicCache[time], additionalSanctions...)
	return basicCache
}
