package iigointernal

import (
	"reflect" // Used to compare two maps
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
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
				ID:              shared.Team1,
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

func checkIfInList(s string, lst []string) bool {
	for _, i := range lst {
		if s == i {
			return true
		}
	}
	return false
}
func TestPickRuleToVote(t *testing.T) {
	cases := []struct {
		name  string
		input []string
		reply string
		want  bool
	}{
		{
			name:  "Basic rule",
			input: []string{"rule"},
			reply: "rule",
			want:  true,
		},
		{
			name:  "Empty string",
			input: []string{""},
			reply: "",
			want:  true,
		},
		{
			name:  "Longer list",
			input: []string{"Somas", "2020", "Internal", "Server", "Roles", "President"},
			reply: "",
			want:  true,
		},
		{
			name:  "Empty list",
			input: []string{},
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
					if got.ProposedRule != "" {
						t.Errorf("%v - Failed. Returned '%v', but expectd an empty string", tc.name, got.ProposedRule)
					}
				} else if !checkIfInList(got.ProposedRule, tc.input) {
					t.Errorf("%v - Failed. Returned '%v', expected '%v'", tc.name, got.ProposedRule, tc.input)
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
		input    []string
		expected []string
	}{
		{
			name:     "Single Rule",
			input:    []string{"rule"},
			expected: []string{"rule"},
		},
		{
			name:     "No Rule",
			input:    []string{""},
			expected: []string{""},
		},
		{
			name:     "Multiple Rules",
			input:    []string{"Somas", "2020", "Internal", "Server", "Roles", "President"},
			expected: []string{"Somas", "2020", "Internal", "Server", "Roles", "President"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			executive := &executive{
				ID:             shared.Team1,
				RulesProposals: []string{},
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
				ID:              3,
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
				ID:              3,
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
				ID:              4,
				clientPresident: &baseclient.BasePresident{},
			},
			expectedLength: 6,
		},
		{
			name:  "Client empty tax map base",
			input: map[shared.ClientID]shared.ResourcesReport{},
			bPresident: executive{
				ID:              5,
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
				ID:              5,
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
				ID:              5,
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
				ID:              5,
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
		IIGORolesBudget: map[string]shared.Resources{
			"president": 10,
			"speaker":   10,
			"judge":     10,
		},
	}
	cases := []struct {
		name       string
		bPresident executive // base
		expected   []string
	}{
		{
			name: "Empty tax map base",
			bPresident: executive{
				ID:              3,
				RulesProposals:  []string{},
				clientPresident: &baseclient.BasePresident{},
				gameState:       &fakeGameState,
			},
			expected: []string{""},
		},
		{
			name: "Short tax map base",
			bPresident: executive{
				ID:              3,
				RulesProposals:  []string{"test"},
				clientPresident: &baseclient.BasePresident{},
				gameState:       &fakeGameState,
			},
			expected: []string{"test"},
		},
		{
			name: "Long tax map base",
			bPresident: executive{
				ID:              3,
				RulesProposals:  []string{"Somas", "2020", "Internal", "Server", "Roles", "President"},
				clientPresident: &baseclient.BasePresident{},
				gameState:       &fakeGameState,
			},
			expected: []string{"Somas", "2020", "Internal", "Server", "Roles", "President"},
		},
		{
			name: "Client empty tax map base",
			bPresident: executive{
				ID:              5,
				RulesProposals:  []string{"Somas", "2020", "Internal", "Server", "Roles", "President"},
				clientPresident: &baseclient.BasePresident{},
				gameState:       &fakeGameState,
			},
			expected: []string{"Somas", "2020", "Internal", "Server", "Roles", "President"},
		},
		{
			name: "Client tax map base override",
			bPresident: executive{
				ID:              5,
				RulesProposals:  []string{"Somas", "2020", "Internal", "Server", "Roles", "President"},
				clientPresident: &baseclient.BasePresident{},
				gameState:       &fakeGameState,
			},
			expected: []string{"Somas", "2020", "Internal", "Server", "Roles", "President"},
		},
	}

	wantPresidentReturnType := shared.PresidentRuleProposal
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			got := tc.bPresident.getRuleForSpeaker()

			if got.ContentType != wantPresidentReturnType {
				t.Errorf("%v - Failed. Expected action type  %v, got action %v", tc.name, wantPresidentReturnType, got.ContentType)
			}

			if len(tc.expected) == 0 {
				if got.ProposedRule != "" {
					t.Errorf("%v - Failed. Returned '%v', but expectd an empty string", tc.name, got.ProposedRule)
				}
			} else if !checkIfInList(got.ProposedRule, tc.expected) {
				t.Errorf("%v - Failed. Returned '%v', expected '%v'", tc.name, got.ProposedRule, tc.expected)
			}
		})
	}
}

func TestGetAllocationRequests(t *testing.T) {
	cases := []struct {
		name       string
		bPresident executive // base
		input      shared.Resources
		expected   map[shared.ClientID]shared.Resources
	}{
		{
			name: "Limited Resources",
			bPresident: executive{
				ID: 3,
				ResourceRequests: map[shared.ClientID]shared.Resources{
					shared.Team1: 5,
					shared.Team2: 10,
					shared.Team3: 15,
					shared.Team4: 20,
					shared.Team5: 25,
					shared.Team6: 30,
				},
				clientPresident: &baseclient.BasePresident{},
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
				ID: 3,
				ResourceRequests: map[shared.ClientID]shared.Resources{
					shared.Team1: 5,
					shared.Team2: 10,
					shared.Team3: 15,
					shared.Team4: 20,
					shared.Team5: 25,
					shared.Team6: 30,
				},
				clientPresident: &baseclient.BasePresident{},
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
				ID: 5,
				ResourceRequests: map[shared.ClientID]shared.Resources{
					shared.Team1: 5,
					shared.Team2: 10,
					shared.Team3: 15,
					shared.Team4: 20,
					shared.Team5: 25,
					shared.Team6: 30,
				},
				clientPresident: &baseclient.BasePresident{},
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
				ID: 5,
				ResourceRequests: map[shared.ClientID]shared.Resources{
					shared.Team1: 5,
					shared.Team2: 10,
					shared.Team3: 15,
					shared.Team4: 20,
					shared.Team5: 25,
					shared.Team6: 30},
				clientPresident: &baseclient.BasePresident{},
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
				ID:              5,
				clientPresident: &baseclient.BasePresident{},
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
				ID:              1,
				clientPresident: &baseclient.BasePresident{},
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
				ID:              3,
				clientPresident: &baseclient.BasePresident{},
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

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			fakeClientMap := map[shared.ClientID]baseclient.Client{}
			fakeGameState := gamestate.GameState{
				CommonPool: tc.commonPool,
				IIGORolesBudget: map[string]shared.Resources{
					"president": 10,
					"speaker":   10,
					"judge":     10,
				},
			}

			aliveID := []shared.ClientID{}

			for clientID := range tc.clientRequests {
				aliveID = append(aliveID, clientID)
				fakeClientMap[clientID] = baseclient.NewClient(clientID)
			}

			setIIGOClients(&fakeClientMap)
			tc.bPresident.setGameState(&fakeGameState)
			tc.bPresident.setAllocationRequest(tc.clientRequests)
			tc.bPresident.replyAllocationRequest(tc.commonPool, aliveID)

			for clientID, expectedAllocation := range tc.expected {
				communicationGot := *(fakeClientMap[clientID]).GetCommunications()
				communicationExpected := map[shared.ClientID][]map[shared.CommunicationFieldName]shared.CommunicationContent{
					tc.bPresident.ID: {
						{shared.AllocationAmount: {T: shared.CommunicationInt, IntegerData: int(expectedAllocation)}},
					},
				}
				if !reflect.DeepEqual(communicationGot, communicationExpected) {
					t.Errorf("Allocation request failed. Expected communication: %v,\n Got communication : %v", communicationExpected, communicationGot)
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
				ID: shared.Team1,
				gameState: &gamestate.GameState{
					CommonPool: 400,
					IIGORolesBudget: map[string]shared.Resources{
						"president": 100,
						"speaker":   10,
						"judge":     10,
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
				ID: shared.Team1,
				gameState: &gamestate.GameState{
					CommonPool: 400,
					IIGORolesBudget: map[string]shared.Resources{
						"president": 10,
						"speaker":   10,
						"judge":     10,
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
				ID: shared.Team1,
				gameState: &gamestate.GameState{
					CommonPool: 40,
					IIGORolesBudget: map[string]shared.Resources{
						"president": 10,
						"speaker":   10,
						"judge":     10,
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
			presidentBudget := tc.bPresident.gameState.IIGORolesBudget["president"]
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
