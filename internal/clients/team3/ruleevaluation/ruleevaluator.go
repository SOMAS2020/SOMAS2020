package ruleevaluation

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/pkg/errors"
	"gonum.org/v1/gonum/mat"
)

func ruleMul(variableFormalVect mat.VecDense, ApplicableMatrix mat.Dense) *mat.VecDense {
	nRows, _ := ApplicableMatrix.Dims()

	actual := make([]float64, nRows)

	c := mat.NewVecDense(nRows, actual)

	c.MulVec(&ApplicableMatrix, &variableFormalVect)

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
			return nil, errors.Errorf("At auxillary vector entry: '%v' aux value outside of 0-3: '%v' was found", i, interpret)
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
			return nil, 0.0, errors.Errorf("At auxillary vector entry: '%v' aux value outside of 0-3: '%v' was found", i, interpret)
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
func BasicBooleanRuleEvaluator(matrixRule rules.RuleMatrix, currentVals mat.VecDense) (bool, error) {
	resultVect := ruleMul(currentVals, matrixRule.ApplicableMatrix)
	resultsBool, _ := genResult(matrixRule.AuxiliaryVector, resultVect)
	return checkForFalse(resultsBool), nil
}

// BasicRealValuedRuleEvaluator implements real valued rule evaluation in the same form as the boolean one
func BasicRealValuedRuleEvaluator(matrixRule rules.RuleMatrix, currentVals mat.VecDense) (bool, float64, error) {
	resultVect := ruleMul(currentVals, matrixRule.ApplicableMatrix)
	resultsBool, _, _ := genRealResult(matrixRule.AuxiliaryVector, resultVect)
	return checkForFalse(resultsBool), 0.0, nil
}

func RuleEvaluation(matrixRule rules.RuleMatrix, currentVals mat.VecDense) bool {
	tense := false
	nRows, _ := matrixRule.AuxiliaryVector.Dims()
	output := false
	for i := 0; i < nRows; i++ {
		if matrixRule.AuxiliaryVector.AtVec(i) == 4 {
			output, _, _ = BasicRealValuedRuleEvaluator(matrixRule, currentVals)
			tense = true
		}
	}
	if !tense {
		output, _ = BasicBooleanRuleEvaluator(matrixRule, currentVals)
	}

	return output
}
