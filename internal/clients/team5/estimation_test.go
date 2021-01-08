package team5

import (
	"math"
	"testing"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/stat/distuv"
)

// this test only really tests a small helper function. The real KDE test was done by looking at output
// values and comparing to samples from the generating distribution. See "Estimation" in README.md in this package.
func TestKDE(t *testing.T) {

	for _, nSamples := range []int{10, 20, 50, 100, 1000} {
		// distNorm := distuv.Normal{Mu: 10, Sigma: 2, Src: rand.NewSource(1)}
		distExp := distuv.Exponential{Rate: 1}

		dist := distExp
		obs := make([]float64, nSamples)
		for i := 0; i < nSamples; i++ {
			obs[i] = dist.Rand()
		}
		m := kdeModel{observations: obs, weights: nil}

		xMin := 0.0
		xMax := 5.0
		step := 0.1
		result, _ := m.getPDF(uint(nSamples))

		expSize := int(math.Round((xMax - xMin) / step))
		if len(result) != expSize {
			t.Errorf("Got solution vector of unexpected length: want %v, got %v.", expSize, len(result))
		}
		t.Logf("Result (%v samples): %v, captured variance: %v", nSamples, result, step*floats.Sum(result)/1)
	}
	// t.Error("Dummy error to force output log") // uncomment to see output
}

func TestKDE2(t *testing.T) {

	obs := []float64{2, 2, 2, 2, 2}
	n := len(obs)
	m := kdeModel{observations: obs, weights: nil}

	result, _ := m.getPDF(uint(n))

	t.Logf("Result (%v samples): %v, exp: %v", n, result, m.getStatistics(uint(n)))
	t.Error("Dummy error to force output log") // uncomment to see output
}
