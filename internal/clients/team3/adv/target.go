package adv

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"gonum.org/v1/gonum/mat"
)

type Target struct {
	id             shared.ClientID
	TargetID       shared.ClientID
	rulesToPropose map[string]rules.RuleMatrix
}

func (m *Target) Initialise(id shared.ClientID) {
	m.id = id
}

func (m *Target) ProposeRule(availableRules map[string]rules.RuleMatrix) (rules.RuleMatrix, bool) {
	return rules.RuleMatrix{}, false
}

func (m *Target) GetRuleViolationSeverity() (map[string]shared.IIGOSanctionsScore, bool) {
	return map[string]shared.IIGOSanctionsScore{
		"check_taxation_rule": shared.IIGOSanctionsScore(50),
	}, true
}

func (m *Target) SetTaxationAmount(islandsResources map[shared.ClientID]shared.ResourcesReport) (shared.PresidentReturnContent, bool) {
	generatedMap := map[shared.ClientID]shared.Resources{}
	for id, val := range islandsResources {
		if id == m.TargetID {
			if val.Reported {
				generatedMap[id] = val.ReportedAmount
			} else {
				generatedMap[id] = 500
			}
		} else {
			generatedMap[id] = 20
		}
	}
	return shared.PresidentReturnContent{
		ContentType: shared.PresidentTaxation,
		ResourceMap: generatedMap,
		ActionTaken: true,
	}, true
}

func (m *Target) InspectHistory(iigoHistory []shared.Accountability, turnsAgo int) (map[shared.ClientID]shared.EvaluationReturn, bool, bool) {
	returnMap := make(map[shared.ClientID]shared.EvaluationReturn)
	for _, entry := range iigoHistory {
		clientID := entry.ClientID
		if _, ok := returnMap[clientID]; !ok {
			returnMap[clientID] = shared.EvaluationReturn{
				Rules:       []rules.RuleMatrix{},
				Evaluations: []bool{},
			}
		}
		curr := returnMap[clientID]
		if clientID == m.TargetID {
			curr.Rules = append(curr.Rules, fetchCheckTaxationRule())
			curr.Evaluations = append(curr.Evaluations, false)
		}
		returnMap[clientID] = curr
	}
	return returnMap, true, true
}

func (m *Target) GetPardonedIslands(currentSanctions map[int][]shared.Sanction) (map[int][]bool, bool) {
	pardons := make(map[int][]bool)
	for key, sanctionList := range currentSanctions {
		lst := make([]bool, len(sanctionList))
		pardons[key] = lst
		for index, sanction := range sanctionList {
			if sanction.ClientID == m.TargetID {
				pardons[key][index] = false
			} else {
				pardons[key][index] = true
			}
		}
	}
	return pardons, true
}

func (m *Target) CallPresidentElection(monitoring shared.MonitorResult, turnsInPower int, allIslands []shared.ClientID) (shared.ElectionSettings, bool) {
	return shared.ElectionSettings{
		VotingMethod:  shared.InstantRunoff,
		IslandsToVote: allIslands,
		HoldElection:  true,
	}, false
}

func (m *Target) DecideNextPresident(winner shared.ClientID) (shared.ClientID, bool) {
	return m.id, true
}

func (m *Target) CallJudgeElection(monitoring shared.MonitorResult, turnsInPower int, allIslands []shared.ClientID) (shared.ElectionSettings, bool) {
	return shared.ElectionSettings{
		VotingMethod:  shared.InstantRunoff,
		IslandsToVote: allIslands,
		HoldElection:  true,
	}, false
}

func (m *Target) DecideNextJudge(winner shared.ClientID) (shared.ClientID, bool) {
	return m.id, true
}

func (m *Target) CallSpeakerElection(monitoring shared.MonitorResult, turnsInPower int, allIslands []shared.ClientID) (shared.ElectionSettings, bool) {
	return shared.ElectionSettings{
		VotingMethod:  shared.InstantRunoff,
		IslandsToVote: allIslands,
		HoldElection:  true,
	}, false
}

func (m *Target) DecideNextSpeaker(winner shared.ClientID) (shared.ClientID, bool) {
	return m.id, true
}

func fetchCheckTaxationRule() rules.RuleMatrix {
	return rules.RuleMatrix{
		RuleName: "check_taxation_rule",
		RequiredVariables: []rules.VariableFieldName{
			rules.IslandTaxContribution,
			rules.ExpectedTaxContribution,
		},
		ApplicableMatrix: *mat.NewDense(1, 3, []float64{1, -1, 0}),
		AuxiliaryVector:  *mat.NewVecDense(1, []float64{2}),
		Mutable:          false,
		Link:             rules.RuleLink{},
	}
}
