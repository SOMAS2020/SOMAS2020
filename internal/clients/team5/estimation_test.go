package team5

import (
	"math"
	"testing"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/stat/distuv"
)

// this test only really tests a small helper function. The real KDE test was done by looking at output
// values and comparing to samples from the generating distribution. See "Estimation" in README.md in this package.
func TestKDE(t *testing.T) {

	nSamples := 10
	dist := distuv.Normal{Mu: 10, Sigma: 2, Src: rand.NewSource(1)}
	obs := make([]float64, nSamples)
	for i := 0; i < nSamples; i++ {
		obs[i] = dist.Rand()
	}
	m := kdeModel{observations: obs, bandwidth: 1.0, weights: nil}

	xMin := 0.0
	xMax := 20.0
	step := 0.5
	result := m.getPDF(xMin, xMax, step)

	expSize := int(math.Round((xMax - xMin) / step))
	if len(result) != expSize {
		t.Errorf("Got solution vector of unexpected length: want %v, got %v.", expSize, len(result))
	}

	t.Logf("Result: %v, captured variance: %v", result, floats.Sum(result)/1)
	// t.Error("Dummy error to force output log") // uncomment to see output
}
