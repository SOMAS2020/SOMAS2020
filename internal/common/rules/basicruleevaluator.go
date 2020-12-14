package rules

import (
	"fmt"
	"github.com/pkg/errors"
	"gonum.org/v1/gonum/mat"
)

func createVarList(variables []string) ([]float64, error) {
	var variableVect []float64

	for _, v := range variables {
		if val, varOk := VariableMap[v]; varOk {
			variableVect = append(variableVect, val.Values...)
		} else {
			return nil, errors.New(fmt.Sprintf("Variable: '%v' not found in global variable cache", v))
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
			return nil, errors.New(fmt.Sprintf("At auxillary vector entry: '%v' aux value outside of 0-3: '%v' was found", i, interpret))
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
		case 4:
			outputVal = c.AtVec(i)
		default:
			return nil, 0.0, errors.New(fmt.Sprintf("At auxillary vector entry: '%v' aux value outside of 0-3: '%v' was found", i, interpret))
		}
		resultVect = append(resultVect, res)
	}

	return resultVect, outputVal, nil
}

func checkForFalse(resultVect []bool) bool {

	var finalBool = true

	for _, v := range resultVect {
		finalBool = finalBool && v
	}

	return finalBool

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
	} else {
		return false, errors.Errorf("rule name: '%v' provided doesn't exist in global rule list", ruleName)
	}
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
	} else {
		return false, 0, errors.Errorf("rule name: '%v' provided doesn't exist in global rule list", ruleName)
	}
}
