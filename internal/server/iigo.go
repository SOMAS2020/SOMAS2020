package server

import "github.com/SOMAS2020/SOMAS2020/internal/common"

// runIIGO : IIGO decides rule changes, elections, sanctions
func (s *SOMASServer) runIIGO() ([]common.Action, error) {
	s.logf("start runIIGO")
	defer s.logf("finish runIIGO")
	// TODO:- IIGO team
	return nil, nil
}
