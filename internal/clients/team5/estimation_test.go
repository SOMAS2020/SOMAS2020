package team5

import (
	"sort"
	"testing"

	"github.com/aclements/go-moremath/stats"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/stat"
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

		// xMin := 0.0
		// xMax := 5.0
		step := (floats.Max(obs) - floats.Min(obs)) / float64(nSamples)
		result, _ := m.getPDF(uint(nSamples))

		if len(result) != nSamples {
			t.Errorf("Got solution vector of unexpected length: want %v, got %v.", nSamples, len(result))
		}
		t.Logf("Result (%v samples): %v, captured variance: %v", nSamples, result, step*floats.Sum(result)/1)
	}
	// t.Error("Dummy error to force output log") // uncomment to see output
}

func TestConfidence(t *testing.T) {
	n := 100
	pdf := distuv.Exponential{Rate: 1}
	samples := make([]float64, n)
	for i := range samples {
		samples[i] = pdf.Rand()
	}

	quant := func(q float64, x []float64) float64 { return stat.Quantile(q, stat.LinInterp, x, nil) }

	sort.Float64s(samples)
	conf := 1 - (stats.StdDev(samples))/(quant(0.97, samples)-quant(0.03, samples))
	t.Logf("conf: %+v", conf)
	// t.Error("some error :)")

}

func TestStatistics(t *testing.T) {

	tests := []struct {
		name                          string
		obs                           []float64
		expectedMean, expectedVar     float64
		strictEquality, errorExpected bool // whether or not exp. mean and variance should strictly equal result and whether or not we expect an error
	}{
		{"empty obs set", []float64{}, 0, 0, false, true},
		{"single obs", []float64{2}, 2, 0, true, false},
		{"multiple obs periodic, sample only", []float64{3, 3, 3}, 3, 0, true, false},
		{"multiple obs periodic, KDE", []float64{3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3}, 3, 0, true, false},
		{"multiple obs not periodic, KDE", []float64{10, 3, 2, 3, 20, 3, 8, 3, 5, 3, 3, 2, 3}, 3, 0, false, false},
	}
	for _, tc := range tests {
		testname := tc.name
		t.Run(testname, func(t *testing.T) {

			model := newKdeModel(tc.obs)
			stats, err := model.getStatistics(uint(len(tc.obs)))

			if err != nil && !tc.errorExpected {
				t.Errorf("received unexpected error: %v with input: %v", err, tc.obs)
			}

			if stats.mean != tc.expectedMean && tc.strictEquality {
				t.Errorf("incorrect mean: got %.1f, want %.1f", stats.mean, tc.expectedMean)
			}
			if stats.variance != tc.expectedVar && tc.strictEquality {
				t.Errorf("incorrect variance: got %.1f, want %.1f", stats.variance, tc.expectedVar)
			}

		})
	}
}
