package dynamics

import (
	"github.com/SOMAS2020/SOMAS2020/internal/clients/team3/ruleevaluation"
	"github.com/SOMAS2020/SOMAS2020/internal/common/rules"
	"github.com/pkg/errors"
	"gonum.org/v1/gonum/mat"
	"math"
)

func GetDistanceToSubspace(dynamics []dynamic, location mat.VecDense) float64 {
	if len(dynamics) == 0 {
		return 0.0
	}
	var distances []float64
	satisfied := true
	for _, dyn := range dynamics {
		satisfied = satisfied && satisfy(applyDynamic(dyn, location), dyn.aux)
	}
	if satisfied {
		return -1
	}
	for _, v := range dynamics {
		distances = append(distances, findDistanceToHyperplane(v.w, v.b, location))
	}
	if distances != nil {
		return distances[getSmallest(distances)]
	}
	return 0.0
}

// findClosestApproachInSubspace works out the closest point in the rule subspace to the current location
func FindClosestApproachInSubspace(matrixOfRules rules.RuleMatrix, dynamics []dynamic, location mat.VecDense) mat.VecDense {
	if len(dynamics) == 0 {
		return location
	}
	var distances []float64

	for _, v := range dynamics {
		distances = append(distances, findDistanceToHyperplane(v.w, v.b, location))
	}

	indexOfSmall := getSmallest(distances)
	closestApproach := calculateClosestApproach(dynamics[indexOfSmall], location)
	if ruleevaluation.RuleEvaluation(matrixOfRules, location) {
		return closestApproach
	} else {
		xSlice := dynamics[:indexOfSmall]
		ySlice := dynamics[indexOfSmall+1:]
		return FindClosestApproachInSubspace(matrixOfRules, append(xSlice, ySlice...), location)
	}
}

// calculateClosestApproach works out the least squares closes point on hyperplane to current location
func calculateClosestApproach(constraint dynamic, location mat.VecDense) mat.VecDense {
	denom := math.Pow(float64(constraint.w.Len()), 2)
	numer := -1 * constraint.b
	nRows, _ := location.Dims()
	for i := 0; i < nRows; i++ {
		numer -= location.AtVec(i) * constraint.w.AtVec(i)
	}
	t := numer / denom
	closestApproach := mat.VecDenseCopyOf(&location)
	closestApproach.AddScaledVec(closestApproach, t, &constraint.w)
	return *closestApproach
}

// buildSelectedDynamics depending on list of indexes, this function will build dynamics
func BuildSelectedDynamics(matrixOfRules rules.RuleMatrix, auxiliaryVector mat.VecDense, selectedRules []int) []dynamic {
	matrixVal := matrixOfRules.ApplicableMatrix
	nRows, _ := matrixVal.Dims()
	var returnDynamics []dynamic
	for i := 0; i < nRows; i++ {
		if findInSlice(selectedRules, i) {
			tempWeights := mat.VecDenseCopyOf(matrixVal.RowView(i))
			returnDynamics = append(returnDynamics, constructSingleDynamic(*tempWeights, auxiliaryVector.AtVec(i)))
		}
	}
	return returnDynamics
}

// buildAllDynamics takes an entire rule matrix and returns all dynamics from it
func BuildAllDynamics(matrixOfRules rules.RuleMatrix, auxiliaryVector mat.VecDense) []dynamic {
	matrixVal := matrixOfRules.ApplicableMatrix
	nRows, _ := matrixVal.Dims()
	var returnDynamics []dynamic
	for i := 0; i < nRows; i++ {
		tempWeights := mat.VecDenseCopyOf(matrixVal.RowView(i))
		returnDynamics = append(returnDynamics, constructSingleDynamic(*tempWeights, auxiliaryVector.AtVec(i)))
	}
	return returnDynamics
}

// constructSingleDynamic builds a dynamic from a row of the rule matrix
func constructSingleDynamic(ruleRow mat.VecDense, auxCode float64) dynamic {
	nRows, _ := ruleRow.Dims()

	offset := ruleRow.AtVec(nRows - 1)

	weights := mat.VecDenseCopyOf(ruleRow.SliceVec(0, nRows-1))

	return createDynamic(*weights, offset, auxCode)

}

// findMinimumRequirements takes in addresses of deficient variables and calculates the minimum values required at those entries to satisfy the rule
func FindMinimumRequirements(deficients []int, aux mat.VecDense) []float64 {
	var outputMins []float64
	for _, v := range deficients {
		switch interpret := aux.AtVec(v); interpret {
		case 0:
			outputMins = append(outputMins, 0)
		case 1:
			outputMins = append(outputMins, 1)
		case 2:
			outputMins = append(outputMins, 0)
		case 3:
			outputMins = append(outputMins, 1)
		default:
			outputMins = append(outputMins, 0)
		}
	}
	return outputMins
}

// identifyDeficiencies checks result of rule evaluation and finds entries of result vector that do not comply
func IdentifyDeficiencies(b mat.VecDense, aux mat.VecDense) ([]int, error) {
	if checkDimensions(b, aux) {
		nRows, _ := b.Dims()
		var outputData []int
		for i := 0; i < nRows; i++ {
			if aux.AtVec(i) > 4 || aux.AtVec(i) < 0 {
				return []int{}, errors.Errorf("Auxilliary vector at entry '%v' has aux code out of range: '%v'", i, aux.AtVec(i))
			}
			if !satisfy(b.AtVec(i), aux.AtVec(i)) {
				outputData = append(outputData, i)
			}
		}
		return outputData, nil
	} else {
		return []int{}, errors.Errorf("Vectors '%v' and '%v' do not have the same dimensions", b, aux)
	}
}

// satisfy checks whether a condition is met based on result vector value and auxiliary code
func satisfy(x float64, a float64) bool {
	switch a {
	case 0:
		return x == 0
	case 1:
		return x > 0
	case 2:
		return x >= 0
	case 3:
		return x != 0
	case 4:
		return true
	default:
		return false
	}
}
