package rules

import (
	"fmt"

	"github.com/SOMAS2020/SOMAS2020/pkg/miscutils"
)

// RuleErrorType is a non-critical issue which can be caused by an island trying to modify, register or pick rules which isn't mechanically feasible
type RuleErrorType int

const (
	RuleNotInAvailableRulesCache RuleErrorType = iota
	ModifiedRuleMatrixDimensionMismatch
	AuxVectorDimensionDontMatchRuleMatrix
	RuleRequestedForModificationWasImmutable
	TriedToReRegisterRule
	RuleIsAlreadyInPlay
	RuleIsNotInPlay
)

func (r RuleErrorType) String() string {
	strs := [...]string{
		"RuleNotInAvailableRulesCache",
		"ModifiedRuleMatrixDimensionMismatch",
		"AuxVectorDimensionDontMatchRuleMatrix",
		"RuleRequestedForModificationWasImmutable",
		"TriedToReRegisterRule",
		"RuleIsAlreadyInPlay",
		"RuleIsNotInPlay",
	}

	if r >= 0 && int(r) < len(strs) {
		return strs[r]
	}
	return fmt.Sprintf("UNKNOWN RuleErrorType '%v'", int(r))
}

// GoString implements GoStringer
func (r RuleErrorType) GoString() string {
	return r.String()
}

// MarshalText implements TextMarshaler
func (r RuleErrorType) MarshalText() ([]byte, error) {
	return miscutils.MarshalTextForString(r.String())
}

// MarshalJSON implements RawMessage
func (r RuleErrorType) MarshalJSON() ([]byte, error) {
	return miscutils.MarshalJSONForString(r.String())
}

// RuleError provides a packaged version of the RuleErrorType for clients to deal with

type RuleError struct {
	errorType RuleErrorType
	err       error
}

func (e *RuleError) Error() string {
	return e.err.Error()
}

func (e *RuleError) Type() RuleErrorType {
	return e.errorType
}
