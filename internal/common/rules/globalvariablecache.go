package rules

import (
	"fmt"

	"github.com/SOMAS2020/SOMAS2020/pkg/miscutils"
	"github.com/pkg/errors"
)

type VariableValuePair struct {
	VariableName VariableFieldName
	Values       []float64
}

// RegisterNewVariableInternal provides primal register logic for any variable cache
func RegisterNewVariableInternal(pair VariableValuePair, variableStore map[VariableFieldName]VariableValuePair) error {
	if _, ok := variableStore[pair.VariableName]; ok {
		return errors.Errorf("attempted to re-register a variable that had already been registered")
	}
	variableStore[pair.VariableName] = pair
	return nil
}

// UpdateVariableInternal provides primal update logic for any variable cache
func UpdateVariableInternal(variableName VariableFieldName, newValue VariableValuePair, variableStore map[VariableFieldName]VariableValuePair) bool {
	if _, ok := variableStore[variableName]; ok {
		variableStore[variableName] = newValue
		return true
	}
	return false
}

// CopyVariableMap easily copies variable cache
func CopyVariableMap(varMap map[VariableFieldName]VariableValuePair) map[VariableFieldName]VariableValuePair {
	newMap := make(map[VariableFieldName]VariableValuePair)
	for key, value := range varMap {
		newMap[key] = value
	}
	return newMap
}

type VariableFieldName int

const (
	NumberOfIslandsContributingToCommonPool VariableFieldName = iota
	NumberOfFailedForages
	NumberOfBrokenAgreements
	MaxSeverityOfSanctions
	NumberOfIslandsAlive
	NumberOfBallotsCast
	NumberOfAllocationsSent
	AllocationRequestsMade
	AllocationMade
	IslandsAlive
	SpeakerSalary
	SpeakerPayment
	SpeakerPaid
	SpeakerBudgetIncrement
	JudgeSalary
	JudgePayment
	JudgePaid
	JudgeBudgetIncrement
	PresidentSalary
	PresidentPayment
	PresidentPaid
	PresidentBudgetIncrement
	RuleSelected
	VoteCalled
	ExpectedTaxContribution
	ExpectedAllocation
	IslandTaxContribution
	IslandAllocation
	IslandReportedResources
	ConstSanctionAmount
	TurnsLeftOnSanction
	SanctionPaid
	SanctionExpected
	TestVariable
	JudgeInspectionPerformed
	TaxDecisionMade
	MonitorRoleAnnounce
	MonitorRoleDecideToMonitor
	MonitorRoleEvalResult
	MonitorRoleEvalResultDecide
	VoteResultAnnounced
	AllIslandsAllowedToVote
	SpeakerProposedPresidentRule
	PresidentRuleProposal
	RuleChosenFromProposalList
	AnnouncementRuleMatchesVote
	AnnouncementResultMatchesVote
	PresidentLeftoverBudget
	SpeakerLeftoverBudget
	JudgeLeftoverBudget
	IslandsProposedRules
	HasIslandReportPrivateResources
	IslandActualPrivateResources
	IslandReportedPrivateResources
	JudgeHistoricalRetributionPerformed
	TermEnded
	ElectionHeld
	AppointmentMatchesVote
)

func (v VariableFieldName) String() string {
	strs := [...]string{
		"NumberOfIslandsContributingToCommonPool",
		"NumberOfFailedForages",
		"NumberOfBrokenAgreements",
		"MaxSeverityOfSanctions",
		"NumberOfIslandsAlive",
		"NumberOfBallotsCast",
		"NumberOfAllocationsSent",
		"AllocationRequestsMade",
		"AllocationMade",
		"IslandsAlive",
		"SpeakerSalary",
		"SpeakerPayment",
		"SpeakerPaid",
		"SpeakerBudgetIncrement",
		"JudgeSalary",
		"JudgePayment",
		"JudgePaid",
		"JudgeBudgetIncrement",
		"PresidentSalary",
		"PresidentPayment",
		"PresidentPaid",
		"PresidentBudgetIncrement",
		"RuleSelected",
		"VoteCalled",
		"ExpectedTaxContribution",
		"ExpectedAllocation",
		"IslandTaxContribution",
		"IslandAllocation",
		"IslandReportedResources",
		"ConstSanctionAmount",
		"TurnsLeftOnSanction",
		"SanctionPaid",
		"SanctionExpected",
		"TestVariable",
		"JudgeInspectionPerformed",
		"TaxDecisionMade",
		"MonitorRoleAnnounce",
		"MonitorRoleDecideToMonitor",
		"MonitorRoleEvalResult",
		"MonitorRoleEvalResultDecide",
		"VoteResultAnnounced",
		"AllIslandsAllowedToVote",
		"SpeakerProposedPresidentRule",
		"PresidentRuleProposal",
		"RuleChosenFromProposalList",
		"AnnouncementRuleMatchesVote",
		"AnnouncementResultMatchesVote",
		"PresidentLeftoverBudget",
		"SpeakerLeftoverBudget",
		"JudgeLeftoverBudget",
		"IslandsProposedRules",
		"HasIslandReportPrivateResources",
		"IslandActualPrivateResources",
		"IslandReportedPrivateResources",
		"JudgeHistoricalRetributionPerformed",
		"TermEnded",
		"ElectionHeld",
		"AppointmentMatchesVote",
	}
	if v >= 0 && int(v) < len(strs) {
		return strs[v]
	}
	return fmt.Sprintf("UNKNOWN VariableFieldName '%v'", int(v))
}

// GoString implements GoStringer
func (v VariableFieldName) GoString() string {
	return v.String()
}

// MarshalText implements TextMarshaler
func (v VariableFieldName) MarshalText() ([]byte, error) {
	return miscutils.MarshalTextForString(v.String())
}

// MarshalJSON implements RawMessage
func (v VariableFieldName) MarshalJSON() ([]byte, error) {
	return miscutils.MarshalJSONForString(v.String())
}
