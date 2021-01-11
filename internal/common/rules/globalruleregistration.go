package rules

import (
	"fmt"

	"gonum.org/v1/gonum/mat"
)

// Init Registers all global scoped rules
func InitialRuleRegistration(startWithRules bool) (AvailableRules map[string]RuleMatrix, RulesInPlay map[string]RuleMatrix) {
	availableRules := make(map[string]RuleMatrix)
	rulesInPlay := make(map[string]RuleMatrix)
	availableRules = registerDemoRule(availableRules)
	availableRules = registerRulesByMass(availableRules)
	if startWithRules {
		rulesInPlay = CopyRulesMap(availableRules)
	}
	return availableRules, rulesInPlay
}

// registerDemoRule Defines and registers demo rule
func registerDemoRule(AvailableRules map[string]RuleMatrix) map[string]RuleMatrix {

	//A very contrived rule//
	name := "Kinda Complicated Rule"
	reqVar := []VariableFieldName{
		NumberOfIslandsContributingToCommonPool,
		NumberOfFailedForages,
		NumberOfBrokenAgreements,
		MaxSeverityOfSanctions,
	}

	v := []float64{1, 0, 0, 0, -4, 0, -1, -1, 0, 2, 0, 0, 0, 1, -2, 0, 0, 1, 0, -1}
	CoreMatrix := mat.NewDense(4, 5, v)
	aux := []float64{1, 1, 2, 0}
	AuxiliaryVector := mat.NewVecDense(4, aux)

	_, ruleErr := RegisterNewRuleInternal(name, reqVar, *CoreMatrix, *AuxiliaryVector, AvailableRules, false, RuleLink{
		Linked: false,
	})
	if ruleErr != nil {
		panic(ruleErr.Error())
	}
	// Check internal/clients/team3/client.go for an implementation of a basic evaluator for this rule
	return AvailableRules
}

// RawRuleSpecification allows a user to use the CompileRuleCase function to build a rule matrix
type RawRuleSpecification struct {
	Name       string
	ReqVar     []VariableFieldName
	Values     []float64
	Aux        []float64
	Mutable    bool
	Linked     bool
	LinkType   LinkTypeOption
	LinkedRule string
}

func registerRulesByMass(availableRules map[string]RuleMatrix) map[string]RuleMatrix {
	ruleSpecs := []RawRuleSpecification{
		{
			//Deprecated
			Name: "inspect_ballot_rule",
			ReqVar: []VariableFieldName{
				NumberOfIslandsAlive,
				NumberOfBallotsCast,
			},
			Values:  []float64{1, -1, 0},
			Aux:     []float64{0},
			Mutable: false,
			Linked:  false,
		},
		{
			Name: "allocations_made_rule",
			ReqVar: []VariableFieldName{
				AllocationMade,
			},
			Values:  []float64{1, -1},
			Aux:     []float64{0},
			Mutable: false,
			Linked:  false,
		},
		{
			Name: "judge_inspection_rule",
			ReqVar: []VariableFieldName{
				JudgeInspectionPerformed,
			},
			Values:  []float64{1, -1},
			Aux:     []float64{0},
			Mutable: false,
			Linked:  false,
		},
		{
			Name: "check_taxation_rule",
			ReqVar: []VariableFieldName{
				IslandTaxContribution,
				ExpectedTaxContribution,
			},
			Values:  []float64{1, -1, 0},
			Aux:     []float64{2},
			Mutable: false,
			Linked:  false,
		},
		{
			Name: "check_allocation_rule",
			ReqVar: []VariableFieldName{
				IslandAllocation,
				ExpectedAllocation,
			},
			Values:  []float64{-1, 1, 0},
			Aux:     []float64{2},
			Mutable: false,
			Linked:  false,
		},
		{
			Name: "vote_called_rule",
			ReqVar: []VariableFieldName{
				RuleSelected,
				VoteCalled,
			},
			Values:  []float64{1, -1, 0},
			Aux:     []float64{0},
			Mutable: false,
			Linked:  false,
		},
		{
			Name: "vote_result_rule",
			ReqVar: []VariableFieldName{
				VoteResultAnnounced,
				VoteCalled,
			},
			Values:  []float64{1, -1, 0},
			Aux:     []float64{0},
			Mutable: false,
			Linked:  false,
		},
		{
			Name: "islands_allowed_to_vote_rule",
			ReqVar: []VariableFieldName{
				AllIslandsAllowedToVote,
			},
			Values:  []float64{1, -1},
			Aux:     []float64{0},
			Mutable: false,
			Linked:  false,
		},
		{
			Name: "rule_to_vote_on_rule",
			ReqVar: []VariableFieldName{

				SpeakerProposedPresidentRule,
			},
			Values:  []float64{1, -1},
			Aux:     []float64{0},
			Mutable: false,
			Linked:  false,
		},
		{
			Name: "iigo_economic_sanction_1",
			ReqVar: []VariableFieldName{
				IslandReportedResources,
				ConstSanctionAmount,
				TurnsLeftOnSanction,
			},
			Values:  []float64{0, 0, 1, 0, 0, 0, 0, 0},
			Aux:     []float64{1, 4},
			Mutable: true,
			Linked:  false,
		},
		{
			Name: "iigo_economic_sanction_2",
			ReqVar: []VariableFieldName{
				IslandReportedResources,
				ConstSanctionAmount,
				TurnsLeftOnSanction,
			},
			Values:  []float64{0, 0, 1, 0, 0.1, 1, 0, 0},
			Aux:     []float64{1, 4},
			Mutable: true,
			Linked:  false,
		},
		{
			Name: "iigo_economic_sanction_3",
			ReqVar: []VariableFieldName{
				IslandReportedResources,
				ConstSanctionAmount,
				TurnsLeftOnSanction,
			},
			Values:  []float64{0, 0, 1, 0, 0.3, 1, 0, 0},
			Aux:     []float64{1, 4},
			Mutable: true,
			Linked:  false,
		},
		{
			Name: "iigo_economic_sanction_4",
			ReqVar: []VariableFieldName{
				IslandReportedResources,
				ConstSanctionAmount,
				TurnsLeftOnSanction,
			},
			Values:  []float64{0, 0, 1, 0, 0.5, 1, 0, 0},
			Aux:     []float64{1, 4},
			Mutable: true,
			Linked:  false,
		},
		{
			Name: "iigo_economic_sanction_5",
			ReqVar: []VariableFieldName{
				IslandReportedResources,
				ConstSanctionAmount,
				TurnsLeftOnSanction,
			},
			Values:  []float64{0, 0, 1, 0, 0.8, 1, 0, 0},
			Aux:     []float64{1, 4},
			Mutable: true,
			Linked:  false,
		},
		{
			Name: "check_sanction_rule",
			ReqVar: []VariableFieldName{
				SanctionPaid,
				SanctionExpected,
			},
			Values:  []float64{1, -1, 0},
			Aux:     []float64{0},
			Mutable: false,
			Linked:  false,
		},
		{
			Name: "rule_chosen_from_proposal_list",
			ReqVar: []VariableFieldName{
				RuleChosenFromProposalList,
			},
			Values:  []float64{-1, 1},
			Aux:     []float64{0},
			Mutable: false,
			Linked:  false,
		},
		{
			Name: "announcement_matches_vote",
			ReqVar: []VariableFieldName{
				AnnouncementRuleMatchesVote,
				AnnouncementResultMatchesVote,
			},
			Values:  []float64{0, -1, 1, -1, 0, 1},
			Aux:     []float64{0, 0},
			Mutable: false,
			Linked:  false,
		},
		{
			Name: "president_over_budget",
			ReqVar: []VariableFieldName{
				PresidentLeftoverBudget,
			},
			Values:  []float64{1, 0},
			Aux:     []float64{2},
			Mutable: false,
			Linked:  false,
		},
		{
			Name: "speaker_over_budget",
			ReqVar: []VariableFieldName{
				SpeakerLeftoverBudget,
			},
			Values:  []float64{1, 0},
			Aux:     []float64{2},
			Mutable: false,
			Linked:  false,
		},
		{
			Name: "judge_over_budget",
			ReqVar: []VariableFieldName{
				JudgeLeftoverBudget,
			},
			Values:  []float64{1, 0},
			Aux:     []float64{2},
			Mutable: false,
			Linked:  false,
		},
		{
			Name: "obl_to_propose_rule_if_some_are_given",
			ReqVar: []VariableFieldName{
				IslandsProposedRules,
				PresidentRuleProposal,
			},
			Values:  []float64{1, -1, 0},
			Aux:     []float64{0},
			Mutable: false,
			Linked:  false,
		},
		{
			Name: "iigo_monitor_rule_permission_1",
			ReqVar: []VariableFieldName{
				MonitorRoleDecideToMonitor,
				MonitorRoleAnnounce,
			},
			Values:  []float64{1, -1, 0},
			Aux:     []float64{0},
			Mutable: true,
			Linked:  false,
		},
		{
			Name: "iigo_monitor_rule_permission_2",
			ReqVar: []VariableFieldName{
				MonitorRoleEvalResult,
				MonitorRoleEvalResultDecide,
			},
			Values:  []float64{1, -1, 0},
			Aux:     []float64{0},
			Mutable: true,
		},
		{
			Name: "island_must_report_private_resource",
			ReqVar: []VariableFieldName{
				HasIslandReportPrivateResources,
			},
			Values:  []float64{1, -1},
			Aux:     []float64{0},
			Mutable: false,
			Linked:  false,
		},
		{
			Name: "island_must_report_actual_private_resource",
			ReqVar: []VariableFieldName{
				IslandActualPrivateResources,
				IslandReportedPrivateResources,
			},
			Values:  []float64{1, -1, 0},
			Aux:     []float64{0},
			Mutable: false,
			Linked:  false,
		},
		{
			Name: "salary_cycle_speaker",
			ReqVar: []VariableFieldName{
				SpeakerPayment,
			},
			Values:  []float64{-1, 10},
			Aux:     []float64{0},
			Mutable: true,
			Linked:  false,
		},
		{
			Name: "salary_cycle_judge",
			ReqVar: []VariableFieldName{
				JudgePayment,
			},
			Values:  []float64{-1, 10},
			Aux:     []float64{0},
			Mutable: true,
			Linked:  false,
		},
		{
			Name: "salary_cycle_president",
			ReqVar: []VariableFieldName{
				PresidentPayment,
			},
			Values:  []float64{-1, 10},
			Aux:     []float64{0},
			Mutable: true,
			Linked:  false,
		},
		{
			Name: "salary_paid_speaker",
			ReqVar: []VariableFieldName{
				SpeakerPaid,
			},
			Values:  []float64{1, -1},
			Aux:     []float64{0},
			Mutable: false,
			Linked:  false,
		},
		{
			Name: "salary_paid_judge",
			ReqVar: []VariableFieldName{
				JudgePaid,
			},
			Values:  []float64{1, -1},
			Aux:     []float64{0},
			Mutable: false,
			Linked:  false,
		},
		{
			Name: "salary_paid_president",
			ReqVar: []VariableFieldName{
				PresidentPaid,
			},
			Values:  []float64{1, -1},
			Aux:     []float64{0},
			Mutable: false,
			Linked:  false,
		},
		{
			Name: "judge_historical_retribution_permission",
			ReqVar: []VariableFieldName{
				JudgeHistoricalRetributionPerformed,
			},
			Values:  []float64{1, 0},
			Aux:     []float64{0},
			Mutable: false,
			Linked:  false,
		},
		{
			Name: "roles_must_hold_election",
			ReqVar: []VariableFieldName{
				TermEnded,
				ElectionHeld,
			},
			Values:  []float64{1, -1, 0},
			Aux:     []float64{0},
			Mutable: false,
			Linked:  false,
		},
		{
			Name: "must_appoint_elected_island",
			ReqVar: []VariableFieldName{
				AppointmentMatchesVote,
			},
			Values:  []float64{1, -1},
			Aux:     []float64{0},
			Mutable: false,
			Linked:  false,
		},
		{
			Name: "increment_budget_speaker",
			ReqVar: []VariableFieldName{
				SpeakerBudgetIncrement,
			},
			Values:  []float64{-1, 100},
			Aux:     []float64{0},
			Mutable: true,
			Linked:  false,
		},
		{
			Name: "increment_budget_judge",
			ReqVar: []VariableFieldName{
				JudgeBudgetIncrement,
			},
			Values:  []float64{-1, 100},
			Aux:     []float64{0},
			Mutable: true,
			Linked:  false,
		},
		{
			Name: "increment_budget_president",
			ReqVar: []VariableFieldName{
				PresidentBudgetIncrement,
			},
			Values:  []float64{-1, 100},
			Aux:     []float64{0},
			Mutable: true,
			Linked:  false,
		},
		{
			Name: "tax_decision",
			ReqVar: []VariableFieldName{
				TaxDecisionMade,
			},
			Values:     []float64{1, -1},
			Aux:        []float64{0},
			Mutable:    false,
			Linked:     true,
			LinkType:   ParentFailAutoRulePass,
			LinkedRule: "check_taxation_rule",
		},
		{
			Name: "allocation_decision",
			ReqVar: []VariableFieldName{
				AllocationMade,
			},
			Values:     []float64{1, -1},
			Aux:        []float64{0},
			Mutable:    false,
			Linked:     true,
			LinkType:   ParentFailAutoRulePass,
			LinkedRule: "check_allocation_rule",
		},
	}

	for _, rs := range ruleSpecs {
		rowLength := len(rs.ReqVar) + 1
		if len(rs.Values)%rowLength != 0 {
			panic(fmt.Sprintf("Rule '%v' was registered without correct matrix dimensions", rs.Name))
		}
		nrows := len(rs.Values) / rowLength
		CoreMatrix := mat.NewDense(nrows, rowLength, rs.Values)
		AuxiliaryVector := mat.NewVecDense(nrows, rs.Aux)
		var ruleLink RuleLink
		if !rs.Linked {
			ruleLink = RuleLink{
				Linked: false,
			}
		} else {
			ruleLink = RuleLink{
				Linked:     rs.Linked,
				LinkType:   rs.LinkType,
				LinkedRule: rs.LinkedRule,
			}
		}
		_, ruleError := RegisterNewRuleInternal(rs.Name, rs.ReqVar, *CoreMatrix, *AuxiliaryVector, availableRules, rs.Mutable, ruleLink)
		if ruleError != nil {
			panic(ruleError.Error())
		}
	}
	return availableRules
}

// CompileRuleCase allows an agent to quickly build a RuleMatrix using the RawRuleSpecification
func CompileRuleCase(spec RawRuleSpecification) (RuleMatrix, bool) {
	rowLength := len(spec.ReqVar) + 1
	if len(spec.Values)%rowLength != 0 {
		return RuleMatrix{}, false
	}
	nrows := len(spec.Values) / rowLength
	CoreMatrix := mat.NewDense(nrows, rowLength, spec.Values)
	AuxiliaryVector := mat.NewVecDense(nrows, spec.Aux)
	var ruleLink RuleLink
	if !spec.Linked {
		ruleLink = RuleLink{
			Linked: false,
		}
	} else {
		ruleLink = RuleLink{
			Linked:     spec.Linked,
			LinkType:   spec.LinkType,
			LinkedRule: spec.LinkedRule,
		}
	}
	finalRuleMatrix := RuleMatrix{
		RuleName:          spec.Name,
		RequiredVariables: spec.ReqVar,
		ApplicableMatrix:  *CoreMatrix,
		AuxiliaryVector:   *AuxiliaryVector,
		Mutable:           spec.Mutable,
		Link:              ruleLink,
	}
	return finalRuleMatrix, true
}
