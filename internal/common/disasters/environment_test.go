package disasters

import (
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

func TestSamplingOfCertainties(t *testing.T) {
	clientIDs := []shared.ClientID{shared.Team1, shared.Team2} // arbitrarily chosen for test

	env := InitEnvironment(clientIDs)
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
