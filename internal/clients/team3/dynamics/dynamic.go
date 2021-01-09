package dynamics

import (
	"gonum.org/v1/gonum/mat"
)

type dynamic struct {
	w   mat.VecDense // Hyperplane weights
	b   float64      // Offset
	aux float64      // Aux code
}

func createDynamic(weights mat.VecDense, offset float64, auxCode float64) dynamic {
	return dynamic{
		w:   weights,
		b:   offset,
		aux: auxCode,
	}
}
