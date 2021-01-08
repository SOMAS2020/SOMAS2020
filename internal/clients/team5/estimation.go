package team5

import (
	"math"

	"github.com/aclements/go-moremath/stats"
	"github.com/google/go-cmp/cmp"
)

type kdeModel struct {
	observations []float64
	weights      []float64
	estimator    stats.KDE
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

func (m *kdeModel) getPDF(xMin, xMax, step float64) (pdf []float64) {

	if cmp.Equal(m.estimator, stats.KDE{}) {
		m.fitModel()
	}
	for _, x := range makeRange(xMin, xMax, step) {
		pdf = append(pdf, m.estimator.PDF(x))
	}
	return pdf
}

func createBasicKDE(observations []float64, bandwidth float64) kdeModel {
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
