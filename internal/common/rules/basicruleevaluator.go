package rules

import (
	"gonum.org/v1/gonum/mat"

	"github.com/pkg/errors"
)

func createVarList(variables []VariableFieldName, varCache map[VariableFieldName]VariableValuePair) ([]float64, error) {
	var variableVect []float64

	for _, v := range variables {
		if val, varOk := varCache[v]; varOk {
			variableVect = append(variableVect, val.Values...)
		} else {
			return nil, errors.Errorf("Variable: '%v' not found in variable cache %v", v, varCache)
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

func checkForCode4(auxVect mat.VecDense) bool {
	nRows, _ := auxVect.Dims()
	for i := 0; i < nRows; i++ {
		if auxVect.AtVec(i) == 4.0 {
			return true
		}
	}
	return false
}

// basicBooleanRuleEvaluator implements a basic version of the Matrix rule evaluator, provides single boolean output (and error if present)
func basicBooleanRuleEvaluator(rule RuleMatrix, variableCache map[VariableFieldName]VariableValuePair) (bool, error) {

	variableVect, err := createVarList(rule.RequiredVariables, variableCache)
	if err != nil {
		return false, &RuleError{
			ErrorType: VariableCacheDidNotHaveAllRequiredVariables,
			Err:       errors.Errorf("Variable cache did not contain all required variables"),
		}
	}

	//Checking dimensions line up
	_, nCols := rule.ApplicableMatrix.Dims()

	if nCols != len(variableVect) {
		return false, &RuleError{
			ErrorType: VariableVectDimsDoNotMatchRuleMatrix,
			Err:       errors.Errorf("Variable Vector Dimensions do not match the rule matrix"),
		}
	}

	c := ruleMul(variableVect, rule.ApplicableMatrix)

	resultVect, err := genResult(rule.AuxiliaryVector, c)
	if err != nil {
		return false, &RuleError{
			ErrorType: AuxVectorCodeOutOfRange,
			Err:       err,
		}
	}

	return checkForFalse(resultVect), nil

}

// basicRealValuedRuleEvaluator implements real valued rule evaluation in the same form as the boolean one
func basicRealValuedRuleEvaluator(rule RuleMatrix, variableCache map[VariableFieldName]VariableValuePair) (bool, float64, error) {

	variableVect, err := createVarList(rule.RequiredVariables, variableCache)
	if err != nil {
		return false, 0, &RuleError{
			ErrorType: VariableCacheDidNotHaveAllRequiredVariables,
			Err:       errors.Errorf("Variable cache did not contain all required variables"),
		}
	}

	//Checking dimensions line up
	_, nCols := rule.ApplicableMatrix.Dims()

	if nCols != len(variableVect) {
		return false, 0, &RuleError{
			ErrorType: VariableVectDimsDoNotMatchRuleMatrix,
			Err:       errors.Errorf("Variable Vector Dimensions do not match the rule matrix"),
		}

	}

	c := ruleMul(variableVect, rule.ApplicableMatrix)

	resultVect, outputVal, err := genRealResult(rule.AuxiliaryVector, c)
	if err != nil {
		return false, 0, &RuleError{
			ErrorType: AuxVectorCodeOutOfRange,
			Err:       err,
		}
	}

	return checkForFalse(resultVect), outputVal, nil

}

// basicLinkedRuleEvaluator evaluates linked rules using the above two evaluators
func basicLinkedRuleEvaluator(rule RuleMatrix, childRule RuleMatrix, variableCache map[VariableFieldName]VariableValuePair) (bool, error) {
	link := rule.Link
	if !link.Linked {
		return basicBooleanRuleEvaluator(rule, variableCache)
	}
	if link.LinkType == ParentFailAutoRulePass {
		parentPass, parentErr := basicBooleanRuleEvaluator(rule, variableCache)
		if parentErr != nil {
			return false, errors.Errorf("Parent Rule errored out with : %v", parentErr)
		}
		if !parentPass {
			return true, nil
		}
		childPass, childErr := basicBooleanRuleEvaluator(childRule, variableCache)
		if childErr != nil {
			return false, errors.Errorf("Parent Rule errored out with : %v", childErr)
		}
		return childPass, nil
	}
	return false, errors.Errorf("Unrecognised rule linking %v", link.LinkType)
}

// RuleEvaluationReturn provides a wrapped for the results of a rule evaluation
type RuleEvaluationReturn struct {
	RulePasses    bool
	IsRealOutput  bool
	RealOutputVal float64
	EvalError     error
}

func EvaluateRuleFromCaches(ruleName string, rulesCache map[string]RuleMatrix, variableCache map[VariableFieldName]VariableValuePair) RuleEvaluationReturn {
	if rule, ok := rulesCache[ruleName]; ok {
		if checkAllVariablesAvailable(rule.RequiredVariables, variableCache) {
			auxVect := rule.AuxiliaryVector
			linked := rule.Link.Linked
			if linked {
				if childRule, isInCache := rulesCache[rule.Link.LinkedRule]; isInCache {
					eval, err := basicLinkedRuleEvaluator(rule, childRule, variableCache)
					return RuleEvaluationReturn{
						RulePasses:    eval,
						IsRealOutput:  false,
						RealOutputVal: 0,
						EvalError:     err,
					}
				}
				return RuleEvaluationReturn{
					RulePasses:    false,
					IsRealOutput:  false,
					RealOutputVal: 0,
					EvalError: &RuleError{
						ErrorType: ChildRuleNotFound,
						Err:       errors.Errorf("Child rule was not found in cache"),
					},
				}
			}
			isRealValued := checkForCode4(auxVect)
			if isRealValued {
				eval, res, err := basicRealValuedRuleEvaluator(rule, variableCache)
				return RuleEvaluationReturn{
					RulePasses:    eval,
					IsRealOutput:  true,
					RealOutputVal: res,
					EvalError:     err,
				}
			}
			eval, err := basicBooleanRuleEvaluator(rule, variableCache)
			return RuleEvaluationReturn{
				RulePasses:    eval,
				IsRealOutput:  false,
				RealOutputVal: 0,
				EvalError:     err,
			}
		}
		return RuleEvaluationReturn{
			RulePasses:    false,
			IsRealOutput:  false,
			RealOutputVal: 0,
			EvalError: &RuleError{
				ErrorType: VariableCacheDidNotHaveAllRequiredVariables,
				Err:       errors.Errorf("Provided variable cache did not contain all required variables"),
			},
		}
	}
	return RuleEvaluationReturn{
		RulePasses:    false,
		IsRealOutput:  false,
		RealOutputVal: 0,
		EvalError: &RuleError{
			ErrorType: RuleNotInAvailableRulesCache,
			Err:       errors.Errorf("Rule was not found in the given cache"),
		},
	}
}
