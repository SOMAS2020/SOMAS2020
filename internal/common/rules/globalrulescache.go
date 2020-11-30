package rules

import "gonum.org/v1/gonum/mat"

var AvailableRules = map[string]RuleMatrix{}

func RegisterNewRule(ruleName string, requiredVariables []string, applicableMatrix mat.Dense, auxiliaryVector mat.VecDense) *RuleMatrix {
	if _, ok := AvailableRules[ruleName]; ok {
		//Some form of anger at this point
	}

	rm := RuleMatrix{ruleName: ruleName, RequiredVariables: requiredVariables, ApplicableMatrix: applicableMatrix, AuxiliaryVector: auxiliaryVector}
	AvailableRules[ruleName] = rm
	return &rm
}
