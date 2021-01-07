package dynamics

import (
	"github.com/pkg/errors"
	"gonum.org/v1/gonum/mat"
	"math"
)

func getSmallestNonZero(values []float64) int {
	newList := []float64{}
	for _, val := range values {
		if val != 0 {
			newList = append(newList, val)
		}
	}
	return getSmallest(newList)
}

func getSmallest(values []float64) int {
	if len(values) == 0 {
		return -1
	}
	currentLowest := values[0]
	currentLowestIndex := 0
	for i, v := range values {
		if v <= currentLowest {
			currentLowest = v
			currentLowestIndex = i
		}
	}
	return currentLowestIndex
}

func findInSlice(input []int, val int) bool {
	for _, v := range input {
		if v == val {
			return true
		}
	}
	return false
}

// findDistanceToHyperplane calculates distance between x and hyperplane defined by w and b in the form w*x + b = 0
func findDistanceToHyperplane(w mat.VecDense, b float64, x mat.VecDense) float64 {
	v := mat.Dot(&w, &x) + b
	v = math.Abs(v)
	scaling := float64(mat.Norm(&w, 2))
	return v / scaling
}

func applyDynamic(dyn dynamic, location mat.VecDense) float64 {
	nRows, _ := dyn.w.Dims()
	pRows, _ := location.Dims()
	if nRows == pRows{
		return mat.Dot(&dyn.w, &location) + dyn.b
	}
	return 0
}

// calculateDistanceVecVec calculates the euclidean distance between two vectors
func calculateDistanceVecVec(a mat.VecDense, b mat.VecDense) (float64, error) {
	if checkDimensions(a, b) {
		v := a
		v.SubVec(&a, &b)
		return float64(v.Len()), nil
	} else {
		return 0.0, errors.Errorf("Vectors '%v' and '%v' do not share the same dimensions", a, b)
	}
}

// checkDimensions will compare two vectors to see if they have the same dimensions
func checkDimensions(a mat.VecDense, b mat.VecDense) bool {
	ax, ay := a.Dims()
	bx, by := b.Dims()
	return ax == bx && ay == by
}
