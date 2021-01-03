package rules

import (
	"gonum.org/v1/gonum/mat"
)

// LinkTypeOption gives an enumerated type for the various link types available for rules
type LinkTypeOption int

const (
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
