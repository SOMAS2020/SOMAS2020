package rules

import (
	"gonum.org/v1/gonum/mat"
)

const SingleValueVariableEntry = 0

// LinkTypeOption gives an enumerated type for the various link types available for rules
type LinkTypeOption int

const (
	// ParentFailAutoRulePass allows for NOT(Parent Passes) || Parent and Child pass
	// Useful for cases where if a condition isn't met we don't want to evaluate a rule
	ParentFailAutoRulePass LinkTypeOption = iota
	NoLink
)

// RuleLink provides a containerised package for all linked rules
type RuleLink struct {
	Linked     bool
	LinkType   LinkTypeOption
	LinkedRule string
}

// RuleMatrix provides a container for our matrix based rules
type RuleMatrix struct {
	RuleName          string
	RequiredVariables []VariableFieldName
	ApplicableMatrix  mat.Dense
	AuxiliaryVector   mat.VecDense
	Mutable           bool
	Link              RuleLink
}

// RuleMatrixIsEmpty returns true is the RuleMatrix is uninitialised
func (r *RuleMatrix) RuleMatrixIsEmpty() bool {
	if r.RuleName == "" &&
		len(r.RequiredVariables) == 0 &&
		!r.Mutable &&
		r.Link == (RuleLink{}) {
		// if r.ApplicableMatrix != nil && r.AuxiliaryVector != nil {
		r1, c1 := r.ApplicableMatrix.Dims()
		r2, c2 := r.AuxiliaryVector.Dims()
		if r1 == 0 && c1 == 0 && r2 == 0 && c2 == 0 {
			return true
			// }
		}
	}
	return false

}
