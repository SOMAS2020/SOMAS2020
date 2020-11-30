package rules

import (
	"errors"
	"gonum.org/v1/gonum/mat"
)

type RuleMatrix struct {
	ruleName          string
	RequiredVariables []string
	ApplicableMatrix  mat.Dense
	AuxiliaryVector   mat.VecDense
}

var AvailableRules = map[string]RuleMatrix{}

func registerNewRule(ruleName string, requiredVariables []string, applicableMatrix mat.Dense, auxiliaryVector mat.VecDense) *RuleMatrix {
	if _, ok := AvailableRules[ruleName]; ok {
		//Some form of anger at this point
	}

	rm := RuleMatrix{ruleName: ruleName, RequiredVariables: requiredVariables, ApplicableMatrix: applicableMatrix, AuxiliaryVector: auxiliaryVector}
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

func registerNewSliceVariable(variableName string, value []float64) {
	if _, ok := AvailableRules[variableName]; ok {
		//Some form of anger at this point
	}

	VariableMap[variableName] = value
}

func modifyVariable(variableName string, newValue float64) error {
	if _, ok := VariableMap[variableName]; ok {
		VariableMap[variableName] = []float64{newValue}
		return nil
	} else {
		return errors.New("attempted to modify a variable has not been defined")
	}
}
