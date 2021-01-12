package adv

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"gonum.org/v1/gonum/mat"
)

type Malice struct {
	Id             shared.ClientID
	rulesToPropose map[string]rules.RuleMatrix
}

func (m *Malice) Initialise(id shared.ClientID) {
	m.Id = id
	m.rulesToPropose = map[string]rules.RuleMatrix{
		"iigo_economic_sanction_2": {
			RuleName: "iigo_economic_sanction_2",
			RequiredVariables: []rules.VariableFieldName{
				rules.IslandReportedResources,
				rules.ConstSanctionAmount,
				rules.TurnsLeftOnSanction,
			},
			ApplicableMatrix: *mat.NewDense(2, 4, []float64{0, 0, 1, 0, 0, 0, 0, 0}),
			AuxiliaryVector:  *mat.NewVecDense(2, []float64{1, 4}),
			Mutable:          false,
			Link:             rules.RuleLink{},
		},
		"iigo_economic_sanction_3": {
			RuleName: "iigo_economic_sanction_3",
			RequiredVariables: []rules.VariableFieldName{
				rules.IslandReportedResources,
				rules.ConstSanctionAmount,
				rules.TurnsLeftOnSanction,
			},
			ApplicableMatrix: *mat.NewDense(2, 4, []float64{0, 0, 1, 0, 0, 0, 0, 0}),
			AuxiliaryVector:  *mat.NewVecDense(2, []float64{1, 4}),
			Mutable:          false,
			Link:             rules.RuleLink{},
		},
		"iigo_economic_sanction_4": {
			RuleName: "iigo_economic_sanction_4",
			RequiredVariables: []rules.VariableFieldName{
				rules.IslandReportedResources,
				rules.ConstSanctionAmount,
				rules.TurnsLeftOnSanction,
			},
			ApplicableMatrix: *mat.NewDense(2, 4, []float64{0, 0, 1, 0, 0, 0, 0, 0}),
			AuxiliaryVector:  *mat.NewVecDense(2, []float64{1, 4}),
			Mutable:          false,
			Link:             rules.RuleLink{},
		},
		"iigo_economic_sanction_5": {
			RuleName: "iigo_economic_sanction_5",
			RequiredVariables: []rules.VariableFieldName{
				rules.IslandReportedResources,
				rules.ConstSanctionAmount,
				rules.TurnsLeftOnSanction,
			},
			ApplicableMatrix: *mat.NewDense(2, 4, []float64{0, 0, 1, 0, 0, 0, 0, 0}),
			AuxiliaryVector:  *mat.NewVecDense(2, []float64{1, 4}),
			Mutable:          false,
			Link:             rules.RuleLink{},
		},
		"iigo_monitor_rule_permission_1": {
			RuleName: "iigo_monitor_rule_permission_1",
			RequiredVariables: []rules.VariableFieldName{
				rules.MonitorRoleDecideToMonitor,
				rules.MonitorRoleAnnounce,
			},
			ApplicableMatrix: *mat.NewDense(1, 3, []float64{1, 0, 0}),
			AuxiliaryVector:  *mat.NewVecDense(1, []float64{0}),
			Mutable:          false,
			Link:             rules.RuleLink{},
		},
		"iigo_monitor_rule_permission_2": {
			RuleName: "iigo_monitor_rule_permission_2",
			RequiredVariables: []rules.VariableFieldName{
				rules.MonitorRoleEvalResult,
				rules.MonitorRoleEvalResultDecide,
			},
			ApplicableMatrix: *mat.NewDense(1, 3, []float64{1, 0, 0}),
			AuxiliaryVector:  *mat.NewVecDense(1, []float64{0}),
			Mutable:          false,
			Link:             rules.RuleLink{},
		},
	}
}

func (m *Malice) SetTaxationAmount(islandsResources map[shared.ClientID]shared.ResourcesReport) (shared.PresidentReturnContent, bool) {
	generatedMap := map[shared.ClientID]shared.Resources{}
	for id, val := range islandsResources {
		if id != m.Id {
			if val.Reported {
				generatedMap[id] = val.ReportedAmount
			} else {
				generatedMap[id] = 500
			}
		} else {
			generatedMap[id] = 0
		}
	}
	return shared.PresidentReturnContent{
		ContentType: shared.PresidentTaxation,
		ResourceMap: generatedMap,
		ActionTaken: true,
	}, true
}

func (m *Malice) ProposeRule(availableRules map[string]rules.RuleMatrix) (rules.RuleMatrix, bool) {
	for key, val := range m.rulesToPropose {
		value := copySingleRuleMatrix(val)
		delete(m.rulesToPropose, key)
		return value, true
	}
	return rules.RuleMatrix{
		RuleName: "check_allocation_rule",
		RequiredVariables: []rules.VariableFieldName{
			rules.IslandAllocation,
			rules.ExpectedAllocation,
		},
		ApplicableMatrix: *mat.NewDense(1, 3, []float64{1, -1, 0}),
		AuxiliaryVector:  *mat.NewVecDense(1, []float64{0}),
		Mutable:          false,
		Link:             rules.RuleLink{},
	}, false
}

func (m *Malice) GetRuleViolationSeverity() (map[string]shared.IIGOSanctionsScore, bool) {
	return map[string]shared.IIGOSanctionsScore{
		"check_taxation_rule": shared.IIGOSanctionsScore(50),
	}, true
}

func (m *Malice) GetPardonedIslands(currentSanctions map[int][]shared.Sanction) (map[int][]bool, bool) {
	pardons := make(map[int][]bool)
	for key, sanctionList := range currentSanctions {
		lst := make([]bool, len(sanctionList))
		pardons[key] = lst
		for index, sanction := range sanctionList {
			if sanction.ClientID == m.Id {
				pardons[key][index] = true
			} else {
				pardons[key][index] = false
			}
		}
	}
	return pardons, true
}

func (m *Malice) CallPresidentElection(monitoring shared.MonitorResult, turnsInPower int, allIslands []shared.ClientID) (shared.ElectionSettings, bool) {
	return shared.ElectionSettings{
		VotingMethod:  shared.InstantRunoff,
		IslandsToVote: allIslands,
		HoldElection:  true,
	}, true
}

func (m *Malice) DecideNextPresident(winner shared.ClientID) (shared.ClientID, bool) {
	return m.Id, true
}

func (m *Malice) CallJudgeElection(monitoring shared.MonitorResult, turnsInPower int, allIslands []shared.ClientID) (shared.ElectionSettings, bool) {
	return shared.ElectionSettings{
		VotingMethod:  shared.InstantRunoff,
		IslandsToVote: allIslands,
		HoldElection:  true,
	}, true
}

func (m *Malice) DecideNextJudge(winner shared.ClientID) (shared.ClientID, bool) {
	return m.Id, true
}

func (m *Malice) CallSpeakerElection(monitoring shared.MonitorResult, turnsInPower int, allIslands []shared.ClientID) (shared.ElectionSettings, bool) {
	return shared.ElectionSettings{
		VotingMethod:  shared.InstantRunoff,
		IslandsToVote: allIslands,
		HoldElection:  true,
	}, true
}

func (m *Malice) DecideNextSpeaker(winner shared.ClientID) (shared.ClientID, bool) {
	return m.Id, true
}

func (m *Malice) InspectHistory(iigoHistory []shared.Accountability, turnsAgo int) (map[shared.ClientID]shared.EvaluationReturn, bool, bool) {
	return map[shared.ClientID]shared.EvaluationReturn{}, false, false
}
