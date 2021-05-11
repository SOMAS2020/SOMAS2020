package server

import (
	"testing"

	"github.com/SOMAS2020/SOMAS2020/internal/common/baseclient"
	"github.com/SOMAS2020/SOMAS2020/internal/common/config"
	"github.com/SOMAS2020/SOMAS2020/internal/common/disasters"
	"github.com/SOMAS2020/SOMAS2020/internal/common/gamestate"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// note: individual effects of disaster on islands already tested in environment_test.go. This test is just to test
// that CP resources are deducted correctly.
func TestDisasterEffects(t *testing.T) {

	initialCP := shared.Resources(100)

	envConf := config.DisasterConfig{
		XMax:                        10,
		YMax:                        10,
		MagnitudeLambda:             1.0,
		CommonpoolThreshold:         0,
		MagnitudeResourceMultiplier: 100,
	}

	clientIDs := []shared.ClientID{shared.Teams["Team1"], shared.Teams["Team2"]}
	clientMap := map[shared.ClientID]baseclient.Client{}
	clientInfos := map[shared.ClientID]gamestate.ClientInfo{}
	for _, cid := range clientIDs {
		clientMap[cid] = baseclient.NewClient(cid)
		clientInfos[cid] = gamestate.ClientInfo{LifeStatus: shared.Alive, Resources: 100}
	}

	s := SOMASServer{
		gameState: gamestate.GameState{
			CommonPool:  initialCP,
			ClientInfos: clientInfos,
			Environment: disasters.InitEnvironment(clientIDs, envConf),
		},
		clientMap: clientMap,
	}
	s.gameConfig.DisasterConfig = envConf // update so that server methods relying on this get the same conf
	s.gameState.Environment.LastDisasterReport = disasters.DisasterReport{Magnitude: 0.5, X: 0, Y: 0}

	totalDamage := s.getTotalDamage(envConf)
	totalDamage /= 2 // since CPThreshold is zero, we need to halve totalDamage

	s.applyDisasterEffects() // actually apply effects now

	if s.gameState.CommonPool != (initialCP - shared.Resources(totalDamage)) {
		t.Errorf("Incorrect CP after disaster. Expected %v, got %v", initialCP-shared.Resources(totalDamage), s.gameState.CommonPool)
	}

	envConf.CommonpoolThreshold = 10000                                                             // now test CP damage with a very high CP threshold
	s.gameState.Environment.LastDisasterReport = disasters.DisasterReport{Magnitude: 1, X: 5, Y: 5} // test another location
	s.gameConfig.DisasterConfig = envConf
	s.gameState.CommonPool = 100 // reset CP
	totalDamage = s.getTotalDamage(envConf)

	s.applyDisasterEffects() // apply effects

	if s.gameState.CommonPool != (initialCP - shared.Resources(totalDamage)) {
		t.Errorf("Incorrect CP after disaster. Expected %v, got %v", initialCP-shared.Resources(totalDamage), s.gameState.CommonPool)
	}
}

func (s SOMASServer) getTotalDamage(conf config.DisasterConfig) float64 {
	effects := s.gameState.Environment.ComputeDisasterEffects(s.gameState.CommonPool, s.gameConfig.DisasterConfig)
	totalDamage := 0.0

	for _, effect := range effects.Absolute {
		totalDamage += effect * conf.MagnitudeResourceMultiplier
	}
	return totalDamage

}
