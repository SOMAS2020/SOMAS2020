package rules

import "gonum.org/v1/gonum/mat"

var AvailableRules = map[string]RuleMatrix{}

// RegisterNewRule Creates and registers new rule based on inputs
func RegisterNewRule(ruleName string, requiredVariables []string, applicableMatrix mat.Dense, auxiliaryVector mat.VecDense) *RuleMatrix {
	if _, ok := AvailableRules[ruleName]; ok {
		panic("attempt to re-register rule with name " + ruleName)
	}

	rm := RuleMatrix{ruleName: ruleName, RequiredVariables: requiredVariables, ApplicableMatrix: applicableMatrix, AuxiliaryVector: auxiliaryVector}
	AvailableRules[ruleName] = rm
	return &rm
}
