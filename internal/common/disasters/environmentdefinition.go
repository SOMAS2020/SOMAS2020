package disasters

import (
	"math"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"gonum.org/v1/gonum/stat/distuv"
)

// IslandLocationInfo captures the location of a single island
type IslandLocationInfo struct {
	id   shared.ClientID
	x, y shared.Coordinate
}

// ArchipelagoGeography captures the collection of island geographies including bounding region of whole archipelago
type ArchipelagoGeography struct {
	islands                map[shared.ClientID]IslandLocationInfo
	xMin, xMax, yMin, yMax shared.Coordinate
}

// disasterParameters encapsulates a disaster's information - when and how it occurs. Disaster occurring is a Bernoulli random var with `p`=GlobalProb
// SpatialPDF determines the distribution type for the XY-location of the disaster peak. MagnitudeLambda is the lambda param in an exponential distr.
type disasterParameters struct {
	globalProb      float64
	spatialPDF      string
	magnitudeLambda float64
}

// DisasterReport encapsulates a disaster location and magnitude. Note: magnitude of 0 => no disaster
type DisasterReport struct {
	Magnitude, X, Y float64
}

// Environment holds the state of the enivornment
type Environment struct {
	geography          ArchipelagoGeography
	disasterParams     disasterParameters
	lastDisasterReport DisasterReport
}

// SampleForDisaster samples the stochastic disaster process to see if a disaster occurred
func (e Environment) SampleForDisaster() DisasterReport {
	// spatial distr info
	pdfX := distuv.Uniform{Min: e.geography.xMin, Max: e.geography.xMax}
	pdfY := distuv.Uniform{Min: e.geography.yMin, Max: e.geography.yMax}

	pdfMag := distuv.Exponential{Rate: e.disasterParams.magnitudeLambda} // Rate = lambda
	pdfGlobal := distuv.Bernoulli{P: e.disasterParams.globalProb}        // Bernoulli RV where `P` = P(X=1)

	dR := DisasterReport{0, -1, -1} // default: no disaster. Zero magnitude with arb co-ords

	if pdfGlobal.Rand() == 1.0 { // D Day
		dR = DisasterReport{Magnitude: pdfMag.Rand(), X: pdfX.Rand(), Y: pdfY.Rand()}
	}
	e.lastDisasterReport = dR // record last report in env state // TODO: change - won't currently do anything
	return dR
}

// DisasterEffects returns the effects of the most recent DisasterReport held in the environment state
func (e Environment) DisasterEffects() map[shared.ClientID]shared.Magnitude {
	out := map[shared.ClientID]shared.Magnitude{}
	epiX, epiY := e.lastDisasterReport.X, e.lastDisasterReport.Y // epicentre of the disaster (peak mag)
	for _, island := range e.geography.islands {
		out[island.id] = e.lastDisasterReport.Magnitude / (math.Sqrt(math.Pow(island.x-epiX, 2) + math.Pow(island.y-epiY, 2))) // effect on island i is inverse prop. to square of distance to epicentre
	}
	return out
}
