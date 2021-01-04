package rules

import "fmt"

// init Registers all variables defined in Static variables list
func init() {
	for _, v := range StaticVariables {
		e := RegisterNewVariable(v)
		if e != nil {
			panic(fmt.Sprintf("variable registration gone wrong, variable: '%v' has been registered multiple times", v.VariableName))
		}
	}
}

// StaticVariables holds all globally defined variables
var StaticVariables = []VariableValuePair{
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
		VariableName: JudgeSalary,
		Values:       []float64{50},
	},
	{
		VariableName: PresidentSalary,
		Values:       []float64{50},
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
		VariableName: IslandsAllowedToVote,
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
}
