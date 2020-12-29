package rules

import (
	"gonum.org/v1/gonum/mat"
)

type RuleMatrix struct {
	RuleName          string
	RequiredVariables []VariableFieldName
	ApplicableMatrix  mat.Dense
	AuxiliaryVector   mat.VecDense
	Mutable           bool
}
