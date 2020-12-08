package server

import "github.com/SOMAS2020/SOMAS2020/internal/common/action"

// runIIGO : IIGO decides rule changes, elections, sanctions
func (s *SOMASServer) runIIGO() ([]action.Action, error) {
	s.logf("start runIITO")
	defer s.logf("finish runIITO")
	// TOOD:- IIGO team
	return nil, nil
}
