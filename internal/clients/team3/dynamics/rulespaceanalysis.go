package dynamics

import (
	"github.com/SOMAS2020/SOMAS2020/internal/clients/team3/ruleevaluation"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"gonum.org/v1/gonum/mat"
)

func CalculateDistanceFromRuleSpace(allRules []rules.RuleMatrix, namedInputs map[rules.VariableFieldName]Input) float64 {
	if !checkAllCompliant(allRules, namedInputs) {
		distances := make(map[int]float64)
		for entry, rule := range allRules {
			dist := getDistance(rule, namedInputs)
			distances[entry] = dist
		}
		compliant := false
		for i := 0; !compliant; i++ {
			if len(distances) == 0 {
				compliant = true
			}
			currentSmallest := getSmallestInMap(distances)
			if currentSmallest == -1 {
				return -1
			}
			valRule := allRules[currentSmallest]
			reqMap := sourceRequiredInputs(valRule, namedInputs)
			newPoint := FindClosestApproach(valRule, reqMap)
			newMap := softMergeInputs(namedInputs, newPoint)
			if checkAllCompliant(allRules, newMap) {
				return distances[currentSmallest]
			} else {
				delete(distances, currentSmallest)
			}
		}
		return -1
	} else {
		return -1
	}
}

func removeFromDistances(s []float64, i int) []float64 {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

func getSmallestInMap(inputMap map[int]float64) int {
	if len(inputMap) > 0 {
		base := inputMap[0]
		index := 0
		for key, val := range inputMap {
			if val <= base {
				base = val
				index = key
			}
		}
		return index
	}
	return -1
}

func getDistance(rule rules.RuleMatrix, namedInputs map[rules.VariableFieldName]Input) float64 {
	reqMap := sourceRequiredInputs(rule, namedInputs)
	rawMap := dropAllInputStructs(reqMap)
	vectData := decodeValues(rule, rawMap)
	allDynamics := BuildAllDynamics(rule, rule.AuxiliaryVector)
	return GetDistanceToSubspace(allDynamics, *vectData)
}

func sourceRequiredInputs(rule rules.RuleMatrix, namedInputs map[rules.VariableFieldName]Input) map[rules.VariableFieldName]Input {
	newMap := make(map[rules.VariableFieldName]Input)
	for _, val := range rule.RequiredVariables {
		newMap[val] = namedInputs[val]
	}
	return newMap
}

func softMergeInputs(namedInputs map[rules.VariableFieldName]Input, newVals map[rules.VariableFieldName]Input) map[rules.VariableFieldName]Input {
	newMap := make(map[rules.VariableFieldName]Input)
	for key, val := range namedInputs {
		if newVal, ok := newVals[key]; ok {
			newMap[key] = newVal
		} else {
			newMap[key] = val
		}
	}
	return newMap
}

func checkAllCompliant(allRules []rules.RuleMatrix, namedInputs map[rules.VariableFieldName]Input) bool {
	cast := true
	for _, rule := range allRules {
		cast = cast && checkForCompliance(rule, namedInputs)
	}
	return cast
}

func checkForCompliance(rule rules.RuleMatrix, namedInputs map[rules.VariableFieldName]Input) bool {
	required := rule.RequiredVariables
	data := []float64{}
	for _, val := range required {
		if entry, ok := namedInputs[val]; ok {
			data = append(data, entry.Value...)
		} else {
			return false
		}
	}
	return ruleevaluation.RuleEvaluation(rule, *mat.NewVecDense(len(data), data))
}
