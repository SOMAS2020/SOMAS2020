package server

import (
	"fmt"
)
// probeDisaster checks if a disaster occurs this turn
func (s *SOMASServer) probeDisaster() (bool, error) {
	s.logf("start probeDisaster")

	s.gameState.Environment.SampleForDisaster()

	disasterReport, leftover_damage := s.gameState.Environment.DisplayReport()
	fmt.Println(disasterReport)
	s.islandDistribute(s.gameState.Environment.CommonPool.Resource, disasterReport)
	s.islandDeplete(leftover_damage)


	defer s.logf("finish probeDisaster")
	// TODO:- env team
	return false, nil
}
