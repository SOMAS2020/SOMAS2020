package team5

import (
	"testing"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/stat/distuv"
)

func TestKDE(t *testing.T) {

	nSamples := 1000
	dist := distuv.Normal{Mu: 10, Sigma: 2, Src: rand.NewSource(1)}
	obs := make([]float64, nSamples)
	for i := 0; i < nSamples; i++ {
		obs[i] = dist.Rand()
	}
	m := kdeModel{observations: obs, bandwidth: 1.0, weights: nil}

	result := m.getPDF(0, 20)

	t.Logf("Result: %v, sum: %v", result, floats.Sum(result))

	t.Error("some error")
}
