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

	cpResources := shared.Resources(50)
	env := InitEnvironment(clientIDs, disasterConf)
	env.LastDisasterReport = DisasterReport{Magnitude: 1.0, X: env.Geography.XMax, Y: 0}
	t.Log(env.DisplayReport(cpResources, disasterConf))
	t.Error("dummy error")
}
