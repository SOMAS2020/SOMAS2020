package team5

import (
	"math"
	"sort"

	"github.com/aclements/go-moremath/stats"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/stat"
)

type kdeModel struct {
	observations []float64
	weights      []float64
	estimator    stats.KDE
}

// convenience method to create a new model in the way that will most commonly be required - not
// messing with bandwidth or weight. Also calls m.fitModel for you
func newKdeModel(observations []float64) kdeModel {
	m := kdeModel{observations: observations, weights: nil}
	m.fitModel()
	return m
}

type modelStats struct {
	mean, variance, meanConfidence float64 // TODO: decide which other pertinent stats to incorporate
}

func (m *kdeModel) updateModel(newSamples []float64) {
	for _, s := range newSamples {
		m.observations = append(m.observations, s)
	}
	m.fitModel()
}

// allows you to modify bandwidth if you think you're better than Scott estimator
func (m *kdeModel) setBandwidth(bw float64) {
	m.estimator.Bandwidth = bw
}

// allows you to modify bandwidth if you think you're better than Scott estimator
func (m *kdeModel) setWeights(w []float64) {
	m.weights = w
}

func (m *kdeModel) fitModel() {
	s := stats.Sample{Xs: m.observations, Weights: m.weights, Sorted: false}
	m.estimator = stats.KDE{
		Sample:         s,
		Kernel:         stats.GaussianKernel,
		Bandwidth:      0, // if zero, uses Scott BW estimator based on data
		BoundaryMethod: stats.BoundaryReflect,
		BoundaryMin:    0,
		BoundaryMax:    math.Inf(1),
	}
}

// get kernel-smoothed KDE PDF over specific range
func (m *kdeModel) getConstrainedPDF(xMin, xMax, step float64) (pdf, xrange []float64, err error) {

	if xMin == xMax {
		return pdf, xrange, errors.Errorf("xMin and xMax cannot be the same!")
	}
	if step <= 0 {
		return pdf, xrange, errors.Errorf("Step must be positive non-zero!")
	}

	if cmp.Equal(m.estimator, stats.KDE{}) {
		m.fitModel()
	}
	xrange = makeRange(xMin, xMax, step)
	for _, x := range xrange {
		pdf = append(pdf, m.estimator.PDF(x))
	}
	return pdf, xrange, nil
}

// get kernel-smoothed KDE PDF
func (m *kdeModel) getPDF(nSamples uint) (pdf, xrange []float64) {

	// update model if it has not been fitted yet
	if cmp.Equal(m.estimator, stats.KDE{}) {
		m.fitModel()
	}

	X := m.estimator.Sample.Xs

	if len(X) == 0 || nSamples == 0 {
		return pdf, xrange
	}

	xMax := floats.Max(X)
	xMin := floats.Min(X)

	if xMax == xMin || stats.Variance(X) == 0 { // special case where bounds are same or zero variance (all same values)
		return []float64{1.0}, []float64{xMax} // trivial PDF
	}
	step := (xMax - xMin) / float64(nSamples)              // 0 check above
	pdf, xrange, _ = m.getConstrainedPDF(xMin, xMax, step) // don't
	return pdf, xrange
}

func (m *kdeModel) getStatistics(nSamples uint) (modelStats, error) {

	mean, variance := 0.0, 0.0

	X := m.estimator.Sample.Xs

	if len(X) == 0 {
		return modelStats{}, errors.Errorf("Sample length 0 - cannot calculate statistics.") // zeros
	} else if len(X) < 10 {
		// if we have too few samples, simply use observed sample statistics
		mean = stats.Mean(X)
		variance = stats.Variance(X)

	} else { // we have bare min required to use KDE estimated distribution now
		pdf, xrange := m.getPDF(nSamples) // TODO: simply use sample means for < 5 samples since KDE method is a bit unstable
		E, E2 := 0.0, 0.0
		for i, p := range pdf {
			E += xrange[i] * p               // expected value
			E2 += math.Pow(xrange[i], 2) * p // variance
		}
		mean = E // use statistics (expectation) from our estimated distribution
		variance = E2 - math.Pow(E, 2)
	}

	return modelStats{
		mean:           mean,
		variance:       variance,
		meanConfidence: naiveConfidence(X),
	}, nil
}

// returns the confidence of using the mean value as a the predictor (based on dispersion)
func naiveConfidence(x []float64) float64 {
	if len(x) < 2 {
		return 0
	}

	quant := func(q float64, x []float64) float64 { return stat.Quantile(q, stat.LinInterp, x, nil) }
	sort.Float64s(x)
	qLo := quant(0.03, x) // 3% percentile as robust approx of minimum
	qHi := quant(0.97, x) // 97 percentile as robust approx of max (avoid outliers)
	if qHi == qLo {
		return 0.99 // can be pretty certain if max \approx min
	}
	conf := 1 - math.Min((stats.StdDev(x))/(qHi-qLo), 1) // compare std dev to range
	return conf
}

func makeRange(min, max, step float64) []float64 {
	size := int(math.Round(((max - min) / step)))
	a := make([]float64, size)
	for i := range a {
		a[i] = min + step*float64(i)
	}
	return a
}
