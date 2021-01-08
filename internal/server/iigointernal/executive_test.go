package iigointernal

import (
	"reflect" // Used to compare two maps
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

type allocInput struct {
	resourceRequests   map[shared.ClientID]shared.Resources
	commonPoolResource shared.Resources
}

func TestAllocationRequests(t *testing.T) {
	cases := []struct {
		name  string
		input allocInput
		reply map[shared.ClientID]shared.Resources
		want  bool
	}{
		{
			name: "Limited resources",
			input: allocInput{resourceRequests: map[shared.ClientID]shared.Resources{
				shared.Team1: 5,
				shared.Team2: 10,
				shared.Team3: 15,
				shared.Team4: 20,
				shared.Team5: 25,
				shared.Team6: 30,
			}, commonPoolResource: 100},
			reply: map[shared.ClientID]shared.Resources{
				shared.Team1: calcExpectedVal(5, 105, 100),
				shared.Team2: calcExpectedVal(10, 105, 100),
				shared.Team3: calcExpectedVal(15, 105, 100),
				shared.Team4: calcExpectedVal(20, 105, 100),
				shared.Team5: calcExpectedVal(25, 105, 100),
				shared.Team6: calcExpectedVal(30, 105, 100),
			},
			want: true,
		},
		{
			name: "Enough resources",
			input: allocInput{resourceRequests: map[shared.ClientID]shared.Resources{
				shared.Team1: 5,
				shared.Team2: 10,
				shared.Team3: 15,
				shared.Team4: 20,
				shared.Team5: 25,
				shared.Team6: 30,
			}, commonPoolResource: 150},
			reply: map[shared.ClientID]shared.Resources{
				shared.Team1: 5,
				shared.Team2: 10,
				shared.Team3: 15,
				shared.Team4: 20,
				shared.Team5: 25,
				shared.Team6: 30,
			},
			want: true,
		},
		{
			name: "No resources",
			input: allocInput{resourceRequests: map[shared.ClientID]shared.Resources{
				shared.Team1: 5,
				shared.Team2: 10,
				shared.Team3: 15,
				shared.Team4: 20,
				shared.Team5: 25,
				shared.Team6: 30,
			}, commonPoolResource: 0},
			reply: map[shared.ClientID]shared.Resources{
				shared.Team1: 0,
				shared.Team2: 0,
				shared.Team3: 0,
				shared.Team4: 0,
				shared.Team5: 0,
				shared.Team6: 0,
			},
			want: true,
		},
	}

	wantPresidentReturnType := shared.PresidentAllocation
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			executive := &executive{
				PresidentID:     shared.Team1,
				clientPresident: &baseclient.BasePresident{},
			}

			got := executive.clientPresident.EvaluateAllocationRequests(tc.input.resourceRequests, tc.input.commonPoolResource)

			if got.ContentType != wantPresidentReturnType {
				t.Errorf("%v - Failed. Expected action type  %v, got action %v", tc.name, wantPresidentReturnType, got.ContentType)
			}

			if got.ActionTaken && tc.want {
				if !reflect.DeepEqual(tc.reply, got.ResourceMap) {
					t.Errorf("%v - Failed. Expected %v, got %v", tc.name, tc.reply, got.ResourceMap)
				}
			} else {
				if got.ActionTaken != tc.want {
					t.Errorf("%v - Failed. Returned error '%v' which was not equal to expected '%v'", tc.name, got.ActionTaken, tc.want)
				}
			}
		})
	}
}

func calcExpectedVal(resourcesRequested shared.Resources, totalResourcesRequested shared.Resources, availableInPool shared.Resources) shared.Resources {
	return 0.75 * (resourcesRequested * availableInPool) / totalResourcesRequested
}

func checkIfInList(s rules.RuleMatrix, lst []rules.RuleMatrix) bool {
	for _, i := range lst {
		if reflect.DeepEqual(s, i) {
			return true
		}
	}
	return false
}
func TestPickRuleToVote(t *testing.T) {
	cases := []struct {
		name  string
		input []rules.RuleMatrix
		reply string
		want  bool
	}{
		{
			name:  "Basic rule",
			input: []rules.RuleMatrix{{RuleName: "rule"}},
			reply: "rule",
			want:  true,
		},
		{
			name:  "Empty rule",
			input: []rules.RuleMatrix{{}},
			reply: "",
			want:  true,
		},
		{
			name:  "Longer list",
			input: []rules.RuleMatrix{{RuleName: "Somas"}, {RuleName: "2020"}, {RuleName: "Internal"}, {RuleName: "Server"}, {RuleName: "Roles"}, {RuleName: "President"}},
			reply: "",
			want:  true,
		},
		{
			name:  "Empty list",
			input: []rules.RuleMatrix{},
			reply: "",
			want:  false,
		},
	}

	wantPresidentReturnType := shared.PresidentRuleProposal
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			executive := &executive{
				clientPresident: &baseclient.BasePresident{},
			}

			got := executive.clientPresident.PickRuleToVote(tc.input)

			if got.ContentType != wantPresidentReturnType {
				t.Errorf("%v - Failed. Expected return type  %v, got %v", tc.name, wantPresidentReturnType, got.ContentType)
			}

			if got.ActionTaken && tc.want {
				if len(tc.input) == 0 {
					// if got.ProposedRuleMatrix != "" {
					// 	t.Errorf("%v - Failed. Returned '%v', but expectd an empty string", tc.name, got.ProposedRuleMatrix)
					if !got.ProposedRuleMatrix.RuleMatrixIsEmpty() { //
						t.Errorf("%v - Failed. Returned '%v', but expectd an empty ruleMatrix", tc.name, got.ProposedRuleMatrix)
					}
				} else if !checkIfInList(got.ProposedRuleMatrix, tc.input) {
					t.Errorf("%v - Failed. Returned '%v', expected '%v'", tc.name, got.ProposedRuleMatrix, tc.input)
				}
			} else {
				if got.ActionTaken != tc.want {
					t.Errorf("%v - Failed. Returned error '%v' which was not equal to expected '%v'", tc.name, got.ActionTaken, tc.want)
				}
			}
		})
	}
}

func TestSetRuleProposals(t *testing.T) {
	cases := []struct {
		name     string
		input    []rules.RuleMatrix
		expected []rules.RuleMatrix
	}{
		{
			name:     "Single Rule",
			input:    []rules.RuleMatrix{{RuleName: "rule"}},
			expected: []rules.RuleMatrix{{RuleName: "rule"}},
		},
		{
			name:     "No Rule",
			input:    []rules.RuleMatrix{{RuleName: ""}},
			expected: []rules.RuleMatrix{{RuleName: ""}},
		},
		{
			name:     "Multiple Rules",
			input:    []rules.RuleMatrix{{RuleName: "Somas"}, {RuleName: "2020"}, {RuleName: "Internal"}, {RuleName: "Server"}, {RuleName: "Roles"}, {RuleName: "President"}},
			expected: []rules.RuleMatrix{{RuleName: "Somas"}, {RuleName: "2020"}, {RuleName: "Internal"}, {RuleName: "Server"}, {RuleName: "Roles"}, {RuleName: "President"}},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			executive := &executive{
				PresidentID:    shared.Team1,
				RulesProposals: []rules.RuleMatrix{},
			}

			executive.setRuleProposals(tc.input)

			if !reflect.DeepEqual(executive.RulesProposals, tc.expected) {
				t.Errorf("%v - Failed. RulesProposals set to '%v', expected '%v'", tc.name, executive.RulesProposals, tc.expected)
			}

		})
	}
}

func TestGetTaxMap(t *testing.T) {
	cases := []struct {
		name           string
		input          map[shared.ClientID]shared.ResourcesReport
		bPresident     executive // base
		expectedLength int
		expected       map[shared.ClientID]shared.Resources
	}{
		{
			name:  "Empty tax map base",
			input: map[shared.ClientID]shared.ResourcesReport{},
			bPresident: executive{
				PresidentID:     3,
				clientPresident: &baseclient.BasePresident{},
			},
			expectedLength: 0,
		},
		{
			name: "Short tax map base",
			input: map[shared.ClientID]shared.ResourcesReport{
				1: {ReportedAmount: 5, Reported: true},
			},
			bPresident: executive{
				PresidentID:     3,
				clientPresident: &baseclient.BasePresident{},
			},
			expectedLength: 1,
		},
		{
			name: "Long tax map base",
			input: map[shared.ClientID]shared.ResourcesReport{
				1: {ReportedAmount: 5, Reported: true},
				2: {ReportedAmount: 10, Reported: true},
				3: {ReportedAmount: 15, Reported: true},
				4: {ReportedAmount: 20, Reported: true},
				5: {ReportedAmount: 25, Reported: true},
				6: {ReportedAmount: 30, Reported: true},
			},
			bPresident: executive{
				PresidentID:     4,
				clientPresident: &baseclient.BasePresident{},
			},
			expectedLength: 6,
		},
		{
			name:  "Client empty tax map base",
			input: map[shared.ClientID]shared.ResourcesReport{},
			bPresident: executive{
				PresidentID:     5,
				clientPresident: &baseclient.BasePresident{},
			},
			expectedLength: 0,
		},
		{
			name: "Client short tax map base",
			input: map[shared.ClientID]shared.ResourcesReport{
				1: {ReportedAmount: 5, Reported: true},
			},
			bPresident: executive{
				PresidentID:     5,
				clientPresident: &baseclient.BasePresident{},
			},
			expectedLength: 1,
		},
		{
			name: "Client long tax map base",
			input: map[shared.ClientID]shared.ResourcesReport{
				1: {ReportedAmount: 5, Reported: true},
				2: {ReportedAmount: 10, Reported: true},
				3: {ReportedAmount: 15, Reported: true},
				4: {ReportedAmount: 20, Reported: true},
				5: {ReportedAmount: 25, Reported: true},
				6: {ReportedAmount: 30, Reported: true},
			},
			bPresident: executive{
				PresidentID:     5,
				clientPresident: &baseclient.BasePresident{},
			},
			expectedLength: 6,
		},
		{
			name: "Clients not reporting",
			input: map[shared.ClientID]shared.ResourcesReport{
				1: {ReportedAmount: 5, Reported: false},
				2: {ReportedAmount: 10, Reported: false},
				3: {ReportedAmount: 15, Reported: false},
				4: {ReportedAmount: 20, Reported: false},
				5: {ReportedAmount: 25, Reported: false},
				6: {ReportedAmount: 30, Reported: false},
			},
			bPresident: executive{
				PresidentID:     5,
				clientPresident: &baseclient.BasePresident{},
			},
			expected: map[shared.ClientID]shared.Resources{
				1: 15,
				2: 15,
				3: 15,
				4: 15,
				5: 15,
				6: 15,
			},
		},
	}

	wantPresidentReturnType := shared.PresidentTaxation
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.bPresident.getTaxMap(tc.input)

			if got.ContentType != wantPresidentReturnType {
				t.Errorf("%v - Failed. Expected action type  %v, got action %v", tc.name, wantPresidentReturnType, got.ContentType)
			}

			if tc.expected == nil {
				if len(got.ResourceMap) != tc.expectedLength {
					t.Errorf("%v - Failed. TaxMap set to '%v', expected '%v'", tc.name, got.ResourceMap, tc.input)
				}
			} else {
				if !reflect.DeepEqual(got.ResourceMap, tc.expected) {
					t.Errorf("%v - Failed. TaxMap set to '%v', expected '%v'", tc.name, got.ResourceMap, tc.expected)
				}
			}

		})
	}
}

func TestGetRuleForSpeaker(t *testing.T) {

	fakeGameState := gamestate.GameState{
		CommonPool: 400,
		IIGORolesBudget: map[shared.Role]shared.Resources{
			shared.President: 10,
			shared.Speaker:   10,
			shared.Judge:     10,
		},
		IIGORoleMonitoringCache: []shared.Accountability{},
	}
	cases := []struct {
		name       string
		bPresident executive // base
		expected   []rules.RuleMatrix
	}{
		{
			name: "Empty tax map base",
			bPresident: executive{
				PresidentID:     3,
				RulesProposals:  []rules.RuleMatrix{},
				clientPresident: &baseclient.BasePresident{},
				gameState:       &fakeGameState,
				gameConf:        &config.IIGOConfig{},
				monitoring: &monitor{
					gameState: &fakeGameState,
				},
			},
			expected: []rules.RuleMatrix{{RuleName: ""}},
		},
		{
			name: "Short tax map base",
			bPresident: executive{
				PresidentID:     3,
				RulesProposals:  []rules.RuleMatrix{{RuleName: "test"}},
				clientPresident: &baseclient.BasePresident{},
				gameState:       &fakeGameState,
				gameConf:        &config.IIGOConfig{},
				monitoring: &monitor{
					gameState: &fakeGameState,
				},
			},
			expected: []rules.RuleMatrix{{RuleName: "test"}},
		},
		{
			name: "Long tax map base",
			bPresident: executive{
				PresidentID:     3,
				RulesProposals:  []rules.RuleMatrix{{RuleName: "Somas"}, {RuleName: "2020"}, {RuleName: "Internal"}, {RuleName: "Server"}, {RuleName: "Roles"}, {RuleName: "President"}},
				clientPresident: &baseclient.BasePresident{},
				gameState:       &fakeGameState,
				gameConf:        &config.IIGOConfig{},
				monitoring: &monitor{
					gameState: &fakeGameState,
				},
			},
			expected: []rules.RuleMatrix{{RuleName: "Somas"}, {RuleName: "2020"}, {RuleName: "Internal"}, {RuleName: "Server"}, {RuleName: "Roles"}, {RuleName: "President"}},
		},
		{
			name: "Client empty tax map base",
			bPresident: executive{
				PresidentID:     5,
				RulesProposals:  []rules.RuleMatrix{{RuleName: "Somas"}, {RuleName: "2020"}, {RuleName: "Internal"}, {RuleName: "Server"}, {RuleName: "Roles"}, {RuleName: "President"}},
				clientPresident: &baseclient.BasePresident{},
				gameState:       &fakeGameState,
				gameConf:        &config.IIGOConfig{},
				monitoring: &monitor{
					gameState: &fakeGameState,
				},
			},
			expected: []rules.RuleMatrix{{RuleName: "Somas"}, {RuleName: "2020"}, {RuleName: "Internal"}, {RuleName: "Server"}, {RuleName: "Roles"}, {RuleName: "President"}},
		},
		{
			name: "Client tax map base override",
			bPresident: executive{
				PresidentID:     5,
				RulesProposals:  []rules.RuleMatrix{{RuleName: "Somas"}, {RuleName: "2020"}, {RuleName: "Internal"}, {RuleName: "Server"}, {RuleName: "Roles"}, {RuleName: "President"}},
				clientPresident: &baseclient.BasePresident{},
				gameState:       &fakeGameState,
				gameConf:        &config.IIGOConfig{},
				monitoring: &monitor{
					gameState: &fakeGameState,
				},
			},
			expected: []rules.RuleMatrix{{RuleName: "Somas"}, {RuleName: "2020"}, {RuleName: "Internal"}, {RuleName: "Server"}, {RuleName: "Roles"}, {RuleName: "President"}},
		},
	}

	wantPresidentReturnType := shared.PresidentRuleProposal
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fakeGameConfig := config.IIGOConfig{
				GetRuleForSpeakerActionCost: 1,
			}
			fakeGameState := gamestate.GameState{
				CommonPool: 100,
				IIGORolesBudget: map[shared.Role]shared.Resources{
					shared.President: 10,
					shared.Speaker:   10,
					shared.Judge:     10,
				},
			}

			tc.bPresident.syncWithGame(&fakeGameState, &fakeGameConfig)
			got, _ := tc.bPresident.getRuleForSpeaker()

			if got.ContentType != wantPresidentReturnType {
				t.Errorf("%v - Failed. Expected action type  %v, got action %v", tc.name, wantPresidentReturnType, got.ContentType)
			}

			if len(tc.expected) == 0 {
				// if got.ProposedRuleMatrix != "" {
				// 	t.Errorf("%v - Failed. Returned '%v', but expectd an empty string", tc.name, got.ProposedRuleMatrix)
				if !got.ProposedRuleMatrix.RuleMatrixIsEmpty() { //
					t.Errorf("%v - Failed. Returned '%v', but expectd an empty ruleMatrix", tc.name, got.ProposedRuleMatrix)
				}
			} else if !checkIfInList(got.ProposedRuleMatrix, tc.expected) {
				t.Errorf("%v - Failed. Returned '%v', expected '%v'", tc.name, got.ProposedRuleMatrix, tc.expected)
			}
		})
	}
}

func TestGetAllocationRequests(t *testing.T) {
	fakeClientMap := map[shared.ClientID]baseclient.Client{}
	cases := []struct {
		name       string
		bPresident executive // base
		input      shared.Resources
		expected   map[shared.ClientID]shared.Resources
	}{
		{
			name: "Limited Resources",
			bPresident: executive{
				PresidentID: 3,
				ResourceRequests: map[shared.ClientID]shared.Resources{
					shared.Team1: 5,
					shared.Team2: 10,
					shared.Team3: 15,
					shared.Team4: 20,
					shared.Team5: 25,
					shared.Team6: 30,
				},
				clientPresident: &baseclient.BasePresident{},
				iigoClients:     fakeClientMap,
			},
			input: 100,
			expected: map[shared.ClientID]shared.Resources{
				shared.Team1: calcExpectedVal(5, 105, 100),
				shared.Team2: calcExpectedVal(10, 105, 100),
				shared.Team3: calcExpectedVal(15, 105, 100),
				shared.Team4: calcExpectedVal(20, 105, 100),
				shared.Team5: calcExpectedVal(25, 105, 100),
				shared.Team6: calcExpectedVal(30, 105, 100),
			},
		},
		{
			name: "Excess Resources",
			bPresident: executive{
				PresidentID: 3,
				ResourceRequests: map[shared.ClientID]shared.Resources{
					shared.Team1: 5,
					shared.Team2: 10,
					shared.Team3: 15,
					shared.Team4: 20,
					shared.Team5: 25,
					shared.Team6: 30,
				},
				clientPresident: &baseclient.BasePresident{},
				iigoClients:     fakeClientMap,
			},
			input: 150,
			expected: map[shared.ClientID]shared.Resources{
				shared.Team1: 5,
				shared.Team2: 10,
				shared.Team3: 15,
				shared.Team4: 20,
				shared.Team5: 25,
				shared.Team6: 30},
		},
		{
			name: "Client override limited resourceRequests",
			bPresident: executive{
				PresidentID: 5,
				ResourceRequests: map[shared.ClientID]shared.Resources{
					shared.Team1: 5,
					shared.Team2: 10,
					shared.Team3: 15,
					shared.Team4: 20,
					shared.Team5: 25,
					shared.Team6: 30,
				},
				clientPresident: &baseclient.BasePresident{},
				iigoClients:     fakeClientMap,
			},
			input: 100,
			expected: map[shared.ClientID]shared.Resources{
				shared.Team1: calcExpectedVal(5, 105, 100),
				shared.Team2: calcExpectedVal(10, 105, 100),
				shared.Team3: calcExpectedVal(15, 105, 100),
				shared.Team4: calcExpectedVal(20, 105, 100),
				shared.Team5: calcExpectedVal(25, 105, 100),
				shared.Team6: calcExpectedVal(30, 105, 100),
			},
		},
		{
			name: "Client override excess resourceRequests",
			bPresident: executive{
				PresidentID: 5,
				ResourceRequests: map[shared.ClientID]shared.Resources{
					shared.Team1: 5,
					shared.Team2: 10,
					shared.Team3: 15,
					shared.Team4: 20,
					shared.Team5: 25,
					shared.Team6: 30},
				clientPresident: &baseclient.BasePresident{},
				iigoClients:     fakeClientMap,
			},
			input: 150,
			expected: map[shared.ClientID]shared.Resources{
				shared.Team1: 5,
				shared.Team2: 10,
				shared.Team3: 15,
				shared.Team4: 20,
				shared.Team5: 25,
				shared.Team6: 30,
			},
		},
	}

	wantPresidentReturnType := shared.PresidentAllocation
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			got := tc.bPresident.getAllocationRequests(tc.input)

			if got.ContentType != wantPresidentReturnType {
				t.Errorf("%v - Failed. Expected action type  %v, got action %v", tc.name, wantPresidentReturnType, got.ContentType)
			}

			if !reflect.DeepEqual(got.ResourceMap, tc.expected) {
				t.Errorf("%v - Failed. Got '%v', but expected '%v'", tc.name, got.ResourceMap, tc.expected)
			}
		})
	}
}

func TestReplyAllocationRequest(t *testing.T) {
	fakeClientMap := map[shared.ClientID]baseclient.Client{}
	var logging shared.Logger = func(format string, a ...interface{}) {}
	cases := []struct {
		name           string
		bPresident     executive // base
		commonPool     shared.Resources
		clientRequests map[shared.ClientID]shared.Resources
		expected       map[shared.ClientID]shared.Resources
	}{
		{
			name: "Excess resources",
			bPresident: executive{
				PresidentID:     5,
				clientPresident: &baseclient.BasePresident{},
				gameConf:        &config.IIGOConfig{},
				logger:          logging,
				iigoClients:     fakeClientMap,
			},
			clientRequests: map[shared.ClientID]shared.Resources{
				shared.Team1: 5,
				shared.Team2: 10,
				shared.Team3: 15,
				shared.Team4: 20,
				shared.Team5: 25,
				shared.Team6: 30,
			},
			commonPool: 150,
			expected: map[shared.ClientID]shared.Resources{
				shared.Team1: 5,
				shared.Team2: 10,
				shared.Team3: 15,
				shared.Team4: 20,
				shared.Team5: 25,
				shared.Team6: 30,
			},
		},
		{
			name: "Limited resources",
			bPresident: executive{
				PresidentID:     1,
				clientPresident: &baseclient.BasePresident{},
				gameConf:        &config.IIGOConfig{},
				logger:          logging,
				iigoClients:     fakeClientMap,
			},
			clientRequests: map[shared.ClientID]shared.Resources{
				shared.Team1: 5,
				shared.Team2: 10,
				shared.Team3: 15,
				shared.Team4: 20,
				shared.Team5: 25,
				shared.Team6: 30,
			},
			commonPool: 100,
			expected: map[shared.ClientID]shared.Resources{
				shared.Team1: calcExpectedVal(5, 105, 100),
				shared.Team2: calcExpectedVal(10, 105, 100),
				shared.Team3: calcExpectedVal(15, 105, 100),
				shared.Team4: calcExpectedVal(20, 105, 100),
				shared.Team5: calcExpectedVal(25, 105, 100),
				shared.Team6: calcExpectedVal(30, 105, 100),
			},
		},
		{
			name: "Zero requests",
			bPresident: executive{
				PresidentID:     3,
				clientPresident: &baseclient.BasePresident{},
				gameConf:        &config.IIGOConfig{},
				logger:          logging,
				iigoClients:     fakeClientMap,
			},
			clientRequests: map[shared.ClientID]shared.Resources{
				shared.Team1: 0,
				shared.Team2: 0,
				shared.Team3: 0,
				shared.Team4: 0,
				shared.Team5: 0,
				shared.Team6: 0,
			},
			commonPool: 100,
			expected: map[shared.ClientID]shared.Resources{
				shared.Team1: 0,
				shared.Team2: 0,
				shared.Team3: 0,
				shared.Team4: 0,
				shared.Team5: 0,
				shared.Team6: 0,
			},
		},
	}

	rulesInPlay := map[string]rules.RuleMatrix{}

	err := rules.PullRuleIntoPlayInternal("allocation_decision", generateRuleStore(), rulesInPlay)
	if err != nil {
		t.Errorf("Unexpected Error when pulling rule into play: %v", err)
	}
	err = rules.PullRuleIntoPlayInternal("check_allocation_rule", generateRuleStore(), rulesInPlay)
	if err != nil {
		t.Errorf("Unexpected Error when pulling rule into play: %v", err)
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			fakeClientMap := map[shared.ClientID]baseclient.Client{}
			fakeGameState := gamestate.GameState{
				CommonPool: tc.commonPool,
				IIGORolesBudget: map[shared.Role]shared.Resources{
					shared.President: 10,
					shared.Speaker:   10,
					shared.Judge:     10,
				},
			}
			fakeGameConfig := config.IIGOConfig{
				ReplyAllocationRequestsActionCost: 1,
			}
			fakeServer := fakeServerHandle{
				PresidentID: tc.bPresident.PresidentID,
				RuleInfo: gamestate.RulesContext{
					CurrentRulesInPlay: rulesInPlay,
				},
			}

			aliveID := []shared.ClientID{}

			for clientID := range tc.clientRequests {
				aliveID = append(aliveID, clientID)
				newClient := baseclient.NewClient(clientID)
				newClient.Initialise(fakeServer)
				fakeClientMap[clientID] = newClient
			}
			tc.bPresident.iigoClients = fakeClientMap
			tc.bPresident.syncWithGame(&fakeGameState, &fakeGameConfig)
			tc.bPresident.setAllocationRequest(tc.clientRequests)
			tc.bPresident.replyAllocationRequest(tc.commonPool)

			for clientID, expectedAllocation := range tc.expected {
				communicationsGot := *(fakeClientMap[clientID]).GetCommunications()
				presidentCommunication := communicationsGot[tc.bPresident.PresidentID][0]

				allocation := presidentCommunication[shared.IIGOAllocationDecision].IIGOValueData
				allocationAmount := allocation.Amount
				expectedAllocVar := allocation.Expected
				decidedVar := allocation.DecisionMade

				if presidentCommunication[shared.IIGOAllocationDecision].T != shared.CommunicationIIGOValue {
					t.Errorf("Allocation failed for client %v. Rule type is %v", clientID, presidentCommunication[shared.IIGOAllocationDecision].T)
				}

				if reflect.DeepEqual(expectedAllocVar, rules.VariableValuePair{}) {
					t.Errorf("Allocation failed for client %v. Expected Tax Variable entry is empty. Got communication : %v", clientID, communicationsGot)
				}

				if reflect.DeepEqual(decidedVar, rules.VariableValuePair{}) {
					t.Errorf("Allocation failed for client %v. Decided Tax Variable entry is empty. Got communication : %v", clientID, communicationsGot)
				}

				if allocationAmount != expectedAllocation {
					t.Errorf("Allocation failed for client %v. Expected tax: %v, evaluated tax: %v", clientID, expectedAllocation, allocationAmount)
				}

				allocationRequest := fakeClientMap[clientID].RequestAllocation()

				if allocationRequest != expectedAllocation {
					t.Errorf("Allocation failed for client %v. Expected allocation request <= %v , got allocation request: %v", clientID, expectedAllocation, allocationRequest)
				}

			}
		})
	}
}

func TestPresidentIncurServiceCharge(t *testing.T) {
	cases := []struct {
		name                    string
		bPresident              executive // base
		input                   shared.Resources
		expectedReturn          bool
		expectedCommonPool      shared.Resources
		expectedPresidentBudget shared.Resources
	}{
		{
			name: "Excess pay",
			bPresident: executive{
				PresidentID: shared.Team1,
				gameState: &gamestate.GameState{
					CommonPool: 400,
					IIGORolesBudget: map[shared.Role]shared.Resources{
						shared.President: 100,
						shared.Speaker:   10,
						shared.Judge:     10,
					},
				},
			},
			input:                   50,
			expectedReturn:          true,
			expectedCommonPool:      350,
			expectedPresidentBudget: 50,
		},
		{
			name: "Negative Budget",
			bPresident: executive{
				PresidentID: shared.Team1,
				gameState: &gamestate.GameState{
					CommonPool: 400,
					IIGORolesBudget: map[shared.Role]shared.Resources{
						shared.President: 10,
						shared.Speaker:   10,
						shared.Judge:     10,
					},
				},
			},
			input:                   50,
			expectedReturn:          true,
			expectedCommonPool:      350,
			expectedPresidentBudget: -40,
		},
		{
			name: "Limited common pool",
			bPresident: executive{
				PresidentID: shared.Team1,
				gameState: &gamestate.GameState{
					CommonPool: 40,
					IIGORolesBudget: map[shared.Role]shared.Resources{
						shared.President: 10,
						shared.Speaker:   10,
						shared.Judge:     10,
					},
				},
			},
			input:                   50,
			expectedReturn:          false,
			expectedCommonPool:      40,
			expectedPresidentBudget: 10,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			returned := tc.bPresident.incurServiceCharge(tc.input)
			commonPool := tc.bPresident.gameState.CommonPool
			presidentBudget := tc.bPresident.gameState.IIGORolesBudget[shared.President]
			if returned != tc.expectedReturn ||
				commonPool != tc.expectedCommonPool ||
				presidentBudget != tc.expectedPresidentBudget {
				t.Errorf("%v - Failed. Got '%v, %v, %v', but expected '%v, %v, %v'",
					tc.name, returned, commonPool, presidentBudget,
					tc.expectedReturn, tc.expectedCommonPool, tc.expectedPresidentBudget)
			}
		})
	}
}

func expectedTax(r shared.ResourcesReport) shared.Resources {
	if r.Reported {
		return 0.1 * r.ReportedAmount
	} else {
		return 15
	}
}

func TestBroadcastTaxation(t *testing.T) {
	rulesInPlay := map[string]rules.RuleMatrix{}
	var logging shared.Logger = func(format string, a ...interface{}) {}
	cases := []struct {
		name          string
		bPresident    executive // base
		commonPool    shared.Resources
		clientReports map[shared.ClientID]shared.ResourcesReport
	}{
		{
			name: "Simple test",
			bPresident: executive{
				PresidentID:     shared.Team4,
				clientPresident: &baseclient.BasePresident{},
				logger:          logging,
				gameState: &gamestate.GameState{
					RulesInfo: gamestate.RulesContext{
						CurrentRulesInPlay: rulesInPlay,
					},
				},
			},
			clientReports: map[shared.ClientID]shared.ResourcesReport{
				shared.Team1: {ReportedAmount: 30, Reported: true},
				shared.Team2: {ReportedAmount: 9, Reported: true},
				shared.Team3: {ReportedAmount: 15, Reported: true},
				shared.Team4: {ReportedAmount: 20, Reported: true},
				shared.Team5: {ReportedAmount: 25, Reported: true},
				shared.Team6: {ReportedAmount: 40, Reported: true},
			},
			commonPool: 150,
		},
		{
			name: "Some non-reports test",
			bPresident: executive{
				PresidentID:     shared.Team1,
				clientPresident: &baseclient.BasePresident{},
				logger:          logging,
				gameState: &gamestate.GameState{
					RulesInfo: gamestate.RulesContext{
						CurrentRulesInPlay: rulesInPlay,
					},
				},
			},
			clientReports: map[shared.ClientID]shared.ResourcesReport{
				shared.Team1: {ReportedAmount: 30, Reported: true},
				shared.Team2: {Reported: false},
				shared.Team3: {ReportedAmount: 15, Reported: true},
				shared.Team4: {Reported: false},
				shared.Team5: {ReportedAmount: 25, Reported: true},
				shared.Team6: {Reported: false},
			},
			commonPool: 150,
		},
	}

	err := rules.PullRuleIntoPlayInternal("tax_decision", generateRuleStore(), rulesInPlay)
	if err != nil {
		t.Errorf("Unexpected Error when pulling rule into play: %v", err)
	}
	err = rules.PullRuleIntoPlayInternal("check_taxation_rule", generateRuleStore(), rulesInPlay)
	if err != nil {
		t.Errorf("Unexpected Error when pulling rule into play: %v", err)
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			fakeClientMap := map[shared.ClientID]baseclient.Client{}
			fakeGameState := gamestate.GameState{
				CommonPool: tc.commonPool,
				IIGORolesBudget: map[shared.Role]shared.Resources{
					shared.President: 10,
					shared.Speaker:   10,
					shared.Judge:     10,
				},
			}
			fakeGameConfig := config.IIGOConfig{
				ReplyAllocationRequestsActionCost: 1,
				BroadcastTaxationActionCost:       1,
			}
			fakeServer := fakeServerHandle{
				PresidentID: tc.bPresident.PresidentID,
				RuleInfo: gamestate.RulesContext{
					CurrentRulesInPlay: rulesInPlay,
				},
			}

			aliveID := []shared.ClientID{}

			for clientID := range tc.clientReports {
				aliveID = append(aliveID, clientID)
				newClient := baseclient.NewClient(clientID)
				newClient.Initialise(fakeServer)
				fakeClientMap[clientID] = newClient
			}

			tc.bPresident.iigoClients = fakeClientMap
			tc.bPresident.syncWithGame(&fakeGameState, &fakeGameConfig)
			tc.bPresident.broadcastTaxation(tc.clientReports, aliveID)

			for clientID, resources := range tc.clientReports {
				communicationsGot := *(fakeClientMap[clientID]).GetCommunications()
				presidentCommunication := communicationsGot[tc.bPresident.PresidentID][0]

				taxDecision := presidentCommunication[shared.IIGOTaxDecision].IIGOValueData
				taxAmount := taxDecision.Amount
				taxExpectedVar := taxDecision.Expected
				decided := taxDecision.DecisionMade

				if presidentCommunication[shared.IIGOTaxDecision].T != shared.CommunicationIIGOValue {
					t.Errorf("Taxation failed for client %v. Rule type is %v", clientID, presidentCommunication[shared.IIGOTaxDecision].T)
				}

				if reflect.DeepEqual(taxExpectedVar, rules.VariableValuePair{}) {
					t.Errorf("Taxation failed for client %v. Expected Tax Variable entry is empty. Got communication : %v", clientID, communicationsGot)
				}

				if reflect.DeepEqual(decided, rules.VariableValuePair{}) {
					t.Errorf("Taxation failed for client %v. Decided Tax Variable entry is empty. Got communication : %v", clientID, communicationsGot)
				}

				expectedTax := expectedTax(resources)

				if taxAmount != expectedTax {
					t.Errorf("Taxation failed for client %v. Expected tax: %v, evaluated tax: %v", clientID, expectedTax, taxAmount)
				}

				// Check client return
				paidTax := fakeClientMap[clientID].GetTaxContribution()

				if expectedTax != paidTax {
					t.Errorf("Taxation failed for client %v. expected to pay at least %v, got tax contribution %v", clientID, expectedTax, paidTax)
				}
			}

		})
	}
}

type fakeServerHandle struct {
	PresidentID shared.ClientID
	JudgeID     shared.ClientID
	SpeakerID   shared.ClientID
	RuleInfo    gamestate.RulesContext
}

func (s fakeServerHandle) GetGameState() gamestate.ClientGameState {
	return gamestate.ClientGameState{
		SpeakerID:   s.SpeakerID,
		JudgeID:     s.JudgeID,
		PresidentID: s.PresidentID,
		RulesInfo:   s.RuleInfo,
	}
}

func (s fakeServerHandle) GetGameConfig() config.ClientConfig {
	return config.ClientConfig{}
}
