package dynamics

import (
	"fmt"
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
		if len(distances) == 0 {
			return -1
		}
		return distances[getSmallest(distances)]
	}
	return 0.0
}

type Input struct {
	Name             rules.VariableFieldName
	ClientAdjustable bool
	Value            []float64
}

func FindClosestApproach(ruleMatrix rules.RuleMatrix, namedInputs map[rules.VariableFieldName]Input) (namedOutputs map[rules.VariableFieldName]Input) {
	// Get raw data for processing
	droppedInputs := DropAllInputStructs(namedInputs)
	// Evaluate the rule on this data
	results, success := evaluateSingle(ruleMatrix, droppedInputs)
	if success {
		// Find any results in the vector that indicate a failure condition
		deficient, err := IdentifyDeficiencies(results, ruleMatrix.AuxiliaryVector)
		if len(deficient) == 0 {
			// If none are given, the submitted position is good, return
			return namedInputs
		}
		if err == nil {
			// Build selected dynamic structures
			selectedDynamics := BuildSelectedDynamics(ruleMatrix, ruleMatrix.AuxiliaryVector, deficient)
			// We now need to add the dynamics implied by the present inputs (things we can't change)
			// First work out which inputs these are
			immutables := fetchImmutableInputs(namedInputs)
			// Collapse map to list
			allInputs := collapseInputsMap(namedInputs)
			// Work out dimensions of w
			fullSize := calculateFullSpaceSize(allInputs)
			// Select the immutable inputs
			immutableInputs := selectInputs(namedInputs, immutables)
			// Calculate all the extra dynamics
			extraDynamics, success2 := constructAllImmutableDynamics(immutableInputs, ruleMatrix, fullSize)
			if triggerZeroCheck(namedInputs) {
				laced, success := zeroSatisfyCheck(ruleMatrix, namedInputs)
				if success {
					return laced
				}
			}
			if success2 {
				// Append the restrictions to the dynamics
				selectedDynamics = append(selectedDynamics, extraDynamics...)
			}
			copyOfSelectedDynamics := copySelectedDynamics(selectedDynamics)
			// We are ready to use the findClosestApproach internal function
			bestPosition := findClosestApproachInSubspace(ruleMatrix, selectedDynamics, *DecodeValues(ruleMatrix, droppedInputs))
			if ruleevaluation.RuleEvaluation(ruleMatrix, bestPosition) {
				return laceOutputs(bestPosition, ruleMatrix, droppedInputs, namedInputs)
			}
			resolved, left, success3 := combineDefiniteApproaches(copyOfSelectedDynamics, fullSize)
			if success3 {
				bestPosition = findClosestApproachInSubspace(ruleMatrix, append(left, resolved), *DecodeValues(ruleMatrix, droppedInputs))
			}
			return laceOutputs(bestPosition, ruleMatrix, droppedInputs, namedInputs)
		}
	}
	return namedInputs
}

func Shift(ruleMatrix rules.RuleMatrix, namedInputs map[rules.VariableFieldName]Input) (newMatrix rules.RuleMatrix, edited bool) {
	droppedInputs := DropAllInputStructs(namedInputs)
	results, success := evaluateSingle(ruleMatrix, droppedInputs)
	if success {
		// Find any results in the vector that indicate a failure condition
		deficient, err := IdentifyDeficiencies(results, ruleMatrix.AuxiliaryVector)
		if err != nil {
			return ruleMatrix, false
		}
		if len(deficient) == 0 {
			// If none are given, the submitted position is good, return
			return ruleMatrix, false
		}
		allDynamics := BuildAllDynamics(ruleMatrix, ruleMatrix.AuxiliaryVector)
		if len(allDynamics) == 0 {
			return ruleMatrix, false
		}
		for _, row := range deficient {
			copyOfMatrix := mat.DenseCopyOf(&ruleMatrix.ApplicableMatrix)
			if row < len(allDynamics) {
				newDynamic := shiftDynamic(allDynamics[row], DecodeValues(ruleMatrix, droppedInputs))
				newRow := translateDynamic(newDynamic)
				copyOfMatrix.SetRow(row, newRow)
				ruleMatrix.ApplicableMatrix = *copyOfMatrix
			}
		}
		return ruleMatrix, true
	}
	return ruleMatrix, false
}

func translateDynamic(inputDyn dynamic) []float64 {
	w := inputDyn.w
	nRows, _ := w.Dims()
	outputVect := []float64{}
	for i := 0; i < nRows; i++ {
		outputVect = append(outputVect, w.AtVec(i))
	}
	outputVect = append(outputVect, inputDyn.b)
	return outputVect
}

func shiftDynamic(inputDyn dynamic, input mat.Vector) (outputDyn dynamic) {
	res := mat.Dot(&inputDyn.w, input) + inputDyn.b
	switch inputDyn.aux {
	case 0:
		inputDyn.b -= res
		return inputDyn
	case 1:
		inputDyn.b -= res
		inputDyn.b += 1
		return inputDyn
	case 2:
		inputDyn.b -= res
		return inputDyn
	case 3:
		inputDyn.b += 1
		return inputDyn
	default:
		return inputDyn
	}
}

// findClosestApproachInSubspace works out the closest point in the rule subspace to the current location
func findClosestApproachInSubspace(matrixOfRules rules.RuleMatrix, dynamics []dynamic, location mat.VecDense) mat.VecDense {
	if len(dynamics) == 0 {
		return location
	}
	var distances []float64

	for _, v := range dynamics {
		distances = append(distances, findDistanceToHyperplane(v.w, v.b, location))
	}

	indexOfSmall := getSmallestNonZero(distances)
	if indexOfSmall == -1 {
		return location
	}
	closestApproach := calculateClosestApproach(dynamics[indexOfSmall], location)
	if ruleevaluation.RuleEvaluation(matrixOfRules, closestApproach) {
		return closestApproach
	} else {
		xSlice := dynamics[:indexOfSmall]
		ySlice := dynamics[indexOfSmall+1:]
		return findClosestApproachInSubspace(matrixOfRules, append(xSlice, ySlice...), closestApproach)
	}
}

func triggerZeroCheck(namedInputs map[rules.VariableFieldName]Input) bool {
	for _, inp := range namedInputs {
		if !inp.ClientAdjustable {
			if inp.Value[rules.SingleValueVariableEntry] == 0 {
				return true
			}
		}
	}
	return false
}

func zeroSatisfyCheck(ruleMatrix rules.RuleMatrix, namedInputs map[rules.VariableFieldName]Input) (map[rules.VariableFieldName]Input, bool) {
	newMap := make(map[rules.VariableFieldName][]float64)
	for keys, values := range namedInputs {
		if values.ClientAdjustable {
			newMap[keys] = generateEmpty(len(values.Value))
		} else {
			newMap[keys] = values.Value
		}
	}
	if ruleevaluation.RuleEvaluation(ruleMatrix, *DecodeValues(ruleMatrix, newMap)) {
		laced := laceOutputs(*DecodeValues(ruleMatrix, newMap), ruleMatrix, DropAllInputStructs(namedInputs), namedInputs)
		return laced, true
	}
	return map[rules.VariableFieldName]Input{}, false
}

func copySelectedDynamics(orig []dynamic) []dynamic {
	cpy := make([]dynamic, len(orig))
	for i, val := range orig {
		cpy[i] = val
	}
	return cpy
}

// combineDefiniteApproaches is a last ditch attempt to get a defined position if all else fails
func combineDefiniteApproaches(selectedDynamics []dynamic, size int) (dynamic, []dynamic, bool) {
	leftOver := []dynamic{}
	if len(selectedDynamics) != 0 {
		baseDynamic := dynamic{
			w:   *mat.NewVecDense(size, generateEmpty(size)),
			b:   0,
			aux: 12,
		}
		for _, val := range selectedDynamics {
			if val.aux == 0 {
				if baseDynamic.aux == 12 {
					baseDynamic = val
				} else {
					baseDynamic.w.AddScaledVec(&baseDynamic.w, -1, &val.w)
					baseDynamic.b -= val.b
				}
			} else {
				leftOver = append(leftOver, val)
			}
		}
		return baseDynamic, leftOver, true
	}
	return dynamic{}, selectedDynamics, false
}

func generateEmpty(size int) []float64 {
	final := []float64{}
	for i := 0; i < size; i++ {
		final = append(final, 0.0)
	}
	return final
}

// calculateClosestApproach works out the least squares closes point on hyperplane to current location
func calculateClosestApproach(constraint dynamic, location mat.VecDense) mat.VecDense {
	denom := math.Pow(mat.Norm(&constraint.w, 2), 2)
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
	returnDynamics := []dynamic{}
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

func laceOutputs(vector mat.VecDense, ruleMat rules.RuleMatrix, droppedInputs map[rules.VariableFieldName][]float64, original map[rules.VariableFieldName]Input) (namedOutputs map[rules.VariableFieldName]Input) {
	outputs := make(map[rules.VariableFieldName]Input)
	mark := 0
	for _, name := range ruleMat.RequiredVariables {
		input := droppedInputs[name]
		output := []float64{}
		for i := 0; i < len(input); i++ {
			output = append(output, vector.AtVec(mark+i))
		}
		mark = len(input)
		inp := Input{
			Name:             name,
			ClientAdjustable: original[name].ClientAdjustable,
			Value:            output,
		}
		outputs[name] = inp
	}
	return outputs
}

func selectInputs(data map[rules.VariableFieldName]Input, needed []rules.VariableFieldName) []Input {
	outputs := []Input{}
	for _, name := range needed {
		outputs = append(outputs, data[name])
	}
	return outputs
}

func collapseInputsMap(data map[rules.VariableFieldName]Input) []Input {
	inputs := []Input{}
	for _, val := range data {
		inputs = append(inputs, val)
	}
	return inputs
}

func calculateFullSpaceSize(inputs []Input) int {
	size := 0
	for _, input := range inputs {
		size += len(input.Value)
	}
	return size
}

func constructAllImmutableDynamics(immutables []Input, ruleMat rules.RuleMatrix, fullSize int) ([]dynamic, bool) {
	success := true
	allDynamics := []dynamic{}
	for _, immutable := range immutables {
		dyn, succ := constructImmutableDynamic(immutable, ruleMat.RequiredVariables, fullSize)
		success = success && succ
		allDynamics = append(allDynamics, dyn)
	}
	return allDynamics, success
}

func constructImmutableDynamic(immutable Input, reqVar []rules.VariableFieldName, fullSize int) (dynamic, bool) {
	index := getIndexOfVar(immutable.Name, reqVar)
	if index != -1 {
		newDynamic := dynamic{
			aux: 0,
		}
		deltaVect := []float64{}
		for i := 0; i < fullSize; i++ {
			deltaVect = append(deltaVect, 0)
		}
		for i := 0; i < len(immutable.Value); i++ {
			deltaVect[i+index] = 1.0
		}
		w := mat.NewVecDense(len(deltaVect), deltaVect)
		b := sumBValues(immutable.Value)
		newDynamic.w = *w
		newDynamic.b = b
		return newDynamic, true
	}
	return dynamic{}, false
}

func sumBValues(values []float64) float64 {
	val := 0.0
	for _, v := range values {
		val += v
	}
	return -1 * val
}

func getIndexOfVar(name rules.VariableFieldName, reqVar []rules.VariableFieldName) int {
	for index, v := range reqVar {
		if v == name {
			return index
		}
	}
	return -1
}

func fetchImmutableInputs(namedInputs map[rules.VariableFieldName]Input) []rules.VariableFieldName {
	immutables := []rules.VariableFieldName{}
	for name, input := range namedInputs {
		if !input.ClientAdjustable {
			immutables = append(immutables, name)
		}
	}
	return immutables
}

func DropAllInputStructs(inputs map[rules.VariableFieldName]Input) map[rules.VariableFieldName][]float64 {
	outputMap := make(map[rules.VariableFieldName][]float64)
	for key, val := range inputs {
		outputMap[key] = dropInputStruct(val)
	}
	return outputMap
}

func dropInputStruct(input Input) []float64 {
	return input.Value
}

func DecodeValues(rm rules.RuleMatrix, values map[rules.VariableFieldName][]float64) *mat.VecDense {
	var finalVariableVect []float64
	for _, varName := range rm.RequiredVariables {
		if value, ok := values[varName]; ok {
			finalVariableVect = append(finalVariableVect, value...)
		} else {
			return nil
		}
	}
	if len(finalVariableVect) == 0 {
		_ = fmt.Sprintf("Hello")
	}
	varVect := mat.NewVecDense(len(finalVariableVect), finalVariableVect)
	return varVect
}

func decodeWithConst(rm rules.RuleMatrix, values map[rules.VariableFieldName][]float64) *mat.VecDense {
	var finalVariableVect []float64
	for _, varName := range rm.RequiredVariables {
		if value, ok := values[varName]; ok {
			finalVariableVect = append(finalVariableVect, value...)
		} else {
			return nil
		}
	}
	finalVariableVect = append(finalVariableVect, 1)
	varVect := mat.NewVecDense(len(finalVariableVect), finalVariableVect)
	return varVect
}

func evaluateSingle(rm rules.RuleMatrix, values map[rules.VariableFieldName][]float64) (outputVec mat.VecDense, success bool) {
	varVect := decodeWithConst(rm, values)
	if varVect != nil {
		nrows, ncols := rm.ApplicableMatrix.Dims()
		if nrows == 0 || ncols == 0 {
			return mat.VecDense{}, false
		}
		results := ruleevaluation.RuleMul(*varVect, rm.ApplicableMatrix)
		return *results, true
	}
	return mat.VecDense{}, false
}

// identifyDeficiencies checks result of rule evaluation and finds entries of result vector that do not comply
func IdentifyDeficiencies(b mat.VecDense, aux mat.VecDense) ([]int, error) {
	if checkDimensions(b, aux) {
		nRows, _ := b.Dims()
		outputData := []int{}
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

func CollapseRuleMap(input map[string]rules.RuleMatrix) []rules.RuleMatrix {
	newInput := []rules.RuleMatrix{}
	for key, inp := range input {
		if inp.RuleName == "" {
			fmt.Sprintf("hello %v", key)
		}
		newInput = append(newInput, inp)
	}
	return newInput
}

func RemoveFromMap(input map[string]rules.RuleMatrix, ruleName string) []rules.RuleMatrix {
	returnList := []rules.RuleMatrix{}
	for key, val := range input {
		if key != ruleName {
			returnList = append(returnList, val)
		}
	}
	return returnList
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
