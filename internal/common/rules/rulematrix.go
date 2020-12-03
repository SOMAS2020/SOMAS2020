package rules

import (
	"gonum.org/v1/gonum/mat"
)

type RuleMatrix struct {
	ruleName          string
	RequiredVariables []string
	ApplicableMatrix  mat.Dense
	AuxiliaryVector   mat.VecDense
}
