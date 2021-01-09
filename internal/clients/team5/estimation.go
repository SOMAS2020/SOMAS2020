package team5

import (
	"math"

	"github.com/aclements/go-moremath/stats"
	"github.com/google/go-cmp/cmp"
	"gonum.org/v1/gonum/floats"
)

type kdeModel struct {
	observations []float64
	weights      []float64
	estimator    stats.KDE
}

type modelStats struct {
	mean, variance float64 // TODO: decide which other pertinent stats to incorporate
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
func (m *kdeModel) getConstrainedPDF(xMin, xMax, step float64) (pdf, xrange []float64) {

	if cmp.Equal(m.estimator, stats.KDE{}) {
		m.fitModel()
	}
	xRange := makeRange(xMin, xMax, step)
	for _, x := range xRange {
		pdf = append(pdf, m.estimator.PDF(x))
	}
	return pdf, xRange
}

// get kernel-smoothed KDE PDF
func (m *kdeModel) getPDF(nSamples uint) (pdf, xrange []float64) {

	if nSamples == 0 {
		return pdf, xrange
	}

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
		for i := 0; i < int(nSamples); i++ {
			xrange = append(xrange, X[0])
			pdf = append(pdf, X[0])
		}
		floats.Scale(1/floats.Sum(pdf), pdf)
		return pdf, xrange
	}
	step := (xMax - xMin) / float64(nSamples) // 0 check above
	pdf, xrange = m.getConstrainedPDF(xMin, xMax, step)
	return pdf, xrange
}

func (m *kdeModel) getStatistics(nSamples uint) modelStats {

	mean, variance := 0.0, 0.0

	X := m.estimator.Sample.Xs

	if len(X) < 2 {

	} else if len(X) < 10 {
		// if we have too few samples, simply use observed sample statistics
		mean = stats.Mean(X)
		variance = stats.Variance(X)

	} else { // we have bare min required to use KDE estimated distribution now
		pdf, xrange := m.getPDF(nSamples) // TODO: simply use sample means for < 5 samples since KDE method is a bit unstable
		E := 0.0
		E2 := 0.0
		for i, p := range pdf {
			E += xrange[i] * p               // expected value
			E2 += math.Pow(xrange[i], 2) * p // variance
		}
		mean = E // use statistics (expectation) from our estimated distribution
		variance = E2 - math.Pow(E, 2)
	}

	return modelStats{
		mean:     mean,
		variance: variance,
	}
}

func createBasicKDE(observations []float64) kdeModel {
	m := kdeModel{observations: observations, weights: nil}
	m.fitModel()
	return m
}

func makeRange(min, max, step float64) []float64 {
	size := int(math.Round(((max - min) / step)))
	a := make([]float64, size)
	for i := range a {
		a[i] = min + step*float64(i)
	}
	return a
}
