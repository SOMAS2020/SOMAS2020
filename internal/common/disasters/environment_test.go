package disasters

import (
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func TestSamplingOfCertainties(t *testing.T) {
	clientIDs := []shared.ClientID{shared.Team1, shared.Team2} // arbitrarily chosen for test

	disasterConf := config.DisasterConfig{
		XMin:            0.0,
		XMax:            10.0, // chosen quite arbitrarily
		YMin:            0.0,
		YMax:            10.0,
		GlobalProb:      0.1,
		SpatialPDFType:  shared.Uniform,
		MagnitudeLambda: 1.0,
	}

	env := InitEnvironment(clientIDs, disasterConf)
	env.DisasterParams.globalProb = 0
	updatedEnv := env.SampleForDisaster()
	if updatedEnv.LastDisasterReport.Magnitude > 0.0 {
		t.Error("Disaster struck despite global probability set to zero")
	}

	env.DisasterParams.globalProb = 1
	updatedEnv = env.SampleForDisaster()
	if updatedEnv.LastDisasterReport.Magnitude == 0.0 {
		t.Error("No disaster recorded despite global prob. set to one")
	}
	x := updatedEnv.LastDisasterReport.X
	y := updatedEnv.LastDisasterReport.Y

	if x < disasterConf.XMin || x > disasterConf.XMax {
		t.Error("Disaster location outside of config x bounds")
	}

	if y < disasterConf.YMin || y > disasterConf.YMax {
		t.Error("Disaster location outside of config y bounds")
	}
}

func TestDisasterEffects(t *testing.T) {

	clientIDs := []shared.ClientID{shared.Team1, shared.Team2, shared.Team3} // arbitrarily chosen for test

	disasterConf := config.DisasterConfig{
		XMin:                        0.0,
		XMax:                        10.0,
		YMin:                        0.0,
		YMax:                        10.0,
		GlobalProb:                  1.0,
		SpatialPDFType:              shared.Uniform,
		MagnitudeLambda:             1.0,
		CommonpoolThreshold:         10.0,
		MagnitudeResourceMultiplier: 100,
	}

	env := InitEnvironment(clientIDs, disasterConf)
	env.LastDisasterReport = DisasterReport{Magnitude: 1.0, X: env.Geography.XMax, Y: 0} // right on team 3

	// test complete mitigation
	cpResources := shared.Resources(500)
	effects := env.ComputeDisasterEffects(cpResources, disasterConf)
	if !zeroEffects(effects.CommonPoolMitigated) {
		t.Error("Expected full mitigation by common pool but got non-zero disaster effects")
	}

	// test partial mitigation
	cpResources = shared.Resources(50)
	effects = env.ComputeDisasterEffects(cpResources, disasterConf)
	if zeroEffects(effects.CommonPoolMitigated) {
		t.Error("Unexpected full disaster mitigation despite deficient common pool")
	}
	if !correctlyMitigatedEffects(effects, disasterConf.MagnitudeResourceMultiplier) {
		t.Log(effects.Absolute)
		t.Error("Expected mitigated effects to be less than or equal to original effects for each island")
	}

	// test differential impact
	eA := effects.Absolute
	if !(eA[shared.Team1] < eA[shared.Team2] && eA[shared.Team2] < eA[shared.Team3]) {
		t.Error("Expected (descending) order of abs effects to be Team3, Team2, Team1")
	}
	eA = effects.CommonPoolMitigated
	if !(eA[shared.Team1] < eA[shared.Team2] && eA[shared.Team2] < eA[shared.Team3]) {
		t.Error("Expected (descending) order of CP-mitigated effects to be Team3, Team2, Team1")
	}

	t.Log(env.DisplayReport(cpResources, disasterConf)) // in case of an error
}

// check if effect for every island is zero
func zeroEffects(effects map[shared.ClientID]shared.Magnitude) bool {
	allZero := true
	for _, mag := range effects {
		allZero = allZero && (mag == 0)
	}
	return allZero
}

// check that mitigated effects are indeed less than original effects
func correctlyMitigatedEffects(de DisasterEffects, magResourceMult float64) bool {
	correct := true
	for id, mag := range de.CommonPoolMitigated {
		correct = correct && (mag <= de.Absolute[id]*magResourceMult)
	}
	return correct
}
