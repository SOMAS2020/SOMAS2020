package disasters

import (
	"math"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
	"gonum.org/v1/gonum/stat/distuv"
)

// IslandLocationInfo captures the location of a single island
type IslandLocationInfo struct {
	ID   shared.ClientID
	X, Y shared.Coordinate
}

// ArchipelagoGeography captures the collection of island geographies including bounding region of whole archipelago
type ArchipelagoGeography struct {
	Islands                map[shared.ClientID]IslandLocationInfo
	XMin, XMax, YMin, YMax shared.Coordinate
}

// disasterParameters encapsulates a disaster's information - when and how it occurs. Disaster occurring is a Bernoulli random var with `p`=GlobalProb
// SpatialPDF determines the distribution type for the XY-location of the disaster peak. MagnitudeLambda is the lambda param in an exponential distr.
type disasterParameters struct {
	globalProb      float64
	spatialPDF      shared.SpatialPDFType
	magnitudeLambda float64
}

// DisasterReport encapsulates a disaster location and magnitude. Note: magnitude of 0 => no disaster
type DisasterReport struct {
	Magnitude shared.Magnitude
	X, Y      shared.Coordinate
}

// Environment holds the state of the enivornment
type Environment struct {
	Geography          ArchipelagoGeography
	DisasterParams     disasterParameters
	LastDisasterReport DisasterReport
	CommonPool		   Commonpool
}

// SampleForDisaster samples the stochastic disaster process to see if a disaster occurred
func (e Environment) SampleForDisaster() Environment {
	// spatial distr info
	pdfX := distuv.Uniform{Min: e.Geography.XMin, Max: e.Geography.XMax}
	pdfY := distuv.Uniform{Min: e.Geography.YMin, Max: e.Geography.YMax}

	pdfMag := distuv.Exponential{Rate: e.DisasterParams.magnitudeLambda} // Rate = lambda
	pdfGlobal := distuv.Bernoulli{P: e.DisasterParams.globalProb}        // Bernoulli RV where `P` = P(X=1)

	dR := DisasterReport{Magnitude: 0, X: -1, Y: -1} // default: no disaster. Zero magnitude with arb co-ords

	if pdfGlobal.Rand() == 1.0 { // D Day
		dR = DisasterReport{Magnitude: pdfMag.Rand(), X: pdfX.Rand(), Y: pdfY.Rand()}
	}
	e.LastDisasterReport = dR // record last report in env state
	return e                  // return same env back but with updated disaster report
}

// DisasterEffects returns the individual effects and proportional effect (compared to total damage on 6 island) 
// 		of the most recent DisasterReport held in the environment state
func (e Environment) DisasterEffects() (map[shared.ClientID]shared.Magnitude, map[shared.ClientID]shared.Magnitude) {
	individualEffect := map[shared.ClientID]shared.Magnitude{}
	proportionalEffect := map[shared.ClientID]shared.Magnitude{}
	totalEffect := 0.0 

	epiX, epiY := e.LastDisasterReport.X, e.LastDisasterReport.Y // epicentre of the disaster (peak mag)
	for _, island := range e.Geography.Islands {
		effect := e.LastDisasterReport.Magnitude / math.Hypot(island.X-epiX, island.Y-epiY) // effect on island i is inverse prop. to square of distance to epicentre
		individualEffect[island.ID] = math.Min(effect, e.LastDisasterReport.Magnitude)                   // to prevent divide by zero -> inf
		totalEffect = totalEffect + individualEffect[island.ID]
	}

	for _, island := range e.Geography.Islands {
		proportionalEffect[island.ID] =  individualEffect[island.ID] / totalEffect
	}

	return individualEffect, proportionalEffect
}
