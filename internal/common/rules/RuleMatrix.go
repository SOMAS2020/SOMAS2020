package rules

import "gonum.org/v1/gonum/mat"

type RuleMatrix struct {
	ruleName          string
	RequiredVariables []string
	ApplicableMatrix  mat.Dense
	AuxillaryVector   mat.VecDense
}

var AvailableRules = map[string]RuleMatrix{}

func registerNewRule(ruleName string, requiredVariables []string, applicableMatrix mat.Dense, auxillaryVector mat.VecDense) *RuleMatrix {
	if _, ok := AvailableRules[ruleName]; ok {
		//Some form of anger at this point
	}

	rm := RuleMatrix{ruleName: ruleName, RequiredVariables: requiredVariables, ApplicableMatrix: applicableMatrix, AuxillaryVector: auxillaryVector}
	AvailableRules[ruleName] = rm
	return &rm
}

var VariableMap = map[string][]float64{}

func registerNewVariable(variableName string, value float64) {

	if _, ok := AvailableRules[variableName]; ok {
		//Some form of anger at this point
	}

	VariableMap[variableName] = []float64{value}

}

func modifyVariable(variableName string, newValue float64) {
	if _, ok := VariableMap[variableName]; ok {
		VariableMap[variableName] = []float64{newValue}
	}
}
