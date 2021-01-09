package rules

import (
	"fmt"
)

// init Registers all variables defined in Static variables list
func InitialVarRegistration() map[VariableFieldName]VariableValuePair {
	baseCache := make(map[VariableFieldName]VariableValuePair)
	for _, v := range StaticVariables {
		e := RegisterNewVariableInternal(v, baseCache)
		if e != nil {
			panic(fmt.Sprintf("variable registration gone wrong, variable: '%v' has been registered multiple times", v.VariableName))
		}
	}
	return baseCache
}

// StaticVariables holds all globally defined variables
var StaticVariables = [...]VariableValuePair{
	{
		VariableName: NumberOfIslandsContributingToCommonPool,
		Values:       []float64{5},
	},
	{
		VariableName: NumberOfFailedForages,
		Values:       []float64{0.5},
	},
	{
		VariableName: NumberOfBrokenAgreements,
		Values:       []float64{1},
	},
	{
		VariableName: MaxSeverityOfSanctions,
		Values:       []float64{2},
	},
	{
		VariableName: NumberOfIslandsAlive,
		Values:       []float64{6},
	},
	{
		VariableName: NumberOfBallotsCast,
		Values:       []float64{6},
	},
	{
		VariableName: NumberOfAllocationsSent,
		Values:       []float64{6},
	},
	{
		VariableName: IslandsAlive,
		Values:       []float64{0, 1, 2, 3, 4, 5},
	},
	{
		VariableName: SpeakerSalary,
		Values:       []float64{50},
	},
	{
		VariableName: SpeakerPayment,
		Values:       []float64{50},
	},
	{
		VariableName: SpeakerPaid,
		Values:       []float64{0},
	},
	{
		VariableName: SpeakerBudgetIncrement,
		Values:       []float64{100},
	},
	{
		VariableName: JudgeSalary,
		Values:       []float64{50},
	},
	{
		VariableName: JudgePayment,
		Values:       []float64{50},
	},
	{
		VariableName: JudgePaid,
		Values:       []float64{0},
	},
	{
		VariableName: JudgeBudgetIncrement,
		Values:       []float64{100},
	},
	{
		VariableName: PresidentSalary,
		Values:       []float64{50},
	},
	{
		VariableName: PresidentPayment,
		Values:       []float64{50},
	},
	{
		VariableName: PresidentPaid,
		Values:       []float64{0},
	},
	{
		VariableName: PresidentBudgetIncrement,
		Values:       []float64{100},
	},
	{
		VariableName: ExpectedTaxContribution,
		Values:       []float64{0},
	},
	{
		VariableName: ExpectedAllocation,
		Values:       []float64{0},
	},
	{
		VariableName: IslandTaxContribution,
		Values:       []float64{0},
	},
	{
		VariableName: IslandAllocation,
		Values:       []float64{0},
	},
	{
		VariableName: RuleSelected,
		Values:       []float64{0},
	},
	{
		VariableName: VoteCalled,
		Values:       []float64{0},
	},
	{
		VariableName: IslandReportedResources,
		Values:       []float64{0},
	},
	{
		VariableName: ConstSanctionAmount,
		Values:       []float64{0},
	},
	{
		VariableName: TurnsLeftOnSanction,
		Values:       []float64{0},
	},
	{
		VariableName: SanctionPaid,
		Values:       []float64{0},
	},
	{
		VariableName: SanctionExpected,
		Values:       []float64{0},
	},
	{
		VariableName: AllocationRequestsMade,
		Values:       []float64{1},
	},
	{
		VariableName: AllocationMade,
		Values:       []float64{0},
	},
	{
		VariableName: JudgeInspectionPerformed,
		Values:       []float64{0},
	},
	{
		VariableName: MonitorRoleAnnounce,
		Values:       []float64{0},
	},
	{
		VariableName: MonitorRoleDecideToMonitor,
		Values:       []float64{0},
	},
	{
		VariableName: MonitorRoleEvalResult,
		Values:       []float64{0},
	},
	{
		VariableName: MonitorRoleEvalResultDecide,
		Values:       []float64{0},
	},
	{
		VariableName: VoteResultAnnounced,
		Values:       []float64{0},
	},
	{
		VariableName: AllIslandsAllowedToVote,
		Values:       []float64{0},
	},
	{
		VariableName: SpeakerProposedPresidentRule,
		Values:       []float64{0},
	},
	{
		VariableName: PresidentRuleProposal,
		Values:       []float64{0},
	},
	{
		VariableName: RuleChosenFromProposalList,
		Values:       []float64{0},
	},
	{
		VariableName: AnnouncementRuleMatchesVote,
		Values:       []float64{0},
	},
	{
		VariableName: AnnouncementResultMatchesVote,
		Values:       []float64{0},
	},
	{
		VariableName: PresidentLeftoverBudget,
		Values:       []float64{0},
	},
	{
		VariableName: SpeakerLeftoverBudget,
		Values:       []float64{0},
	},
	{
		VariableName: JudgeLeftoverBudget,
		Values:       []float64{0},
	},
	{
		VariableName: IslandsProposedRules,
		Values:       []float64{0},
	},
	{
		VariableName: HasIslandReportPrivateResources,
		Values:       []float64{0},
	},
	{
		VariableName: IslandActualPrivateResources,
		Values:       []float64{0},
	},
	{
		VariableName: IslandReportedPrivateResources,
		Values:       []float64{0},
	},
	{
		VariableName: JudgeHistoricalRetributionPerformed,
		Values:       []float64{0},
	},
	{
		VariableName: TermEnded,
		Values:       []float64{0},
	},
	{
		VariableName: ElectionHeld,
		Values:       []float64{0},
	},
	{
		VariableName: AppointmentMatchesVote,
		Values:       []float64{0},
	},
	{
		VariableName: TaxDecisionMade,
		Values:       []float64{1},
	},
}
