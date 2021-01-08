package team5

import (
	"math"

	"github.com/aclements/go-moremath/stats"
	"github.com/google/go-cmp/cmp"
)

type kdeModel struct {
	observations []float64
	bandwidth    float64
	weights      []float64
	estimator    stats.KDE
}

func (m *kdeModel) updateModel(newSamples []float64) {
	for _, s := range newSamples {
		m.observations = append(m.observations, s)
	}
	m.fitModel()
}

func (m *kdeModel) fitModel() {
	s := stats.Sample{Xs: m.observations, Weights: m.weights, Sorted: false}
	m.estimator = stats.KDE{
		Sample:         s,
		Kernel:         stats.GaussianKernel,
		Bandwidth:      m.bandwidth,
		BoundaryMethod: stats.BoundaryReflect,
		BoundaryMin:    0,
		BoundaryMax:    math.Inf(1),
	}
}

func (m *kdeModel) getPDF(xMin, xMax float64) (pdf []float64) {

	if cmp.Equal(m.estimator, stats.KDE{}) {
		m.fitModel()
	}
	for _, x := range makeRange(xMin, xMax) {
		pdf = append(pdf, m.estimator.PDF(x))
	}
	return pdf
}

func createBasicKDE(observations []float64, bandwidth float64) kdeModel {
	m := kdeModel{observations: observations, bandwidth: bandwidth, weights: nil}
	m.fitModel()
	return m
}

func makeRange(min, max float64) []float64 {
	a := make([]float64, int(max-min+1))
	for i := range a {
		a[i] = min + float64(i)
	}
	return a
}
