package server

import "github.com/SOMAS2020/SOMAS2020/internal/common/disasters"

// probeDisaster checks if a disaster occurs this turn
func (s *SOMASServer) probeDisaster() (disasters.Environment, error) {
	s.logf("start probeDisaster")
	defer s.logf("finish probeDisaster")

	e := s.gameState.Environment
	e = e.SampleForDisaster() // update env instance with sampled disaster info

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
	s.islandDeplete(effects.CommonPoolMitigated)                                             //island's resource will be depleted by disaster only when disaster happens and cp cannot fully mitigate
}
