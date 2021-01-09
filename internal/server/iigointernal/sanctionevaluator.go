package iigointernal

import (
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"github.com/pkg/errors"
	"gonum.org/v1/gonum/mat"
)

// DefaultInitLocalSanctionCache generates a blank sanction cache
func DefaultInitLocalSanctionCache(depth int) map[int][]shared.Sanction {
	returnMap := map[int][]shared.Sanction{}
	for i := 0; i < depth; i++ {
		returnMap[i] = []shared.Sanction{}
	}
	return returnMap
}

// DefaultInitLocalHistoryCache generates a blank history cache
func DefaultInitLocalHistoryCache(depth int) map[int][]shared.Accountability {
	returnMap := map[int][]shared.Accountability{}
	for i := 0; i < depth; i++ {
		returnMap[i] = []shared.Accountability{}
	}
	return returnMap
}

func evaluateSanction(sanction rules.RuleMatrix, localVariableCache map[rules.VariableFieldName]rules.VariableValuePair) shared.Resources {
	reqVar := sanction.RequiredVariables
	if checkAllRequiredVariablesAreAvailable(reqVar, localVariableCache) {
		loadVariableValues := sourceAndAppendAllValues(reqVar, localVariableCache)
		results := matrixOp(loadVariableValues, sanction)
		valid, sanctionRes, err := genRealResult(sanction.AuxiliaryVector, results)
		if err != nil {
			return shared.Resources(0)
		}
		if checkAllValid(valid) {
			return shared.Resources(sanctionRes)
		}
	}
	return 0
}

func checkAllRequiredVariablesAreAvailable(reqVar []rules.VariableFieldName, localVarCache map[rules.VariableFieldName]rules.VariableValuePair) bool {
	for _, v := range reqVar {
		if _, ok := localVarCache[v]; !ok {
			return false
		}
	}
	return true
}

func sourceAndAppendAllValues(reqVar []rules.VariableFieldName, localVarCache map[rules.VariableFieldName]rules.VariableValuePair) []float64 {
	var outputArray []float64
	for _, v := range reqVar {
		outputArray = append(outputArray, localVarCache[v].Values...)
	}
	outputArray = append(outputArray, 1.0)
	return outputArray
}

func matrixOp(variables []float64, matrix rules.RuleMatrix) *mat.VecDense {
	nRows, _ := matrix.ApplicableMatrix.Dims()

	variableFormalVect := mat.NewVecDense(len(variables), variables)

	actual := make([]float64, nRows)

	c := mat.NewVecDense(nRows, actual)

	c.MulVec(&matrix.ApplicableMatrix, variableFormalVect)

	return c
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

func checkAllValid(arr []bool) bool {
	for _, v := range arr {
		if !v {
			return false
		}
	}
	return true
}
