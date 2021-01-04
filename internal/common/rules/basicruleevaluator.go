package rules

import (
	"github.com/pkg/errors"
	"gonum.org/v1/gonum/mat"
)

func createVarList(variables []VariableFieldName) ([]float64, error) {
	return createVarListInternal(variables, VariableMap)
}

func createVarListInternal(variables []VariableFieldName, variableMap map[VariableFieldName]VariableValuePair) ([]float64, error) {
	var variableVect []float64

	for _, v := range variables {
		if val, varOk := variableMap[v]; varOk {
			variableVect = append(variableVect, val.Values...)
		} else {
			return nil, errors.Errorf("Variable: '%v' not found in variable cache %v", v, variableMap)
		}
	}

	variableVect = append(variableVect, 1)
	return variableVect, nil
}

func ruleMul(variableVect []float64, ApplicableMatrix mat.Dense) *mat.VecDense {
	nRows, _ := ApplicableMatrix.Dims()

	variableFormalVect := mat.NewVecDense(len(variableVect), variableVect)

	actual := make([]float64, nRows)

	c := mat.NewVecDense(nRows, actual)

	c.MulVec(&ApplicableMatrix, variableFormalVect)

	return c

}

func genResult(aux mat.VecDense, c *mat.VecDense) ([]bool, error) {
	var resultVect []bool
	nRows, _ := c.Dims()

	for i := 0; i < nRows; i++ {
		res := false
		switch interpret := aux.AtVec(i); interpret {
		case 0:
			res = c.AtVec(i) == 0
		case 1:
			res = c.AtVec(i) > 0
		case 2:
			res = c.AtVec(i) >= 0
		case 3:
			res = c.AtVec(i) != 0
		default:
			return nil, errors.Errorf("At auxiliary vector entry: '%v' aux value outside of 0-3: "+
				"'%v' was found", i, interpret)
		}
		resultVect = append(resultVect, res)
	}

	return resultVect, nil

}

func genRealResult(aux mat.VecDense, c *mat.VecDense) ([]bool, float64, error) {
	var resultVect []bool
	nRows, _ := c.Dims()

	outputVal := 0.0

	for i := 0; i < nRows; i++ {
		res := true
		switch interpret := aux.AtVec(i); interpret {
		case 0:
			res = c.AtVec(i) == 0
		case 1:
			res = c.AtVec(i) > 0
		case 2:
			res = c.AtVec(i) >= 0
		case 3:
			res = c.AtVec(i) != 0
		case 4:
			outputVal = c.AtVec(i)
		default:
			return nil, 0.0, errors.Errorf("At auxiliary vector entry: '%v' aux value outside of 0-3: "+
				"'%v' was found", i, interpret)
		}
		resultVect = append(resultVect, res)
	}

	return resultVect, outputVal, nil
}

func checkForFalse(resultVect []bool) bool {

	for _, v := range resultVect {
		if !v {
			return false
		}
	}

	return true

}

// BasicBooleanRuleEvaluator implements a basic version of the Matrix rule evaluator, provides single boolean output (and error if present)
func BasicBooleanRuleEvaluator(ruleName string) (bool, error) {

	if rm, ok := AvailableRules[ruleName]; ok {
		variableVect, err := createVarList(rm.RequiredVariables)
		if err != nil {
			return false, err
		}

		//Checking dimensions line up
		_, nCols := rm.ApplicableMatrix.Dims()

		if nCols != len(variableVect) {
			return false, errors.Errorf(
				"dimension mismatch in evaluating rule: '%v' rule matrix has '%v' columns, while we sourced '%v' variables",
				ruleName,
				nCols,
				len(variableVect),
			)
		}

		c := ruleMul(variableVect, rm.ApplicableMatrix)

		resultVect, err := genResult(rm.AuxiliaryVector, c)
		if err != nil {
			return false, err
		}

		return checkForFalse(resultVect), nil
	}
	return false, errors.Errorf("rule name: '%v' provided doesn't exist in global rule list", ruleName)
}

// BasicRealValuedRuleEvaluator implements real valued rule evaluation in the same form as the boolean one
func BasicRealValuedRuleEvaluator(ruleName string) (bool, float64, error) {
	if rm, ok := AvailableRules[ruleName]; ok {
		variableVect, err := createVarList(rm.RequiredVariables)
		if err != nil {
			return false, 0, err
		}

		//Checking dimensions line up
		_, nCols := rm.ApplicableMatrix.Dims()

		if nCols != len(variableVect) {
			return false, 0, errors.Errorf(
				"dimension mismatch in evaluating rule: '%v' rule matrix has '%v' columns, while we sourced '%v' variables",
				ruleName,
				nCols,
				len(variableVect),
			)
		}

		c := ruleMul(variableVect, rm.ApplicableMatrix)

		resultVect, outputVal, err := genRealResult(rm.AuxiliaryVector, c)
		if err != nil {
			return false, 0, err
		}

		return checkForFalse(resultVect), outputVal, nil
	}
	return false, 0, errors.Errorf("rule name: '%v' provided doesn't exist in global rule list", ruleName)
}

// BasicLinkedRuleEvaluator evaluates linked rules using the above two evaluators
func BasicLinkedRuleEvaluator(ruleName string) (bool, error) {
	if rm, ok := AvailableRules[ruleName]; ok {
		link := rm.Link
		if !link.Linked {
			return BasicBooleanRuleEvaluator(ruleName)
		}
		if link.LinkType == ParentFailAutoRulePass {
			childRule := link.LinkedRule
			parentPass, parentErr := BasicBooleanRuleEvaluator(ruleName)
			if parentErr != nil {
				return false, errors.Errorf("Parent Rule errored out with : %v", parentErr)
			}
			if !parentPass {
				return true, nil
			}
			childPass, childErr := BasicBooleanRuleEvaluator(childRule)
			if childErr != nil {
				return false, errors.Errorf("Parent Rule errored out with : %v", childErr)
			}
			return childPass, nil
		}
		return false, errors.Errorf("Unrecognised rule linking %v", link.LinkType)
	}
	return false, errors.Errorf("rule name: '%v' provided doesn't exist in global rule list", ruleName)
}

// BasicLocalBooleanRuleEvaluator implements a basic version local rule evaluator, provides single boolean output (and error if present)
func BasicLocalBooleanRuleEvaluator(rule RuleMatrix, vars map[VariableFieldName]VariableValuePair) (bool, error) {
	//Checking dimensions line up
	variableVect, err := createVarListInternal(rule.RequiredVariables, vars)
	if err != nil {
		return false, err
	}

	_, nCols := rule.ApplicableMatrix.Dims()

	if nCols != len(variableVect) {
		return false, errors.Errorf(
			"dimension mismatch in evaluating rule: '%v' rule matrix has '%v' columns, while we sourced '%v' variables",
			rule.RuleName,
			nCols,
			len(vars),
		)
	}

	c := ruleMul(variableVect, rule.ApplicableMatrix)

	resultVect, err := genResult(rule.AuxiliaryVector, c)
	if err != nil {
		return false, err
	}

	return checkForFalse(resultVect), nil
}
