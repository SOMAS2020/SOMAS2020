package rules

import (
	"github.com/pkg/errors"
	"gonum.org/v1/gonum/mat"
)

// ComplianceCheck Gives a simple check for compliance to clients, they just need to feed in a RuleMatrix and the
// variable cache they want to feed from
func ComplianceCheck(rule RuleMatrix, variables map[VariableFieldName]VariableValuePair, inPlayRules map[string]RuleMatrix) (compliant bool, ruleError error) {
	if checkAllVariablesAvailable(rule.RequiredVariables, variables) {
		ruleCache := map[string]RuleMatrix{}
		ruleCache[rule.RuleName] = rule
		if rule.Link.Linked {
			ruleCache[rule.Link.LinkedRule] = inPlayRules[rule.Link.LinkedRule]
		}
		evalResult := EvaluateRuleFromCaches(rule.RuleName, ruleCache, variables)
		return evalResult.RulePasses, evalResult.EvalError
	}
	return false, &RuleError{
		ErrorType: VariableCacheDidNotHaveAllRequiredVariables,
		Err:       errors.Errorf("Variable cache provided didn't have all the errors for this rule"),
	}
}

// ComplianceRecommendation *attempts* to calculate a set of variables that comply with a rule if the
// input variables don't satisfy it already
func ComplianceRecommendation(rule RuleMatrix, variables map[VariableFieldName]VariableValuePair) (map[VariableFieldName]VariableValuePair, bool) {
	slicedValues, fixed, success := fetchRequiredVariables(rule.RequiredVariables, variables)
	if success {
		nRows, _ := rule.ApplicableMatrix.Dims()
		for i := 0; i < nRows; i++ {
			row := rule.ApplicableMatrix.RowView(i)
			newVals, changed := analysisEngine(row, slicedValues, fixed, rule.AuxiliaryVector.AtVec(i))
			if changed {
				slicedValues = newVals
			}
		}
		return unfetchRequiredVariables(rule.RequiredVariables, variables, slicedValues), true

	}
	return variables, false
}

func analysisEngine(line mat.Vector, currentVal []float64, fixed []bool, auxCode float64) ([]float64, bool) {
	forAnalysis := deepCopyList(currentVal)
	forAnalysis = append(forAnalysis, 1)
	currentVect := mat.NewVecDense(len(forAnalysis), forAnalysis)
	res := mat.Dot(line, currentVect)
	if satisfy(res, auxCode) {
		return currentVal, true
	}
	adjustable := findFirst(true, fixed, line)
	if adjustable == -1 {
		return currentVal, false
	}
	newEntryVal := recommendCorrection(currentVal[adjustable], res, auxCode, line.AtVec(adjustable))
	currentVal[adjustable] = newEntryVal
	return currentVal, true
}

func deepCopyList(inp []float64) []float64 {
	newList := make([]float64, len(inp))
	copy(newList, inp)
	return newList
}

func findFirst(search bool, values []bool, line mat.Vector) (index int) {
	for i, val := range values {
		if search == val && line.AtVec(i) != 0 {
			return i
		}
	}
	return -1
}

func satisfy(res float64, aux float64) bool {
	if aux == 0 {
		return res == 0
	}
	if aux == 1 {
		return res > 0
	}
	if aux == 2 {
		return res >= 0
	}
	if aux == 3 {
		return res != 0
	}
	if aux == 4 {
		return true
	}
	return false
}

func recommendCorrection(prevVal float64, res float64, aux float64, multiplier float64) float64 {
	calc := res - prevVal*multiplier
	if aux == 0 {
		return (-1 * calc) / multiplier
	}
	if aux == 1 {
		return (-2*calc + 1) / multiplier
	}
	if aux == 2 {
		return (-1 * calc) / multiplier
	}
	if aux == 3 {
		return (calc + 1) / multiplier
	}
	if aux == 4 {
		return calc / multiplier
	}
	return prevVal
}

func fetchRequiredVariables(reqVariables []VariableFieldName, variables map[VariableFieldName]VariableValuePair) ([]float64, []bool, bool) {
	finalSlice := []float64{}
	returnFix := []bool{}
	for _, variable := range reqVariables {
		if value, ok := variables[variable]; ok {
			finalSlice = append(finalSlice, value.Values...)
			returnFix = append(returnFix, generateBoolList(IsChangeable()[variable], len(value.Values))...)
		} else {
			return finalSlice, returnFix, false
		}
	}
	return finalSlice, returnFix, true
}

func unfetchRequiredVariables(reqVariables []VariableFieldName, variables map[VariableFieldName]VariableValuePair, newData []float64) map[VariableFieldName]VariableValuePair {
	pointer := 0
	for _, val := range reqVariables {
		pair := variables[val]
		length := len(pair.Values)
		pair.Values = newData[pointer : pointer+length]
		pointer += length
		variables[val] = pair
	}
	return variables
}

func generateBoolList(initial bool, length int) []bool {
	finalList := []bool{}
	for i := 0; i < length; i++ {
		finalList = append(finalList, initial)
	}
	return finalList
}

func IsChangeable() map[VariableFieldName]bool {
	return map[VariableFieldName]bool{
		NumberOfIslandsContributingToCommonPool: false,
		NumberOfFailedForages:                   false,
		NumberOfBrokenAgreements:                false,
		MaxSeverityOfSanctions:                  false,
		NumberOfIslandsAlive:                    false,
		NumberOfBallotsCast:                     false,
		NumberOfAllocationsSent:                 false,
		AllocationRequestsMade:                  false,
		AllocationMade:                          false,
		IslandsAlive:                            false,
		SpeakerSalary:                           true,
		JudgeSalary:                             true,
		PresidentSalary:                         true,
		RuleSelected:                            false,
		VoteCalled:                              false,
		ExpectedTaxContribution:                 false,
		ExpectedAllocation:                      false,
		IslandTaxContribution:                   true,
		IslandAllocation:                        true,
		IslandReportedResources:                 false,
		ConstSanctionAmount:                     false,
		TurnsLeftOnSanction:                     false,
		SanctionPaid:                            true,
		SanctionExpected:                        false,
		TestVariable:                            false,
		JudgeInspectionPerformed:                false,
	}
}
