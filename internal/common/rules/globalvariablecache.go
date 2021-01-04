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

var VariableMap = map[VariableFieldName]VariableValuePair{}

// RegisterNewVariable Registers the provided variable in the global variable cache
func RegisterNewVariable(pair VariableValuePair) error {
	return registerNewVariableInternal(pair, VariableMap)
}

// registerNewVariableInternal provides primal register logic for any variable cache
func registerNewVariableInternal(pair VariableValuePair, variableStore map[VariableFieldName]VariableValuePair) error {
	if _, ok := variableStore[pair.VariableName]; ok {
		return errors.Errorf("attempted to re-register a variable that had already been registered")
	}
	variableStore[pair.VariableName] = pair
	return nil
}

// UpdateVariable Updates variable in global cache with new value
func UpdateVariable(variableName VariableFieldName, newValue VariableValuePair) bool {
	return updateVariableInternal(variableName, newValue, VariableMap)
}

// updateVariableInternal provides primal update logic for any variable cache
func updateVariableInternal(variableName VariableFieldName, newValue VariableValuePair, variableStore map[VariableFieldName]VariableValuePair) bool {
	if _, ok := variableStore[variableName]; ok {
		variableStore[variableName] = newValue
		return true
	}
	return false
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
	JudgeSalary
	PresidentSalary
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
	VoteResultAnnounced
	IslandsAllowedToVote
	SpeakerProposedPresidentRule
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
		"JudgeSalary",
		"PresidentSalary",
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
		"VoteResultAnnounced",
		"IslandsAllowedToVote",
		"SpeakerProposedPresidentRule",
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
