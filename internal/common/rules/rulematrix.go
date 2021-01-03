package rules

import (
	"gonum.org/v1/gonum/mat"
)

type LinkTypeOption int

const (
	ParentFailAutoRulePass LinkTypeOption = iota
	NoLink
)

type RuleLink struct {
	Linked     bool
	LinkType   LinkTypeOption
	LinkedRule string
}

type RuleMatrix struct {
	RuleName          string
	RequiredVariables []VariableFieldName
	ApplicableMatrix  mat.Dense
	AuxiliaryVector   mat.VecDense
	Mutable           bool
	Link              RuleLink
}
