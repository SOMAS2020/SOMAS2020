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
}
