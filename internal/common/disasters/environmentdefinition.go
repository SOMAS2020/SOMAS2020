package disasters

import (
	"math"

	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
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
	Effects   DisasterEffects
}

// DisasterEffects encapsulates various types of effects on each island after a disaster. These include
// the absolute magnitude felt, the proportional mag. relative to that of other islands and the proportional
// magnitude felt after applying common pool mitigation
type DisasterEffects struct {
	Absolute, Proportional, CommonPoolMitigated map[shared.ClientID]shared.Magnitude
}

// Environment holds the state of the enivornment
type Environment struct {
	Geography          ArchipelagoGeography
	LastDisasterReport DisasterReport
}

// SampleForDisaster samples the stochastic disaster process to see if a disaster occurred
func (e Environment) SampleForDisaster(dConf config.DisasterConfig, turn uint) Environment {
	// spatial distr info
	pdfX := distuv.Uniform{Min: e.Geography.XMin, Max: e.Geography.XMax}
	pdfY := distuv.Uniform{Min: e.Geography.YMin, Max: e.Geography.YMax}

	pdfMag := distuv.Exponential{Rate: dConf.MagnitudeLambda} // Rate = lambda

	dR := DisasterReport{Magnitude: 0, X: -1, Y: -1} // default: no disaster. Zero magnitude with arb co-ords

	if dConf.StochasticPeriod {
		// if T is the disaster period (time between occurrences), we need:
		// E[T] = T (stochastic and deterministic cases respectively). Since
		// T is a geometric RV in the stochastic case, p = 1/E[T]
		p := 1 / float64(dConf.Period)
		pdfGlobal := distuv.Bernoulli{P: p} // Bernoulli RV where `P` = P(X=1)

		if pdfGlobal.Rand() == 1.0 { // D Day
			dR = DisasterReport{Magnitude: pdfMag.Rand(), X: pdfX.Rand(), Y: pdfY.Rand()}
		}
	} else {
		if turn%uint(dConf.Period) == 0 && turn > 0 {
			dR = DisasterReport{Magnitude: pdfMag.Rand(), X: pdfX.Rand(), Y: pdfY.Rand()}
		}
	}

	e.LastDisasterReport = dR // record last report in env state
	return e                  // return same env back but with updated disaster report
}

func (e Environment) computeUnmitigatedDisasterEffects() DisasterEffects {
	individualEffect := map[shared.ClientID]shared.Magnitude{}
	proportionalEffect := map[shared.ClientID]shared.Magnitude{}
	totalEffect := 0.0

	epiX, epiY := e.LastDisasterReport.X, e.LastDisasterReport.Y // epicentre of the disaster (peak mag)
	for _, island := range e.Geography.Islands {
		effect := e.LastDisasterReport.Magnitude / math.Hypot(island.X-epiX, island.Y-epiY) // effect on island i is inverse prop. to square of distance to epicentre
		individualEffect[island.ID] = math.Min(effect, e.LastDisasterReport.Magnitude)      // to prevent divide by zero -> inf
		totalEffect = totalEffect + individualEffect[island.ID]
	}
	if totalEffect == 0 {
		for _, island := range e.Geography.Islands {
			proportionalEffect[island.ID] = individualEffect[island.ID]
		}
	} else {
		for _, island := range e.Geography.Islands {
			proportionalEffect[island.ID] = individualEffect[island.ID] / totalEffect
		}
	}
	return DisasterEffects{Absolute: individualEffect, Proportional: proportionalEffect} // ommit CP mitigated effect here - not relevant
}

// ComputeDisasterEffects returns the individual (absolute) effects and proportional effect (compared to total damage on each island)
// This method uses the latest disaster report stored in environment
func (e Environment) ComputeDisasterEffects(cpResources shared.Resources, dConf config.DisasterConfig) DisasterEffects {

	unmitigatedEffects := e.computeUnmitigatedDisasterEffects()
	mitigatedEffects := e.MitigateDisaster(cpResources, unmitigatedEffects, dConf)

	return DisasterEffects{Absolute: unmitigatedEffects.Absolute, Proportional: unmitigatedEffects.Proportional, CommonPoolMitigated: mitigatedEffects}
}
