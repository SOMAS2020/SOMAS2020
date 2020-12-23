package server

import "github.com/SOMAS2020/SOMAS2020/internal/common/disasters"

// probeDisaster checks if a disaster occurs this turn
func (s *SOMASServer) probeDisaster() (disasters.Environment, error) {
	s.logf("start probeDisaster")
	defer s.logf("finish probeDisaster")

	e := s.gameState.Environment.SampleForDisaster()
	disasterReport, leftover_damage := s.gameState.Environment.DisplayReport()
	s.logf(disasterReport)
	s.islandDeplete(leftover_damage)				//island will be further depleted by disaster only when disaster happens and cp does not have enough resource

	return e, nil
}
