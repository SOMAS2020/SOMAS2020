package server

import (
	"math"

	"github.com/SOMAS2020/SOMAS2020/internal/common/disasters"
	"github.com/SOMAS2020/SOMAS2020/internal/common/shared"
)

// probeDisaster checks if a disaster occurs this turn
func (s *SOMASServer) probeDisaster() (disasters.Environment, error) {
	s.logf("start probeDisaster")
	defer s.logf("finish probeDisaster")

	e := s.gameState.Environment
	e = e.SampleForDisaster(s.gameConfig.DisasterConfig, s.gameState.Turn) // update env instance with sampled disaster info
	e.LastDisasterReport.Effects = e.ComputeDisasterEffects(s.gameState.CommonPool, s.gameConfig.DisasterConfig)

	disasterReport := e.DisplayReport(s.gameState.CommonPool, s.gameConfig.DisasterConfig) // displays disaster info and effects
	s.logf(disasterReport)

	return e, nil
}

// probeDisaster checks if a disaster occurs this turn
func (s *SOMASServer) applyDisasterEffects() {
	s.logf("start applyDisasterEffects")
	defer s.logf("finish applyDisasterEffects")

	e := s.gameState.Environment
	effects := e.ComputeDisasterEffects(s.gameState.CommonPool, s.gameConfig.DisasterConfig) // get disaster effects - absolute, proportional and CP-mitigated
	totalResourceImpact := disasters.GetDisasterResourceImpact(s.gameState.CommonPool, effects, s.gameConfig.DisasterConfig)
	s.islandDeplete(effects.CommonPoolMitigated)
	s.logf("*** impact: %v, CP: %v, conf: %+v", totalResourceImpact, s.gameState.CommonPool, s.gameConfig.DisasterConfig) //island's resource will be depleted by disaster only when disaster happens and cp cannot fully mitigate
	s.gameState.CommonPool = shared.Resources(math.Max(float64(s.gameState.CommonPool)-float64(totalResourceImpact), 0))  // deduct disaster damage from CP
}
