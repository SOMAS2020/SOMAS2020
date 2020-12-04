package disasters

import (
	"math"

	"gonum.org/v1/gonum/stat/distuv"
)

// Island captures the location of a single island
type Island struct {
	name string
	x, y float64
}

// ArchipeligoGeography captures the collection of island geographies including bounding region of whole archipelago
type ArchipeligoGeography struct {
	islands          []Island
	xBounds, ybounds [2]float64
}

// IslandLocation is a convenience method to extract an island's location given its index
func (a ArchipeligoGeography) IslandLocation(index int) []float64 {
	island := a.islands[index]
	return []float64{island.x, island.y}
}

// DisasterParameters encapsulates a disaster's information - when and how it occurs. Disaster occurring is a Bernoulli random var with `p`=GlobalProb
// SpatialPDF determines the distribution type for the XY-location of the disaster peak. MagnitudeLambda is the lambda param in an exponential distr.
type DisasterParameters struct {
	GlobalProb      float64
	SpatialPDF      string
	MagnitudeLambda float64
}

// DisasterReport encapsulates a disaster location and magnitude. Note: magnitude of 0 => no disaster
type DisasterReport struct {
	magnitude, x, y float64
}

// Environment holds the state of the enivornment
type Environment struct {
	geography          ArchipeligoGeography
	disasterParams     DisasterParameters
	lastDisasterReport DisasterReport
}

// SampleForDisaster samples the stochastic disaster process to see if a disaster occurred
func (e *Environment) SampleForDisaster() DisasterReport {
	xBounds := e.geography.xBounds
	yBounds := e.geography.ybounds

	// spatial distr info
	pdfX := distuv.Uniform{Min: xBounds[0], Max: xBounds[1]}
	pdfY := distuv.Uniform{Min: yBounds[0], Max: yBounds[1]}

	pdfMag := distuv.Exponential{Rate: e.disasterParams.MagnitudeLambda} // Rate = lambda
	pdfGlobal := distuv.Bernoulli{P: e.disasterParams.GlobalProb}        // Bernoulli RV where `P` = P(X=1)

	dR := DisasterReport{0, -1, -1} // default: no disaster. Zero magnitude with arb co-ords

	if pdfGlobal.Rand() == 1.0 { // D Day
		dR = DisasterReport{pdfMag.Rand(), pdfX.Rand(), pdfY.Rand()}
	}
	e.lastDisasterReport = dR // record last report in env state
	return dR
}

// DisasterEffects returns the effects of the most recent DisasterReport held in the environment state
func (e Environment) DisasterEffects() map[string]float64 {
	out := map[string]float64{}                                  // TODO: change key type to ClientID
	epiX, epiY := e.lastDisasterReport.x, e.lastDisasterReport.x // epicentre of the disaster (peak mag)
	for _, island := range e.geography.islands {
		out[island.name] = e.lastDisasterReport.magnitude / (math.Sqrt(math.Pow(island.x-epiX, 2) + math.Pow(island.y-epiY, 2))) // effect on island i is inverse prop. to square of distance to epicentre
	}
	return out
}
