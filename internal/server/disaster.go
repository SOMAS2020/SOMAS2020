package server

import "github.com/SOMAS2020/SOMAS2020/internal/common/disasters"

// probeDisaster checks if a disaster occurs this turn
func (s *SOMASServer) probeDisaster() (disasters.DisasterReport, error) {
	s.logf("start probeDisaster")
	defer s.logf("finish probeDisaster")

	disasterReport := s.gameState.Environment.SampleForDisaster()
	return disasterReport, nil
}
